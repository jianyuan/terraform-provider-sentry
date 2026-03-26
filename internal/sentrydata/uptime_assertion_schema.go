package sentrydata

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jianyuan/go-utils/must"
)

var uptimeAssertionSchema = &jsonschema.Schema{
	Defs: map[string]*jsonschema.Schema{
		"OpType": {
			Type: "string",
			Enum: []any{
				any("and"),
				any("or"),
				any("not"),
				any("status_code_check"),
				any("json_path"),
				any("header_check"),
			},
		},
		"ComparisonType": {
			Type: "string",
			Enum: []any{
				any("equals"),
				any("not_equal"),
				any("less_than"),
				any("greater_than"),
				any("always"),
				any("never"),
			},
		},
		"Comparison": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"cmp"},
			Properties: map[string]*jsonschema.Schema{
				"cmp": {
					Ref: "#/$defs/ComparisonType",
				},
			},
		},
		"Op": {
			OneOf: []*jsonschema.Schema{
				{
					Anchor: "OpAnd",
					Ref:    "#/$defs/OpAnd",
				},
				{
					Anchor: "OpOr",
					Ref:    "#/$defs/OpOr",
				},
				{
					Anchor: "OpNot",
					Ref:    "#/$defs/OpNot",
				},
				{
					Anchor: "OpStatusCode",
					Ref:    "#/$defs/OpStatusCode",
				},
				{
					Anchor: "OpHeaderCheck",
					Ref:    "#/$defs/OpHeaderCheck",
				},
				{
					Anchor: "OpJsonPath",
					Ref:    "#/$defs/OpJsonPath",
				},
			},
		},
		"OpAnd": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"op", "children"},
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Type: "string",
				},
				"op": {
					Const: new(any("and")),
				},
				"children": {
					Type: "array",
					Items: &jsonschema.Schema{
						Ref: "#/$defs/Op",
					},
				},
			},
		},
		"OpOr": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"op", "children"},
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Type: "string",
				},
				"op": {
					Const: new(any("or")),
				},
				"children": {
					Type: "array",
					Items: &jsonschema.Schema{
						Ref: "#/$defs/Op",
					},
				},
			},
		},
		"OpNot": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"op", "operand"},
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Type: "string",
				},
				"op": {
					Const: new(any("not")),
				},
				"operand": {
					Ref: "#/$defs/Op",
				},
			},
		},
		"OpStatusCode": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"op", "operator", "value"},
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Type: "string",
				},
				"op": {
					Const: new(any("status_code_check")),
				},
				"operator": {
					Ref: "#/$defs/Comparison",
				},
				"value": {
					Type: "integer",
				},
			},
		},

		"OpHeaderOperand": {
			OneOf: []*jsonschema.Schema{
				{
					Anchor: "OpHeaderOperandLiteral",
					Ref:    "#/$defs/OpHeaderOperandLiteral",
				},
				{
					Anchor: "OpHeaderOperandGlob",
					Ref:    "#/$defs/OpHeaderOperandGlob",
				},
			},
		},
		"OpHeaderOperandLiteral": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"header_op", "value"},
			Properties: map[string]*jsonschema.Schema{
				"header_op": {
					Const: new(any("literal")),
				},
				"value": {
					Type: "string",
				},
			},
		},
		"OpHeaderOperandGlob": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"header_op", "pattern"},
			Properties: map[string]*jsonschema.Schema{
				"header_op": {
					Const: new(any("glob")),
				},
				"pattern": {
					Type:                 "object",
					AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
					Required:             []string{"value"},
					Properties: map[string]*jsonschema.Schema{
						"value": {
							Type: "string",
						},
					},
				},
			},
		},
		"OpHeaderCheck": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"op", "key_op", "key_operand", "value_op", "value_operand"},
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Type: "string",
				},
				"op": {
					Const: new(any("header_check")),
				},
				"key_op": {
					Ref: "#/$defs/Comparison",
				},
				"key_operand": {
					Ref: "#/$defs/OpHeaderOperand",
				},
				"value_op": {
					Ref: "#/$defs/Comparison",
				},
				"value_operand": {
					Ref: "#/$defs/OpHeaderOperand",
				},
			},
		},

		"OpJsonPathOperand": {
			OneOf: []*jsonschema.Schema{
				{
					Anchor: "OpJsonPathOperandLiteral",
					Ref:    "#/$defs/OpJsonPathOperandLiteral",
				},
				{
					Anchor: "OpJsonPathOperandGlob",
					Ref:    "#/$defs/OpJsonPathOperandGlob",
				},
			},
		},
		"OpJsonPathOperandLiteral": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"jsonpath_op", "value"},
			Properties: map[string]*jsonschema.Schema{
				"jsonpath_op": {
					Const: new(any("literal")),
				},
				"value": {
					Type: "string",
				},
			},
		},
		"OpJsonPathOperandGlob": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"jsonpath_op", "pattern"},
			Properties: map[string]*jsonschema.Schema{
				"jsonpath_op": {
					Const: new(any("glob")),
				},
				"pattern": {
					Type:                 "object",
					AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
					Required:             []string{"value"},
					Properties: map[string]*jsonschema.Schema{
						"value": {
							Type: "string",
						},
					},
				},
			},
		},
		"OpJsonPath": {
			Type:                 "object",
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"op", "operand", "operator", "value"},
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Type: "string",
				},
				"op": {
					Const: new(any("json_path")),
				},
				"operand": {
					Ref: "#/$defs/OpJsonPathOperand",
				},
				"operator": {
					Ref: "#/$defs/Comparison",
				},
				"value": {
					Type: "string",
				},
			},
		},
	},
}

var uptimeSchemaRegistry sync.Map

func GetResolvedUptimeAssertionSchemaForDefinition(def string) (*jsonschema.Resolved, error) {
	resolve, _ := uptimeSchemaRegistry.LoadOrStore(def, sync.OnceValues(func() (*jsonschema.Resolved, error) {
		schema := &jsonschema.Schema{
			Ref: fmt.Sprintf("main.json#/$defs/%s", def),
		}
		return schema.Resolve(&jsonschema.ResolveOptions{
			Loader: func(uri *url.URL) (*jsonschema.Schema, error) {
				if uri.Path == "/main.json" {
					return uptimeAssertionSchema, nil
				}
				return nil, fmt.Errorf("cannot resolve %s", uri)
			},
			ValidateDefaults: true,
		})
	}))

	return resolve.(func() (*jsonschema.Resolved, error))()
}

func MustResolvedUptimeAssertionSchemaForDefinition(def string) *jsonschema.Resolved {
	return must.Get(GetResolvedUptimeAssertionSchemaForDefinition(def))
}

func ValidateUptimeAssertionForDefinition(def string, instance any) error {
	resolved, err := GetResolvedUptimeAssertionSchemaForDefinition(def)
	if err != nil {
		return err
	}

	return resolved.Validate(instance)
}

func ValidateJSONUptimeAssertionForDefinition(def string, instance []byte) error {
	resolved, err := GetResolvedUptimeAssertionSchemaForDefinition(def)
	if err != nil {
		return err
	}

	var instanceJson any
	if err := json.Unmarshal(instance, &instanceJson); err != nil {
		return err
	}

	return resolved.Validate(instanceJson)
}
