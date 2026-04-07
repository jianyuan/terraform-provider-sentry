import dedent from "dedent";
import type { DataSource } from "../schema";

export default {
  name: "cron_monitor",
  description: dedent.withOptions({ trimWhitespace: true })`
      ⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

      Retrieve a Cron Monitor.
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
      name: "checkin_margin_minutes",
      type: "int",
      description:
        "Grace period. The number of minutes before a check-in is considered missed.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "failure_issue_threshold",
      type: "int",
      description:
        "Failure tolerance. Create a new issue when this many consecutive missed or error check-ins are processed.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "max_runtime_minutes",
      type: "int",
      description:
        "Maximum runtime. The number of minutes before an in-progress check-in is marked timed out.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "recovery_threshold",
      type: "int",
      description:
        "Recovery Tolerance. Resolve the issue when this many consecutive healthy check-ins are processed. Either `crontab` or `interval_value` and `interval_unit` must be provided.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "schedule",
      type: "single_nested",
      description: "Schedule for the cron monitor.",
      computedOptionalRequired: "computed",
      skipFill: true,
      attributes: [
        {
          name: "crontab",
          type: "string",
          description: "Crontab expression for the cron monitor.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
        {
          name: "interval_value",
          type: "int",
          description: "Interval value for the cron monitor.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
        {
          name: "interval_unit",
          type: "string",
          description: "Interval unit for the cron monitor.",
          computedOptionalRequired: "computed",
          skipFill: true,
        },
      ],
    },
    {
      name: "timezone",
      type: "string",
      description: "The timezone of the cron monitor.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
  ],
} satisfies DataSource;
