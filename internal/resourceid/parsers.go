package resourceid

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ParsePath extracts 1 identifier part from a short key ("partA").
func ParsePath(rawInput, labelA string) (string, error) {
	parts, err := SplitPath(rawInput, labelA)
	if err != nil {
		return "", err
	}
	return parts[0], nil
}

// Split2Path extracts 2 identifier parts from a short key ("partA/partB").
func Split2Path(rawInput, labelA, labelB string) (string, string, error) {
	parts, err := SplitPath(rawInput, labelA, labelB)
	if err != nil {
		return "", "", err
	}
	return parts[0], parts[1], nil
}

// Split3Path extracts 3 identifier parts from a short key ("partA/partB/partC").
func Split3Path(rawInput, labelA, labelB, labelC string) (string, string, string, error) {
	parts, err := SplitPath(rawInput, labelA, labelB, labelC)
	if err != nil {
		return "", "", "", err
	}
	return parts[0], parts[1], parts[2], nil
}

// Split4Path extracts 4 identifier parts from a short key ("partA/partB/partC/partD").
func Split4Path(rawInput, labelA, labelB, labelC, labelD string) (string, string, string, string, error) {
	parts, err := SplitPath(rawInput, labelA, labelB, labelC, labelD)
	if err != nil {
		return "", "", "", "", err
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

// SplitPath extracts N slash-separated parts given expected label names.
func SplitPath(rawInput string, labels ...string) ([]string, error) {
	expectedCount := len(labels)
	input := strings.TrimSpace(rawInput)
	expectedFormat := strings.Join(labels, "/")

	if input == "" {
		return nil, fmt.Errorf("import ID cannot be empty, expected format %s", expectedFormat)
	}

	parts := strings.Split(input, "/")
	if len(parts) != expectedCount {
		return nil, fmt.Errorf("unexpected ID format (%s), expected %s", rawInput, expectedFormat)
	}

	for i, part := range parts {
		if strings.TrimSpace(part) == "" {
			return nil, fmt.Errorf("unexpected ID format (%s), segment for %s cannot be empty", rawInput, labels[i])
		}
	}

	return parts, nil
}

// Parse extracts 1 identifier part matching the specified label token in rawURLTemplate.
func Parse(rawInput, rawURLTemplate, labelA string) (string, error) {
	parts, err := Split(rawInput, rawURLTemplate, labelA)
	if err != nil {
		return "", err
	}
	return parts[0], nil
}

// Split2 extracts 2 identifier parts matching labelA and labelB in rawURLTemplate.
func Split2(rawInput, rawURLTemplate, labelA, labelB string) (string, string, error) {
	parts, err := Split(rawInput, rawURLTemplate, labelA, labelB)
	if err != nil {
		return "", "", err
	}
	return parts[0], parts[1], nil
}

// Split3 extracts 3 identifier parts matching labelA, labelB, and labelC in rawURLTemplate.
func Split3(rawInput, rawURLTemplate, labelA, labelB, labelC string) (string, string, string, error) {
	parts, err := Split(rawInput, rawURLTemplate, labelA, labelB, labelC)
	if err != nil {
		return "", "", "", err
	}
	return parts[0], parts[1], parts[2], nil
}

// Split4 extracts 4 identifier parts matching labelA, labelB, labelC, and labelD in rawURLTemplate.
func Split4(rawInput, rawURLTemplate, labelA, labelB, labelC, labelD string) (string, string, string, string, error) {
	parts, err := Split(rawInput, rawURLTemplate, labelA, labelB, labelC, labelD)
	if err != nil {
		return "", "", "", "", err
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

// Split extracts N identifier parts matching the requested labels in order.
func Split(rawInput, rawURLTemplate string, labels ...string) ([]string, error) {
	expectedCount := len(labels)
	input := strings.TrimSpace(rawInput)
	expectedFormat := strings.Join(labels, "/")

	labelIndices := make(map[string]int, expectedCount)
	for i, l := range labels {
		labelIndices[l] = i
	}

	// 1. Handle Full URL Inputs
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		inputURL, err := url.Parse(input)
		if err != nil {
			return nil, fmt.Errorf("invalid URL (%s): %w", rawInput, err)
		}

		// Sanitize host placeholders so url.Parse does not fail on '{' or '}'
		sanitizedTmpl := sanitizeTemplateForParsing(rawURLTemplate)
		tmplURL, err := url.Parse(sanitizedTmpl)
		if err != nil {
			return nil, fmt.Errorf("invalid template URL (%s): %w", rawURLTemplate, err)
		}

		results := make([]string, expectedCount)

		// 1a. Match Subdomain / Host placeholders
		inputHost := strings.Split(inputURL.Hostname(), ".")
		tmplHost := strings.Split(tmplURL.Hostname(), ".")
		if len(inputHost) == len(tmplHost) {
			for i, tmplPart := range tmplHost {
				if label := restoreLabel(tmplPart); label != "" {
					if idx, exists := labelIndices[label]; exists {
						results[idx] = inputHost[i]
					}
				}
			}
		}

		// 1b. Match Path segment placeholders
		inputPath := strings.Split(strings.Trim(inputURL.Path, "/"), "/")
		tmplPath := strings.Split(strings.Trim(tmplURL.Path, "/"), "/")
		if len(inputPath) == len(tmplPath) {
			for i, tmplSegment := range tmplPath {
				if label := restoreLabel(tmplSegment); label != "" {
					if idx, exists := labelIndices[label]; exists {
						results[idx] = inputPath[i]
					}
				}
			}
		}

		// 1c. Match Query Parameter placeholders
		inputQuery := inputURL.Query()
		tmplQuery := tmplURL.Query()

		for paramKey, tmplValues := range tmplQuery {
			for _, tmplVal := range tmplValues {
				if label := restoreLabel(tmplVal); label != "" {
					if idx, exists := labelIndices[label]; exists {
						if actualVal := inputQuery.Get(paramKey); actualVal != "" {
							results[idx] = actualVal
						}
					}
				}
			}
		}

		// Verify all requested labels were matched
		for idx, val := range results {
			if val == "" {
				return nil, fmt.Errorf("could not extract placeholder {%s} from URL (%s) using template (%s)", labels[idx], rawInput, rawURLTemplate)
			}
		}

		return results, nil
	}

	// 2. Handle Short Path Inputs ("partA/partB/...")
	parts := strings.Split(input, "/")
	if len(parts) != expectedCount {
		return nil, fmt.Errorf("unexpected format (%s), expected %s or full URL (%s)", rawInput, expectedFormat, rawURLTemplate)
	}

	for i, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("unexpected format (%s), segment for %s cannot be empty", rawInput, labels[i])
		}
	}

	return parts, nil
}

// Regex to find {placeholder} inside URLs
var placeholderRegex = regexp.MustCompile(`\{([a-zA-Z0-9_-]+)\}`)

// Replaces `{label}` with `tmpl-placeholder-label` so url.Parse accepts the host component
func sanitizeTemplateForParsing(tmpl string) string {
	return placeholderRegex.ReplaceAllString(tmpl, "tmpl-placeholder-$1")
}

// Converts `tmpl-placeholder-label` back to `label` (or handles `{label}` if still present)
func restoreLabel(token string) string {
	if strings.HasPrefix(token, "tmpl-placeholder-") {
		return strings.TrimPrefix(token, "tmpl-placeholder-")
	}
	if strings.HasPrefix(token, "{") && strings.HasSuffix(token, "}") {
		return strings.Trim(token, "{}")
	}
	return ""
}
