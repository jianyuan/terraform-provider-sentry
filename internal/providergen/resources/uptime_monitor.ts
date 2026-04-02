import dedent from "dedent";
import type { Resource } from "../schema";

export default {
  name: "uptime_monitor",
  description: dedent.withOptions({ trimWhitespace: true })`
      ⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

      Create an Uptime Monitor for a Project.

      The \`assertion_json\` argument is a JSON string that represents the assertion to use for the monitor. It is a JSON object with a single key \`root\` whose value is the root operation of the assertion. The assertion is a tree of operations that are evaluated in order. Operations may be constructed using the \`op_\` functions.
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
      name: "url",
      type: "string",
      description: "The URL to monitor.",
      computedOptionalRequired: "required",
    },
    {
      name: "method",
      type: "string",
      description: "The HTTP method to use for the request.",
      computedOptionalRequired: "required",
      enum: "sentrydata.UptimeSubscriptionSupportedHttpMethods",
    },
    {
      name: "body",
      type: "string",
      customType: {
        type: "sentrytypes.TrimmedStringType{}",
        value: "sentrytypes.TrimmedString",
      },
      description:
        "The request body to send. Only applicable for methods that support a body.",
      computedOptionalRequired: "optional",
      nullable: true,
      validators: [
        // https://github.com/getsentry/sentry/blob/master/static/app/views/detectors/components/forms/uptime/detect/index.tsx#L23
        `fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("method"), []attr.Value{supertypes.NewStringValue("GET"), supertypes.NewStringValue("HEAD"), supertypes.NewStringValue("OPTIONS")})`,
      ],
    },
    {
      name: "headers",
      type: "map",
      elementType: "string",
      description: "The headers to send with the request.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "interval_seconds",
      type: "int",
      description: "The amount of time between each uptime check request.",
      computedOptionalRequired: "required",
      enum: "sentrydata.UptimeSubscriptionIntervalSeconds",
    },
    {
      name: "timeout_ms",
      type: "int",
      description: "The request timeout in milliseconds.",
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
        "Number of consecutive successful checks required to mark monitor as recovered. Defaults to `1`.",
      computedOptionalRequired: "computed_optional",
      default: `int64default.StaticInt64(1)`,
    },
    {
      name: "downtime_threshold",
      type: "int",
      description:
        "Number of consecutive failed checks required to mark monitor as down. Defaults to `3`.",
      computedOptionalRequired: "computed_optional",
      default: `int64default.StaticInt64(3)`,
    },
    {
      name: "assertion_json",
      type: "string",
      customType: {
        type: "jsontypes.NormalizedType{}",
        value: "jsontypes.Normalized",
      },
      description:
        "Define conditions that must be met for the check to be considered successful.",
      computedOptionalRequired: "optional",
    },
  ],
} satisfies Resource;
