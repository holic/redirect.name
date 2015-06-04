package main

import (
	"bytes"
	"regexp"
	"strings"
)

type Redirect struct {
	Location string
	Status   int
}

func Translate(uri string, config *Config) *Redirect {
	if uri == "" {
		return nil
	}
	if config == nil {
		return nil
	}
	if config.To == "" {
		return nil
	}

	redirect := &Redirect{Location: config.To}

	switch config.RedirectState {
	case "301", "permanently":
		redirect.Status = 301
	case "302", "temporarily":
		redirect.Status = 302
	default:
		redirect.Status = 302
	}

	// no `From` assumes catch-all, so redirect immediately to `Location`
	if config.From == "" {
		return redirect
	}

	count := strings.Count(config.From, `*`)

	var exp bytes.Buffer
	exp.WriteString(`^`)
	exp.WriteString(strings.Replace(regexp.QuoteMeta(config.From), `\*`, `(.*)`, 1))
	exp.WriteString(`$`)

	fromRE := regexp.MustCompile(exp.String())

	// if we can't find the pattern, return to continue to next record
	if !fromRE.MatchString(uri) {
		return nil
	}

	// wildcard replacement of `uri` if there's a wildcard in our `From` path
	if count > 0 {
		redirect.Location = fromRE.ReplaceAllString(uri, strings.Replace(redirect.Location, `*`, `${1}`, 1))
	}

	return redirect
}
