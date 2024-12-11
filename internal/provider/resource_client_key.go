package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/maputils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type ClientKeyJavascriptLoaderScriptResourceModel struct {
	BrowserSdkVersion            types.String `tfsdk:"browser_sdk_version"`
	PerformanceMonitoringEnabled types.Bool   `tfsdk:"performance_monitoring_enabled"`
	SessionReplayEnabled         types.Bool   `tfsdk:"session_replay_enabled"`
	DebugEnabled                 types.Bool   `tfsdk:"debug_enabled"`
}

type ClientKeyResourceModel struct {
	Id                     types.String                                  `tfsdk:"id"`
	Organization           types.String                                  `tfsdk:"organization"`
	Project                types.String                                  `tfsdk:"project"`
	ProjectId              types.String                                  `tfsdk:"project_id"`
	Name                   types.String                                  `tfsdk:"name"`
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

func (m *ClientKeyResourceModel) Fill(key apiclient.ProjectKey) error {
	m.Id = types.StringValue(key.Id)
	m.ProjectId = types.StringValue(key.ProjectId.String())
	m.Name = types.StringValue(key.Name)

	if key.RateLimit == nil {
		m.RateLimitWindow = types.Int64Null()
		m.RateLimitCount = types.Int64Null()
	} else {
		m.RateLimitWindow = types.Int64Value(key.RateLimit.Window)
		m.RateLimitCount = types.Int64Value(key.RateLimit.Count)
	}

	if m.JavascriptLoaderScript != nil {
		if !m.JavascriptLoaderScript.BrowserSdkVersion.IsNull() {
			m.JavascriptLoaderScript.BrowserSdkVersion = types.StringValue(key.BrowserSdkVersion)
		}

		if !m.JavascriptLoaderScript.PerformanceMonitoringEnabled.IsNull() {
			m.JavascriptLoaderScript.PerformanceMonitoringEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasPerformance)
		}

		if !m.JavascriptLoaderScript.SessionReplayEnabled.IsNull() {
			m.JavascriptLoaderScript.SessionReplayEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasReplay)
		}

		if !m.JavascriptLoaderScript.DebugEnabled.IsNull() {
			m.JavascriptLoaderScript.DebugEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasDebug)
		}
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

func (m *ClientKeyResourceModel) FillAll(organization string, project string, key apiclient.ProjectKey) error {
	m.Id = types.StringValue(key.Id)
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(project)
	m.ProjectId = types.StringValue(key.ProjectId.String())
	m.Name = types.StringValue(key.Name)

	if key.RateLimit == nil {
		m.RateLimitWindow = types.Int64Null()
		m.RateLimitCount = types.Int64Null()
	} else {
		m.RateLimitWindow = types.Int64Value(key.RateLimit.Window)
		m.RateLimitCount = types.Int64Value(key.RateLimit.Count)
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

var _ resource.Resource = &ClientKeyResource{}
var _ resource.ResourceWithConfigure = &ClientKeyResource{}
var _ resource.ResourceWithConfigValidators = &ClientKeyResource{}
var _ resource.ResourceWithImportState = &ClientKeyResource{}

func NewClientKeyResource() resource.Resource {
	return &ClientKeyResource{}
}

type ClientKeyResource struct {
	baseResource
}

func (r *ClientKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *ClientKeyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.RequiredTogether(
			path.MatchRoot("rate_limit_window"),
			path.MatchRoot("rate_limit_count"),
		),
	}
}

func (r *ClientKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Return a client key bound to a project.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the client key.",
				Required:            true,
			},
			"rate_limit_window": schema.Int64Attribute{
				MarkdownDescription: "Length of time in seconds that will be considered when checking the rate limit.",
				Optional:            true,
			},
			"rate_limit_count": schema.Int64Attribute{
				MarkdownDescription: "Number of events that can be reported within the rate limit window.",
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
		Blocks: map[string]schema.Block{
			"javascript_loader_script": schema.SingleNestedBlock{
				MarkdownDescription: "The JavaScript loader script configuration.",
				Attributes: map[string]schema.Attribute{
					"browser_sdk_version": schema.StringAttribute{
						MarkdownDescription: "The version of the browser SDK to load. Valid values are `7.x` and `8.x`.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("7.x", "8.x"),
						},
					},
					"performance_monitoring_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether performance monitoring is enabled for this key.",
						Optional:            true,
					},
					"session_replay_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether session replay is enabled for this key.",
						Optional:            true,
					},
					"debug_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether debug bundles & logging are enabled for this key.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *ClientKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClientKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.CreateProjectClientKeyJSONRequestBody{
		Name: data.Name.ValueString(),
	}

	if !data.RateLimitWindow.IsNull() || !data.RateLimitCount.IsNull() {
		body.RateLimit = &struct {
			Count  int64 `json:"count"`
			Window int64 `json:"window"`
		}{
			Count:  data.RateLimitCount.ValueInt64(),
			Window: data.RateLimitWindow.ValueInt64(),
		}
	}

	if data.JavascriptLoaderScript != nil {
		body.BrowserSdkVersion = data.JavascriptLoaderScript.BrowserSdkVersion.ValueStringPointer()

		if !data.JavascriptLoaderScript.SessionReplayEnabled.IsNull() ||
			!data.JavascriptLoaderScript.PerformanceMonitoringEnabled.IsNull() ||
			!data.JavascriptLoaderScript.DebugEnabled.IsNull() {
			body.DynamicSdkLoaderOptions = &struct {
				HasDebug       *bool `json:"hasDebug,omitempty"`
				HasPerformance *bool `json:"hasPerformance,omitempty"`
				HasReplay      *bool `json:"hasReplay,omitempty"`
			}{
				HasReplay:      data.JavascriptLoaderScript.SessionReplayEnabled.ValueBoolPointer(),
				HasDebug:       data.JavascriptLoaderScript.DebugEnabled.ValueBoolPointer(),
				HasPerformance: data.JavascriptLoaderScript.PerformanceMonitoringEnabled.ValueBoolPointer(),
			}
		}
	}

	httpResp, err := r.apiClient.CreateProjectClientKeyWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if err := data.Fill(*httpResp.JSON201); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClientKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetProjectClientKeyWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("client key"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ClientKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.UpdateProjectClientKeyJSONRequestBody{
		Name: plan.Name.ValueStringPointer(),
	}

	if !plan.RateLimitWindow.Equal(state.RateLimitWindow) || !plan.RateLimitCount.Equal(state.RateLimitCount) {
		body.RateLimit = &struct {
			Count  int64 `json:"count"`
			Window int64 `json:"window"`
		}{
			Count:  plan.RateLimitCount.ValueInt64(),
			Window: plan.RateLimitWindow.ValueInt64(),
		}
	}

	if plan.JavascriptLoaderScript != nil {
		if state.JavascriptLoaderScript == nil {
			body.BrowserSdkVersion = plan.JavascriptLoaderScript.BrowserSdkVersion.ValueStringPointer()

			if !plan.JavascriptLoaderScript.SessionReplayEnabled.IsNull() ||
				!plan.JavascriptLoaderScript.PerformanceMonitoringEnabled.IsNull() ||
				!plan.JavascriptLoaderScript.DebugEnabled.IsNull() {
				body.DynamicSdkLoaderOptions = &struct {
					HasDebug       *bool `json:"hasDebug,omitempty"`
					HasPerformance *bool `json:"hasPerformance,omitempty"`
					HasReplay      *bool `json:"hasReplay,omitempty"`
				}{
					HasReplay:      plan.JavascriptLoaderScript.SessionReplayEnabled.ValueBoolPointer(),
					HasDebug:       plan.JavascriptLoaderScript.DebugEnabled.ValueBoolPointer(),
					HasPerformance: plan.JavascriptLoaderScript.PerformanceMonitoringEnabled.ValueBoolPointer(),
				}
			}
		} else {
			if !plan.JavascriptLoaderScript.BrowserSdkVersion.Equal(state.JavascriptLoaderScript.BrowserSdkVersion) {
				body.BrowserSdkVersion = plan.JavascriptLoaderScript.BrowserSdkVersion.ValueStringPointer()
			}

			if !plan.JavascriptLoaderScript.SessionReplayEnabled.Equal(state.JavascriptLoaderScript.SessionReplayEnabled) ||
				!plan.JavascriptLoaderScript.PerformanceMonitoringEnabled.Equal(state.JavascriptLoaderScript.PerformanceMonitoringEnabled) ||
				!plan.JavascriptLoaderScript.DebugEnabled.Equal(state.JavascriptLoaderScript.DebugEnabled) {
				body.DynamicSdkLoaderOptions = &struct {
					HasDebug       *bool `json:"hasDebug,omitempty"`
					HasPerformance *bool `json:"hasPerformance,omitempty"`
					HasReplay      *bool `json:"hasReplay,omitempty"`
				}{
					HasReplay:      plan.JavascriptLoaderScript.SessionReplayEnabled.ValueBoolPointer(),
					HasDebug:       plan.JavascriptLoaderScript.DebugEnabled.ValueBoolPointer(),
					HasPerformance: plan.JavascriptLoaderScript.PerformanceMonitoringEnabled.ValueBoolPointer(),
				}
			}
		}
	}

	httpResp, err := r.apiClient.UpdateProjectClientKeyWithResponse(
		ctx,
		plan.Organization.ValueString(),
		plan.Project.ValueString(),
		plan.Id.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("client key"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if err := plan.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClientKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClientKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteProjectClientKeyWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *ClientKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateThreePartId(ctx, "organization", "project", req, resp)
}
