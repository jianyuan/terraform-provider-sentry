package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func (r *CronMonitorResource) getCreateJSONRequestBody(ctx context.Context, data CronMonitorResourceModel) (*apiclient.CreateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	var outDs apiclient.ProjectMonitorDataSourceConfigCron

	inSchedule := tfutils.MergeDiagnostics(data.Schedule.Get(ctx))(&diags)
	if diags.HasError() {
		return nil, diags
	}

	switch {
	case inSchedule.Crontab.IsKnown():
		outDs.FromProjectMonitorDataSourceConfigCronCrontab(apiclient.ProjectMonitorDataSourceConfigCronCrontab{
			CheckinMargin:         data.CheckinMargin.Get(),
			FailureIssueThreshold: data.FailureIssueThreshold.Get(),
			MaxRuntime:            data.MaxRuntime.Get(),
			RecoveryThreshold:     data.RecoveryThreshold.Get(),
			Timezone:              data.Timezone.Get(),
			Schedule:              inSchedule.Crontab.Get(),
		})
	case inSchedule.IntervalUnit.IsKnown() && inSchedule.IntervalValue.IsKnown():
		var intervalValue apiclient.ProjectMonitorDataSourceConfigCronInterval_Schedule_Item
		intervalValue.FromProjectMonitorDataSourceConfigCronIntervalValue(inSchedule.IntervalValue.Get())
		var intervalUnit apiclient.ProjectMonitorDataSourceConfigCronInterval_Schedule_Item
		intervalUnit.FromProjectMonitorDataSourceConfigCronIntervalUnit(apiclient.ProjectMonitorDataSourceConfigCronIntervalUnit(inSchedule.IntervalUnit.Get()))

		outDs.FromProjectMonitorDataSourceConfigCronInterval(apiclient.ProjectMonitorDataSourceConfigCronInterval{
			CheckinMargin:         data.CheckinMargin.Get(),
			FailureIssueThreshold: data.FailureIssueThreshold.Get(),
			MaxRuntime:            data.MaxRuntime.Get(),
			RecoveryThreshold:     data.RecoveryThreshold.Get(),
			Timezone:              data.Timezone.Get(),
			Schedule: []apiclient.ProjectMonitorDataSourceConfigCronInterval_Schedule_Item{
				intervalValue,
				intervalUnit,
			},
		})
	}

	out := apiclient.ProjectMonitorRequestMonitorCheckInFailure{
		Name:      data.Name.Get(),
		ProjectId: data.Project.Get(),
		DataSources: []apiclient.ProjectMonitorDataSourceCron{
			{
				Name:   data.Name.Get(),
				Config: outDs,
			},
		},
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
	if err := req.FromProjectMonitorRequestMonitorCheckInFailure(out); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}
	return &req, nil
}

func (r *CronMonitorResource) getUpdateJSONRequestBody(ctx context.Context, data CronMonitorResourceModel) (*apiclient.UpdateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	return r.getCreateJSONRequestBody(ctx, data)
}

func (m *CronMonitorResourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
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

		outOwner := &CronMonitorResourceModelOwner{}

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

	schedule := &CronMonitorResourceModelSchedule{}

	switch configValue := configValue.(type) {
	case apiclient.ProjectMonitorDataSourceConfigCronCrontab:
		m.CheckinMargin.Set(configValue.CheckinMargin)
		m.FailureIssueThreshold.Set(configValue.FailureIssueThreshold)
		m.MaxRuntime.Set(configValue.MaxRuntime)
		m.RecoveryThreshold.Set(configValue.RecoveryThreshold)
		m.Timezone.Set(configValue.Timezone)
		schedule.Crontab.Set(configValue.Schedule)

	case apiclient.ProjectMonitorDataSourceConfigCronInterval:
		m.CheckinMargin.Set(configValue.CheckinMargin)
		m.FailureIssueThreshold.Set(configValue.FailureIssueThreshold)
		m.MaxRuntime.Set(configValue.MaxRuntime)
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
