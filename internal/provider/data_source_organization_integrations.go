package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

type OrganizationIntegrationsDataSourceModel struct {
	ProviderKey  types.String                                                                          `tfsdk:"provider_key"`
	Organization types.String                                                                          `tfsdk:"organization"`
	Integrations supertypes.SetNestedObjectValueOf[OrganizationIntegrationsDataSourceModelIntegration] `tfsdk:"integrations"`
}

func (m *OrganizationIntegrationsDataSourceModel) Fill(ctx context.Context, data []apiclient.OrganizationIntegration) (diags diag.Diagnostics) {
	m.Integrations = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(data, func(item apiclient.OrganizationIntegration, _ int) OrganizationIntegrationsDataSourceModelIntegration {
		var model OrganizationIntegrationsDataSourceModelIntegration
		diags.Append(model.Fill(ctx, item)...)
		return model
	}))
	return
}

type OrganizationIntegrationsDataSourceModelIntegration struct {
	Id          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	ProviderKey types.String         `tfsdk:"provider_key"`
	RawJson     jsontypes.Normalized `tfsdk:"raw_json"`
}

func (m *OrganizationIntegrationsDataSourceModelIntegration) Fill(ctx context.Context, data apiclient.OrganizationIntegration) (diags diag.Diagnostics) {
	m.Id = types.StringValue(data.Id)
	m.Name = types.StringValue(data.Name)
	m.ProviderKey = types.StringValue(data.Provider.Key)
	m.RawJson = func() jsontypes.Normalized {
		b, err := data.MarshalJSON()
		if err != nil {
			diags.AddError("failed to marshal organization integration", err.Error())
			return jsontypes.NewNormalizedUnknown()
		}
		return jsontypes.NewNormalizedValue(string(b))
	}()
	return
}

var _ datasource.DataSource = &OrganizationIntegrationsDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationIntegrationsDataSource{}

func NewOrganizationIntegrationsDataSource() datasource.DataSource {
	return &OrganizationIntegrationsDataSource{}
}

type OrganizationIntegrationsDataSource struct {
	baseDataSource
}

func (d *OrganizationIntegrationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_integrations"
}

func (d *OrganizationIntegrationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of all integrations connected to an organization. For more information, see the [Sentry documentation](https://docs.sentry.io/api/integrations/list-an-organizations-available-integrations/).\n\nUse this data source to fetch integrations in bulk. You can optionally filter by a `provider_key` or omit it to retrieve all integrations across all providers.\n\nTo look up a single integration by name, use the [`sentry_organization_integration`](organization_integration.md) data source instead.",

		Attributes: map[string]schema.Attribute{
			"organization": DataSourceOrganizationAttribute(),
			"provider_key": schema.StringAttribute{
				MarkdownDescription: "Specific integration provider to filter by such as `slack`. See [the list of supported providers](https://docs.sentry.io/integrations/).",
				Optional:            true,
			},
			"integrations": schema.SetNestedAttribute{
				MarkdownDescription: "List of integrations.",
				Computed:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[OrganizationIntegrationsDataSourceModelIntegration](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the integration.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the integration.",
							Computed:            true,
						},
						"provider_key": schema.StringAttribute{
							MarkdownDescription: "The provider of the integration such as `slack`. See [the list of supported providers](https://docs.sentry.io/integrations/).",
							Computed:            true,
						},
						"raw_json": schema.StringAttribute{
							MarkdownDescription: "Raw JSON representation of the integration. Use this if you need to access fields that are not explicitly exposed by the provider.",
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationIntegrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationIntegrationsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var integrations []apiclient.OrganizationIntegration
	params := &apiclient.ListOrganizationIntegrationsParams{
		ProviderKey:   data.ProviderKey.ValueStringPointer(),
		IncludeConfig: new(apiclient.ListOrganizationIntegrationsParamsIncludeConfigFalse),
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

		integrations = append(integrations, *httpResp.JSON200...)

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	resp.Diagnostics.Append(data.Fill(ctx, integrations)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
