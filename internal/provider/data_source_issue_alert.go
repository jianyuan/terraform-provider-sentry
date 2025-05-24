package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

var _ datasource.DataSource = &IssueAlertDataSource{}
var _ datasource.DataSourceWithConfigure = &IssueAlertDataSource{}

func NewIssueAlertDataSource() datasource.DataSource {
	return &IssueAlertDataSource{}
}

type IssueAlertDataSource struct {
	baseDataSource
}

func (d *IssueAlertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_alert"
}

func (d *IssueAlertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	stringAttribute := schema.StringAttribute{
		Computed: true,
	}
	int64Attribute := schema.Int64Attribute{
		Computed: true,
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Issue Alert data source. See the [Sentry documentation](https://docs.sentry.io/api/alerts/retrieve-an-issue-alert-rule-for-a-project/) for more information.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Required:    true,
			},
			"organization": DataSourceOrganizationAttribute(),
			"project":      DataSourceProjectAttribute(),
			"name": schema.StringAttribute{
				MarkdownDescription: "The issue alert name.",
				Computed:            true,
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "List of conditions. In JSON string format.",
				Computed:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"conditions_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of triggers that determine when the rule fires.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"first_seen_event": schema.SingleNestedAttribute{
							MarkdownDescription: "A new issue is created.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionFirstSeenEventModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name": stringAttribute,
							},
						},
						"regression_event": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue changes state from resolved to unresolved.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionRegressionEventModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name": stringAttribute,
							},
						},
						"reappeared_event": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue changes state from ignored to unresolved.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionReappearedEventModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name": stringAttribute,
							},
						},
						"new_high_priority_issue": schema.SingleNestedAttribute{
							MarkdownDescription: "Sentry marks a new issue as high priority.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionNewHighPriorityIssueModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name": stringAttribute,
							},
						},
						"existing_high_priority_issue": schema.SingleNestedAttribute{
							MarkdownDescription: "Sentry marks an existing issue as high priority.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionExistingHighPriorityIssueModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name": stringAttribute,
							},
						},
						"event_frequency": schema.SingleNestedAttribute{
							MarkdownDescription: "When the `comparison_type` is `count`, the number of events in an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of events in an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionEventFrequencyModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name":                stringAttribute,
								"comparison_type":     stringAttribute,
								"comparison_interval": stringAttribute,
								"value":               int64Attribute,
								"interval":            stringAttribute,
							},
						},
						"event_unique_user_frequency": schema.SingleNestedAttribute{
							MarkdownDescription: "When the `comparison_type` is `count`, the number of users affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of users affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionEventUniqueUserFrequencyModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name":                stringAttribute,
								"comparison_type":     stringAttribute,
								"comparison_interval": stringAttribute,
								"value":               int64Attribute,
								"interval":            stringAttribute,
							},
						},
						"event_frequency_percent": schema.SingleNestedAttribute{
							MarkdownDescription: "When the `comparison_type` is `count`, the percent of sessions affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the percent of sessions affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[IssueAlertConditionEventFrequencyPercentModel](ctx),
							Attributes: map[string]schema.Attribute{
								"name":                stringAttribute,
								"comparison_type":     stringAttribute,
								"comparison_interval": stringAttribute,
								"value": schema.Float64Attribute{
									Computed: true,
								},
								"interval": stringAttribute,
							},
						},
					},
				},
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met. In JSON string format.",
				Computed:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"filters_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"age_comparison": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue is older or newer than `value` `time`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":            stringAttribute,
								"comparison_type": stringAttribute,
								"value":           int64Attribute,
								"time":            stringAttribute,
							},
						},
						"issue_occurrences": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue has happened at least `value` times (Note: this is approximate).",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":  stringAttribute,
								"value": int64Attribute,
							},
						},
						"assigned_to": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue is assigned to no one, team, or member.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":              stringAttribute,
								"target_type":       int64Attribute,
								"target_identifier": int64Attribute,
							},
						},
						"latest_adopted_release": schema.SingleNestedAttribute{
							MarkdownDescription: "The {oldest_or_newest} adopted release associated with the event's issue is {older_or_newer} than the latest adopted release in {environment}.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":             stringAttribute,
								"oldest_or_newest": int64Attribute,
								"older_or_newer":   int64Attribute,
								"environment":      int64Attribute,
							},
						},
						"latest_release": schema.SingleNestedAttribute{
							MarkdownDescription: "The event is from the latest release.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": stringAttribute,
							},
						},
						"issue_category": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue's category is equal to `value`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":  stringAttribute,
								"value": stringAttribute,
							},
						},
						"event_attribute": schema.SingleNestedAttribute{
							MarkdownDescription: "The event's `attribute` value `match` `value`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":      stringAttribute,
								"attribute": stringAttribute,
								"match":     stringAttribute,
								"value":     stringAttribute,
							},
						},
						"tagged_event": schema.SingleNestedAttribute{
							MarkdownDescription: "The event's tags match `key` `match` `value`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":  stringAttribute,
								"key":   stringAttribute,
								"match": stringAttribute,
								"value": stringAttribute,
							},
						},
						"level": schema.SingleNestedAttribute{
							MarkdownDescription: "The event's level is `match` `level`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":  stringAttribute,
								"match": stringAttribute,
								"level": stringAttribute,
							},
						},
					},
				},
			},
			"actions": schema.StringAttribute{
				MarkdownDescription: "List of actions. In JSON string format.",
				Computed:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"actions_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of actions that take place when all required conditions and filters for the rule are met.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"notify_email": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to `target_type` and if none can be found then send a notification to `fallthrough_type`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":              stringAttribute,
								"target_type":       stringAttribute,
								"target_identifier": stringAttribute,
								"fallthrough_type":  stringAttribute,
							},
						},
						"notify_event": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to all legacy integrations.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": stringAttribute,
							},
						},
						"notify_event_service": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification via an integration.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":    stringAttribute,
								"service": stringAttribute,
							},
						},
						"notify_event_sentry_app": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to a Sentry app.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":                         stringAttribute,
								"sentry_app_installation_uuid": stringAttribute,
								"settings": schema.MapAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"opsgenie_notify_team": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to Opsgenie account `account` and team `team` with `priority` priority.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":     stringAttribute,
								"account":  stringAttribute,
								"team":     stringAttribute,
								"priority": stringAttribute,
							},
						},
						"pagerduty_notify_service": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to PagerDuty account `account` and service `service` with `severity` severity.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":     stringAttribute,
								"account":  stringAttribute,
								"service":  stringAttribute,
								"severity": stringAttribute,
							},
						},
						"slack_notify_service": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to the `workspace` Slack workspace to `channel` (optionally, an ID: `channel_id`) and show tags `tags` and notes `notes` in notification.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":       stringAttribute,
								"workspace":  stringAttribute,
								"channel":    stringAttribute,
								"channel_id": stringAttribute,
								"tags": schema.SetAttribute{
									Computed: true,
									CustomType: sentrytypes.StringSetType{
										SetType: types.SetType{
											ElemType: types.StringType,
										},
									},
								},
								"notes": stringAttribute,
							},
						},
						"msteams_notify_service": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to the `team` Team to `channel`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":       stringAttribute,
								"team":       stringAttribute,
								"channel":    stringAttribute,
								"channel_id": stringAttribute,
							},
						},
						"discord_notify_service": schema.SingleNestedAttribute{
							MarkdownDescription: "Send a notification to the `server` Discord server in the channel with ID or URL: `channel_id` and show tags `tags` in the notification.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":       stringAttribute,
								"server":     stringAttribute,
								"channel_id": stringAttribute,
								"tags": schema.SetAttribute{
									Computed: true,
									CustomType: sentrytypes.StringSetType{
										SetType: types.SetType{
											ElemType: types.StringType,
										},
									},
								},
							},
						},
						"jira_create_ticket": schema.SingleNestedAttribute{
							MarkdownDescription: "Create a Jira issue in `integration`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":        stringAttribute,
								"integration": stringAttribute,
								"project":     stringAttribute,
								"issue_type":  stringAttribute,
							},
						},
						"jira_server_create_ticket": schema.SingleNestedAttribute{
							MarkdownDescription: "Create a Jira Server issue in `integration`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":        stringAttribute,
								"integration": stringAttribute,
								"project":     stringAttribute,
								"issue_type":  stringAttribute,
							},
						},
						"github_create_ticket": schema.SingleNestedAttribute{
							MarkdownDescription: "Create a GitHub issue in `integration`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":        stringAttribute,
								"integration": stringAttribute,
								"repo":        stringAttribute,
								"assignee":    stringAttribute,
								"labels": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"github_enterprise_create_ticket": schema.SingleNestedAttribute{
							MarkdownDescription: "Create a GitHub Enterprise issue in `integration`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":        stringAttribute,
								"integration": stringAttribute,
								"repo":        stringAttribute,
								"assignee":    stringAttribute,
								"labels": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"azure_devops_create_ticket": schema.SingleNestedAttribute{
							MarkdownDescription: "Create an Azure DevOps work item in `integration`.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":           stringAttribute,
								"integration":    stringAttribute,
								"work_item_type": stringAttribute,
							},
						},
					},
				},
			},
			"action_match": schema.StringAttribute{
				MarkdownDescription: "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.",
				Computed:            true,
			},
			"filter_match": schema.StringAttribute{
				MarkdownDescription: "A string determining which filters need to be true before any actions take place. Required when a value is provided for `filters`.",
				Computed:            true,
			},
			"frequency": schema.Int64Attribute{
				MarkdownDescription: "Perform actions at most once every `X` minutes for this issue.",
				Computed:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Perform issue alert in a specific environment.",
				Computed:            true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The ID of the team or user that owns the rule.",
				Computed:            true,
			},
		},
	}
}

func (d *IssueAlertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.apiClient.GetProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("issue alert"))
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
