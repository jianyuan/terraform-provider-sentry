package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/go-utils/sliceutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
)

type IssueAlertConditionFirstSeenEventModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionFirstSeenEventModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionFirstSeenEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionFirstSeenEventModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionFirstSeenEvent(apiclient.ProjectRuleConditionFirstSeenEvent{
		Name: m.Name.ValueStringPointer(),
	}))
	return v
}

type IssueAlertConditionRegressionEventModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionRegressionEventModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionRegressionEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionRegressionEventModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionRegressionEvent(apiclient.ProjectRuleConditionRegressionEvent{
		Name: m.Name.ValueStringPointer(),
	}))
	return v
}

type IssueAlertConditionReappearedEventModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionReappearedEventModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionReappearedEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionReappearedEventModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionReappearedEvent(apiclient.ProjectRuleConditionReappearedEvent{
		Name: m.Name.ValueStringPointer(),
	}))
	return v
}

type IssueAlertConditionNewHighPriorityIssueModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionNewHighPriorityIssueModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionNewHighPriorityIssue) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionNewHighPriorityIssueModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionNewHighPriorityIssue(apiclient.ProjectRuleConditionNewHighPriorityIssue{
		Name: m.Name.ValueStringPointer(),
	}))
	return v
}

type IssueAlertCondtionExistingHighPriorityIssueModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertCondtionExistingHighPriorityIssueModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionExistingHighPriorityIssue) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertCondtionExistingHighPriorityIssueModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionExistingHighPriorityIssue(apiclient.ProjectRuleConditionExistingHighPriorityIssue{
		Name: m.Name.ValueStringPointer(),
	}))
	return v
}

type IssueAlertConditionEventFrequencyModel struct {
	Name               types.String `tfsdk:"name"`
	ComparisonType     types.String `tfsdk:"comparison_type"`
	ComparisonInterval types.String `tfsdk:"comparison_interval"`
	Value              types.Int64  `tfsdk:"value"`
	Interval           types.String `tfsdk:"interval"`
}

func (m *IssueAlertConditionEventFrequencyModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionEventFrequency) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	m.ComparisonType = types.StringValue(condition.ComparisonType)
	m.ComparisonInterval = types.StringPointerValue(condition.ComparisonInterval)
	m.Value = types.Int64Value(condition.Value)
	m.Interval = types.StringValue(condition.Interval)
	return
}

func (m IssueAlertConditionEventFrequencyModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionEventFrequency(apiclient.ProjectRuleConditionEventFrequency{
		Name:               m.Name.ValueStringPointer(),
		ComparisonType:     m.ComparisonType.ValueString(),
		ComparisonInterval: m.ComparisonInterval.ValueStringPointer(),
		Value:              m.Value.ValueInt64(),
		Interval:           m.Interval.ValueString(),
	}))
	return v
}

type IssueAlertConditionEventUniqueUserFrequencyModel struct {
	Name               types.String `tfsdk:"name"`
	ComparisonType     types.String `tfsdk:"comparison_type"`
	ComparisonInterval types.String `tfsdk:"comparison_interval"`
	Value              types.Int64  `tfsdk:"value"`
	Interval           types.String `tfsdk:"interval"`
}

func (m *IssueAlertConditionEventUniqueUserFrequencyModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionEventUniqueUserFrequency) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	m.ComparisonType = types.StringValue(condition.ComparisonType)
	m.ComparisonInterval = types.StringPointerValue(condition.ComparisonInterval)
	m.Value = types.Int64Value(condition.Value)
	m.Interval = types.StringValue(condition.Interval)
	return
}

func (m IssueAlertConditionEventUniqueUserFrequencyModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionEventUniqueUserFrequency(apiclient.ProjectRuleConditionEventUniqueUserFrequency{
		Name:               m.Name.ValueStringPointer(),
		ComparisonType:     m.ComparisonType.ValueString(),
		ComparisonInterval: m.ComparisonInterval.ValueStringPointer(),
		Value:              m.Value.ValueInt64(),
		Interval:           m.Interval.ValueString(),
	}))
	return v
}

type IssueAlertConditionEventFrequencyPercentModel struct {
	Name               types.String  `tfsdk:"name"`
	ComparisonType     types.String  `tfsdk:"comparison_type"`
	ComparisonInterval types.String  `tfsdk:"comparison_interval"`
	Value              types.Float64 `tfsdk:"value"`
	Interval           types.String  `tfsdk:"interval"`
}

func (m *IssueAlertConditionEventFrequencyPercentModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionEventFrequencyPercent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	m.ComparisonType = types.StringValue(condition.ComparisonType)
	m.ComparisonInterval = types.StringPointerValue(condition.ComparisonInterval)
	m.Value = types.Float64Value(condition.Value)
	m.Interval = types.StringValue(condition.Interval)
	return
}

func (m IssueAlertConditionEventFrequencyPercentModel) ToApi() apiclient.ProjectRuleCondition {
	var v apiclient.ProjectRuleCondition
	must.Do(v.FromProjectRuleConditionEventFrequencyPercent(apiclient.ProjectRuleConditionEventFrequencyPercent{
		Name:               m.Name.ValueStringPointer(),
		ComparisonType:     m.ComparisonType.ValueString(),
		ComparisonInterval: m.ComparisonInterval.ValueStringPointer(),
		Value:              m.Value.ValueFloat64(),
		Interval:           m.Interval.ValueString(),
	}))
	return v
}

type IssueAlertConditionModel struct {
	FirstSeenEvent            *IssueAlertConditionFirstSeenEventModel           `tfsdk:"first_seen_event"`
	RegressionEvent           *IssueAlertConditionRegressionEventModel          `tfsdk:"regression_event"`
	ReappearedEvent           *IssueAlertConditionReappearedEventModel          `tfsdk:"reappeared_event"`
	NewHighPriorityIssue      *IssueAlertConditionNewHighPriorityIssueModel     `tfsdk:"new_high_priority_issue"`
	ExistingHighPriorityIssue *IssueAlertCondtionExistingHighPriorityIssueModel `tfsdk:"existing_high_priority_issue"`
	EventFrequency            *IssueAlertConditionEventFrequencyModel           `tfsdk:"event_frequency"`
	EventUniqueUserFrequency  *IssueAlertConditionEventUniqueUserFrequencyModel `tfsdk:"event_unique_user_frequency"`
	EventFrequencyPercent     *IssueAlertConditionEventFrequencyPercentModel    `tfsdk:"event_frequency_percent"`
}

func (m IssueAlertConditionModel) ToApi() apiclient.ProjectRuleCondition {
	if m.FirstSeenEvent != nil {
		return m.FirstSeenEvent.ToApi()
	} else if m.RegressionEvent != nil {
		return m.RegressionEvent.ToApi()
	} else if m.ReappearedEvent != nil {
		return m.ReappearedEvent.ToApi()
	} else if m.NewHighPriorityIssue != nil {
		return m.NewHighPriorityIssue.ToApi()
	} else if m.ExistingHighPriorityIssue != nil {
		return m.ExistingHighPriorityIssue.ToApi()
	} else if m.EventFrequency != nil {
		return m.EventFrequency.ToApi()
	} else if m.EventUniqueUserFrequency != nil {
		return m.EventUniqueUserFrequency.ToApi()
	} else if m.EventFrequencyPercent != nil {
		return m.EventFrequencyPercent.ToApi()
	}

	panic("unsupported condition")
}

func (m *IssueAlertConditionModel) FromApi(ctx context.Context, condition apiclient.ProjectRuleCondition) (diags diag.Diagnostics) {
	conditionValue, err := condition.ValueByDiscriminator()
	if err != nil {
		diags.AddError("Invalid condition", err.Error())
		return
	}

	switch conditionValue := conditionValue.(type) {
	case apiclient.ProjectRuleConditionFirstSeenEvent:
		m.FirstSeenEvent = &IssueAlertConditionFirstSeenEventModel{}
		diags.Append(m.FirstSeenEvent.Fill(ctx, conditionValue)...)
	case apiclient.ProjectRuleConditionRegressionEvent:
		m.RegressionEvent = &IssueAlertConditionRegressionEventModel{}
		diags.Append(m.RegressionEvent.Fill(ctx, conditionValue)...)
	case apiclient.ProjectRuleConditionReappearedEvent:
		m.ReappearedEvent = &IssueAlertConditionReappearedEventModel{}
		diags.Append(m.ReappearedEvent.Fill(ctx, conditionValue)...)
	case apiclient.ProjectRuleConditionNewHighPriorityIssue:
		m.NewHighPriorityIssue = &IssueAlertConditionNewHighPriorityIssueModel{}
		diags.Append(m.NewHighPriorityIssue.Fill(ctx, conditionValue)...)
	case apiclient.ProjectRuleConditionExistingHighPriorityIssue:
		m.ExistingHighPriorityIssue = &IssueAlertCondtionExistingHighPriorityIssueModel{}
		diags.Append(m.ExistingHighPriorityIssue.Fill(ctx, conditionValue)...)
	case apiclient.ProjectRuleConditionEventFrequency:
		m.EventFrequency = &IssueAlertConditionEventFrequencyModel{}
		diags.Append(m.EventFrequency.Fill(ctx, conditionValue)...)
	case apiclient.ProjectRuleConditionEventUniqueUserFrequency:
		m.EventUniqueUserFrequency = &IssueAlertConditionEventUniqueUserFrequencyModel{}
		diags.Append(m.EventUniqueUserFrequency.Fill(ctx, conditionValue)...)
	case apiclient.ProjectRuleConditionEventFrequencyPercent:
		m.EventFrequencyPercent = &IssueAlertConditionEventFrequencyPercentModel{}
		diags.Append(m.EventFrequencyPercent.Fill(ctx, conditionValue)...)
	default:
		diags.AddError("Unsupported condition", fmt.Sprintf("Unsupported condition type %T", conditionValue))
	}

	return
}

type IssueAlertModel struct {
	Id           types.String                `tfsdk:"id"`
	Organization types.String                `tfsdk:"organization"`
	Project      types.String                `tfsdk:"project"`
	Name         types.String                `tfsdk:"name"`
	Conditions   sentrytypes.LossyJson       `tfsdk:"conditions"`
	Filters      sentrytypes.LossyJson       `tfsdk:"filters"`
	Actions      sentrytypes.LossyJson       `tfsdk:"actions"`
	ActionMatch  types.String                `tfsdk:"action_match"`
	FilterMatch  types.String                `tfsdk:"filter_match"`
	Frequency    types.Int64                 `tfsdk:"frequency"`
	Environment  types.String                `tfsdk:"environment"`
	Owner        types.String                `tfsdk:"owner"`
	ConditionsV2 *[]IssueAlertConditionModel `tfsdk:"conditions_v2"`
}

func (m *IssueAlertModel) Fill(ctx context.Context, alert apiclient.ProjectRule) (diags diag.Diagnostics) {
	m.Id = types.StringValue(alert.Id)

	if len(alert.Projects) != 1 {
		diags.AddError("Invalid project count", fmt.Sprintf("Expected 1 project, got %d", len(alert.Projects)))
		return
	}
	m.Project = types.StringValue(alert.Projects[0])
	m.Name = types.StringValue(alert.Name)
	m.ActionMatch = types.StringValue(alert.ActionMatch)
	m.FilterMatch = types.StringValue(alert.FilterMatch)
	m.Frequency = types.Int64Value(alert.Frequency)
	m.Environment = types.StringPointerValue(alert.Environment)
	m.Owner = types.StringPointerValue(alert.Owner)

	if !m.Conditions.IsNull() {
		m.Conditions = sentrytypes.NewLossyJsonNull()
		if len(alert.Conditions) > 0 {
			if conditions, err := json.Marshal(alert.Conditions); err == nil {
				m.Conditions = sentrytypes.NewLossyJsonValue(string(conditions))
			} else {
				diags.AddError("Invalid conditions", err.Error())
				return
			}
		}
	} else if m.ConditionsV2 != nil {
		m.ConditionsV2 = ptr.Ptr(sliceutils.Map(func(condition apiclient.ProjectRuleCondition) IssueAlertConditionModel {
			var conditionModel IssueAlertConditionModel
			diags.Append(conditionModel.FromApi(ctx, condition)...)
			return conditionModel
		}, alert.Conditions))

		if diags.HasError() {
			return
		}
	}

	m.Filters = sentrytypes.NewLossyJsonNull()
	if len(alert.Filters) > 0 {
		if filters, err := json.Marshal(alert.Filters); err == nil {
			m.Filters = sentrytypes.NewLossyJsonValue(string(filters))
		} else {
			diags.AddError("Invalid filters", err.Error())
		}
	}

	m.Actions = sentrytypes.NewLossyJsonNull()
	if len(alert.Actions) > 0 {
		if actions, err := json.Marshal(alert.Actions); err == nil && len(actions) > 0 {
			m.Actions = sentrytypes.NewLossyJsonValue(string(actions))
		} else {
			diags.AddError("Invalid actions", err.Error())
		}
	}

	return
}
