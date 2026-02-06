package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"github.com/oapi-codegen/nullable"
)

type MonitorConfigResourceModel struct {
	ScheduleCrontab       types.String `tfsdk:"schedule_crontab"`
	ScheduleInterval      types.Object `tfsdk:"schedule_interval"`
	CheckinMargin         types.Int64  `tfsdk:"checkin_margin"`
	MaxRuntime            types.Int64  `tfsdk:"max_runtime"`
	Timezone              types.String `tfsdk:"timezone"`
	FailureIssueThreshold types.Int64  `tfsdk:"failure_issue_threshold"`
	RecoveryThreshold     types.Int64  `tfsdk:"recovery_threshold"`
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
				stringvalidator.LengthAtLeast(1),
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
			Computed: true,
			Validators: []validator.Int64{
				int64validator.Between(0, 40320),
			},
		},
		"max_runtime": schema.Int64Attribute{
			Optional: true,
			Computed: true,
			Validators: []validator.Int64{
				int64validator.Between(1, 40320),
			},
		},
		"timezone": schema.StringAttribute{
			Optional: true,
		},
		"failure_issue_threshold": schema.Int64Attribute{
			Optional: true,
			Computed: true,
			Validators: []validator.Int64{
				int64validator.Between(1, 720),
			},
		},
		"recovery_threshold": schema.Int64Attribute{
			Optional: true,
			Computed: true,
			Validators: []validator.Int64{
				int64validator.Between(1, 720),
			},
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
	}
}

func (m *MonitorConfigResourceModel) ToMonitorRequest(ctx context.Context, path path.Path) (apiclient.MonitorConfigRequest, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	var scheduleType apiclient.MonitorConfigScheduleType
	var configSchedule apiclient.MonitorConfigRequest_Schedule

	if !m.ScheduleCrontab.IsUnknown() && !m.ScheduleCrontab.IsNull() {
		scheduleType = apiclient.MonitorConfigScheduleTypeCrontab

		if err := configSchedule.FromMonitorConfigScheduleString(m.ScheduleCrontab.ValueString()); err != nil {
			diags.AddAttributeError(path.AtName("schedule_crontab"), "Invalid schedule", err.Error())
		}
	} else if !m.ScheduleInterval.IsUnknown() && !m.ScheduleInterval.IsNull() {
		scheduleType = apiclient.MonitorConfigScheduleTypeInterval

		var scheduleIntervalModel MonitorConfigScheduleIntervalResourceModel
		diags.Append(m.ScheduleInterval.As(ctx, &scheduleIntervalModel, basetypes.ObjectAsOptions{})...)

		scheduleInterval, scheduleIntervalDiags := formatMonitorConfigScheduleInterval(scheduleIntervalModel)
		diags.Append(scheduleIntervalDiags...)

		if scheduleInterval == nil {
			diags.AddAttributeError(path.AtName("schedule_interval"), "Missing schedule interval", "Exactly one of year, month, week, day, hour, or minute must be set.")
		} else if err := configSchedule.FromMonitorConfigScheduleInterval(scheduleInterval); err != nil {
			diags.AddAttributeError(path.AtName("schedule_interval"), "Invalid schedule", err.Error())
		}
	}

	var scheduleTypePtr *apiclient.MonitorConfigScheduleType
	if scheduleType != "" {
		scheduleTypePtr = &scheduleType
	}

	var timezone *string
	if !m.Timezone.IsNull() && !m.Timezone.IsUnknown() {
		value := m.Timezone.ValueString()
		timezone = &value
	}

	return apiclient.MonitorConfigRequest{
		ScheduleType:          scheduleTypePtr,
		Schedule:              configSchedule,
		Timezone:              timezone,
		CheckinMargin:         int64ToNullable(m.CheckinMargin),
		MaxRuntime:            int64ToNullable(m.MaxRuntime),
		FailureIssueThreshold: int64ToNullable(m.FailureIssueThreshold),
		RecoveryThreshold:     int64ToNullable(m.RecoveryThreshold),
	}, diags
}

func (m *MonitorConfigResourceModel) Fill(ctx context.Context, path path.Path, config apiclient.MonitorConfig) (diags diag.Diagnostics) {
	switch config.ScheduleType {
	case apiclient.MonitorConfigScheduleTypeCrontab:
		schedule, scheduleErr := config.Schedule.AsMonitorConfigScheduleString()
		if scheduleErr != nil {
			diags.AddAttributeError(path.AtName("schedule"), "Invalid schedule", scheduleErr.Error())
			break
		}
		m.ScheduleCrontab = types.StringValue(schedule)
		m.ScheduleInterval = types.ObjectNull((&MonitorConfigScheduleIntervalResourceModel{}).AttributeTypes())
	case apiclient.MonitorConfigScheduleTypeInterval:
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

	m.CheckinMargin = types.Int64PointerValue(nullableInt64ToPointer(config.CheckinMargin))
	m.MaxRuntime = types.Int64PointerValue(nullableInt64ToPointer(config.MaxRuntime))
	timezone := nullableStringToPointer(config.Timezone)
	if timezone == nil || *timezone == "" {
		m.Timezone = types.StringNull()
	} else {
		m.Timezone = types.StringValue(*timezone)
	}
	m.FailureIssueThreshold = types.Int64PointerValue(nullableInt64ToPointer(config.FailureIssueThreshold))
	m.RecoveryThreshold = types.Int64PointerValue(nullableInt64ToPointer(config.RecoveryThreshold))

	return
}

func int64ToNullable(value types.Int64) nullable.Nullable[int64] {
	var result nullable.Nullable[int64]
	if value.IsNull() || value.IsUnknown() {
		return result
	}
	result.Set(value.ValueInt64())
	return result
}

func nullableInt64ToPointer(value nullable.Nullable[int64]) *int64 {
	if value.IsNull() || !value.IsSpecified() {
		return nil
	}
	parsed, err := value.Get()
	if err != nil {
		return nil
	}
	return &parsed
}

func nullableStringToPointer(value nullable.Nullable[string]) *string {
	if value.IsNull() || !value.IsSpecified() {
		return nil
	}
	parsed, err := value.Get()
	if err != nil {
		return nil
	}
	return &parsed
}
