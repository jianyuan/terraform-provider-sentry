package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// stringSetValue returns an empty set for a missing list.
func stringSetValue(ctx context.Context, values *[]string) (types.Set, diag.Diagnostics) {
	elements := []string{}
	if values != nil {
		elements = append(elements, *values...)
	}
	return types.SetValueFrom(ctx, types.StringType, elements)
}

// stringSetElements returns a non-nil slice, because Sentry rejects null for list fields.
func stringSetElements(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	elements := []string{}
	diags := set.ElementsAs(ctx, &elements, false)
	return elements, diags
}
