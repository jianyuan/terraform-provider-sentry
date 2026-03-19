package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/go-utils/ptr"
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

	inActionFilters := data.ActionFilters.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}
	var outActionFilters []apiclient.OrganizationWorkflowActionFilter
	for _, inActionFilter := range inActionFilters {
		inActions := inActionFilter.Actions.DiagsGet(ctx, diags)
		if diags.HasError() {
			return nil, diags
		}

		var outActions []apiclient.OrganizationWorkflowActionFilterAction
		for _, inAction := range inActions {
			var outAction apiclient.OrganizationWorkflowActionFilterAction
			switch {
			case inAction.Email.IsKnown():
				inEmail := inAction.Email.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEmail apiclient.OrganizationWorkflowActionFilterActionEmail
				outEmail.Config.TargetType = apiclient.OrganizationWorkflowActionFilterActionEmailConfigTargetType(inEmail.TargetType.Get())
				outEmail.Config.TargetIdentifier = inEmail.TargetId.GetPtr()
				if inEmail.FallthroughType.IsKnown() {
					outEmail.Data.FallthroughType = ptr.Ptr(apiclient.OrganizationWorkflowActionFilterActionEmailDataFallthroughType(inEmail.FallthroughType.Get()))
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionEmail(outEmail); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Plugin.IsKnown():
				if err := outAction.FromOrganizationWorkflowActionFilterActionPlugin(apiclient.OrganizationWorkflowActionFilterActionPlugin{
					Data:   map[string]any{},
					Config: map[string]any{},
				}); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Slack.IsKnown():
				inSlack := inAction.Slack.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outSlack apiclient.OrganizationWorkflowActionFilterActionSlack
				outSlack.IntegrationId = inSlack.IntegrationId.Get()
				outSlack.Config.TargetType = "specific"
				outSlack.Config.TargetIdentifier = inSlack.ChannelId.Get()
				outSlack.Config.TargetDisplay = inSlack.ChannelName.Get()
				if inSlack.Tags.IsKnown() {
					outSlack.Data.Tags = ptr.Ptr(inSlack.Tags.Get())
				}
				if inSlack.Notes.IsKnown() {
					outSlack.Data.Notes = ptr.Ptr(inSlack.Notes.Get())
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionSlack(outSlack); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Pagerduty.IsKnown():
				inPagerduty := inAction.Pagerduty.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outPagerduty apiclient.OrganizationWorkflowActionFilterActionPagerDuty
				outPagerduty.IntegrationId = inPagerduty.IntegrationId.Get()
				outPagerduty.Config.TargetType = "specific"
				outPagerduty.Config.TargetIdentifier = inPagerduty.ServiceId.Get()
				outPagerduty.Config.TargetDisplay = inPagerduty.ServiceName.Get()
				outPagerduty.Data.Priority = inPagerduty.Severity.GetPtr()

				if err := outAction.FromOrganizationWorkflowActionFilterActionPagerDuty(outPagerduty); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Discord.IsKnown():
				inDiscord := inAction.Discord.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outDiscord apiclient.OrganizationWorkflowActionFilterActionDiscord
				outDiscord.IntegrationId = inDiscord.IntegrationId.Get()
				outDiscord.Config.TargetType = "specific"
				outDiscord.Config.TargetIdentifier = inDiscord.ChannelId.Get()
				if inDiscord.Tags.IsKnown() {
					outDiscord.Data.Tags = ptr.Ptr(inDiscord.Tags.Get())
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionDiscord(outDiscord); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Msteams.IsKnown():
				inMsteams := inAction.Msteams.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outMsteams apiclient.OrganizationWorkflowActionFilterActionMsTeams
				outMsteams.IntegrationId = inMsteams.IntegrationId.Get()
				outMsteams.Config.TargetType = "specific"
				outMsteams.Config.TargetIdentifier = inMsteams.TeamId.Get()
				outMsteams.Config.TargetDisplay = inMsteams.ChannelName.Get()

				if err := outAction.FromOrganizationWorkflowActionFilterActionMsTeams(outMsteams); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Opsgenie.IsKnown():
				inOpsgenie := inAction.Opsgenie.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outOpsgenie apiclient.OrganizationWorkflowActionFilterActionOpsgenie
				outOpsgenie.IntegrationId = inOpsgenie.IntegrationId.Get()
				outOpsgenie.Config.TargetType = "specific"
				outOpsgenie.Config.TargetIdentifier = inOpsgenie.TeamId.Get()
				outOpsgenie.Config.TargetDisplay = inOpsgenie.TeamName.Get()
				outOpsgenie.Data.Priority = inOpsgenie.Priority.GetPtr()

				if err := outAction.FromOrganizationWorkflowActionFilterActionOpsgenie(outOpsgenie); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Vsts.IsKnown():
				inVsts := inAction.Vsts.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outVsts apiclient.OrganizationWorkflowActionFilterActionVsts
				outVsts.IntegrationId = inVsts.IntegrationId.Get()
				outVsts.Config.TargetType = "specific"
				if inVsts.Data.IsKnown() {
					outVsts.Data = inVsts.Data.DiagsGet(ctx, diags)
				} else {
					outVsts.Data = map[string]string{}
				}
				if diags.HasError() {
					return nil, diags
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionVsts(outVsts); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Jira.IsKnown():
				inJira := inAction.Jira.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outJira apiclient.OrganizationWorkflowActionFilterActionJira
				outJira.IntegrationId = inJira.IntegrationId.Get()
				outJira.Config.TargetType = "specific"
				if inJira.Data.IsKnown() {
					outJira.Data = inJira.Data.DiagsGet(ctx, diags)
				} else {
					outJira.Data = map[string]string{}
				}
				if diags.HasError() {
					return nil, diags
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionJira(outJira); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.JiraServer.IsKnown():
				inJiraServer := inAction.JiraServer.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outJiraServer apiclient.OrganizationWorkflowActionFilterActionJiraServer
				outJiraServer.IntegrationId = inJiraServer.IntegrationId.Get()
				outJiraServer.Config.TargetType = "specific"
				if inJiraServer.Data.IsKnown() {
					outJiraServer.Data = inJiraServer.Data.DiagsGet(ctx, diags)
				} else {
					outJiraServer.Data = map[string]string{}
				}
				if diags.HasError() {
					return nil, diags
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionJiraServer(outJiraServer); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Github.IsKnown():
				inGithub := inAction.Github.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outGithub apiclient.OrganizationWorkflowActionFilterActionGitHub
				outGithub.IntegrationId = inGithub.IntegrationId.Get()
				outGithub.Config.TargetType = "specific"
				outGithub.Data.AdditionalFields.Repo = inGithub.Repo.Get()
				outGithub.Data.AdditionalFields.Assignee = inGithub.Assignee.Get()
				outGithub.Data.AdditionalFields.Labels = inGithub.Labels.DiagsGet(ctx, diags)
				outGithub.Data.AdditionalFields.Integration = inGithub.IntegrationId.Get()
				if diags.HasError() {
					return nil, diags
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionGitHub(outGithub); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			}

			outActions = append(outActions, outAction)
		}

		outActionFilters = append(outActionFilters, apiclient.OrganizationWorkflowActionFilter{
			LogicType:  apiclient.OrganizationWorkflowActionFilterLogicType(inActionFilter.LogicType.Get()),
			Conditions: []interface{}{},
			Actions:    outActions,
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
		ActionFilters: outActionFilters,
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

	var outActionFilters []AlertResourceModelActionFiltersItem
	for _, actionFilter := range data.ActionFilters {
		var outActions []AlertResourceModelActionFiltersItemActionsItem
		for _, action := range actionFilter.Actions {
			outAction := AlertResourceModelActionFiltersItemActionsItem{
				Email:      supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemEmail](ctx),
				Plugin:     supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemPlugin](ctx),
				Slack:      supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemSlack](ctx),
				Pagerduty:  supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemPagerduty](ctx),
				Discord:    supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemDiscord](ctx),
				Msteams:    supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemMsteams](ctx),
				Opsgenie:   supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemOpsgenie](ctx),
				Vsts:       supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemVsts](ctx),
				Jira:       supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemJira](ctx),
				JiraServer: supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemJiraServer](ctx),
				Github:     supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemGithub](ctx),
			}

			actionValue, err := action.ValueByDiscriminator()
			if err != nil {
				diags.AddError("Failed to get action value", err.Error())
				return
			}

			switch actionValue := actionValue.(type) {
			case apiclient.OrganizationWorkflowActionFilterActionEmail:
				var outEmail AlertResourceModelActionFiltersItemActionsItemEmail
				outEmail.TargetType = supertypes.NewStringValue(string(actionValue.Config.TargetType))
				if actionValue.Config.TargetIdentifier != nil {
					outEmail.TargetId = supertypes.NewStringPointerValue(actionValue.Config.TargetIdentifier)
				}
				if actionValue.Data.FallthroughType != nil {
					outEmail.FallthroughType = supertypes.NewStringValue(string(*actionValue.Data.FallthroughType))
				}

				outAction.Email = supertypes.NewSingleNestedObjectValueOf(ctx, &outEmail)

			case apiclient.OrganizationWorkflowActionFilterActionPlugin:
				var outPlugin AlertResourceModelActionFiltersItemActionsItemPlugin

				outAction.Plugin = supertypes.NewSingleNestedObjectValueOf(ctx, &outPlugin)

			case apiclient.OrganizationWorkflowActionFilterActionSlack:
				var outSlack AlertResourceModelActionFiltersItemActionsItemSlack
				outSlack.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outSlack.ChannelId = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)
				outSlack.ChannelName = supertypes.NewStringValue(actionValue.Config.TargetDisplay)
				if actionValue.Data.Tags != nil {
					outSlack.Tags = supertypes.NewStringValue(*actionValue.Data.Tags)
				}
				if actionValue.Data.Notes != nil {
					outSlack.Notes = supertypes.NewStringValue(*actionValue.Data.Notes)
				}

				outAction.Slack = supertypes.NewSingleNestedObjectValueOf(ctx, &outSlack)

			case apiclient.OrganizationWorkflowActionFilterActionPagerDuty:
				var outPagerduty AlertResourceModelActionFiltersItemActionsItemPagerduty
				outPagerduty.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outPagerduty.ServiceName = supertypes.NewStringValue(actionValue.Config.TargetDisplay)
				outPagerduty.ServiceId = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)
				outPagerduty.Severity = supertypes.NewStringValue(string(*actionValue.Data.Priority))

				outAction.Pagerduty = supertypes.NewSingleNestedObjectValueOf(ctx, &outPagerduty)

			case apiclient.OrganizationWorkflowActionFilterActionDiscord:
				var outDiscord AlertResourceModelActionFiltersItemActionsItemDiscord
				outDiscord.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outDiscord.ChannelId = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)
				if actionValue.Data.Tags != nil {
					outDiscord.Tags = supertypes.NewStringValue(*actionValue.Data.Tags)
				}

				outAction.Discord = supertypes.NewSingleNestedObjectValueOf(ctx, &outDiscord)

			case apiclient.OrganizationWorkflowActionFilterActionMsTeams:
				var outMsteams AlertResourceModelActionFiltersItemActionsItemMsteams
				outMsteams.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outMsteams.TeamId = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)
				outMsteams.ChannelName = supertypes.NewStringValue(actionValue.Config.TargetDisplay)

				outAction.Msteams = supertypes.NewSingleNestedObjectValueOf(ctx, &outMsteams)

			case apiclient.OrganizationWorkflowActionFilterActionOpsgenie:
				var outOpsgenie AlertResourceModelActionFiltersItemActionsItemOpsgenie
				outOpsgenie.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outOpsgenie.TeamId = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)
				outOpsgenie.TeamName = supertypes.NewStringValue(actionValue.Config.TargetDisplay)
				outOpsgenie.Priority = supertypes.NewStringValue(string(*actionValue.Data.Priority))

				outAction.Opsgenie = supertypes.NewSingleNestedObjectValueOf(ctx, &outOpsgenie)

			case apiclient.OrganizationWorkflowActionFilterActionVsts:
				var outVsts AlertResourceModelActionFiltersItemActionsItemVsts
				outVsts.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outVsts.Data.DiagsSet(ctx, diags, actionValue.Data)
				if diags.HasError() {
					return
				}

				outAction.Vsts = supertypes.NewSingleNestedObjectValueOf(ctx, &outVsts)

			case apiclient.OrganizationWorkflowActionFilterActionJira:
				var outJira AlertResourceModelActionFiltersItemActionsItemJira
				outJira.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outJira.Data.DiagsSet(ctx, diags, actionValue.Data)
				if diags.HasError() {
					return
				}

				outAction.Jira = supertypes.NewSingleNestedObjectValueOf(ctx, &outJira)

			case apiclient.OrganizationWorkflowActionFilterActionJiraServer:
				var outJiraServer AlertResourceModelActionFiltersItemActionsItemJiraServer
				outJiraServer.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outJiraServer.Data.DiagsSet(ctx, diags, actionValue.Data)
				if diags.HasError() {
					return
				}

				outAction.JiraServer = supertypes.NewSingleNestedObjectValueOf(ctx, &outJiraServer)

			case apiclient.OrganizationWorkflowActionFilterActionGitHub:
				var outGithub AlertResourceModelActionFiltersItemActionsItemGithub
				outGithub.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outGithub.Repo = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Repo)
				if actionValue.Data.AdditionalFields.Assignee != "" {
					outGithub.Assignee = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Assignee)
				}
				outGithub.Labels = supertypes.NewSetValueOfSlice(ctx, actionValue.Data.AdditionalFields.Labels)
				if diags.HasError() {
					return
				}

				outAction.Github = supertypes.NewSingleNestedObjectValueOf(ctx, &outGithub)
			}

			outActions = append(outActions, outAction)
		}

		outActionFilters = append(outActionFilters, AlertResourceModelActionFiltersItem{
			LogicType:  supertypes.NewStringValue(string(actionFilter.LogicType)),
			Conditions: supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelActionFiltersItemConditionsItem{}),
			Actions:    supertypes.NewListNestedObjectValueOfValueSlice(ctx, outActions),
		})
	}
	m.ActionFilters = supertypes.NewListNestedObjectValueOfValueSlice(ctx, outActionFilters)

	return
}
