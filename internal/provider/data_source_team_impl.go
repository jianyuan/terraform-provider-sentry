package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mzglinski/terraform-provider-sentry/internal/apiclient"
)

func (m *TeamDataSourceModel) Fill(ctx context.Context, data apiclient.Team) (diags diag.Diagnostics) {
	m.Slug.Set(data.Slug)
	m.InternalId.Set(data.Id)
	m.Name.Set(data.Name)
	m.Id.Set(data.Slug)                // Deprecated
	m.HasAccess.SetPtr(data.HasAccess) // Deprecated
	m.IsPending.SetPtr(data.IsPending) // Deprecated
	m.IsMember.SetPtr(data.IsMember)   // Deprecated
	return
}
