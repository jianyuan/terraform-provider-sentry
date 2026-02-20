package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/sliceutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

// Conditions

type IssueAlertConditionFirstSeenEventModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionFirstSeenEventModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionFirstSeenEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionFirstSeenEventModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionFirstSeenEvent(apiclient.ProjectRuleConditionFirstSeenEvent{
		Name: m.Name.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertConditionRegressionEventModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionRegressionEventModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionRegressionEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionRegressionEventModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionRegressionEvent(apiclient.ProjectRuleConditionRegressionEvent{
		Name: m.Name.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertConditionReappearedEventModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionReappearedEventModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionReappearedEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionReappearedEventModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionReappearedEvent(apiclient.ProjectRuleConditionReappearedEvent{
		Name: m.Name.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertConditionNewHighPriorityIssueModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionNewHighPriorityIssueModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionNewHighPriorityIssue) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionNewHighPriorityIssueModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionNewHighPriorityIssue(apiclient.ProjectRuleConditionNewHighPriorityIssue{
		Name: m.Name.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertConditionExistingHighPriorityIssueModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertConditionExistingHighPriorityIssueModel) Fill(ctx context.Context, condition apiclient.ProjectRuleConditionExistingHighPriorityIssue) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(condition.Name)
	return
}

func (m IssueAlertConditionExistingHighPriorityIssueModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionExistingHighPriorityIssue(apiclient.ProjectRuleConditionExistingHighPriorityIssue{
		Name: m.Name.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
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

func (m IssueAlertConditionEventFrequencyModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionEventFrequency(apiclient.ProjectRuleConditionEventFrequency{
		Name:               m.Name.ValueStringPointer(),
		ComparisonType:     m.ComparisonType.ValueString(),
		ComparisonInterval: m.ComparisonInterval.ValueStringPointer(),
		Value:              m.Value.ValueInt64(),
		Interval:           m.Interval.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
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

func (m IssueAlertConditionEventUniqueUserFrequencyModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionEventUniqueUserFrequency(apiclient.ProjectRuleConditionEventUniqueUserFrequency{
		Name:               m.Name.ValueStringPointer(),
		ComparisonType:     m.ComparisonType.ValueString(),
		ComparisonInterval: m.ComparisonInterval.ValueStringPointer(),
		Value:              m.Value.ValueInt64(),
		Interval:           m.Interval.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
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

func (m IssueAlertConditionEventFrequencyPercentModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleCondition
	err := v.FromProjectRuleConditionEventFrequencyPercent(apiclient.ProjectRuleConditionEventFrequencyPercent{
		Name:               m.Name.ValueStringPointer(),
		ComparisonType:     m.ComparisonType.ValueString(),
		ComparisonInterval: m.ComparisonInterval.ValueStringPointer(),
		Value:              m.Value.ValueFloat64(),
		Interval:           m.Interval.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertConditionModel struct {
	FirstSeenEvent            *IssueAlertConditionFirstSeenEventModel            `tfsdk:"first_seen_event"`
	RegressionEvent           *IssueAlertConditionRegressionEventModel           `tfsdk:"regression_event"`
	ReappearedEvent           *IssueAlertConditionReappearedEventModel           `tfsdk:"reappeared_event"`
	NewHighPriorityIssue      *IssueAlertConditionNewHighPriorityIssueModel      `tfsdk:"new_high_priority_issue"`
	ExistingHighPriorityIssue *IssueAlertConditionExistingHighPriorityIssueModel `tfsdk:"existing_high_priority_issue"`
	EventFrequency            *IssueAlertConditionEventFrequencyModel            `tfsdk:"event_frequency"`
	EventUniqueUserFrequency  *IssueAlertConditionEventUniqueUserFrequencyModel  `tfsdk:"event_unique_user_frequency"`
	EventFrequencyPercent     *IssueAlertConditionEventFrequencyPercentModel     `tfsdk:"event_frequency_percent"`
}

func (m IssueAlertConditionModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	if m.FirstSeenEvent != nil {
		return m.FirstSeenEvent.ToApi(ctx)
	} else if m.RegressionEvent != nil {
		return m.RegressionEvent.ToApi(ctx)
	} else if m.ReappearedEvent != nil {
		return m.ReappearedEvent.ToApi(ctx)
	} else if m.NewHighPriorityIssue != nil {
		return m.NewHighPriorityIssue.ToApi(ctx)
	} else if m.ExistingHighPriorityIssue != nil {
		return m.ExistingHighPriorityIssue.ToApi(ctx)
	} else if m.EventFrequency != nil {
		return m.EventFrequency.ToApi(ctx)
	} else if m.EventUniqueUserFrequency != nil {
		return m.EventUniqueUserFrequency.ToApi(ctx)
	} else if m.EventFrequencyPercent != nil {
		return m.EventFrequencyPercent.ToApi(ctx)
	} else {
		var diags diag.Diagnostics
		diags.AddError("Exactly one condition must be set", "Exactly one condition must be set")
		return nil, diags
	}
}

func (m *IssueAlertConditionModel) Fill(ctx context.Context, condition apiclient.ProjectRuleCondition) (diags diag.Diagnostics) {
	conditionValue, err := condition.ValueByDiscriminator()
	if err != nil {
		diags.AddError("Invalid condition", err.Error())
		return
	}

	m.FirstSeenEvent = nil
	m.RegressionEvent = nil
	m.ReappearedEvent = nil
	m.NewHighPriorityIssue = nil
	m.ExistingHighPriorityIssue = nil
	m.EventFrequency = nil
	m.EventUniqueUserFrequency = nil
	m.EventFrequencyPercent = nil

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
		m.ExistingHighPriorityIssue = &IssueAlertConditionExistingHighPriorityIssueModel{}
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

// Filters

type IssueAlertFilterAgeComparisonModel struct {
	Name           types.String `tfsdk:"name"`
	ComparisonType types.String `tfsdk:"comparison_type"`
	Value          types.Int64  `tfsdk:"value"`
	Time           types.String `tfsdk:"time"`
}

func (m *IssueAlertFilterAgeComparisonModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterAgeComparison) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	m.ComparisonType = types.StringValue(filter.ComparisonType)
	m.Value = types.Int64Value(filter.Value)
	m.Time = types.StringValue(filter.Time)
	return
}

func (m IssueAlertFilterAgeComparisonModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterAgeComparison(apiclient.ProjectRuleFilterAgeComparison{
		Name:           m.Name.ValueStringPointer(),
		ComparisonType: m.ComparisonType.ValueString(),
		Value:          m.Value.ValueInt64(),
		Time:           m.Time.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterIssueOccurrencesModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.Int64  `tfsdk:"value"`
}

func (m *IssueAlertFilterIssueOccurrencesModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterIssueOccurrences) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	m.Value = types.Int64Value(filter.Value)
	return
}

func (m IssueAlertFilterIssueOccurrencesModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterIssueOccurrences(apiclient.ProjectRuleFilterIssueOccurrences{
		Name:  m.Name.ValueStringPointer(),
		Value: m.Value.ValueInt64(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterAssignedToModel struct {
	Name             types.String `tfsdk:"name"`
	TargetType       types.String `tfsdk:"target_type"`
	TargetIdentifier types.String `tfsdk:"target_identifier"`
}

func (m *IssueAlertFilterAssignedToModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterAssignedTo) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	m.TargetType = types.StringValue(filter.TargetType)

	if filter.TargetIdentifier == nil {
		m.TargetIdentifier = types.StringNull()
	} else if v, err := filter.TargetIdentifier.AsProjectRuleFilterAssignedToTargetIdentifier0(); err == nil {
		if v == "" {
			m.TargetIdentifier = types.StringNull()
		} else {
			m.TargetIdentifier = types.StringValue(v)
		}
	} else if v, err := filter.TargetIdentifier.AsProjectRuleFilterAssignedToTargetIdentifier1(); err == nil {
		m.TargetIdentifier = types.StringValue(v.String())
	}

	return
}

func (m IssueAlertFilterAssignedToModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	var targetIdentifier *apiclient.ProjectRuleFilterAssignedTo_TargetIdentifier

	if !m.TargetIdentifier.IsNull() {
		targetIdentifier = &apiclient.ProjectRuleFilterAssignedTo_TargetIdentifier{}
		err := targetIdentifier.FromProjectRuleFilterAssignedToTargetIdentifier0(m.TargetIdentifier.ValueString())
		if err != nil {
			diags.AddError("Failed to convert to API model", err.Error())
			return nil, diags
		}
	}

	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterAssignedTo(apiclient.ProjectRuleFilterAssignedTo{
		Name:             m.Name.ValueStringPointer(),
		TargetType:       m.TargetType.ValueString(),
		TargetIdentifier: targetIdentifier,
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterLatestAdoptedReleaseModel struct {
	Name           types.String `tfsdk:"name"`
	OldestOrNewest types.String `tfsdk:"oldest_or_newest"`
	OlderOrNewer   types.String `tfsdk:"older_or_newer"`
	Environment    types.String `tfsdk:"environment"`
}

func (m *IssueAlertFilterLatestAdoptedReleaseModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterLatestAdoptedRelease) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	m.OldestOrNewest = types.StringValue(filter.OldestOrNewest)
	m.OlderOrNewer = types.StringValue(filter.OlderOrNewer)
	m.Environment = types.StringValue(filter.Environment)
	return
}

func (m IssueAlertFilterLatestAdoptedReleaseModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterLatestAdoptedRelease(apiclient.ProjectRuleFilterLatestAdoptedRelease{
		Name:           m.Name.ValueStringPointer(),
		OldestOrNewest: m.OldestOrNewest.ValueString(),
		OlderOrNewer:   m.OlderOrNewer.ValueString(),
		Environment:    m.Environment.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterLatestReleaseModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertFilterLatestReleaseModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterLatestRelease) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	return
}

func (m IssueAlertFilterLatestReleaseModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterLatestRelease(apiclient.ProjectRuleFilterLatestRelease{
		Name: m.Name.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterIssueCategoryModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (m *IssueAlertFilterIssueCategoryModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterIssueCategory) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)

	value, ok := sentrydata.IssueGroupCategoryIdToName[filter.Value]
	if !ok {
		diags.AddError("Invalid issue category", fmt.Sprintf("Invalid issue category %q. Please report this to the provider developers.", filter.Value))
		return
	}
	m.Value = types.StringValue(value)

	return
}

func (m IssueAlertFilterIssueCategoryModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterIssueCategory(apiclient.ProjectRuleFilterIssueCategory{
		Name:  m.Name.ValueStringPointer(),
		Value: sentrydata.IssueGroupCategoryNameToId[m.Value.ValueString()],
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterEventAttributeModel struct {
	Name      types.String `tfsdk:"name"`
	Attribute types.String `tfsdk:"attribute"`
	Match     types.String `tfsdk:"match"`
	Value     types.String `tfsdk:"value"`
}

func (m *IssueAlertFilterEventAttributeModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterEventAttribute) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	m.Attribute = types.StringValue(filter.Attribute)

	match, ok := sentrydata.MatchTypeIdToName[filter.Match]
	if !ok {
		diags.AddError("Invalid match type", fmt.Sprintf("Invalid match type %q. Please report this to the provider developers.", filter.Match))
		return
	}
	m.Match = types.StringValue(match)

	if filter.Value == nil || *filter.Value == "" {
		m.Value = types.StringNull()
	} else {
		m.Value = types.StringValue(*filter.Value)
	}
	return
}

func (m IssueAlertFilterEventAttributeModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterEventAttribute(apiclient.ProjectRuleFilterEventAttribute{
		Name:      m.Name.ValueStringPointer(),
		Attribute: m.Attribute.ValueString(),
		Match:     sentrydata.MatchTypeNameToId[m.Match.ValueString()],
		Value:     m.Value.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterTaggedEventModel struct {
	Name  types.String `tfsdk:"name"`
	Key   types.String `tfsdk:"key"`
	Match types.String `tfsdk:"match"`
	Value types.String `tfsdk:"value"`
}

func (m *IssueAlertFilterTaggedEventModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterTaggedEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	m.Key = types.StringValue(filter.Key)

	match, ok := sentrydata.MatchTypeIdToName[filter.Match]
	if !ok {
		diags.AddError("Invalid match type", fmt.Sprintf("Invalid match type %q. Please report this to the provider developers.", filter.Match))
		return
	}
	m.Match = types.StringValue(match)

	if filter.Value == nil || *filter.Value == "" {
		m.Value = types.StringNull()
	} else {
		m.Value = types.StringValue(*filter.Value)
	}
	return
}

func (m IssueAlertFilterTaggedEventModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterTaggedEvent(apiclient.ProjectRuleFilterTaggedEvent{
		Name:  m.Name.ValueStringPointer(),
		Key:   m.Key.ValueString(),
		Match: sentrydata.MatchTypeNameToId[m.Match.ValueString()],
		Value: m.Value.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterLevelModel struct {
	Name  types.String `tfsdk:"name"`
	Match types.String `tfsdk:"match"`
	Level types.String `tfsdk:"level"`
}

func (m *IssueAlertFilterLevelModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterLevel) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)

	match, ok := sentrydata.MatchTypeIdToName[filter.Match]
	if !ok {
		diags.AddError("Invalid match type", fmt.Sprintf("Invalid match type %q. Please report this to the provider developers.", filter.Match))
		return
	}
	m.Match = types.StringValue(match)

	level, ok := sentrydata.LogLevelIdToName[filter.Level]
	if !ok {
		diags.AddError("Invalid level", fmt.Sprintf("Invalid level %q. Please report this to the provider developers.", filter.Level))
		return
	}
	m.Level = types.StringValue(level)
	return
}

func (m IssueAlertFilterLevelModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter
	err := v.FromProjectRuleFilterLevel(apiclient.ProjectRuleFilterLevel{
		Name:  m.Name.ValueStringPointer(),
		Match: sentrydata.MatchTypeNameToId[m.Match.ValueString()],
		Level: sentrydata.LogLevelNameToId[m.Level.ValueString()],
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertFilterModel struct {
	AgeComparison        *IssueAlertFilterAgeComparisonModel        `tfsdk:"age_comparison"`
	IssueOccurrences     *IssueAlertFilterIssueOccurrencesModel     `tfsdk:"issue_occurrences"`
	AssignedTo           *IssueAlertFilterAssignedToModel           `tfsdk:"assigned_to"`
	LatestAdoptedRelease *IssueAlertFilterLatestAdoptedReleaseModel `tfsdk:"latest_adopted_release"`
	LatestRelease        *IssueAlertFilterLatestReleaseModel        `tfsdk:"latest_release"`
	IssueCategory        *IssueAlertFilterIssueCategoryModel        `tfsdk:"issue_category"`
	EventAttribute       *IssueAlertFilterEventAttributeModel       `tfsdk:"event_attribute"`
	TaggedEvent          *IssueAlertFilterTaggedEventModel          `tfsdk:"tagged_event"`
	Level                *IssueAlertFilterLevelModel                `tfsdk:"level"`
}

func (m IssueAlertFilterModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	if m.AgeComparison != nil {
		return m.AgeComparison.ToApi(ctx)
	} else if m.IssueOccurrences != nil {
		return m.IssueOccurrences.ToApi(ctx)
	} else if m.AssignedTo != nil {
		return m.AssignedTo.ToApi(ctx)
	} else if m.LatestAdoptedRelease != nil {
		return m.LatestAdoptedRelease.ToApi(ctx)
	} else if m.LatestRelease != nil {
		return m.LatestRelease.ToApi(ctx)
	} else if m.IssueCategory != nil {
		return m.IssueCategory.ToApi(ctx)
	} else if m.EventAttribute != nil {
		return m.EventAttribute.ToApi(ctx)
	} else if m.TaggedEvent != nil {
		return m.TaggedEvent.ToApi(ctx)
	} else if m.Level != nil {
		return m.Level.ToApi(ctx)
	} else {
		var diags diag.Diagnostics
		diags.AddError("Exactly one filter must be set", "Exactly one filter must be set")
		return nil, diags
	}
}

func (m *IssueAlertFilterModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilter) (diags diag.Diagnostics) {
	filterValue, err := filter.ValueByDiscriminator()
	if err != nil {
		diags.AddError("Invalid filter", err.Error())
		return
	}

	m.AgeComparison = nil
	m.IssueOccurrences = nil
	m.AssignedTo = nil
	m.LatestAdoptedRelease = nil
	m.LatestRelease = nil
	m.IssueCategory = nil
	m.EventAttribute = nil
	m.TaggedEvent = nil
	m.Level = nil

	switch filterValue := filterValue.(type) {
	case apiclient.ProjectRuleFilterAgeComparison:
		m.AgeComparison = &IssueAlertFilterAgeComparisonModel{}
		diags.Append(m.AgeComparison.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterIssueOccurrences:
		m.IssueOccurrences = &IssueAlertFilterIssueOccurrencesModel{}
		diags.Append(m.IssueOccurrences.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterAssignedTo:
		m.AssignedTo = &IssueAlertFilterAssignedToModel{}
		diags.Append(m.AssignedTo.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterLatestAdoptedRelease:
		m.LatestAdoptedRelease = &IssueAlertFilterLatestAdoptedReleaseModel{}
		diags.Append(m.LatestAdoptedRelease.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterLatestRelease:
		m.LatestRelease = &IssueAlertFilterLatestReleaseModel{}
		diags.Append(m.LatestRelease.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterIssueCategory:
		m.IssueCategory = &IssueAlertFilterIssueCategoryModel{}
		diags.Append(m.IssueCategory.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterEventAttribute:
		m.EventAttribute = &IssueAlertFilterEventAttributeModel{}
		diags.Append(m.EventAttribute.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterTaggedEvent:
		m.TaggedEvent = &IssueAlertFilterTaggedEventModel{}
		diags.Append(m.TaggedEvent.Fill(ctx, filterValue)...)
	case apiclient.ProjectRuleFilterLevel:
		m.Level = &IssueAlertFilterLevelModel{}
		diags.Append(m.Level.Fill(ctx, filterValue)...)
	default:
		diags.AddError("Unsupported filter", fmt.Sprintf("Unsupported filter type %T", filterValue))
	}

	return
}

// Actions

type IssueAlertActionNotifyEmailModel struct {
	Name             types.String `tfsdk:"name"`
	TargetType       types.String `tfsdk:"target_type"`
	TargetIdentifier types.String `tfsdk:"target_identifier"`
	FallthroughType  types.String `tfsdk:"fallthrough_type"`
}

func (m *IssueAlertActionNotifyEmailModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionNotifyEmail) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.TargetType = types.StringValue(action.TargetType)

	if action.TargetIdentifier == nil {
		m.TargetIdentifier = types.StringNull()
	} else if v, err := action.TargetIdentifier.AsProjectRuleActionNotifyEmailTargetIdentifier0(); err == nil {
		if v == "" {
			m.TargetIdentifier = types.StringNull()
		} else {
			m.TargetIdentifier = types.StringValue(v)
		}
	} else if v, err := action.TargetIdentifier.AsProjectRuleActionNotifyEmailTargetIdentifier1(); err == nil {
		m.TargetIdentifier = types.StringValue(v.String())
	}

	// Only set FallthroughType for IssueOwners
	if action.TargetType == "IssueOwners" {
		m.FallthroughType = types.StringPointerValue(action.FallthroughType)
	} else {
		m.FallthroughType = types.StringNull()
	}

	return
}

func (m IssueAlertActionNotifyEmailModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var targetIdentifier *apiclient.ProjectRuleActionNotifyEmail_TargetIdentifier

	if !m.TargetIdentifier.IsNull() {
		targetIdentifier = &apiclient.ProjectRuleActionNotifyEmail_TargetIdentifier{}
		err := targetIdentifier.FromProjectRuleActionNotifyEmailTargetIdentifier0(m.TargetIdentifier.ValueString())
		if err != nil {
			diags.AddError("Failed to convert to API model", err.Error())
			return nil, diags
		}
	}

	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionNotifyEmail(apiclient.ProjectRuleActionNotifyEmail{
		Name:             m.Name.ValueStringPointer(),
		TargetType:       m.TargetType.ValueString(),
		TargetIdentifier: targetIdentifier,
		FallthroughType:  m.FallthroughType.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionNotifyEventModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertActionNotifyEventModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionNotifyEvent) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	return
}

func (m IssueAlertActionNotifyEventModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionNotifyEvent(apiclient.ProjectRuleActionNotifyEvent{
		Name: m.Name.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionNotifyEventServiceModel struct {
	Name    types.String `tfsdk:"name"`
	Service types.String `tfsdk:"service"`
}

func (m *IssueAlertActionNotifyEventServiceModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionNotifyEventService) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Service = types.StringValue(action.Service)
	return
}

func (m IssueAlertActionNotifyEventServiceModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionNotifyEventService(apiclient.ProjectRuleActionNotifyEventService{
		Name:    m.Name.ValueStringPointer(),
		Service: m.Service.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionNotifyEventSentryAppModel struct {
	Name                      types.String `tfsdk:"name"`
	SentryAppInstallationUuid types.String `tfsdk:"sentry_app_installation_uuid"`
	Settings                  types.Map    `tfsdk:"settings"`
}

func (m *IssueAlertActionNotifyEventSentryAppModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionNotifyEventSentryApp) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.SentryAppInstallationUuid = types.StringValue(action.SentryAppInstallationUuid)

	if action.Settings == nil {
		m.Settings = types.MapNull(types.StringType)
	} else {
		var settingsMap = make(map[string]attr.Value)
		for _, setting := range *action.Settings {
			settingsMap[setting.Name] = types.StringValue(setting.Value)
		}
		m.Settings = types.MapValueMust(types.StringType, settingsMap)
	}
	return
}

func (m IssueAlertActionNotifyEventSentryAppModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction

	var settings *[]struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	if !m.Settings.IsNull() {
		elements := make(map[string]string, len(m.Settings.Elements()))
		diags.Append(m.Settings.ElementsAs(ctx, &elements, false)...)
		if diags.HasError() {
			return nil, diags
		}

		settings = &[]struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{}

		for k, v := range elements {
			*settings = append(*settings, struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}{
				Name:  k,
				Value: v,
			})
		}
	}

	err := v.FromProjectRuleActionNotifyEventSentryApp(apiclient.ProjectRuleActionNotifyEventSentryApp{
		Name:                      m.Name.ValueStringPointer(),
		SentryAppInstallationUuid: m.SentryAppInstallationUuid.ValueString(),
		Settings:                  settings,
		HasSchemaFormConfig:       true,
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionOpsgenieNotifyTeam struct {
	Name     types.String `tfsdk:"name"`
	Account  types.String `tfsdk:"account"`
	Team     types.String `tfsdk:"team"`
	Priority types.String `tfsdk:"priority"`
}

func (m *IssueAlertActionOpsgenieNotifyTeam) Fill(ctx context.Context, action apiclient.ProjectRuleActionOpsgenieNotifyTeam) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Account = types.StringValue(action.Account)
	m.Team = types.StringValue(action.Team)
	m.Priority = types.StringValue(action.Priority)
	return
}

func (m IssueAlertActionOpsgenieNotifyTeam) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionOpsgenieNotifyTeam(apiclient.ProjectRuleActionOpsgenieNotifyTeam{
		Name:     m.Name.ValueStringPointer(),
		Account:  m.Account.ValueString(),
		Team:     m.Team.ValueString(),
		Priority: m.Priority.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionPagerDutyNotifyServiceModel struct {
	Name     types.String `tfsdk:"name"`
	Account  types.String `tfsdk:"account"`
	Service  types.String `tfsdk:"service"`
	Severity types.String `tfsdk:"severity"`
}

func (m *IssueAlertActionPagerDutyNotifyServiceModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionPagerDutyNotifyService) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Account = types.StringValue(action.Account)
	m.Service = types.StringValue(action.Service)
	m.Severity = types.StringValue(action.Severity)
	return
}

func (m IssueAlertActionPagerDutyNotifyServiceModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionPagerDutyNotifyService(apiclient.ProjectRuleActionPagerDutyNotifyService{
		Name:     m.Name.ValueStringPointer(),
		Account:  m.Account.ValueString(),
		Service:  m.Service.ValueString(),
		Severity: m.Severity.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionSlackNotifyServiceModel struct {
	Name      types.String          `tfsdk:"name"`
	Workspace types.String          `tfsdk:"workspace"`
	Channel   types.String          `tfsdk:"channel"`
	ChannelId types.String          `tfsdk:"channel_id"`
	Tags      sentrytypes.StringSet `tfsdk:"tags"`
	Notes     types.String          `tfsdk:"notes"`
}

func (m *IssueAlertActionSlackNotifyServiceModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionSlackNotifyService) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Workspace = types.StringValue(action.Workspace)
	m.Channel = types.StringValue(action.Channel)
	m.ChannelId = types.StringPointerValue(action.ChannelId)
	m.Tags = tfutils.MergeDiagnostics(sentrytypes.StringSetPointerValue(action.Tags))(&diags)
	m.Notes = types.StringPointerValue(action.Notes)
	return
}

func (m IssueAlertActionSlackNotifyServiceModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionSlackNotifyService(apiclient.ProjectRuleActionSlackNotifyService{
		Name:      m.Name.ValueStringPointer(),
		Workspace: m.Workspace.ValueString(),
		Channel:   m.Channel.ValueString(),
		ChannelId: m.ChannelId.ValueStringPointer(),
		Tags:      tfutils.MergeDiagnostics(m.Tags.ValueStringPointer(ctx))(&diags),
		Notes:     m.Notes.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionMsTeamsNotifyServiceModel struct {
	Name      types.String `tfsdk:"name"`
	Team      types.String `tfsdk:"team"`
	Channel   types.String `tfsdk:"channel"`
	ChannelId types.String `tfsdk:"channel_id"`
}

func (m *IssueAlertActionMsTeamsNotifyServiceModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionMsTeamsNotifyService) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Team = types.StringValue(action.Team)
	m.Channel = types.StringValue(action.Channel)
	m.ChannelId = types.StringPointerValue(action.ChannelId)
	return
}

func (m IssueAlertActionMsTeamsNotifyServiceModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionMsTeamsNotifyService(apiclient.ProjectRuleActionMsTeamsNotifyService{
		Name:      m.Name.ValueStringPointer(),
		Team:      m.Team.ValueString(),
		Channel:   m.Channel.ValueString(),
		ChannelId: m.ChannelId.ValueStringPointer(),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionDiscordNotifyServiceModel struct {
	Name      types.String          `tfsdk:"name"`
	Server    types.String          `tfsdk:"server"`
	ChannelId types.String          `tfsdk:"channel_id"`
	Tags      sentrytypes.StringSet `tfsdk:"tags"`
}

func (m *IssueAlertActionDiscordNotifyServiceModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionDiscordNotifyService) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Server = types.StringValue(action.Server)
	m.ChannelId = types.StringValue(action.ChannelId)
	m.Tags = tfutils.MergeDiagnostics(sentrytypes.StringSetPointerValue(action.Tags))(&diags)
	return
}

func (m IssueAlertActionDiscordNotifyServiceModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionDiscordNotifyService(apiclient.ProjectRuleActionDiscordNotifyService{
		Name:      m.Name.ValueStringPointer(),
		Server:    m.Server.ValueString(),
		ChannelId: m.ChannelId.ValueString(),
		Tags:      tfutils.MergeDiagnostics(m.Tags.ValueStringPointer(ctx))(&diags),
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionJiraCreateTicketModel struct {
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Project     types.String `tfsdk:"project"`
	IssueType   types.String `tfsdk:"issue_type"`
}

func (m *IssueAlertActionJiraCreateTicketModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionJiraCreateTicket) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Integration = types.StringValue(action.Integration)
	m.Project = types.StringValue(action.Project)
	m.IssueType = types.StringValue(action.IssueType)
	return
}

func (m IssueAlertActionJiraCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionJiraCreateTicket(apiclient.ProjectRuleActionJiraCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       m.Integration.ValueString(),
		Project:           m.Project.ValueString(),
		IssueType:         m.IssueType.ValueString(),
		DynamicFormFields: []map[string]interface{}{{"dummy": "dummy"}},
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionJiraServerCreateTicketModel struct {
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Project     types.String `tfsdk:"project"`
	IssueType   types.String `tfsdk:"issue_type"`
}

func (m *IssueAlertActionJiraServerCreateTicketModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionJiraServerCreateTicket) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Integration = types.StringValue(action.Integration)
	m.Project = types.StringValue(action.Project)
	m.IssueType = types.StringValue(action.IssueType)
	return
}

func (m IssueAlertActionJiraServerCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionJiraServerCreateTicket(apiclient.ProjectRuleActionJiraServerCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       m.Integration.ValueString(),
		Project:           m.Project.ValueString(),
		IssueType:         m.IssueType.ValueString(),
		DynamicFormFields: []map[string]interface{}{{"dummy": "dummy"}},
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionGitHubCreateTicketModel struct {
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Repo        types.String `tfsdk:"repo"`
	Assignee    types.String `tfsdk:"assignee"`
	Labels      types.Set    `tfsdk:"labels"`
}

func (m *IssueAlertActionGitHubCreateTicketModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionGitHubCreateTicket) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Integration = types.StringValue(action.Integration)
	m.Repo = types.StringValue(action.Repo)
	m.Assignee = types.StringPointerValue(action.Assignee)

	if action.Labels == nil {
		m.Labels = types.SetNull(types.StringType)
	} else {
		m.Labels = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
			return types.StringValue(v)
		}, *action.Labels))
	}
	return
}

func (m IssueAlertActionGitHubCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := apiclient.ProjectRuleActionGitHubCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       m.Integration.ValueString(),
		Repo:              m.Repo.ValueString(),
		Assignee:          m.Assignee.ValueStringPointer(),
		DynamicFormFields: []map[string]interface{}{{"dummy": "dummy"}},
	}

	if !m.Labels.IsNull() {
		var labels []string
		diags.Append(m.Labels.ElementsAs(ctx, &labels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		body.Labels = &labels
	}

	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionGitHubCreateTicket(body)
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionGitHubEnterpriseCreateTicketModel struct {
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Repo        types.String `tfsdk:"repo"`
	Assignee    types.String `tfsdk:"assignee"`
	Labels      types.Set    `tfsdk:"labels"`
}

func (m *IssueAlertActionGitHubEnterpriseCreateTicketModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionGitHubEnterpriseCreateTicket) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Integration = types.StringValue(action.Integration)
	m.Repo = types.StringValue(action.Repo)
	m.Assignee = types.StringPointerValue(action.Assignee)

	if action.Labels == nil {
		m.Labels = types.SetNull(types.StringType)
	} else {
		m.Labels = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
			return types.StringValue(v)
		}, *action.Labels))
	}
	return
}

func (m IssueAlertActionGitHubEnterpriseCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := apiclient.ProjectRuleActionGitHubEnterpriseCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       m.Integration.ValueString(),
		Repo:              m.Repo.ValueString(),
		Assignee:          m.Assignee.ValueStringPointer(),
		DynamicFormFields: []map[string]interface{}{{"dummy": "dummy"}},
	}

	if !m.Labels.IsNull() {
		var labels []string
		diags.Append(m.Labels.ElementsAs(ctx, &labels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		body.Labels = &labels
	}

	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionGitHubEnterpriseCreateTicket(body)
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionAzureDevopsCreateTicketModel struct {
	Name         types.String `tfsdk:"name"`
	Integration  types.String `tfsdk:"integration"`
	Project      types.String `tfsdk:"project"`
	WorkItemType types.String `tfsdk:"work_item_type"`
}

func (m *IssueAlertActionAzureDevopsCreateTicketModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionAzureDevopsCreateTicket) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Integration = types.StringValue(action.Integration)
	m.Project = types.StringValue(action.Project)
	m.WorkItemType = types.StringValue(action.WorkItemType)
	return
}

func (m IssueAlertActionAzureDevopsCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionAzureDevopsCreateTicket(apiclient.ProjectRuleActionAzureDevopsCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       m.Integration.ValueString(),
		Project:           m.Project.ValueString(),
		WorkItemType:      m.WorkItemType.ValueString(),
		DynamicFormFields: []map[string]interface{}{{"dummy": "dummy"}},
	})
	if err != nil {
		diags.AddError("Failed to convert to API model", err.Error())
		return nil, diags
	}
	return &v, diags
}

type IssueAlertActionModel struct {
	NotifyEmail                  *IssueAlertActionNotifyEmailModel                  `tfsdk:"notify_email"`
	NotifyEvent                  *IssueAlertActionNotifyEventModel                  `tfsdk:"notify_event"`
	NotifyEventService           *IssueAlertActionNotifyEventServiceModel           `tfsdk:"notify_event_service"`
	NotifyEventSentryApp         *IssueAlertActionNotifyEventSentryAppModel         `tfsdk:"notify_event_sentry_app"`
	OpsgenieNotifyTeam           *IssueAlertActionOpsgenieNotifyTeam                `tfsdk:"opsgenie_notify_team"`
	PagerDutyNotifyService       *IssueAlertActionPagerDutyNotifyServiceModel       `tfsdk:"pagerduty_notify_service"`
	SlackNotifyService           *IssueAlertActionSlackNotifyServiceModel           `tfsdk:"slack_notify_service"`
	MsTeamsNotifyService         *IssueAlertActionMsTeamsNotifyServiceModel         `tfsdk:"msteams_notify_service"`
	DiscordNotifyService         *IssueAlertActionDiscordNotifyServiceModel         `tfsdk:"discord_notify_service"`
	JiraCreateTicket             *IssueAlertActionJiraCreateTicketModel             `tfsdk:"jira_create_ticket"`
	JiraServerCreateTicket       *IssueAlertActionJiraServerCreateTicketModel       `tfsdk:"jira_server_create_ticket"`
	GitHubCreateTicket           *IssueAlertActionGitHubCreateTicketModel           `tfsdk:"github_create_ticket"`
	GitHubEnterpriseCreateTicket *IssueAlertActionGitHubEnterpriseCreateTicketModel `tfsdk:"github_enterprise_create_ticket"`
	AzureDevopsCreateTicket      *IssueAlertActionAzureDevopsCreateTicketModel      `tfsdk:"azure_devops_create_ticket"`
}

func (m IssueAlertActionModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	if m.NotifyEmail != nil {
		return m.NotifyEmail.ToApi(ctx)
	} else if m.NotifyEvent != nil {
		return m.NotifyEvent.ToApi(ctx)
	} else if m.NotifyEventService != nil {
		return m.NotifyEventService.ToApi(ctx)
	} else if m.NotifyEventSentryApp != nil {
		return m.NotifyEventSentryApp.ToApi(ctx)
	} else if m.OpsgenieNotifyTeam != nil {
		return m.OpsgenieNotifyTeam.ToApi(ctx)
	} else if m.PagerDutyNotifyService != nil {
		return m.PagerDutyNotifyService.ToApi(ctx)
	} else if m.SlackNotifyService != nil {
		return m.SlackNotifyService.ToApi(ctx)
	} else if m.MsTeamsNotifyService != nil {
		return m.MsTeamsNotifyService.ToApi(ctx)
	} else if m.DiscordNotifyService != nil {
		return m.DiscordNotifyService.ToApi(ctx)
	} else if m.JiraCreateTicket != nil {
		return m.JiraCreateTicket.ToApi(ctx)
	} else if m.JiraServerCreateTicket != nil {
		return m.JiraServerCreateTicket.ToApi(ctx)
	} else if m.GitHubCreateTicket != nil {
		return m.GitHubCreateTicket.ToApi(ctx)
	} else if m.GitHubEnterpriseCreateTicket != nil {
		return m.GitHubEnterpriseCreateTicket.ToApi(ctx)
	} else if m.AzureDevopsCreateTicket != nil {
		return m.AzureDevopsCreateTicket.ToApi(ctx)
	} else {
		var diags diag.Diagnostics
		diags.AddError("Exactly one action must be set", "Exactly one action must be set")
		return nil, diags
	}
}

func (m *IssueAlertActionModel) Fill(ctx context.Context, action apiclient.ProjectRuleAction) (diags diag.Diagnostics) {
	actionValue, err := action.ValueByDiscriminator()
	if err != nil {
		diags.AddError("Invalid action", err.Error())
		return
	}

	m.NotifyEmail = nil
	m.NotifyEvent = nil
	m.NotifyEventService = nil
	m.NotifyEventSentryApp = nil
	m.OpsgenieNotifyTeam = nil
	m.PagerDutyNotifyService = nil
	m.SlackNotifyService = nil
	m.MsTeamsNotifyService = nil
	m.DiscordNotifyService = nil
	m.JiraCreateTicket = nil
	m.JiraServerCreateTicket = nil
	m.GitHubCreateTicket = nil
	m.GitHubEnterpriseCreateTicket = nil
	m.AzureDevopsCreateTicket = nil

	switch actionValue := actionValue.(type) {
	case apiclient.ProjectRuleActionNotifyEmail:
		m.NotifyEmail = &IssueAlertActionNotifyEmailModel{}
		diags.Append(m.NotifyEmail.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionNotifyEvent:
		m.NotifyEvent = &IssueAlertActionNotifyEventModel{}
		diags.Append(m.NotifyEvent.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionNotifyEventService:
		m.NotifyEventService = &IssueAlertActionNotifyEventServiceModel{}
		diags.Append(m.NotifyEventService.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionNotifyEventSentryApp:
		m.NotifyEventSentryApp = &IssueAlertActionNotifyEventSentryAppModel{}
		diags.Append(m.NotifyEventSentryApp.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionOpsgenieNotifyTeam:
		m.OpsgenieNotifyTeam = &IssueAlertActionOpsgenieNotifyTeam{}
		diags.Append(m.OpsgenieNotifyTeam.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionPagerDutyNotifyService:
		m.PagerDutyNotifyService = &IssueAlertActionPagerDutyNotifyServiceModel{}
		diags.Append(m.PagerDutyNotifyService.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionSlackNotifyService:
		m.SlackNotifyService = &IssueAlertActionSlackNotifyServiceModel{}
		diags.Append(m.SlackNotifyService.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionMsTeamsNotifyService:
		m.MsTeamsNotifyService = &IssueAlertActionMsTeamsNotifyServiceModel{}
		diags.Append(m.MsTeamsNotifyService.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionDiscordNotifyService:
		m.DiscordNotifyService = &IssueAlertActionDiscordNotifyServiceModel{}
		diags.Append(m.DiscordNotifyService.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionJiraCreateTicket:
		m.JiraCreateTicket = &IssueAlertActionJiraCreateTicketModel{}
		diags.Append(m.JiraCreateTicket.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionJiraServerCreateTicket:
		m.JiraServerCreateTicket = &IssueAlertActionJiraServerCreateTicketModel{}
		diags.Append(m.JiraServerCreateTicket.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionGitHubCreateTicket:
		m.GitHubCreateTicket = &IssueAlertActionGitHubCreateTicketModel{}
		diags.Append(m.GitHubCreateTicket.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionGitHubEnterpriseCreateTicket:
		m.GitHubEnterpriseCreateTicket = &IssueAlertActionGitHubEnterpriseCreateTicketModel{}
		diags.Append(m.GitHubEnterpriseCreateTicket.Fill(ctx, actionValue)...)
	case apiclient.ProjectRuleActionAzureDevopsCreateTicket:
		m.AzureDevopsCreateTicket = &IssueAlertActionAzureDevopsCreateTicketModel{}
		diags.Append(m.AzureDevopsCreateTicket.Fill(ctx, actionValue)...)
	default:
		diags.AddError("Unsupported action", fmt.Sprintf("Unsupported action type %T", actionValue))
	}

	return
}

// Model

type IssueAlertModel struct {
	Id           types.String          `tfsdk:"id"`
	Organization types.String          `tfsdk:"organization"`
	Project      types.String          `tfsdk:"project"`
	Name         types.String          `tfsdk:"name"`
	Conditions   sentrytypes.LossyJson `tfsdk:"conditions"`
	Filters      sentrytypes.LossyJson `tfsdk:"filters"`
	Actions      sentrytypes.LossyJson `tfsdk:"actions"`
	ActionMatch  types.String          `tfsdk:"action_match"`
	FilterMatch  types.String          `tfsdk:"filter_match"`
	Frequency    types.Int64           `tfsdk:"frequency"`
	Environment  types.String          `tfsdk:"environment"`
	Owner        types.String          `tfsdk:"owner"`
	ConditionsV2 types.List            `tfsdk:"conditions_v2"`
	FiltersV2    types.List            `tfsdk:"filters_v2"`
	ActionsV2    types.List            `tfsdk:"actions_v2"`
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

	if v, err := alert.Environment.Get(); err == nil {
		m.Environment = types.StringValue(v)
	} else {
		m.Environment = types.StringNull()
	}

	if v, err := alert.Owner.Get(); err == nil {
		m.Owner = types.StringValue(v)
	} else {
		m.Owner = types.StringNull()
	}

	if !m.Conditions.IsNull() {
		if conditions, err := json.Marshal(alert.Conditions); err == nil {
			m.Conditions = sentrytypes.NewLossyJsonValue(string(conditions))
		} else {
			diags.AddError("Invalid conditions", err.Error())
			return
		}
	} else if !m.ConditionsV2.IsNull() {
		conditions := sliceutils.Map(func(condition apiclient.ProjectRuleCondition) IssueAlertConditionModel {
			var conditionModel IssueAlertConditionModel
			diags.Append(conditionModel.Fill(ctx, condition)...)
			return conditionModel
		}, alert.Conditions)

		if diags.HasError() {
			return
		}

		conditionsV2, d := types.ListValueFrom(ctx, issueAlertConditionV2ElemType, conditions)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
		m.ConditionsV2 = conditionsV2
	}

	if !m.Filters.IsNull() {
		if filters, err := json.Marshal(alert.Filters); err == nil {
			m.Filters = sentrytypes.NewLossyJsonValue(string(filters))
		} else {
			diags.AddError("Invalid filters", err.Error())
		}
	} else if !m.FiltersV2.IsNull() {
		filters := sliceutils.Map(func(filter apiclient.ProjectRuleFilter) IssueAlertFilterModel {
			var filterModel IssueAlertFilterModel
			diags.Append(filterModel.Fill(ctx, filter)...)
			return filterModel
		}, alert.Filters)

		if diags.HasError() {
			return
		}

		filtersV2, d := types.ListValueFrom(ctx, issueAlertFilterV2ElemType, filters)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
		m.FiltersV2 = filtersV2
	}

	if !m.Actions.IsNull() {
		if actions, err := json.Marshal(alert.Actions); err == nil && len(actions) > 0 {
			m.Actions = sentrytypes.NewLossyJsonValue(string(actions))
		} else {
			diags.AddError("Invalid actions", err.Error())
		}
	} else if !m.ActionsV2.IsNull() {
		actions := sliceutils.Map(func(action apiclient.ProjectRuleAction) IssueAlertActionModel {
			var actionModel IssueAlertActionModel
			diags.Append(actionModel.Fill(ctx, action)...)
			return actionModel
		}, alert.Actions)

		if diags.HasError() {
			return
		}

		actionsV2, d := types.ListValueFrom(ctx, issueAlertActionV2ElemType, actions)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
		m.ActionsV2 = actionsV2
	}

	return
}
