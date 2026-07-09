import dedent from "dedent";
import type { DataSource } from "../schema";

export default {
  name: "all_projects",
  description: "List of projects in an organization.",
  api: {
    model: "Project",
    readStrategy: "paginate",
    readMethod: "ListOrganizationProjects",
    readRequestAttributes: ["organization"],
  },
  generate: {
    modelFillers: true,
  },
  attributes: [
    {
      name: "organization",
      type: "string",
      description: "The organization slug or internal ID to list projects for.",
      computedOptionalRequired: "required",
      skipFill: true,
    },
    {
      name: "project_slugs",
      type: "set",
      description: "The set of project slugs in this organization.",
      computedOptionalRequired: "computed",
      elementType: "string",
      deprecationMessage: "Use `projects[*].slug` instead.",
      customFill: dedent.withOptions({ trimWhitespace: true })`
          m.ProjectSlugs = supertypes.NewSetValueOfSlice(ctx, lo.Map(data, func(item apiclient.Project, _ int) string {
            return item.Slug
          }))
        `,
    },
    {
      name: "projects",
      type: "set_nested",
      description: "The projects in this organization.",
      computedOptionalRequired: "computed",
      model: "Project",
      sourceAttribute: [],
      attributes: [
        {
          name: "slug",
          type: "string",
          description: "The unique URL slug for the project.",
          computedOptionalRequired: "computed",
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
          name: "color",
          type: "string",
          description: "The color of this project.",
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
      ],
    },
  ],
} satisfies DataSource;
