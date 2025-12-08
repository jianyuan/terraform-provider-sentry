import { camelize } from "inflection";
import { DATASOURCES, RESOURCES } from "./settings";
import type { DataSource, Attribute, Resource } from "./schema";
import { match, P } from "ts-pattern";
import { parseArgs } from "util";
import dedent from "dedent";

function generateTerraformAttribute({
  parent,
  attribute,
}: {
  parent: string;
  attribute: Attribute;
}) {
  let description = attribute.description;
  if (attribute.deprecationMessage) {
    description += ` **Deprecated** ${attribute.deprecationMessage}`;
  }

  const commonParts: string[] = [];
  commonParts.push(`MarkdownDescription: ${JSON.stringify(description)},`);
  if (attribute.deprecationMessage) {
    commonParts.push(
      `DeprecationMessage: ${JSON.stringify(attribute.deprecationMessage)},`,
    );
  }
  commonParts.push(
    match(attribute.computedOptionalRequired)
      .with("required", () => "Required: true,")
      .with("computed", () => "Computed: true,")
      .with("computed_optional", () => "Optional: true,\nComputed: true,")
      .with("optional", () => "Optional: true,")
      .exhaustive(),
  );
  if (attribute.sensitive) {
    commonParts.push("Sensitive: true,");
  }

  return match(attribute)
    .with({ type: "string" }, () => {
      const parts: string[] = [];
      parts.push("schema.StringAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.StringType{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.String{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.String{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`),
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "int" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.Int64Attribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.Int64Type{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.Int64{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Int64{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`),
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "bool" }, () => {
      const parts: string[] = [];
      parts.push("schema.BoolAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.BoolType{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.Bool{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Bool{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`),
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "list", elementType: "string" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.ListAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.NewListTypeOf[string](ctx),");
      if (attribute.validators) {
        parts.push("Validators: []validator.List{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.List{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`),
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "set", elementType: "string" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.SetAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.NewSetTypeOf[string](ctx),");
      if (attribute.validators) {
        parts.push("Validators: []validator.Set{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Set{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`),
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "set_nested" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.SetNestedAttribute{");
      parts.push(...commonParts);
      parts.push(
        `CustomType: supertypes.NewSetNestedObjectTypeOf[${parent}${camelize(
          attribute.name,
        )}Item](ctx),`,
      );
      parts.push("NestedObject: schema.NestedAttributeObject{");
      parts.push("Attributes: map[string]schema.Attribute{");
      for (const nestedAttribute of attribute.attributes) {
        parts.push(
          `"${nestedAttribute.name}": ${generateTerraformAttribute({
            parent: `${parent}${camelize(attribute.name)}Item`,
            attribute: nestedAttribute,
          })},`,
        );
      }
      parts.push("},");
      parts.push("},");
      parts.push("}");
      return parts.join("\n");
    })
    .exhaustive();
}

function generateTerraformValueType({
  parent,
  attribute,
}: {
  parent: string;
  attribute: Attribute;
}) {
  return match(attribute)
    .with({ type: "string" }, () => "supertypes.StringValue")
    .with({ type: "int" }, () => "supertypes.Int64Value")
    .with({ type: "bool" }, () => "supertypes.BoolValue")
    .with(
      { type: "list", elementType: "string" },
      () => "supertypes.ListValueOf[string]",
    )
    .with(
      { type: "set", elementType: "string" },
      () => "supertypes.SetValueOf[string]",
    )
    .with(
      { type: "set_nested" },
      () =>
        `supertypes.SetNestedObjectValueOf[${parent}${camelize(
          attribute.name,
        )}Item]`,
    )
    .exhaustive();
}

function generateTerraformToPrimitive({
  attribute,
  srcVar,
}: {
  attribute: Attribute;
  srcVar: string;
}) {
  const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
  return match(attribute)
    .with({ type: "string" }, () => `${srcVarName}.ValueString()`)
    .with({ type: "int" }, () => `${srcVarName}.ValueInt64()`)
    .with({ type: "bool" }, () => `${srcVarName}.ValueBool()`)
    .exhaustive();
}

function generatePrimitiveToTerraform({
  name,
  attribute,
  srcVar,
  destVar,
}: {
  name: string;
  attribute: Attribute;
  srcVar: string;
  destVar: string;
}) {
  const srcVarName = [
    srcVar,
    ...(Array.isArray(attribute.sourceAttribute)
      ? attribute.sourceAttribute
      : [camelize(attribute.name)]),
  ].join(".");
  const destVarName = [
    destVar,
    ...(Array.isArray(attribute.destinationAttribute)
      ? attribute.destinationAttribute
      : [camelize(attribute.name)]),
  ].join(".");
  return match(attribute)
    .with(
      { type: "string", nullable: true },
      () => dedent`
      if v, err := ${srcVarName}.Get(); err == nil {
        ${destVarName} = supertypes.NewStringValueOrNull(v)
      } else {
        ${destVarName} = supertypes.NewStringNull()
      }
      `,
    )
    .with(
      { type: "string", sourceType: "time" },
      () =>
        `${destVarName} = supertypes.NewStringValue(${srcVarName}.String())`,
    )
    .with(
      { type: "string" },
      () => `${destVarName} = supertypes.NewStringValue(${srcVarName})`,
    )
    .with(
      { type: "int" },
      () => `${destVarName} = supertypes.NewInt64Value(${srcVarName})`,
    )
    .with(
      { type: "bool" },
      () => `${destVarName} = supertypes.NewBoolValue(${srcVarName})`,
    )
    .with(
      { type: "list", elementType: "string" },
      () =>
        `${destVarName} = supertypes.NewListValueOfSlice(ctx, ${srcVarName})`,
    )
    .with(
      { type: "set", elementType: "string" },
      () =>
        `${destVarName} = supertypes.NewSetValueOfSlice(ctx, ${srcVarName})`,
    )
    .with(
      { type: "set_nested" },
      (attribute) =>
        `${destVarName} = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, sliceutils.Map(func(item apiclient.${attribute.model}) ${name}${camelize(attribute.name)}Item {
          var model ${name}${camelize(attribute.name)}Item
          diags.Append(model.Fill(ctx, item)...)
          return model
        }, ${srcVarName}))`,
    )
    .exhaustive();
}

function generateModel({
  name,
  attributes,
  srcModel,
  generateFillers,
}: {
  name: string;
  attributes: Array<Attribute>;
  srcModel: string;
  generateFillers: boolean;
}) {
  const structLines: string[] = [];
  const fillerLines: string[] = [];
  const extras: string[] = [];

  for (const attribute of attributes) {
    structLines.push(
      `${camelize(attribute.name)} ${generateTerraformValueType({
        parent: name,
        attribute,
      })} \`tfsdk:"${attribute.name}"\``,
    );

    if (!attribute.skipFill) {
      fillerLines.push(
        `${generatePrimitiveToTerraform({
          name,
          attribute,
          srcVar: `data`,
          destVar: `m`,
        })}${attribute.deprecationMessage ? " // Deprecated" : ""}`,
      );
    }

    extras.push(
      ...match(attribute)
        .with({ type: "set_nested" }, (attribute) => [
          generateModel({
            name: `${name}${camelize(attribute.name)}Item`,
            attributes: attribute.attributes,
            srcModel: `apiclient.${attribute.model}`,
            generateFillers,
          }),
        ])
        .otherwise(() => []),
    );
  }

  return `
type ${name} struct {
  ${structLines.join("\n")}
}

${
  generateFillers
    ? dedent`
      func (m *${name}) Fill(ctx context.Context, data ${srcModel}) (diags diag.Diagnostics) {
        ${fillerLines.join("\n")}
        return
      }
    `
    : ""
}

${extras.join("\n\n")}
`;
}

function generateDataSourceModel({ dataSource }: { dataSource: DataSource }) {
  const modelName = `${camelize(dataSource.name)}DataSourceModel`;
  const srcModel = match(dataSource.api)
    .with({ readStrategy: "paginate" }, (api) => `[]apiclient.${api.model}`)
    .with({ readStrategy: "simple" }, (api) => `apiclient.${api.model}`)
    .exhaustive();
  return generateModel({
    name: modelName,
    attributes: dataSource.attributes,
    srcModel,
    generateFillers: dataSource.generate?.modelFillers ?? false,
  });
}

function generateDataSourceSchemaAttributes({
  dataSource,
}: {
  dataSource: DataSource;
}) {
  const lines: string[] = [];

  for (const attribute of dataSource.attributes) {
    lines.push(
      `"${attribute.name}": ${generateTerraformAttribute({
        parent: `${camelize(dataSource.name)}DataSourceModel`,
        attribute,
      })},`,
    );
  }

  return lines.join("\n");
}

function generateDataSource({ dataSource }: { dataSource: DataSource }) {
  console.log(`Generating data source - ${dataSource.name}`);

  const dataSourceName = `${camelize(dataSource.name)}DataSource`;
  const modelName = `${camelize(dataSource.name)}DataSourceModel`;

  const readRequestParams = ["ctx"];
  readRequestParams.push(
    ...match(dataSource.api)
      .with({ readStrategy: "paginate" }, (api) => {
        const parts: string[] = [];
        if (api.readRequestAttributes) {
          parts.push(
            ...api.readRequestAttributes.map((param) => {
              const attribute = dataSource.attributes.find(
                (attribute) => attribute.name === param,
              );
              if (!attribute) {
                throw new Error(
                  `Attribute ${param} not found in data source ${dataSource.name}`,
                );
              }
              return generateTerraformToPrimitive({
                attribute,
                srcVar: "data",
              });
            }),
          );
        }
        parts.push("params");
        return parts;
      })
      .with({ readStrategy: "simple" }, (api) => {
        const parts: string[] = [];
        if (api.readRequestAttributes) {
          parts.push(
            ...api.readRequestAttributes.map((param) => {
              const attribute = dataSource.attributes.find(
                (attribute) => attribute.name === param,
              );
              if (!attribute) {
                throw new Error(
                  `Attribute ${param} not found in data source ${dataSource.name}`,
                );
              }
              return generateTerraformToPrimitive({
                attribute,
                srcVar: "data",
              });
            }),
          );
        }
        return parts;
      })
      .exhaustive(),
  );

  const read = match(dataSource.api)
    .with(
      { readStrategy: "paginate" },
      (api) => `
    var modelInstances []apiclient.${api.readModel ?? api.model}
    params := &apiclient.${api.readMethod}Params{}

    ${api.readInitLoop ?? ""}

    for {
      ${api.readPreIterate ?? ""}

      httpResp, err := d.apiClient.${api.readMethod}WithResponse(${readRequestParams.join(",")})
      if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
        return
      } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
        return
      } else if httpResp.JSON200 == nil {
        resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response body")
        return
      }

      modelInstances = append(modelInstances, *httpResp.JSON200...)

      params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
      if params.Cursor == nil {
        break
      }

      ${api.readPostIterate ?? ""}
    }

    resp.Diagnostics.Append(data.Fill(ctx, modelInstances)...)
    if resp.Diagnostics.HasError() {
      return
    }
    `,
    )
    .with(
      { readStrategy: "simple" },
      (api) => `
    httpResp, err := d.apiClient.${api.readMethod}WithResponse(${readRequestParams.join(",")})
    if err != nil {
      resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
      return
    } else if httpResp.StatusCode() != http.StatusOK {
      resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
      return
    } else if httpResp.JSON200 == nil {
      resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response body")
      return
    }

    resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
    if resp.Diagnostics.HasError() {
      return
    }
    `,
    )
    .exhaustive();

  return `
// Code generated by providergen. DO NOT EDIT.
package provider

import (
  "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &${dataSourceName}{}

func New${dataSourceName}() datasource.DataSource {
  return &${dataSourceName}{}
}

type ${dataSourceName} struct {
  baseDataSource
}

func (d *${dataSourceName}) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_${dataSource.name}"
}

func (d *${dataSourceName}) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
    MarkdownDescription: ${JSON.stringify(dataSource.description)},
    Attributes: map[string]schema.Attribute{
      ${generateDataSourceSchemaAttributes({ dataSource })}
    },
  }
}

func (d *${dataSourceName}) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
  var data ${modelName}

  resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }

  ${read}

  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

${generateDataSourceModel({ dataSource })}
`;
}

function generateResourceModel({ resource }: { resource: Resource }) {
  const modelName = `${camelize(resource.name)}ResourceModel`;
  return generateModel({
    name: modelName,
    attributes: resource.attributes,
  });
}

function generateResourceSchemaAttributes({
  resource,
}: {
  resource: Resource;
}) {
  const lines: string[] = [];

  for (const attribute of resource.attributes) {
    lines.push(
      `"${attribute.name}": ${generateTerraformAttribute({
        parent: `${camelize(resource.name)}ResourceModel`,
        attribute,
      })},`,
    );
  }

  return lines.join("\n");
}

function generateResource({ resource }: { resource: Resource }) {
  console.log(`Generating resource - ${resource.name}`);

  const resourceName = `${camelize(resource.name)}Resource`;
  const modelName = `${camelize(resource.name)}ResourceModel`;

  const createRequestParams = ["ctx"];
  if (resource.api.createRequestAttributes) {
    createRequestParams.push(
      ...resource.api.createRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param,
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`,
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      }),
    );
  }
  createRequestParams.push("body");

  const readRequestParams = ["ctx"];
  if (resource.api.readRequestAttributes) {
    readRequestParams.push(
      ...resource.api.readRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param,
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`,
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      }),
    );
  }
  if (resource.api.readStrategy === "paginate") {
    readRequestParams.push("params");
  }

  const updateRequestParams = ["ctx"];
  if (resource.api.updateRequestAttributes) {
    updateRequestParams.push(
      ...resource.api.updateRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param,
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`,
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      }),
    );
  }
  updateRequestParams.push("body");

  const deleteRequestParams = ["ctx"];
  if (resource.api.deleteRequestAttributes) {
    deleteRequestParams.push(
      ...resource.api.deleteRequestAttributes.map((param) => {
        const attribute = resource.attributes.find(
          (attribute) => attribute.name === param,
        );
        if (!attribute) {
          throw new Error(
            `Attribute ${param} not found in resource ${resource.name}`,
          );
        }
        return generateTerraformToPrimitive({
          attribute,
          srcVar: "data",
        });
      }),
    );
  }

  return `
// Code generated by providergen. DO NOT EDIT.
package provider

import (
  "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &${resourceName}{}
${
  resource.importStateAttributes
    ? `var _ resource.ResourceWithImportState = &${resourceName}{}`
    : ""
}

func New${resourceName}() resource.Resource {
  return &${resourceName}{}
}

type ${resourceName} struct {
  baseResource
}

func (r *${resourceName}) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_${resource.name}"
}

func (r *${resourceName}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
  resp.Schema = schema.Schema{
    MarkdownDescription: ${JSON.stringify(resource.description)},
    Attributes: map[string]schema.Attribute{
      ${generateResourceSchemaAttributes({ resource })}
    },
  }
}

func (r *${resourceName}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
  var data ${modelName}

  resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }

  body, diags := r.getCreateJSONRequestBody(ctx, data)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }

  httpResp, err := r.client.${
    resource.api.createMethod
  }WithResponse(${createRequestParams.join(",")})
  if err != nil {
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got error: %s", err))
    return
  } else if httpResp.StatusCode() != http.StatusOK {
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
    return
  } else if httpResp.JSON200 == nil {
    resp.Diagnostics.AddError("Client Error", "Unable to create, got empty response body")
    return
  }

  resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
  if resp.Diagnostics.HasError() {
    return
  }

  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
  var data ${modelName}

  resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }

  ${match(resource.api)
    .with(
      { readStrategy: "paginate" },
      (api) => dedent`
        var responseData *apiclient.${api.readModel ?? api.model}

        err := retry.Do(
          func() error {
            params := &apiclient.${api.readMethod}Params{
              Limit: ptr.Ptr(int64(100)),
            }

            for {
              httpResp, err := r.client.${
                api.readMethod
              }WithResponse(${readRequestParams.join(",")})
              if err != nil {
                return fmt.Errorf("Unable to read, got error: %s", err)
              } else if httpResp.StatusCode() != http.StatusOK {
                return fmt.Errorf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
              } else if httpResp.JSON200 == nil {
                return fmt.Errorf("Unable to read, got empty response body")
              }

              for _, responseDataItem := range httpResp.JSON200.Data {
                if r.resourceMatch(data, responseDataItem) {
                  responseData = &responseDataItem
                  break
                }
              }

              if v := getBool(httpResp.JSON200.HasMore); !v {
                break
              }

              if v := getString(httpResp.JSON200.${
                api.readCursorParam ?? "LastId"
              }); v != "" {
                params.After = &v
              }
            }

            if responseData == nil {
              return fmt.Errorf("Unable to read, could not find resource in the list")
            }

            return nil
          },
          retry.Delay(5*time.Second),
        )

        if err != nil {
          resp.Diagnostics.AddError("Client Error", err.Error())
          return
        }
      `,
    )
    .otherwise(
      (api) => dedent`
        httpResp, err := r.client.${
          api.readMethod
        }WithResponse(${readRequestParams.join(",")})
        if err != nil {
          resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
          return
        } else if httpResp.StatusCode() != http.StatusOK {
          resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
          return
        } else if httpResp.JSON200 == nil {
          resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response body")
          return
        }

        responseData := httpResp.JSON200
      `,
    )}

  if responseData == nil {
    resp.Diagnostics.AddError("Client Error", "Unable to read, could not find resource in the list")
    return
  }

  resp.Diagnostics.Append(data.Fill(ctx, *responseData)...)
  if resp.Diagnostics.HasError() {
    return
  }

  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
  ${
    resource.api.updateMethod
      ? dedent`
      var data ${modelName}

      resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
      if resp.Diagnostics.HasError() {
        return
      }

      body, diags := r.getUpdateJSONRequestBody(ctx, data)
      resp.Diagnostics.Append(diags...)
      if resp.Diagnostics.HasError() {
        return
      }

      httpResp, err := r.client.${
        resource.api.updateMethod
      }WithResponse(${updateRequestParams.join(",")})
      if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got error: %s", err))
        return
      } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
        return
      } else if httpResp.JSON200 == nil {
        resp.Diagnostics.AddError("Client Error", "Unable to update, got empty response body")
        return
      }

      resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
      if resp.Diagnostics.HasError() {
        return
      }

      resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
      `.trim()
      : dedent`
      resp.Diagnostics.AddError("Not Supported", "Update is not supported for this resource")
      `
  }
}

func (r *${resourceName}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
  ${
    resource.api.deleteMethod
      ? dedent`
      var data ${modelName}

      resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
      if resp.Diagnostics.HasError() {
        return
      }

      httpResp, err := r.client.${
        resource.api.deleteMethod
      }WithResponse(${deleteRequestParams.join(",")})
      if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got error: %s", err))
        return
      } else if httpResp.StatusCode() == http.StatusNotFound {
        return
      } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
        return
      }
      `
      : dedent`
      resp.Diagnostics.AddWarning("Not Supported", "Delete is not supported for this resource. Please manually delete the resource.")
      `
  }
}

${match(resource.importStateAttributes)
  .with([P.any], (attributes) => {
    return `
        func (r *${resourceName}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
          resource.ImportStatePassthroughID(ctx, path.Root("${attributes[0]}"), req, resp)
        }
      `;
  })
  .with([P.any, P.any], (attributes) => {
    return `
        func (r *${resourceName}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
          first, second, err := tfutils.SplitTwoPartId(req.ID, "${attributes[0]}", "${attributes[1]}")
          if err != nil {
            resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
            return
          }

          resp.Diagnostics.Append(resp.State.SetAttribute(
            ctx, path.Root("${attributes[0]}"), first,
          )...)
          resp.Diagnostics.Append(resp.State.SetAttribute(
            ctx, path.Root("${attributes[1]}"), second,
          )...)
        }
      `;
  })
  .with([P.any, P.any, P.any], (attributes) => {
    return `
        func (r *${resourceName}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
          first, second, third, err := tfutils.SplitThreePartId(req.ID, "${attributes[0]}", "${attributes[1]}", "${attributes[2]}")
          if err != nil {
            resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
            return
          }

          resp.Diagnostics.Append(resp.State.SetAttribute(
            ctx, path.Root("${attributes[0]}"), first,
          )...)
          resp.Diagnostics.Append(resp.State.SetAttribute(
            ctx, path.Root("${attributes[1]}"), second,
          )...)
          resp.Diagnostics.Append(resp.State.SetAttribute(
            ctx, path.Root("${attributes[2]}"), third,
          )...)
        }
      `;
  })
  .otherwise(() => "")}

${generateResourceModel({ resource })}
`;
}

function generateProvider({
  resources,
  dataSources,
}: {
  resources: Array<Resource>;
  dataSources: Array<DataSource>;
}) {
  console.log("Generating provider...");
  return `
// Code generated by providergen. DO NOT EDIT.
package provider

var (
	AutoGeneratedResources = []func() resource.Resource{
		${resources
      .sort((a, b) => a.name.localeCompare(b.name))
      .map((resource) => `New${camelize(resource.name)}Resource,`)
      .join("\n")}
	}
	AutoGeneratedDataSources = []func() datasource.DataSource{
		${dataSources
      .sort((a, b) => a.name.localeCompare(b.name))
      .map((dataSource) => `New${camelize(dataSource.name)}DataSource,`)
      .join("\n")}
	}
)
`;
}

async function writeAndFormatGoFile(destination: URL, code: string) {
  await Bun.write(destination, code);
  await Bun.$`go fmt ${destination.pathname}`;
  await Bun.$`go tool goimports -w ${destination.pathname}`;
}

async function main() {
  const { values } = parseArgs({
    args: Bun.argv,
    options: {
      filter: {
        type: "string",
      },
    },
    strict: true,
    allowPositionals: true,
  });

  await Promise.all([
    ...DATASOURCES.map((dataSource) => {
      if (values.filter && values.filter !== dataSource.name) {
        return;
      }

      const code = generateDataSource({ dataSource });
      return writeAndFormatGoFile(
        new URL(
          `../provider/data_source_${dataSource.name}.go`,
          import.meta.url,
        ),
        code,
      );
    }),
    ...RESOURCES.map((resource) => {
      if (values.filter && values.filter !== resource.name) {
        return;
      }

      const code = generateResource({ resource });
      return writeAndFormatGoFile(
        new URL(`../provider/resource_${resource.name}.go`, import.meta.url),
        code,
      );
    }),
    async () => {
      const code = generateProvider({
        resources: RESOURCES,
        dataSources: DATASOURCES,
      });
      await writeAndFormatGoFile(
        new URL(`../provider/provider_gen.go`, import.meta.url),
        code,
      );
    },
  ]);
}

await main()
  .then(() => {
    console.log("✨ Done");
    process.exit(0);
  })
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
