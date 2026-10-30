package provider

import (
	"context"
	"testing"

	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func TestOrganizationUserMappingFillKeepsUserIDWhenOmitted(t *testing.T) {
	m := OrganizationUserMappingResourceModel{
		Organization: supertypes.NewStringValue("acme"),
		UserId:       supertypes.NewInt64Value(2),
	}

	diags := m.Fill(context.Background(), apiclient.ExternalUser{
		Id:            "4",
		ExternalName:  "@renamed-user",
		Provider:      "gitlab",
		IntegrationId: "9",
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if m.UserId.ValueInt64() != 2 {
		t.Fatalf("user id = %d, want 2", m.UserId.ValueInt64())
	}
	if m.ExternalName.ValueString() != "@renamed-user" {
		t.Fatalf("external name = %q, want @renamed-user", m.ExternalName.ValueString())
	}
	if m.Id.ValueString() != "acme/4" {
		t.Fatalf("id = %q, want acme/4", m.Id.ValueString())
	}
}

func TestOrganizationUserMappingFillErrorsWhenUserIDMissing(t *testing.T) {
	m := OrganizationUserMappingResourceModel{
		Organization: supertypes.NewStringValue("acme"),
	}

	diags := m.Fill(context.Background(), apiclient.ExternalUser{
		Id:            "4",
		ExternalName:  "@renamed-user",
		Provider:      "gitlab",
		IntegrationId: "9",
	})
	if !diags.HasError() {
		t.Fatal("expected an error when userId is omitted and no id is already set")
	}
}
