package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (m *AlertDataSourceModel) fill(ctx context.Context, data apiclient.OrganizationWorkflow) (diags diag.Diagnostics) {
	triggersJson, err := data.Triggers.MarshalJSON()
	if err != nil {
		diags.AddError("Failed to marshal triggers", err.Error())
		return diags
	}
	m.TriggersJson = supertypes.NewStringValue(string(triggersJson))

	actionFiltersJson, err := data.ActionFilters.MarshalJSON()
	if err != nil {
		diags.AddError("Failed to marshal action filters", err.Error())
		return diags
	}

	m.ActionFiltersJson = supertypes.NewStringValue(string(actionFiltersJson))

	return
}
