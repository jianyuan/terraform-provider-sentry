package client_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type JavascriptLoaderScriptModel struct {
	BrowserSdkVersion            types.String `tfsdk:"browser_sdk_version"`
	PerformanceMonitoringEnabled types.Bool   `tfsdk:"performance_monitoring_enabled"`
	SessionReplayEnabled         types.Bool   `tfsdk:"session_replay_enabled"`
	DebugEnabled                 types.Bool   `tfsdk:"debug_enabled"`
}

func (m *JavascriptLoaderScriptModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
	m.BrowserSdkVersion = types.StringValue(key.BrowserSdkVersion)
	m.PerformanceMonitoringEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasPerformance)
	m.SessionReplayEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasReplay)
	m.DebugEnabled = types.BoolValue(key.DynamicSdkLoaderOptions.HasDebug)
	return
}

type ResourceModel struct {
	Id                     types.String                                                      `tfsdk:"id"`
	Organization           types.String                                                      `tfsdk:"organization"`
	Project                types.String                                                      `tfsdk:"project"`
	ProjectId              types.String                                                      `tfsdk:"project_id"`
	Name                   types.String                                                      `tfsdk:"name"`
	RateLimitWindow        types.Int64                                                       `tfsdk:"rate_limit_window"`
	RateLimitCount         types.Int64                                                       `tfsdk:"rate_limit_count"`
	JavascriptLoaderScript supertypes.SingleNestedObjectValueOf[JavascriptLoaderScriptModel] `tfsdk:"javascript_loader_script"`
	Public                 types.String                                                      `tfsdk:"public"`
	Secret                 types.String                                                      `tfsdk:"secret"`
	Dsn                    supertypes.MapValueOf[string]                                     `tfsdk:"dsn"`
	DsnPublic              types.String                                                      `tfsdk:"dsn_public"`
	DsnSecret              types.String                                                      `tfsdk:"dsn_secret"`
	DsnCsp                 types.String                                                      `tfsdk:"dsn_csp"`
}

func (m *ResourceModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
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

	var javascriptLoaderScript JavascriptLoaderScriptModel
	diags.Append(javascriptLoaderScript.Fill(ctx, key)...)
	if diags.HasError() {
		return
	}

	m.JavascriptLoaderScript = supertypes.NewSingleNestedObjectValueOfNull[JavascriptLoaderScriptModel](ctx)
	diags.Append(m.JavascriptLoaderScript.Set(ctx, &javascriptLoaderScript)...)
	if diags.HasError() {
		return
	}

	m.Public = types.StringValue(key.Public)
	m.Secret = types.StringValue(key.Secret)
	m.Dsn = tfutils.MergeDiagnostics(supertypes.NewMapValueOfMap(ctx, key.Dsn))(&diags)

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

type DataSourceModel struct {
	Organization           types.String                                                      `tfsdk:"organization"`
	Project                types.String                                                      `tfsdk:"project"`
	Id                     types.String                                                      `tfsdk:"id"`
	Name                   types.String                                                      `tfsdk:"name"`
	First                  types.Bool                                                        `tfsdk:"first"`
	ProjectId              types.String                                                      `tfsdk:"project_id"`
	RateLimitWindow        types.Int64                                                       `tfsdk:"rate_limit_window"`
	RateLimitCount         types.Int64                                                       `tfsdk:"rate_limit_count"`
	JavascriptLoaderScript supertypes.SingleNestedObjectValueOf[JavascriptLoaderScriptModel] `tfsdk:"javascript_loader_script"`
	Public                 types.String                                                      `tfsdk:"public"`
	Secret                 types.String                                                      `tfsdk:"secret"`
	Dsn                    supertypes.MapValueOf[string]                                     `tfsdk:"dsn"`
	DsnPublic              types.String                                                      `tfsdk:"dsn_public"`
	DsnSecret              types.String                                                      `tfsdk:"dsn_secret"`
	DsnCsp                 types.String                                                      `tfsdk:"dsn_csp"`
}

func (m *DataSourceModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
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

	var javascriptLoaderScript JavascriptLoaderScriptModel
	diags.Append(javascriptLoaderScript.Fill(ctx, key)...)
	if diags.HasError() {
		return
	}

	m.JavascriptLoaderScript = supertypes.NewSingleNestedObjectValueOfNull[JavascriptLoaderScriptModel](ctx)
	diags.Append(m.JavascriptLoaderScript.Set(ctx, &javascriptLoaderScript)...)
	if diags.HasError() {
		return
	}

	m.Public = types.StringValue(key.Public)
	m.Secret = types.StringValue(key.Secret)
	m.Dsn = tfutils.MergeDiagnostics(supertypes.NewMapValueOfMap(ctx, key.Dsn))(&diags)

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
