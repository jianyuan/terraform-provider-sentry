package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/samber/lo"
)

// Framework values have unexported fields, so cmp needs their Equal methods.
var dataScrubbingCmpOptions = cmp.Options{
	cmp.Comparer(func(a, b types.String) bool { return a.Equal(b) }),
	cmp.Comparer(func(a, b types.Bool) bool { return a.Equal(b) }),
	cmp.Comparer(func(a, b types.Set) bool { return a.Equal(b) }),
	cmp.Comparer(func(a, b jsontypes.Normalized) bool { return a.Equal(b) }),
}

func dataScrubbingStringSet(values ...string) types.Set {
	return types.SetValueMust(types.StringType, lo.Map(values, func(v string, _ int) attr.Value { return types.StringValue(v) }))
}

func TestOrganizationDataScrubbingResourceModel_Fill(t *testing.T) {
	const rules = `{"applications":{"$message":["0"]},"rules":{"0":{"redaction":{"method":"mask"},"type":"creditcard"}}}`

	for _, tt := range []struct {
		name       string
		configured string
		response   string
		expected   OrganizationDataScrubbingResourceModel
	}{
		{
			name:       "all settings",
			configured: "my-org",
			response: `{
				"id": "1234", "slug": "my-org", "name": "My Org", "orgRoleList": [], "teamRoleList": [],
				"dataScrubber": true, "dataScrubberDefaults": true, "scrubIPAddresses": true,
				"sensitiveFields": ["email", "IBAN"], "safeFields": ["order_id"],
				"relayPiiConfig": ` + string(lo.Must(json.Marshal(rules))) + `
			}`,
			expected: OrganizationDataScrubbingResourceModel{
				Id:                   types.StringValue("my-org"),
				Organization:         types.StringValue("my-org"),
				DataScrubber:         types.BoolValue(true),
				DataScrubberDefaults: types.BoolValue(true),
				ScrubIpAddresses:     types.BoolValue(true),
				SensitiveFields:      dataScrubbingStringSet("email", "IBAN"),
				SafeFields:           dataScrubbingStringSet("order_id"),
				RelayPiiConfig:       jsontypes.NewNormalizedValue(rules),
			},
		},
		{
			name:       "Sentry defaults, organization configured by internal ID",
			configured: "1234",
			response: `{
				"id": "1234", "slug": "my-org", "name": "My Org", "orgRoleList": [], "teamRoleList": [],
				"dataScrubber": false, "dataScrubberDefaults": false, "scrubIPAddresses": false,
				"sensitiveFields": [], "safeFields": [], "relayPiiConfig": null
			}`,
			expected: OrganizationDataScrubbingResourceModel{
				Id:                   types.StringValue("my-org"),
				Organization:         types.StringValue("1234"),
				DataScrubber:         types.BoolValue(false),
				DataScrubberDefaults: types.BoolValue(false),
				ScrubIpAddresses:     types.BoolValue(false),
				SensitiveFields:      dataScrubbingStringSet(),
				SafeFields:           dataScrubbingStringSet(),
				RelayPiiConfig:       jsontypes.NewNormalizedNull(),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var organization apiclient.Organization
			if err := json.Unmarshal([]byte(tt.response), &organization); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			data := OrganizationDataScrubbingResourceModel{Organization: types.StringValue(tt.configured)}
			if diags := data.Fill(context.Background(), organization); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if diff := cmp.Diff(tt.expected, data, dataScrubbingCmpOptions); diff != "" {
				t.Errorf("unexpected model (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOrganizationDataScrubbingResourceModel_UpdateBody(t *testing.T) {
	for _, tt := range []struct {
		name     string
		data     OrganizationDataScrubbingResourceModel
		expected string
	}{
		{
			name: "settings",
			data: OrganizationDataScrubbingResourceModel{
				Organization:         types.StringValue("my-org"),
				DataScrubber:         types.BoolValue(true),
				DataScrubberDefaults: types.BoolValue(false),
				ScrubIpAddresses:     types.BoolValue(true),
				SensitiveFields:      dataScrubbingStringSet("email"),
				SafeFields:           dataScrubbingStringSet(),
				RelayPiiConfig:       jsontypes.NewNormalizedValue(`{"rules":{},"applications":{}}`),
			},
			expected: `{
				"dataScrubber": true, "dataScrubberDefaults": false, "scrubIPAddresses": true,
				"sensitiveFields": ["email"], "safeFields": [],
				"relayPiiConfig": "{\"rules\":{},\"applications\":{}}"
			}`,
		},
		{
			// Destroy sends this body: empty lists, not null, and an empty string to clear the rules.
			name: "defaults",
			data: organizationDataScrubbingDefaults(),
			expected: `{
				"dataScrubber": false, "dataScrubberDefaults": false, "scrubIPAddresses": false,
				"sensitiveFields": [], "safeFields": [], "relayPiiConfig": ""
			}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			body, diags := tt.data.UpdateBody(context.Background())
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			var got, want any
			lo.Must0(json.Unmarshal(lo.Must(json.Marshal(body)), &got))
			lo.Must0(json.Unmarshal([]byte(tt.expected), &want))
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("unexpected request body (-want +got):\n%s", diff)
			}
		})
	}
}
