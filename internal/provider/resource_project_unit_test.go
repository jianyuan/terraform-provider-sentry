package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func TestProjectResourceModel_FillDataScrubbing(t *testing.T) {
	const base = `"id": "1", "slug": "web", "name": "Web", "platform": "javascript",
		"organization": {"id": "1234", "slug": "my-org", "name": "My Org", "orgRoleList": [], "teamRoleList": []},
		"teams": [], "features": [], "options": {
			"filters:blacklisted_ips": "", "filters:releases": "", "filters:error_messages": ""
		}`

	type scrubbing struct {
		DataScrubber         types.Bool
		DataScrubberDefaults types.Bool
		ScrubIpAddresses     types.Bool
		SensitiveFields      types.Set
		SafeFields           types.Set
		RelayPiiConfig       jsontypes.Normalized
	}

	for _, tt := range []struct {
		name     string
		response string
		expected scrubbing
	}{
		{
			name: "all settings",
			response: `{` + base + `,
				"dataScrubber": true, "dataScrubberDefaults": false, "scrubIPAddresses": true,
				"sensitiveFields": ["email"], "safeFields": ["order_id"],
				"relayPiiConfig": "{\"rules\": {}, \"applications\": {}}"
			}`,
			expected: scrubbing{
				DataScrubber:         types.BoolValue(true),
				DataScrubberDefaults: types.BoolValue(false),
				ScrubIpAddresses:     types.BoolValue(true),
				SensitiveFields:      dataScrubbingStringSet("email"),
				SafeFields:           dataScrubbingStringSet("order_id"),
				RelayPiiConfig:       jsontypes.NewNormalizedValue(`{"applications":{},"rules":{}}`),
			},
		},
		{
			// A project whose rules were cleared stores an empty string.
			name: "no rules",
			response: `{` + base + `,
				"dataScrubber": true, "dataScrubberDefaults": true, "scrubIPAddresses": false,
				"sensitiveFields": [], "safeFields": [], "relayPiiConfig": ""
			}`,
			expected: scrubbing{
				DataScrubber:         types.BoolValue(true),
				DataScrubberDefaults: types.BoolValue(true),
				ScrubIpAddresses:     types.BoolValue(false),
				SensitiveFields:      dataScrubbingStringSet(),
				SafeFields:           dataScrubbingStringSet(),
				RelayPiiConfig:       jsontypes.NewNormalizedNull(),
			},
		},
		{
			name:     "settings missing from the response",
			response: `{` + base + `}`,
			expected: scrubbing{
				DataScrubber:         types.BoolNull(),
				DataScrubberDefaults: types.BoolNull(),
				ScrubIpAddresses:     types.BoolNull(),
				SensitiveFields:      types.SetNull(types.StringType),
				SafeFields:           types.SetNull(types.StringType),
				RelayPiiConfig:       jsontypes.NewNormalizedNull(),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var project apiclient.Project
			if err := json.Unmarshal([]byte(tt.response), &project); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var data ProjectResourceModel
			if diags := data.Fill(context.Background(), project); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			got := scrubbing{
				DataScrubber:         data.DataScrubber,
				DataScrubberDefaults: data.DataScrubberDefaults,
				ScrubIpAddresses:     data.ScrubIpAddresses,
				SensitiveFields:      data.SensitiveFields,
				SafeFields:           data.SafeFields,
				RelayPiiConfig:       data.RelayPiiConfig,
			}
			if diff := cmp.Diff(tt.expected, got, dataScrubbingCmpOptions); diff != "" {
				t.Errorf("unexpected model (-want +got):\n%s", diff)
			}
		})
	}
}
