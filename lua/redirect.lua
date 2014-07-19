local resolver = require "resty.dns.resolver"
local r, err = resolver:new { nameservers = { "8.8.8.8", "8.8.4.4" } }

if not r then
    ngx.status = ngx.HTTP_INTERNAL_SERVER_ERROR
    return ngx.say("failed to instantiate the resolver:", err)
end

-- local start = os.clock()
local answers, err = r:query(ngx.var.host, { qtype = resolver.TYPE_TXT })
-- ngx.header["X-Query-Time"] = os.clock() - start

if not answers then
    ngx.status = ngx.HTTP_INTERNAL_SERVER_ERROR
    return ngx.say("failed to query the DNS server:", err)
end


local codes = {
    ["301"] = 301,
    ["302"] = 302,
    -- ["307"] = 307, -- not supported by openresty
    default = 302
}

-- escape patterns
function escape (s)
    return string.gsub(s, "[%^%$%(%)%%%.%[%]%*%+%-%?]", "%%%0")
end


if not answers.errcode then
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
end

-- fall back to redirect.name for documentation
return ngx.redirect("http://redirect.name/")
