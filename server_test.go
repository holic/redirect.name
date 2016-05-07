package main

import "testing"

func TestGetRedirectSimple(t *testing.T) {
	var redirect *Redirect
	var err error

	dnsTXT := []string{
		"Redirects from /test/* to https://github.com/holic/*",
	}

	redirect, err = getRedirect(dnsTXT, "/test/")
	assertEqual(t, err, nil)
	assertEqual(t, redirect.Location, "https://github.com/holic/")

	redirect, err = getRedirect(dnsTXT, "/test/success")
	assertEqual(t, err, nil)
	assertEqual(t, redirect.Location, "https://github.com/holic/success")

	redirect, err = getRedirect(dnsTXT, "/should/fail")
	assertEqual(t, err.Error(), "No paths matched")
}

func TestGetRedirectComplex(t *testing.T) {
	// Tests that catchalls (even interspersed in the TXT records) apply
	// only after more specific matches
	var redirect *Redirect
	var err error

	dnsTXT := []string{
		"Redirects from /test/* to https://github.com/holic/*",
		"Redirects to https://github.com/holic",
		"Redirects from /noglob/ to https://github.com/holic/noglob",
	}

	redirect, err = getRedirect(dnsTXT, "/")
	assertEqual(t, err, nil)
	assertEqual(t, redirect.Location, "https://github.com/holic")

	redirect, err = getRedirect(dnsTXT, "/test/somepath")
	assertEqual(t, err, nil)
	assertEqual(t, redirect.Location, "https://github.com/holic/somepath")

	redirect, err = getRedirect(dnsTXT, "/noglob/")
	assertEqual(t, err, nil)
	assertEqual(t, redirect.Location, "https://github.com/holic/noglob")

	redirect, err = getRedirect(dnsTXT, "/catch/all")
	assertEqual(t, err, nil)
	assertEqual(t, redirect.Location, "https://github.com/holic")
}
