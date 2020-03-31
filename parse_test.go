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

	config = Parse("Redirect to ftp://github.com")
	assertEqual(t, config.From, "")
	assertEqual(t, config.To, "ftp://github.com")
	assertEqual(t, config.RedirectState, "")

	config = Parse("Redirect to mailto:test@example.com")
	assertEqual(t, config.From, "")
	assertEqual(t, config.To, "mailto:test@example.com")
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

	// Test status codes

	config = Parse("Redirect to http://github.com/holic with 301")
	assertEqual(t, config.From, "")
	assertEqual(t, config.To, "http://github.com/holic")
	assertEqual(t, config.RedirectState, "301")

	config = Parse("Redirect with 302 from / to http://github.com/holic")
	assertEqual(t, config.From, "/")
	assertEqual(t, config.To, "http://github.com/holic")
	assertEqual(t, config.RedirectState, "302")

	config = Parse("Redirects to /new from /old with 307")
	assertEqual(t, config.From, "/old")
	assertEqual(t, config.To, "/new")
	assertEqual(t, config.RedirectState, "307")

	config = Parse("Redirects with 308 from /old to /new")
	assertEqual(t, config.From, "/old")
	assertEqual(t, config.To, "/new")
	assertEqual(t, config.RedirectState, "308")

	// Test that we get the first parsed value when multiple values
	// of the same type are specified (e.g. `permanently` and `with 308`)

	config = Parse("Redirects permanently from /old to /new with 302")
	assertEqual(t, config.From, "/old")
	assertEqual(t, config.To, "/new")
	assertEqual(t, config.RedirectState, "permanently")

	config = Parse("Redirects with 307 from /old to /new permanently")
	assertEqual(t, config.From, "/old")
	assertEqual(t, config.To, "/new")
	assertEqual(t, config.RedirectState, "307")
}
