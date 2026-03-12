import dedent from "dedent";
import type { DataSource, Resource } from "./schema";

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
          m.ProjectSlugs = supertypes.NewSetValueOfSlice(ctx, lo.Map(data, func(item apiclient.Project, _ int) string {
            return item.Slug
          }))
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
            description: "TODO",
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
];
