package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

var _ datasource.DataSource = &OrganizationDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationDataSource{}

type OrganizationDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Slug       types.String `tfsdk:"slug"`
	Name       types.String `tfsdk:"name"`
	InternalId types.String `tfsdk:"internal_id"`
}

func (m *OrganizationDataSourceModel) Fill(org apiclient.Organization) error {
	m.Id = types.StringValue(org.Slug)
	m.Slug = types.StringValue(org.Slug)
	m.Name = types.StringValue(org.Name)
	m.InternalId = types.StringValue(org.Id)

	return nil
}

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
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Organization data source.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique URL slug for this organization.",
				Computed:            true,
			},
			"slug": DataSourceOrganizationAttribute(),
			"internal_id": schema.StringAttribute{
				MarkdownDescription: "The internal ID for this organization.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human readable name for this organization.",
				Computed:            true,
			},
		},
	}
}

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.apiClient.GetOrganizationWithResponse(
		ctx,
		data.Slug.ValueString(),
	)
	if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("organization"))
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
