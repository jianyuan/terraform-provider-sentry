package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/kr/pretty"
)

type IssueAlertConditionFirstSeenEventResourceModel struct {
	Name types.String `tfsdk:"name"`
}

func (m IssueAlertConditionFirstSeenEventResourceModel) SentryId() string {
	return "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
}

func (m IssueAlertConditionFirstSeenEventResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id": m.SentryId(),
	}
}

type IssueAlertConditionRegressionEventResourceModel struct {
	Name types.String `tfsdk:"name"`
}

func (m IssueAlertConditionRegressionEventResourceModel) SentryId() string {
	return "sentry.rules.conditions.regression_event.RegressionEventCondition"
}

func (m IssueAlertConditionRegressionEventResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id": m.SentryId(),
	}
}

type IssueAlertConditionEventFrequencyResourceModel struct {
	Name     types.String `tfsdk:"name"`
	Value    types.Int64  `tfsdk:"value"`
	Interval types.String `tfsdk:"interval"`
}

func (m IssueAlertConditionEventFrequencyResourceModel) SentryId() string {
	return "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
}

func (m IssueAlertConditionEventFrequencyResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":       m.SentryId(),
		"value":    m.Value.ValueInt64(),
		"interval": m.Interval.ValueString(),
	}
}

type IssueAlertConditionEventUniqueUserFrequencyResourceModel struct {
	Name     types.String `tfsdk:"name"`
	Value    types.Int64  `tfsdk:"value"`
	Interval types.String `tfsdk:"interval"`
}

func (m IssueAlertConditionEventUniqueUserFrequencyResourceModel) SentryId() string {
	return "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"
}

func (m IssueAlertConditionEventUniqueUserFrequencyResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":       m.SentryId(),
		"value":    m.Value.ValueInt64(),
		"interval": m.Interval.ValueString(),
	}
}

type IssueAlertConditionEventFrequencyPercentResourceModel struct {
	Name     types.String  `tfsdk:"name"`
	Value    types.Float64 `tfsdk:"value"`
	Interval types.String  `tfsdk:"interval"`
}

func (m IssueAlertConditionEventFrequencyPercentResourceModel) SentryId() string {
	return "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"
}

func (m IssueAlertConditionEventFrequencyPercentResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":       m.SentryId(),
		"value":    m.Value.ValueFloat64(),
		"interval": m.Interval.ValueString(),
	}
}

type IssueAlertConditionResourceModel struct {
	FirstSeenEvent           []IssueAlertConditionFirstSeenEventResourceModel           `tfsdk:"first_seen_event"`
	RegressionEvent          []IssueAlertConditionRegressionEventResourceModel          `tfsdk:"regression_event"`
	EventFrequency           []IssueAlertConditionEventFrequencyResourceModel           `tfsdk:"event_frequency"`
	EventUniqueUserFrequency []IssueAlertConditionEventUniqueUserFrequencyResourceModel `tfsdk:"event_unique_user_frequency"`
	EventFrequencyPercent    []IssueAlertConditionEventFrequencyPercentResourceModel    `tfsdk:"event_frequency_percent"`
}

func (m IssueAlertConditionResourceModel) ToSentry() []map[string]interface{} {
	var conditions []map[string]interface{}

	for _, condition := range m.FirstSeenEvent {
		conditions = append(conditions, condition.ToSentry())
	}
	for _, condition := range m.RegressionEvent {
		conditions = append(conditions, condition.ToSentry())
	}
	for _, condition := range m.EventFrequency {
		conditions = append(conditions, condition.ToSentry())
	}
	for _, condition := range m.EventUniqueUserFrequency {
		conditions = append(conditions, condition.ToSentry())
	}
	for _, condition := range m.EventFrequencyPercent {
		conditions = append(conditions, condition.ToSentry())
	}

	return conditions
}

type IssueAlertFilterAgeComparisonResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ComparisonType types.String `tfsdk:"comparison_type"`
	Value          types.Int64  `tfsdk:"value"`
	Time           types.String `tfsdk:"time"`
}

type IssueAlertFilterResourceModel struct {
	AgeComparison []IssueAlertFilterAgeComparisonResourceModel `tfsdk:"age_comparison"`
}

type IssueAlertActionNotifyEmailResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	TargetType       types.String `tfsdk:"target_type"`
	TargetIdentifier types.String `tfsdk:"target_identifier"`
	FallthroughType  types.String `tfsdk:"fallthrough_type"`
}

func (m IssueAlertActionNotifyEmailResourceModel) SentryId() string {
	return "sentry.mail.actions.NotifyEmailAction"
}

func (m IssueAlertActionNotifyEmailResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":               m.SentryId(),
		"uuid":             m.Id.ValueStringPointer(),
		"targetType":       m.TargetType.ValueStringPointer(),
		"targetIdentifier": m.TargetIdentifier.ValueStringPointer(),
		"fallthroughType":  m.FallthroughType.ValueStringPointer(),
	}
}

type IssueAlertActionSlackNotifyServiceResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Workspace types.String `tfsdk:"workspace"`
	Channel   types.String `tfsdk:"channel"`
	ChannelId types.String `tfsdk:"channel_id"`
	Tags      types.String `tfsdk:"tags"`
	Notes     types.String `tfsdk:"notes"`
}

func (m IssueAlertActionSlackNotifyServiceResourceModel) SentryId() string {
	return "sentry.integrations.slack.notify_action.SlackNotifyServiceAction"
}

func (m IssueAlertActionSlackNotifyServiceResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":         m.SentryId(),
		"uuid":       m.Id.ValueStringPointer(),
		"workspace":  m.Workspace.ValueStringPointer(),
		"channel":    m.Channel.ValueStringPointer(),
		"channel_id": m.ChannelId.ValueStringPointer(),
		"tags":       m.Tags.ValueStringPointer(),
		"notes":      m.Notes.ValueStringPointer(),
	}
}

type IssueAlertActionMsTeamsNotifyServiceResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Team    types.String `tfsdk:"team"`
	Channel types.String `tfsdk:"channel"`
}

func (m IssueAlertActionMsTeamsNotifyServiceResourceModel) SentryId() string {
	return "sentry.integrations.msteams.notify_action.MsTeamsNotifyServiceAction"
}

func (m IssueAlertActionMsTeamsNotifyServiceResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":      m.SentryId(),
		"uuid":    m.Id.ValueStringPointer(),
		"team":    m.Team.ValueStringPointer(),
		"channel": m.Channel.ValueStringPointer(),
	}
}

type IssueAlertActionDiscordNotifyServiceResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Server    types.String `tfsdk:"server"`
	ChannelId types.String `tfsdk:"channel_id"`
	Tags      types.String `tfsdk:"tags"`
}

func (m IssueAlertActionDiscordNotifyServiceResourceModel) SentryId() string {
	return "sentry.integrations.discord.notify_action.DiscordNotifyServiceAction"
}

func (m IssueAlertActionDiscordNotifyServiceResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":         m.SentryId(),
		"uuid":       m.Id.ValueStringPointer(),
		"server":     m.Server.ValueStringPointer(),
		"channel_id": m.ChannelId.ValueStringPointer(),
		"tags":       m.Tags.ValueStringPointer(),
	}
}

type IssueAlertActionJiraCreateTicketResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Project     types.String `tfsdk:"project"`
	IssueType   types.String `tfsdk:"issue_type"`
}

func (m IssueAlertActionJiraCreateTicketResourceModel) SentryId() string {
	return "sentry.integrations.jira.notify_action.JiraCreateTicketAction"
}

func (m IssueAlertActionJiraCreateTicketResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":          m.SentryId(),
		"uuid":        m.Id.ValueStringPointer(),
		"integration": m.Integration.ValueStringPointer(),
		"project":     m.Project.ValueStringPointer(),
		"issuetype":   m.IssueType.ValueStringPointer(),
	}
}

type IssueAlertActionJiraServerCreateTicketResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Project     types.String `tfsdk:"project"`
	IssueType   types.String `tfsdk:"issue_type"`
}

func (m IssueAlertActionJiraServerCreateTicketResourceModel) SentryId() string {
	return "sentry.integrations.jira_server.notify_action.JiraServerCreateTicketAction"
}

func (m IssueAlertActionJiraServerCreateTicketResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":          m.SentryId(),
		"uuid":        m.Id.ValueStringPointer(),
		"integration": m.Integration.ValueStringPointer(),
		"project":     m.Project.ValueStringPointer(),
		"issuetype":   m.IssueType.ValueStringPointer(),
	}
}

type IssueAlertActionGitHubCreateTicketResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Repo        types.String `tfsdk:"repo"`
	Title       types.String `tfsdk:"title"`
	Body        types.String `tfsdk:"body"`
	Assignee    types.String `tfsdk:"assignee"`
	Labels      types.Set    `tfsdk:"labels"`
}

func (m IssueAlertActionGitHubCreateTicketResourceModel) SentryId() string {
	return "sentry.integrations.github.notify_action.GitHubCreateTicketAction"
}

func (m IssueAlertActionGitHubCreateTicketResourceModel) ToSentry() map[string]interface{} {
	labels := []string{}
	if !m.Labels.IsNull() {
		m.Labels.ElementsAs(context.Background(), &labels, false)
	}

	return map[string]interface{}{
		"id":          m.SentryId(),
		"uuid":        m.Id.ValueStringPointer(),
		"integration": m.Integration.ValueStringPointer(),
		"repo":        m.Repo.ValueStringPointer(),
		"title":       m.Title.ValueStringPointer(),
		"body":        m.Body.ValueStringPointer(),
		"assignee":    m.Assignee.ValueStringPointer(),
		"labels":      labels,
		"dynamic_form_fields": []map[string]interface{}{
			{"ok": "ok"}, // Must be truthy
		},
	}
}

type IssueAlertActionGitHubEnterpriseCreateTicketResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Integration types.String `tfsdk:"integration"`
	Repo        types.String `tfsdk:"repo"`
	Title       types.String `tfsdk:"title"`
	Body        types.String `tfsdk:"body"`
	Assignee    types.String `tfsdk:"assignee"`
	Labels      types.Set    `tfsdk:"labels"`
}

func (m IssueAlertActionGitHubEnterpriseCreateTicketResourceModel) SentryId() string {
	return "sentry.integrations.github_enterprise.notify_action.GitHubEnterpriseCreateTicketAction"
}

func (m IssueAlertActionGitHubEnterpriseCreateTicketResourceModel) ToSentry() map[string]interface{} {
	labels := []string{}
	if !m.Labels.IsNull() {
		m.Labels.ElementsAs(context.Background(), &labels, false)
	}

	return map[string]interface{}{
		"id":          m.SentryId(),
		"uuid":        m.Id.ValueStringPointer(),
		"integration": m.Integration.ValueStringPointer(),
		"repo":        m.Repo.ValueStringPointer(),
		"title":       m.Title.ValueStringPointer(),
		"body":        m.Body.ValueStringPointer(),
		"assignee":    m.Assignee.ValueStringPointer(),
		"labels":      labels,
		"dynamic_form_fields": []map[string]interface{}{
			{"ok": "ok"}, // Must be truthy
		},
	}
}

type IssueAlertActionAzureDevopsCreateTicketResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Integration  types.String `tfsdk:"integration"`
	Project      types.String `tfsdk:"project"`
	WorkItemType types.String `tfsdk:"work_item_type"`
}

func (m IssueAlertActionAzureDevopsCreateTicketResourceModel) SentryId() string {
	return "sentry.integrations.vsts.notify_action.AzureDevopsCreateTicketAction"
}

func (m IssueAlertActionAzureDevopsCreateTicketResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":             m.SentryId(),
		"uuid":           m.Id.ValueStringPointer(),
		"integration":    m.Integration.ValueStringPointer(),
		"project":        m.Project.ValueStringPointer(),
		"work_item_type": m.WorkItemType.ValueStringPointer(),
		"dynamic_form_fields": []map[string]interface{}{
			{"ok": "ok"}, // Must be truthy
		},
	}
}

type IssueAlertActionPagerDutyNotifyServiceResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Account  types.String `tfsdk:"account"`
	Service  types.String `tfsdk:"service"`
	Severity types.String `tfsdk:"severity"`
}

func (m IssueAlertActionPagerDutyNotifyServiceResourceModel) SentryId() string {
	return "sentry.integrations.pagerduty.notify_action.PagerDutyNotifyServiceAction"
}

func (m IssueAlertActionPagerDutyNotifyServiceResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":       m.SentryId(),
		"uuid":     m.Id.ValueStringPointer(),
		"account":  m.Account.ValueStringPointer(),
		"service":  m.Service.ValueStringPointer(),
		"severity": m.Severity.ValueStringPointer(),
	}
}

type IssueAlertActionOpsgenieNotifyTeamResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Account  types.String `tfsdk:"account"`
	Team     types.String `tfsdk:"team"`
	Priority types.String `tfsdk:"priority"`
}

func (m IssueAlertActionOpsgenieNotifyTeamResourceModel) SentryId() string {
	return "sentry.integrations.opsgenie.notify_action.OpsgenieNotifyTeamAction"
}

func (m IssueAlertActionOpsgenieNotifyTeamResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":       m.SentryId(),
		"uuid":     m.Id.ValueStringPointer(),
		"account":  m.Account.ValueStringPointer(),
		"team":     m.Team.ValueStringPointer(),
		"priority": m.Priority.ValueStringPointer(),
	}
}

type IssueAlertActionNotifyEventResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (m IssueAlertActionNotifyEventResourceModel) SentryId() string {
	return "sentry.rules.actions.notify_event.NotifyEventAction"
}

func (m IssueAlertActionNotifyEventResourceModel) ToSentry() map[string]interface{} {
	return map[string]interface{}{
		"id":   m.SentryId(),
		"uuid": m.Id.ValueStringPointer(),
	}
}

type IssueAlertActionResourceModel struct {
	NotifyEmail                  []IssueAlertActionNotifyEmailResourceModel                  `tfsdk:"notify_email"`
	SlackNotifyService           []IssueAlertActionSlackNotifyServiceResourceModel           `tfsdk:"slack_notify_service"`
	MsTeamsNotifyService         []IssueAlertActionMsTeamsNotifyServiceResourceModel         `tfsdk:"ms_teams_notify_service"`
	DiscordNotifyService         []IssueAlertActionDiscordNotifyServiceResourceModel         `tfsdk:"discord_notify_service"`
	JiraCreateTicket             []IssueAlertActionJiraCreateTicketResourceModel             `tfsdk:"jira_create_ticket"`
	JiraServerCreateTicket       []IssueAlertActionJiraServerCreateTicketResourceModel       `tfsdk:"jira_server_create_ticket"`
	GitHubCreateTicket           []IssueAlertActionGitHubCreateTicketResourceModel           `tfsdk:"github_create_ticket"`
	GitHubEnterpriseCreateTicket []IssueAlertActionGitHubEnterpriseCreateTicketResourceModel `tfsdk:"github_enterprise_create_ticket"`
	AzureDevopsCreateTicket      []IssueAlertActionAzureDevopsCreateTicketResourceModel      `tfsdk:"azure_devops_create_ticket"`
	PagerDutyNotifyService       []IssueAlertActionPagerDutyNotifyServiceResourceModel       `tfsdk:"pagerduty_notify_service"`
	OpsgenieNotifyTeam           []IssueAlertActionOpsgenieNotifyTeamResourceModel           `tfsdk:"opsgenie_notify_team"`
	NotifyEvent                  []IssueAlertActionNotifyEventResourceModel                  `tfsdk:"notify_event"`
}

func (m IssueAlertActionResourceModel) ToSentry() []map[string]interface{} {
	var actions []map[string]interface{}

	for _, action := range m.NotifyEmail {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.SlackNotifyService {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.MsTeamsNotifyService {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.DiscordNotifyService {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.JiraCreateTicket {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.JiraServerCreateTicket {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.GitHubCreateTicket {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.GitHubEnterpriseCreateTicket {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.AzureDevopsCreateTicket {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.PagerDutyNotifyService {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.OpsgenieNotifyTeam {
		actions = append(actions, action.ToSentry())
	}
	for _, action := range m.NotifyEvent {
		actions = append(actions, action.ToSentry())
	}

	return actions
}

type IssueAlertResourceModel struct {
	Id           types.String                       `tfsdk:"id"`
	Organization types.String                       `tfsdk:"organization"`
	Project      types.String                       `tfsdk:"project"`
	Name         types.String                       `tfsdk:"name"`
	Conditions   sentrytypes.LossyJson              `tfsdk:"conditions"`
	Filters      sentrytypes.LossyJson              `tfsdk:"filters"`
	Actions      sentrytypes.LossyJson              `tfsdk:"actions"`
	ActionMatch  types.String                       `tfsdk:"action_match"`
	FilterMatch  types.String                       `tfsdk:"filter_match"`
	Frequency    types.Int64                        `tfsdk:"frequency"`
	Environment  types.String                       `tfsdk:"environment"`
	Owner        types.String                       `tfsdk:"owner"`
	Condition    []IssueAlertConditionResourceModel `tfsdk:"condition"`
	Filter       []IssueAlertFilterResourceModel    `tfsdk:"filter"`
	Action       []IssueAlertActionResourceModel    `tfsdk:"action"`
}

func (m *IssueAlertResourceModel) Fill(organization string, alert sentry.IssueAlert) error {
	m.Id = types.StringPointerValue(alert.ID)
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(alert.Projects[0])
	m.Name = types.StringPointerValue(alert.Name)
	m.ActionMatch = types.StringPointerValue(alert.ActionMatch)
	m.FilterMatch = types.StringPointerValue(alert.FilterMatch)
	m.Owner = types.StringPointerValue(alert.Owner)

	m.Conditions = sentrytypes.NewLossyJsonNull()
	m.Condition = []IssueAlertConditionResourceModel{}
	if !m.Conditions.IsNull() {
		m.Conditions = sentrytypes.NewLossyJsonValue("[]")
		if len(alert.Conditions) > 0 {
			if conditions, err := json.Marshal(alert.Conditions); err == nil {
				m.Conditions = sentrytypes.NewLossyJsonValue(string(conditions))
			} else {
				return err
			}
		}
	} else {
		for _, condition := range alert.Conditions {
			switch condition["id"] {
			case IssueAlertConditionFirstSeenEventResourceModel{}.SentryId():
				m.Condition = append(m.Condition, IssueAlertConditionResourceModel{
					FirstSeenEvent: []IssueAlertConditionFirstSeenEventResourceModel{
						{
							Name: types.StringValue(condition["name"].(string)),
						},
					},
				})
			case IssueAlertConditionRegressionEventResourceModel{}.SentryId():
				m.Condition = append(m.Condition, IssueAlertConditionResourceModel{
					RegressionEvent: []IssueAlertConditionRegressionEventResourceModel{
						{
							Name: types.StringValue(condition["name"].(string)),
						},
					},
				})
			case IssueAlertConditionEventFrequencyResourceModel{}.SentryId():
				m.Condition = append(m.Condition, IssueAlertConditionResourceModel{
					EventFrequency: []IssueAlertConditionEventFrequencyResourceModel{
						{
							Name:     types.StringValue(condition["name"].(string)),
							Value:    types.Int64Value(must.Get(condition["value"].(json.Number).Int64())),
							Interval: types.StringValue(condition["interval"].(string)),
						},
					},
				})
			case IssueAlertConditionEventUniqueUserFrequencyResourceModel{}.SentryId():
				m.Condition = append(m.Condition, IssueAlertConditionResourceModel{
					EventUniqueUserFrequency: []IssueAlertConditionEventUniqueUserFrequencyResourceModel{
						{
							Name:     types.StringValue(condition["name"].(string)),
							Value:    types.Int64Value(must.Get(condition["value"].(json.Number).Int64())),
							Interval: types.StringValue(condition["interval"].(string)),
						},
					},
				})
			case IssueAlertConditionEventFrequencyPercentResourceModel{}.SentryId():
				m.Condition = append(m.Condition, IssueAlertConditionResourceModel{
					EventFrequencyPercent: []IssueAlertConditionEventFrequencyPercentResourceModel{
						{
							Name:     types.StringValue(condition["name"].(string)),
							Value:    types.Float64Value(must.Get(condition["value"].(json.Number).Float64())),
							Interval: types.StringValue(condition["interval"].(string)),
						},
					},
				})
			default:
				return fmt.Errorf("unsupported condition: %s", condition["id"])
			}
		}
	}

	m.Filters = sentrytypes.NewLossyJsonNull()
	m.Filter = []IssueAlertFilterResourceModel{}
	if !m.Filters.IsNull() {
		if len(alert.Filters) > 0 {
			if filters, err := json.Marshal(alert.Filters); err == nil {
				m.Filters = sentrytypes.NewLossyJsonValue(string(filters))
			} else {
				return err
			}
		}
	} else {
		// TODO
	}

	m.Actions = sentrytypes.NewLossyJsonNull()
	m.Action = []IssueAlertActionResourceModel{}
	if !m.Actions.IsNull() {
		if len(alert.Actions) > 0 {
			if actions, err := json.Marshal(alert.Actions); err == nil && len(actions) > 0 {
				m.Actions = sentrytypes.NewLossyJsonValue(string(actions))
			} else {
				return err
			}
		}
	} else {
		for _, action := range alert.Actions {
			switch action["id"] {
			case IssueAlertActionNotifyEmailResourceModel{}.SentryId():
				notifyEmailAction := IssueAlertActionNotifyEmailResourceModel{
					Id:              types.StringValue(action["uuid"].(string)),
					Name:            types.StringValue(action["name"].(string)),
					TargetType:      types.StringValue(action["targetType"].(string)),
					FallthroughType: types.StringValue(action["fallthroughType"].(string)),
				}

				switch value := action["targetIdentifier"].(type) {
				case string:
					if value == "" {
						notifyEmailAction.TargetIdentifier = types.StringNull()
					} else {
						notifyEmailAction.TargetIdentifier = types.StringValue(value)
					}
				case json.Number:
					notifyEmailAction.TargetIdentifier = types.StringValue(value.String())
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					NotifyEmail: []IssueAlertActionNotifyEmailResourceModel{notifyEmailAction},
				})
			case IssueAlertActionSlackNotifyServiceResourceModel{}.SentryId():
				slackNotifyServiceAction := IssueAlertActionSlackNotifyServiceResourceModel{
					Id:        types.StringValue(action["uuid"].(string)),
					Name:      types.StringValue(action["name"].(string)),
					Workspace: types.StringValue(action["workspace"].(string)),
					Channel:   types.StringValue(action["channel"].(string)),
				}
				if value, ok := action["channel_id"].(string); ok {
					slackNotifyServiceAction.ChannelId = types.StringValue(value)
				} else {
					slackNotifyServiceAction.ChannelId = types.StringNull()
				}

				if value, ok := action["tags"].(string); ok {
					slackNotifyServiceAction.Tags = types.StringValue(value)
				} else {
					slackNotifyServiceAction.Tags = types.StringNull()
				}

				if value, ok := action["notes"].(string); ok {
					slackNotifyServiceAction.Notes = types.StringValue(value)
				} else {
					slackNotifyServiceAction.Notes = types.StringNull()
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					SlackNotifyService: []IssueAlertActionSlackNotifyServiceResourceModel{slackNotifyServiceAction},
				})
			case IssueAlertActionDiscordNotifyServiceResourceModel{}.SentryId():
				discordNotifyServiceAction := IssueAlertActionDiscordNotifyServiceResourceModel{
					Id:        types.StringValue(action["uuid"].(string)),
					Name:      types.StringValue(action["name"].(string)),
					Server:    types.StringValue(action["server"].(string)),
					ChannelId: types.StringValue(action["channel_id"].(string)),
				}
				if value, ok := action["tags"].(string); ok {
					discordNotifyServiceAction.Tags = types.StringValue(value)
				} else {
					discordNotifyServiceAction.Tags = types.StringNull()
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					DiscordNotifyService: []IssueAlertActionDiscordNotifyServiceResourceModel{discordNotifyServiceAction},
				})
			case IssueAlertActionJiraCreateTicketResourceModel{}.SentryId():
				jiraCreateTicketAction := IssueAlertActionJiraCreateTicketResourceModel{
					Id:          types.StringValue(action["uuid"].(string)),
					Name:        types.StringValue(action["name"].(string)),
					Integration: types.StringValue(action["integration"].(string)),
					Project:     types.StringValue(action["project"].(string)),
					IssueType:   types.StringValue(action["issuetype"].(string)),
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					JiraCreateTicket: []IssueAlertActionJiraCreateTicketResourceModel{jiraCreateTicketAction},
				})
			case IssueAlertActionJiraServerCreateTicketResourceModel{}.SentryId():
				jiraServerCreateTicketAction := IssueAlertActionJiraServerCreateTicketResourceModel{
					Id:          types.StringValue(action["uuid"].(string)),
					Name:        types.StringValue(action["name"].(string)),
					Integration: types.StringValue(action["integration"].(string)),
					Project:     types.StringValue(action["project"].(string)),
					IssueType:   types.StringValue(action["issuetype"].(string)),
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					JiraServerCreateTicket: []IssueAlertActionJiraServerCreateTicketResourceModel{jiraServerCreateTicketAction},
				})
			case IssueAlertActionMsTeamsNotifyServiceResourceModel{}.SentryId():
				msTeamsNotifyServiceAction := IssueAlertActionMsTeamsNotifyServiceResourceModel{
					Id:      types.StringValue(action["uuid"].(string)),
					Team:    types.StringValue(action["team"].(string)),
					Channel: types.StringValue(action["channel"].(string)),
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					MsTeamsNotifyService: []IssueAlertActionMsTeamsNotifyServiceResourceModel{msTeamsNotifyServiceAction},
				})
			case IssueAlertActionGitHubCreateTicketResourceModel{}.SentryId():
				gitHubCreateTicketAction := IssueAlertActionGitHubCreateTicketResourceModel{
					Id:          types.StringValue(action["uuid"].(string)),
					Name:        types.StringValue(action["name"].(string)),
					Integration: types.StringValue(action["integration"].(string)),
					Repo:        types.StringValue(action["repo"].(string)),
					Title:       types.StringValue(action["title"].(string)),
				}

				if value, ok := action["body"].(string); ok {
					gitHubCreateTicketAction.Body = types.StringValue(value)
				} else {
					gitHubCreateTicketAction.Body = types.StringNull()
				}

				if value, ok := action["assignee"].(string); ok {
					gitHubCreateTicketAction.Assignee = types.StringValue(value)
				} else {
					gitHubCreateTicketAction.Assignee = types.StringNull()
				}

				if value, ok := action["labels"].([]interface{}); ok {
					labelElements := make([]attr.Value, 0, len(value))
					for _, element := range value {
						labelElements = append(labelElements, types.StringValue(element.(string)))
					}
					gitHubCreateTicketAction.Labels = types.SetValueMust(types.StringType, labelElements)
				} else {
					gitHubCreateTicketAction.Labels = types.SetNull(types.StringType)
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					GitHubCreateTicket: []IssueAlertActionGitHubCreateTicketResourceModel{gitHubCreateTicketAction},
				})
			case IssueAlertActionGitHubEnterpriseCreateTicketResourceModel{}.SentryId():
				gitHubEnterpriseCreateTicketAction := IssueAlertActionGitHubEnterpriseCreateTicketResourceModel{
					Id:          types.StringValue(action["uuid"].(string)),
					Name:        types.StringValue(action["name"].(string)),
					Integration: types.StringValue(action["integration"].(string)),
					Repo:        types.StringValue(action["repo"].(string)),
					Title:       types.StringValue(action["title"].(string)),
				}

				if value, ok := action["body"].(string); ok {
					gitHubEnterpriseCreateTicketAction.Body = types.StringValue(value)
				} else {
					gitHubEnterpriseCreateTicketAction.Body = types.StringNull()
				}

				if value, ok := action["assignee"].(string); ok {
					gitHubEnterpriseCreateTicketAction.Assignee = types.StringValue(value)
				} else {
					gitHubEnterpriseCreateTicketAction.Assignee = types.StringNull()
				}

				if value, ok := action["labels"].([]interface{}); ok {
					labelElements := make([]attr.Value, 0, len(value))
					for _, element := range value {
						labelElements = append(labelElements, types.StringValue(element.(string)))
					}
					gitHubEnterpriseCreateTicketAction.Labels = types.SetValueMust(types.StringType, labelElements)
				} else {
					gitHubEnterpriseCreateTicketAction.Labels = types.SetNull(types.StringType)
				}

				m.Action = append(m.Action, IssueAlertActionResourceModel{
					GitHubEnterpriseCreateTicket: []IssueAlertActionGitHubEnterpriseCreateTicketResourceModel{gitHubEnterpriseCreateTicketAction},
				})
			case IssueAlertActionAzureDevopsCreateTicketResourceModel{}.SentryId():
				m.Action = append(m.Action, IssueAlertActionResourceModel{
					AzureDevopsCreateTicket: []IssueAlertActionAzureDevopsCreateTicketResourceModel{
						{
							Id:           types.StringValue(action["uuid"].(string)),
							Name:         types.StringValue(action["name"].(string)),
							Integration:  types.StringValue(action["integration"].(string)),
							Project:      types.StringValue(action["project"].(string)),
							WorkItemType: types.StringValue(action["work_item_type"].(string)),
						},
					},
				})
			case IssueAlertActionPagerDutyNotifyServiceResourceModel{}.SentryId():
				m.Action = append(m.Action, IssueAlertActionResourceModel{
					PagerDutyNotifyService: []IssueAlertActionPagerDutyNotifyServiceResourceModel{
						{
							Id:       types.StringValue(action["uuid"].(string)),
							Name:     types.StringValue(action["name"].(string)),
							Account:  types.StringValue(action["account"].(string)),
							Service:  types.StringValue(action["service"].(string)),
							Severity: types.StringValue(action["severity"].(string)),
						},
					},
				})
			case IssueAlertActionOpsgenieNotifyTeamResourceModel{}.SentryId():
				m.Action = append(m.Action, IssueAlertActionResourceModel{
					OpsgenieNotifyTeam: []IssueAlertActionOpsgenieNotifyTeamResourceModel{
						{
							Id:       types.StringValue(action["uuid"].(string)),
							Name:     types.StringValue(action["name"].(string)),
							Account:  types.StringValue(action["account"].(string)),
							Team:     types.StringValue(action["team"].(string)),
							Priority: types.StringValue(action["priority"].(string)),
						},
					},
				})
			case IssueAlertActionNotifyEventResourceModel{}.SentryId():
				m.Action = append(m.Action, IssueAlertActionResourceModel{
					NotifyEvent: []IssueAlertActionNotifyEventResourceModel{
						{
							Id:   types.StringValue(action["uuid"].(string)),
							Name: types.StringValue(action["name"].(string)),
						},
					},
				})
			default:
				return fmt.Errorf("unsupported action: %s", action["id"])
			}
		}
	}

	frequency, err := alert.Frequency.Int64()
	if err != nil {
		return err
	}
	m.Frequency = types.Int64Value(frequency)

	m.Environment = types.StringPointerValue(alert.Environment)
	m.Owner = types.StringPointerValue(alert.Owner)

	return nil
}

var _ resource.Resource = &IssueAlertResource{}
var _ resource.ResourceWithConfigure = &IssueAlertResource{}
var _ resource.ResourceWithImportState = &IssueAlertResource{}
var _ resource.ResourceWithValidateConfig = &IssueAlertResource{}
var _ resource.ResourceWithUpgradeState = &IssueAlertResource{}

func NewIssueAlertResource() resource.Resource {
	return &IssueAlertResource{}
}

type IssueAlertResource struct {
	baseResource
}

func (r *IssueAlertResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_alert"
}

func (r *IssueAlertResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	intervalStringAttribute := schema.StringAttribute{
		MarkdownDescription: "Valid values are `1m`, `5m`, `15m`, `1h`, `1d`, `1w` and `30d` (`m` for minutes, `h` for hours, `d` for days, and `w` for weeks).",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("1m", "5m", "15m", "1h", "1d", "1w", "30d"),
		},
	}
	idStringAttribute := schema.StringAttribute{
		Computed: true,
	}
	nameStringAttribute := schema.StringAttribute{
		Computed: true,
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: `Create an Issue Alert Rule for a Project. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/) for more information.

Please note the following changes since v0.12.0:
- The attributes ` + "`conditions`" + `, ` + "`filters`" + `, and ` + "`actions`" + ` are in JSON string format. The types must match the Sentry API, otherwise Terraform will incorrectly detect a drift. Use ` + "`parseint(\"string\", 10)`" + ` to convert a string to an integer. Avoid using ` + "`jsonencode()`" + ` as it is unable to distinguish between an integer and a float.
- The attribute ` + "`internal_id`" + ` has been removed. Use ` + "`id`" + ` instead.
- The attribute ` + "`id`" + ` is now the ID of the issue alert. Previously, it was a combination of the organization, project, and issue alert ID.
		`,

		Version: 2,

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
			"name": schema.StringAttribute{
				MarkdownDescription: "The issue alert name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "**Deprecated** in favor of `condition`. A list of triggers that determine when the rule fires. In JSON string format.",
				DeprecationMessage:  "Use `condition` instead.",
				Optional:            true,
				CustomType: sentrytypes.LossyJsonType{
					IgnoreKeys: []string{"name"},
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("condition")),
				},
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "**Deprecated** in favor of `filter`. A list of filters that determine if a rule fires after the necessary conditions have been met. In JSON string format.",
				DeprecationMessage:  "Use `filter` instead.",
				Optional:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"actions": schema.StringAttribute{
				MarkdownDescription: "**Deprecated** in favor of `action`. A list of actions that take place when all required conditions and filters for the rule are met. In JSON string format.",
				DeprecationMessage:  "Use `action` instead.",
				Optional:            true,
				CustomType:          sentrytypes.LossyJsonType{},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("action")),
				},
			},
			"action_match": schema.StringAttribute{
				MarkdownDescription: "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "any"),
				},
			},
			"filter_match": schema.StringAttribute{
				MarkdownDescription: "A string determining which filters need to be true before any actions take place. Required when a value is provided for `filters`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "any", "none"),
				},
			},
			"frequency": schema.Int64Attribute{
				MarkdownDescription: "Perform actions at most once every `X` minutes for this issue.",
				Required:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Perform issue alert in a specific environment.",
				Optional:            true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The ID of the team or user that owns the rule.",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"condition": schema.ListNestedBlock{
				MarkdownDescription: "A list of triggers that determine when the rule fires.",
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("conditions")),
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"first_seen_event": schema.ListNestedBlock{
							MarkdownDescription: "A new issue is created.",
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": nameStringAttribute,
								},
							},
						},
						"regression_event": schema.ListNestedBlock{
							MarkdownDescription: "The issue changes state from resolved to unresolved.",
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": nameStringAttribute,
								},
							},
						},
						"event_frequency": schema.ListNestedBlock{
							MarkdownDescription: "The issue is seen more than `value` times in `interval`.",
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": nameStringAttribute,
									"value": schema.Int64Attribute{
										Required: true,
									},
									"interval": intervalStringAttribute,
								},
							},
						},
						"event_unique_user_frequency": schema.ListNestedBlock{
							MarkdownDescription: "The issue is seen by more than `value` users in `interval`.",
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": nameStringAttribute,
									"value": schema.Int64Attribute{
										Required: true,
									},
									"interval": intervalStringAttribute,
								},
							},
						},
						"event_frequency_percent": schema.ListNestedBlock{
							MarkdownDescription: "The issue affects more than `value` percent of sessions in `interval`.",
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": nameStringAttribute,
									"value": schema.Float64Attribute{
										Required: true,
									},
									"interval": schema.StringAttribute{
										MarkdownDescription: "Valid values are `5m`, `10m`, `30m`, and `1h` (`m` for minutes, `h` for hours).",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("5m", "10m", "30m", "1h"),
										},
									},
								},
							},
						},
					},
				},
			},
			"filter": schema.ListNestedBlock{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met.",
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("filters")),
				},
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"age_comparison": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"comparison_type": schema.StringAttribute{
										MarkdownDescription: "One of `older` or `newer`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("older", "newer"),
										},
									},
									"value": schema.Int64Attribute{
										Required: true,
									},
									"time": schema.StringAttribute{
										MarkdownDescription: "The unit of time. Valid values are `minute`, `hour`, `day`, and `week`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("minute", "hour", "day", "week"),
										},
									},
								},
							},
						},
					},
				},
			},
			"action": schema.ListNestedBlock{
				MarkdownDescription: "A list of actions that take place when all required conditions and filters for the rule are met.",
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("actions")),
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"notify_email": schema.ListNestedBlock{
							MarkdownDescription: "Send a notification to Suggested Assignees.",
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"target_type": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("IssueOwners", "Team", "Member"),
										},
									},
									"target_identifier": schema.StringAttribute{
										Optional: true,
									},
									"fallthrough_type": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("AllMembers", "ActiveMembers", "NoOne"),
										},
									},
								},
							},
						},
						"slack_notify_service": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"workspace": schema.StringAttribute{
										Required: true,
									},
									"channel": schema.StringAttribute{
										Required: true,
									},
									"channel_id": schema.StringAttribute{
										Computed: true,
									},
									"tags": schema.StringAttribute{
										Optional: true,
									},
									"notes": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
						"ms_teams_notify_service": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"team": schema.StringAttribute{
										Required: true,
									},
									"channel": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"discord_notify_service": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"server": schema.StringAttribute{
										Required: true,
									},
									"channel_id": schema.StringAttribute{
										Required: true,
									},
									"tags": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
						"jira_create_ticket": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"integration": schema.StringAttribute{
										Required: true,
									},
									"project": schema.StringAttribute{
										Required: true,
									},
									"issue_type": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"jira_server_create_ticket": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"integration": schema.StringAttribute{
										Required: true,
									},
									"project": schema.StringAttribute{
										Required: true,
									},
									"issue_type": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"github_create_ticket": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"integration": schema.StringAttribute{
										Required: true,
									},
									"repo": schema.StringAttribute{
										Required: true,
									},
									"title": schema.StringAttribute{
										Required: true,
									},
									"body": schema.StringAttribute{
										Optional: true,
									},
									"assignee": schema.StringAttribute{
										Optional: true,
									},
									"labels": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
						"github_enterprise_create_ticket": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"integration": schema.StringAttribute{
										Required: true,
									},
									"repo": schema.StringAttribute{
										Required: true,
									},
									"title": schema.StringAttribute{
										Required: true,
									},
									"body": schema.StringAttribute{
										Optional: true,
									},
									"assignee": schema.StringAttribute{
										Optional: true,
									},
									"labels": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
						"azure_devops_create_ticket": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"integration": schema.StringAttribute{
										Required: true,
									},
									"project": schema.StringAttribute{
										Required: true,
									},
									"work_item_type": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"pagerduty_notify_service": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"account": schema.StringAttribute{
										Required: true,
									},
									"service": schema.StringAttribute{
										Required: true,
									},
									"severity": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"opsgenie_notify_team": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
									"account": schema.StringAttribute{
										Required: true,
									},
									"team": schema.StringAttribute{
										Required: true,
									},
									"priority": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"notify_event": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id":   idStringAttribute,
									"name": nameStringAttribute,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *IssueAlertResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// TODO: Implement validation
}

func (r *IssueAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.IssueAlert{
		Name:        data.Name.ValueStringPointer(),
		ActionMatch: data.ActionMatch.ValueStringPointer(),
		FilterMatch: data.FilterMatch.ValueStringPointer(),
		Frequency:   sentry.JsonNumber(json.Number(data.Frequency.String())),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.String()},
	}

	if len(data.Condition) > 0 {
		for _, condition := range data.Condition {
			params.Conditions = append(params.Conditions, condition.ToSentry()...)
		}
	} else if !data.Conditions.IsNull() {
		resp.Diagnostics.Append(data.Conditions.Unmarshal(&params.Conditions)...)
	} else {
		// TODO
		// resp.Diagnostics.AddError("Missing required block", "The `condition` block is required.")
	}

	if len(data.Filter) > 0 {

	} else if !data.Filters.IsNull() {
		resp.Diagnostics.Append(data.Filters.Unmarshal(&params.Filters)...)
	}

	if len(data.Action) > 0 {
		for _, action := range data.Action {
			params.Actions = append(params.Actions, action.ToSentry()...)
		}
	} else if !data.Actions.IsNull() {
		resp.Diagnostics.Append(data.Actions.Unmarshal(&params.Actions)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	action, _, err := r.client.IssueAlerts.Create(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	}

	// TODO
	pretty.Println(*action)

	if err := data.Fill(data.Organization.ValueString(), *action); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, apiResp, err := r.client.IssueAlerts.Get(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("issue alert"))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *action); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.IssueAlert{
		Name:        data.Name.ValueStringPointer(),
		ActionMatch: data.ActionMatch.ValueStringPointer(),
		FilterMatch: data.FilterMatch.ValueStringPointer(),
		Frequency:   sentry.JsonNumber(json.Number(data.Frequency.String())),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.String()},
	}
	if !data.Conditions.IsNull() {
		resp.Diagnostics.Append(data.Conditions.Unmarshal(&params.Conditions)...)
	}
	if !data.Filters.IsNull() {
		resp.Diagnostics.Append(data.Filters.Unmarshal(&params.Filters)...)
	}
	if !data.Actions.IsNull() {
		resp.Diagnostics.Append(data.Actions.Unmarshal(&params.Actions)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	action, apiResp, err := r.client.IssueAlerts.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
		params,
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("issue alert"))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *action); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.IssueAlerts.Delete(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	}
}

func (r *IssueAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, actionId, err := splitThreePartID(req.ID, "organization", "project-slug", "alert-id")
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("project"), project,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), actionId,
	)...)
}

func (r *IssueAlertResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	type modelV0 struct {
		Id           types.String `tfsdk:"id"`
		Organization types.String `tfsdk:"organization"`
		Project      types.String `tfsdk:"project"`
		Name         types.String `tfsdk:"name"`
		Conditions   types.List   `tfsdk:"conditions"`
		Filters      types.List   `tfsdk:"filters"`
		Actions      types.List   `tfsdk:"actions"`
		ActionMatch  types.String `tfsdk:"action_match"`
		FilterMatch  types.String `tfsdk:"filter_match"`
		Frequency    types.Int64  `tfsdk:"frequency"`
		Environment  types.String `tfsdk:"environment"`
	}

	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				// No-op
			},
		},
		1: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"organization": schema.StringAttribute{
						Required: true,
					},
					"project": schema.StringAttribute{
						Required: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"conditions": schema.ListAttribute{
						ElementType: types.MapType{
							ElemType: types.StringType,
						},
						Required: true,
					},
					"filters": schema.ListAttribute{
						ElementType: types.MapType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
					"actions": schema.ListAttribute{
						ElementType: types.MapType{
							ElemType: types.StringType,
						},
						Required: true,
					},
					"action_match": schema.StringAttribute{
						Optional: true,
					},
					"filter_match": schema.StringAttribute{
						Optional: true,
					},
					"frequency": schema.Int64Attribute{
						Optional: true,
					},
					"environment": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData modelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
				if resp.Diagnostics.HasError() {
					return
				}

				organization, project, actionId, err := splitThreePartID(priorStateData.Id.ValueString(), "organization", "project-slug", "alert-id")
				if err != nil {
					resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
					return
				}

				upgradedStateData := IssueAlertResourceModel{
					Id:           types.StringValue(actionId),
					Organization: types.StringValue(organization),
					Project:      types.StringValue(project),
					Name:         priorStateData.Name,
					ActionMatch:  priorStateData.ActionMatch,
					FilterMatch:  priorStateData.FilterMatch,
					Frequency:    priorStateData.Frequency,
					Environment:  priorStateData.Environment,
				}

				upgradedStateData.Conditions = sentrytypes.NewLossyJsonNull()
				if !priorStateData.Conditions.IsNull() {
					conditions := []map[string]string{}
					resp.Diagnostics.Append(priorStateData.Conditions.ElementsAs(ctx, &conditions, false)...)
					if resp.Diagnostics.HasError() {
						return
					}

					if len(conditions) > 0 {
						upgradedStateData.Conditions = sentrytypes.NewLossyJsonValue(string(must.Get(json.Marshal(conditions))))
					}
				}

				upgradedStateData.Filters = sentrytypes.NewLossyJsonNull()
				if !priorStateData.Filters.IsNull() {
					filters := []map[string]string{}
					resp.Diagnostics.Append(priorStateData.Filters.ElementsAs(ctx, &filters, false)...)
					if resp.Diagnostics.HasError() {
						return
					}

					if len(filters) > 0 {
						upgradedStateData.Filters = sentrytypes.NewLossyJsonValue(string(must.Get(json.Marshal(filters))))
					}
				}

				upgradedStateData.Actions = sentrytypes.NewLossyJsonNull()
				if !priorStateData.Actions.IsNull() {
					actions := []map[string]string{}
					resp.Diagnostics.Append(priorStateData.Actions.ElementsAs(ctx, &actions, false)...)
					if resp.Diagnostics.HasError() {
						return
					}

					if len(actions) > 0 {
						upgradedStateData.Actions = sentrytypes.NewLossyJsonValue(string(must.Get(json.Marshal(actions))))
					}
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, &upgradedStateData)...)
			},
		},
	}
}
