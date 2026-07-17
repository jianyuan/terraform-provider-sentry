import type { Resource } from "../schema";

export default {
  name: "organization",
  description: "Sentry Organization resource.",
  api: {
    model: "Organization",
    createMethod: "CreateOrganization",
    readMethod: "GetOrganization",
    readRequestAttributes: ["id"],
    updateMethod: "UpdateOrganization",
    updateRequestAttributes: ["id"],
    deleteMethod: "DeleteOrganization",
    deleteRequestAttributes: ["id"],
  },
  generate: {
    modelFillers: false,
  },
  importStateAttributes: ["id"],
  attributes: [
    {
      name: "id",
      type: "string",
      description: "The unique URL slug for this organization.",
      computedOptionalRequired: "computed",
      planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
    },
    {
      name: "name",
      type: "string",
      description: "The human readable name for the organization.",
      computedOptionalRequired: "required",
    },
    {
      name: "slug",
      type: "string",
      description: "The unique URL slug for this organization.",
      computedOptionalRequired: "computed_optional",
      planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
    },
    {
      name: "agree_terms",
      type: "bool",
      description:
        "You agree to the applicable terms of service and privacy policy. This is only used for creation.",
      computedOptionalRequired: "required",
      planModifiers: ["boolplanmodifier.RequiresReplace()"],
    },
    {
      name: "internal_id",
      type: "string",
      description: "The internal ID for this organization.",
      computedOptionalRequired: "computed",
      planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
    },
    {
      name: "is_early_adopter",
      type: "bool",
      description: "Opt-in to new features before they're released to the public.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "hide_ai_features",
      type: "bool",
      description: "Hide AI features from the organization.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "codecov_access",
      type: "bool",
      description:
        "Enable Code Coverage Insights. This feature is only available for organizations on the Team plan and above.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "default_role",
      type: "string",
      description:
        "The default role new members will receive. Valid values are `member`, `admin`, `manager`, `owner`.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "open_membership",
      type: "bool",
      description: "Allow organization members to freely join any team.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "events_member_admin",
      type: "bool",
      description:
        "Allow members to delete events by granting them the `event:admin` scope.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "alerts_member_write",
      type: "bool",
      description:
        "Allow members to create, edit, and delete alert rules by granting them the `alerts:write` scope.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "attachments_role",
      type: "string",
      description:
        "The role required to download event attachments. Valid values are `member`, `admin`, `manager`, `owner`.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "debug_files_role",
      type: "string",
      description:
        "The role required to download debug information files. Valid values are `member`, `admin`, `manager`, `owner`.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "avatar_type",
      type: "string",
      description:
        "The type of display picture for the organization. Valid values are `letter_avatar`, `upload`.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "avatar",
      type: "string",
      description:
        "The image to upload as the organization avatar, in base64. Required if `avatar_type` is `upload`.",
      computedOptionalRequired: "optional",
      skipFill: true,
    },
    {
      name: "require_2fa",
      type: "bool",
      description: "Require and enforce two-factor authentication for all members.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "allow_shared_issues",
      type: "bool",
      description:
        "Allow sharing of limited details on issues to anonymous users.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "enhanced_privacy",
      type: "bool",
      description:
        "Enable enhanced privacy controls to limit personally identifiable information (PII) as well as source code in things like notifications.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "scrape_javascript",
      type: "bool",
      description:
        "Allow Sentry to scrape missing JavaScript source context when possible.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "store_crash_reports",
      type: "int64",
      description:
        "How many native crash reports to store per issue. Valid values are `0`, `1`, `5`, `10`, `20`, `50`, `100`, `-1` (unlimited).",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "allow_join_requests",
      type: "bool",
      description: "Allow users to request to join your organization.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "data_scrubber",
      type: "bool",
      description: "Require server-side data scrubbing for all projects.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "data_scrubber_defaults",
      type: "bool",
      description:
        "Apply the default scrubbers to prevent things like passwords and credit cards from being stored for all projects.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "sensitive_fields",
      type: "list",
      elementType: "string",
      description:
        "A list of additional global field names to match against when scrubbing data for all projects.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "safe_fields",
      type: "list",
      elementType: "string",
      description: "A list of global field names which data scrubbers should ignore.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "scrub_ip_addresses",
      type: "bool",
      description:
        "Prevent IP addresses from being stored for new events on all projects.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "relay_pii_config",
      type: "string",
      description:
        "Advanced data scrubbing rules that can be configured for each project as a JSON string.",
      computedOptionalRequired: "computed_optional",
      nullable: true,
    },
    {
      name: "trusted_relays",
      type: "list_nested",
      description: "A list of local Relays registered for the organization.",
      computedOptionalRequired: "computed_optional",
      attributes: [
        {
          name: "name",
          type: "string",
          description: "The name of the relay.",
          computedOptionalRequired: "required",
        },
        {
          name: "public_key",
          type: "string",
          description: "The public key of the relay.",
          computedOptionalRequired: "required",
        },
        {
          name: "description",
          type: "string",
          description: "A description for the relay.",
          computedOptionalRequired: "optional",
        },
      ],
    },
    {
      name: "github_pr_bot",
      type: "bool",
      description:
        "Allow Sentry to comment on recent pull requests suspected of causing issues.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "github_open_pr_bot",
      type: "bool",
      description:
        "Allow Sentry to comment on open pull requests to show recent error issues for the code being changed.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "github_nudge_invite",
      type: "bool",
      description:
        "Allow Sentry to detect users committing to your GitHub repositories that are not part of your Sentry organization.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "gitlab_pr_bot",
      type: "bool",
      description:
        "Allow Sentry to comment on recent pull requests suspected of causing issues.",
      computedOptionalRequired: "computed_optional",
    },
    {
      name: "allow_member_project_creation",
      type: "bool",
      description: "Allow members to create projects.",
      computedOptionalRequired: "computed_optional",
    },
  ],
} satisfies Resource;
