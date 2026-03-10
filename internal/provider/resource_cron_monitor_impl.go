package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/oapi-codegen/nullable"
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
