function parse (record)
    -- TODO: add ^ and $ once resolver is fixed
    -- https://github.com/openresty/lua-resty-dns/issues/4
    local value = string.match(record, "Redirects?(%s+.*)")
    if not value then return {} end

    local args = {}

    -- parse target
    local target = string.match(value, "%s+to%s+(%S+)")
    if target and (string.match(target, "^/") or string.match(target, "^https?://")) then
        args.target = target
    else
        return {}
    end

    -- parse path
    local path = string.match(value, "%s+from%s+(%S+)")
    if path then
        args.path = path
    end

    -- parse status
    local status = string.match(value, "%s+(permanently)") or string.match(value, "%s+(temporarily)")
    if status then
        args.status = status
    end

    return args
end

return parse
