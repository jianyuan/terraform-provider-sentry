package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

var _ resource.Resource = &IssueAlertResource{}
var _ resource.ResourceWithConfigValidators = &IssueAlertResource{}
var _ resource.ResourceWithValidateConfig = &IssueAlertResource{}
var _ resource.ResourceWithConfigure = &IssueAlertResource{}
var _ resource.ResourceWithImportState = &IssueAlertResource{}
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
	nameStringAttribute := schema.StringAttribute{
		Computed: true,
	}
	intervalStringAttribute := tfutils.WithEnumStringAttribute(schema.StringAttribute{
		MarkdownDescription: "`m` for minutes, `h` for hours, `d` for days, and `w` for weeks.",
		Optional:            true,
	}, []string{"1m", "5m", "15m", "1h", "1d", "1w", "30d"})
	conditionComparisonTypeStringAttribute := tfutils.WithEnumStringAttribute(schema.StringAttribute{
		Required: true,
	}, []string{"count", "percent"})
	conditionComparisonIntervalStringAttribute := tfutils.WithEnumStringAttribute(schema.StringAttribute{
		MarkdownDescription: "`m` for minutes, `h` for hours, `d` for days, and `w` for weeks.",
		Optional:            true,
	}, []string{"5m", "15m", "1h", "1d", "1w", "30d"})

	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an Issue Alert Rule for a Project. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/) for more information.\n\n" +
			"**NOTE:** Since v0.15.0, the `conditions`, `filters`, and `actions` attributes which are JSON strings have been deprecated in favor of `conditions_v2`, `filters_v2`, and `actions_v2` which are lists of objects.",

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
				MarkdownDescription: "**Deprecated** in favor of `conditions_v2`. A list of triggers that determine when the rule fires. In JSON string format.",
				DeprecationMessage:  "Use `conditions_v2` instead.",
				Optional:            true,
				CustomType: sentrytypes.LossyJsonType{
					IgnoreKeys: []string{"name"},
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("conditions_v2")),
				},
			},
			"conditions_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of triggers that determine when the rule fires.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("conditions")),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: tfutils.WithMutuallyExclusiveValidator(map[string]schema.SingleNestedAttribute{
						"first_seen_event": {
							MarkdownDescription: "A new issue is created.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"regression_event": {
							MarkdownDescription: "The issue changes state from resolved to unresolved.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"reappeared_event": {
							MarkdownDescription: "The issue changes state from ignored to unresolved.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"new_high_priority_issue": {
							MarkdownDescription: "Sentry marks a new issue as high priority.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"existing_high_priority_issue": {
							MarkdownDescription: "Sentry marks an existing issue as high priority.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"event_frequency": {
							MarkdownDescription: "When the `comparison_type` is `count`, the number of events in an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of events in an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Int64Attribute{
									Required: true,
								},
								"interval": intervalStringAttribute,
							},
						},
						"event_unique_user_frequency": {
							MarkdownDescription: "When the `comparison_type` is `count`, the number of users affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of users affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Int64Attribute{
									Required: true,
								},
								"interval": intervalStringAttribute,
							},
						},
						"event_frequency_percent": {
							MarkdownDescription: "When the `comparison_type` is `count`, the percent of sessions affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the percent of sessions affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Float64Attribute{
									Required: true,
								},
								"interval": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									MarkdownDescription: "`m` for minutes, `h` for hours.",
									Required:            true,
								}, []string{"5m", "10m", "30m", "1h"}),
							},
						},
					}),
				},
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "**Deprecated** in favor of `filters_v2`. A list of filters that determine if a rule fires after the necessary conditions have been met. In JSON string format.",
				DeprecationMessage:  "Use `filters_v2` instead.",
				Optional:            true,
				CustomType:          sentrytypes.LossyJsonType{},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("filters_v2")),
				},
			},
			"filters_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("filters")),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: tfutils.WithMutuallyExclusiveValidator(map[string]schema.SingleNestedAttribute{
						"age_comparison": {
							MarkdownDescription: "The issue is older or newer than `value` `time`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"comparison_type": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, []string{"older", "newer"}),
								"value": schema.Int64Attribute{
									Required: true,
								},
								"time": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, []string{"minute", "hour", "day", "week"}),
							},
						},
						"issue_occurrences": {
							MarkdownDescription: "The issue has happened at least `value` times (Note: this is approximate).",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"value": schema.Int64Attribute{
									Required: true,
								},
							},
						},
						"assigned_to": {
							MarkdownDescription: "The issue is assigned to no one, team, or member.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"target_type": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, []string{"Unassigned", "Team", "Member"}),
								"target_identifier": schema.StringAttribute{
									MarkdownDescription: "The target's ID. Only required when `target_type` is `Team` or `Member`.",
									Optional:            true,
								},
							},
						},
						"latest_adopted_release": {
							MarkdownDescription: "The {oldest_or_newest} adopted release associated with the event's issue is {older_or_newer} than the latest adopted release in {environment}.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"oldest_or_newest": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, []string{"oldest", "newest"}),
								"older_or_newer": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, []string{"older", "newer"}),
								"environment": schema.StringAttribute{
									Required: true,
								},
							},
						},
						"latest_release": {
							MarkdownDescription: "The event is from the latest release.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"issue_category": {
							MarkdownDescription: "The issue's category is equal to `value`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"value": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, sentrydata.IssueGroupCategories),
							},
						},
						"event_attribute": {
							MarkdownDescription: "The event's `attribute` value `match` `value`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"attribute": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, sentrydata.EventAttributes),
								"match": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									MarkdownDescription: "The comparison operator.",
									Required:            true,
								}, sentrydata.MatchTypes),
								"value": schema.StringAttribute{
									Optional: true,
								},
							},
						},
						"tagged_event": {
							MarkdownDescription: "The event's tags match `key` `match` `value`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"key": schema.StringAttribute{
									MarkdownDescription: "The tag.",
									Required:            true,
								},
								"match": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									MarkdownDescription: "The comparison operator.",
									Required:            true,
								}, sentrydata.MatchTypes),
								"value": schema.StringAttribute{
									Optional: true,
								},
							},
						},
						"level": {
							MarkdownDescription: "The event's level is `match` `level`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"match": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									MarkdownDescription: "The comparison operator.",
									Required:            true,
								}, sentrydata.LevelMatchTypes),
								"level": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, sentrydata.LogLevels),
							},
						},
					}),
				},
			},
			"actions": schema.StringAttribute{
				MarkdownDescription: "**Deprecated** in favor of `actions_v2`. A list of actions that take place when all required conditions and filters for the rule are met. In JSON string format.",
				DeprecationMessage:  "Use `actions_v2` instead.",
				Optional:            true,
				CustomType:          sentrytypes.LossyJsonType{},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("actions_v2")),
				},
			},
			"actions_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of actions that take place when all required conditions and filters for the rule are met.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("actions")),
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: tfutils.WithMutuallyExclusiveValidator(map[string]schema.SingleNestedAttribute{
						"notify_email": {
							MarkdownDescription: "Send a notification to `target_type` and if none can be found then send a notification to `fallthrough_type`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"target_type": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									Required: true,
								}, []string{"IssueOwners", "Team", "Member"}),
								"target_identifier": schema.StringAttribute{
									MarkdownDescription: "The ID of the Member or Team the notification should be sent to. Only required when `target_type` is `Team` or `Member`.",
									Optional:            true,
								},
								"fallthrough_type": tfutils.WithEnumStringAttribute(schema.StringAttribute{
									MarkdownDescription: "Who the notification should be sent to if there are no suggested assignees.",
									Optional:            true,
								}, []string{"AllMembers", "ActiveMembers", "NoOne"}),
							},
						},
						"notify_event": {
							MarkdownDescription: "Send a notification to all legacy integrations.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"notify_event_service": {
							MarkdownDescription: "Send a notification via an integration.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"service": schema.StringAttribute{
									MarkdownDescription: "The slug of the integration service. Sourced from `https://terraform-provider-sentry.sentry.io/settings/developer-settings/<service>/`.",
									Required:            true,
								},
							},
						},
						"notify_event_sentry_app": {
							MarkdownDescription: "Send a notification to a Sentry app.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"sentry_app_installation_uuid": schema.StringAttribute{
									Required: true,
								},
								"settings": schema.MapAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
							},
						},
						"opsgenie_notify_team": {
							MarkdownDescription: "Send a notification to Opsgenie account `account` and team `team` with `priority` priority.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
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
						"pagerduty_notify_service": {
							MarkdownDescription: "Send a notification to PagerDuty account `account` and service `service` with `severity` severity.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
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
						"slack_notify_service": {
							MarkdownDescription: "Send a notification to the `workspace` Slack workspace to `channel` (optionally, an ID: `channel_id`) and show tags `tags` and notes `notes` in notification.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"workspace": schema.StringAttribute{
									MarkdownDescription: "The integration ID associated with the Slack workspace.",
									Required:            true,
								},
								"channel": schema.StringAttribute{
									MarkdownDescription: "The name of the channel to send the notification to (e.g., #critical, Jane Schmidt).",
									Required:            true,
								},
								"channel_id": schema.StringAttribute{
									MarkdownDescription: "The ID of the channel to send the notification to.",
									Computed:            true,
								},
								"tags": schema.SetAttribute{
									MarkdownDescription: "A string of tags to show in the notification.",
									Optional:            true,
									CustomType: sentrytypes.StringSetType{
										SetType: types.SetType{
											ElemType: types.StringType,
										},
									},
								},
								"notes": schema.StringAttribute{
									MarkdownDescription: "Text to show alongside the notification. To @ a user, include their user id like `@<USER_ID>`. To include a clickable link, format the link and title like `<http://example.com|Click Here>`.",
									Optional:            true,
								},
							},
						},
						"msteams_notify_service": {
							MarkdownDescription: "Send a notification to the `team` Team to `channel`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"team": schema.StringAttribute{
									MarkdownDescription: "The integration ID associated with the Microsoft Teams team.",
									Required:            true,
								},
								"channel": schema.StringAttribute{
									MarkdownDescription: "The name of the channel to send the notification to.",
									Required:            true,
								},
								"channel_id": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"discord_notify_service": {
							MarkdownDescription: "Send a notification to the `server` Discord server in the channel with ID or URL: `channel_id` and show tags `tags` in the notification.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"server": schema.StringAttribute{
									MarkdownDescription: "The integration ID associated with the Discord server.",
									Required:            true,
								},
								"channel_id": schema.StringAttribute{
									MarkdownDescription: "The ID of the channel to send the notification to. You must enter either a channel ID or a channel URL, not a channel name",
									Required:            true,
								},
								"tags": schema.SetAttribute{
									MarkdownDescription: "A string of tags to show in the notification.",
									Optional:            true,
									CustomType: sentrytypes.StringSetType{
										SetType: types.SetType{
											ElemType: types.StringType,
										},
									},
								},
							},
						},
						"jira_create_ticket": {
							MarkdownDescription: "Create a Jira issue in `integration`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"integration": schema.StringAttribute{
									MarkdownDescription: "The integration ID associated with Jira.",
									Required:            true,
								},
								"project": schema.StringAttribute{
									MarkdownDescription: "The ID of the Jira project.",
									Required:            true,
								},
								"issue_type": schema.StringAttribute{
									MarkdownDescription: "The ID of the type of issue that the ticket should be created as.",
									Required:            true,
								},
							},
						},
						"jira_server_create_ticket": {
							MarkdownDescription: "Create a Jira Server issue in `integration`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"integration": schema.StringAttribute{
									MarkdownDescription: "The integration ID associated with Jira Server.",
									Required:            true,
								},
								"project": schema.StringAttribute{
									MarkdownDescription: "The ID of the Jira Server project.",
									Required:            true,
								},
								"issue_type": schema.StringAttribute{
									MarkdownDescription: "The ID of the type of issue that the ticket should be created as.",
									Required:            true,
								},
							},
						},
						"github_create_ticket": {
							MarkdownDescription: "Create a GitHub issue in `integration`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"integration": schema.StringAttribute{
									MarkdownDescription: "The integration ID associated with GitHub.",
									Required:            true,
								},
								"repo": schema.StringAttribute{
									MarkdownDescription: "The name of the repository to create the issue in.",
									Required:            true,
								},
								"assignee": schema.StringAttribute{
									MarkdownDescription: "The GitHub user to assign the issue to.",
									Optional:            true,
								},
								"labels": schema.SetAttribute{
									MarkdownDescription: "A list of labels to assign to the issue.",
									Optional:            true,
									ElementType:         types.StringType,
								},
							},
						},
						"github_enterprise_create_ticket": {
							MarkdownDescription: "Create a GitHub Enterprise issue in `integration`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"integration": schema.StringAttribute{
									MarkdownDescription: "The integration ID associated with GitHub Enterprise.",
									Required:            true,
								},
								"repo": schema.StringAttribute{
									MarkdownDescription: "The name of the repository to create the issue in.",
									Required:            true,
								},
								"assignee": schema.StringAttribute{
									MarkdownDescription: "The GitHub user to assign the issue to.",
									Optional:            true,
								},
								"labels": schema.SetAttribute{
									MarkdownDescription: "A list of labels to assign to the issue.",
									Optional:            true,
									ElementType:         types.StringType,
								},
							},
						},
						"azure_devops_create_ticket": {
							MarkdownDescription: "Create an Azure DevOps work item in `integration`.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
								"integration": schema.StringAttribute{
									MarkdownDescription: "The integration ID.",
									Required:            true,
								},
								"project": schema.StringAttribute{
									MarkdownDescription: "The ID of the Azure DevOps project.",
									Required:            true,
								},
								"work_item_type": schema.StringAttribute{
									MarkdownDescription: "The type of work item to create.",
									Required:            true,
								},
							},
						},
					}),
				},
			},
			"action_match": tfutils.WithEnumStringAttribute(schema.StringAttribute{
				MarkdownDescription: "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.",
				Required:            true,
			}, []string{"all", "any"}),
			"filter_match": tfutils.WithEnumStringAttribute(schema.StringAttribute{
				MarkdownDescription: "A string determining which filters need to be true before any actions take place. Required when a value is provided for `filters`.",
				Optional:            true,
			}, []string{"all", "any", "none"}),
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
	}
}

func (r *IssueAlertResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("actions"),
			path.MatchRoot("actions_v2"),
		),
	}
}

func (r *IssueAlertResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ConditionsV2 != nil {
		for i, item := range *data.ConditionsV2 {
			if _, diags := item.ToApi(ctx); diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("conditions_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert condition: %s", diags),
				)
			}
		}
	}

	if data.FiltersV2 != nil {
		for i, item := range *data.FiltersV2 {
			if _, diags := item.ToApi(ctx); diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("filters_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert filter: %s", diags),
				)
			}
		}
	}

	if !data.Actions.IsNull() {
		if ok, _ := data.Actions.StringSemanticEquals(ctx, sentrytypes.NewLossyJsonValue(`[]`)); ok {
			resp.Diagnostics.AddAttributeError(
				path.Root("actions"),
				"Missing attribute configuration",
				"You must add an action for this alert to fire",
			)
		}
	} else if data.ActionsV2 != nil {
		if len(*data.ActionsV2) == 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("actions_v2"),
				"Missing attribute configuration",
				"You must add an action for this alert to fire",
			)
		}

		for i, item := range *data.ActionsV2 {
			if _, diags := item.ToApi(ctx); diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("actions_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert action: %s", diags),
				)
			}
		}
	}
}

func (r *IssueAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.CreateProjectRuleJSONRequestBody{
		Name:        data.Name.ValueString(),
		ActionMatch: data.ActionMatch.ValueString(),
		FilterMatch: data.FilterMatch.ValueString(),
		Frequency:   data.Frequency.ValueInt64(),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.ValueString()},
	}

	if !data.Conditions.IsNull() {
		resp.Diagnostics.Append(data.Conditions.Unmarshal(&body.Conditions)...)
	} else if data.ConditionsV2 != nil {
		body.Conditions = []apiclient.ProjectRuleCondition{}
		for i, item := range *data.ConditionsV2 {
			condition, diags := item.ToApi(ctx)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("conditions_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert condition: %s", diags),
				)
				return
			}
			body.Conditions = append(body.Conditions, *condition)
		}
	} else {
		body.Conditions = []apiclient.ProjectRuleCondition{}
	}

	if !data.Filters.IsNull() {
		resp.Diagnostics.Append(data.Filters.Unmarshal(&body.Filters)...)
	} else if data.FiltersV2 != nil {
		body.Filters = []apiclient.ProjectRuleFilter{}
		for i, item := range *data.FiltersV2 {
			filter, diags := item.ToApi(ctx)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("filters_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert filter: %s", diags),
				)
				return
			}
			body.Filters = append(body.Filters, *filter)
		}
	} else {
		body.Filters = []apiclient.ProjectRuleFilter{}
	}

	if !data.Actions.IsNull() {
		resp.Diagnostics.Append(data.Actions.Unmarshal(&body.Actions)...)
	} else if data.ActionsV2 != nil {
		body.Actions = []apiclient.ProjectRuleAction{}
		for i, item := range *data.ActionsV2 {
			action, diags := item.ToApi(ctx)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("actions_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert action: %s", diags),
				)
				return
			}
			body.Actions = append(body.Actions, *action)
		}
	} else {
		body.Actions = []apiclient.ProjectRuleAction{}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CreateProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.UpdateProjectRuleJSONRequestBody{
		Name:        data.Name.ValueString(),
		ActionMatch: data.ActionMatch.ValueString(),
		FilterMatch: data.FilterMatch.ValueString(),
		Frequency:   data.Frequency.ValueInt64(),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.ValueString()},
	}

	if !data.Conditions.IsNull() {
		resp.Diagnostics.Append(data.Conditions.Unmarshal(&body.Conditions)...)
	} else if data.ConditionsV2 != nil {
		body.Conditions = []apiclient.ProjectRuleCondition{}
		for i, item := range *data.ConditionsV2 {
			condition, diags := item.ToApi(ctx)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("conditions_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert condition: %s", diags),
				)
				return
			}
			body.Conditions = append(body.Conditions, *condition)
		}
	} else {
		body.Conditions = []apiclient.ProjectRuleCondition{}
	}

	if !data.Filters.IsNull() {
		resp.Diagnostics.Append(data.Filters.Unmarshal(&body.Filters)...)
	} else if data.FiltersV2 != nil {
		body.Filters = []apiclient.ProjectRuleFilter{}
		for i, item := range *data.FiltersV2 {
			filter, diags := item.ToApi(ctx)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("filters_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert filter: %s", diags),
				)
				return
			}
			body.Filters = append(body.Filters, *filter)
		}
	} else {
		body.Filters = []apiclient.ProjectRuleFilter{}
	}

	if !data.Actions.IsNull() {
		resp.Diagnostics.Append(data.Actions.Unmarshal(&body.Actions)...)
	} else if data.ActionsV2 != nil {
		body.Actions = []apiclient.ProjectRuleAction{}
		for i, item := range *data.ActionsV2 {
			action, diags := item.ToApi(ctx)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(
					path.Root("actions_v2").AtListIndex(i),
					"Missing attribute configuration",
					fmt.Sprintf("Failed to convert action: %s", diags),
				)
				return
			}
			body.Actions = append(body.Actions, *action)
		}
	} else {
		body.Actions = []apiclient.ProjectRuleAction{}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UpdateProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("issue alert"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() != http.StatusAccepted {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *IssueAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateThreePartId(ctx, "organization", "project", req, resp)
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

				organization, project, actionId, err := tfutils.SplitThreePartId(priorStateData.Id.ValueString(), "organization", "project-slug", "alert-id")
				if err != nil {
					resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
					return
				}

				upgradedStateData := IssueAlertModel{
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
