package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func (m *CronMonitorDataSourceModel) fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
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

		outOwner := &CronMonitorDataSourceModelOwner{}

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

	if diags.HasError() {
		return
	}

	if len(data.DataSources) != 1 {
		diags.AddError("Invalid data source", fmt.Sprintf("Expected 1 data source, got %d", len(data.DataSources)))
		return
	}

	dataSource, err := data.DataSources[0].AsProjectMonitorDataSourceWrapperCronMonitor()
	if err != nil {
		diags.AddError("Invalid config", err.Error())
		return
	}

	configValue, err := dataSource.QueryObj.Config.ValueByDiscriminator()
	if err != nil {
		diags.AddError("Invalid config", err.Error())
		return
	}

	schedule := &CronMonitorDataSourceModelSchedule{}

	switch configValue := configValue.(type) {
	case apiclient.ProjectMonitorDataSourceConfigCronCrontab:
		m.CheckinMarginMinutes.Set(configValue.CheckinMargin)
		m.FailureIssueThreshold.Set(configValue.FailureIssueThreshold)
		m.MaxRuntimeMinutes.Set(configValue.MaxRuntime)
		m.RecoveryThreshold.Set(configValue.RecoveryThreshold)
		m.Timezone.Set(configValue.Timezone)
		schedule.Crontab.Set(configValue.Schedule)

	case apiclient.ProjectMonitorDataSourceConfigCronInterval:
		m.CheckinMarginMinutes.Set(configValue.CheckinMargin)
		m.FailureIssueThreshold.Set(configValue.FailureIssueThreshold)
		m.MaxRuntimeMinutes.Set(configValue.MaxRuntime)
		m.RecoveryThreshold.Set(configValue.RecoveryThreshold)
		m.Timezone.Set(configValue.Timezone)

		if len(configValue.Schedule) != 2 {
			diags.AddError("Invalid schedule", fmt.Sprintf("Expected 2 items, got %d", len(configValue.Schedule)))

			return
		}

		if intervalValue, err := configValue.Schedule[0].AsProjectMonitorDataSourceConfigCronIntervalValue(); err == nil {
			schedule.IntervalValue.Set(intervalValue)
		} else {
			diags.AddError("Invalid schedule", "Invalid interval value")
			return
		}

		if intervalUnit, err := configValue.Schedule[1].AsProjectMonitorDataSourceConfigCronIntervalUnit(); err == nil {
			schedule.IntervalUnit.Set(intervalUnit)
		} else {
			diags.AddError("Invalid schedule", "Invalid interval unit")
			return
		}
	}

	diags.Append(m.Schedule.Set(ctx, schedule)...)

	return
}
