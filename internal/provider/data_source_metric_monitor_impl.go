package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

func (m *MetricMonitorDataSourceModel) fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	if data.Owner.IsSpecified() && !data.Owner.IsNull() {
		inOwner, err := data.Owner.Get()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		inOwnerValue, err := inOwner.ValueByDiscriminator()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		outOwner := &MetricMonitorDataSourceModelOwner{}

		switch inOwnerValue := inOwnerValue.(type) {
		case apiclient.ProjectMonitorOwnerUser:
			outOwner.UserId.Set(inOwnerValue.Id)
			diags.Append(m.Owner.Set(ctx, outOwner)...)
		case apiclient.ProjectMonitorOwnerTeam:
			outOwner.TeamId.Set(inOwnerValue.Id)
			diags.Append(m.Owner.Set(ctx, outOwner)...)
		default:
			m.Owner.SetNull(ctx)
		}
	} else {
		m.Owner.SetNull(ctx)
	}

	outConditions := make([]*MetricMonitorDataSourceModelConditionGroupConditionsItem, 0, len(data.ConditionGroup.Conditions))
	for _, inCondition := range data.ConditionGroup.Conditions {
		var outCondition MetricMonitorDataSourceModelConditionGroupConditionsItem
		outCondition.Type.Set(inCondition.Type)

		if inComparison, err := inCondition.Comparison.AsProjectMonitorConditionGroupConditionComparison1(); err == nil {
			outCondition.Comparison.Set(inComparison)
		} else if inComparison, err := inCondition.Comparison.AsProjectMonitorConditionGroupConditionComparison2(); err == nil {
			outCondition.ComparisonSensitivity.Set(inComparison.Sensitivity)
			outCondition.ComparisonThresholdType.Set(sentrydata.AlertRuleThresholdTypeIdToName[inComparison.ThresholdType])
		} else {
			diags.AddError("Invalid comparison", "Unable to unmarshal comparison")
			return
		}

		outCondition.ConditionResult.Set(inCondition.ConditionResult)
		outConditions = append(outConditions, &outCondition)
	}

	var conditionGroup MetricMonitorDataSourceModelConditionGroup
	conditionGroup.LogicType.Set(string(data.ConditionGroup.LogicType))
	diags.Append(conditionGroup.Conditions.Set(ctx, outConditions)...)
	if diags.HasError() {
		return
	}

	diags.Append(m.ConditionGroup.Set(ctx, &conditionGroup)...)
	if diags.HasError() {
		return
	}

	if len(data.DataSources) != 1 {
		diags.AddError("Invalid data source", fmt.Sprintf("Expected 1 data source, got %d", len(data.DataSources)))
		return
	}

	dataSource, err := data.DataSources[0].AsProjectMonitorDataSourceWrapperSnubaQuerySubscription()
	if err != nil {
		diags.AddError("Error unmarshalling JSON", err.Error())
	}

	m.Aggregate.Set(dataSource.QueryObj.SnubaQuery.Aggregate)
	m.Dataset.Set(dataSource.QueryObj.SnubaQuery.Dataset)
	if v, err := dataSource.QueryObj.SnubaQuery.Environment.Get(); err == nil {
		m.Environment.Set(v)
	} else {
		m.Environment.SetNull()
	}

	diags.Append(m.EventTypes.Set(ctx, dataSource.QueryObj.SnubaQuery.EventTypes)...)
	if diags.HasError() {
		return
	}

	if v, err := dataSource.QueryObj.SnubaQuery.Query.Get(); err == nil {
		m.Query.Set(v)
	} else {
		m.Query.SetNull()
	}
	if v, err := dataSource.QueryObj.SnubaQuery.QueryType.Get(); err == nil {
		m.QueryType.Set(sentrydata.SnubaQueryTypeIdToName[v])
	} else {
		// BUG?
		m.QueryType.Set(sentrydata.SnubaQueryTypeIdToName[0])
	}
	if v, err := dataSource.QueryObj.SnubaQuery.TimeWindow.Get(); err == nil {
		m.TimeWindowSeconds.Set(v)
	} else {
		m.TimeWindowSeconds.SetNull()
	}
	if v, err := dataSource.QueryObj.SnubaQuery.ExtrapolationMode.Get(); err == nil {
		m.ExtrapolationMode.Set(v)
	} else {
		m.ExtrapolationMode.SetNull()
	}

	if inConfig, err := data.Config.AsProjectMonitorConfigMetricIssue(); err == nil {
		var issueDetection MetricMonitorDataSourceModelIssueDetection
		issueDetection.Type.SetPtr(inConfig.DetectionType)
		issueDetection.ComparisonDelta.SetPtr(inConfig.ComparisonDelta)

		diags.Append(m.IssueDetection.Set(ctx, &issueDetection)...)
		if diags.HasError() {
			return
		}
	} else {
		diags.AddError("Invalid config", "Invalid config")
		m.IssueDetection.SetNull(ctx)
	}

	return
}
