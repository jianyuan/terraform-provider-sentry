package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-plugin-framework-utils/fwdiag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/must"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (r *AlertResource) getActionFilters(ctx context.Context, data AlertResourceModel) ([]apiclient.OrganizationWorkflowActionFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	inActionFilters := fwdiag.Merge(data.ActionFilters.Get(ctx))(&diags)
	if diags.HasError() {
		return nil, diags
	}
	var outActionFilters []apiclient.OrganizationWorkflowActionFilter
	for _, inActionFilter := range inActionFilters {
		// Conditions
		inConditions := fwdiag.Merge(inActionFilter.Conditions.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}

		outConditions := []apiclient.OrganizationWorkflowActionFilterCondition{}
		for _, inCondition := range inConditions {
			var outCondition apiclient.OrganizationWorkflowActionFilterCondition
			switch {
			case inCondition.AgeComparison.IsKnown():
				inAgeComparison := fwdiag.Merge(inCondition.AgeComparison.Get(ctx))(&diags)
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
				inAssignedTo := fwdiag.Merge(inCondition.AssignedTo.Get(ctx))(&diags)
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
				inIssueCategory := fwdiag.Merge(inCondition.IssueCategory.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outIssueCategory apiclient.OrganizationWorkflowActionFilterConditionIssueCategory
				outIssueCategory.Comparison.Value = inIssueCategory.Value.Get()
				outIssueCategory.Comparison.Include = new(inIssueCategory.Include.Get())
				outIssueCategory.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionIssueCategory(outIssueCategory); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.IssueOccurrences.IsKnown():
				inIssueOccurrences := fwdiag.Merge(inCondition.IssueOccurrences.Get(ctx))(&diags)
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
				inIssuePriorityDeescalating := fwdiag.Merge(inCondition.IssuePriorityDeescalating.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outIssuePriorityDeescalating apiclient.OrganizationWorkflowActionFilterConditionIssuePriorityDeescalating
				outIssuePriorityDeescalating.Comparison = inIssuePriorityDeescalating.Comparison.Get()
				outIssuePriorityDeescalating.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionIssuePriorityDeescalating(outIssuePriorityDeescalating); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.IssuePriorityGreaterOrEqual.IsKnown():
				inIssuePriorityGreaterOrEqual := fwdiag.Merge(inCondition.IssuePriorityGreaterOrEqual.Get(ctx))(&diags)
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
				inEventUniqueUserFrequencyCount := fwdiag.Merge(inCondition.EventUniqueUserFrequencyCount.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				inFilters := fwdiag.Merge(inEventUniqueUserFrequencyCount.Filters.Get(ctx))(&diags)
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
				inEventFrequencyCount := fwdiag.Merge(inCondition.EventFrequencyCount.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				inFrequencyCountFilters := fwdiag.Merge(inEventFrequencyCount.Filters.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEventFrequencyCount apiclient.OrganizationWorkflowActionFilterConditionEventFrequencyCount
				outEventFrequencyCount.Comparison.Value = inEventFrequencyCount.Value.Get()
				outEventFrequencyCount.Comparison.Filters = lo.ToPtr(lo.Map(inFrequencyCountFilters, func(inFilter *AlertResourceModelActionFiltersItemConditionsItemEventFrequencyCountFiltersItem, _ int) apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter {
					return apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter{
						Attribute: inFilter.Attribute.GetPtr(),
						Key:       inFilter.Key.GetPtr(),
						Match:     inFilter.Match.GetPtr(),
						Value:     inFilter.Value.GetPtr(),
					}
				}))
				outEventFrequencyCount.Comparison.Interval = inEventFrequencyCount.Interval.Get()
				outEventFrequencyCount.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionEventFrequencyCount(outEventFrequencyCount); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.EventFrequencyPercent.IsKnown():
				inEventFrequencyPercent := fwdiag.Merge(inCondition.EventFrequencyPercent.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				inFrequencyPercentFilters := fwdiag.Merge(inEventFrequencyPercent.Filters.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEventFrequencyPercent apiclient.OrganizationWorkflowActionFilterConditionEventFrequencyPercent
				outEventFrequencyPercent.Comparison.Value = inEventFrequencyPercent.Value.Get()
				outEventFrequencyPercent.Comparison.Filters = lo.ToPtr(lo.Map(inFrequencyPercentFilters, func(inFilter *AlertResourceModelActionFiltersItemConditionsItemEventFrequencyPercentFiltersItem, _ int) apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter {
					return apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter{
						Attribute: inFilter.Attribute.GetPtr(),
						Key:       inFilter.Key.GetPtr(),
						Match:     inFilter.Match.GetPtr(),
						Value:     inFilter.Value.GetPtr(),
					}
				}))
				outEventFrequencyPercent.Comparison.Interval = inEventFrequencyPercent.Interval.Get()
				outEventFrequencyPercent.Comparison.ComparisonInterval = inEventFrequencyPercent.ComparisonInterval.Get()
				outEventFrequencyPercent.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionEventFrequencyPercent(outEventFrequencyPercent); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.PercentSessionsCount.IsKnown():
				inPercentSessionsCount := fwdiag.Merge(inCondition.PercentSessionsCount.Get(ctx))(&diags)
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
				inPercentSessionsPercent := fwdiag.Merge(inCondition.PercentSessionsPercent.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				inSessionsPercentFilters := fwdiag.Merge(inPercentSessionsPercent.Filters.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outPercentSessionsPercent apiclient.OrganizationWorkflowActionFilterConditionPercentSessionsPercent
				outPercentSessionsPercent.Comparison.Value = inPercentSessionsPercent.Value.Get()
				outPercentSessionsPercent.Comparison.Filters = lo.ToPtr(lo.Map(inSessionsPercentFilters, func(inFilter *AlertResourceModelActionFiltersItemConditionsItemPercentSessionsPercentFiltersItem, _ int) apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter {
					return apiclient.OrganizationWorkflowActionFilterConditionEventUniqueUserFrequencyCountFilter{
						Attribute: inFilter.Attribute.GetPtr(),
						Key:       inFilter.Key.GetPtr(),
						Match:     inFilter.Match.GetPtr(),
						Value:     inFilter.Value.GetPtr(),
					}
				}))
				outPercentSessionsPercent.Comparison.Interval = inPercentSessionsPercent.Interval.Get()
				outPercentSessionsPercent.Comparison.ComparisonInterval = inPercentSessionsPercent.ComparisonInterval.Get()
				outPercentSessionsPercent.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionPercentSessionsPercent(outPercentSessionsPercent); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.EventAttribute.IsKnown():
				inEventAttribute := fwdiag.Merge(inCondition.EventAttribute.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEventAttribute apiclient.OrganizationWorkflowActionFilterConditionEventAttribute
				outEventAttribute.Comparison.Attribute = inEventAttribute.Attribute.Get()
				outEventAttribute.Comparison.Match = inEventAttribute.Match.Get()
				outEventAttribute.Comparison.Value = inEventAttribute.Value.GetPtr()
				outEventAttribute.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionEventAttribute(outEventAttribute); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}

			case inCondition.TaggedEvent.IsKnown():
				inTaggedEvent := fwdiag.Merge(inCondition.TaggedEvent.Get(ctx))(&diags)
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
				inLatestAdoptedRelease := fwdiag.Merge(inCondition.LatestAdoptedRelease.Get(ctx))(&diags)
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
				inLevel := fwdiag.Merge(inCondition.Level.Get(ctx))(&diags)
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

			case inCondition.IssueType.IsKnown():
				inIssueType := fwdiag.Merge(inCondition.IssueType.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outIssueType apiclient.OrganizationWorkflowActionFilterConditionIssueType
				outIssueType.Comparison.Value = inIssueType.Value.Get()
				outIssueType.Comparison.Include = inIssueType.Include.Get()
				outIssueType.ConditionResult = true

				if err := outCondition.FromOrganizationWorkflowActionFilterConditionIssueType(outIssueType); err != nil {
					diags.AddError("Failed to create condition", err.Error())
					return nil, diags
				}
			}
			outConditions = append(outConditions, outCondition)
		}

		// Actions
		inActions := fwdiag.Merge(inActionFilter.Actions.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}

		var outActions []apiclient.OrganizationWorkflowActionFilterAction
		for _, inAction := range inActions {
			var outAction apiclient.OrganizationWorkflowActionFilterAction
			switch {
			case inAction.Email.IsKnown():
				inEmail := fwdiag.Merge(inAction.Email.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outEmail apiclient.OrganizationWorkflowActionFilterActionEmail
				outEmail.Config.TargetType = apiclient.OrganizationWorkflowActionFilterActionEmailConfigTargetType(inEmail.TargetType.Get())
				outEmail.Config.TargetIdentifier = inEmail.TargetId.GetPtr()
				if inEmail.FallthroughType.IsKnown() {
					outEmail.Data.FallthroughType = new(apiclient.OrganizationWorkflowActionFilterActionEmailDataFallthroughType(inEmail.FallthroughType.Get()))
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
				inSlack := fwdiag.Merge(inAction.Slack.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outSlack apiclient.OrganizationWorkflowActionFilterActionSlack
				outSlack.IntegrationId = inSlack.IntegrationId.Get()
				outSlack.Config.TargetType = "specific"
				outSlack.Config.TargetIdentifier = inSlack.ChannelId.Get()
				outSlack.Config.TargetDisplay = inSlack.ChannelName.ValueString()
				if inSlack.Tags.IsKnown() {
					outSlack.Data.Tags = new(inSlack.Tags.Get())
				}
				if inSlack.Notes.IsKnown() {
					outSlack.Data.Notes = new(inSlack.Notes.Get())
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionSlack(outSlack); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Pagerduty.IsKnown():
				inPagerduty := fwdiag.Merge(inAction.Pagerduty.Get(ctx))(&diags)
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
				inDiscord := fwdiag.Merge(inAction.Discord.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outDiscord apiclient.OrganizationWorkflowActionFilterActionDiscord
				outDiscord.IntegrationId = inDiscord.IntegrationId.Get()
				outDiscord.Config.TargetType = "specific"
				outDiscord.Config.TargetIdentifier = inDiscord.ChannelId.Get()
				if inDiscord.Tags.IsKnown() {
					outDiscord.Data.Tags = new(inDiscord.Tags.Get())
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionDiscord(outDiscord); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Msteams.IsKnown():
				inMsteams := fwdiag.Merge(inAction.Msteams.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outMsteams apiclient.OrganizationWorkflowActionFilterActionMsTeams
				outMsteams.IntegrationId = inMsteams.IntegrationId.Get()
				outMsteams.Config.TargetType = "specific"
				outMsteams.Config.TargetIdentifier = inMsteams.TeamId.ValueString()
				outMsteams.Config.TargetDisplay = inMsteams.ChannelName.Get()
				outMsteams.Data = map[string]interface{}{}

				if err := outAction.FromOrganizationWorkflowActionFilterActionMsTeams(outMsteams); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Opsgenie.IsKnown():
				inOpsgenie := fwdiag.Merge(inAction.Opsgenie.Get(ctx))(&diags)
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
				inVsts := fwdiag.Merge(inAction.Vsts.Get(ctx))(&diags)
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
				inJira := fwdiag.Merge(inAction.Jira.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outJira apiclient.OrganizationWorkflowActionFilterActionJira
				outJira.IntegrationId = inJira.IntegrationId.Get()
				outJira.Config.TargetType = "specific"
				outJira.Data.AdditionalFields.Project = inJira.Project.Get()
				outJira.Data.AdditionalFields.Issuetype = inJira.IssueType.Get()
				outJira.Data.AdditionalFields.Labels = inJira.Labels.GetPtr()
				outJira.Data.AdditionalFields.Priority = inJira.Priority.GetPtr()
				outJira.Data.AdditionalFields.Reporter = inJira.Reporter.GetPtr()
				if inJira.Components.IsKnown() {
					components := fwdiag.Merge(inJira.Components.Get(ctx))(&diags)
					outJira.Data.AdditionalFields.Components = &components
				}
				if inJira.AdditionalFields.IsKnown() {
					inAdditionalFields := fwdiag.Merge(inJira.AdditionalFields.Get(ctx))(&diags)
					for k, v := range inAdditionalFields {
						outJira.Data.AdditionalFields.Set(k, v)
					}
				}
				if diags.HasError() {
					return nil, diags
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionJira(outJira); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.JiraServer.IsKnown():
				inJiraServer := fwdiag.Merge(inAction.JiraServer.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outJiraServer apiclient.OrganizationWorkflowActionFilterActionJiraServer
				outJiraServer.IntegrationId = inJiraServer.IntegrationId.Get()
				outJiraServer.Config.TargetType = "specific"
				outJiraServer.Data.AdditionalFields.Project = inJiraServer.Project.Get()
				outJiraServer.Data.AdditionalFields.Issuetype = inJiraServer.IssueType.Get()
				outJiraServer.Data.AdditionalFields.Labels = inJiraServer.Labels.GetPtr()
				outJiraServer.Data.AdditionalFields.Priority = inJiraServer.Priority.GetPtr()
				outJiraServer.Data.AdditionalFields.Reporter = inJiraServer.Reporter.GetPtr()
				if inJiraServer.Components.IsKnown() {
					components := fwdiag.Merge(inJiraServer.Components.Get(ctx))(&diags)
					outJiraServer.Data.AdditionalFields.Components = &components
				}
				if inJiraServer.AdditionalFields.IsKnown() {
					inAdditionalFields := fwdiag.Merge(inJiraServer.AdditionalFields.Get(ctx))(&diags)
					for k, v := range inAdditionalFields {
						outJiraServer.Data.AdditionalFields.Set(k, v)
					}
				}
				if diags.HasError() {
					return nil, diags
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionJiraServer(outJiraServer); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Github.IsKnown():
				inGithub := fwdiag.Merge(inAction.Github.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outGithub apiclient.OrganizationWorkflowActionFilterActionGitHub
				outGithub.IntegrationId = inGithub.IntegrationId.Get()
				outGithub.Config.TargetType = "specific"
				outGithub.Data.AdditionalFields.Repo = inGithub.Repo.Get()
				outGithub.Data.AdditionalFields.Assignee = inGithub.Assignee.Get()
				outGithub.Data.AdditionalFields.Labels = fwdiag.Merge(inGithub.Labels.Get(ctx))(&diags)
				outGithub.Data.AdditionalFields.Integration = inGithub.IntegrationId.Get()
				if diags.HasError() {
					return nil, diags
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionGitHub(outGithub); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.SentryApp.IsKnown():
				inSentryApp := fwdiag.Merge(inAction.SentryApp.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outSentryApp apiclient.OrganizationWorkflowActionFilterActionSentryApp
				outSentryApp.Config.TargetType = "sentry_app"
				outSentryApp.Config.TargetIdentifier = inSentryApp.SentryAppId.Get()

				if inSentryApp.Settings.IsKnown() {
					inSettings := fwdiag.Merge(inSentryApp.Settings.Get(ctx))(&diags)
					if diags.HasError() {
						return nil, diags
					}
					settings := make([]struct {
						Label *string `json:"label,omitempty"`
						Name  string  `json:"name"`
						Value string  `json:"value"`
					}, 0, len(inSettings))
					for _, s := range inSettings {
						settings = append(settings, struct {
							Label *string `json:"label,omitempty"`
							Name  string  `json:"name"`
							Value string  `json:"value"`
						}{
							Name:  s.Name.Get(),
							Value: s.Value.Get(),
							Label: s.Label.GetPtr(),
						})
					}
					outSentryApp.Data.Settings = &settings
				}

				if err := outAction.FromOrganizationWorkflowActionFilterActionSentryApp(outSentryApp); err != nil {
					diags.AddError("Failed to create action", err.Error())
					return nil, diags
				}

			case inAction.Webhook.IsKnown():
				inWebhook := fwdiag.Merge(inAction.Webhook.Get(ctx))(&diags)
				if diags.HasError() {
					return nil, diags
				}

				var outWebhook apiclient.OrganizationWorkflowActionFilterActionWebhook
				outWebhook.Data = map[string]any{}
				outWebhook.Config.TargetIdentifier = inWebhook.Service.Get()

				if err := outAction.FromOrganizationWorkflowActionFilterActionWebhook(outWebhook); err != nil {
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

	inTriggerConditions := fwdiag.Merge(data.TriggerConditions.Get(ctx))(&diags)
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
		case triggerCondition.EventFrequencyCount.IsKnown():
			in := fwdiag.Merge(triggerCondition.EventFrequencyCount.Get(ctx))(&diags)
			if diags.HasError() {
				return nil, diags
			}

			comparison := map[string]any{
				"interval": in.Interval.Get(),
				"value":    in.Value.Get(),
			}

			if err := outTriggerCondition.Comparison.FromOrganizationWorkflowTriggerConditionComparison1(comparison); err != nil {
				diags.AddError("Failed to create event_frequency_count trigger condition", err.Error())
				return nil, diags
			}
			outTriggerCondition.Type = "event_frequency_count"
		}

		outTriggerConditions = append(outTriggerConditions, outTriggerCondition)
	}

	if data.LegacyTriggerConditions.IsKnown() {
		inLegacyTriggerConditions := fwdiag.Merge(data.LegacyTriggerConditions.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}

		for _, inLegacyTriggerCondition := range inLegacyTriggerConditions {
			var comp apiclient.OrganizationWorkflowTriggerCondition_Comparison
			if err := comp.FromOrganizationWorkflowTriggerConditionComparison0(true); err != nil {
				diags.AddError("Failed to build legacy trigger condition", err.Error())
				return nil, diags
			}
			outTriggerConditions = append(outTriggerConditions, apiclient.OrganizationWorkflowTriggerCondition{
				Type:            inLegacyTriggerCondition,
				Comparison:      comp,
				ConditionResult: true,
			})
		}
	}

	return outTriggerConditions, diags
}

func parseEventFrequencyCountTriggerComparison(comparison map[string]any) (string, int64, error) {
	interval, ok := comparison["interval"].(string)
	if !ok {
		return "", 0, fmt.Errorf("expected interval to be a string, got %T", comparison["interval"])
	}

	value, ok := comparison["value"].(float64)
	if !ok {
		return "", 0, fmt.Errorf("expected value to be a number, got %T", comparison["value"])
	} else if value < 0 || value >= float64(math.MaxInt64) || math.Trunc(value) != value {
		return "", 0, fmt.Errorf("expected value to be a non-negative integer, got %v", value)
	}

	if rawFilters, exists := comparison["filters"]; exists && rawFilters != nil {
		filters, ok := rawFilters.([]any)
		if !ok {
			return "", 0, fmt.Errorf("expected filters to be a list, got %T", rawFilters)
		}
		if len(filters) > 0 {
			return "", 0, fmt.Errorf("event_frequency_count filters are not supported")
		}
	}

	return interval, int64(value), nil
}

func (r *AlertResource) getCreateJSONRequestBody(ctx context.Context, data AlertResourceModel) (*apiclient.CreateOrganizationWorkflowJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	monitorIds := fwdiag.Merge(data.MonitorIds.Get(ctx))(&diags)
	if diags.HasError() {
		return nil, diags
	}

	triggerConditions := append(
		[]apiclient.OrganizationWorkflowTriggerCondition{},
		fwdiag.Merge(r.getTriggerConditions(ctx, data))(&diags)...,
	)

	req := apiclient.CreateOrganizationWorkflowJSONRequestBody{
		Name:        data.Name.Get(),
		Enabled:     data.Enabled.Get(),
		Environment: nullableFromPtr(data.Environment.GetPtr()),
		Config: apiclient.OrganizationWorkflowConfig{
			Frequency: data.FrequencyMinutes.Get(),
		},
		DetectorIds: monitorIds,
		Triggers: apiclient.OrganizationWorkflowTrigger{
			LogicType:  apiclient.OrganizationWorkflowTriggerLogicTypeAnyShort,
			Conditions: triggerConditions,
		},
		ActionFilters: fwdiag.Merge(r.getActionFilters(ctx, data))(&diags),
	}

	return &req, diags
}

func (r *AlertResource) getUpdateJSONRequestBody(ctx context.Context, data AlertResourceModel) (*apiclient.UpdateOrganizationWorkflowJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	monitorIds := fwdiag.Merge(data.MonitorIds.Get(ctx))(&diags)
	if diags.HasError() {
		return nil, diags
	}

	triggerConditions := append(
		[]apiclient.OrganizationWorkflowTriggerCondition{},
		fwdiag.Merge(r.getTriggerConditions(ctx, data))(&diags)...,
	)

	req := apiclient.UpdateOrganizationWorkflowJSONRequestBody{
		Id:          data.Id.Get(),
		Name:        data.Name.Get(),
		Enabled:     data.Enabled.Get(),
		Environment: nullableFromPtr(data.Environment.GetPtr()),
		Config: apiclient.OrganizationWorkflowConfig{
			Frequency: data.FrequencyMinutes.Get(),
		},
		DetectorIds: monitorIds,
		Triggers: apiclient.OrganizationWorkflowTrigger{
			LogicType:  apiclient.OrganizationWorkflowTriggerLogicTypeAnyShort,
			Conditions: triggerConditions,
		},
		ActionFilters: fwdiag.Merge(r.getActionFilters(ctx, data))(&diags),
	}

	return &req, diags
}

func (m *AlertResourceModel) Fill(ctx context.Context, data apiclient.OrganizationWorkflow) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringValue(data.Name)
	m.Enabled = supertypes.NewBoolValue(data.Enabled)
	if v, err := data.Environment.Get(); err == nil {
		m.Environment = supertypes.NewStringValueOrNull(v)
	} else {
		m.Environment = supertypes.NewStringNull()
	}
	m.FrequencyMinutes = supertypes.NewInt64Value(data.Config.Frequency)
	m.MonitorIds = supertypes.NewSetValueOfSlice(ctx, data.DetectorIds)

	triggers, err := data.Triggers.AsOrganizationWorkflowTrigger()
	if err != nil {
		diags.AddError("Failed to parse triggers", err.Error())
		return diags
	}

	// NOTE: The API returns conditions in a random order, so we need to sort them by ID to ensure that the
	// order is deterministic.
	slices.SortFunc(triggers.Conditions, func(a, b apiclient.OrganizationWorkflowTriggerCondition) int {
		aId := int(must.Get(a.Id.Int64()))
		bId := int(must.Get(b.Id.Int64()))
		return aId - bId
	})

	triggerConditions := []AlertResourceModelTriggerConditionsItem{}
	var legacyTriggerConditions []string
	for _, triggerCondition := range triggers.Conditions {
		outTriggerCondition := AlertResourceModelTriggerConditionsItem{
			FirstSeenEvent:       supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemFirstSeenEvent](ctx),
			IssueResolvedTrigger: supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemIssueResolvedTrigger](ctx),
			ReappearedEvent:      supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemReappearedEvent](ctx),
			RegressionEvent:      supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemRegressionEvent](ctx),
			EventFrequencyCount:  supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventFrequencyCount](ctx),
		}
		switch triggerCondition.Type {
		case "first_seen_event":
			outTriggerCondition.FirstSeenEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemFirstSeenEvent{})
			triggerConditions = append(triggerConditions, outTriggerCondition)
		case "issue_resolved_trigger":
			outTriggerCondition.IssueResolvedTrigger = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemIssueResolvedTrigger{})
			triggerConditions = append(triggerConditions, outTriggerCondition)
		case "reappeared_event":
			outTriggerCondition.ReappearedEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemReappearedEvent{})
			triggerConditions = append(triggerConditions, outTriggerCondition)
		case "regression_event":
			outTriggerCondition.RegressionEvent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemRegressionEvent{})
			triggerConditions = append(triggerConditions, outTriggerCondition)
		case "event_frequency_count":
			comparison, err := triggerCondition.Comparison.AsOrganizationWorkflowTriggerConditionComparison1()
			if err != nil {
				if _, boolErr := triggerCondition.Comparison.AsOrganizationWorkflowTriggerConditionComparison0(); boolErr == nil {
					legacyTriggerConditions = append(legacyTriggerConditions, triggerCondition.Type)
					continue
				}
				diags.AddError("Failed to parse event_frequency_count trigger condition", err.Error())
				return diags
			}
			interval, value, err := parseEventFrequencyCountTriggerComparison(comparison)
			if err != nil {
				diags.AddError("Failed to parse event_frequency_count trigger condition", err.Error())
				return diags
			}
			outTriggerCondition.EventFrequencyCount = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemEventFrequencyCount{
				Interval: supertypes.NewStringValue(interval),
				Value:    supertypes.NewInt64Value(value),
			})
			triggerConditions = append(triggerConditions, outTriggerCondition)
		default:
			legacyTriggerConditions = append(legacyTriggerConditions, triggerCondition.Type)
		}
	}
	m.TriggerConditions = supertypes.NewListNestedObjectValueOfValueSlice(ctx, triggerConditions)
	if len(legacyTriggerConditions) == 0 {
		m.LegacyTriggerConditions.SetNull(ctx)
	} else {
		diags.Append(m.LegacyTriggerConditions.Set(ctx, legacyTriggerConditions)...)
	}

	actionFilters, err := data.ActionFilters.AsOrganizationWorkflowActionFilters0()
	if err != nil {
		diags.AddError("Failed to parse action filters", err.Error())
		return diags
	}

	var outActionFilters []AlertResourceModelActionFiltersItem
	for _, actionFilter := range actionFilters {
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
				IssueType:                     supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemConditionsItemIssueType](ctx),
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
				if conditionValue.Comparison.Include != nil {
					issueCategory.Include = supertypes.NewBoolValue(*conditionValue.Comparison.Include)
				} else {
					issueCategory.Include = supertypes.NewBoolValue(true)
				}

				outCondition.IssueCategory = supertypes.NewSingleNestedObjectValueOf(ctx, &issueCategory)

			case apiclient.OrganizationWorkflowActionFilterConditionIssueOccurrences:
				var issueOccurrences AlertResourceModelActionFiltersItemConditionsItemIssueOccurrences
				issueOccurrences.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)

				outCondition.IssueOccurrences = supertypes.NewSingleNestedObjectValueOf(ctx, &issueOccurrences)

			case apiclient.OrganizationWorkflowActionFilterConditionIssuePriorityDeescalating:
				var issuePriorityDeescalating AlertResourceModelActionFiltersItemConditionsItemIssuePriorityDeescalating
				issuePriorityDeescalating.Comparison = supertypes.NewInt64Value(conditionValue.Comparison)

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

				outFrequencyCountFilters := []AlertResourceModelActionFiltersItemConditionsItemEventFrequencyCountFiltersItem{}
				if conditionValue.Comparison.Filters != nil {
					for _, filter := range *conditionValue.Comparison.Filters {
						outFrequencyCountFilters = append(outFrequencyCountFilters, AlertResourceModelActionFiltersItemConditionsItemEventFrequencyCountFiltersItem{
							Attribute: supertypes.NewStringPointerValueOrNull(filter.Attribute),
							Key:       supertypes.NewStringPointerValueOrNull(filter.Key),
							Match:     supertypes.NewStringPointerValueOrNull(filter.Match),
							Value:     supertypes.NewStringPointerValueOrNull(filter.Value),
						})
					}
				}
				eventFrequencyCount.Filters = supertypes.NewListNestedObjectValueOfValueSlice(ctx, outFrequencyCountFilters)

				outCondition.EventFrequencyCount = supertypes.NewSingleNestedObjectValueOf(ctx, &eventFrequencyCount)

			case apiclient.OrganizationWorkflowActionFilterConditionEventFrequencyPercent:
				var eventFrequencyPercent AlertResourceModelActionFiltersItemConditionsItemEventFrequencyPercent
				eventFrequencyPercent.Value = supertypes.NewInt64Value(conditionValue.Comparison.Value)
				eventFrequencyPercent.Interval = supertypes.NewStringValue(conditionValue.Comparison.Interval)
				eventFrequencyPercent.ComparisonInterval = supertypes.NewStringValue(conditionValue.Comparison.ComparisonInterval)

				outFrequencyPercentFilters := []AlertResourceModelActionFiltersItemConditionsItemEventFrequencyPercentFiltersItem{}
				if conditionValue.Comparison.Filters != nil {
					for _, filter := range *conditionValue.Comparison.Filters {
						outFrequencyPercentFilters = append(outFrequencyPercentFilters, AlertResourceModelActionFiltersItemConditionsItemEventFrequencyPercentFiltersItem{
							Attribute: supertypes.NewStringPointerValueOrNull(filter.Attribute),
							Key:       supertypes.NewStringPointerValueOrNull(filter.Key),
							Match:     supertypes.NewStringPointerValueOrNull(filter.Match),
							Value:     supertypes.NewStringPointerValueOrNull(filter.Value),
						})
					}
				}
				eventFrequencyPercent.Filters = supertypes.NewListNestedObjectValueOfValueSlice(ctx, outFrequencyPercentFilters)

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

				outSessionsPercentFilters := []AlertResourceModelActionFiltersItemConditionsItemPercentSessionsPercentFiltersItem{}
				if conditionValue.Comparison.Filters != nil {
					for _, filter := range *conditionValue.Comparison.Filters {
						outSessionsPercentFilters = append(outSessionsPercentFilters, AlertResourceModelActionFiltersItemConditionsItemPercentSessionsPercentFiltersItem{
							Attribute: supertypes.NewStringPointerValueOrNull(filter.Attribute),
							Key:       supertypes.NewStringPointerValueOrNull(filter.Key),
							Match:     supertypes.NewStringPointerValueOrNull(filter.Match),
							Value:     supertypes.NewStringPointerValueOrNull(filter.Value),
						})
					}
				}
				percentSessionsPercent.Filters = supertypes.NewListNestedObjectValueOfValueSlice(ctx, outSessionsPercentFilters)

				outCondition.PercentSessionsPercent = supertypes.NewSingleNestedObjectValueOf(ctx, &percentSessionsPercent)

			case apiclient.OrganizationWorkflowActionFilterConditionEventAttribute:
				var eventAttribute AlertResourceModelActionFiltersItemConditionsItemEventAttribute
				eventAttribute.Attribute = supertypes.NewStringValue(conditionValue.Comparison.Attribute)
				eventAttribute.Match = supertypes.NewStringValue(conditionValue.Comparison.Match)
				eventAttribute.Value = supertypes.NewStringPointerValueOrNull(conditionValue.Comparison.Value)

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

			case apiclient.OrganizationWorkflowActionFilterConditionIssueType:
				var issueType AlertResourceModelActionFiltersItemConditionsItemIssueType
				issueType.Value = supertypes.NewStringValue(conditionValue.Comparison.Value)
				issueType.Include = supertypes.NewBoolValue(conditionValue.Comparison.Include)

				outCondition.IssueType = supertypes.NewSingleNestedObjectValueOf(ctx, &issueType)
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
				SentryApp:  supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemSentryApp](ctx),
				Webhook:    supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemWebhook](ctx),
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
				if actionValue.Config.TargetIdentifier != nil && *actionValue.Config.TargetIdentifier != "" {
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
				outSlack.ChannelName = sentrytypes.NewSlackChannelValue(actionValue.Config.TargetDisplay)
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
				outMsteams.TeamId = sentrytypes.NewMsTeamsTeamIdValue(actionValue.Config.TargetIdentifier)
				outMsteams.ChannelName = supertypes.NewStringValue(actionValue.Config.TargetDisplay)
				outMsteams.TeamThreadId = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)

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
				outJira.Labels = newOptionalTicketString(actionValue.Data.AdditionalFields.Labels)
				outJira.Priority = newOptionalTicketString(actionValue.Data.AdditionalFields.Priority)
				outJira.Reporter = newOptionalTicketString(actionValue.Data.AdditionalFields.Reporter)
				outJira.Components = newOptionalTicketStringSet(ctx, actionValue.Data.AdditionalFields.Components)
				outJira.AdditionalFields = newOptionalTicketAdditionalFields(ctx, actionValue.Data.AdditionalFields.AdditionalProperties, &diags)

				outAction.Jira = supertypes.NewSingleNestedObjectValueOf(ctx, &outJira)

			case apiclient.OrganizationWorkflowActionFilterActionJiraServer:
				var outJiraServer AlertResourceModelActionFiltersItemActionsItemJiraServer
				outJiraServer.IntegrationId = supertypes.NewStringValue(actionValue.IntegrationId)
				outJiraServer.Project = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Project)
				outJiraServer.IssueType = supertypes.NewStringValue(actionValue.Data.AdditionalFields.Issuetype)
				outJiraServer.Labels = newOptionalTicketString(actionValue.Data.AdditionalFields.Labels)
				outJiraServer.Priority = newOptionalTicketString(actionValue.Data.AdditionalFields.Priority)
				outJiraServer.Reporter = newOptionalTicketString(actionValue.Data.AdditionalFields.Reporter)
				outJiraServer.Components = newOptionalTicketStringSet(ctx, actionValue.Data.AdditionalFields.Components)
				outJiraServer.AdditionalFields = newOptionalTicketAdditionalFields(ctx, actionValue.Data.AdditionalFields.AdditionalProperties, &diags)

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

			case apiclient.OrganizationWorkflowActionFilterActionSentryApp:
				var outSentryApp AlertResourceModelActionFiltersItemActionsItemSentryApp
				outSentryApp.SentryAppId = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)

				if actionValue.Data.Settings != nil {
					items := []AlertResourceModelActionFiltersItemActionsItemSentryAppSettingsItem{}
					for _, s := range *actionValue.Data.Settings {
						items = append(items, AlertResourceModelActionFiltersItemActionsItemSentryAppSettingsItem{
							Name:  supertypes.NewStringValue(s.Name),
							Value: supertypes.NewStringValue(s.Value),
							Label: supertypes.NewStringPointerValueOrNull(s.Label),
						})
					}
					outSentryApp.Settings = supertypes.NewListNestedObjectValueOfValueSlice(ctx, items)
				} else {
					outSentryApp.Settings = supertypes.NewListNestedObjectValueOfNull[AlertResourceModelActionFiltersItemActionsItemSentryAppSettingsItem](ctx)
				}

				outAction.SentryApp = supertypes.NewSingleNestedObjectValueOf(ctx, &outSentryApp)

			case apiclient.OrganizationWorkflowActionFilterActionWebhook:
				var outWebhook AlertResourceModelActionFiltersItemActionsItemWebhook
				outWebhook.Service = supertypes.NewStringValue(actionValue.Config.TargetIdentifier)

				outAction.Webhook = supertypes.NewSingleNestedObjectValueOf(ctx, &outWebhook)
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

// newOptionalTicketString converts an optional ticket `additionalFields` string
// into a Terraform value. Sentry stores unset ticket fields inconsistently --
// the key may be absent, or present with an empty string -- and both mean "not
// set". Normalising both to null keeps an omitted attribute from showing
// permanent drift after apply.
func newOptionalTicketString(value *string) supertypes.StringValue {
	if value == nil || *value == "" {
		return supertypes.NewStringNull()
	}
	return supertypes.NewStringValue(*value)
}

// newOptionalTicketStringSet converts an optional ticket `additionalFields`
// string list into a Terraform set, mapping both an absent key and an empty
// list to null for the same reason as newOptionalTicketString.
func newOptionalTicketStringSet(ctx context.Context, values *[]string) supertypes.SetValueOf[string] {
	if values == nil || len(*values) == 0 {
		return supertypes.NewSetValueOfNull[string](ctx)
	}
	return supertypes.NewSetValueOfSlice(ctx, *values)
}

// newOptionalTicketAdditionalFields converts the passthrough `additionalFields`
// keys that have no dedicated attribute into a Terraform map. Values are
// stringified because the schema models them as a map of strings; Sentry echoes
// back whatever JSON scalar it was given, so a number written as "3" returns as
// 3 and must be normalised to avoid drift.
func newOptionalTicketAdditionalFields(ctx context.Context, values map[string]interface{}, diags *diag.Diagnostics) supertypes.MapValueOf[string] {
	if len(values) == 0 {
		return supertypes.NewMapValueOfNull[string](ctx)
	}

	out := make(map[string]string, len(values))
	for k, v := range values {
		switch tv := v.(type) {
		case nil:
			continue
		case string:
			out[k] = tv
		case json.Number:
			out[k] = tv.String()
		case float64:
			out[k] = strconv.FormatFloat(tv, 'f', -1, 64)
		case bool:
			out[k] = strconv.FormatBool(tv)
		default:
			b, err := json.Marshal(tv)
			if err != nil {
				continue
			}
			out[k] = string(b)
		}
	}

	if len(out) == 0 {
		return supertypes.NewMapValueOfNull[string](ctx)
	}
	v, d := supertypes.NewMapValueOfMap(ctx, out)
	diags.Append(d...)
	return v
}
