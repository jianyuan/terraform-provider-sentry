package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/oapi-codegen/nullable"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *MetricMonitorResource) getCreateJSONRequestBody(ctx context.Context, data MetricMonitorResourceModel) (*apiclient.CreateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	outDs := apiclient.ProjectMonitorDataSourceSnubaQuerySubscription{
		Aggregate:  data.Aggregate.Get(),
		Dataset:    data.Dataset.Get(),
		EventTypes: data.EventTypes.DiagsGet(ctx, diags),
	}
	if data.Environment.IsKnown() {
		outDs.Environment = nullable.NewNullableWithValue(data.Environment.Get())
	} else {
		outDs.Environment = nullable.NewNullNullable[string]()
	}
	if data.ExtrapolationMode.IsKnown() {
		outDs.ExtrapolationMode = nullable.NewNullableWithValue(data.ExtrapolationMode.Get())
	} else {
		outDs.ExtrapolationMode = nullable.NewNullNullable[string]()
	}
	if diags.HasError() {
		return nil, diags
	}

	inConditionGroup := data.ConditionGroup.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}

	inConditions := inConditionGroup.Conditions.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}
	var outConditions []apiclient.ProjectMonitorConditionGroupCondition
	for _, inCondition := range inConditions {
		var outComparison apiclient.ProjectMonitorConditionGroupCondition_Comparison
		if err := outComparison.FromProjectMonitorConditionGroupConditionComparison1(inCondition.Comparison.Get()); err != nil {
			diags.AddError("Error marshalling JSON", err.Error())
			return nil, diags
		}
		outConditions = append(outConditions, apiclient.ProjectMonitorConditionGroupCondition{
			Type:            inCondition.Type.Get(),
			Comparison:      outComparison,
			ConditionResult: inCondition.ConditionResult.Get(),
		})
	}

	inConfig := data.IssueDetection.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}

	var outConfig apiclient.ProjectMonitorConfig
	if err := outConfig.FromProjectMonitorConfigMetricIssue(apiclient.ProjectMonitorConfigMetricIssue{
		DetectionType:   inConfig.Type.GetPtr(),
		ComparisonDelta: inConfig.ComparisonDelta.GetPtr(),
	}); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}

	out := apiclient.ProjectMonitorRequestMetricIssue{
		Name:      data.Name.Get(),
		ProjectId: data.Project.Get(),
		DataSources: []apiclient.ProjectMonitorDataSourceSnubaQuerySubscription{
			outDs,
		},
		ConditionGroup: apiclient.ProjectMonitorConditionGroup{
			LogicType:  apiclient.ProjectMonitorConditionGroupLogicType(inConditionGroup.LogicType.Get()),
			Conditions: outConditions,
		},
		Config: &outConfig,
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
	if err := req.FromProjectMonitorRequestMetricIssue(out); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}
	return &req, nil
}

func (m *MetricMonitorResourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
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

		defaultAssignee := &MetricMonitorResourceModelDefaultAssignee{}

		switch ownerValue := ownerValue.(type) {
		case apiclient.ProjectMonitorOwnerUser:
			defaultAssignee.UserId = supertypes.NewStringValue(ownerValue.Id)
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOf(ctx, defaultAssignee)
		case apiclient.ProjectMonitorOwnerTeam:
			defaultAssignee.TeamId = supertypes.NewStringValue(ownerValue.Id)
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOf(ctx, defaultAssignee)
		default:
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOfNull[MetricMonitorResourceModelDefaultAssignee](ctx)
		}
	} else {
		m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOfNull[MetricMonitorResourceModelDefaultAssignee](ctx)
	}

	var outConditions []MetricMonitorResourceModelConditionGroupConditionsItem
	for _, inCondition := range data.ConditionGroup.Conditions {
		inComparison, err := inCondition.Comparison.AsProjectMonitorConditionGroupConditionComparison1()
		if err != nil {
			diags.AddError("Invalid comparison", err.Error())
			return
		}
		outConditions = append(outConditions, MetricMonitorResourceModelConditionGroupConditionsItem{
			Type:            supertypes.NewStringValue(inCondition.Type),
			Comparison:      supertypes.NewInt64Value(inComparison),
			ConditionResult: supertypes.NewInt64Value(inCondition.ConditionResult),
		})
	}

	m.ConditionGroup = supertypes.NewSingleNestedObjectValueOf(ctx, &MetricMonitorResourceModelConditionGroup{
		LogicType:  supertypes.NewStringValue(string(data.ConditionGroup.LogicType)),
		Conditions: supertypes.NewListNestedObjectValueOfValueSlice(ctx, outConditions),
	})

	if len(data.DataSources) != 1 {
		diags.AddError("Invalid data source", fmt.Sprintf("Expected 1 data source, got %d", len(data.DataSources)))
		return
	}

	dataSource, err := data.DataSources[0].AsProjectMonitorDataSourceWrapperSnubaQuerySubscription()
	if err != nil {
		diags.AddError("Error unmarshalling JSON", err.Error())
	}

	m.Aggregate = supertypes.NewStringValue(dataSource.QueryObj.SnubaQuery.Aggregate)
	m.Dataset = supertypes.NewStringValue(dataSource.QueryObj.SnubaQuery.Dataset)
	if v, err := dataSource.QueryObj.SnubaQuery.Environment.Get(); err == nil {
		m.Environment = supertypes.NewStringValue(v)
	} else {
		m.Environment = supertypes.NewStringNull()
	}
	m.EventTypes = supertypes.NewSetValueOfSlice(ctx, dataSource.QueryObj.SnubaQuery.EventTypes)
	if v, err := dataSource.QueryObj.SnubaQuery.ExtrapolationMode.Get(); err == nil {
		m.ExtrapolationMode = supertypes.NewStringValue(v)
	} else {
		m.ExtrapolationMode = supertypes.NewStringNull()
	}

	if inConfig, err := data.Config.AsProjectMonitorConfigMetricIssue(); err == nil {
		m.IssueDetection = supertypes.NewSingleNestedObjectValueOf(ctx, &MetricMonitorResourceModelIssueDetection{
			Type:            supertypes.NewStringPointerValueOrNull(inConfig.DetectionType),
			ComparisonDelta: supertypes.NewInt64PointerValueOrNull(inConfig.ComparisonDelta),
		})
	} else {
		diags.AddError("Invalid config", "Invalid config")
		m.IssueDetection = supertypes.NewSingleNestedObjectValueOfNull[MetricMonitorResourceModelIssueDetection](ctx)
	}

	return
}
