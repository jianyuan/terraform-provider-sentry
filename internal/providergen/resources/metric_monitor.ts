import dedent from "dedent";
import type { Resource } from "../schema";

export default {
  name: "metric_monitor",
  description: dedent.withOptions({ trimWhitespace: true })`
      ⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

      Create a Metric Monitor for a Project.

      For more information about configuring metric monitors, see [Create a Monitor for a Project](https://docs.sentry.io/api/monitors/create-a-monitor-for-a-project/).
    `,
  api: {
    model: "ProjectMonitor",
    createMethod: "CreateProjectMonitor",
    createRequestAttributes: ["organization", "project"],
    readMethod: "GetProjectMonitor",
    readRequestAttributes: ["organization", "id"],
    updateMethod: "UpdateProjectMonitor",
    updateRequestAttributes: ["organization", "id"],
    deleteMethod: "DeleteProjectMonitor",
    deleteRequestAttributes: ["organization", "id"],
  },
  generate: {
    modelFillers: false,
  },
  importStateAttributes: ["organization", "project", "id"],
  attributes: [
    {
      name: "id",
      type: "string",
      description: "The internal ID of this monitor.",
      computedOptionalRequired: "computed",
      sourceAttribute: ["Id"],
      planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
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
      description: "The project slug or internal ID to create the monitor for.",
      computedOptionalRequired: "required",
      planModifiers: ["stringplanmodifier.RequiresReplace()"],
    },
    {
      name: "enabled",
      type: "bool",
      description: "Whether the monitor is enabled. Defaults to `true`.",
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
      name: "owner",
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
            `stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("team_id"))`,
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
      enum: "sentrydata.SnubaQueryEventTypes",
    },
    {
      name: "query",
      type: "string",
      description:
        "An event search query to subscribe to and monitor for alerts. For example, to filter transactions so that only those with status code 400 are included, you could use `http.status_code:400`.",
      computedOptionalRequired: "optional",
    },
    {
      name: "query_type",
      type: "string",
      description:
        "The type of query. If no value is provided, `query_type` is set to the default for the specified `dataset.`",
      computedOptionalRequired: "computed_optional",
      enum: "sentrydata.SnubaQueryTypes",
    },
    {
      name: "time_window_seconds",
      type: "int",
      description: "The time window in seconds to use for the aggregate query.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "extrapolation_mode",
      type: "string",
      description: "Extrapolation mode to use for the aggregate query.",
      computedOptionalRequired: "computed_optional",
      enum: "sentrydata.SnubaExtrapolationModes",
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
          description:
            "The comparison delta in seconds to use for the aggregate query. Only required for `percent` type.",
          computedOptionalRequired: "optional",
          validators: [
            `fint64validator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), []attr.Value{supertypes.NewStringValue("percent")})`,
            `fint64validator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), []attr.Value{supertypes.NewStringValue("static"), supertypes.NewStringValue("dynamic")})`,
          ],
        },
      ],
    },
    {
      name: "condition_group",
      type: "single_nested",
      description: "Issue detection condition group configuration.",
      computedOptionalRequired: "required",
      attributes: [
        {
          name: "logic_type",
          type: "string",
          description:
            "The logic to apply to the conditions. `any` will evaluate all conditions, and return true if any of those are met. `any-short` will stop evaluating conditions as soon as one is met. `all` will evaluate all conditions, and return true if all of those are met. `none` will return true if none of the conditions are met, will return false immediately if any are met.",
          computedOptionalRequired: "computed_optional",
          default: `stringdefault.StaticString("any")`,
          enum: "sentrydata.DataConditionGroupTypes",
        },
        {
          name: "conditions",
          type: "list_nested",
          description: "Issue detection conditions.",
          computedOptionalRequired: "required",
          attributes: [
            {
              name: "type",
              type: "string",
              description: "The type of condition.",
              computedOptionalRequired: "required",
              enum: "sentrydata.DataConditionTypes",
            },
            {
              name: "comparison",
              type: "int",
              description:
                "The value to compare against. Only required for types other than `anomaly_detection`.",
              computedOptionalRequired: "optional",
              validators: [
                `fint64validator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), []attr.Value{supertypes.NewStringValue("anomaly_detection")})`,
                `fint64validator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), []attr.Value{supertypes.NewStringValue("eq"), supertypes.NewStringValue("gte"), supertypes.NewStringValue("gt"), supertypes.NewStringValue("lte"), supertypes.NewStringValue("lt"), supertypes.NewStringValue("ne")})`,
              ],
            },
            {
              name: "comparison_sensitivity",
              type: "string",
              description:
                "Choose your level of anomaly responsiveness. Higher thresholds means alerts for most anomalies. Lower thresholds means alerts only for larger ones. Only required for `anomaly_detection` type.",
              computedOptionalRequired: "optional",
              enum: "sentrydata.AlertRuleSensitivities",
              validators: [
                `fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), []attr.Value{supertypes.NewStringValue("anomaly_detection")})`,
              ],
            },
            {
              name: "comparison_threshold_type",
              type: "string",
              description:
                "Decide if you want to be alerted to anomalies that are moving above, below, or in both directions in relation to your threshold. Only required for `anomaly_detection` type.",
              computedOptionalRequired: "optional",
              enum: "sentrydata.AlertRuleThresholdTypes",
              validators: [
                `fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), []attr.Value{supertypes.NewStringValue("anomaly_detection")})`,
              ],
            },
            {
              name: "condition_result",
              type: "int",
              description:
                "When the condition is met, the result will be set to this value.",
              computedOptionalRequired: "required",
            },
          ],
        },
      ],
    },
  ],
} satisfies Resource;
