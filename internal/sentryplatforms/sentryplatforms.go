package sentryplatforms

import (
	_ "embed"
	"strings"
)

//go:generate bash ./generate.sh

//go:embed platforms.txt
var rawPlatforms string

// Platforms is a list of valid Sentry platforms.
var Platforms = strings.Split(strings.TrimSpace(rawPlatforms), "\n")
