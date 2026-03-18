package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *AlertResource) getCreateJSONRequestBody(ctx context.Context, data AlertResourceModel) (*apiclient.CreateOrganizationAlertJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	monitorIds := data.MonitorIds.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}

	inTriggerConditions := data.TriggerConditions.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}
	var outTriggerConditions []apiclient.OrganizationWorkflowTriggerCondition
	for _, triggerCondition := range inTriggerConditions {
		outTriggerConditions = append(outTriggerConditions, apiclient.OrganizationWorkflowTriggerCondition{
			Type:            triggerCondition,
			Comparison:      true,
			ConditionResult: true,
		})
	}

	req := apiclient.CreateOrganizationAlertJSONRequestBody{
		Name:        data.Name.Get(),
		Enabled:     data.Enabled.Get(),
		Environment: data.Environment.Get(),
		Config: apiclient.OrganizationWorkflowConfig{
			Frequency: data.FrequencyMinutes.Get(),
		},
		DetectorIds: monitorIds,
		Triggers: apiclient.OrganizationWorkflowTrigger{
			LogicType:  apiclient.OrganizationWorkflowTriggerLogicTypeAnyShort,
			Conditions: outTriggerConditions,
		},
		AdditionalProperties: map[string]any{
			"actionFilters": []map[string]any{},
		},
	}

	return &req, nil
}

func (m *AlertResourceModel) Fill(ctx context.Context, data apiclient.OrganizationWorkflow) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringValue(data.Name)
	m.Enabled = supertypes.NewBoolValue(data.Enabled)
	m.Environment = supertypes.NewStringValue(data.Environment)
	m.FrequencyMinutes = supertypes.NewInt64Value(data.Config.Frequency)
	m.MonitorIds = supertypes.NewSetValueOfSlice(ctx, data.DetectorIds)

	var triggerConditions []string
	for _, triggerCondition := range data.Triggers.Conditions {
		triggerConditions = append(triggerConditions, triggerCondition.Type)
	}
	m.TriggerConditions = supertypes.NewSetValueOfSlice(ctx, triggerConditions)

	return
}
