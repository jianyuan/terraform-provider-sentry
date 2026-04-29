import dedent from "dedent";
import type { Resource } from "../schema";
import { withExactlyOneAttribute } from "../utils";

export default {
  name: "alert",
  description: dedent.withOptions({ trimWhitespace: true })`
      ⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

      Create an Alert for a Monitor in an Organization. Monitors must be created separately using the [\`sentry_cron_monitor\`](cron_monitor.md), [\`sentry_metric_monitor\`](metric_monitor.md), or [\`sentry_uptime_monitor\`](uptime_monitor.md) resources.
    `,
  api: {
    model: "OrganizationWorkflow",
    createMethod: "CreateOrganizationWorkflow",
    createRequestAttributes: ["organization"],
    readMethod: "GetOrganizationWorkflow",
    readRequestAttributes: ["organization", "id"],
    updateMethod: "UpdateOrganizationWorkflow",
    updateRequestAttributes: ["organization", "id"],
    deleteMethod: "DeleteOrganizationWorkflow",
    deleteRequestAttributes: ["organization", "id"],
  },
  generate: {
    modelFillers: false,
  },
  importStateAttributes: ["organization", "id"],
  attributes: [
    {
      name: "id",
      type: "string",
      description: "The internal ID of this alert.",
      computedOptionalRequired: "computed",
      sourceAttribute: ["Id"],
      planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
    },
    {
      name: "organization",
      type: "string",
      description:
        "The organization slug or internal ID to create the alert for.",
      computedOptionalRequired: "required",
      planModifiers: ["stringplanmodifier.RequiresReplace()"],
    },
    {
      name: "enabled",
      type: "bool",
      description: "Whether the alert is enabled. Defaults to `true`.",
      computedOptionalRequired: "computed_optional",
      default: `booldefault.StaticBool(true)`,
    },
    {
      name: "name",
      type: "string",
      description: "The name of this alert.",
      computedOptionalRequired: "required",
    },
    {
      name: "environment",
      type: "string",
      description: "Name of the environment to create alerts in.",
      computedOptionalRequired: "required",
    },
    {
      name: "monitor_ids",
      type: "set",
      description: "The IDs of the monitors to create alerts for.",
      computedOptionalRequired: "required",
      elementType: "string",
    },
    {
      name: "frequency_minutes",
      type: "int",
      description: "How often the alert should fire in minutes.",
      computedOptionalRequired: "required",
    },
    {
      name: "trigger_conditions",
      type: "list_nested",
      description: "The conditions on which the alert will trigger.",
      computedOptionalRequired: "required",
      validators: ["listvalidator.SizeAtLeast(1)"],
      attributes: withExactlyOneAttribute([
        {
          name: "first_seen_event",
          type: "single_nested",
          description: "A new issue is created.",
          computedOptionalRequired: "optional",
          attributes: [],
        },
        {
          name: "issue_resolved_trigger",
          type: "single_nested",
          description: "An issue is resolved.",
          computedOptionalRequired: "optional",
          attributes: [],
        },
        {
          name: "reappeared_event",
          type: "single_nested",
          description: "An issue escalates.",
          computedOptionalRequired: "optional",
          attributes: [],
        },
        {
          name: "regression_event",
          type: "single_nested",
          description: "A resolved issue becomes unresolved.",
          computedOptionalRequired: "optional",
          attributes: [],
        },
      ]),
    },
    {
      name: "action_filters",
      type: "list_nested",
      description:
        "The filters to run before the action will fire and the action(s) to fire.",
      computedOptionalRequired: "required",
      validators: ["listvalidator.SizeAtLeast(1)"],
      attributes: [
        {
          name: "logic_type",
          type: "string",
          description:
            "The logic to apply to the conditions. `any` will evaluate all conditions, and return true if any of those are met. `any-short` will stop evaluating conditions as soon as one is met. `all` will evaluate all conditions, and return true if all of those are met. `none` will return true if none of the conditions are met, will return false immediately if any are met.",
          computedOptionalRequired: "required",
          enum: "sentrydata.DataConditionGroupTypes",
        },
        {
          name: "conditions",
          type: "list_nested",
          description: "The conditions to evaluate.",
          computedOptionalRequired: "computed_optional",
          attributes: withExactlyOneAttribute([
            {
              name: "age_comparison",
              type: "single_nested",
              description: "Issue age.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "time",
                  type: "string",
                  description: "The unit of time for the age comparison.",
                  computedOptionalRequired: "required",
                  enum: `[]string{"minute", "hour", "day", "week"}`,
                },
                {
                  name: "value",
                  type: "int",
                  description: "The value of the age comparison.",
                  computedOptionalRequired: "required",
                  validators: ["int64validator.AtLeast(1)"],
                },
                {
                  name: "comparison_type",
                  type: "string",
                  description: "The type of comparison to perform.",
                  computedOptionalRequired: "required",
                  enum: `[]string{"older", "newer"}`,
                },
              ],
            },
            {
              name: "assigned_to",
              type: "single_nested",
              description: "Issue assignment.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "target_type",
                  type: "string",
                  description: "Who the issue is assigned to.",
                  computedOptionalRequired: "optional",
                  enum: `[]string{"Unassigned", "Member", "Team"}`,
                },
                {
                  name: "target_id",
                  type: "string",
                  description:
                    "The internal ID of the user or team. Only required if the target type is `Member` or `Team`.",
                  computedOptionalRequired: "optional",
                  validators: [
                    `fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("target_type"), []attr.Value{supertypes.NewStringValue("Member"), supertypes.NewStringValue("Team")})`,
                    `fstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("target_type"), []attr.Value{supertypes.NewStringValue("Unassigned")})`,
                  ],
                },
              ],
            },
            {
              name: "issue_category",
              type: "single_nested",
              description: "Issue category.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "value",
                  type: "int",
                  description: "The issue category to filter to.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "issue_occurrences",
              type: "single_nested",
              description: "Issue frequency.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "value",
                  type: "int",
                  description:
                    "A positive integer representing how many times the issue has to happen before the alert will fire.",
                  computedOptionalRequired: "required",
                  validators: ["int64validator.AtLeast(1)"],
                },
              ],
            },
            {
              name: "issue_priority_deescalating",
              type: "single_nested",
              description: "De-escalation.",
              computedOptionalRequired: "optional",
              attributes: [],
            },
            {
              name: "issue_priority_greater_or_equal",
              type: "single_nested",
              description: "Issue priority.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "comparison",
                  type: "int",
                  description:
                    "he priority the issue must be for the alert to fire.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "event_unique_user_frequency_count",
              type: "single_nested",
              description: "Number of users affected.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "value",
                  type: "int",
                  description:
                    "A positive integer representing the number of users that must be affected before the alert will fire.",
                  computedOptionalRequired: "required",
                  validators: ["int64validator.AtLeast(1)"],
                },
                {
                  name: "filters",
                  type: "list_nested",
                  description:
                    "A list of additional sub-filters to evaluate before the alert will fire.",
                  computedOptionalRequired: "computed_optional",
                  attributes: [
                    {
                      name: "key",
                      type: "string",
                      description:
                        "The key of the filter. Conflicts with `attribute`.",
                      computedOptionalRequired: "optional",
                      validators: [
                        `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("attribute"))`,
                      ],
                    },
                    {
                      name: "attribute",
                      type: "string",
                      description:
                        "The attribute of the filter. Conflicts with `key`.",
                      computedOptionalRequired: "optional",
                      validators: [
                        `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("key"))`,
                      ],
                    },
                    {
                      name: "match",
                      type: "string",
                      description: "The match type of the filter.",
                      computedOptionalRequired: "optional",
                    },
                    {
                      name: "value",
                      type: "string",
                      description: "The value of the filter.",
                      computedOptionalRequired: "optional",
                    },
                  ],
                },
                {
                  name: "interval",
                  type: "string",
                  description:
                    "The time period in which to evaluate the value. e.g. Number of users affected by an issue is more than `value` in `interval`.",
                  computedOptionalRequired: "required",
                  enum: `sentrydata.EventFrequencyStandardIntervals`,
                },
              ],
            },
            {
              name: "event_frequency_count",
              type: "single_nested",
              description: "Number of events.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "value",
                  type: "int",
                  description:
                    "A positive integer representing the number of events in an issue that must come in before the alert will fire.",
                  computedOptionalRequired: "required",
                  validators: ["int64validator.AtLeast(1)"],
                },
                {
                  name: "interval",
                  type: "string",
                  description:
                    "The time period in which to evaluate the value. e.g. Number of events in an issue is more than `value` in `interval`.",
                  computedOptionalRequired: "required",
                  enum: `sentrydata.EventFrequencyStandardIntervals`,
                },
              ],
            },
            {
              name: "event_frequency_percent",
              type: "single_nested",
              description: "Percent of events.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "value",
                  type: "int",
                  description:
                    "A positive integer representing the number of events in an issue that must come in before the alert will fire.",
                  computedOptionalRequired: "required",
                  validators: ["int64validator.AtLeast(1)"],
                },
                {
                  name: "interval",
                  type: "string",
                  description:
                    "The time period in which to evaluate the value. e.g. Number of events in an issue is `comparisonInterval` percent higher `value` compared to `interval`.",
                  computedOptionalRequired: "required",
                  enum: `sentrydata.EventFrequencyStandardIntervals`,
                },
                {
                  name: "comparison_interval",
                  type: "string",
                  description: "The time period to compare against.",
                  computedOptionalRequired: "required",
                  enum: `sentrydata.EventFrequencyStandardIntervals`,
                },
              ],
            },
            {
              name: "percent_sessions_count",
              type: "single_nested",
              description: "Percentage of sessions affected count.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "value",
                  type: "int",
                  description:
                    "A positive integer representing the number of events in an issue that must come in before the alert will fire.",
                  computedOptionalRequired: "required",
                  validators: ["int64validator.AtLeast(1)"],
                },
                {
                  name: "interval",
                  type: "string",
                  description:
                    "The time period in which to evaluate the value. e.g. Percentage of sessions affected by an issue is more than `value` in `interval`.",
                  computedOptionalRequired: "required",
                  enum: `sentrydata.EventFrequencyStandardIntervals`,
                },
              ],
            },
            {
              name: "percent_sessions_percent",
              type: "single_nested",
              description: "Percentage of sessions affected percent.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "value",
                  type: "int",
                  description:
                    "A positive integer representing the number of events in an issue that must come in before the alert will fire.",
                  computedOptionalRequired: "required",
                  validators: ["int64validator.AtLeast(1)"],
                },
                {
                  name: "interval",
                  type: "string",
                  description:
                    "The time period in which to evaluate the value. e.g. Percentage of sessions affected by an issue is `comparisonInterval` percent higher `value` compared to `interval`.",
                  computedOptionalRequired: "required",
                  enum: `sentrydata.EventFrequencyStandardIntervals`,
                },
                {
                  name: "comparison_interval",
                  type: "string",
                  description: "The time period to compare against.",
                  computedOptionalRequired: "required",
                  enum: `sentrydata.EventFrequencyStandardIntervals`,
                },
              ],
            },
            {
              name: "event_attribute",
              type: "single_nested",
              description: "The event's `attribute` value `match` `value`.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "attribute",
                  type: "string",
                  description: "The attribute to evaluate.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "match",
                  type: "string",
                  description: "The match type.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "value",
                  type: "string",
                  description: "The value to compare against.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "tagged_event",
              type: "single_nested",
              description: "The event's tags `key` match `value`.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "key",
                  type: "string",
                  description: "The tag value.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "match",
                  type: "string",
                  description: "The comparison operator.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "value",
                  type: "string",
                  description:
                    "A string. Not required when match is `is` or `ns`.",
                  computedOptionalRequired: "optional",
                  validators: [
                    // TODO: Require value when match is not `is` or `ns`
                    `fstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("match"), []attr.Value{supertypes.NewStringValue("is"), supertypes.NewStringValue("ns")})`,
                  ],
                },
              ],
            },
            {
              name: "latest_release",
              type: "single_nested",
              description: "The event is from the latest release.",
              computedOptionalRequired: "optional",
              attributes: [],
            },
            {
              name: "latest_adopted_release",
              type: "single_nested",
              description:
                "The `release_age_type` adopted release associated with the event's issue is `age_comparison` than the latest adopted release in `environment`.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "environment",
                  type: "string",
                  description: "The environment to compare against.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "age_comparison",
                  type: "string",
                  description: "The age comparison to use.",
                  computedOptionalRequired: "required",
                  enum: `[]string{"older", "newer"}`,
                },
                {
                  name: "release_age_type",
                  type: "string",
                  description: "The release age type to use.",
                  computedOptionalRequired: "required",
                  enum: `[]string{"oldest", "newest"}`,
                },
              ],
            },
            {
              name: "level",
              type: "single_nested",
              description: "The event's level match `level`.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "match",
                  type: "string",
                  description: "The comparison operator.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "level",
                  type: "int",
                  description: "The level to compare against.",
                  computedOptionalRequired: "required",
                },
              ],
            },
          ]),
        },
        {
          name: "actions",
          type: "list_nested",
          description: "The actions to perform.",
          computedOptionalRequired: "required",
          validators: ["listvalidator.SizeAtLeast(1)"],
          attributes: withExactlyOneAttribute([
            {
              name: "email",
              type: "single_nested",
              description: "Notify on Preferred Channel.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "target_type",
                  type: "string",
                  description: "The type of recipient to notify.",
                  computedOptionalRequired: "required",
                  enum: `[]string{"issue_owners", "team", "user"}`,
                },
                {
                  name: "target_id",
                  type: "string",
                  description:
                    "The internal ID of the user or team. Only required if the target type is `team` or `user`.",
                  computedOptionalRequired: "optional",
                  validators: [
                    `fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("target_type"), []attr.Value{supertypes.NewStringValue("team"), supertypes.NewStringValue("user")})`,
                    `fstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("target_type"), []attr.Value{supertypes.NewStringValue("issue_owners")})`,
                  ],
                },
                {
                  name: "fallthrough_type",
                  type: "string",
                  description:
                    "The type of fallthrough to apply when choosing to notify issue owners. Only required if the target type is `issue_owners`.",
                  computedOptionalRequired: "optional",
                  enum: `[]string{"AllMembers", "ActiveMembers", "NoOne"}`,
                  validators: [
                    `fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("target_type"), []attr.Value{supertypes.NewStringValue("issue_owners")})`,
                    `fstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("target_type"), []attr.Value{supertypes.NewStringValue("team"), supertypes.NewStringValue("user")})`,
                  ],
                },
              ],
            },
            {
              name: "plugin",
              type: "single_nested",
              description:
                "Send a notification to all legacy integrations (plugins).",
              computedOptionalRequired: "optional",
              attributes: [],
            },
            {
              name: "slack",
              type: "single_nested",
              description: "Notify on Slack.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the Slack integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "channel_name",
                  type: "string",
                  description:
                    "The name of the Slack channel to send the notification to (e.g., #critical, Jane Schmidt).",
                  computedOptionalRequired: "required",
                },
                {
                  name: "channel_id",
                  type: "string",
                  description:
                    "The Slack channel ID to send the notification to. This is an optional field that can be used to avoid rate-limiting.",
                  computedOptionalRequired: "computed_optional",
                },
                {
                  name: "tags",
                  type: "string",
                  description: "A list of tags to show in the notification.",
                  computedOptionalRequired: "computed_optional",
                },
                {
                  name: "notes",
                  type: "string",
                  description:
                    "Text to show alongside the notification. To @ a user, include their user id like `@<USER_ID>`. To include a clickable link, format the link and title like `<http://example.com|Click Here>`.",
                  computedOptionalRequired: "optional",
                },
              ],
            },
            {
              name: "pagerduty",
              type: "single_nested",
              description: "Notify on PagerDuty.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the PagerDuty integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "service_name",
                  type: "string",
                  description:
                    "The name of the service to create the ticket in.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "service_id",
                  type: "string",
                  description: "The ID of the PagerDuty service.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "severity",
                  type: "string",
                  description:
                    "The PagerDuty severity level for the notification.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "discord",
              type: "single_nested",
              description: "Notify on Discord.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the Discord integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "channel_id",
                  type: "string",
                  description:
                    "The ID of the Discord channel to send the notification to.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "tags",
                  type: "string",
                  description: "A list of tags to show in the notification.",
                  computedOptionalRequired: "optional",
                },
              ],
            },
            {
              name: "msteams",
              type: "single_nested",
              description: "Notify on Microsoft Teams.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the Microsoft Teams integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "team_id",
                  type: "string",
                  description:
                    "The integration ID associated with the Microsoft Teams team.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "channel_name",
                  type: "string",
                  description:
                    "The name of the Microsoft Teams channel to send the notification to.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "opsgenie",
              type: "single_nested",
              description: "Notify on OpsGenie.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the OpsGenie integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "team_name",
                  type: "string",
                  description: "The name of the Opsgenie team.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "team_id",
                  type: "string",
                  description:
                    "The ID of the Opsgenie team to send the notification to.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "priority",
                  type: "string",
                  description: "The priority level for the notification.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "vsts",
              type: "single_nested",
              description: "Notify on Azure DevOps.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the OpsGenie integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "project",
                  type: "string",
                  description: "The ID of the Azure DevOps project.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "work_item_type",
                  type: "string",
                  description: "The type of work item to create.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "jira",
              type: "single_nested",
              description: "Create a Jira ticket.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the Jira integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "project",
                  type: "string",
                  description: "The ID of the Jira project.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "issue_type",
                  type: "string",
                  description:
                    "The ID of the type of issue that the ticket should be created as.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "jira_server",
              type: "single_nested",
              description: "Create a Jira Server ticket.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the Jira Server integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "project",
                  type: "string",
                  description: "The ID of the Jira project.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "issue_type",
                  type: "string",
                  description:
                    "The ID of the type of issue that the ticket should be created as.",
                  computedOptionalRequired: "required",
                },
              ],
            },
            {
              name: "github",
              type: "single_nested",
              description: "Create a GitHub issue.",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "integration_id",
                  type: "string",
                  description: "The ID of the GitHub integration.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "repo",
                  type: "string",
                  description:
                    "The name of the repository to create the issue in.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "assignee",
                  type: "string",
                  description: "The GitHub user to assign the issue to.",
                  computedOptionalRequired: "optional",
                },
                {
                  name: "labels",
                  type: "set",
                  description: "A list of labels to assign to the issue.",
                  computedOptionalRequired: "optional",
                  elementType: "string",
                },
              ],
            },
            {
              name: "sentry_app",
              type: "single_nested",
              description: "Trigger an action in a Sentry App (e.g. Rootly).",
              computedOptionalRequired: "optional",
              attributes: [
                {
                  name: "sentry_app_id",
                  type: "string",
                  description:
                    "The numeric Sentry App ID. Use `tostring(data.sentry_app_installation.<name>.sentry_app_id)` to source this value.",
                  computedOptionalRequired: "required",
                },
                {
                  name: "settings",
                  type: "list_nested",
                  description:
                    "Key-value settings passed to the Sentry App action. Specifying `label` preserves the human-readable display name in the Sentry UI for async select fields whose options are paginated by the third-party app.",
                  computedOptionalRequired: "optional",
                  attributes: [
                    {
                      name: "name",
                      type: "string",
                      description: "The name of the setting field.",
                      computedOptionalRequired: "required",
                    },
                    {
                      name: "value",
                      type: "string",
                      description: "The value of the setting field.",
                      computedOptionalRequired: "required",
                    },
                    {
                      name: "label",
                      type: "string",
                      description:
                        "The human-readable display label for the value. Required for async select fields whose option list is paginated by the third-party app — without it the field may appear blank under certain conditions in the Sentry UI after apply.",
                      computedOptionalRequired: "computed_optional",
                    },
                  ],
                },
              ],
            },
          ]),
        },
      ],
    },
  ],
} satisfies Resource;
