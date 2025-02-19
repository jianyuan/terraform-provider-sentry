package client_key

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/services"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithConfigure = &Resource{}
var _ resource.ResourceWithConfigValidators = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	services.BaseResource
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
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

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = Schema().GetResource(ctx)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ResourceModel

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

	httpResp, err := r.ApiClient.CreateProjectClientKeyWithResponse(
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

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.ApiClient.GetProjectClientKeyWithResponse(
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

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceModel

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

	httpResp, err := r.ApiClient.UpdateProjectClientKeyWithResponse(
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

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.ApiClient.DeleteProjectClientKeyWithResponse(
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

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateThreePartId(ctx, "organization", "project", req, resp)
}
