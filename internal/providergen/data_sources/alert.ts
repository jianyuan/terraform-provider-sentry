import type { DataSource } from "../schema";

export default {
  name: "alert",
  description: "Retrieve an Alert for a Monitor in an Organization.",
  api: {
    model: "OrganizationWorkflow",
    readStrategy: "simple",
    readMethod: "GetOrganizationWorkflow",
    readRequestAttributes: ["organization", "id"],
  },
  generate: {
    modelFillers: true,
  },
  attributes: [
    {
      name: "organization",
      type: "string",
      description: "The organization slug or internal ID of the alert.",
      computedOptionalRequired: "required",
      skipFill: true,
    },
    {
      name: "id",
      type: "string",
      description: "The internal ID of the alert.",
      computedOptionalRequired: "required",
    },

    {
      name: "enabled",
      type: "bool",
      description: "Whether the alert is enabled. Defaults to `true`.",
      computedOptionalRequired: "computed",
    },
    {
      name: "name",
      type: "string",
      description: "The name of this alert.",
      computedOptionalRequired: "computed",
    },
    {
      name: "environment",
      type: "string",
      description: "Name of the environment for this alert.",
      computedOptionalRequired: "computed",
      nullable: true,
    },
    {
      name: "monitor_ids",
      type: "set",
      description: "The IDs of the monitors for this alert.",
      computedOptionalRequired: "computed",
      elementType: "string",
      sourceAttribute: ["DetectorIds"],
    },
    {
      name: "frequency_minutes",
      type: "int",
      description: "How often the alert should fire in minutes.",
      computedOptionalRequired: "computed",
      sourceAttribute: ["Config", "Frequency"],
    },
    {
      name: "triggers_json",
      type: "string",
      description: "The triggers for this alert.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
    {
      name: "action_filters_json",
      type: "string",
      description:
        "The filters to run before the action will fire and the action(s) to fire.",
      computedOptionalRequired: "computed",
      skipFill: true,
    },
  ],
} satisfies DataSource;
