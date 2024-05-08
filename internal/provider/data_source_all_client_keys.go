package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &AllClientKeysDataSource{}
var _ datasource.DataSourceWithConfigure = &AllClientKeysDataSource{}

func NewAllClientKeysDataSource() datasource.DataSource {
	return &AllClientKeysDataSource{}
}

type AllClientKeysDataSource struct {
	baseDataSource
}

type AllClientKeysDataSourceModel struct {
	Organization types.String             `tfsdk:"organization"`
	Project      types.String             `tfsdk:"project"`
	FilterStatus types.String             `tfsdk:"filter_status"`
	Keys         []ClientKeyResourceModel `tfsdk:"keys"`
}

func (m *AllClientKeysDataSourceModel) Fill(organization string, project string, filterStatus *string, keys []*sentry.ProjectKey) error {
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(project)
	m.FilterStatus = types.StringPointerValue(filterStatus)

	m.Keys = []ClientKeyResourceModel{}
	for _, key := range keys {
		var model ClientKeyResourceModel
		if err := model.Fill(organization, project, *key); err != nil {
			return err
		}

		m.Keys = append(m.Keys, model)
	}

	return nil
}

func (d *AllClientKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_keys"
}

func (d *AllClientKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Project's Client Keys.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the resource belongs to.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The slug of the project the resource belongs to.",
				Required:            true,
			},
			"filter_status": schema.StringAttribute{
				MarkdownDescription: "Filter client keys by `active` or `inactive`. Defaults to returning all keys if not specified.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"active",
						"inactive",
					),
				},
			},
			"keys": schema.ListNestedAttribute{
				MarkdownDescription: "The list of client keys.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of this resource.",
							Computed:            true,
						},
						"organization": schema.StringAttribute{
							MarkdownDescription: "The slug of the organization the resource belongs to.",
							Computed:            true,
						},
						"project": schema.StringAttribute{
							MarkdownDescription: "The slug of the project the resource belongs to.",
							Computed:            true,
						},
						"project_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the project that the key belongs to.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the client key.",
							Computed:            true,
						},
						"public": schema.StringAttribute{
							MarkdownDescription: "The public key.",
							Computed:            true,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "The secret key.",
							Computed:            true,
						},
						"rate_limit_window": schema.Int64Attribute{
							MarkdownDescription: "Length of time that will be considered when checking the rate limit.",
							Computed:            true,
						},
						"rate_limit_count": schema.Int64Attribute{
							MarkdownDescription: "Number of events that can be reported within the rate limit window.",
							Computed:            true,
						},
						"dsn_public": schema.StringAttribute{
							MarkdownDescription: "The DSN tells the SDK where to send the events to.",
							Computed:            true,
						},
						"dsn_secret": schema.StringAttribute{
							MarkdownDescription: "Deprecated DSN includes a secret which is no longer required by newer SDK versions. If you are unsure which to use, follow installation instructions for your language.",
							Computed:            true,
						},
						"dsn_csp": schema.StringAttribute{
							MarkdownDescription: "Security header endpoint for features like CSP and Expect-CT reports.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AllClientKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllClientKeysDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allKeys []*sentry.ProjectKey
	params := &sentry.ListProjectKeysParams{
		Status: data.FilterStatus.ValueStringPointer(),
	}
	for {
		keys, apiResp, err := d.client.ProjectKeys.List(ctx, data.Organization.ValueString(), data.Project.ValueString(), params)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err))
			return
		}

		allKeys = append(allKeys, keys...)

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	if err := data.Fill(data.Organization.ValueString(), data.Project.ValueString(), data.FilterStatus.ValueStringPointer(), allKeys); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
