package main

import "testing"

func TestTranslate(t *testing.T) {
	var redirect *Redirect

	redirect = Translate("GET", "/", nil)
	if redirect != nil {
		t.Errorf("Expected %#v to be %#v", redirect, nil)
	}

	redirect = Translate("GET", "/", &Config{To: "https://example.com/"})
	assertEqual(t, redirect.Location, "https://example.com/")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("GET", "/", &Config{To: "https://example.com/", RedirectState: "301"})
	assertEqual(t, redirect.Location, "https://example.com/")
	assertEqual(t, redirect.Status, 301)

	redirect = Translate("GET", "/", &Config{From: "/twitter", To: "https://example.com/", RedirectState: "permanently"})
	if redirect != nil {
		t.Errorf("Expected %#v to be %#v", redirect, nil)
	}

	redirect = Translate("GET", "/", &Config{From: "/", To: "https://example.com/", RedirectState: "permanently"})
	assertEqual(t, redirect.Location, "https://example.com/")
	assertEqual(t, redirect.Status, 301)
}

func TestTranslateWildcard(t *testing.T) {
	var redirect *Redirect

	redirect = Translate("GET", "/about-us", &Config{From: "/*", To: "http://example.com/"})
	assertEqual(t, redirect.Location, "http://example.com/")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("GET", "/about-us", &Config{From: "/*", To: "http://example.com/*"})
	assertEqual(t, redirect.Location, "http://example.com/about-us")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("GET", "/about-us", &Config{From: "/*", To: "http://example.com/*"})
	assertEqual(t, redirect.Location, "http://example.com/about-us")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("GET", "/blog/1", &Config{From: "/*/1", To: "http://example.com/*", RedirectState: "temporarily"})
	assertEqual(t, redirect.Location, "http://example.com/blog")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("GET", "/wildcard", &Config{From: "/*", To: "http://example.com/**"})
	assertEqual(t, redirect.Location, "http://example.com/wildcard*")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("GET", "/wildcard", &Config{From: "/**", To: "http://example.com/*"})
	if redirect != nil {
		t.Errorf("Expected %#v to be %#v", redirect, nil)
	}

	redirect = Translate("GET", "/wildcard*", &Config{From: "/**", To: "http://example.com/*"})
	assertEqual(t, redirect.Location, "http://example.com/wildcard")
	assertEqual(t, redirect.Status, 302)
}

func TestTranslateStatus(t *testing.T)  {
	var redirect *Redirect

	redirect = Translate("POST", "/", &Config{From: "/*", To: "http://example.com/*"})
	assertEqual(t, redirect.Status, 307)

	redirect = Translate("POST", "/", &Config{From: "/*", To: "http://example.com/*", RedirectState: "permanently"})
	assertEqual(t, redirect.Status, 307)

	redirect = Translate("POST", "/", &Config{From: "/*", To: "http://example.com/*", RedirectState: "temporarily"})
	assertEqual(t, redirect.Status, 307)
}
