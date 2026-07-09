import type { Attribute } from "./schema";

export function tfAttributeDescription(attribute: Attribute) {
  let description = attribute.description;
  if (attribute.deprecationMessage) {
    description += ` **Deprecated** ${attribute.deprecationMessage}`;
  }

  return description;
}
