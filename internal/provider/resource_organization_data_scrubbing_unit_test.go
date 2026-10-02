package provider

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/samber/lo"
)

// Framework values have unexported fields, so cmp needs their Equal methods.
var organizationDataScrubbingCmpOptions = cmp.Options{
	cmp.Comparer(func(a, b types.String) bool { return a.Equal(b) }),
	cmp.Comparer(func(a, b types.Bool) bool { return a.Equal(b) }),
	cmp.Comparer(func(a, b types.Set) bool { return a.Equal(b) }),
}

func organizationDataScrubbingStringSet(values ...string) types.Set {
	elements := lo.Map(values, func(v string, _ int) attr.Value { return types.StringValue(v) })
	return types.SetValueMust(types.StringType, elements)
}

func TestOrganizationDataScrubbingResourceModel_Fill(t *testing.T) {
	for _, tt := range []struct {
		name         string
		configured   types.String
		organization apiclient.Organization
		expected     OrganizationDataScrubbingResourceModel
	}{
		{
			name:       "all settings",
			configured: types.StringValue("my-org"),
			organization: apiclient.Organization{
				Slug:                 "my-org",
				DataScrubber:         lo.ToPtr(true),
				DataScrubberDefaults: lo.ToPtr(true),
				SensitiveFields:      lo.ToPtr([]string{"email", "IBAN"}),
				SafeFields:           lo.ToPtr([]string{"order_id"}),
				ScrubIPAddresses:     lo.ToPtr(true),
			},
			expected: OrganizationDataScrubbingResourceModel{
				Id:                   types.StringValue("my-org"),
				Organization:         types.StringValue("my-org"),
				DataScrubber:         types.BoolValue(true),
				DataScrubberDefaults: types.BoolValue(true),
				SensitiveFields:      organizationDataScrubbingStringSet("email", "IBAN"),
				SafeFields:           organizationDataScrubbingStringSet("order_id"),
				ScrubIpAddresses:     types.BoolValue(true),
			},
		},
		{
			name:         "missing settings",
			configured:   types.StringValue("my-org"),
			organization: apiclient.Organization{Slug: "my-org"},
			expected: OrganizationDataScrubbingResourceModel{
				Id:                   types.StringValue("my-org"),
				Organization:         types.StringValue("my-org"),
				DataScrubber:         types.BoolValue(false),
				DataScrubberDefaults: types.BoolValue(false),
				SensitiveFields:      organizationDataScrubbingStringSet(),
				SafeFields:           organizationDataScrubbingStringSet(),
				ScrubIpAddresses:     types.BoolValue(false),
			},
		},
		{
			name:         "organization configured by internal ID",
			configured:   types.StringValue("1234"),
			organization: apiclient.Organization{Id: "1234", Slug: "my-org"},
			expected: OrganizationDataScrubbingResourceModel{
				Id:                   types.StringValue("my-org"),
				Organization:         types.StringValue("1234"),
				DataScrubber:         types.BoolValue(false),
				DataScrubberDefaults: types.BoolValue(false),
				SensitiveFields:      organizationDataScrubbingStringSet(),
				SafeFields:           organizationDataScrubbingStringSet(),
				ScrubIpAddresses:     types.BoolValue(false),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			data := OrganizationDataScrubbingResourceModel{Organization: tt.configured}
			if diags := data.Fill(context.Background(), tt.organization); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if diff := cmp.Diff(tt.expected, data, organizationDataScrubbingCmpOptions); diff != "" {
				t.Errorf("unexpected model (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOrganizationDataScrubbingResourceModel_UpdateBody(t *testing.T) {
	for _, tt := range []struct {
		name     string
		data     OrganizationDataScrubbingResourceModel
		expected apiclient.UpdateOrganizationJSONRequestBody
	}{
		{
			name: "settings",
			data: OrganizationDataScrubbingResourceModel{
				Organization:         types.StringValue("my-org"),
				DataScrubber:         types.BoolValue(true),
				DataScrubberDefaults: types.BoolValue(false),
				SensitiveFields:      organizationDataScrubbingStringSet("email"),
				SafeFields:           organizationDataScrubbingStringSet(),
				ScrubIpAddresses:     types.BoolValue(true),
			},
			expected: apiclient.UpdateOrganizationJSONRequestBody{
				DataScrubber:         lo.ToPtr(true),
				DataScrubberDefaults: lo.ToPtr(false),
				SensitiveFields:      &[]string{"email"},
				SafeFields:           &[]string{},
				ScrubIPAddresses:     lo.ToPtr(true),
			},
		},
		{
			name: "defaults",
			data: organizationDataScrubbingDefaults(),
			expected: apiclient.UpdateOrganizationJSONRequestBody{
				DataScrubber:         lo.ToPtr(false),
				DataScrubberDefaults: lo.ToPtr(false),
				SensitiveFields:      &[]string{},
				SafeFields:           &[]string{},
				ScrubIPAddresses:     lo.ToPtr(false),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			body, diags := tt.data.UpdateBody(context.Background())
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if diff := cmp.Diff(tt.expected, body); diff != "" {
				t.Errorf("unexpected body (-want +got):\n%s", diff)
			}
		})
	}
}
