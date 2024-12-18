package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"

	// "github.com/jianyuan/terraform-provider-sentry/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitorResourceModelToMonitorRequestScheduleCrontab(t *testing.T) {
	config := MonitorConfigResourceModel{
		ScheduleCrontab:       types.StringValue("* * * * *"),
		ScheduleInterval:      types.ObjectNull((*MonitorConfigScheduleIntervalResourceModel)(nil).AttributeTypes()),
		Timezone:              types.StringValue("UTC"),
		CheckinMargin:         types.Int64Value(10),
		MaxRuntime:            types.Int64Value(20),
		FailureIssueThreshold: types.Int64Value(5),
		RecoveryThreshold:     types.Int64Value(10),
		AlertRuleId:           types.Int64Value(123),
	}

	configObject, configObjectDiags := types.ObjectValueFrom(context.Background(), (*MonitorConfigResourceModel)(nil).AttributeTypes(), config)
	require.Empty(t, configObjectDiags)

	model := MonitorResourceModel{
		Organization: types.StringValue("sentry-org"),
		Id:           types.StringValue("monitor_id"),
		Project:      types.StringValue("sentry-project"),
		Name:         types.StringValue("monitor name"),
		Slug:         types.StringValue("monitor-slug"),
		Owner:        types.StringValue("team:123"),
		IsMuted:      types.BoolValue(false),
		Status:       types.StringValue("active"),
		Config:       configObject,
	}

	monitorRequest, monitorRequestDiags := model.ToMonitorRequest(context.Background())

	require.Empty(t, monitorRequestDiags)

	assert.Equal(t, "monitor name", monitorRequest.Name)
	assert.Equal(t, "monitor-slug", *monitorRequest.Slug)
	assert.Equal(t, "sentry-project", monitorRequest.Project)
	assert.Equal(t, "team:123", *monitorRequest.Owner)
	assert.Equal(t, apiclient.MonitorRequestStatusActive, *monitorRequest.Status)
	assert.Equal(t, false, *monitorRequest.IsMuted)

	assert.Equal(t, apiclient.Crontab, monitorRequest.Config.ScheduleType)

	monitorRequestConfigScheduleCrontab, monitorRequestConfigScheduleCrontabErr := monitorRequest.Config.Schedule.AsMonitorConfigScheduleString()
	assert.NoError(t, monitorRequestConfigScheduleCrontabErr)
	assert.Equal(t, "* * * * *", monitorRequestConfigScheduleCrontab)

	_, monitorRequestConfigScheduleIntervalErr := monitorRequest.Config.Schedule.AsMonitorConfigScheduleInterval()
	assert.Error(t, monitorRequestConfigScheduleIntervalErr)

	assert.Equal(t, "UTC", monitorRequest.Config.Timezone)
	assert.Equal(t, int64(10), *monitorRequest.Config.CheckinMargin)
	assert.Equal(t, int64(20), *monitorRequest.Config.MaxRuntime)
	assert.Equal(t, int64(5), *monitorRequest.Config.FailureIssueThreshold)
	assert.Equal(t, int64(10), *monitorRequest.Config.RecoveryThreshold)
	assert.Equal(t, int64(123), *monitorRequest.Config.AlertRuleId)
}

func TestMonitorResourceModelToMonitorRequestScheduleInterval(t *testing.T) {
	scheduleInterval := MonitorConfigScheduleIntervalResourceModel{
		Day: types.Int64Value(1),
	}

	scheduleIntervalObject, scheduleIntervalObjectDiags := types.ObjectValueFrom(context.Background(), (*MonitorConfigScheduleIntervalResourceModel)(nil).AttributeTypes(), scheduleInterval)
	require.Empty(t, scheduleIntervalObjectDiags)

	config := MonitorConfigResourceModel{
		ScheduleCrontab:       types.StringNull(),
		ScheduleInterval:      scheduleIntervalObject,
		Timezone:              types.StringValue("UTC"),
		CheckinMargin:         types.Int64Value(10),
		MaxRuntime:            types.Int64Value(20),
		FailureIssueThreshold: types.Int64Value(5),
		RecoveryThreshold:     types.Int64Value(10),
		AlertRuleId:           types.Int64Value(123),
	}

	configObject, configObjectDiags := types.ObjectValueFrom(context.Background(), (*MonitorConfigResourceModel)(nil).AttributeTypes(), config)
	require.Empty(t, configObjectDiags)

	model := MonitorResourceModel{
		Organization: types.StringValue("sentry-org"),
		Id:           types.StringValue("monitor_id"),
		Project:      types.StringValue("sentry-project"),
		Name:         types.StringValue("monitor name"),
		Slug:         types.StringValue("monitor-slug"),
		Owner:        types.StringValue("team:123"),
		IsMuted:      types.BoolValue(false),
		Status:       types.StringValue("active"),
		Config:       configObject,
	}

	monitorRequest, monitorRequestDiags := model.ToMonitorRequest(context.Background())

	require.Empty(t, monitorRequestDiags)

	assert.Equal(t, "monitor name", monitorRequest.Name)
	assert.Equal(t, "monitor-slug", *monitorRequest.Slug)
	assert.Equal(t, "sentry-project", monitorRequest.Project)
	assert.Equal(t, "team:123", *monitorRequest.Owner)
	assert.Equal(t, apiclient.MonitorRequestStatusActive, *monitorRequest.Status)
	assert.Equal(t, false, *monitorRequest.IsMuted)

	assert.Equal(t, apiclient.Interval, monitorRequest.Config.ScheduleType)

	_, monitorRequestConfigScheduleCrontabErr := monitorRequest.Config.Schedule.AsMonitorConfigScheduleString()
	assert.Error(t, monitorRequestConfigScheduleCrontabErr)

	monitorRequestConfigScheduleInterval, monitorRequestConfigScheduleIntervalErr := monitorRequest.Config.Schedule.AsMonitorConfigScheduleInterval()
	assert.NoError(t, monitorRequestConfigScheduleIntervalErr)
	assert.Equal(t, []any{float64(1), "day"}, monitorRequestConfigScheduleInterval)

	assert.Equal(t, "UTC", monitorRequest.Config.Timezone)
	assert.Equal(t, int64(10), *monitorRequest.Config.CheckinMargin)
	assert.Equal(t, int64(20), *monitorRequest.Config.MaxRuntime)
	assert.Equal(t, int64(5), *monitorRequest.Config.FailureIssueThreshold)
	assert.Equal(t, int64(10), *monitorRequest.Config.RecoveryThreshold)
	assert.Equal(t, int64(123), *monitorRequest.Config.AlertRuleId)
}
