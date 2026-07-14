import type { DataSource } from "../schema";

export default {
  name: "team",
  description: "Retrieves a Team",
  api: {
    model: "Team",
    readStrategy: "simple",
    readMethod: "GetOrganizationTeam",
    readRequestAttributes: ["organization", "slug"],
  },
  generate: {
    modelFillers: false,
  },
  attributes: [
    {
      name: "organization",
      type: "string",
      description: "The organization slug or internal ID of the organization.",
      computedOptionalRequired: "required",
    },
    {
      name: "slug",
      type: "string",
      description: "The team slug.",
      computedOptionalRequired: "required",
    },
    {
      name: "id",
      type: "string",
      description: "The unique URL slug for this team.",
      deprecationMessage: "Use `slug` instead.",
      computedOptionalRequired: "computed",
    },
    {
      name: "has_access",
      type: "bool",
      description: "Whether the API key user has access to this team.",
      deprecationMessage:
        "This field is deprecated and will be removed in a future version.",
      computedOptionalRequired: "computed",
    },
    {
      name: "internal_id",
      type: "string",
      description: "The internal ID for this team.",
      computedOptionalRequired: "computed",
    },
    {
      name: "is_member",
      type: "bool",
      description: "Whether the API key user is a member of this team.",
      deprecationMessage:
        "This field is deprecated and will be removed in a future version.",
      computedOptionalRequired: "computed",
    },
    {
      name: "is_pending",
      type: "bool",
      description: "Whether the API key user is pending on this team.",
      deprecationMessage:
        "This field is deprecated and will be removed in a future version.",
      computedOptionalRequired: "computed",
    },
    {
      name: "name",
      type: "string",
      description: "The human readable name for this team.",
      computedOptionalRequired: "computed",
    },
  ],
} satisfies DataSource;
