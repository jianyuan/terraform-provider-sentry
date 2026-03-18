package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/oapi-codegen/nullable"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *CronMonitorResource) getCreateJSONRequestBody(ctx context.Context, data CronMonitorResourceModel) (*apiclient.CreateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	var outDs apiclient.ProjectMonitorDataSourceConfigCron

	inSchedule := data.Schedule.DiagsGet(ctx, diags)
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
		WorkflowIds: []string{},
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
		defaultAssignee := data.DefaultAssignee.MustGet(ctx)
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
	if err := req.FromProjectMonitorRequestMonitorCheckInFailure(out); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}
	return &req, nil
}

func (m *CronMonitorResourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringValue(data.Name)
	if v, err := data.Description.Get(); err == nil {
		m.Description = supertypes.NewStringValueOrNull(v)
	} else {
		m.Description = supertypes.NewStringNull()
	}
	m.Enabled = supertypes.NewBoolValue(data.Enabled)

	if data.Owner.IsSpecified() && !data.Owner.IsNull() {
		ownerValue, err := data.Owner.MustGet().ValueByDiscriminator()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		defaultAssignee := &CronMonitorResourceModelDefaultAssignee{}

		switch ownerValue := ownerValue.(type) {
		case apiclient.ProjectMonitorOwnerUser:
			defaultAssignee.UserId = supertypes.NewStringValue(ownerValue.Id)
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOf(ctx, defaultAssignee)
		case apiclient.ProjectMonitorOwnerTeam:
			defaultAssignee.TeamId = supertypes.NewStringValue(ownerValue.Id)
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOf(ctx, defaultAssignee)
		default:
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOfNull[CronMonitorResourceModelDefaultAssignee](ctx)
		}
	} else {
		m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOfNull[CronMonitorResourceModelDefaultAssignee](ctx)
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
		m.CheckinMargin = supertypes.NewInt64Value(configValue.CheckinMargin)
		m.FailureIssueThreshold = supertypes.NewInt64Value(configValue.FailureIssueThreshold)
		m.MaxRuntime = supertypes.NewInt64Value(configValue.MaxRuntime)
		m.RecoveryThreshold = supertypes.NewInt64Value(configValue.RecoveryThreshold)
		m.Timezone = supertypes.NewStringValue(configValue.Timezone)
		schedule.Crontab = supertypes.NewStringValue(configValue.Schedule)

	case apiclient.ProjectMonitorDataSourceConfigCronInterval:
		m.CheckinMargin = supertypes.NewInt64Value(configValue.CheckinMargin)
		m.FailureIssueThreshold = supertypes.NewInt64Value(configValue.FailureIssueThreshold)
		m.MaxRuntime = supertypes.NewInt64Value(configValue.MaxRuntime)
		m.RecoveryThreshold = supertypes.NewInt64Value(configValue.RecoveryThreshold)
		m.Timezone = supertypes.NewStringValue(configValue.Timezone)

		if len(configValue.Schedule) == 2 {
			if intervalValue, err := configValue.Schedule[0].AsProjectMonitorDataSourceConfigCronIntervalValue(); err == nil {
				schedule.IntervalValue = supertypes.NewInt64Value(intervalValue)
			} else {
				diags.AddError("Invalid schedule", "Invalid schedule")
			}

			if intervalUnit, err := configValue.Schedule[1].AsProjectMonitorDataSourceConfigCronIntervalUnit(); err == nil {
				schedule.IntervalUnit = supertypes.NewStringValue(string(intervalUnit))
			} else {
				diags.AddError("Invalid schedule", "Invalid schedule")
			}
		} else {
			diags.AddError("Invalid schedule", fmt.Sprintf("Expected 2 items, got %d", len(configValue.Schedule)))
		}
	}

	m.Schedule = supertypes.NewSingleNestedObjectValueOf(ctx, schedule)

	return
}
