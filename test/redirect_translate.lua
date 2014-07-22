local expect = require "expect"
local translate = require "redirect.translate"
local location, status

location, status = translate("/")
expect(location):to_be(nil)
expect(status):to_be(nil)

location, status = translate("/", { target = "https://example.com/" })
expect(location):to_be("https://example.com/")
expect(status):to_be(302)

location, status = translate("/", { target = "https://example.com/", status = "301" })
expect(location):to_be("https://example.com/")
expect(status):to_be(301)

location, status = translate("/", { path = "/twitter", target = "http://example.com/", status = "301" })
expect(location):to_be(nil)
expect(status):to_be(nil)

location, status = translate("/", { path = "/", target = "http://example.com/" })
expect(location):to_be("http://example.com/")
expect(status):to_be(302)

location, status = translate("/about-us", { path = "/*", target = "http://example.com/" })
expect(location):to_be("http://example.com/")
expect(status):to_be(302)

location, status = translate("/about-us", { path = "/*", target = "http://example.com/*" })
expect(location):to_be("http://example.com/about-us")
expect(status):to_be(302)

location, status = translate("/blog/1", { path = "/*/1", target = "http://example.com/*" })
expect(location):to_be("http://example.com/blog")
expect(status):to_be(302)
