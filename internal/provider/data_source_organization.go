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

var _ datasource.DataSource = &OrganizationDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationDataSource{}

func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

type OrganizationDataSource struct {
	baseDataSource
}

func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = gen.OrganizationDataSourceSchema(ctx)
}

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data gen.OrganizationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.apiClient.GetOrganizationWithResponse(
		ctx,
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

	resp.Diagnostics.Append(fillDataSourceOrganizationModel(ctx, &data, *httpResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fillDataSourceOrganizationModel(ctx context.Context, m *gen.OrganizationModel, resp apiclient.GetOrganizationResponse) (diags diag.Diagnostics) {
	m.Slug = types.StringValue(resp.JSON200.Slug)
	m.Name = types.StringValue(resp.JSON200.Name)
	m.InternalId = types.StringValue(resp.JSON200.Id)
	m.Id = types.StringValue(resp.JSON200.Slug) // Deprecated
	return
}
