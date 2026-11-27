package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type ClientKeyDataSourceModel struct {
	ClientKeyLookupModel
	DsnPublic types.String `tfsdk:"dsn_public"`
	DsnSecret types.String `tfsdk:"dsn_secret"`
	DsnCsp    types.String `tfsdk:"dsn_csp"`
}

func (m *ClientKeyDataSourceModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
	diags.Append(m.ClientKeyLookupModel.Fill(ctx, key)...)
	if diags.HasError() {
		return
	}

	if v, ok := key.Dsn["public"]; ok {
		m.DsnPublic = types.StringValue(v)
	} else {
		m.DsnPublic = types.StringNull()
	}

	if v, ok := key.Dsn["secret"]; ok {
		m.DsnSecret = types.StringValue(v)
	} else {
		m.DsnSecret = types.StringNull()
	}

	if v, ok := key.Dsn["csp"]; ok {
		m.DsnCsp = types.StringValue(v)
	} else {
		m.DsnCsp = types.StringNull()
	}

	return
}

var _ datasource.DataSource = &ClientKeyDataSource{}
var _ datasource.DataSourceWithConfigure = &ClientKeyDataSource{}
var _ datasource.DataSourceWithConfigValidators = &ClientKeyDataSource{}

func NewClientKeyDataSource() datasource.DataSource {
	return &ClientKeyDataSource{}
}

type ClientKeyDataSource struct {
	baseDataSource
}

func (d *ClientKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *ClientKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve a Project's Client Key.",
		Attributes: map[string]schema.Attribute{
			"organization": DataSourceOrganizationAttribute(),
			"project":      DataSourceProjectAttribute(),
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the client key.",
				Optional:            true,
			},
			"first": schema.BoolAttribute{
				MarkdownDescription: "Return the first key of the returned keys.",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project that the key belongs to.",
				Computed:            true,
			},
			"public": schema.StringAttribute{
				MarkdownDescription: "The public key.",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "The secret key.",
				Computed:            true,
				Sensitive:           true,
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
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[ClientKeyJavascriptLoaderScriptModel](ctx),
				Attributes: map[string]schema.Attribute{
					"browser_sdk_version": schema.StringAttribute{
						MarkdownDescription: "The version of the browser SDK to load.",
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
				ElementType:         types.StringType,
				CustomType:          supertypes.NewMapTypeOf[string](ctx),
				Computed:            true,
				Sensitive:           true,
			},
			"dsn_public": schema.StringAttribute{
				MarkdownDescription: "The DSN tells the SDK where to send the events to. **Deprecated** Use `dsn[\"public\"]` instead.",
				DeprecationMessage:  "This field is deprecated and will be removed in a future version. Use `dsn[\"public\"]` instead.",
				Computed:            true,
			},
			"dsn_secret": schema.StringAttribute{
				MarkdownDescription: "Deprecated DSN includes a secret which is no longer required by newer SDK versions. If you are unsure which to use, follow installation instructions for your language. **Deprecated** Use `dsn[\"secret\"]` instead.",
				DeprecationMessage:  "This field is deprecated and will be removed in a future version. Use `dsn[\"secret\"]` instead.",
				Computed:            true,
				Sensitive:           true,
			},
			"dsn_csp": schema.StringAttribute{
				MarkdownDescription: "Security header endpoint for features like CSP and Expect-CT reports. **Deprecated** Use `dsn[\"csp\"]` instead.",
				DeprecationMessage:  "This field is deprecated and will be removed in a future version. Use `dsn[\"csp\"]` instead.",
				Computed:            true,
			},
		},
	}
}

func (d *ClientKeyDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
			path.MatchRoot("first"),
		),
	}
}

func (d *ClientKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClientKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	foundKey, diags := findClientKey(ctx, d.apiClient, data.ClientKeyLookupModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *foundKey)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
