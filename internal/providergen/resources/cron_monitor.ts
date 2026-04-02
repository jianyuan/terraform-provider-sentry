import dedent from "dedent";
import type { Resource } from "../schema";

export default {
  name: "cron_monitor",
  description: dedent.withOptions({ trimWhitespace: true })`
      ⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

      Create a Cron Monitor for a Project.
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
      name: "checkin_margin_minutes",
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
      name: "max_runtime_minutes",
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
      description: "The timezone of the cron monitor.",
      computedOptionalRequired: "computed_optional",
      default: `stringdefault.StaticString("UTC")`,
      enum: "sentrydata.Timezones",
    },
  ],
} satisfies Resource;
