package main

import "testing"

func assertEqual(t *testing.T, value interface{}, expectation interface{}) {
	if value != expectation {
		t.Errorf("Expected %#v to be %#v", value, expectation)
	}
}

func TestParse(t *testing.T) {
	var config *Config

	config = Parse("This is not a valid statement")
	if config != nil {
		t.Errorf("Expected %#v to be %#v", config, nil)
	}

	config = Parse("Redirect to http://github.com/holic")
	assertEqual(t, config.From, "")
	assertEqual(t, config.To, "http://github.com/holic")
	assertEqual(t, config.RedirectState, "")

	config = Parse("Redirect from / to http://github.com/holic")
	assertEqual(t, config.From, "/")
	assertEqual(t, config.To, "http://github.com/holic")
	assertEqual(t, config.RedirectState, "")

	config = Parse("Redirects to /new from /old")
	assertEqual(t, config.From, "/old")
	assertEqual(t, config.To, "/new")
	assertEqual(t, config.RedirectState, "")

	config = Parse("Redirects from /old to /new permanently")
	assertEqual(t, config.From, "/old")
	assertEqual(t, config.To, "/new")
	assertEqual(t, config.RedirectState, "permanently")

	config = Parse("Redirects temporarily to /new")
	assertEqual(t, config.From, "")
	assertEqual(t, config.To, "/new")
	assertEqual(t, config.RedirectState, "temporarily")
}
