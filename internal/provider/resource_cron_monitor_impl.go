package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/oapi-codegen/nullable"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *CronMonitorResource) getCreateJSONRequestBody(ctx context.Context, data CronMonitorResourceModel) (apiclient.CreateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	dataSourceConfig := apiclient.ProjectMonitorDataSourceConfigCron{
		CheckinMargin:         data.CheckinMargin.GetInt64(),
		FailureIssueThreshold: data.FailureIssueThreshold.GetInt64(),
		MaxRuntime:            data.MaxRuntime.GetInt64(),
		RecoveryThreshold:     data.RecoveryThreshold.GetInt64(),
		Timezone:              data.Timezone.Get(),
	}

	// dataSourceConfig.MergeProjectMonitorDataSourceConfigCronCrontab(apiclient.ProjectMonitorDataSourceConfigCronCrontab{
	// 	ScheduleType: apiclient.Crontab,
	// 	Schedule:     "0 0 * * *",
	// })

	var scheduleIntervalValue apiclient.ProjectMonitorDataSourceConfigCronInterval_Schedule_Item
	scheduleIntervalValue.FromProjectMonitorDataSourceConfigCronIntervalValue(1)
	var scheduleIntervalUnit apiclient.ProjectMonitorDataSourceConfigCronInterval_Schedule_Item
	scheduleIntervalUnit.FromProjectMonitorDataSourceConfigCronIntervalUnit(apiclient.Day)

	dataSourceConfig.MergeProjectMonitorDataSourceConfigCronInterval(apiclient.ProjectMonitorDataSourceConfigCronInterval{
		ScheduleType: apiclient.Interval,
		Schedule:     []apiclient.ProjectMonitorDataSourceConfigCronInterval_Schedule_Item{scheduleIntervalValue, scheduleIntervalUnit},
	})

	out := apiclient.CreateProjectMonitorJSONRequestBody{
		Type:      apiclient.MonitorCheckInFailure,
		Name:      data.Name.ValueString(),
		ProjectId: data.Project.ValueString(),
		DataSources: []apiclient.ProjectMonitorDataSource{
			{
				Name:   data.Name.ValueString(),
				Config: dataSourceConfig,
			},
		},
		WorkflowIds: []string{},
	}

	if data.Description.IsNull() {
		out.Description = nullable.NewNullNullable[string]()
	} else {
		out.Description = nullable.NewNullableWithValue(data.Description.Get())
	}

	return out, nil
}

func (m *CronMonitorResourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringValue(data.Name)
	if v, err := data.Description.Get(); err == nil {
		m.Description = supertypes.NewStringValueOrNull(v)
	} else {
		m.Description = supertypes.NewStringNull()
	}

	if len(data.DataSources) == 1 {
		m.CheckinMargin = supertypes.NewInt64Value(data.DataSources[0].QueryObj.Config.CheckinMargin)
		m.FailureIssueThreshold = supertypes.NewInt64Value(data.DataSources[0].QueryObj.Config.FailureIssueThreshold)
		m.MaxRuntime = supertypes.NewInt64Value(data.DataSources[0].QueryObj.Config.MaxRuntime)
		m.RecoveryThreshold = supertypes.NewInt64Value(data.DataSources[0].QueryObj.Config.RecoveryThreshold)
		m.Timezone = supertypes.NewStringValue(data.DataSources[0].QueryObj.Config.Timezone)
	} else {
		m.CheckinMargin = supertypes.NewInt64Null()
		m.FailureIssueThreshold = supertypes.NewInt64Null()
		m.MaxRuntime = supertypes.NewInt64Null()
		m.RecoveryThreshold = supertypes.NewInt64Null()
		m.Timezone = supertypes.NewStringNull()
		// TODO
		// diags.AddError("")
	}

	return
}
