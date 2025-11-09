package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/provider/gen"
)

var _ datasource.DataSource = &TeamDataSource{}
var _ datasource.DataSourceWithConfigure = &TeamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

type TeamDataSource struct {
	baseDataSource
}

func (d *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = gen.TeamDataSourceSchema(ctx)
}

func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data gen.TeamModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.apiClient.GetOrganizationTeamWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Slug.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(fillDataSourceTeamModel(ctx, &data, *httpResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fillDataSourceTeamModel(ctx context.Context, m *gen.TeamModel, resp apiclient.GetOrganizationTeamResponse) (diags diag.Diagnostics) {
	m.Slug = types.StringValue(resp.JSON200.Slug)
	m.InternalId = types.StringValue(resp.JSON200.Id)
	m.Name = types.StringValue(resp.JSON200.Name)
	m.Id = types.StringValue(resp.JSON200.Slug)           // Deprecated
	m.HasAccess = types.BoolValue(resp.JSON200.HasAccess) // Deprecated
	m.IsPending = types.BoolValue(resp.JSON200.IsPending) // Deprecated
	m.IsMember = types.BoolValue(resp.JSON200.IsMember)   // Deprecated
	return
}
