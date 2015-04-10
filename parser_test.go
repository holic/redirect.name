package parser

import "testing"

func assertEqual(t *testing.T, value interface{}, expectation interface{}) {
	if value != expectation {
		t.Errorf("Expected %#v to be %#v", value, expectation)
	}
}

func TestParse(t *testing.T) {
	var config *Config

	config = Parse("Redirect to http://github.com/holic")
	assertEqual(t, config.Path, "")
	assertEqual(t, config.Target, "http://github.com/holic")
	assertEqual(t, config.Status, "")

	config = Parse("Redirect from / to http://github.com/holic")
	assertEqual(t, config.Path, "/")
	assertEqual(t, config.Target, "http://github.com/holic")
	assertEqual(t, config.Status, "")

	config = Parse("Redirects to /new from /old")
	assertEqual(t, config.Path, "/old")
	assertEqual(t, config.Target, "/new")
	assertEqual(t, config.Status, "")

	config = Parse("Redirects from /old to /new permanently")
	assertEqual(t, config.Path, "/old")
	assertEqual(t, config.Target, "/new")
	assertEqual(t, config.Status, "permanently")

	config = Parse("Redirects temporarily to /new")
	assertEqual(t, config.Path, "")
	assertEqual(t, config.Target, "/new")
	assertEqual(t, config.Status, "temporarily")
}
