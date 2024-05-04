package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &OrganizationIntegrationDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationIntegrationDataSource{}

func NewOrganizationIntegrationDataSource() datasource.DataSource {
	return &OrganizationIntegrationDataSource{}
}

type OrganizationIntegrationDataSource struct {
	baseDataSource
}

type OrganizationIntegrationDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	InternalId   types.String `tfsdk:"internal_id"`
	Organization types.String `tfsdk:"organization"`
	ProviderKey  types.String `tfsdk:"provider_key"`
	Name         types.String `tfsdk:"name"`
}

func (m *OrganizationIntegrationDataSourceModel) Fill(organizationSlug string, d sentry.OrganizationIntegration) error {
	m.Id = types.StringValue(d.ID)
	m.InternalId = types.StringValue(d.ID)
	m.Organization = types.StringValue(organizationSlug)
	m.ProviderKey = types.StringValue(d.Provider.Key)
	m.Name = types.StringValue(d.Name)

	return nil
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
			"organization": schema.StringAttribute{
				Description: "The slug of the organization.",
				Required:    true,
			},
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

	var matchedIntegrations []*sentry.OrganizationIntegration
	params := &sentry.ListOrganizationIntegrationsParams{
		ListCursorParams: sentry.ListCursorParams{},
		ProviderKey:      data.ProviderKey.ValueString(),
	}
	for {
		integrations, apiResp, err := d.client.OrganizationIntegrations.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read organization integrations, got error: %s", err))
			return
		}

		for _, integration := range integrations {
			if integration.Name == data.Name.ValueString() {
				matchedIntegrations = append(matchedIntegrations, integration)
			}
		}

		if apiResp.Cursor == "" {
			break
		}
		params.ListCursorParams.Cursor = apiResp.Cursor
	}

	if len(matchedIntegrations) == 0 {
		resp.Diagnostics.AddError("Not Found", "No matching organization integrations found")
		return
	} else if len(matchedIntegrations) > 1 {
		resp.Diagnostics.AddError("Not Unique", "More than one matching organization integration found")
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *matchedIntegrations[0]); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling organization integration: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
