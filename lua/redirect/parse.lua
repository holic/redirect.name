function parse (record)
    -- TODO: add ^ and $ once resolver is fixed
    -- https://github.com/openresty/lua-resty-dns/issues/4
    local value = string.match(record, "redirect%.name=(.*)")
    if not value then return {} end

    -- split record value into a table of args
    local args = {}
    local i = 0
    for arg in string.gmatch(value, "%S+") do
        table.insert(args, arg)
        i = i + 1
        -- break after first three args
        if i >= 3 then break end
    end
    if not #args then return {} end

    -- parse status code
    if #args > 1 and string.match(args[#args], "^30[12]$") then
        args.status = args[#args]
        table.remove(args, #args)
    end

    -- parse path
    if #args == 2 then
        if string.match(args[1], "^/") then
            args.path = args[1]
        end
        table.remove(args, 1)
    end

    -- parse target
    if string.match(args[1], "^/") or string.match(args[1], "^https?://") then
        args.target = args[1]
    end
    table.remove(args, 1)

    return args
end

return parse
