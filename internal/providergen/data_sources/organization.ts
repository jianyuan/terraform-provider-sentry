import type { DataSource } from "../schema";

export default {
  name: "organization",
  description: "Retrieves an organization.",
  api: {
    model: "Organization",
    readStrategy: "simple",
    readMethod: "GetOrganization",
    readRequestAttributes: ["slug"],
  },
  generate: {
    modelFillers: true,
  },
  attributes: [
    {
      name: "slug",
      type: "string",
      description: "The unique URL slug for the organization.",
      computedOptionalRequired: "required",
    },
    {
      name: "internal_id",
      type: "string",
      description: "The internal ID for this organization.",
      computedOptionalRequired: "computed",
      sourceAttribute: ["Id"],
    },
    {
      name: "name",
      type: "string",
      description: "The human readable name for this organization.",
      computedOptionalRequired: "computed",
    },
    {
      name: "id",
      type: "string",
      description: "The unique URL slug for this organization.",
      deprecationMessage: "Use `slug` instead.",
      computedOptionalRequired: "computed",
      sourceAttribute: ["Slug"],
    },
  ],
} satisfies DataSource;
