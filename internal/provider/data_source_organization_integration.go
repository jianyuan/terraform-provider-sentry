package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

type OrganizationIntegrationDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	InternalId   types.String `tfsdk:"internal_id"`
	Organization types.String `tfsdk:"organization"`
	ProviderKey  types.String `tfsdk:"provider_key"`
	Name         types.String `tfsdk:"name"`
}

func (m *OrganizationIntegrationDataSourceModel) Fill(ctx context.Context, d apiclient.OrganizationIntegration) (diags diag.Diagnostics) {
	m.Id = types.StringValue(d.Id)
	m.InternalId = types.StringValue(d.Id)
	m.ProviderKey = types.StringValue(d.Provider.Key)
	m.Name = types.StringValue(d.Name)
	return
}

var _ datasource.DataSource = &OrganizationIntegrationDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationIntegrationDataSource{}

func NewOrganizationIntegrationDataSource() datasource.DataSource {
	return &OrganizationIntegrationDataSource{}
}

type OrganizationIntegrationDataSource struct {
	baseDataSource
}

func (d *OrganizationIntegrationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_integration"
}

func (d *OrganizationIntegrationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Organization Integration data source. See the [Sentry documentation](https://docs.sentry.io/api/integrations/list-an-organizations-available-integrations/) for more information.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
			},
			"internal_id": schema.StringAttribute{
				MarkdownDescription: "The internal ID for this organization integration. **Deprecated** Use `id` instead.",
				Computed:            true,
				DeprecationMessage:  "This field is deprecated and will be removed in a future version. Use `id` instead.",
			},
			"organization": DataSourceOrganizationAttribute(),
			"provider_key": schema.StringAttribute{
				Description: "Specific integration provider to filter by such as `slack`. See [the list of supported providers](https://docs.sentry.io/product/integrations/).",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the integration.",
				Required:    true,
			},
		},
	}
}

func (d *OrganizationIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationIntegrationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var matchedIntegrations []apiclient.OrganizationIntegration
	params := &apiclient.ListOrganizationIntegrationsParams{
		ProviderKey: data.ProviderKey.ValueStringPointer(),
	}

	for {
		httpResp, err := d.apiClient.ListOrganizationIntegrationsWithResponse(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return
		}

		for _, integration := range *httpResp.JSON200 {
			if integration.Name == data.Name.ValueString() {
				matchedIntegrations = append(matchedIntegrations, integration)
			}
		}

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	if len(matchedIntegrations) == 0 {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("organization integration"))
		return
	} else if len(matchedIntegrations) > 1 {
		resp.Diagnostics.AddError("Not unique", "More than one matching organization integration found")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, matchedIntegrations[0])...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
