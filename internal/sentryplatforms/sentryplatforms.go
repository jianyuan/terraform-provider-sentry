package sentryplatforms

import (
	_ "embed"
	"strings"
)

//go:generate bash ./generate.sh

//go:embed platforms.txt
var rawPlatforms string
var platforms = strings.Split(strings.TrimSpace(rawPlatforms), "\n")

// Validate checks if a platform is valid from the list loaded from platforms.txt.
func Validate(platform string) bool {
	// "other" is a special case that is always valid
	if platform == "other" {
		return true
	}

	for _, p := range platforms {
		if p == platform {
			return true
		}
	}

	return false
}
