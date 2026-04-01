package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
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
	if !data.Body.IsNull() && !data.Body.IsUnknown() {
		outDs.Body.Set(data.Body.ValueString())
	} else {
		outDs.Body.SetNull()
	}
	if data.Headers.IsKnown() {
		inHeaders := tfutils.MergeDiagnostics(data.Headers.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}
		for key, value := range inHeaders {
			outDs.Headers = append(outDs.Headers, []string{key, value})
		}
	}
	if !data.AssertionJson.IsNull() && !data.AssertionJson.IsUnknown() {
		outDs.Assertion.Set(json.RawMessage(data.AssertionJson.ValueString()))
	} else {
		outDs.Assertion.SetNull()
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
		out.Enabled.Set(data.Enabled.Get())
	} else {
		out.Enabled.SetNull()
	}

	if data.Description.IsKnown() {
		out.Description.Set(data.Description.Get())
	} else {
		out.Description.SetNull()
	}

	if data.Owner.IsKnown() {
		owner := tfutils.MergeDiagnostics(data.Owner.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}

		switch {
		case owner.TeamId.IsKnown():
			out.Owner.Set(fmt.Sprintf("team:%s", owner.TeamId.Get()))
		case owner.UserId.IsKnown():
			out.Owner.Set(fmt.Sprintf("user:%s", owner.UserId.Get()))
		default:
			out.Owner.SetNull()
		}
	} else {
		out.Owner.SetNull()
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
		inOwner, err := data.Owner.Get()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		inOwnerValue, err := inOwner.ValueByDiscriminator()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		outOwner := &UptimeMonitorResourceModelOwner{}

		switch inOwnerValue := inOwnerValue.(type) {
		case apiclient.ProjectMonitorOwnerUser:
			outOwner.UserId.Set(inOwnerValue.Id)
			diags.Append(m.Owner.Set(ctx, outOwner)...)
		case apiclient.ProjectMonitorOwnerTeam:
			outOwner.TeamId.Set(inOwnerValue.Id)
			diags.Append(m.Owner.Set(ctx, outOwner)...)
		default:
			m.Owner.SetNull(ctx)
		}
	} else {
		m.Owner.SetNull(ctx)
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
		m.Body = sentrytypes.TrimmedStringValue(v)
	} else {
		m.Body = sentrytypes.TrimmedStringNull()
	}

	headers := make(map[string]string, len(dataSource.QueryObj.Headers))
	for _, headerValues := range dataSource.QueryObj.Headers {
		if len(headerValues) != 2 {
			diags.AddError("Invalid header", fmt.Sprintf("Expected 2 elements in header, got %d", len(headerValues)))
			return
		}
		headers[headerValues[0]] = headerValues[1]
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

	if dataSource.QueryObj.Assertion.IsSpecified() && !dataSource.QueryObj.Assertion.IsNull() {
		assertion, err := dataSource.QueryObj.Assertion.Get()
		if err != nil {
			diags.AddError("Invalid assertion", err.Error())
			return
		}
		m.AssertionJson = jsontypes.NewNormalizedValue(string(assertion))
	} else {
		m.AssertionJson = jsontypes.NewNormalizedNull()
	}

	return
}
