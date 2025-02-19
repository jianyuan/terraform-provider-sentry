package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

type DataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Slug       types.String `tfsdk:"slug"`
	Name       types.String `tfsdk:"name"`
	InternalId types.String `tfsdk:"internal_id"`
}

func (m *DataSourceModel) Fill(ctx context.Context, org apiclient.Organization) (diags diag.Diagnostics) {
	m.Id = types.StringValue(org.Slug)
	m.Slug = types.StringValue(org.Slug)
	m.Name = types.StringValue(org.Name)
	m.InternalId = types.StringValue(org.Id)
	return
}
