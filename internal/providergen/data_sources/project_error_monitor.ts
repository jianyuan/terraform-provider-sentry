import dedent from "dedent";
import type { DataSource } from "../schema";

export default {
  name: "project_error_monitor",
  description: dedent.withOptions({ trimWhitespace: true })`
      Retrieve a Project Error Monitor by project ID or slug. This is helpful for managing [default monitors](https://docs.sentry.io/product/new-monitors-and-alerts/monitors/#default-monitors) that were created by Sentry outside of Terraform. You can then map these IDs into \`sentry_alert.monitor_ids\` to define [alert rules](../resources/alert.md) for those monitors.

      **Note:** When multiple monitors are found, the \`first\` attribute can be set to \`true\` to return the first monitor found. If \`first\` is not set to \`true\` and multiple monitors are found, the data source will return an error.
    `,
  api: {
    model: "ProjectMonitor",
    readStrategy: "custom",
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
      name: "project",
      type: "string",
      description: "The project slug or internal ID of the monitor.",
      computedOptionalRequired: "required",
      skipFill: true,
    },
    {
      name: "first",
      type: "bool",
      description: "Return the first monitor found.",
      computedOptionalRequired: "optional",
      skipFill: true,
    },
    {
      name: "id",
      type: "string",
      description: "The internal ID of this monitor.",
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
  ],
} satisfies DataSource;
