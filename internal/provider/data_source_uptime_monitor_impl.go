package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func (m *UptimeMonitorDataSourceModel) fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
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

		outOwner := &UptimeMonitorDataSourceModelOwner{}

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
		m.Body.Set(v)
	} else {
		m.Body.SetNull()
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
		m.AssertionJson.Set(string(assertion))
	} else {
		m.AssertionJson.SetNull()
	}

	return
}
