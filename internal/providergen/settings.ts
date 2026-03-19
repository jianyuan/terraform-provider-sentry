import dedent from "dedent";
import type { DataSource, Resource } from "./schema";
import { decapsulate } from "crypto";

export const DATASOURCES: Array<DataSource> = [
  {
    name: "organization",
    description: "Retrieves an organization.",
    api: {
      model: "Organization",
      readStrategy: "simple",
      readMethod: "GetOrganization",
      readRequestAttributes: ["slug"],
    },
    generate: {
      modelFillers: true,
    },
    attributes: [
      {
        name: "slug",
        type: "string",
        description: "The unique URL slug for the organization.",
        computedOptionalRequired: "required",
      },
      {
        name: "internal_id",
        type: "string",
        description: "The internal ID for this organization.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Id"],
      },
      {
        name: "name",
        type: "string",
        description: "The human readable name for this organization.",
        computedOptionalRequired: "computed",
      },
      {
        name: "id",
        type: "string",
        description: "The unique URL slug for this organization.",
        deprecationMessage: "Use `slug` instead.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Slug"],
      },
    ],
  },
  {
    name: "project",
    description: "Retrieves a project.",
    api: {
      model: "Project",
      readStrategy: "simple",
      readMethod: "GetOrganizationProject",
      readRequestAttributes: ["organization", "slug"],
    },
    generate: {
      modelFillers: true,
    },
    attributes: [
      {
        name: "organization",
        type: "string",
        description: "The organization slug.",
        computedOptionalRequired: "required",
        sourceAttribute: ["Organization", "Slug"],
      },
      {
        name: "slug",
        type: "string",
        description: "The unique URL slug for the project.",
        computedOptionalRequired: "required",
      },
      {
        name: "internal_id",
        type: "string",
        description: "The internal ID of this project.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Id"],
      },
      {
        name: "name",
        type: "string",
        description: "The name of this project.",
        computedOptionalRequired: "computed",
      },
      {
        name: "platform",
        type: "string",
        description: "The platform of this project.",
        computedOptionalRequired: "computed",
        nullable: true,
      },
      {
        name: "subject_template",
        type: "string",
        description: "The subject template of this project.",
        computedOptionalRequired: "computed",
      },
      {
        name: "color",
        type: "string",
        description: "The color of this project.",
        computedOptionalRequired: "computed",
      },
      {
        name: "is_public",
        type: "bool",
        description: "Whether this project is public.",
        computedOptionalRequired: "computed",
      },
      {
        name: "date_created",
        type: "string",
        description: "The date this project was created.",
        computedOptionalRequired: "computed",
        sourceType: "time",
      },
      {
        name: "features",
        type: "set",
        description: "The features of this project.",
        computedOptionalRequired: "computed",
        elementType: "string",
      },
      {
        name: "teams",
        type: "set_nested",
        description: "The teams of this project.",
        computedOptionalRequired: "computed",
        model: "Team",
        attributes: [
          {
            name: "internal_id",
            type: "string",
            description: "The internal ID of this team.",
            computedOptionalRequired: "computed",
            sourceAttribute: ["Id"],
          },
          {
            name: "name",
            type: "string",
            description: "The name of this team.",
            computedOptionalRequired: "computed",
          },
          {
            name: "slug",
            type: "string",
            description: "The slug of this team.",
            computedOptionalRequired: "computed",
          },
        ],
      },
      {
        name: "id",
        type: "string",
        description: "The unique URL slug for this project.",
        deprecationMessage: "Use `slug` instead.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Slug"],
      },
    ],
  },
  {
    name: "all_projects",
    description: "List of projects in an organization.",
    api: {
      model: "Project",
      readStrategy: "paginate",
      readMethod: "ListOrganizationProjects",
      readRequestAttributes: ["organization"],
    },
    generate: {
      modelFillers: true,
    },
    attributes: [
      {
        name: "organization",
        type: "string",
        description:
          "The organization slug or internal ID to list projects for.",
        computedOptionalRequired: "required",
        skipFill: true,
      },
      {
        name: "project_slugs",
        type: "set",
        description: "The set of project slugs in this organization.",
        computedOptionalRequired: "computed",
        elementType: "string",
        deprecationMessage: "Use `projects[*].slug` instead.",
        customFill: dedent.withOptions({ trimWhitespace: true })`
          m.ProjectSlugs = supertypes.NewSetValueOfSlice(ctx, sliceutils.Map(func(item apiclient.Project) string {
            return item.Slug
          }, data))
        `,
      },
      {
        name: "projects",
        type: "set_nested",
        description: "The projects in this organization.",
        computedOptionalRequired: "computed",
        model: "Project",
        sourceAttribute: [],
        attributes: [
          {
            name: "slug",
            type: "string",
            description: "The unique URL slug for the project.",
            computedOptionalRequired: "computed",
          },
          {
            name: "internal_id",
            type: "string",
            description: "The internal ID of this project.",
            computedOptionalRequired: "computed",
            sourceAttribute: ["Id"],
          },
          {
            name: "name",
            type: "string",
            description: "The name of this project.",
            computedOptionalRequired: "computed",
          },
          {
            name: "platform",
            type: "string",
            description: "The platform of this project.",
            computedOptionalRequired: "computed",
            nullable: true,
          },
          {
            name: "color",
            type: "string",
            description: "The color of this project.",
            computedOptionalRequired: "computed",
          },
          {
            name: "date_created",
            type: "string",
            description: "The date this project was created.",
            computedOptionalRequired: "computed",
            sourceType: "time",
          },
          {
            name: "features",
            type: "set",
            description: "The features of this project.",
            computedOptionalRequired: "computed",
            elementType: "string",
          },
          {
            name: "teams",
            type: "set_nested",
            description: "The teams of this project.",
            computedOptionalRequired: "computed",
            model: "Team",
            attributes: [
              {
                name: "internal_id",
                type: "string",
                description: "The internal ID of this team.",
                computedOptionalRequired: "computed",
                sourceAttribute: ["Id"],
              },
              {
                name: "name",
                type: "string",
                description: "The name of this team.",
                computedOptionalRequired: "computed",
              },
              {
                name: "slug",
                type: "string",
                description: "The slug of this team.",
                computedOptionalRequired: "computed",
              },
            ],
          },
        ],
      },
    ],
  },
];
export const RESOURCES: Array<Resource> = [
  {
    name: "metric_monitor",
    description: "Create a Metric Monitor for a Project.",
    api: {
      model: "ProjectMonitor",
      createMethod: "CreateProjectMonitor",
      createRequestAttributes: ["organization", "project"],
      readMethod: "GetProjectMonitor",
      readRequestAttributes: ["organization", "id"],
      deleteMethod: "DeleteProjectMonitor",
      deleteRequestAttributes: ["organization", "id"],
    },
    generate: {
      modelFillers: false,
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "The internal ID of this monitor.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Id"],
      },
      {
        name: "organization",
        type: "string",
        description:
          "The organization slug or internal ID to create the monitor for.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "project",
        type: "string",
        description:
          "The project slug or internal ID to create the monitor for.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "enabled",
        type: "bool",
        description: "Whether the monitor is enabled. Defaults to true.",
        computedOptionalRequired: "computed_optional",
        default: `booldefault.StaticBool(true)`,
      },
      {
        name: "name",
        type: "string",
        description: "The name of this monitor.",
        computedOptionalRequired: "required",
      },
      {
        name: "description",
        type: "string",
        description:
          "A description of the monitor. Will be used in the resulting issue.",
        computedOptionalRequired: "optional",
        nullable: true,
      },
      {
        name: "default_assignee",
        type: "single_nested",
        description: "Sentry will assign new issues to this assignee.",
        computedOptionalRequired: "optional",
        nullable: true,
        attributes: [
          {
            name: "user_id",
            type: "string",
            description:
              "The user ID to assign new issues to. Conflicts with `team_id`.",
            computedOptionalRequired: "optional",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("team_id"))`,
            ],
          },
          {
            name: "team_id",
            type: "string",
            description:
              "The team internal ID to assign new issues to. Conflicts with `user_id`.",
            computedOptionalRequired: "optional",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("user_id"))`,
            ],
          },
        ],
      },
      {
        name: "aggregate",
        type: "string",
        description: "Aggregate query to run on the metric.",
        computedOptionalRequired: "required",
      },
      {
        name: "dataset",
        type: "string",
        description: "Dataset to run the aggregate query on.",
        computedOptionalRequired: "required",
        enum: "sentrydata.SnubaDatasets",
      },
      {
        name: "environment",
        type: "string",
        description: "Environment to run the aggregate query on.",
        computedOptionalRequired: "optional",
      },
      {
        name: "event_types",
        type: "set",
        elementType: "string",
        description: "Event types to run the aggregate query on.",
        computedOptionalRequired: "required",
      },
      {
        name: "extrapolation_mode",
        type: "string",
        description: "Extrapolation mode to use for the aggregate query.",
        computedOptionalRequired: "computed_optional",
        enum: "sentrydata.ExtrapolationModes",
      },
      {
        name: "issue_detection",
        type: "single_nested",
        description: "The issue detection type configuration.",
        computedOptionalRequired: "required",
        attributes: [
          {
            name: "type",
            type: "string",
            description:
              "`static`: Threshold based monitor; `percent`: Change based monitor; `dynamic`: Dynamic monitor.",
            computedOptionalRequired: "required",
            enum: "sentrydata.AlertRuleDetectionTypes",
          },
          {
            name: "comparison_delta",
            type: "int",
            description: "TODO",
            computedOptionalRequired: "optional",
          },
        ],
      },
      {
        name: "condition_group",
        type: "single_nested",
        description: "TODO",
        computedOptionalRequired: "required",
        attributes: [
          {
            name: "logic_type",
            type: "string",
            description: "The logic to apply to the conditions.",
            computedOptionalRequired: "computed_optional",
            default: `stringdefault.StaticString("any")`,
            enum: "sentrydata.DataConditionGroupTypes",
          },
          {
            name: "conditions",
            type: "list_nested",
            description: "TODO",
            computedOptionalRequired: "required",
            attributes: [
              {
                name: "type",
                type: "string",
                description: "TODO",
                computedOptionalRequired: "required",
                enum: "sentrydata.DataConditionTypes",
              },
              {
                name: "comparison",
                type: "int",
                description: "TODO",
                computedOptionalRequired: "required",
              },
              {
                name: "condition_result",
                type: "int",
                description: "TODO",
                computedOptionalRequired: "required",
              },
            ],
          },
        ],
      },
    ],
  },
  {
    name: "cron_monitor",
    description: "Create a Cron Monitor for a Project.",
    api: {
      model: "ProjectMonitor",
      createMethod: "CreateProjectMonitor",
      createRequestAttributes: ["organization", "project"],
      readMethod: "GetProjectMonitor",
      readRequestAttributes: ["organization", "id"],
      deleteMethod: "DeleteProjectMonitor",
      deleteRequestAttributes: ["organization", "id"],
    },
    generate: {
      modelFillers: false,
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "The internal ID of this monitor.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Id"],
      },
      {
        name: "organization",
        type: "string",
        description:
          "The organization slug or internal ID to create the monitor for.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "project",
        type: "string",
        description:
          "The project slug or internal ID to create the monitor for.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "enabled",
        type: "bool",
        description: "Whether the monitor is enabled. Defaults to true.",
        computedOptionalRequired: "computed_optional",
        default: `booldefault.StaticBool(true)`,
      },
      {
        name: "name",
        type: "string",
        description: "The name of this monitor.",
        computedOptionalRequired: "required",
      },
      {
        name: "description",
        type: "string",
        description:
          "A description of the monitor. Will be used in the resulting issue.",
        computedOptionalRequired: "optional",
        nullable: true,
      },
      {
        name: "default_assignee",
        type: "single_nested",
        description: "Sentry will assign new issues to this assignee.",
        computedOptionalRequired: "optional",
        nullable: true,
        attributes: [
          {
            name: "user_id",
            type: "string",
            description:
              "The user ID to assign new issues to. Conflicts with `team_id`.",
            computedOptionalRequired: "optional",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("team_id"))`,
            ],
          },
          {
            name: "team_id",
            type: "string",
            description:
              "The team internal ID to assign new issues to. Conflicts with `user_id`.",
            computedOptionalRequired: "optional",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("user_id"))`,
            ],
          },
        ],
      },
      {
        name: "checkin_margin",
        type: "int",
        description:
          "Grace period. The number of minutes before a check-in is considered missed.",
        computedOptionalRequired: "required",
      },
      {
        name: "failure_issue_threshold",
        type: "int",
        description:
          "Failure tolerance. Create a new issue when this many consecutive missed or error check-ins are processed.",
        computedOptionalRequired: "required",
      },
      {
        name: "max_runtime",
        type: "int",
        description:
          "Maximum runtime. The number of minutes before an in-progress check-in is marked timed out.",
        computedOptionalRequired: "required",
      },
      {
        name: "recovery_threshold",
        type: "int",
        description:
          "Recovery Tolerance. Resolve the issue when this many consecutive healthy check-ins are processed. Either `crontab` or `interval_value` and `interval_unit` must be provided.",
        computedOptionalRequired: "required",
      },
      {
        name: "schedule",
        type: "single_nested",
        description: "Set your schedule.",
        computedOptionalRequired: "required",
        attributes: [
          {
            name: "crontab",
            type: "string",
            description:
              "Use the crontab syntax (e.g. `0 0 * * *`). Conflicts with `interval_value` and `interval_unit`.",
            computedOptionalRequired: "optional",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("interval_value"), path.MatchRelative().AtParent().AtName("interval_unit"))`,
              `stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("interval_value"))`,
              `stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("interval_unit"))`,
            ],
          },
          {
            name: "interval_value",
            type: "int",
            description:
              "Interval value. Conflicts with `crontab`. Must be provided with `interval_unit`.",
            computedOptionalRequired: "optional",
            validators: [
              `int64validator.ConflictsWith(path.MatchRelative().AtParent().AtName("crontab"))`,
              `int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("interval_unit"))`,
            ],
          },
          {
            name: "interval_unit",
            type: "string",
            description:
              "Interval unit. Conflicts with `crontab`. Must be provided with `interval_value`.",
            computedOptionalRequired: "optional",
            enum: "sentrydata.Intervals",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("crontab"))`,
              `stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("interval_value"))`,
            ],
          },
        ],
      },
      {
        name: "timezone",
        type: "string",
        description: "Timezone.",
        computedOptionalRequired: "computed_optional",
        default: `stringdefault.StaticString("UTC")`,
        enum: "sentrydata.Timezones",
      },
    ],
  },
  {
    name: "uptime_monitor",
    description: "Create an Uptime Monitor for a Project.",
    api: {
      model: "ProjectMonitor",
      createMethod: "CreateProjectMonitor",
      createRequestAttributes: ["organization", "project"],
      readMethod: "GetProjectMonitor",
      readRequestAttributes: ["organization", "id"],
      deleteMethod: "DeleteProjectMonitor",
      deleteRequestAttributes: ["organization", "id"],
    },
    generate: {
      modelFillers: false,
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "The internal ID of this monitor.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Id"],
      },
      {
        name: "organization",
        type: "string",
        description:
          "The organization slug or internal ID to create the monitor for.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "project",
        type: "string",
        description:
          "The project slug or internal ID to create the monitor for.",
        computedOptionalRequired: "required",
        planModifiers: ["stringplanmodifier.RequiresReplace()"],
      },
      {
        name: "enabled",
        type: "bool",
        description: "Whether the monitor is enabled. Defaults to true.",
        computedOptionalRequired: "computed_optional",
        default: `booldefault.StaticBool(true)`,
      },
      {
        name: "name",
        type: "string",
        description: "The name of this monitor.",
        computedOptionalRequired: "required",
      },
      {
        name: "description",
        type: "string",
        description:
          "A description of the monitor. Will be used in the resulting issue.",
        computedOptionalRequired: "optional",
        nullable: true,
      },
      {
        name: "default_assignee",
        type: "single_nested",
        description: "Sentry will assign new issues to this assignee.",
        computedOptionalRequired: "optional",
        nullable: true,
        attributes: [
          {
            name: "user_id",
            type: "string",
            description:
              "The user ID to assign new issues to. Conflicts with `team_id`.",
            computedOptionalRequired: "optional",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("team_id"))`,
            ],
          },
          {
            name: "team_id",
            type: "string",
            description:
              "The team internal ID to assign new issues to. Conflicts with `user_id`.",
            computedOptionalRequired: "optional",
            validators: [
              `stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("user_id"))`,
            ],
          },
        ],
      },
      {
        name: "url",
        type: "string",
        description: "TODO",
        computedOptionalRequired: "required",
      },
      {
        name: "method",
        type: "string",
        description: "TODO",
        computedOptionalRequired: "required",
        enum: "sentrydata.UptimeSubscriptionSupportedHttpMethods",
      },
      {
        name: "body",
        type: "string",
        description: "TODO",
        computedOptionalRequired: "optional",
        nullable: true,
      },
      {
        name: "headers",
        type: "list_nested",
        description: "TODO",
        computedOptionalRequired: "optional",
        attributes: [
          {
            name: "key",
            type: "string",
            description: "TODO",
            computedOptionalRequired: "required",
          },
          {
            name: "value",
            type: "string",
            description: "TODO",
            computedOptionalRequired: "required",
          },
        ],
      },
      {
        name: "interval_seconds",
        type: "int",
        description: "TODO",
        computedOptionalRequired: "required",
        enum: "sentrydata.UptimeSubscriptionIntervalSeconds",
      },
      {
        name: "timeout_ms",
        type: "int",
        description: "TODO",
        computedOptionalRequired: "required",
      },
      {
        name: "environment",
        type: "string",
        description: "Name of the environment to create uptime issues in.",
        computedOptionalRequired: "required",
      },
      {
        name: "recovery_threshold",
        type: "int",
        description:
          "Number of consecutive successful checks required to mark monitor as recovered. Defaults to 1.",
        computedOptionalRequired: "computed_optional",
        default: `int64default.StaticInt64(1)`,
      },
      {
        name: "downtime_threshold",
        type: "int",
        description:
          "Number of consecutive failed checks required to mark monitor as down. Defaults to 3.",
        computedOptionalRequired: "computed_optional",
        default: `int64default.StaticInt64(3)`,
      },
    ],
  },
  {
    name: "alert",
    description: "Create an Alert for an Organization.",
    api: {
      model: "OrganizationWorkflow",
      createMethod: "CreateOrganizationAlert",
      createRequestAttributes: ["organization"],
      readMethod: "GetOrganizationAlert",
      readRequestAttributes: ["organization", "id"],
      deleteMethod: "DeleteOrganizationAlert",
      deleteRequestAttributes: ["organization", "id"],
    },
    generate: {
      modelFillers: false,
    },
    attributes: [
      {
        name: "id",
        type: "string",
        description: "The internal ID of this alert.",
        computedOptionalRequired: "computed",
        sourceAttribute: ["Id"],
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
        description: "Whether the alert is enabled. Defaults to true.",
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
        type: "set",
        description: "The conditions on which the alert will trigger.",
        computedOptionalRequired: "required",
        elementType: "string",
        enum: "sentrydata.TriggerConditionTypes",
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
            description: "The logic to apply to the conditions.",
            computedOptionalRequired: "required",
            enum: "sentrydata.DataConditionGroupTypes",
          },
          {
            name: "conditions",
            type: "list_nested",
            description: "TODO",
            computedOptionalRequired: "computed_optional",
            attributes: [
              {
                name: "age_comparison",
                type: "single_nested",
                description: "Issue age.",
                computedOptionalRequired: "optional",
                attributes: [
                  {
                    name: "time",
                    type: "string",
                    description: "TODO",
                    computedOptionalRequired: "required",
                    enum: `[]string{"minute", "hour", "day", "week"}`,
                  },
                  {
                    name: "value",
                    type: "int",
                    description: "TODO",
                    computedOptionalRequired: "required",
                    validators: ["int64validator.AtLeast(1)"],
                  },
                  {
                    name: "comparison_type",
                    type: "string",
                    description: "",
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
            ],
          },
          {
            name: "actions",
            type: "list_nested",
            description: "TODO",
            computedOptionalRequired: "required",
            validators: ["listvalidator.SizeAtLeast(1)"],
            attributes: [
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
                    name: "data",
                    type: "map",
                    description:
                      "A list of any fields you want to include in the ticket as objects.",
                    computedOptionalRequired: "optional",
                    elementType: "string",
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
                    name: "data",
                    type: "map",
                    description:
                      "A list of any fields you want to include in the ticket as objects.",
                    computedOptionalRequired: "optional",
                    elementType: "string",
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
                    name: "data",
                    type: "map",
                    description:
                      "A list of any fields you want to include in the ticket as objects.",
                    computedOptionalRequired: "optional",
                    elementType: "string",
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
            ],
          },
        ],
      },
    ],
  },
];
