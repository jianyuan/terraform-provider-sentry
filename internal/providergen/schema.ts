export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type Attribute =
  | StringAttribute
  | Int64Attribute
  | Float64Attribute
  | BoolAttribute
  | ListAttribute
  | ListNestedAttribute
  | SetAttribute
  | SetNestedAttribute
  | MapAttribute
  | ObjectAttribute
  | SingleNestedAttribute;

export interface BaseAttribute {
  name: string;
  customType?: {
    type: string;
    value: string;
  };
  description: string;
  deprecationMessage?: string;
  computedOptionalRequired: ComputedOptionalRequired;
  sensitive?: boolean;
  default?: string;
  planModifiers?: Array<string>;
  enum?: string;
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

export interface Int64Attribute extends BaseAttribute {
  type: "int64";
}

export interface Float64Attribute extends BaseAttribute {
  type: "float64";
}

export interface BoolAttribute extends BaseAttribute {
  type: "bool";
}

export interface ListAttribute extends BaseAttribute {
  type: "list";
  elementType: "string";
}

export interface ListNestedAttribute extends BaseAttribute {
  type: "list_nested";
  attributes: Array<Attribute>;
  model?: string;
}

export interface SetAttribute extends BaseAttribute {
  type: "set";
  elementType: "string";
}

export interface SetNestedAttribute extends BaseAttribute {
  type: "set_nested";
  attributes: Array<Attribute>;
  model?: string;
}

export interface MapAttribute extends BaseAttribute {
  type: "map";
  elementType: "string";
}

export interface ObjectAttribute extends BaseAttribute {
  type: "object";
  attributes: Array<Attribute>;
}

export interface SingleNestedAttribute extends BaseAttribute {
  type: "single_nested";
  attributes: Array<Attribute>;
  model?: string;
}

export interface BaseDataSourceApiStrategy {
  model: string;
}

export type DataSourceApiStrategy =
  | SimpleDataSourceApiStrategy
  | PaginateDataSourceApiStrategy
  | CustomDataSourceApiStrategy;

export interface SimpleDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  readStrategy: "simple";
  readMethod: string;
  readRequestAttributes?: Array<string>;
}

export interface PaginateDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  readStrategy: "paginate";
  readMethod: string;
  readRequestAttributes?: Array<string>;
  readModel?: string;
  readCursorParam?: string;
  readInitLoop?: string;
  readPreIterate?: string;
  readPostIterate?: string;
}

export interface CustomDataSourceApiStrategy extends BaseDataSourceApiStrategy {
  readStrategy: "custom";
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
  /** Required unless readStrategy is "custom". */
  readMethod?: string;
  readRequestAttributes?: Array<string>;
  readStrategy?: "paginate" | "custom";
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
  generate?: {
    modelFillers?: boolean;
  };
  importStateAttributes?: Array<string>;
  attributes: Array<Attribute>;
}
