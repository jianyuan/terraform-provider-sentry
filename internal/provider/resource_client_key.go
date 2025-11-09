package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/maputils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type ClientKeyResourceModel struct {
	Id                     types.String                                                               `tfsdk:"id"`
	Organization           types.String                                                               `tfsdk:"organization"`
	Project                types.String                                                               `tfsdk:"project"`
	ProjectId              types.String                                                               `tfsdk:"project_id"`
	Name                   types.String                                                               `tfsdk:"name"`
	RateLimitWindow        types.Int64                                                                `tfsdk:"rate_limit_window"`
	RateLimitCount         types.Int64                                                                `tfsdk:"rate_limit_count"`
	JavascriptLoaderScript supertypes.SingleNestedObjectValueOf[ClientKeyJavascriptLoaderScriptModel] `tfsdk:"javascript_loader_script"`
	Public                 types.String                                                               `tfsdk:"public"`
	Secret                 types.String                                                               `tfsdk:"secret"`
	Dsn                    types.Map                                                                  `tfsdk:"dsn"`
	DsnPublic              types.String                                                               `tfsdk:"dsn_public"`
	DsnSecret              types.String                                                               `tfsdk:"dsn_secret"`
	DsnCsp                 types.String                                                               `tfsdk:"dsn_csp"`
}

func (m *ClientKeyResourceModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
	m.Id = types.StringValue(key.Id)
	m.ProjectId = types.StringValue(key.ProjectId.String())
	m.Name = types.StringValue(key.Name)

	if v, err := key.RateLimit.Get(); err == nil {
		m.RateLimitWindow = types.Int64Value(v.Window)
		m.RateLimitCount = types.Int64Value(v.Count)

	} else {
		m.RateLimitWindow = types.Int64Null()
		m.RateLimitCount = types.Int64Null()
	}

	var javascriptLoaderScript ClientKeyJavascriptLoaderScriptModel
	diags.Append(javascriptLoaderScript.Fill(ctx, key)...)
	if diags.HasError() {
		return
	}

	m.JavascriptLoaderScript = supertypes.NewSingleNestedObjectValueOfNull[ClientKeyJavascriptLoaderScriptModel](ctx)
	diags.Append(m.JavascriptLoaderScript.Set(ctx, &javascriptLoaderScript)...)
	if diags.HasError() {
		return
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

	return
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
		resourcevalidator.RequiredTogether(
			path.MatchRoot("javascript_loader_script").AtName("performance_monitoring_enabled"),
			path.MatchRoot("javascript_loader_script").AtName("session_replay_enabled"),
			path.MatchRoot("javascript_loader_script").AtName("debug_enabled"),
		),
	}
}

func (r *ClientKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = clientKeySchema().GetResource(ctx)
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

	if !data.RateLimitWindow.IsUnknown() || !data.RateLimitCount.IsUnknown() {
		body.RateLimit = &struct {
			Count  int64 `json:"count"`
			Window int64 `json:"window"`
		}{
			Count:  data.RateLimitCount.ValueInt64(),
			Window: data.RateLimitWindow.ValueInt64(),
		}
	}

	if !data.JavascriptLoaderScript.IsUnknown() {
		javascriptLoaderScript := tfutils.MergeDiagnostics(data.JavascriptLoaderScript.Get(ctx))(&resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if !javascriptLoaderScript.BrowserSdkVersion.IsUnknown() {
			body.BrowserSdkVersion = javascriptLoaderScript.BrowserSdkVersion.ValueStringPointer()
		}

		if !javascriptLoaderScript.SessionReplayEnabled.IsUnknown() &&
			!javascriptLoaderScript.PerformanceMonitoringEnabled.IsUnknown() &&
			!javascriptLoaderScript.DebugEnabled.IsUnknown() {
			body.DynamicSdkLoaderOptions = &struct {
				HasDebug       *bool `json:"hasDebug,omitempty"`
				HasPerformance *bool `json:"hasPerformance,omitempty"`
				HasReplay      *bool `json:"hasReplay,omitempty"`
			}{
				HasReplay:      javascriptLoaderScript.SessionReplayEnabled.ValueBoolPointer(),
				HasDebug:       javascriptLoaderScript.DebugEnabled.ValueBoolPointer(),
				HasPerformance: javascriptLoaderScript.PerformanceMonitoringEnabled.ValueBoolPointer(),
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

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON201)...)
	if resp.Diagnostics.HasError() {
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

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
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

	if !plan.JavascriptLoaderScript.Equal(state.JavascriptLoaderScript) {
		javascriptLoaderScript := tfutils.MergeDiagnostics(plan.JavascriptLoaderScript.Get(ctx))(&resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// NOTE: Both `BrowserSdkVersion` and `DynamicSdkLoaderOptions` must be set together.
		body.BrowserSdkVersion = javascriptLoaderScript.BrowserSdkVersion.ValueStringPointer()
		body.DynamicSdkLoaderOptions = &struct {
			HasDebug       *bool `json:"hasDebug,omitempty"`
			HasPerformance *bool `json:"hasPerformance,omitempty"`
			HasReplay      *bool `json:"hasReplay,omitempty"`
		}{
			HasReplay:      javascriptLoaderScript.SessionReplayEnabled.ValueBoolPointer(),
			HasDebug:       javascriptLoaderScript.DebugEnabled.ValueBoolPointer(),
			HasPerformance: javascriptLoaderScript.PerformanceMonitoringEnabled.ValueBoolPointer(),
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

	resp.Diagnostics.Append(plan.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
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
