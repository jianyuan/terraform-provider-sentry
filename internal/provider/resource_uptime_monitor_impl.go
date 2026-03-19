package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
	"github.com/oapi-codegen/nullable"
)

func (r *UptimeMonitorResource) getCreateJSONRequestBody(ctx context.Context, data UptimeMonitorResourceModel) (*apiclient.CreateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	outDs := apiclient.ProjectMonitorDataSourceUptimeDomainFailure{
		Url:             data.Url.Get(),
		Method:          data.Method.Get(),
		Headers:         [][]string{},
		IntervalSeconds: data.IntervalSeconds.Get(),
		TimeoutMs:       data.TimeoutMs.Get(),
	}
	if data.Body.IsKnown() {
		outDs.Body = nullable.NewNullableWithValue(data.Body.Get())
	} else {
		outDs.Body = nullable.NewNullNullable[string]()
	}
	if data.Headers.IsKnown() {
		inHeaders := tfutils.MergeDiagnostics(data.Headers.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}
		for _, inHeader := range inHeaders {
			outDs.Headers = append(outDs.Headers, []string{inHeader.Key.Get(), inHeader.Value.Get()})
		}
	}

	var outConfig apiclient.ProjectMonitorConfig
	if err := outConfig.FromProjectMonitorConfigUptimeDomainFailure(apiclient.ProjectMonitorConfigUptimeDomainFailure{
		Mode:              apiclient.ProjectMonitorConfigUptimeDomainFailureMode(sentrydata.UptimeMonitorModeNameToId["MANUAL"]),
		Environment:       data.Environment.Get(),
		RecoveryThreshold: data.RecoveryThreshold.Get(),
		DowntimeThreshold: data.DowntimeThreshold.Get(),
	}); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}

	out := apiclient.ProjectMonitorRequestUptimeDomainFailure{
		Name:      data.Name.Get(),
		ProjectId: data.Project.Get(),
		DataSources: []apiclient.ProjectMonitorDataSourceUptimeDomainFailure{
			outDs,
		},
		Config: &outConfig,
	}

	if data.Enabled.IsKnown() {
		out.Enabled = nullable.NewNullableWithValue(data.Enabled.Get())
	} else {
		out.Enabled = nullable.NewNullNullable[bool]()
	}

	if data.Description.IsKnown() {
		out.Description = nullable.NewNullableWithValue(data.Description.Get())
	} else {
		out.Description = nullable.NewNullNullable[string]()
	}

	if data.DefaultAssignee.IsKnown() {
		defaultAssignee := tfutils.MergeDiagnostics(data.DefaultAssignee.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}

		switch {
		case defaultAssignee.TeamId.IsKnown():
			out.Owner = nullable.NewNullableWithValue(fmt.Sprintf("team:%s", defaultAssignee.TeamId.Get()))
		case defaultAssignee.UserId.IsKnown():
			out.Owner = nullable.NewNullableWithValue(fmt.Sprintf("user:%s", defaultAssignee.UserId.Get()))
		default:
			out.Owner = nullable.NewNullNullable[string]()
		}
	} else {
		out.Owner = nullable.NewNullNullable[string]()
	}

	var req apiclient.CreateProjectMonitorJSONRequestBody
	if err := req.FromProjectMonitorRequestUptimeDomainFailure(out); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}
	return &req, nil
}

func (r *UptimeMonitorResource) getUpdateJSONRequestBody(ctx context.Context, data UptimeMonitorResourceModel) (*apiclient.UpdateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	return r.getCreateJSONRequestBody(ctx, data)
}

func (m *UptimeMonitorResourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	m.Id.Set(data.Id)
	m.Name.Set(data.Name)
	if v, err := data.Description.Get(); err == nil {
		m.Description.Set(v)
	} else {
		m.Description.SetNull()
	}
	m.Enabled.Set(data.Enabled)

	if data.Owner.IsSpecified() && !data.Owner.IsNull() {
		ownerValue, err := data.Owner.MustGet().ValueByDiscriminator()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		defaultAssignee := &UptimeMonitorResourceModelDefaultAssignee{}

		switch ownerValue := ownerValue.(type) {
		case apiclient.ProjectMonitorOwnerUser:
			defaultAssignee.UserId.Set(ownerValue.Id)
			diags.Append(m.DefaultAssignee.Set(ctx, defaultAssignee)...)
		case apiclient.ProjectMonitorOwnerTeam:
			defaultAssignee.TeamId.Set(ownerValue.Id)
			diags.Append(m.DefaultAssignee.Set(ctx, defaultAssignee)...)
		default:
			m.DefaultAssignee.SetNull(ctx)
		}
	} else {
		m.DefaultAssignee.SetNull(ctx)
	}

	if len(data.DataSources) != 1 {
		diags.AddError("Invalid data source", fmt.Sprintf("Expected 1 data source, got %d", len(data.DataSources)))
		return
	}

	dataSource, err := data.DataSources[0].AsProjectMonitorDataSourceWrapperUptimeSubscription()
	if err != nil {
		diags.AddError("Invalid data source", err.Error())
		return
	}

	m.Url.Set(dataSource.QueryObj.Url)
	m.Method.Set(dataSource.QueryObj.Method)
	if v, err := dataSource.QueryObj.Body.Get(); err == nil {
		m.Body.Set(v)
	} else {
		m.Body.SetNull()
	}

	headers := make([]*UptimeMonitorResourceModelHeadersItem, 0, len(dataSource.QueryObj.Headers))
	for _, headerValues := range dataSource.QueryObj.Headers {
		if len(headerValues) != 2 {
			diags.AddError("Invalid header", fmt.Sprintf("Expected 2 elements in header, got %d", len(headerValues)))
			return
		}
		var header UptimeMonitorResourceModelHeadersItem
		header.Key.Set(headerValues[0])
		header.Value.Set(headerValues[1])
		headers = append(headers, &header)
	}
	diags.Append(m.Headers.Set(ctx, headers)...)

	m.IntervalSeconds.Set(dataSource.QueryObj.IntervalSeconds)
	m.TimeoutMs.Set(dataSource.QueryObj.TimeoutMs)

	if config, err := data.Config.AsProjectMonitorConfigUptimeDomainFailure(); err == nil {
		m.Environment.Set(config.Environment)
		m.RecoveryThreshold.Set(config.RecoveryThreshold)
		m.DowntimeThreshold.Set(config.DowntimeThreshold)
	} else {
		diags.AddError("Invalid config", err.Error())
		return
	}

	return
}
