package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func (r *MetricMonitorResource) getCreateJSONRequestBody(ctx context.Context, data MetricMonitorResourceModel) (*apiclient.CreateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	outDs := apiclient.ProjectMonitorDataSourceSnubaQuerySubscription{
		Aggregate:  data.Aggregate.Get(),
		Dataset:    data.Dataset.Get(),
		EventTypes: tfutils.MergeDiagnostics(data.EventTypes.Get(ctx))(&diags),
	}
	if data.Environment.IsKnown() {
		outDs.Environment.Set(data.Environment.Get())
	} else {
		outDs.Environment.SetNull()
	}
	if data.ExtrapolationMode.IsKnown() {
		outDs.ExtrapolationMode.Set(data.ExtrapolationMode.Get())
	} else {
		outDs.ExtrapolationMode.SetNull()
	}
	if diags.HasError() {
		return nil, diags
	}

	inConditionGroup := tfutils.MergeDiagnostics(data.ConditionGroup.Get(ctx))(&diags)
	if diags.HasError() {
		return nil, diags
	}

	inConditions := tfutils.MergeDiagnostics(inConditionGroup.Conditions.Get(ctx))(&diags)
	if diags.HasError() {
		return nil, diags
	}

	outConditions := make([]apiclient.ProjectMonitorConditionGroupCondition, 0, len(inConditions))
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

	inConfig := tfutils.MergeDiagnostics(data.IssueDetection.Get(ctx))(&diags)
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
		out.Enabled.Set(data.Enabled.Get())
	} else {
		out.Enabled.SetNull()
	}

	if data.Description.IsKnown() {
		out.Description.Set(data.Description.Get())
	} else {
		out.Description.SetNull()
	}

	if data.Owner.IsKnown() {
		owner := tfutils.MergeDiagnostics(data.Owner.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}

		switch {
		case owner.TeamId.IsKnown():
			out.Owner.Set(fmt.Sprintf("team:%s", owner.TeamId.Get()))
		case owner.UserId.IsKnown():
			out.Owner.Set(fmt.Sprintf("user:%s", owner.UserId.Get()))
		default:
			out.Owner.SetNull()
		}
	} else {
		out.Owner.SetNull()
	}

	var req apiclient.CreateProjectMonitorJSONRequestBody
	if err := req.FromProjectMonitorRequestMetricIssue(out); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}
	return &req, nil
}

func (r *MetricMonitorResource) getUpdateJSONRequestBody(ctx context.Context, data MetricMonitorResourceModel) (*apiclient.UpdateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	return r.getCreateJSONRequestBody(ctx, data)
}

func (m *MetricMonitorResourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	m.Id.Set(data.Id)
	m.Name.Set(data.Name)
	if v, err := data.Description.Get(); err == nil {
		m.Description.Set(v)
	} else {
		m.Description.SetNull()
	}
	m.Enabled.Set(data.Enabled)

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

		outOwner := &MetricMonitorResourceModelOwner{}

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

	outConditions := make([]*MetricMonitorResourceModelConditionGroupConditionsItem, 0, len(data.ConditionGroup.Conditions))
	for _, inCondition := range data.ConditionGroup.Conditions {
		inComparison, err := inCondition.Comparison.AsProjectMonitorConditionGroupConditionComparison1()
		if err != nil {
			diags.AddError("Invalid comparison", err.Error())
			return
		}

		var outCondition MetricMonitorResourceModelConditionGroupConditionsItem
		outCondition.Type.Set(inCondition.Type)
		outCondition.Comparison.Set(inComparison)
		outCondition.ConditionResult.Set(inCondition.ConditionResult)
		outConditions = append(outConditions, &outCondition)
	}

	var conditionGroup MetricMonitorResourceModelConditionGroup
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

	if v, err := dataSource.QueryObj.SnubaQuery.ExtrapolationMode.Get(); err == nil {
		m.ExtrapolationMode.Set(v)
	} else {
		m.ExtrapolationMode.SetNull()
	}

	if inConfig, err := data.Config.AsProjectMonitorConfigMetricIssue(); err == nil {
		var issueDetection MetricMonitorResourceModelIssueDetection
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
