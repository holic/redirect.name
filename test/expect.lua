local function str (value)
	if value == nil then
		return "nil"
	elseif type(value) == "string" then
		return "\"" .. value .. "\""
	else
		return value
	end
end


local Expect = {}
Expect.__index = Expect

function Expect:create (value)
	local expect = {}
	setmetatable(expect, Expect)
	expect.value = value
	return expect
end

function Expect:to_be (value)
	assert(self.value == value, "Expected " .. str(self.value) .. " to be " .. str(value))
end

return function (value)
	return Expect:create(value)
end
