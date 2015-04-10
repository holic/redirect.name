package parser

import "regexp"

type Config struct {
	Target string
	Path   string
	Status string
}

var configRE = regexp.MustCompile(`Redirects?(\s+.*)`)
var targetRE = regexp.MustCompile(`\s+to\s+(https?\://\S+|/\S*)`)
var pathRE = regexp.MustCompile(`\s+from\s+(\S+)`)
var statusRE = regexp.MustCompile(`\s+(permanently|temporarily)`)

func Parse(record string) *Config {
	configMatches := configRE.FindStringSubmatch(record)
	if len(configMatches) == 0 {
		return nil
	}

	targetMatches := targetRE.FindStringSubmatch(configMatches[1])
	pathMatches := pathRE.FindStringSubmatch(configMatches[1])
	statusMatches := statusRE.FindStringSubmatch(configMatches[1])

	config := new(Config)
	if len(targetMatches) > 0 {
		config.Target = targetMatches[1]
	}
	if len(pathMatches) > 0 {
		config.Path = pathMatches[1]
	}
	if len(statusMatches) > 0 {
		config.Status = statusMatches[1]
	}

	return config
}
