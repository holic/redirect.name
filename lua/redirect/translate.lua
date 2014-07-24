-- supported redirect status codes
local codes = {
    ["301"] = 301,
    ["302"] = 302,
    ["permanently"] = 301,
    ["temporarily"] = 302,
    default = 302
}

-- escape pattern characters
local function escape (s)
    return string.gsub(s, "[%^%$%(%)%%%.%[%]%*%+%-%?]", "%%%0")
end

local function translate (uri, args)
    if not uri then return end
    if not args then return end

    local target = args.target
    if not target then return end

    local status = codes[args.status] or codes.default
    -- no `path` assumes catch-all, so redirect immediately to `target`
    if not args.path then return target, status end

    local path, count = string.gsub(escape(args.path), "%%%*", "(.*)", 1)
    path = "^" .. path .. "$"

    -- if we can't find the pattern, break to continue to next record
    if not string.find(uri, path) then return end

    -- wildcard replacement of `uri` if there's a wildcard in our `path`
    if count then
        target = string.gsub(target, "%*", "%%1", 1)
        target = string.gsub(uri, path, target, 1)
    end

    return target, status
end

return translate
