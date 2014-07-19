-- fall back to documentation site
local function fallback (reason)
    local target = "http://redirect.name/"
    if reason then target = target .. "#reason=" .. ngx.escape_uri(reason) end
    ngx.redirect(target)
end


local resolver = require "resty.dns.resolver"
-- instantiate DNS resolver
local r, err = resolver:new { nameservers = { "8.8.8.8", "8.8.4.4" } }
if not r then return fallback("Failed to instantiate the resolver: " .. err) end

-- query hostname for TXT records
local answers, err = r:query(ngx.var.host, { qtype = resolver.TYPE_TXT })
if not answers then return fallback("Failed to query the DNS server: " .. err) end

-- fall back when query returns bad hostname, non-existent hostname, or other errors
if answers.errcode then return fallback("Could not resolve hostname: " .. answers.errcode .. " " .. answers.errstr) end


local parse_redirect = require "parse_redirect"
local uri = ngx.var.request_uri

for i, ans in ipairs(answers) do
    if ans.type == resolver.TYPE_TXT then
        local location, status = parse_redirect(uri, ans.txt)

        if location then
            return ngx.redirect(location, status)
        end
    end
end

-- fall back if no records match
return fallback("No paths matched")
