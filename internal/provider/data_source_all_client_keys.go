package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

type AllClientKeysDataSourceModel struct {
	Organization types.String             `tfsdk:"organization"`
	Project      types.String             `tfsdk:"project"`
	FilterStatus types.String             `tfsdk:"filter_status"`
	Keys         []ClientKeyResourceModel `tfsdk:"keys"`
}

func (m *AllClientKeysDataSourceModel) Fill(ctx context.Context, keys []apiclient.ProjectKey) (diags diag.Diagnostics) {
	m.Keys = make([]ClientKeyResourceModel, len(keys))
	for i, key := range keys {
		m.Keys[i].Organization = types.StringValue(m.Organization.ValueString())
		m.Keys[i].Project = types.StringValue(m.Project.ValueString())
		diags.Append(m.Keys[i].Fill(ctx, key)...)
	}

	return
}

var _ datasource.DataSource = &AllClientKeysDataSource{}
var _ datasource.DataSourceWithConfigure = &AllClientKeysDataSource{}

func NewAllClientKeysDataSource() datasource.DataSource {
	return &AllClientKeysDataSource{}
}

type AllClientKeysDataSource struct {
	baseDataSource
}

func (d *AllClientKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_keys"
}

func (d *AllClientKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Project's Client Keys.",

		Attributes: map[string]schema.Attribute{
			"organization": DataSourceOrganizationAttribute(),
			"project":      DataSourceProjectAttribute(),
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
							MarkdownDescription: "Length of time in seconds that will be considered when checking the rate limit.",
							Computed:            true,
						},
						"rate_limit_count": schema.Int64Attribute{
							MarkdownDescription: "Number of events that can be reported within the rate limit window.",
							Computed:            true,
						},
						"javascript_loader_script": schema.SingleNestedAttribute{
							MarkdownDescription: "The JavaScript loader script configuration.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"browser_sdk_version": schema.StringAttribute{
									MarkdownDescription: "The version of the browser SDK to load. Valid values are `7.x` and `8.x`.",
									Computed:            true,
								},
								"performance_monitoring_enabled": schema.BoolAttribute{
									MarkdownDescription: "Whether performance monitoring is enabled for this key.",
									Computed:            true,
								},
								"session_replay_enabled": schema.BoolAttribute{
									MarkdownDescription: "Whether session replay is enabled for this key.",
									Computed:            true,
								},
								"debug_enabled": schema.BoolAttribute{
									MarkdownDescription: "Whether debug bundles & logging are enabled for this key.",
									Computed:            true,
								},
							},
						},
						"dsn": schema.MapAttribute{
							MarkdownDescription: "This is a map of DSN values. The keys include `public`, `secret`, `csp`, `security`, `minidump`, `nel`, `unreal`, `cdn`, and `crons`.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"dsn_public": schema.StringAttribute{
							MarkdownDescription: "The DSN tells the SDK where to send the events to. **Deprecated** Use `dsn[\"public\"]` instead.",
							DeprecationMessage:  "This field is deprecated and will be removed in a future version. Use `dsn[\"public\"]` instead.",
							Computed:            true,
						},
						"dsn_secret": schema.StringAttribute{
							MarkdownDescription: "Deprecated DSN includes a secret which is no longer required by newer SDK versions. If you are unsure which to use, follow installation instructions for your language. **Deprecated** Use `dsn[\"secret\"] instead.",
							DeprecationMessage:  "This field is deprecated and will be removed in a future version. Use `dsn[\"secret\"]` instead.",
							Computed:            true,
						},
						"dsn_csp": schema.StringAttribute{
							MarkdownDescription: "Security header endpoint for features like CSP and Expect-CT reports. **Deprecated** Use `dsn[\"csp\"]` instead.",
							DeprecationMessage:  "This field is deprecated and will be removed in a future version. Use `dsn[\"csp\"]` instead.",
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

	var allKeys []apiclient.ProjectKey
	params := &apiclient.ListProjectClientKeysParams{
		Status: (*apiclient.ListProjectClientKeysParamsStatus)(data.FilterStatus.ValueStringPointer()),
	}
	for {
		httpResp, err := d.apiClient.ListProjectClientKeysWithResponse(
			ctx,
			data.Organization.ValueString(),
			data.Project.ValueString(),
			params,
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}
		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return
		}

		allKeys = append(allKeys, *httpResp.JSON200...)

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	resp.Diagnostics.Append(data.Fill(ctx, allKeys)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
