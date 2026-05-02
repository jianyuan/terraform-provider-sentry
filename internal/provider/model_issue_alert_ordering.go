package provider

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
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

func issueAlertConditionModelKey(m IssueAlertConditionModel) string {
	switch {
	case m.FirstSeenEvent != nil:
		return "first_seen_event"
	case m.RegressionEvent != nil:
		return "regression_event"
	case m.ReappearedEvent != nil:
		return "reappeared_event"
	case m.NewHighPriorityIssue != nil:
		return "new_high_priority_issue"
	case m.ExistingHighPriorityIssue != nil:
		return "existing_high_priority_issue"
	case m.EventFrequency != nil:
		f := m.EventFrequency
		return fmt.Sprintf("event_frequency\x00%s\x00%s\x00%d\x00%s",
			f.ComparisonType.ValueString(), f.Interval.ValueString(),
			f.Value.ValueInt64(), f.ComparisonInterval.ValueString())
	case m.EventUniqueUserFrequency != nil:
		f := m.EventUniqueUserFrequency
		return fmt.Sprintf("event_unique_user_frequency\x00%s\x00%s\x00%d\x00%s",
			f.ComparisonType.ValueString(), f.Interval.ValueString(),
			f.Value.ValueInt64(), f.ComparisonInterval.ValueString())
	case m.EventFrequencyPercent != nil:
		f := m.EventFrequencyPercent
		return fmt.Sprintf("event_frequency_percent\x00%s\x00%s\x00%g\x00%s",
			f.ComparisonType.ValueString(), f.Interval.ValueString(),
			f.Value.ValueFloat64(), f.ComparisonInterval.ValueString())
	default:
		return ""
	}
}

func issueAlertFilterModelKey(m IssueAlertFilterModel) string {
	switch {
	case m.AgeComparison != nil:
		f := m.AgeComparison
		return fmt.Sprintf("age_comparison\x00%s\x00%d\x00%s",
			f.ComparisonType.ValueString(), f.Value.ValueInt64(), f.Time.ValueString())
	case m.IssueOccurrences != nil:
		return fmt.Sprintf("issue_occurrences\x00%d", m.IssueOccurrences.Value.ValueInt64())
	case m.AssignedTo != nil:
		f := m.AssignedTo
		return fmt.Sprintf("assigned_to\x00%s\x00%s",
			f.TargetType.ValueString(), f.TargetIdentifier.ValueString())
	case m.LatestAdoptedRelease != nil:
		f := m.LatestAdoptedRelease
		return fmt.Sprintf("latest_adopted_release\x00%s\x00%s\x00%s",
			f.OldestOrNewest.ValueString(), f.OlderOrNewer.ValueString(), f.Environment.ValueString())
	case m.LatestRelease != nil:
		return "latest_release"
	case m.IssueCategory != nil:
		return fmt.Sprintf("issue_category\x00%s", m.IssueCategory.Value.ValueString())
	case m.EventAttribute != nil:
		f := m.EventAttribute
		return fmt.Sprintf("event_attribute\x00%s\x00%s\x00%s",
			f.Attribute.ValueString(), f.Match.ValueString(), f.Value.ValueString())
	case m.TaggedEvent != nil:
		f := m.TaggedEvent
		return fmt.Sprintf("tagged_event\x00%s\x00%s\x00%s",
			f.Key.ValueString(), f.Match.ValueString(), f.Value.ValueString())
	case m.Level != nil:
		f := m.Level
		return fmt.Sprintf("level\x00%s\x00%s", f.Match.ValueString(), f.Level.ValueString())
	default:
		return ""
	}
}

func issueAlertActionModelKey(m IssueAlertActionModel) string {
	switch {
	case m.NotifyEmail != nil:
		f := m.NotifyEmail
		return fmt.Sprintf("notify_email\x00%s\x00%s\x00%s",
			f.TargetType.ValueString(), f.TargetIdentifier.ValueString(), f.FallthroughType.ValueString())
	case m.NotifyEvent != nil:
		return "notify_event"
	case m.NotifyEventService != nil:
		return fmt.Sprintf("notify_event_service\x00%s", m.NotifyEventService.Service.ValueString())
	case m.NotifyEventSentryApp != nil:
		return fmt.Sprintf("notify_event_sentry_app\x00%s", m.NotifyEventSentryApp.SentryAppInstallationUuid.ValueString())
	case m.OpsgenieNotifyTeam != nil:
		f := m.OpsgenieNotifyTeam
		return fmt.Sprintf("opsgenie_notify_team\x00%s\x00%s\x00%s",
			f.Account.ValueString(), f.Team.ValueString(), f.Priority.ValueString())
	case m.PagerDutyNotifyService != nil:
		f := m.PagerDutyNotifyService
		return fmt.Sprintf("pagerduty_notify_service\x00%s\x00%s\x00%s",
			f.Account.ValueString(), f.Service.ValueString(), f.Severity.ValueString())
	case m.SlackNotifyService != nil:
		f := m.SlackNotifyService
		return fmt.Sprintf("slack_notify_service\x00%s\x00%s\x00%s",
			f.Workspace.ValueString(), f.Channel.ValueString(), stringSetKey(f.Tags))
	case m.MsTeamsNotifyService != nil:
		f := m.MsTeamsNotifyService
		return fmt.Sprintf("msteams_notify_service\x00%s\x00%s",
			f.Team.ValueString(), f.Channel.ValueString())
	case m.DiscordNotifyService != nil:
		f := m.DiscordNotifyService
		return fmt.Sprintf("discord_notify_service\x00%s\x00%s\x00%s",
			f.Server.ValueString(), f.ChannelId.ValueString(), stringSetKey(f.Tags))
	case m.JiraCreateTicket != nil:
		f := m.JiraCreateTicket
		return fmt.Sprintf("jira_create_ticket\x00%s\x00%s\x00%s",
			f.Integration.ValueString(), f.Project.ValueString(), f.IssueType.ValueString())
	case m.JiraServerCreateTicket != nil:
		f := m.JiraServerCreateTicket
		return fmt.Sprintf("jira_server_create_ticket\x00%s\x00%s\x00%s",
			f.Integration.ValueString(), f.Project.ValueString(), f.IssueType.ValueString())
	case m.GitHubCreateTicket != nil:
		f := m.GitHubCreateTicket
		return fmt.Sprintf("github_create_ticket\x00%s\x00%s\x00%s\x00%s",
			f.Integration.ValueString(), f.Repo.ValueString(),
			f.Assignee.ValueString(), typesSetKey(f.Labels))
	case m.GitHubEnterpriseCreateTicket != nil:
		f := m.GitHubEnterpriseCreateTicket
		return fmt.Sprintf("github_enterprise_create_ticket\x00%s\x00%s\x00%s\x00%s",
			f.Integration.ValueString(), f.Repo.ValueString(),
			f.Assignee.ValueString(), typesSetKey(f.Labels))
	case m.AzureDevopsCreateTicket != nil:
		f := m.AzureDevopsCreateTicket
		return fmt.Sprintf("azure_devops_create_ticket\x00%s\x00%s\x00%s",
			f.Integration.ValueString(), f.Project.ValueString(), f.WorkItemType.ValueString())
	default:
		return ""
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

// typesSetKey produces a stable string key from a types.Set.
func typesSetKey(s types.Set) string {
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
