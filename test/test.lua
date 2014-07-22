-- add path so we can require from lua directory
package.path = "../lua/?.lua;" .. package.path

require "redirect_parse"
require "redirect_translate"

print("All tests passed!")
