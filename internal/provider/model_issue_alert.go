package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
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
	FirstSeenEvent            supertypes.SingleNestedObjectValueOf[IssueAlertConditionFirstSeenEventModel]            `tfsdk:"first_seen_event"`
	RegressionEvent           supertypes.SingleNestedObjectValueOf[IssueAlertConditionRegressionEventModel]           `tfsdk:"regression_event"`
	ReappearedEvent           supertypes.SingleNestedObjectValueOf[IssueAlertConditionReappearedEventModel]           `tfsdk:"reappeared_event"`
	NewHighPriorityIssue      supertypes.SingleNestedObjectValueOf[IssueAlertConditionNewHighPriorityIssueModel]      `tfsdk:"new_high_priority_issue"`
	ExistingHighPriorityIssue supertypes.SingleNestedObjectValueOf[IssueAlertConditionExistingHighPriorityIssueModel] `tfsdk:"existing_high_priority_issue"`
	EventFrequency            supertypes.SingleNestedObjectValueOf[IssueAlertConditionEventFrequencyModel]            `tfsdk:"event_frequency"`
	EventUniqueUserFrequency  supertypes.SingleNestedObjectValueOf[IssueAlertConditionEventUniqueUserFrequencyModel]  `tfsdk:"event_unique_user_frequency"`
	EventFrequencyPercent     supertypes.SingleNestedObjectValueOf[IssueAlertConditionEventFrequencyPercentModel]     `tfsdk:"event_frequency_percent"`
}

func (m IssueAlertConditionModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleCondition, diag.Diagnostics) {
	if m.FirstSeenEvent.IsKnown() {
		return m.FirstSeenEvent.MustGet(ctx).ToApi(ctx)
	} else if m.RegressionEvent.IsKnown() {
		return m.RegressionEvent.MustGet(ctx).ToApi(ctx)
	} else if m.ReappearedEvent.IsKnown() {
		return m.ReappearedEvent.MustGet(ctx).ToApi(ctx)
	} else if m.NewHighPriorityIssue.IsKnown() {
		return m.NewHighPriorityIssue.MustGet(ctx).ToApi(ctx)
	} else if m.ExistingHighPriorityIssue.IsKnown() {
		return m.ExistingHighPriorityIssue.MustGet(ctx).ToApi(ctx)
	} else if m.EventFrequency.IsKnown() {
		return m.EventFrequency.MustGet(ctx).ToApi(ctx)
	} else if m.EventUniqueUserFrequency.IsKnown() {
		return m.EventUniqueUserFrequency.MustGet(ctx).ToApi(ctx)
	} else if m.EventFrequencyPercent.IsKnown() {
		return m.EventFrequencyPercent.MustGet(ctx).ToApi(ctx)
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

	m.FirstSeenEvent = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionFirstSeenEventModel](ctx)
	m.RegressionEvent = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionRegressionEventModel](ctx)
	m.ReappearedEvent = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionReappearedEventModel](ctx)
	m.NewHighPriorityIssue = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionNewHighPriorityIssueModel](ctx)
	m.ExistingHighPriorityIssue = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionExistingHighPriorityIssueModel](ctx)
	m.EventFrequency = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionEventFrequencyModel](ctx)
	m.EventUniqueUserFrequency = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionEventUniqueUserFrequencyModel](ctx)
	m.EventFrequencyPercent = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertConditionEventFrequencyPercentModel](ctx)

	switch conditionValue := conditionValue.(type) {
	case apiclient.ProjectRuleConditionFirstSeenEvent:
		var out IssueAlertConditionFirstSeenEventModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.FirstSeenEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleConditionRegressionEvent:
		var out IssueAlertConditionRegressionEventModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.RegressionEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleConditionReappearedEvent:
		var out IssueAlertConditionReappearedEventModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.ReappearedEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleConditionNewHighPriorityIssue:
		var out IssueAlertConditionNewHighPriorityIssueModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.NewHighPriorityIssue = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleConditionExistingHighPriorityIssue:
		var out IssueAlertConditionExistingHighPriorityIssueModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.ExistingHighPriorityIssue = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleConditionEventFrequency:
		var out IssueAlertConditionEventFrequencyModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.EventFrequency = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleConditionEventUniqueUserFrequency:
		var out IssueAlertConditionEventUniqueUserFrequencyModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.EventUniqueUserFrequency = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleConditionEventFrequencyPercent:
		var out IssueAlertConditionEventFrequencyPercentModel
		diags.Append(out.Fill(ctx, conditionValue)...)
		m.EventFrequencyPercent = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
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

	if filter.Value == nil {
		m.Value = types.StringNull()
	} else if v, err := filter.Value.AsProjectRuleFilterEventAttributeValue0(); err == nil {
		if v == "" {
			m.Value = types.StringNull()
		} else {
			m.Value = types.StringValue(v)
		}
	} else if v, err := filter.Value.AsProjectRuleFilterEventAttributeValue1(); err == nil {
		if v.String() == "" {
			m.Value = types.StringNull()
		} else {
			m.Value = types.StringValue(v.String())
		}
	} else {
		diags.AddError("Invalid event attribute value", fmt.Sprintf("Invalid event attribute value %q. Please report this to the provider developers.", filter.Value))
	}

	return
}

func (m IssueAlertFilterEventAttributeModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleFilter

	var value *apiclient.ProjectRuleFilterEventAttribute_Value
	if !m.Value.IsNull() && !m.Value.IsUnknown() {
		value = &apiclient.ProjectRuleFilterEventAttribute_Value{}
		err := value.FromProjectRuleFilterEventAttributeValue0(m.Value.ValueString())
		if err != nil {
			diags.AddError("Failed to convert to API model", err.Error())
			return nil, diags
		}
	}

	err := v.FromProjectRuleFilterEventAttribute(apiclient.ProjectRuleFilterEventAttribute{
		Name:      m.Name.ValueStringPointer(),
		Attribute: m.Attribute.ValueString(),
		Match:     sentrydata.MatchTypeNameToId[m.Match.ValueString()],
		Value:     value,
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
	AgeComparison        supertypes.SingleNestedObjectValueOf[IssueAlertFilterAgeComparisonModel]        `tfsdk:"age_comparison"`
	IssueOccurrences     supertypes.SingleNestedObjectValueOf[IssueAlertFilterIssueOccurrencesModel]     `tfsdk:"issue_occurrences"`
	AssignedTo           supertypes.SingleNestedObjectValueOf[IssueAlertFilterAssignedToModel]           `tfsdk:"assigned_to"`
	LatestAdoptedRelease supertypes.SingleNestedObjectValueOf[IssueAlertFilterLatestAdoptedReleaseModel] `tfsdk:"latest_adopted_release"`
	LatestRelease        supertypes.SingleNestedObjectValueOf[IssueAlertFilterLatestReleaseModel]        `tfsdk:"latest_release"`
	IssueCategory        supertypes.SingleNestedObjectValueOf[IssueAlertFilterIssueCategoryModel]        `tfsdk:"issue_category"`
	EventAttribute       supertypes.SingleNestedObjectValueOf[IssueAlertFilterEventAttributeModel]       `tfsdk:"event_attribute"`
	TaggedEvent          supertypes.SingleNestedObjectValueOf[IssueAlertFilterTaggedEventModel]          `tfsdk:"tagged_event"`
	Level                supertypes.SingleNestedObjectValueOf[IssueAlertFilterLevelModel]                `tfsdk:"level"`
}

func (m IssueAlertFilterModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleFilter, diag.Diagnostics) {
	if m.AgeComparison.IsKnown() {
		return m.AgeComparison.MustGet(ctx).ToApi(ctx)
	} else if m.IssueOccurrences.IsKnown() {
		return m.IssueOccurrences.MustGet(ctx).ToApi(ctx)
	} else if m.AssignedTo.IsKnown() {
		return m.AssignedTo.MustGet(ctx).ToApi(ctx)
	} else if m.LatestAdoptedRelease.IsKnown() {
		return m.LatestAdoptedRelease.MustGet(ctx).ToApi(ctx)
	} else if m.LatestRelease.IsKnown() {
		return m.LatestRelease.MustGet(ctx).ToApi(ctx)
	} else if m.IssueCategory.IsKnown() {
		return m.IssueCategory.MustGet(ctx).ToApi(ctx)
	} else if m.EventAttribute.IsKnown() {
		return m.EventAttribute.MustGet(ctx).ToApi(ctx)
	} else if m.TaggedEvent.IsKnown() {
		return m.TaggedEvent.MustGet(ctx).ToApi(ctx)
	} else if m.Level.IsKnown() {
		return m.Level.MustGet(ctx).ToApi(ctx)
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

	m.AgeComparison = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterAgeComparisonModel](ctx)
	m.IssueOccurrences = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterIssueOccurrencesModel](ctx)
	m.AssignedTo = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterAssignedToModel](ctx)
	m.LatestAdoptedRelease = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterLatestAdoptedReleaseModel](ctx)
	m.LatestRelease = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterLatestReleaseModel](ctx)
	m.IssueCategory = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterIssueCategoryModel](ctx)
	m.EventAttribute = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterEventAttributeModel](ctx)
	m.TaggedEvent = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterTaggedEventModel](ctx)
	m.Level = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertFilterLevelModel](ctx)

	switch filterValue := filterValue.(type) {
	case apiclient.ProjectRuleFilterAgeComparison:
		var out IssueAlertFilterAgeComparisonModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.AgeComparison = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterIssueOccurrences:
		var out IssueAlertFilterIssueOccurrencesModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.IssueOccurrences = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterAssignedTo:
		var out IssueAlertFilterAssignedToModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.AssignedTo = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterLatestAdoptedRelease:
		var out IssueAlertFilterLatestAdoptedReleaseModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.LatestAdoptedRelease = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterLatestRelease:
		var out IssueAlertFilterLatestReleaseModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.LatestRelease = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterIssueCategory:
		var out IssueAlertFilterIssueCategoryModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.IssueCategory = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterEventAttribute:
		var out IssueAlertFilterEventAttributeModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.EventAttribute = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterTaggedEvent:
		var out IssueAlertFilterTaggedEventModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.TaggedEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleFilterLevel:
		var out IssueAlertFilterLevelModel
		diags.Append(out.Fill(ctx, filterValue)...)
		m.Level = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
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
	Name                      types.String                  `tfsdk:"name"`
	SentryAppInstallationUuid types.String                  `tfsdk:"sentry_app_installation_uuid"`
	Settings                  supertypes.MapValueOf[string] `tfsdk:"settings"`
	SettingsLabels            supertypes.MapValueOf[string] `tfsdk:"settings_labels"`
}

func (m *IssueAlertActionNotifyEventSentryAppModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionNotifyEventSentryApp) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.SentryAppInstallationUuid = types.StringValue(action.SentryAppInstallationUuid)

	if action.Settings == nil {
		m.Settings = supertypes.NewMapValueOfNull[string](ctx)
		m.SettingsLabels = supertypes.NewMapValueOfNull[string](ctx)
	} else {
		var settingsMap = make(map[string]string, len(*action.Settings))
		var labelsMap = make(map[string]string, len(*action.Settings))
		for _, setting := range *action.Settings {
			settingsMap[setting.Name] = setting.Value
			if setting.Label != nil {
				labelsMap[setting.Name] = *setting.Label
			}
		}
		m.Settings = tfutils.MergeDiagnostics(supertypes.NewMapValueOfMap(ctx, settingsMap))(&diags)
		m.SettingsLabels = tfutils.MergeDiagnostics(supertypes.NewMapValueOfMap(ctx, labelsMap))(&diags)
	}
	return
}

func (m IssueAlertActionNotifyEventSentryAppModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction

	var settings *[]struct {
		Label *string `json:"label,omitempty"`
		Name  string  `json:"name"`
		Value string  `json:"value"`
	}

	if m.Settings.IsKnown() {
		elements := tfutils.MergeDiagnostics(m.Settings.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}

		var labels map[string]string
		if m.SettingsLabels.IsKnown() {
			labels = tfutils.MergeDiagnostics(m.SettingsLabels.Get(ctx))(&diags)
			if diags.HasError() {
				return nil, diags
			}
		}

		settings = &[]struct {
			Label *string `json:"label,omitempty"`
			Name  string  `json:"name"`
			Value string  `json:"value"`
		}{}

		for k, val := range elements {
			entry := struct {
				Label *string `json:"label,omitempty"`
				Name  string  `json:"name"`
				Value string  `json:"value"`
			}{
				Name:  k,
				Value: val,
			}
			if lbl, ok := labels[k]; ok {
				entry.Label = &lbl
			}
			*settings = append(*settings, entry)
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
	m.Account = types.StringValue(action.Account.String())
	m.Team = types.StringValue(action.Team)
	m.Priority = types.StringValue(action.Priority)
	return
}

func (m IssueAlertActionOpsgenieNotifyTeam) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionOpsgenieNotifyTeam(apiclient.ProjectRuleActionOpsgenieNotifyTeam{
		Name:     m.Name.ValueStringPointer(),
		Account:  json.Number(m.Account.ValueString()),
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
	m.Account = types.StringValue(action.Account.String())
	m.Service = types.StringValue(action.Service)
	m.Severity = types.StringValue(action.Severity)
	return
}

func (m IssueAlertActionPagerDutyNotifyServiceModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionPagerDutyNotifyService(apiclient.ProjectRuleActionPagerDutyNotifyService{
		Name:     m.Name.ValueStringPointer(),
		Account:  json.Number(m.Account.ValueString()),
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
	Name      types.String             `tfsdk:"name"`
	Workspace types.String             `tfsdk:"workspace"`
	Channel   sentrytypes.SlackChannel `tfsdk:"channel"`
	ChannelId types.String             `tfsdk:"channel_id"`
	Tags      sentrytypes.StringSet    `tfsdk:"tags"`
	Notes     types.String             `tfsdk:"notes"`
}

func (m *IssueAlertActionSlackNotifyServiceModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionSlackNotifyService) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Workspace = types.StringValue(action.Workspace)
	m.Channel = sentrytypes.SlackChannelValue(action.Channel)
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
	m.Server = types.StringValue(action.Server.String())
	m.ChannelId = types.StringValue(action.ChannelId)
	m.Tags = tfutils.MergeDiagnostics(sentrytypes.StringSetPointerValue(action.Tags))(&diags)
	return
}

func (m IssueAlertActionDiscordNotifyServiceModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionDiscordNotifyService(apiclient.ProjectRuleActionDiscordNotifyService{
		Name:      m.Name.ValueStringPointer(),
		Server:    json.Number(m.Server.ValueString()),
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
	m.Integration = types.StringValue(action.Integration.String())
	m.Project = types.StringValue(action.Project)
	m.IssueType = types.StringValue(action.IssueType)
	return
}

func (m IssueAlertActionJiraCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionJiraCreateTicket(apiclient.ProjectRuleActionJiraCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       json.Number(m.Integration.ValueString()),
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
	Name        types.String                  `tfsdk:"name"`
	Integration types.String                  `tfsdk:"integration"`
	Repo        types.String                  `tfsdk:"repo"`
	Assignee    types.String                  `tfsdk:"assignee"`
	Labels      supertypes.SetValueOf[string] `tfsdk:"labels"`
}

func (m *IssueAlertActionGitHubCreateTicketModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionGitHubCreateTicket) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Integration = types.StringValue(action.Integration.String())
	m.Repo = types.StringValue(action.Repo)
	m.Assignee = types.StringPointerValue(action.Assignee)

	if action.Labels == nil {
		m.Labels = supertypes.NewSetValueOfNull[string](ctx)
	} else {
		m.Labels = supertypes.NewSetValueOfSlice(ctx, *action.Labels)
	}
	return
}

func (m IssueAlertActionGitHubCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := apiclient.ProjectRuleActionGitHubCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       json.Number(m.Integration.ValueString()),
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
	Name        types.String                  `tfsdk:"name"`
	Integration types.String                  `tfsdk:"integration"`
	Repo        types.String                  `tfsdk:"repo"`
	Assignee    types.String                  `tfsdk:"assignee"`
	Labels      supertypes.SetValueOf[string] `tfsdk:"labels"`
}

func (m *IssueAlertActionGitHubEnterpriseCreateTicketModel) Fill(ctx context.Context, action apiclient.ProjectRuleActionGitHubEnterpriseCreateTicket) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(action.Name)
	m.Integration = types.StringValue(action.Integration.String())
	m.Repo = types.StringValue(action.Repo)
	m.Assignee = types.StringPointerValue(action.Assignee)

	if action.Labels == nil {
		m.Labels = supertypes.NewSetValueOfNull[string](ctx)
	} else {
		m.Labels = supertypes.NewSetValueOfSlice(ctx, *action.Labels)
	}
	return
}

func (m IssueAlertActionGitHubEnterpriseCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := apiclient.ProjectRuleActionGitHubEnterpriseCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       json.Number(m.Integration.ValueString()),
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
	m.Integration = types.StringValue(action.Integration.String())
	m.Project = types.StringValue(action.Project)
	m.WorkItemType = types.StringValue(action.WorkItemType)
	return
}

func (m IssueAlertActionAzureDevopsCreateTicketModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var v apiclient.ProjectRuleAction
	err := v.FromProjectRuleActionAzureDevopsCreateTicket(apiclient.ProjectRuleActionAzureDevopsCreateTicket{
		Name:              m.Name.ValueStringPointer(),
		Integration:       json.Number(m.Integration.ValueString()),
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
	NotifyEmail                  supertypes.SingleNestedObjectValueOf[IssueAlertActionNotifyEmailModel]                  `tfsdk:"notify_email"`
	NotifyEvent                  supertypes.SingleNestedObjectValueOf[IssueAlertActionNotifyEventModel]                  `tfsdk:"notify_event"`
	NotifyEventService           supertypes.SingleNestedObjectValueOf[IssueAlertActionNotifyEventServiceModel]           `tfsdk:"notify_event_service"`
	NotifyEventSentryApp         supertypes.SingleNestedObjectValueOf[IssueAlertActionNotifyEventSentryAppModel]         `tfsdk:"notify_event_sentry_app"`
	OpsgenieNotifyTeam           supertypes.SingleNestedObjectValueOf[IssueAlertActionOpsgenieNotifyTeam]                `tfsdk:"opsgenie_notify_team"`
	PagerDutyNotifyService       supertypes.SingleNestedObjectValueOf[IssueAlertActionPagerDutyNotifyServiceModel]       `tfsdk:"pagerduty_notify_service"`
	SlackNotifyService           supertypes.SingleNestedObjectValueOf[IssueAlertActionSlackNotifyServiceModel]           `tfsdk:"slack_notify_service"`
	MsTeamsNotifyService         supertypes.SingleNestedObjectValueOf[IssueAlertActionMsTeamsNotifyServiceModel]         `tfsdk:"msteams_notify_service"`
	DiscordNotifyService         supertypes.SingleNestedObjectValueOf[IssueAlertActionDiscordNotifyServiceModel]         `tfsdk:"discord_notify_service"`
	JiraCreateTicket             supertypes.SingleNestedObjectValueOf[IssueAlertActionJiraCreateTicketModel]             `tfsdk:"jira_create_ticket"`
	JiraServerCreateTicket       supertypes.SingleNestedObjectValueOf[IssueAlertActionJiraServerCreateTicketModel]       `tfsdk:"jira_server_create_ticket"`
	GitHubCreateTicket           supertypes.SingleNestedObjectValueOf[IssueAlertActionGitHubCreateTicketModel]           `tfsdk:"github_create_ticket"`
	GitHubEnterpriseCreateTicket supertypes.SingleNestedObjectValueOf[IssueAlertActionGitHubEnterpriseCreateTicketModel] `tfsdk:"github_enterprise_create_ticket"`
	AzureDevopsCreateTicket      supertypes.SingleNestedObjectValueOf[IssueAlertActionAzureDevopsCreateTicketModel]      `tfsdk:"azure_devops_create_ticket"`
}

func (m IssueAlertActionModel) ToApi(ctx context.Context) (*apiclient.ProjectRuleAction, diag.Diagnostics) {
	if m.NotifyEmail.IsKnown() {
		return m.NotifyEmail.MustGet(ctx).ToApi(ctx)
	} else if m.NotifyEvent.IsKnown() {
		return m.NotifyEvent.MustGet(ctx).ToApi(ctx)
	} else if m.NotifyEventService.IsKnown() {
		return m.NotifyEventService.MustGet(ctx).ToApi(ctx)
	} else if m.NotifyEventSentryApp.IsKnown() {
		return m.NotifyEventSentryApp.MustGet(ctx).ToApi(ctx)
	} else if m.OpsgenieNotifyTeam.IsKnown() {
		return m.OpsgenieNotifyTeam.MustGet(ctx).ToApi(ctx)
	} else if m.PagerDutyNotifyService.IsKnown() {
		return m.PagerDutyNotifyService.MustGet(ctx).ToApi(ctx)
	} else if m.SlackNotifyService.IsKnown() {
		return m.SlackNotifyService.MustGet(ctx).ToApi(ctx)
	} else if m.MsTeamsNotifyService.IsKnown() {
		return m.MsTeamsNotifyService.MustGet(ctx).ToApi(ctx)
	} else if m.DiscordNotifyService.IsKnown() {
		return m.DiscordNotifyService.MustGet(ctx).ToApi(ctx)
	} else if m.JiraCreateTicket.IsKnown() {
		return m.JiraCreateTicket.MustGet(ctx).ToApi(ctx)
	} else if m.JiraServerCreateTicket.IsKnown() {
		return m.JiraServerCreateTicket.MustGet(ctx).ToApi(ctx)
	} else if m.GitHubCreateTicket.IsKnown() {
		return m.GitHubCreateTicket.MustGet(ctx).ToApi(ctx)
	} else if m.GitHubEnterpriseCreateTicket.IsKnown() {
		return m.GitHubEnterpriseCreateTicket.MustGet(ctx).ToApi(ctx)
	} else if m.AzureDevopsCreateTicket.IsKnown() {
		return m.AzureDevopsCreateTicket.MustGet(ctx).ToApi(ctx)
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

	m.NotifyEmail = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEmailModel](ctx)
	m.NotifyEvent = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEventModel](ctx)
	m.NotifyEventService = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEventServiceModel](ctx)
	m.NotifyEventSentryApp = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEventSentryAppModel](ctx)
	m.OpsgenieNotifyTeam = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionOpsgenieNotifyTeam](ctx)
	m.PagerDutyNotifyService = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionPagerDutyNotifyServiceModel](ctx)
	m.SlackNotifyService = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionSlackNotifyServiceModel](ctx)
	m.MsTeamsNotifyService = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionMsTeamsNotifyServiceModel](ctx)
	m.DiscordNotifyService = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionDiscordNotifyServiceModel](ctx)
	m.JiraCreateTicket = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionJiraCreateTicketModel](ctx)
	m.JiraServerCreateTicket = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionJiraServerCreateTicketModel](ctx)
	m.GitHubCreateTicket = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionGitHubCreateTicketModel](ctx)
	m.GitHubEnterpriseCreateTicket = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionGitHubEnterpriseCreateTicketModel](ctx)
	m.AzureDevopsCreateTicket = supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionAzureDevopsCreateTicketModel](ctx)

	switch actionValue := actionValue.(type) {
	case apiclient.ProjectRuleActionNotifyEmail:
		var out IssueAlertActionNotifyEmailModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.NotifyEmail = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionNotifyEvent:
		var out IssueAlertActionNotifyEventModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.NotifyEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionNotifyEventService:
		var out IssueAlertActionNotifyEventServiceModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.NotifyEventService = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionNotifyEventSentryApp:
		var out IssueAlertActionNotifyEventSentryAppModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.NotifyEventSentryApp = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionOpsgenieNotifyTeam:
		var out IssueAlertActionOpsgenieNotifyTeam
		diags.Append(out.Fill(ctx, actionValue)...)
		m.OpsgenieNotifyTeam = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionPagerDutyNotifyService:
		var out IssueAlertActionPagerDutyNotifyServiceModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.PagerDutyNotifyService = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionSlackNotifyService:
		var out IssueAlertActionSlackNotifyServiceModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.SlackNotifyService = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionMsTeamsNotifyService:
		var out IssueAlertActionMsTeamsNotifyServiceModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.MsTeamsNotifyService = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionDiscordNotifyService:
		var out IssueAlertActionDiscordNotifyServiceModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.DiscordNotifyService = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionJiraCreateTicket:
		var out IssueAlertActionJiraCreateTicketModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.JiraCreateTicket = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionJiraServerCreateTicket:
		var out IssueAlertActionJiraServerCreateTicketModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.JiraServerCreateTicket = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionGitHubCreateTicket:
		var out IssueAlertActionGitHubCreateTicketModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.GitHubCreateTicket = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionGitHubEnterpriseCreateTicket:
		var out IssueAlertActionGitHubEnterpriseCreateTicketModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.GitHubEnterpriseCreateTicket = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	case apiclient.ProjectRuleActionAzureDevopsCreateTicket:
		var out IssueAlertActionAzureDevopsCreateTicketModel
		diags.Append(out.Fill(ctx, actionValue)...)
		m.AzureDevopsCreateTicket = supertypes.NewSingleNestedObjectValueOf(ctx, &out)
	default:
		diags.AddError("Unsupported action", fmt.Sprintf("Unsupported action type %T", actionValue))
	}

	return
}

// Model

type IssueAlertModel struct {
	Id           types.String                                                 `tfsdk:"id"`
	Organization types.String                                                 `tfsdk:"organization"`
	Project      types.String                                                 `tfsdk:"project"`
	Name         types.String                                                 `tfsdk:"name"`
	Conditions   sentrytypes.LossyJson                                        `tfsdk:"conditions"`
	Filters      sentrytypes.LossyJson                                        `tfsdk:"filters"`
	Actions      sentrytypes.LossyJson                                        `tfsdk:"actions"`
	ActionMatch  types.String                                                 `tfsdk:"action_match"`
	FilterMatch  types.String                                                 `tfsdk:"filter_match"`
	Frequency    types.Int64                                                  `tfsdk:"frequency"`
	Environment  types.String                                                 `tfsdk:"environment"`
	Owner        types.String                                                 `tfsdk:"owner"`
	ConditionsV2 supertypes.ListNestedObjectValueOf[IssueAlertConditionModel] `tfsdk:"conditions_v2"`
	FiltersV2    supertypes.ListNestedObjectValueOf[IssueAlertFilterModel]    `tfsdk:"filters_v2"`
	ActionsV2    supertypes.ListNestedObjectValueOf[IssueAlertActionModel]    `tfsdk:"actions_v2"`
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
		apiConditions := lo.Map(alert.Conditions, func(c apiclient.ProjectRuleCondition, _ int) json.RawMessage {
			b, _ := json.Marshal(c)
			return b
		})
		var priorConditions []json.RawMessage
		if !m.Conditions.IsUnknown() {
			_ = json.Unmarshal([]byte(m.Conditions.ValueString()), &priorConditions)
		}
		apiConditions = reorderToMatchPrior(priorConditions, apiConditions, legacyJsonItemKey)
		if b, err := json.Marshal(apiConditions); err == nil {
			m.Conditions = sentrytypes.NewLossyJsonValue(string(b))
		} else {
			diags.AddError("Invalid conditions", err.Error())
			return
		}
	} else if !m.ConditionsV2.IsNull() {
		var priorConditions []IssueAlertConditionModel
		if !m.ConditionsV2.IsUnknown() {
			diags.Append(m.ConditionsV2.ElementsAs(ctx, &priorConditions, false)...)
			if diags.HasError() {
				return
			}
		}

		conditions := lo.Map(alert.Conditions, func(condition apiclient.ProjectRuleCondition, _ int) IssueAlertConditionModel {
			var conditionModel IssueAlertConditionModel
			diags.Append(conditionModel.Fill(ctx, condition)...)
			return conditionModel
		})

		if diags.HasError() {
			return
		}

		conditions = reorderToMatchPrior(priorConditions, conditions, issueAlertConditionModelKey(ctx))

		m.ConditionsV2 = supertypes.NewListNestedObjectValueOfValueSlice(ctx, conditions)
	}

	if !m.Filters.IsNull() {
		apiFilters := lo.Map(alert.Filters, func(f apiclient.ProjectRuleFilter, _ int) json.RawMessage {
			b, _ := json.Marshal(f)
			return b
		})
		var priorFilters []json.RawMessage
		if !m.Filters.IsUnknown() {
			_ = json.Unmarshal([]byte(m.Filters.ValueString()), &priorFilters)
		}
		apiFilters = reorderToMatchPrior(priorFilters, apiFilters, legacyJsonItemKey)
		if b, err := json.Marshal(apiFilters); err == nil {
			m.Filters = sentrytypes.NewLossyJsonValue(string(b))
		} else {
			diags.AddError("Invalid filters", err.Error())
		}
	} else if !m.FiltersV2.IsNull() {
		var priorFilters []IssueAlertFilterModel
		if !m.FiltersV2.IsUnknown() {
			diags.Append(m.FiltersV2.ElementsAs(ctx, &priorFilters, false)...)
			if diags.HasError() {
				return
			}
		}

		filters := lo.Map(alert.Filters, func(filter apiclient.ProjectRuleFilter, _ int) IssueAlertFilterModel {
			var filterModel IssueAlertFilterModel
			diags.Append(filterModel.Fill(ctx, filter)...)
			return filterModel
		})

		if diags.HasError() {
			return
		}

		filters = reorderToMatchPrior(priorFilters, filters, issueAlertFilterModelKey(ctx))

		m.FiltersV2 = supertypes.NewListNestedObjectValueOfValueSlice(ctx, filters)
	}

	if !m.Actions.IsNull() {
		apiActions := lo.Map(alert.Actions, func(a apiclient.ProjectRuleAction, _ int) json.RawMessage {
			b, _ := json.Marshal(a)
			b, _ = stripLegacyActionDisplayFields(b)
			return b
		})
		var priorActions []json.RawMessage
		if !m.Actions.IsUnknown() {
			_ = json.Unmarshal([]byte(m.Actions.ValueString()), &priorActions)
		}
		apiActions = reorderToMatchPrior(priorActions, apiActions, legacyJsonItemKey)
		if b, err := json.Marshal(apiActions); err != nil {
			diags.AddError("Invalid actions", err.Error())
		} else {
			m.Actions = sentrytypes.NewLossyJsonValue(string(b))
		}
	} else if !m.ActionsV2.IsNull() {
		var priorActions []IssueAlertActionModel
		if !m.ActionsV2.IsUnknown() {
			diags.Append(m.ActionsV2.ElementsAs(ctx, &priorActions, false)...)
			if diags.HasError() {
				return
			}
		}

		actions := lo.Map(alert.Actions, func(action apiclient.ProjectRuleAction, _ int) IssueAlertActionModel {
			var actionModel IssueAlertActionModel
			diags.Append(actionModel.Fill(ctx, action)...)
			return actionModel
		})

		if diags.HasError() {
			return
		}

		actions = reorderToMatchPrior(priorActions, actions, issueAlertActionModelKey(ctx))

		m.ActionsV2 = supertypes.NewListNestedObjectValueOfValueSlice(ctx, actions)
	}

	return
}

// stripLegacyActionDisplayFields removes API-injected display-only keys from the legacy
// action JSON before storing in state. The Sentry GET endpoint unconditionally adds
// "formFields" (live webhook schema) and "name" (generated label) to every action dict.
// These fields are never stored server-side and cannot be stabilized in config.
func stripLegacyActionDisplayFields(actionJSON []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(actionJSON))
	dec.UseNumber()

	var action map[string]interface{}
	if err := dec.Decode(&action); err != nil {
		return nil, err
	}

	delete(action, "formFields")
	delete(action, "name")
	delete(action, "hasSchemaFormConfig")

	return json.Marshal(action)
}
