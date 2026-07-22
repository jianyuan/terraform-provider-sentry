package resourceid

import (
	"fmt"
	"net/url"
	"strings"
)

// BuildPath constructs a short slash-separated ID string ("partA/partB").
func BuildPath(parts ...string) (string, error) {
	for i, p := range parts {
		if strings.TrimSpace(p) == "" {
			return "", fmt.Errorf("part at index %d cannot be empty", i)
		}
	}
	return strings.Join(parts, "/"), nil
}

// Build1 generates a URL replacing labelA in rawURLTemplate.
func Build1(rawURLTemplate, labelA, valueA string) (string, error) {
	return Build(rawURLTemplate, map[string]string{labelA: valueA})
}

// Build2 generates a URL replacing labelA and labelB in rawURLTemplate.
func Build2(rawURLTemplate, labelA, valueA, labelB, valueB string) (string, error) {
	return Build(rawURLTemplate, map[string]string{
		labelA: valueA,
		labelB: valueB,
	})
}

// Build3 generates a URL replacing labelA, labelB, and labelC in rawURLTemplate.
func Build3(rawURLTemplate, labelA, valueA, labelB, valueB, labelC, valueC string) (string, error) {
	return Build(rawURLTemplate, map[string]string{
		labelA: valueA,
		labelB: valueB,
		labelC: valueC,
	})
}

// Build replaces key/value label pairs inside rawURLTemplate.
func Build(rawURLTemplate string, values map[string]string) (string, error) {
	if len(values) == 0 {
		return "", fmt.Errorf("at least one label value map must be provided")
	}

	result := rawURLTemplate
	for label, val := range values {
		if strings.TrimSpace(val) == "" {
			return "", fmt.Errorf("value for label {%s} cannot be empty", label)
		}

		placeholder := fmt.Sprintf("{%s}", label)
		if !strings.Contains(result, placeholder) {
			return "", fmt.Errorf("placeholder %s not found in template URL (%s)", placeholder, rawURLTemplate)
		}

		result = strings.ReplaceAll(result, placeholder, url.QueryEscape(val))
	}

	return result, nil
}
