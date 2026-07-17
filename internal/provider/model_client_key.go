package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
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
