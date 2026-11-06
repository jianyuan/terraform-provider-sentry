package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/oapi-codegen/nullable"
)

func TestRelayPiiConfigValue(t *testing.T) {
	for _, tt := range []struct {
		name     string
		value    nullable.Nullable[string]
		expected jsontypes.Normalized
	}{
		{name: "unspecified", value: nullable.Nullable[string]{}, expected: jsontypes.NewNormalizedNull()},
		{name: "null", value: nullable.NewNullNullable[string](), expected: jsontypes.NewNormalizedNull()},
		{name: "empty string", value: nullable.NewNullableWithValue(""), expected: jsontypes.NewNormalizedNull()},
		{name: "whitespace", value: nullable.NewNullableWithValue(" \n"), expected: jsontypes.NewNormalizedNull()},
		{
			// Matches jsonencode(): compact, sorted keys, array order kept.
			name:     "rules",
			value:    nullable.NewNullableWithValue(`{"rules": {"0": {"type": "anything", "redaction": {"method": "remove"}}}, "applications": {"$message": ["1", "0"]}}`),
			expected: jsontypes.NewNormalizedValue(`{"applications":{"$message":["1","0"]},"rules":{"0":{"redaction":{"method":"remove"},"type":"anything"}}}`),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := relayPiiConfigValue(tt.value)
			if !got.Equal(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

// The configured value is kept when Sentry returns the same rules in another layout, so applies converge without a diff.
func TestRelayPiiConfigValue_semanticEquality(t *testing.T) {
	configured := jsontypes.NewNormalizedValue(`{
		"rules": {"0": {"type": "anything", "redaction": {"method": "remove"}}},
		"applications": {"extra.foo": ["0"], "$message": ["0", "1"]}
	}`)

	for _, tt := range []struct {
		name      string
		returned  string
		wantEqual bool
	}{
		{
			name:      "same rules, other key order and whitespace",
			returned:  `{"applications":{"$message":["0","1"],"extra.foo":["0"]},"rules":{"0":{"redaction":{"method":"remove"},"type":"anything"}}}`,
			wantEqual: true,
		},
		{
			name:      "application order differs",
			returned:  `{"applications":{"$message":["1","0"],"extra.foo":["0"]},"rules":{"0":{"redaction":{"method":"remove"},"type":"anything"}}}`,
			wantEqual: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			returned := relayPiiConfigValue(nullable.NewNullableWithValue(tt.returned))
			equal, diags := configured.StringSemanticEquals(context.Background(), returned)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if equal != tt.wantEqual {
				t.Fatalf("expected semantic equality %v, got %v", tt.wantEqual, equal)
			}
		})
	}
}
