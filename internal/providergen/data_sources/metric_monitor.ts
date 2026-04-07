import dedent from "dedent";
import type { DataSource } from "../schema";

export default {
  name: "metric_monitor",
  description: dedent.withOptions({ trimWhitespace: true })`
      ⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

      Retrieve a Metric Monitor for a Project.
    `,
  api: {
    model: "ProjectMonitor",
    readStrategy: "simple",
    readMethod: "GetProjectMonitor",
    readRequestAttributes: ["organization", "id"],
  },
  generate: {
    modelFillers: true,
  },
  attributes: [
    {
      name: "organization",
      type: "string",
      description: "The organization slug or internal ID of the monitor.",
      computedOptionalRequired: "required",
      skipFill: true,
    },
    {
      name: "id",
      type: "string",
      description: "The internal ID of the monitor.",
      computedOptionalRequired: "required",
    },
    {
      name: "project_id",
      type: "string",
      description: "The internal ID of the project this monitor belongs to.",
      computedOptionalRequired: "computed",
    },
    {
      name: "enabled",
      type: "bool",
      description: "Whether the monitor is enabled.",
      computedOptionalRequired: "computed",
    },
    {
      name: "name",
      type: "string",
      description: "The name of the monitor.",
      computedOptionalRequired: "computed",
    },
    {
      name: "description",
      type: "string",
      description: "The description of the monitor.",
      computedOptionalRequired: "computed",
      nullable: true,
    },
    {
      name: "owner",
      type: "single_nested",
      description: "Sentry will assign new issues to this assignee.",
      computedOptionalRequired: "computed",
      skipFill: true,
      attributes: [
        {
          name: "user_id",
          type: "string",
          description: "The user internal ID to assign new issues to.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
        {
          name: "team_id",
          type: "string",
          description: "The team internal ID to assign new issues to.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
      ],
    },
    {
      name: "aggregate",
      type: "string",
      description: "Aggregate query to run on the metric.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "dataset",
      type: "string",
      description: "Dataset to run the aggregate query on.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "environment",
      type: "string",
      description: "Environment to run the aggregate query on.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "event_types",
      type: "set",
      elementType: "string",
      description: "Event types to run the aggregate query on.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "query",
      type: "string",
      description:
        "An event search query to subscribe to and monitor for alerts.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "query_type",
      type: "string",
      description: "The type of query.`",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "time_window_seconds",
      type: "int",
      description: "The time window in seconds to use for the aggregate query.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "extrapolation_mode",
      type: "string",
      description: "Extrapolation mode to use for the aggregate query.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "issue_detection",
      type: "single_nested",
      description: "The issue detection type configuration.",
      computedOptionalRequired: "computed",
      skipFill: true,
      attributes: [
        {
          name: "type",
          type: "string",
          description:
            "`static`: Threshold based monitor; `percent`: Change based monitor; `dynamic`: Dynamic monitor.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
        {
          name: "comparison_delta",
          type: "int",
          description:
            "The comparison delta in seconds to use for the aggregate query. Only available for `percent` type.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
      ],
    },
    {
      name: "condition_group",
      type: "single_nested",
      description: "Issue detection condition group configuration.",
      computedOptionalRequired: "computed",
      skipFill: true,
      attributes: [
        {
          name: "logic_type",
          type: "string",
          description:
            "The logic to apply to the conditions. `any` will evaluate all conditions, and return true if any of those are met. `any-short` will stop evaluating conditions as soon as one is met. `all` will evaluate all conditions, and return true if all of those are met. `none` will return true if none of the conditions are met, will return false immediately if any are met.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
        {
          name: "conditions",
          type: "list_nested",
          description: "Issue detection conditions.",
          computedOptionalRequired: "computed",
          skipFill: true,
          attributes: [
            {
              name: "type",
              type: "string",
              description: "The type of condition.",
              computedOptionalRequired: "computed",
              skipFill: true,
            },
            {
              name: "comparison",
              type: "int",
              description:
                "The value to compare against. Only available for types other than `anomaly_detection`.",
              computedOptionalRequired: "computed",
              skipFill: true,
            },
            {
              name: "comparison_sensitivity",
              type: "string",
              description:
                "Level of anomaly responsiveness. Higher thresholds means alerts for most anomalies. Lower thresholds means alerts only for larger ones. Only available for `anomaly_detection` type.",
              computedOptionalRequired: "computed",
              skipFill: true,
            },
            {
              name: "comparison_threshold_type",
              type: "string",
              description:
                "Alert to anomalies that are moving above, below, or in both directions in relation to your threshold. Only available for `anomaly_detection` type.",
              computedOptionalRequired: "computed",
              skipFill: true,
            },
            {
              name: "condition_result",
              type: "int",
              description:
                "When the condition is met, the result will be set to this value.",
              computedOptionalRequired: "computed",
              skipFill: true,
            },
          ],
        },
      ],
    },
  ],
} satisfies DataSource;
