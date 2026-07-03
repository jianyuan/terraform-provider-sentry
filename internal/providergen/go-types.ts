import { match, P } from "ts-pattern";
import type { Attribute } from "./schema";
import { camelize } from "inflection";

export function primitiveType(attribute: Pick<Attribute, "type">) {
  return match(attribute.type)
    .with("string", () => "string")
    .otherwise(() => {
      throw new Error(`Unsupported primitive type: ${attribute.type}`);
    });
}

export function modelType(attribute: Attribute, parent: string) {
  return match(attribute)
    .with(
      { type: "list_nested" },
      { type: "set_nested" },
      () => `${parent}${camelize(attribute.name)}Item`,
    )
    .with(
      { type: "single_nested" },
      () => `${parent}${camelize(attribute.name)}`,
    )
    .otherwise(() => {
      throw new Error(`Unsupported model type: ${attribute.type}`);
    });
}

export function tfAttributeType(attribute: Attribute, parent: string) {
  return (
    match(attribute)
      .with(
        { customType: { type: P.any } },
        (attribute) => attribute.customType.type,
      )
      .with({ type: "string" }, () => "supertypes.StringType{}")
      .with({ type: "int" }, () => "supertypes.Int64Type{}")
      // Plain basetypes here, not supertypes.Float64Value: the supertypes wrapper inherits
      // basetypes' Float64SemanticEquals without overriding it, so its type assertion rejects
      // the wrapper on every plan (hashicorp/terraform-plugin-framework#786). Float64 is the
      // only type this hits — Int64/String/Bool have no SemanticEquals to inherit.
      .with({ type: "float64" }, () => "basetypes.Float64Type{}")
      .with({ type: "bool" }, () => "supertypes.BoolType{}")
      .with(
        { type: "list" },
        (attribute) =>
          `supertypes.NewListTypeOf[${primitiveType({ type: attribute.elementType })}](ctx)`,
      )
      .with(
        {
          type: "list_nested",
        },
        (attribute) =>
          `supertypes.NewListNestedObjectTypeOf[${modelType(attribute, parent)}](ctx)`,
      )
      .with(
        { type: "set" },
        (attribute) =>
          `supertypes.NewSetTypeOf[${primitiveType({ type: attribute.elementType })}](ctx)`,
      )
      .with(
        { type: "set_nested" },
        (attribute) =>
          `supertypes.NewSetNestedObjectTypeOf[${modelType(attribute, parent)}](ctx)`,
      )
      .with(
        { type: "single_nested" },
        (attribute) =>
          `supertypes.NewSingleNestedObjectTypeOf[${modelType(attribute, parent)}](ctx)`,
      )
      .with(
        { type: "map" },
        (attribute) =>
          `supertypes.NewMapTypeOf[${primitiveType({ type: attribute.elementType })}](ctx)`,
      )
      .otherwise(() => {
        throw new Error(`Unsupported attribute type: ${attribute.type}`);
      })
  );
}

export function tfEnumWrapperFunction(attribute: Attribute) {
  return match(attribute)
    .with({ type: "string" }, () => "tfutils.WithEnumStringAttribute")
    .with({ type: "int" }, () => "tfutils.WithEnumInt64Attribute")
    .with(
      { type: "set", elementType: "string" },
      () => "tfutils.WithEnumSetAttributeStringElements",
    )
    .otherwise(() => null);
}

export function tfAttributeValueType(attribute: Attribute, parent: string) {
  return match(attribute)
    .with(
      { customType: { value: P.any } },
      (attribute) => attribute.customType.value,
    )
    .with({ type: "string" }, () => "supertypes.StringValue")
    .with({ type: "int" }, () => "supertypes.Int64Value")
    .with({ type: "float64" }, () => "types.Float64")
    .with({ type: "bool" }, () => "supertypes.BoolValue")
    .with(
      { type: "list" },
      (attribute) =>
        `supertypes.ListValueOf[${primitiveType({ type: attribute.elementType })}]`,
    )
    .with(
      { type: "list_nested" },
      () =>
        `supertypes.ListNestedObjectValueOf[${modelType(attribute, parent)}]`,
    )
    .with(
      { type: "set" },
      (attribute) =>
        `supertypes.SetValueOf[${primitiveType({ type: attribute.elementType })}]`,
    )
    .with(
      { type: "set_nested" },
      () =>
        `supertypes.SetNestedObjectValueOf[${modelType(attribute, parent)}]`,
    )
    .with(
      { type: "single_nested" },
      () =>
        `supertypes.SingleNestedObjectValueOf[${modelType(attribute, parent)}]`,
    )
    .with(
      { type: "map" },
      (attribute) =>
        `supertypes.MapValueOf[${primitiveType({ type: attribute.elementType })}]`,
    )
    .otherwise(() => {
      throw new Error(`Unsupported attribute type: ${attribute.type}`);
    });
}

export function tfSchemaAttributeType({ type }: { type: Attribute["type"] }) {
  return match(type)
    .with("string", () => "schema.StringAttribute")
    .with("int", () => "schema.Int64Attribute")
    .with("float64", () => "schema.Float64Attribute")
    .with("bool", () => "schema.BoolAttribute")
    .with("list", () => "schema.ListAttribute")
    .with("list_nested", () => "schema.ListNestedAttribute")
    .with("set", () => "schema.SetAttribute")
    .with("set_nested", () => "schema.SetNestedAttribute")
    .with("single_nested", () => "schema.SingleNestedAttribute")
    .with("map", () => "schema.MapAttribute")
    .with("object", () => "schema.ObjectAttribute")
    .exhaustive();
}

export function tfValidatorType({ type }: { type: Attribute["type"] }) {
  return match(type)
    .with("string", () => "validator.String")
    .with("int", () => "validator.Int64")
    .with("float64", () => "validator.Float64")
    .with("bool", () => "validator.Bool")
    .with("list", () => "validator.List")
    .with("list_nested", () => "validator.List")
    .with("set", () => "validator.Set")
    .with("set_nested", () => "validator.Set")
    .with("single_nested", () => "validator.Object")
    .with("map", () => "validator.Map")
    .with("object", () => "validator.Object")
    .exhaustive();
}

export function tfPlanModifierType({ type }: { type: Attribute["type"] }) {
  return match(type)
    .with("string", () => "planmodifier.String")
    .with("int", () => "planmodifier.Int64")
    .with("float64", () => "planmodifier.Float64")
    .with("bool", () => "planmodifier.Bool")
    .with("list", () => "planmodifier.List")
    .with("list_nested", () => "planmodifier.ListNested")
    .with("set", () => "planmodifier.Set")
    .with("set_nested", () => "planmodifier.SetNested")
    .with("single_nested", () => "planmodifier.SingleNested")
    .with("map", () => "planmodifier.Map")
    .with("object", () => "planmodifier.Object")
    .exhaustive();
}
