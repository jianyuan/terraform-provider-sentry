package provider

import (
	"context"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type ClientKeyJavascriptLoaderScriptModel struct {
	BrowserSdkVersion            types.String `tfsdk:"browser_sdk_version"`
	PerformanceMonitoringEnabled types.Bool   `tfsdk:"performance_monitoring_enabled"`
	SessionReplayEnabled         types.Bool   `tfsdk:"session_replay_enabled"`
	DebugEnabled                 types.Bool   `tfsdk:"debug_enabled"`
}

func (m *ClientKeyJavascriptLoaderScriptModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
	m.BrowserSdkVersion = types.StringValue(key.BrowserSdkVersion)
	m.PerformanceMonitoringEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasPerformance)
	m.SessionReplayEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasReplay)
	m.DebugEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasDebug)
	return
}

// ClientKeyLookupModel holds the client key attributes shared by the
// `sentry_key` data source and ephemeral resource: the lookup inputs and the
// computed key details.
type ClientKeyLookupModel struct {
	Organization           types.String                                                               `tfsdk:"organization"`
	Project                types.String                                                               `tfsdk:"project"`
	Id                     types.String                                                               `tfsdk:"id"`
	Name                   types.String                                                               `tfsdk:"name"`
	First                  types.Bool                                                                 `tfsdk:"first"`
	ProjectId              types.String                                                               `tfsdk:"project_id"`
	RateLimitWindow        types.Int64                                                                `tfsdk:"rate_limit_window"`
	RateLimitCount         types.Int64                                                                `tfsdk:"rate_limit_count"`
	JavascriptLoaderScript supertypes.SingleNestedObjectValueOf[ClientKeyJavascriptLoaderScriptModel] `tfsdk:"javascript_loader_script"`
	Public                 types.String                                                               `tfsdk:"public"`
	Secret                 types.String                                                               `tfsdk:"secret"`
	Dsn                    supertypes.MapValueOf[string]                                              `tfsdk:"dsn"`
}

func (m *ClientKeyLookupModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
	m.Id = types.StringValue(key.Id)
	m.ProjectId = types.StringValue(key.ProjectId.String())
	m.Name = types.StringValue(key.Name)

	if v, err := key.RateLimit.Get(); err == nil {
		m.RateLimitWindow = types.Int64Value(int64(v.Window))
		m.RateLimitCount = types.Int64Value(int64(v.Count))
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

	var dsnDiags diag.Diagnostics
	m.Dsn, dsnDiags = supertypes.NewMapValueOfMap(ctx, key.Dsn)
	diags.Append(dsnDiags...)
	return
}

// findClientKey looks up a project's client key by `id`, by `name`, or, when
// neither is set, returns the only key or the oldest key if `first` is true.
func findClientKey(ctx context.Context, apiClient *apiclient.ClientWithResponses, m ClientKeyLookupModel) (*apiclient.ProjectKey, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !m.Id.IsNull() {
		httpResp, err := apiClient.GetProjectClientKeyWithResponse(
			ctx,
			m.Organization.ValueString(),
			m.Project.ValueString(),
			m.Id.ValueString(),
		)
		if err != nil {
			diags.Append(diagutils.NewClientError("read", err))
			return nil, diags
		}
		if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			diags.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return nil, diags
		}
		return httpResp.JSON200, diags
	}

	var allKeys []apiclient.ProjectKey
	params := &apiclient.ListProjectClientKeysParams{}
	for {
		httpResp, err := apiClient.ListProjectClientKeysWithResponse(
			ctx,
			m.Organization.ValueString(),
			m.Project.ValueString(),
			params,
		)
		if err != nil {
			diags.Append(diagutils.NewClientError("read", err))
			return nil, diags
		}
		if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			diags.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return nil, diags
		}

		allKeys = append(allKeys, *httpResp.JSON200...)

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	if !m.Name.IsNull() {
		for _, key := range allKeys {
			if key.Name == m.Name.ValueString() {
				return new(key), diags
			}
		}
		diags.Append(diagutils.NewNotFoundError("client key"))
		return nil, diags
	}

	switch {
	case len(allKeys) == 0:
		diags.Append(diagutils.NewNotFoundError("client key"))
		return nil, diags
	case len(allKeys) == 1:
		return new(allKeys[0]), diags
	case m.First.ValueBool():
		// Find the first key by date created
		sort.Slice(allKeys, func(i, j int) bool {
			return allKeys[i].DateCreated.Before(allKeys[j].DateCreated)
		})
		return new(allKeys[0]), diags
	default:
		diags.AddError("Client error", "Multiple keys found, please specify the key by `name`, `id`, or set the `first` flag to `true`.")
		return nil, diags
	}
}
