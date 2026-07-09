package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/mzglinski/terraform-provider-sentry/internal/sentrytypes"
)

// reorderToMatchPrior reorders incoming items to match the order of prior items by key.
// Items in incoming that have no matching prior key are appended at the end in their
// original order. Items in prior with no match in incoming are silently skipped.
//
// This is used after reading from the Sentry API to stabilize list order: the API may
// return conditions/filters/actions in a different order than was planned or previously
// stored, causing "Provider produced inconsistent result after apply" errors.
func reorderToMatchPrior[T any, K comparable](prior, incoming []T, key func(T) K) []T {
	if len(prior) == 0 {
		return incoming
	}

	type entry struct {
		item T
		used bool
	}

	pool := make(map[K][]*entry)
	entries := make([]*entry, len(incoming))
	for i, item := range incoming {
		e := &entry{item: item}
		entries[i] = e
		k := key(item)
		pool[k] = append(pool[k], e)
	}

	result := make([]T, 0, len(incoming))
	for _, p := range prior {
		k := key(p)
		for _, e := range pool[k] {
			if !e.used {
				result = append(result, e.item)
				e.used = true
				break
			}
		}
	}
	for _, e := range entries {
		if !e.used {
			result = append(result, e.item)
		}
	}
	return result
}

func issueAlertConditionModelKey(ctx context.Context) func(m IssueAlertConditionModel) string {
	return func(m IssueAlertConditionModel) string {
		switch {
		case m.FirstSeenEvent.IsKnown():
			return "first_seen_event"
		case m.RegressionEvent.IsKnown():
			return "regression_event"
		case m.ReappearedEvent.IsKnown():
			return "reappeared_event"
		case m.NewHighPriorityIssue.IsKnown():
			return "new_high_priority_issue"
		case m.ExistingHighPriorityIssue.IsKnown():
			return "existing_high_priority_issue"
		case m.EventFrequency.IsKnown():
			f := m.EventFrequency.MustGet(ctx)
			return fmt.Sprintf("event_frequency\x00%s\x00%s\x00%d\x00%s",
				f.ComparisonType.ValueString(), f.Interval.ValueString(),
				f.Value.ValueInt64(), f.ComparisonInterval.ValueString())
		case m.EventUniqueUserFrequency.IsKnown():
			f := m.EventUniqueUserFrequency.MustGet(ctx)
			return fmt.Sprintf("event_unique_user_frequency\x00%s\x00%s\x00%d\x00%s",
				f.ComparisonType.ValueString(), f.Interval.ValueString(),
				f.Value.ValueInt64(), f.ComparisonInterval.ValueString())
		case m.EventFrequencyPercent.IsKnown():
			f := m.EventFrequencyPercent.MustGet(ctx)
			return fmt.Sprintf("event_frequency_percent\x00%s\x00%s\x00%g\x00%s",
				f.ComparisonType.ValueString(), f.Interval.ValueString(),
				f.Value.ValueFloat64(), f.ComparisonInterval.ValueString())
		default:
			return ""
		}
	}
}

func issueAlertFilterModelKey(ctx context.Context) func(m IssueAlertFilterModel) string {
	return func(m IssueAlertFilterModel) string {
		switch {
		case m.AgeComparison.IsKnown():
			f := m.AgeComparison.MustGet(ctx)
			return fmt.Sprintf("age_comparison\x00%s\x00%d\x00%s",
				f.ComparisonType.ValueString(), f.Value.ValueInt64(), f.Time.ValueString())
		case m.IssueOccurrences.IsKnown():
			f := m.IssueOccurrences.MustGet(ctx)
			return fmt.Sprintf("issue_occurrences\x00%d", f.Value.ValueInt64())
		case m.AssignedTo.IsKnown():
			f := m.AssignedTo.MustGet(ctx)
			return fmt.Sprintf("assigned_to\x00%s\x00%s",
				f.TargetType.ValueString(), f.TargetIdentifier.ValueString())
		case m.LatestAdoptedRelease.IsKnown():
			f := m.LatestAdoptedRelease.MustGet(ctx)
			return fmt.Sprintf("latest_adopted_release\x00%s\x00%s\x00%s",
				f.OldestOrNewest.ValueString(), f.OlderOrNewer.ValueString(), f.Environment.ValueString())
		case m.LatestRelease.IsKnown():
			return "latest_release"
		case m.IssueCategory.IsKnown():
			f := m.IssueCategory.MustGet(ctx)
			return fmt.Sprintf("issue_category\x00%s", f.Value.ValueString())
		case m.EventAttribute.IsKnown():
			f := m.EventAttribute.MustGet(ctx)
			return fmt.Sprintf("event_attribute\x00%s\x00%s\x00%s",
				f.Attribute.ValueString(), f.Match.ValueString(), f.Value.ValueString())
		case m.TaggedEvent.IsKnown():
			f := m.TaggedEvent.MustGet(ctx)
			return fmt.Sprintf("tagged_event\x00%s\x00%s\x00%s",
				f.Key.ValueString(), f.Match.ValueString(), f.Value.ValueString())
		case m.Level.IsKnown():
			f := m.Level.MustGet(ctx)
			return fmt.Sprintf("level\x00%s\x00%s", f.Match.ValueString(), f.Level.ValueString())
		default:
			return ""
		}
	}
}

func issueAlertActionModelKey(ctx context.Context) func(m IssueAlertActionModel) string {
	return func(m IssueAlertActionModel) string {
		switch {
		case m.NotifyEmail.IsKnown():
			f := m.NotifyEmail.MustGet(ctx)
			return fmt.Sprintf("notify_email\x00%s\x00%s\x00%s",
				f.TargetType.ValueString(), f.TargetType.ValueString(), f.TargetIdentifier.ValueString())
		case m.NotifyEvent.IsKnown():
			return "notify_event"
		case m.NotifyEventService.IsKnown():
			f := m.NotifyEventService.MustGet(ctx)
			return fmt.Sprintf("notify_event_service\x00%s", f.Service.ValueString())
		case m.NotifyEventSentryApp.IsKnown():
			f := m.NotifyEventSentryApp.MustGet(ctx)
			return fmt.Sprintf("notify_event_sentry_app\x00%s", f.SentryAppInstallationUuid.ValueString())
		case m.OpsgenieNotifyTeam.IsKnown():
			f := m.OpsgenieNotifyTeam.MustGet(ctx)
			return fmt.Sprintf("opsgenie_notify_team\x00%s\x00%s\x00%s",
				f.Account.ValueString(), f.Team.ValueString(), f.Priority.ValueString())
		case m.PagerDutyNotifyService.IsKnown():
			f := m.PagerDutyNotifyService.MustGet(ctx)
			return fmt.Sprintf("pagerduty_notify_service\x00%s\x00%s\x00%s",
				f.Account.ValueString(), f.Service.ValueString(), f.Severity.ValueString())
		case m.SlackNotifyService.IsKnown():
			f := m.SlackNotifyService.MustGet(ctx)
			return fmt.Sprintf("slack_notify_service\x00%s\x00%s\x00%s",
				f.Workspace.ValueString(), f.Channel.ValueString(), stringSetKey(f.Tags))
		case m.MsTeamsNotifyService.IsKnown():
			f := m.MsTeamsNotifyService.MustGet(ctx)
			return fmt.Sprintf("msteams_notify_service\x00%s\x00%s",
				f.Team.ValueString(), f.Channel.ValueString())
		case m.DiscordNotifyService.IsKnown():
			f := m.DiscordNotifyService.MustGet(ctx)
			return fmt.Sprintf("discord_notify_service\x00%s\x00%s\x00%s",
				f.Server.ValueString(), f.ChannelId.ValueString(), stringSetKey(f.Tags))
		case m.JiraCreateTicket.IsKnown():
			f := m.JiraCreateTicket.MustGet(ctx)
			return fmt.Sprintf("jira_create_ticket\x00%s\x00%s\x00%s",
				f.Integration.ValueString(), f.Project.ValueString(), f.IssueType.ValueString())
		case m.JiraServerCreateTicket.IsKnown():
			f := m.JiraServerCreateTicket.MustGet(ctx)
			return fmt.Sprintf("jira_server_create_ticket\x00%s\x00%s\x00%s",
				f.Integration.ValueString(), f.Project.ValueString(), f.IssueType.ValueString())
		case m.GitHubCreateTicket.IsKnown():
			f := m.GitHubCreateTicket.MustGet(ctx)
			labels := f.Labels.MustGet(ctx)
			sort.Strings(labels)
			return fmt.Sprintf("github_create_ticket\x00%s\x00%s\x00%s\x00%s",
				f.Integration.ValueString(), f.Repo.ValueString(),
				f.Assignee.ValueString(), strings.Join(labels, ","))
		case m.GitHubEnterpriseCreateTicket.IsKnown():
			f := m.GitHubEnterpriseCreateTicket.MustGet(ctx)
			labels := f.Labels.MustGet(ctx)
			sort.Strings(labels)
			return fmt.Sprintf("github_enterprise_create_ticket\x00%s\x00%s\x00%s\x00%s",
				f.Integration.ValueString(), f.Repo.ValueString(),
				f.Assignee.ValueString(), strings.Join(labels, ","))
		case m.AzureDevopsCreateTicket.IsKnown():
			f := m.AzureDevopsCreateTicket.MustGet(ctx)
			return fmt.Sprintf("azure_devops_create_ticket\x00%s\x00%s\x00%s",
				f.Integration.ValueString(), f.Project.ValueString(), f.WorkItemType.ValueString())
		default:
			return ""
		}
	}
}

// stringSetKey produces a stable string key from a sentrytypes.StringSet.
func stringSetKey(s sentrytypes.StringSet) string {
	if s.IsNull() || s.IsUnknown() {
		return ""
	}
	strs := make([]string, 0, len(s.Elements()))
	for _, e := range s.Elements() {
		strs = append(strs, e.String())
	}
	sort.Strings(strs)
	return strings.Join(strs, ",")
}

// legacyJsonItemKey extracts a stable identity key from a raw JSON object
// by reading its "id" field. Used to reorder legacy conditions/filters/actions
// JSON arrays to match prior state after an API read.
func legacyJsonItemKey(raw json.RawMessage) string {
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.UseNumber()
	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		return string(raw)
	}
	id, _ := m["id"].(string)
	return id
}
