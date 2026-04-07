import dedent from "dedent";
import type { DataSource } from "../schema";

export default {
  name: "uptime_monitor",
  description: dedent.withOptions({ trimWhitespace: true })`
      ⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

      Retrieve an Uptime Monitor for a Project.
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
      name: "url",
      type: "string",
      description: "The URL to monitor.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "method",
      type: "string",
      description: "The HTTP method to use for the request.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "body",
      type: "string",
      description:
        "The request body to send. Only applicable for methods that support a body.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "headers",
      type: "map",
      elementType: "string",
      description: "The headers to send with the request.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "interval_seconds",
      type: "int",
      description: "The amount of time between each uptime check request.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "timeout_ms",
      type: "int",
      description: "The request timeout in milliseconds.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "environment",
      type: "string",
      description: "Name of the environment to create uptime issues in.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "recovery_threshold",
      type: "int",
      description:
        "Number of consecutive successful checks required to mark monitor as recovered.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "downtime_threshold",
      type: "int",
      description:
        "Number of consecutive failed checks required to mark monitor as down.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "assertion_json",
      type: "string",
      description:
        "Define conditions that must be met for the check to be considered successful.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
  ],
} satisfies DataSource;
