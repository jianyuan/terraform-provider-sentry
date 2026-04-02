import type { Attribute } from "./schema";

export function withExactlyOneAttribute(attributes: Attribute[]): Attribute[] {
  const attributeNames = attributes.map((attribute) => attribute.name);

  return attributes.map((attribute) => {
    return {
      ...attribute,
      validators: [
        ...(attribute.validators ?? []),
        `objectvalidator.ConflictsWith(${attributeNames
          .filter((name) => name !== attribute.name)
          .map((name) => `path.MatchRelative().AtParent().AtName("${name}")`)
          .join(", ")})`,
      ],
    };
  });
}
