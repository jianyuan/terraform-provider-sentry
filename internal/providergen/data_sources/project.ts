import type { DataSource } from "../schema";

export default {
  name: "project",
  description: "Retrieves a project.",
  api: {
    model: "Project",
    readStrategy: "simple",
    readMethod: "GetOrganizationProject",
    readRequestAttributes: ["organization", "slug"],
  },
  generate: {
    modelFillers: true,
  },
  attributes: [
    {
      name: "organization",
      type: "string",
      description: "The organization slug.",
      computedOptionalRequired: "required",
      sourceAttribute: ["Organization", "Slug"],
    },
    {
      name: "slug",
      type: "string",
      description: "The unique URL slug for the project.",
      computedOptionalRequired: "required",
    },
    {
      name: "internal_id",
      type: "string",
      description: "The internal ID of this project.",
      computedOptionalRequired: "computed",
      sourceAttribute: ["Id"],
    },
    {
      name: "name",
      type: "string",
      description: "The name of this project.",
      computedOptionalRequired: "computed",
    },
    {
      name: "platform",
      type: "string",
      description: "The platform of this project.",
      computedOptionalRequired: "computed",
      nullable: true,
    },
    {
      name: "subject_template",
      type: "string",
      description: "The subject template of this project.",
      computedOptionalRequired: "computed",
    },
    {
      name: "color",
      type: "string",
      description: "The color of this project.",
      computedOptionalRequired: "computed",
    },
    {
      name: "is_public",
      type: "bool",
      description: "Whether this project is public.",
      computedOptionalRequired: "computed",
    },
    {
      name: "date_created",
      type: "string",
      description: "The date this project was created.",
      computedOptionalRequired: "computed",
      sourceType: "time",
    },
    {
      name: "features",
      type: "set",
      description: "The features of this project.",
      computedOptionalRequired: "computed",
      elementType: "string",
    },
    {
      name: "teams",
      type: "set_nested",
      description: "The teams of this project.",
      computedOptionalRequired: "computed",
      model: "Team",
      attributes: [
        {
          name: "internal_id",
          type: "string",
          description: "The internal ID of this team.",
          computedOptionalRequired: "computed",
          sourceAttribute: ["Id"],
        },
        {
          name: "name",
          type: "string",
          description: "The name of this team.",
          computedOptionalRequired: "computed",
        },
        {
          name: "slug",
          type: "string",
          description: "The slug of this team.",
          computedOptionalRequired: "computed",
        },
      ],
    },
    {
      name: "id",
      type: "string",
      description: "The unique URL slug for this project.",
      deprecationMessage: "Use `slug` instead.",
      computedOptionalRequired: "computed",
      sourceAttribute: ["Slug"],
    },
  ],
} satisfies DataSource;
