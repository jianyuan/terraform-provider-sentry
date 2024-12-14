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
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
)

// Conditions

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

	panic("provider error: unsupported condition")
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

func (m IssueAlertFilterAgeComparisonModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterAgeComparison(apiclient.ProjectRuleFilterAgeComparison{
		Name:           m.Name.ValueStringPointer(),
		ComparisonType: m.ComparisonType.ValueString(),
		Value:          m.Value.ValueInt64(),
		Time:           m.Time.ValueString(),
	}))
	return v
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

func (m IssueAlertFilterIssueOccurrencesModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterIssueOccurrences(apiclient.ProjectRuleFilterIssueOccurrences{
		Name:  m.Name.ValueStringPointer(),
		Value: m.Value.ValueInt64(),
	}))
	return v
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

func (m IssueAlertFilterAssignedToModel) ToApi() apiclient.ProjectRuleFilter {
	var targetIdentifier *apiclient.ProjectRuleFilterAssignedTo_TargetIdentifier

	if !m.TargetIdentifier.IsNull() {
		targetIdentifier = &apiclient.ProjectRuleFilterAssignedTo_TargetIdentifier{}
		must.Do(targetIdentifier.FromProjectRuleFilterAssignedToTargetIdentifier0(m.TargetIdentifier.ValueString()))
	}

	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterAssignedTo(apiclient.ProjectRuleFilterAssignedTo{
		Name:             m.Name.ValueStringPointer(),
		TargetType:       m.TargetType.ValueString(),
		TargetIdentifier: targetIdentifier,
	}))
	return v
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

func (m IssueAlertFilterLatestAdoptedReleaseModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterLatestAdoptedRelease(apiclient.ProjectRuleFilterLatestAdoptedRelease{
		Name:           m.Name.ValueStringPointer(),
		OldestOrNewest: m.OldestOrNewest.ValueString(),
		OlderOrNewer:   m.OlderOrNewer.ValueString(),
		Environment:    m.Environment.ValueString(),
	}))
	return v
}

type IssueAlertFilterLatestReleaseModel struct {
	Name types.String `tfsdk:"name"`
}

func (m *IssueAlertFilterLatestReleaseModel) Fill(ctx context.Context, filter apiclient.ProjectRuleFilterLatestRelease) (diags diag.Diagnostics) {
	m.Name = types.StringPointerValue(filter.Name)
	return
}

func (m IssueAlertFilterLatestReleaseModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterLatestRelease(apiclient.ProjectRuleFilterLatestRelease{
		Name: m.Name.ValueStringPointer(),
	}))
	return v
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

func (m IssueAlertFilterIssueCategoryModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterIssueCategory(apiclient.ProjectRuleFilterIssueCategory{
		Name:  m.Name.ValueStringPointer(),
		Value: sentrydata.IssueGroupCategoryNameToId[m.Value.ValueString()],
	}))
	return v
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

func (m IssueAlertFilterEventAttributeModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterEventAttribute(apiclient.ProjectRuleFilterEventAttribute{
		Name:      m.Name.ValueStringPointer(),
		Attribute: m.Attribute.ValueString(),
		Match:     sentrydata.MatchTypeNameToId[m.Match.ValueString()],
		Value:     m.Value.ValueStringPointer(),
	}))
	return v
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

func (m IssueAlertFilterTaggedEventModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterTaggedEvent(apiclient.ProjectRuleFilterTaggedEvent{
		Name:  m.Name.ValueStringPointer(),
		Key:   m.Key.ValueString(),
		Match: sentrydata.MatchTypeNameToId[m.Match.ValueString()],
		Value: m.Value.ValueStringPointer(),
	}))
	return v
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

func (m IssueAlertFilterLevelModel) ToApi() apiclient.ProjectRuleFilter {
	var v apiclient.ProjectRuleFilter
	must.Do(v.FromProjectRuleFilterLevel(apiclient.ProjectRuleFilterLevel{
		Name:  m.Name.ValueStringPointer(),
		Match: sentrydata.MatchTypeNameToId[m.Match.ValueString()],
		Level: sentrydata.LogLevelNameToId[m.Level.ValueString()],
	}))
	return v
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

func (m IssueAlertFilterModel) ToApi() apiclient.ProjectRuleFilter {
	if m.AgeComparison != nil {
		return m.AgeComparison.ToApi()
	} else if m.IssueOccurrences != nil {
		return m.IssueOccurrences.ToApi()
	} else if m.AssignedTo != nil {
		return m.AssignedTo.ToApi()
	} else if m.LatestAdoptedRelease != nil {
		return m.LatestAdoptedRelease.ToApi()
	} else if m.LatestRelease != nil {
		return m.LatestRelease.ToApi()
	} else if m.IssueCategory != nil {
		return m.IssueCategory.ToApi()
	} else if m.EventAttribute != nil {
		return m.EventAttribute.ToApi()
	} else if m.TaggedEvent != nil {
		return m.TaggedEvent.ToApi()
	} else if m.Level != nil {
		return m.Level.ToApi()
	}

	panic("provider error: unsupported filter")
}

func (m *IssueAlertFilterModel) FromApi(ctx context.Context, filter apiclient.ProjectRuleFilter) (diags diag.Diagnostics) {
	filterValue, err := filter.ValueByDiscriminator()
	if err != nil {
		diags.AddError("Invalid filter", err.Error())
		return
	}

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

// Model

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
	FiltersV2    *[]IssueAlertFilterModel    `tfsdk:"filters_v2"`
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

	if !m.Filters.IsNull() {
		m.Filters = sentrytypes.NewLossyJsonNull()
		if len(alert.Filters) > 0 {
			if filters, err := json.Marshal(alert.Filters); err == nil {
				m.Filters = sentrytypes.NewLossyJsonValue(string(filters))
			} else {
				diags.AddError("Invalid filters", err.Error())
			}
		}
	} else if m.FiltersV2 != nil {
		m.FiltersV2 = ptr.Ptr(sliceutils.Map(func(filter apiclient.ProjectRuleFilter) IssueAlertFilterModel {
			var filterModel IssueAlertFilterModel
			diags.Append(filterModel.FromApi(ctx, filter)...)
			return filterModel
		}, alert.Filters))

		if diags.HasError() {
			return
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
