local resolver = require "resty.dns.resolver"

local codes = {
    ["301"] = 301,
    ["302"] = 302,
    -- ["307"] = 307, -- not supported by openresty
    default = 302
}

function fail (...)
    ngx.status = ngx.HTTP_INTERNAL_SERVER_ERROR
    ngx.say("<h1>Error</h1>")
    ngx.say(...)
end

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

-- fall back to redirect when query returns bad hostname, non-existent hostname, or other errors
if answers.errcode then return fallback() end


for i, ans in ipairs(answers) do
    -- wrap in `repeat` so we can use `break` to mimic Lua's missing `continue`
    repeat
        if ans.type ~= resolver.TYPE_TXT then break end

        local m, err = ngx.re.match(ans.txt, "redirect\\.name=(?:(?<from>/\\S*)\\s+)?(?<to>(?:https?://\\S+|/\\S*))(?:\\s+(?<status>\\d+))?", "i")
        if not m then break end

        local uri = ngx.var.request_uri
        -- ngx.header["X-URI"] = uri
        -- ngx.header["X-From"] = m.from
        local from = escape(m.from or "/*")
        -- ngx.header["X-From-Escaped"] = from
        local pattern = "^" .. string.gsub(from, "%%%*", "(.*)", 1) .. "$"
        -- ngx.header["X-Pattern"] = pattern

        if not string.find(uri, pattern) then break end

        -- ngx.header["X-To"] = m.to
        local to = escape(m.to)
        -- ngx.header["X-To-Escaped"] = to
        local replace = string.gsub(to, "%%%*", "%%1", 1)
        -- ngx.header["X-Replace"] = replace
        local target = string.gsub(uri, pattern, replace)
        -- ngx.header["X-Target"] = target

        local status = codes[m.status] or codes.default

        return ngx.redirect(target, status)
    until true
end


-- fall back to redirect.name for documentation
return fallback()
