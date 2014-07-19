-- supported redirect status codes
local codes = {
    ["301"] = 301,
    ["302"] = 302,
    -- ["307"] = 307, -- not supported by openresty
    default = 302
}

-- escape patterns
local function escape (s)
    return string.gsub(s, "[%^%$%(%)%%%.%[%]%*%+%-%?]", "%%%0")
end

local function parse (uri, config)
    local m, err = ngx.re.match(config, "redirect\\.name=(?:(?<path>/\\S*)\\s+)?(?<target>https?://\\S+|/\\S*)(?:\\s+(?<status>301|302))?", "i")
    if not m then return end

    local target = m.target
    local status = codes[m.status] or codes.default
    -- no `path` assumes catch-all, so redirect immediately to `target`
    if not m.path then return target, status end

    local path, count = string.gsub(escape(m.path), "%%%*", "(.*)", 1)
    path = "^" .. path .. "$"

    -- if we can't find the pattern, break to continue to next record
    if not string.find(uri, path) then return end

    -- wildcard replacement of `uri` if there's a wildcard in our `path`
    if count then
        target = escape(target)
        target = string.gsub(target, "%%%*", "%%1", 1)
        target = string.gsub(uri, from, target)
    end

    return target, status
end

return parse
