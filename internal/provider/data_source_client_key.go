package provider

import (
	"context"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/maputils"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

type ClientKeyJavascriptLoaderScriptDataSourceModel struct {
	BrowserSdkVersion            types.String `tfsdk:"browser_sdk_version"`
	PerformanceMonitoringEnabled types.Bool   `tfsdk:"performance_monitoring_enabled"`
	SessionReplayEnabled         types.Bool   `tfsdk:"session_replay_enabled"`
	DebugEnabled                 types.Bool   `tfsdk:"debug_enabled"`
}

type ClientKeyDataSourceModel struct {
	Organization           types.String                                  `tfsdk:"organization"`
	Project                types.String                                  `tfsdk:"project"`
	Id                     types.String                                  `tfsdk:"id"`
	Name                   types.String                                  `tfsdk:"name"`
	First                  types.Bool                                    `tfsdk:"first"`
	ProjectId              types.String                                  `tfsdk:"project_id"`
	RateLimitWindow        types.Int64                                   `tfsdk:"rate_limit_window"`
	RateLimitCount         types.Int64                                   `tfsdk:"rate_limit_count"`
	JavascriptLoaderScript *ClientKeyJavascriptLoaderScriptResourceModel `tfsdk:"javascript_loader_script"`
	Public                 types.String                                  `tfsdk:"public"`
	Secret                 types.String                                  `tfsdk:"secret"`
	Dsn                    types.Map                                     `tfsdk:"dsn"`
	DsnPublic              types.String                                  `tfsdk:"dsn_public"`
	DsnSecret              types.String                                  `tfsdk:"dsn_secret"`
	DsnCsp                 types.String                                  `tfsdk:"dsn_csp"`
}

func (m *ClientKeyDataSourceModel) Fill(key apiclient.ProjectKey) error {
	m.Id = types.StringValue(key.Id)
	m.ProjectId = types.StringValue(key.ProjectId.String())
	m.Name = types.StringValue(key.Name)

	if key.RateLimit == nil {
		m.RateLimitWindow = types.Int64Null()
		m.RateLimitCount = types.Int64Null()
	} else {
		m.RateLimitWindow = types.Int64Value(int64(key.RateLimit.Window))
		m.RateLimitCount = types.Int64Value(int64(key.RateLimit.Count))
	}

	m.JavascriptLoaderScript = &ClientKeyJavascriptLoaderScriptResourceModel{
		BrowserSdkVersion:            types.StringValue(key.BrowserSdkVersion),
		PerformanceMonitoringEnabled: types.BoolValue(key.DynamicSdkLoaderOptions.HasPerformance),
		SessionReplayEnabled:         types.BoolValue(key.DynamicSdkLoaderOptions.HasReplay),
		DebugEnabled:                 types.BoolValue(key.DynamicSdkLoaderOptions.HasDebug),
	}
	m.Public = types.StringValue(key.Public)
	m.Secret = types.StringValue(key.Secret)

	m.Dsn = types.MapValueMust(types.StringType, maputils.MapValues(key.Dsn, func(v string) attr.Value {
		return types.StringValue(v)
	}))

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

	return nil
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
				MarkdownDescription: "Boolean flag indicating that we want the first key of the returned keys.",
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
				Computed:            true,
				ElementType:         types.StringType,
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

	var foundKey *apiclient.ProjectKey

	if data.Id.IsNull() {
		var allKeys []apiclient.ProjectKey
		params := &apiclient.ListProjectClientKeysParams{}
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

		if data.Name.IsNull() {
			if len(allKeys) == 1 {
				foundKey = ptr.Ptr(allKeys[0])
			} else if !data.First.IsNull() && data.First.ValueBool() {
				// Find the first key

				// Sort keys by date created
				sort.Slice(allKeys, func(i, j int) bool {
					return allKeys[i].DateCreated.Before(allKeys[j].DateCreated)
				})

				foundKey = ptr.Ptr(allKeys[0])
			} else {
				resp.Diagnostics.AddError("Client error", "Multiple keys found, please specify the key by `name`, `id`, or set the `first` flag to `true`.")
				return
			}
		} else {
			// Find the key by name
			for _, key := range allKeys {
				if key.Name == data.Name.ValueString() {
					foundKey = ptr.Ptr(key)
					break
				}
			}
		}

	} else {
		// Get the key by ID
		httpResp, err := d.apiClient.GetProjectClientKeyWithResponse(
			ctx,
			data.Organization.ValueString(),
			data.Project.ValueString(),
			data.Id.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}
		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}

		foundKey = httpResp.JSON200
	}

	if foundKey == nil {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("client key"))
		return
	}

	if err := data.Fill(*foundKey); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
