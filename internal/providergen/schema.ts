export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type Attribute =
  | StringAttribute
  | IntAttribute
  | BoolAttribute
  | ListAttribute
  | SetAttribute
  | SetNestedAttribute
  | ObjectAttribute;

export interface BaseAttribute {
  name: string;
  description: string;
  deprecationMessage?: string;
  computedOptionalRequired: ComputedOptionalRequired;
  sensitive?: boolean;
  planModifiers?: Array<string>;
  validators?: Array<string>;
  nullable?: boolean;
  sourceAttribute?: Array<string>;
  sourceType?: "time";
  destinationAttribute?: Array<string>;
  skipFill?: boolean;
  customFill?: string;
}

export interface StringAttribute extends BaseAttribute {
  type: "string";
}

export interface IntAttribute extends BaseAttribute {
  type: "int";
}

export interface BoolAttribute extends BaseAttribute {
  type: "bool";
}

export interface ListAttribute extends BaseAttribute {
  type: "list";
  elementType: "string";
}

export interface SetAttribute extends BaseAttribute {
  type: "set";
  elementType: "string";
}

export interface SetNestedAttribute extends BaseAttribute {
  type: "set_nested";
  attributes: Array<Attribute>;
  model: string;
}

export interface ObjectAttribute extends BaseAttribute {
  type: "object";
  attributes: Array<Attribute>;
}

export interface BaseDataSourceApiStrategy {
  model: string;
  readMethod: string;
  readRequestAttributes?: Array<string>;
}

export type DataSourceApiStrategy =
  | SimpleDataSourceApiStrategy
  | PaginateDataSourceApiStrategy;

export interface SimpleDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  readStrategy: "simple";
}

export interface PaginateDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  readStrategy: "paginate";
  readModel?: string;
  readCursorParam?: string;
  readInitLoop?: string;
  readPreIterate?: string;
  readPostIterate?: string;
}

export interface DataSource {
  name: string;
  description: string;
  api: DataSourceApiStrategy;
  generate?: {
    modelFillers?: boolean;
  };
  attributes: Array<Attribute>;
}

export interface ResourceApiStrategy {
  model?: string;
  createMethod: string;
  createRequestAttributes?: Array<string>;
  readMethod: string;
  readRequestAttributes?: Array<string>;
  readStrategy?: "paginate";
  readModel?: string;
  readCursorParam?: string;
  updateMethod?: string;
  updateRequestAttributes?: Array<string>;
  deleteMethod?: string;
  deleteRequestAttributes?: Array<string>;
}

export interface Resource {
  name: string;
  description: string;
  api: ResourceApiStrategy;
  importStateAttributes?: Array<string>;
  attributes: Array<Attribute>;
}
