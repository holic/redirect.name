local resolver = require "resty.dns.resolver"

-- supported redirect status codes
local codes = {
    ["301"] = 301,
    ["302"] = 302,
    -- ["307"] = 307, -- not supported by openresty
    default = 302
}

-- hard failures
function fail (...)
    ngx.status = ngx.HTTP_INTERNAL_SERVER_ERROR
    ngx.say("<h1>Error</h1>")
    ngx.say(...)
end

-- fall back to documentation site
function fallback ()
    ngx.redirect("http://redirect.name/")
end

-- escape patterns
function escape (s)
    return string.gsub(s, "[%^%$%(%)%%%.%[%]%*%+%-%?]", "%%%0")
end


local r, err = resolver:new { nameservers = { "8.8.8.8", "8.8.4.4" } }
if not r then return fail("Failed to instantiate the resolver: ", err) end

local answers, err = r:query(ngx.var.host, { qtype = resolver.TYPE_TXT })
if not answers then return fail("Failed to query the DNS server: ", err) end

-- fall back when query returns bad hostname, non-existent hostname, or other errors
if answers.errcode then return fallback() end


for i, ans in ipairs(answers) do
    -- wrap in `repeat` so we can use `break` to mimic Lua's missing `continue`
    repeat
        if ans.type ~= resolver.TYPE_TXT then break end

        local m, err = ngx.re.match(ans.txt, "redirect\\.name=(?:(?<from>/\\S*)\\s+)?(?<to>(?:https?://\\S+|/\\S*))(?:\\s+(?<status>301|302))?", "i")
        if not m then break end

        local to = m.to
        local status = codes[m.status] or codes.default
        -- no `from` assumes catch-all, so redirect immediately to `to`
        if not m.from then return ngx.redirect(to, status) end

        local uri = ngx.var.request_uri
        local from, count = string.gsub(escape(m.from), "%%%*", "(.*)", 1)
        from = "^" .. from .. "$"

        -- if we can't find the pattern, break to continue to next record
        if not string.find(uri, from) then break end

        -- wildcard replacement of `uri` if there's a wildcard in our `from` path
        if count then
            to = escape(to)
            to = string.gsub(to, "%%%*", "%%1", 1)
            to = string.gsub(uri, from, to)
        end

        return ngx.redirect(to, status)
    until true
end


-- fall back if no records match
return fallback()
