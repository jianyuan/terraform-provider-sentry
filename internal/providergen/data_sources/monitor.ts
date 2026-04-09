import type { DataSource } from "../schema";

export default {
  name: "monitor",
  description: "Retrieve a monitor by ID, or by project and type.",
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
      description:
        "The project slug or internal ID of the monitor. Conflicts with `id`.",
      computedOptionalRequired: "optional",
      skipFill: true,
      validators: [`stringvalidator.ConflictsWith(path.MatchRoot("id"))`],
    },
    {
      name: "type",
      type: "string",
      description: "Monitor type to retrieve.",
      computedOptionalRequired: "optional",
      validators: [`stringvalidator.ConflictsWith(path.MatchRoot("id"))`],
    },
    {
      name: "first",
      type: "bool",
      description: "Return the first monitor found.",
      computedOptionalRequired: "optional",
      skipFill: true,
      validators: [`boolvalidator.ConflictsWith(path.MatchRoot("id"))`],
    },
    {
      name: "id",
      type: "string",
      description:
        "The internal ID of the monitor to retrieve. Conflicts with `project`, `type`, and `first`.",
      computedOptionalRequired: "optional",
      validators: [
        `stringvalidator.ConflictsWith(path.MatchRoot("project"), path.MatchRoot("type"), path.MatchRoot("first"))`,
      ],
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
