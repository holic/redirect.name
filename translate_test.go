package main

import "testing"

func TestTranslate(t *testing.T) {
	var redirect *Redirect

	redirect = Translate("/", nil)
	if redirect != nil {
		t.Errorf("Expected %#v to be %#v", redirect, nil)
	}

	redirect = Translate("/", &Config{To: "https://example.com/"})
	assertEqual(t, redirect.Location, "https://example.com/")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("/", &Config{To: "https://example.com/", RedirectState: "301"})
	assertEqual(t, redirect.Location, "https://example.com/")
	assertEqual(t, redirect.Status, 301)

	redirect = Translate("/", &Config{From: "/twitter", To: "https://example.com/", RedirectState: "permanently"})
	if redirect != nil {
		t.Errorf("Expected %#v to be %#v", redirect, nil)
	}

	redirect = Translate("/", &Config{From: "/", To: "https://example.com/", RedirectState: "permanently"})
	assertEqual(t, redirect.Location, "https://example.com/")
	assertEqual(t, redirect.Status, 301)
}

func TestTranslateWildcard(t *testing.T) {
	var redirect *Redirect

	redirect = Translate("/about-us", &Config{From: "/*", To: "http://example.com/"})
	assertEqual(t, redirect.Location, "http://example.com/")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("/about-us", &Config{From: "/*", To: "http://example.com/*"})
	assertEqual(t, redirect.Location, "http://example.com/about-us")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("/about-us", &Config{From: "/*", To: "http://example.com/*"})
	assertEqual(t, redirect.Location, "http://example.com/about-us")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("/blog/1", &Config{From: "/*/1", To: "http://example.com/*", RedirectState: "temporarily"})
	assertEqual(t, redirect.Location, "http://example.com/blog")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("/wildcard", &Config{From: "/*", To: "http://example.com/**"})
	assertEqual(t, redirect.Location, "http://example.com/wildcard*")
	assertEqual(t, redirect.Status, 302)

	redirect = Translate("/wildcard", &Config{From: "/**", To: "http://example.com/*"})
	if redirect != nil {
		t.Errorf("Expected %#v to be %#v", redirect, nil)
	}

	redirect = Translate("/wildcard*", &Config{From: "/**", To: "http://example.com/*"})
	assertEqual(t, redirect.Location, "http://example.com/wildcard")
	assertEqual(t, redirect.Status, 302)
}
