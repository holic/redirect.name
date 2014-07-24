local expect = require "expect"
local parse = require "redirect.parse"
local args

args = parse("Redirect to http://github.com.com/holic")
expect(args.path):to_be(nil)
expect(args.target):to_be("http://github.com.com/holic")
expect(args.status):to_be(nil)

args = parse("Redirect permanently to http://github.com.com/holic")
expect(args.path):to_be(nil)
expect(args.target):to_be("http://github.com.com/holic")
expect(args.status):to_be("permanently")

args = parse("Redirect from / to https://github.com.com/holic")
expect(args.path):to_be("/")
expect(args.target):to_be("https://github.com.com/holic")
expect(args.status):to_be(nil)

args = parse("Redirects to /new from /old")
expect(args.path):to_be("/old")
expect(args.target):to_be("/new")
expect(args.status):to_be(nil)

args = parse("Redirects from /old to /new permanently")
expect(args.path):to_be("/old")
expect(args.target):to_be("/new")
expect(args.status):to_be("permanently")

args = parse("Redirects temporarily to /new")
expect(args.path):to_be(nil)
expect(args.target):to_be("/new")
expect(args.status):to_be("temporarily")
