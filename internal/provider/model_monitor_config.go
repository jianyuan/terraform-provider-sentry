package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type MonitorConfigResourceModel struct {
	ScheduleCrontab       types.String `tfsdk:"schedule_crontab"`
	ScheduleInterval      types.Object `tfsdk:"schedule_interval"`
	CheckinMargin         types.Int64  `tfsdk:"checkin_margin"`
	MaxRuntime            types.Int64  `tfsdk:"max_runtime"`
	Timezone              types.String `tfsdk:"timezone"`
	FailureIssueThreshold types.Int64  `tfsdk:"failure_issue_threshold"`
	RecoveryThreshold     types.Int64  `tfsdk:"recovery_threshold"`
	AlertRuleId           types.Int64  `tfsdk:"alert_rule_id"`
}

func (m MonitorConfigResourceModel) SchemaAttribute(required bool) schema.Attribute {
	return schema.SingleNestedAttribute{
		Required:   required,
		Attributes: m.SchemaAttributes(),
	}
}

func (m MonitorConfigResourceModel) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"schedule_crontab": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("schedule_interval")),
			},
		},
		"schedule_interval": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: MonitorConfigScheduleIntervalResourceModel{}.SchemaAttributes(),
			Validators: []validator.Object{
				objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("schedule_crontab")),
			},
		},
		"checkin_margin": schema.Int64Attribute{
			Optional: true,
		},
		"max_runtime": schema.Int64Attribute{
			Optional: true,
		},
		"timezone": schema.StringAttribute{
			Optional: true,
		},
		"failure_issue_threshold": schema.Int64Attribute{
			Optional: true,
		},
		"recovery_threshold": schema.Int64Attribute{
			Optional: true,
		},
		"alert_rule_id": schema.Int64Attribute{
			Optional: true,
		},
	}
}

func (m *MonitorConfigResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"schedule_crontab":        types.StringType,
		"schedule_interval":       types.ObjectType{AttrTypes: (&MonitorConfigScheduleIntervalResourceModel{}).AttributeTypes()},
		"checkin_margin":          types.Int64Type,
		"max_runtime":             types.Int64Type,
		"timezone":                types.StringType,
		"failure_issue_threshold": types.Int64Type,
		"recovery_threshold":      types.Int64Type,
		"alert_rule_id":           types.Int64Type,
	}
}

func (m *MonitorConfigResourceModel) ToMonitorRequest(ctx context.Context, path path.Path) (apiclient.MonitorConfig, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	var scheduleType apiclient.MonitorConfigScheduleType
	var configSchedule apiclient.MonitorConfig_Schedule

	if !m.ScheduleCrontab.IsUnknown() && !m.ScheduleCrontab.IsNull() {
		scheduleType = apiclient.Crontab

		configSchedule.FromMonitorConfigScheduleString(m.ScheduleCrontab.ValueString())
	} else if !m.ScheduleInterval.IsUnknown() && !m.ScheduleInterval.IsNull() {
		scheduleType = apiclient.Interval

		var scheduleIntervalModel MonitorConfigScheduleIntervalResourceModel
		diags.Append(m.ScheduleInterval.As(ctx, &scheduleIntervalModel, basetypes.ObjectAsOptions{})...)

		scheduleInterval, scheduleIntervalDiags := formatMonitorConfigScheduleInterval(scheduleIntervalModel)
		diags.Append(scheduleIntervalDiags...)

		configSchedule.FromMonitorConfigScheduleInterval(scheduleInterval)
	}

	return apiclient.MonitorConfig{
		ScheduleType:          scheduleType,
		Schedule:              configSchedule,
		Timezone:              m.Timezone.ValueString(),
		CheckinMargin:         m.CheckinMargin.ValueInt64Pointer(),
		MaxRuntime:            m.MaxRuntime.ValueInt64Pointer(),
		FailureIssueThreshold: m.FailureIssueThreshold.ValueInt64Pointer(),
		RecoveryThreshold:     m.RecoveryThreshold.ValueInt64Pointer(),
		AlertRuleId:           m.AlertRuleId.ValueInt64Pointer(),
	}, diags
}

func (m *MonitorConfigResourceModel) Fill(ctx context.Context, path path.Path, config apiclient.MonitorConfig) (diags diag.Diagnostics) {
	switch config.ScheduleType {
	case apiclient.Crontab:
		schedule, scheduleErr := config.Schedule.AsMonitorConfigScheduleString()
		if scheduleErr != nil {
			diags.AddAttributeError(path.AtName("schedule"), "Invalid schedule", scheduleErr.Error())
			break
		}
		m.ScheduleCrontab = types.StringValue(schedule)
		m.ScheduleInterval = types.ObjectNull((&MonitorConfigScheduleIntervalResourceModel{}).AttributeTypes())
	case apiclient.Interval:
		schedule, scheduleErr := config.Schedule.AsMonitorConfigScheduleInterval()
		if scheduleErr != nil {
			diags.AddAttributeError(path.AtName("schedule"), "Invalid schedule", scheduleErr.Error())
			break
		}
		parsedSchedule := tfutils.MergeDiagnostics(parseMonitorConfigScheduleInterval(schedule))(&diags)
		m.ScheduleCrontab = types.StringNull()
		m.ScheduleInterval = tfutils.MergeDiagnostics(types.ObjectValueFrom(ctx, (&MonitorConfigScheduleIntervalResourceModel{}).AttributeTypes(), parsedSchedule))(&diags)
	default:
		diags.AddAttributeError(path.AtName("schedule"), "Invalid schedule type", string(config.ScheduleType))
	}

	m.CheckinMargin = types.Int64PointerValue(config.CheckinMargin)
	m.MaxRuntime = types.Int64PointerValue(config.MaxRuntime)
	m.Timezone = types.StringValue(config.Timezone)
	m.FailureIssueThreshold = types.Int64PointerValue(config.FailureIssueThreshold)
	m.RecoveryThreshold = types.Int64PointerValue(config.RecoveryThreshold)

	m.AlertRuleId = types.Int64PointerValue(config.AlertRuleId)

	return
}
