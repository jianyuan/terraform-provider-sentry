package provider

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/oapi-codegen/nullable"
)

const relayPiiConfigMarkdownDescription = "Advanced data scrubbing rules, as a JSON string holding the `rules` and `applications` objects. Build it with `jsonencode()`. The value replaces all advanced data scrubbing rules. For details, see [Advanced Data Scrubbing](https://docs.sentry.io/security-legal-pii/scrubbing/advanced-datascrubbing/)."

// relayPiiConfigValue maps relayPiiConfig to state. Sentry returns null or an empty string when no rules are set.
// The value is compacted with sorted keys, the layout jsonencode() produces, so imported state matches configuration.
func relayPiiConfigValue(v nullable.Nullable[string]) jsontypes.Normalized {
	s, err := v.Get()
	if err != nil || strings.TrimSpace(s) == "" {
		return jsontypes.NewNormalizedNull()
	}

	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var value any
	if err := dec.Decode(&value); err != nil {
		return jsontypes.NewNormalizedValue(s)
	}
	if b, err := json.Marshal(value); err == nil {
		s = string(b)
	}
	return jsontypes.NewNormalizedValue(s)
}
