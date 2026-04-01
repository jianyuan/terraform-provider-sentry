package provider

import (
	"context"
	"encoding/json"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (r *AlertResource) getActionFilters(ctx context.Context, data AlertResourceModel) ([]apiclient.OrganizationWorkflowActionFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	inActionFilters := data.ActionFilters.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}
	var outActionFilters []apiclient.OrganizationWorkflowActionFilter
	for _, inActionFilter := range inActionFilters {
		// Conditions
		inConditions := inActionFilter.Conditions.DiagsGet(ctx, diags)
		if diags.HasError() {
			return nil, diags
		}

		outConditions := []apiclient.OrganizationWorkflowActionFilterCondition{}
		for _, inCondition := range inConditions {
			var outCondition apiclient.OrganizationWorkflowActionFilterCondition
			switch {
			case inCondition.AgeComparison.IsKnown():
				inAgeComparison := inCondition.AgeComparison.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outAgeComparison apiclient.OrganizationWorkflowActionFilterConditionAgeComparison
				outAgeComparison.Comparison.ComparisonType = apiclient.OrganizationWorkflowActionFilterConditionAgeComparisonComparisonComparisonType(inAgeComparison.ComparisonType.Get())
				outAgeComparison.Comparison.Time = apiclient.OrganizationWorkflowActionFilterConditionAgeComparisonComparisonTime(inAgeComparison.Time.Get())
				outAgeComparison.Comparison.Value = inAgeComparison.Value.Get()
				outAgeComparison.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionAgeComparison(outAgeComparison); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.AssignedTo.IsKnown():
				inAssignedTo := inCondition.AssignedTo.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outAssignedTo apiclient.OrganizationWorkflowActionFilterConditionAssignedTo
				outAssignedTo.Comparison.TargetType = apiclient.OrganizationWorkflowActionFilterConditionAssignedToComparisonTargetType(inAssignedTo.TargetType.Get())
				if inAssignedTo.TargetId.IsKnown() {
					if err := outAssignedTo.Comparison.TargetIdentifier.FromOrganizationWorkflowActionFilterConditionAssignedToComparisonTargetIdentifier0(inAssignedTo.TargetId.Get()); err != nil {
						diags.AddError("Failed to create condition", err.Error())
						return nil, diags
					}
				} else {
					if err := outAssignedTo.Comparison.TargetIdentifier.FromOrganizationWorkflowActionFilterConditionAssignedToComparisonTargetIdentifier0(""); err != nil {
						diags.AddError("Failed to create condition", err.Error())
						return nil, diags
					}
				}
				outAssignedTo.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionAssignedTo(outAssignedTo); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.IssueCategory.IsKnown():
				inIssueCategory := inCondition.IssueCategory.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outIssueCategory apiclient.OrganizationWorkflowActionFilterConditionIssueCategory
				outIssueCategory.Comparison.Value = inIssueCategory.Value.Get()
				outIssueCategory.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionIssueCategory(outIssueCategory); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.IssueOccurrences.IsKnown():
				inIssueOccurrences := inCondition.IssueOccurrences.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outIssueOccurrences apiclient.OrganizationWorkflowActionFilterConditionIssueOccurrences
				outIssueOccurrences.Comparison.Value = inIssueOccurrences.Value.Get()
				outIssueOccurrences.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionIssueOccurrences(outIssueOccurrences); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.IssuePriorityDeescalating.IsKnown():
				var outIssuePriorityDeescalating apiclient.OrganizationWorkflowActionFilterConditionIssuePriorityDeescalating
				outIssuePriorityDeescalating.Comparison = true
				outIssuePriorityDeescalating.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionIssuePriorityDeescalating(outIssuePriorityDeescalating); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.IssuePriorityGreaterOrEqual.IsKnown():
				inIssuePriorityGreaterOrEqual := inCondition.IssuePriorityGreaterOrEqual.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outIssuePriorityGreaterOrEqual apiclient.OrganizationWorkflowActionFilterConditionIssuePriorityGreaterOrEqual
				outIssuePriorityGreaterOrEqual.Comparison = inIssuePriorityGreaterOrEqual.Comparison.Get()
				outIssuePriorityGreaterOrEqual.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionIssuePriorityGreaterOrEqual(outIssuePriorityGreaterOrEqual); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.EventUniqueUserFrequencyCount.IsKnown():
				inEventUniqueUserFrequencyCount := inCondition.EventUniqueUserFrequencyCount.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				inFilters := inEventUniqueUserFrequencyCount.Filters.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEventUniqueUserFrequencyCount apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCount
				outEventUniqueUserFrequencyCount.Comparison.Value = inEventUniqueUserFrequencyCount.Value.Get()
				outEventUniqueUserFrequencyCount.Comparison.Filters = lo.Map(inFilters, func(inFilter *AlertResourceModelActionFiltersItemConditionsItemEventUniqueUserFrequencyCountFiltersItem, _ int) apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter {
					return apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter{
						Attribute: inFilter.Attribute.GetPtr(),
						Key:       inFilter.Key.GetPtr(),
						Match:     inFilter.Match.GetPtr(),
						Value:     inFilter.Value.GetPtr(),
					}
				})
				outEventUniqueUserFrequencyCount.Comparison.Interval = inEventUniqueUserFrequencyCount.Interval.Get()
				outEventUniqueUserFrequencyCount.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCount(outEventUniqueUserFrequencyCount); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.EventFrequencyCount.IsKnown():
				inEventFrequencyCount := inCondition.EventFrequencyCount.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEventFrequencyCount apiclient.OrganizationWorkflowActionFilterConditionEventFrequencyCount
				outEventFrequencyCount.Comparison.Value = inEventFrequencyCount.Value.Get()
				outEventFrequencyCount.Comparison.Interval = inEventFrequencyCount.Interval.Get()
				outEventFrequencyCount.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionEventFrequencyCount(outEventFrequencyCount); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.EventFrequencyPercent.IsKnown():
				inEventFrequencyPercent := inCondition.EventFrequencyPercent.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEventFrequencyPercent apiclient.OrganizationWorkflowActionFilterConditionEventFrequencyPercent
				outEventFrequencyPercent.Comparison.Value = inEventFrequencyPercent.Value.Get()
				outEventFrequencyPercent.Comparison.Interval = inEventFrequencyPercent.Interval.Get()
				outEventFrequencyPercent.Comparison.ComparisonInterval = inEventFrequencyPercent.ComparisonInterval.Get()
				outEventFrequencyPercent.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionEventFrequencyPercent(outEventFrequencyPercent); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.PercentSessionsCount.IsKnown():
				inPercentSessionsCount := inCondition.PercentSessionsCount.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outPercentSessionsCount apiclient.OrganizationWorkflowActionFilterConditionPercentSessionsCount
				outPercentSessionsCount.Comparison.Value = inPercentSessionsCount.Value.Get()
				outPercentSessionsCount.Comparison.Interval = inPercentSessionsCount.Interval.Get()
				outPercentSessionsCount.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionPercentSessionsCount(outPercentSessionsCount); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.PercentSessionsPercent.IsKnown():
				inPercentSessionsPercent := inCondition.PercentSessionsPercent.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outPercentSessionsPercent apiclient.OrganizationWorkflowActionFilterConditionPercentSessionsPercent
				outPercentSessionsPercent.Comparison.Value = inPercentSessionsPercent.Value.Get()
				outPercentSessionsPercent.Comparison.Interval = inPercentSessionsPercent.Interval.Get()
				outPercentSessionsPercent.Comparison.ComparisonInterval = inPercentSessionsPercent.ComparisonInterval.Get()
				outPercentSessionsPercent.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionPercentSessionsPercent(outPercentSessionsPercent); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.EventAttribute.IsKnown():
				inEventAttribute := inCondition.EventAttribute.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEventAttribute apiclient.OrganizationWorkflowActionFilterConditionEventAttribute
				outEventAttribute.Comparison.Attribute = inEventAttribute.Attribute.Get()
				outEventAttribute.Comparison.Match = inEventAttribute.Match.Get()
				outEventAttribute.Comparison.Value = inEventAttribute.Value.Get()
				outEventAttribute.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionEventAttribute(outEventAttribute); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.TaggedEvent.IsKnown():
				inTaggedEvent := inCondition.TaggedEvent.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outTaggedEvent apiclient.OrganizationWorkflowActionFilterConditionTaggedEvent
				outTaggedEvent.Comparison.Key = inTaggedEvent.Key.Get()
				outTaggedEvent.Comparison.Match = inTaggedEvent.Match.Get()
				outTaggedEvent.Comparison.Value = inTaggedEvent.Value.GetPtr()
				outTaggedEvent.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionTaggedEvent(outTaggedEvent); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.LatestRelease.IsKnown():
				var outLatestRelease apiclient.OrganizationWorkflowActionFilterConditionLatestRelease
				outLatestRelease.Comparison = true
				outLatestRelease.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionLatestRelease(outLatestRelease); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.LatestAdoptedRelease.IsKnown():
				inLatestAdoptedRelease := inCondition.LatestAdoptedRelease.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outLatestAdoptedRelease apiclient.OrganizationWorkflowActionFilterConditionLatestAdoptedRelease
				outLatestAdoptedRelease.Comparison.Environment = inLatestAdoptedRelease.Environment.Get()
				outLatestAdoptedRelease.Comparison.AgeComparison = inLatestAdoptedRelease.AgeComparison.Get()
				outLatestAdoptedRelease.Comparison.ReleaseAgeType = inLatestAdoptedRelease.ReleaseAgeType.Get()
				outLatestAdoptedRelease.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionLatestAdoptedRelease(outLatestAdoptedRelease); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.Level.IsKnown():
				inLevel := inCondition.Level.DiagsGet(ctx, diags)
				if diags.HasError() {
					return nil, diags
				}

				var outLevel apiclient.OrganizationWorkflowActionFilterConditionLevel
				outLevel.Comparison.Match = inLevel.Match.Get()
				outLevel.Comparison.Level = inLevel.Level.Get()
				outLevel.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionLevel(outLevel); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}
			}
			outConditions = append(outConditions, outCondition)
		}

		// Actions
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
				outVsts.Data.AdditionalFields.Project = inVsts.Project.Get()
				outVsts.Data.AdditionalFields.WorkItemType = inVsts.WorkItemType.Get()
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
				outJira.Data.AdditionalFields.Project = inJira.Project.Get()
				outJira.Data.AdditionalFields.Issuetype = inJira.IssueType.Get()
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
				outJiraServer.Data.AdditionalFields.Project = inJiraServer.Project.Get()
				outJiraServer.Data.AdditionalFields.Issuetype = inJiraServer.IssueType.Get()
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
			Conditions: outConditions,
			Actions:    outActions,
		})
	}

	return outActionFilters, diags
}

func (r *AlertResource) getTriggerConditions(ctx context.Context, data AlertResourceModel) ([]apiclient.OrganizationWorkflowTriggerCondition, diag.Diagnostics) {
	var diags diag.Diagnostics

	inTriggerConditions := data.TriggerConditions.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}
	var outTriggerConditions []apiclient.OrganizationWorkflowTriggerCondition
	for _, triggerCondition := range inTriggerConditions {
		var outTriggerConditionComparison apiclient.OrganizationWorkflowTriggerCondition_Comparison
		if err := outTriggerConditionComparison.FromOrganizationWorkflowTriggerConditionComparison0(true); err != nil {
			diags.AddError("Failed to create trigger condition", err.Error())
			return nil, diags
		}

		var outTriggerCondition apiclient.OrganizationWorkflowTriggerCondition
		outTriggerCondition.Comparison = outTriggerConditionComparison
		outTriggerCondition.ConditionResult = true

		switch {
		case triggerCondition.FirstSeenEvent.IsKnown():
			outTriggerCondition.Type = "first_seen_event"
		case triggerCondition.IssueResolvedTrigger.IsKnown():
			outTriggerCondition.Type = "issue_resolved_trigger"
		case triggerCondition.ReappearedEvent.IsKnown():
			outTriggerCondition.Type = "reappeared_event"
		case triggerCondition.RegressionEvent.IsKnown():
			outTriggerCondition.Type = "regression_event"
		}

		outTriggerConditions = append(outTriggerConditions, outTriggerCondition)
	}

	return outTriggerConditions, diags
}

func (r *AlertResource) getCreateJSONRequestBody(ctx context.Context, data AlertResourceModel) (*apiclient.CreateOrganizationWorkflowJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	monitorIds := data.MonitorIds.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}

	req := apiclient.CreateOrganizationWorkflowJSONRequestBody{
		Name:        data.Name.Get(),
		Enabled:     data.Enabled.Get(),
		Environment: data.Environment.Get(),
		Config: apiclient.OrganizationWorkflowConfig{
			Frequency: data.FrequencyMinutes.Get(),
		},
		DetectorIds: monitorIds,
		Triggers: apiclient.OrganizationWorkflowTrigger{
			LogicType:  apiclient.OrganizationWorkflowTriggerLogicTypeAnyShort,
			Conditions: tfutils.MergeDiagnostics(r.getTriggerConditions(ctx, data))(&diags),
		},
		ActionFilters: tfutils.MergeDiagnostics(r.getActionFilters(ctx, data))(&diags),
	}

	return &req, nil
}

func (r *AlertResource) getUpdateJSONRequestBody(ctx context.Context, data AlertResourceModel) (*apiclient.UpdateOrganizationWorkflowJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	monitorIds := data.MonitorIds.DiagsGet(ctx, diags)
	if diags.HasError() {
		return nil, diags
	}

	req := apiclient.UpdateOrganizationWorkflowJSONRequestBody{
		Id:          data.Id.Get(),
		Name:        data.Name.Get(),
		Enabled:     data.Enabled.Get(),
		Environment: data.Environment.Get(),
		Config: apiclient.OrganizationWorkflowConfig{
			Frequency: data.FrequencyMinutes.Get(),
		},
		DetectorIds: monitorIds,
		Triggers: apiclient.OrganizationWorkflowTrigger{
			LogicType:  apiclient.OrganizationWorkflowTriggerLogicTypeAnyShort,
			Conditions: tfutils.MergeDiagnostics(r.getTriggerConditions(ctx, data))(&diags),
		},
		ActionFilters: tfutils.MergeDiagnostics(r.getActionFilters(ctx, data))(&diags),
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

	var triggerConditions []AlertResourceModelTriggerConditionsItem
	for _, triggerCondition := range data.Triggers.Conditions {
		outTriggerCondition := AlertResourceModelTriggerConditionsItem{
			FirstSeenEvent:       supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemFirstSeenEvent](ctx),
			IssueResolvedTrigger: supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemIssueResolvedTrigger](ctx),
			ReappearedEvent:      supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemReappearedEvent](ctx),
			RegressionEvent:      supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemRegressionEvent](ctx),
		}
		switch triggerCondition.Type {
		case "first_seen_event":
			outTriggerCondition.FirstSeenEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemFirstSeenEvent{})
		case "issue_resolved_trigger":
			outTriggerCondition.IssueResolvedTrigger = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemIssueResolvedTrigger{})
		case "reappeared_event":
			outTriggerCondition.ReappearedEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemReappearedEvent{})
		case "regression_event":
			outTriggerCondition.RegressionEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemRegressionEvent{})
		}
		triggerConditions = append(triggerConditions, outTriggerCondition)
	}
	m.TriggerConditions = supertypes.NewListNestedObjectValueOfValueSlice(ctx, triggerConditions)

	var outActionFilters []AlertResourceModelActionFiltersItem
	for _, actionFilter := range data.ActionFilters {
		// Conditions

		// NOTE: The API returns conditions in a random order, so we need to sort them by ID to ensure that the
		// order is deterministic.
		slices.SortFunc(actionFilter.Conditions, func(a, b apiclient.OrganizationWorkflowActionFilterCondition) int {
			var aData struct {
				Id string `json:"id"`
			}
			var bData struct {
				Id string `json:"id"`
			}
			must.Do(json.Unmarshal(must.Get(a.MarshalJSON()), &aData))
			must.Do(json.Unmarshal(must.Get(b.MarshalJSON()), &bData))
			aId := must.Get(strconv.Atoi(aData.Id))
			bId := must.Get(strconv.Atoi(bData.Id))
			return aId - bId
		})

		outConditions := []AlertResourceModelActionFiltersItemConditionsItem{}
		for _, condition := range actionFilter.Conditions {
			outCondition := AlertResourceModelActionFiltersItemConditionsItem{
				AgeComparison:                 supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemAgeComparison](ctx),
				AssignedTo:                    supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemAssignedTo](ctx),
				IssueCategory:                 supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemIssueCategory](ctx),
				IssueOccurrences:              supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemIssueOccurrences](ctx),
				IssuePriorityDeescalating:     supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemIssuePriorityDeescalating](ctx),
				IssuePriorityGreaterOrEqual:   supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemIssuePriorityGreaterOrEqual](ctx),
				EventUniqueUserFrequencyCount: supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemEventUniqueUserFrequencyCount](ctx),
				EventFrequencyCount:           supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemEventFrequencyCount](ctx),
				EventFrequencyPercent:         supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemEventFrequencyPercent](ctx),
				PercentSessionsCount:          supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemPercentSessionsCount](ctx),
				PercentSessionsPercent:        supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemPercentSessionsPercent](ctx),
				EventAttribute:                supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemEventAttribute](ctx),
				TaggedEvent:                   supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemTaggedEvent](ctx),
				LatestRelease:                 supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemLatestRelease](ctx),
				LatestAdoptedRelease:          supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemLatestAdoptedRelease](ctx),
				Level:                         supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemLevel](ctx),
			}

			conditionValue, err := condition.ValueByDiscriminator()
			if err != nil {
				diags.AddError("Failed to get condition value", err.Error())
				return
			}

			switch conditionValue := conditionValue.(type) {
			case apiclient.OrganizationWorkflowActionFilterConditionAgeComparison:
				var outAgeComparison AlertResourceModelActionFiltersItemConditionsItemAgeComparison
				outAgeComparison.Time = supertypes.NewStringValue(string(conditionValue.Comparison.Time))
				outAgeComparison.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)
				outAgeComparison.ComparisonType = supertypes.NewStringValue(string(conditionValue.Comparison.ComparisonType))

				outCondition.AgeComparison = supertypes.NewSingleNestedObjectValueOf(ctx, &outAgeComparison)

			case apiclient.OrganizationWorkflowActionFilterConditionAssignedTo:
				var assignedTo AlertResourceModelActionFiltersItemConditionsItemAssignedTo
				assignedTo.TargetType = supertypes.NewStringValue(string(conditionValue.Comparison.TargetType))
				if v, err := conditionValue.Comparison.TargetIdentifier.AsOrganizationWorkflowActionFilterConditionAssignedToComparisonTargetIdentifier0(); err == nil {
					if v == "" {
						assignedTo.TargetId = supertypes.NewStringNull()
					} else {
						assignedTo.TargetId = supertypes.NewStringValue(v)
					}
				} else if v, err := conditionValue.Comparison.TargetIdentifier.AsOrganizationWorkflowActionFilterConditionAssignedToComparisonTargetIdentifier1(); err == nil {
					assignedTo.TargetId = supertypes.NewStringValue(strconv.FormatInt(v, 10))
				}

				outCondition.AssignedTo = supertypes.NewSingleNestedObjectValueOf(ctx, &assignedTo)

			case apiclient.OrganizationWorkflowActionFilterConditionIssueCategory:
				var issueCategory AlertResourceModelActionFiltersItemConditionsItemIssueCategory
				issueCategory.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)

				outCondition.IssueCategory = supertypes.NewSingleNestedObjectValueOf(ctx, &issueCategory)

			case apiclient.OrganizationWorkflowActionFilterConditionIssueOccurrences:
				var issueOccurrences AlertResourceModelActionFiltersItemConditionsItemIssueOccurrences
				issueOccurrences.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)

				outCondition.IssueOccurrences = supertypes.NewSingleNestedObjectValueOf(ctx, &issueOccurrences)

			case apiclient.OrganizationWorkflowActionFilterConditionIssuePriorityDeescalating:
				var issuePriorityDeescalating AlertResourceModelActionFiltersItemConditionsItemIssuePriorityDeescalating

				outCondition.IssuePriorityDeescalating = supertypes.NewSingleNestedObjectValueOf(ctx, &issuePriorityDeescalating)

			case apiclient.OrganizationWorkflowActionFilterConditionIssuePriorityGreaterOrEqual:
				var issuePriorityGreaterOrEqual AlertResourceModelActionFiltersItemConditionsItemIssuePriorityGreaterOrEqual
				issuePriorityGreaterOrEqual.Comparison = supertypes.NewInt64Value(conditionValue.Comparison)

				outCondition.IssuePriorityGreaterOrEqual = supertypes.NewSingleNestedObjectValueOf(ctx, &issuePriorityGreaterOrEqual)

			case apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCount:
				var eventUniqueUserFrequencyCount AlertResourceModelActionFiltersItemConditionsItemEventUniqueUserFrequencyCount
				eventUniqueUserFrequencyCount.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)
				eventUniqueUserFrequencyCount.Interval = supertypes.NewStringValue(conditionValue.Comparison.Interval)

				outFilters := []AlertResourceModelActionFiltersItemConditionsItemEventUniqueUserFrequencyCountFiltersItem{}
				for _, filter := range conditionValue.Comparison.Filters {
					outFilters = append(outFilters, AlertResourceModelActionFiltersItemConditionsItemEventUniqueUserFrequencyCountFiltersItem{
						Attribute: supertypes.NewStringPointerValueOrNull(filter.Attribute),
						Key:       supertypes.NewStringPointerValueOrNull(filter.Key),
						Match:     supertypes.NewStringPointerValueOrNull(filter.Match),
						Value:     supertypes.NewStringPointerValueOrNull(filter.Value),
					})
				}
				eventUniqueUserFrequencyCount.Filters = supertypes.NewListNestedObjectValueOfValueSlice(ctx, outFilters)

				outCondition.EventUniqueUserFrequencyCount = supertypes.NewSingleNestedObjectValueOf(ctx, &eventUniqueUserFrequencyCount)

			case apiclient.OrganizationWorkflowActionFilterConditionEventFrequencyCount:
				var eventFrequencyCount AlertResourceModelActionFiltersItemConditionsItemEventFrequencyCount
				eventFrequencyCount.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)
				eventFrequencyCount.Interval = supertypes.NewStringValue(conditionValue.Comparison.Interval)

				outCondition.EventFrequencyCount = supertypes.NewSingleNestedObjectValueOf(ctx, &eventFrequencyCount)

			case apiclient.OrganizationWorkflowActionFilterConditionEventFrequencyPercent:
				var eventFrequencyPercent AlertResourceModelActionFiltersItemConditionsItemEventFrequencyPercent
				eventFrequencyPercent.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)
				eventFrequencyPercent.Interval = supertypes.NewStringValue(conditionValue.Comparison.Interval)
				eventFrequencyPercent.ComparisonInterval = supertypes.NewStringValue(conditionValue.Comparison.ComparisonInterval)

				outCondition.EventFrequencyPercent = supertypes.NewSingleNestedObjectValueOf(ctx, &eventFrequencyPercent)

			case apiclient.OrganizationWorkflowActionFilterConditionPercentSessionsCount:
				var percentSessionsCount AlertResourceModelActionFiltersItemConditionsItemPercentSessionsCount
				percentSessionsCount.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)
				percentSessionsCount.Interval = supertypes.NewStringValue(conditionValue.Comparison.Interval)

				outCondition.PercentSessionsCount = supertypes.NewSingleNestedObjectValueOf(ctx, &percentSessionsCount)

			case apiclient.OrganizationWorkflowActionFilterConditionPercentSessionsPercent:
				var percentSessionsPercent AlertResourceModelActionFiltersItemConditionsItemPercentSessionsPercent
				percentSessionsPercent.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)
				percentSessionsPercent.Interval = supertypes.NewStringValue(conditionValue.Comparison.Interval)
				percentSessionsPercent.ComparisonInterval = supertypes.NewStringValue(conditionValue.Comparison.ComparisonInterval)

				outCondition.PercentSessionsPercent = supertypes.NewSingleNestedObjectValueOf(ctx, &percentSessionsPercent)

			case apiclient.OrganizationWorkflowActionFilterConditionEventAttribute:
				var eventAttribute AlertResourceModelActionFiltersItemConditionsItemEventAttribute
				eventAttribute.Attribute = supertypes.NewStringValue(conditionValue.Comparison.Attribute)
				eventAttribute.Match = supertypes.NewStringValue(conditionValue.Comparison.Match)
				eventAttribute.Value = supertypes.NewStringValue(conditionValue.Comparison.Value)

				outCondition.EventAttribute = supertypes.NewSingleNestedObjectValueOf(ctx, &eventAttribute)

			case apiclient.OrganizationWorkflowActionFilterConditionTaggedEvent:
				var taggedEvent AlertResourceModelActionFiltersItemConditionsItemTaggedEvent
				taggedEvent.Key = supertypes.NewStringValue(conditionValue.Comparison.Key)
				taggedEvent.Match = supertypes.NewStringValue(conditionValue.Comparison.Match)
				taggedEvent.Value = supertypes.NewStringPointerValueOrNull(conditionValue.Comparison.Value)

				outCondition.TaggedEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &taggedEvent)

			case apiclient.OrganizationWorkflowActionFilterConditionLatestRelease:
				var latestRelease AlertResourceModelActionFiltersItemConditionsItemLatestRelease

				outCondition.LatestRelease = supertypes.NewSingleNestedObjectValueOf(ctx, &latestRelease)

			case apiclient.OrganizationWorkflowActionFilterConditionLatestAdoptedRelease:
				var latestAdoptedRelease AlertResourceModelActionFiltersItemConditionsItemLatestAdoptedRelease
				latestAdoptedRelease.Environment = supertypes.NewStringValue(conditionValue.Comparison.Environment)
				latestAdoptedRelease.AgeComparison = supertypes.NewStringValue(conditionValue.Comparison.AgeComparison)
				latestAdoptedRelease.ReleaseAgeType = supertypes.NewStringValue(conditionValue.Comparison.ReleaseAgeType)

				outCondition.LatestAdoptedRelease = supertypes.NewSingleNestedObjectValueOf(ctx, &latestAdoptedRelease)

			case apiclient.OrganizationWorkflowActionFilterConditionLevel:
				var level AlertResourceModelActionFiltersItemConditionsItemLevel
				level.Match = supertypes.NewStringValue(conditionValue.Comparison.Match)
				level.Level = supertypes.NewInt64Value(conditionValue.Comparison.Level)

				outCondition.Level = supertypes.NewSingleNestedObjectValueOf(ctx, &level)
			}

			outConditions = append(outConditions, outCondition)
		}

		// Actions
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
				outVsts.Project = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Project)
				outVsts.WorkItemType = supertypes.NewStringValue(actionValue.Data.AdditionalFields.WorkItemType)

				outAction.Vsts = supertypes.NewSingleNestedObjectValueOf(ctx, &outVsts)

			case apiclient.OrganizationWorkflowActionFilterActionJira:
				var outJira AlertResourceModelActionFiltersItemActionsItemJira
				outJira.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outJira.Project = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Project)
				outJira.IssueType = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Issuetype)

				outAction.Jira = supertypes.NewSingleNestedObjectValueOf(ctx, &outJira)

			case apiclient.OrganizationWorkflowActionFilterActionJiraServer:
				var outJiraServer AlertResourceModelActionFiltersItemActionsItemJiraServer
				outJiraServer.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outJiraServer.Project = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Project)
				outJiraServer.IssueType = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Issuetype)

				outAction.JiraServer = supertypes.NewSingleNestedObjectValueOf(ctx, &outJiraServer)

			case apiclient.OrganizationWorkflowActionFilterActionGitHub:
				var outGithub AlertResourceModelActionFiltersItemActionsItemGithub
				outGithub.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outGithub.Repo = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Repo)
				if actionValue.Data.AdditionalFields.Assignee != "" {
					outGithub.Assignee = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Assignee)
				}
				outGithub.Labels = supertypes.NewSetValueOfSlice(ctx, actionValue.Data.AdditionalFields.Labels)

				outAction.Github = supertypes.NewSingleNestedObjectValueOf(ctx, &outGithub)
			}

			if diags.HasError() {
				return
			}

			outActions = append(outActions, outAction)
		}

		outActionFilters = append(outActionFilters, AlertResourceModelActionFiltersItem{
			LogicType:  supertypes.NewStringValue(string(actionFilter.LogicType)),
			Conditions: supertypes.NewListNestedObjectValueOfValueSlice(ctx, outConditions),
			Actions:    supertypes.NewListNestedObjectValueOfValueSlice(ctx, outActions),
		})
	}
	m.ActionFilters = supertypes.NewListNestedObjectValueOfValueSlice(ctx, outActionFilters)

	return
}
