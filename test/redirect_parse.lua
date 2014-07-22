local expect = require "expect"
local parse = require "redirect.parse"
local args

args = parse("redirect.name=http://github.com.com/holic")
expect(args.path):to_be(nil)
expect(args.target):to_be("http://github.com.com/holic")
expect(args.status):to_be(nil)

args = parse("redirect.name=http://github.com.com/holic 301")
expect(args.path):to_be(nil)
expect(args.target):to_be("http://github.com.com/holic")
expect(args.status):to_be("301")

args = parse("redirect.name=/ https://github.com.com/holic")
expect(args.path):to_be("/")
expect(args.target):to_be("https://github.com.com/holic")
expect(args.status):to_be(nil)

args = parse("redirect.name=/old /new")
expect(args.path):to_be("/old")
expect(args.target):to_be("/new")
expect(args.status):to_be(nil)

args = parse("redirect.name=/old /new 301")
expect(args.path):to_be("/old")
expect(args.target):to_be("/new")
expect(args.status):to_be("301")

args = parse("redirect.name=/new 302")
expect(args.path):to_be(nil)
expect(args.target):to_be("/new")
expect(args.status):to_be("302")
