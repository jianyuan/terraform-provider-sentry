package sentrydata

import (
	"errors"
	"fmt"
	"strings"

	jsonschemaori "github.com/google/jsonschema-go/jsonschema"
	"github.com/kaptinlin/jsonschema"
	"github.com/samber/lo"
)

var (
	UptimeOpType = jsonschema.Enum(
		any("and"),
		any("or"),
		any("not"),
		any("status_code_check"),
		any("json_path"),
		any("header_check"),
	)
	UptimeComparisonType = jsonschema.Enum(
		any("equals"),
		any("not_equal"),
		any("less_than"),
		any("greater_than"),
		any("always"),
		any("never"),
	)
	UptimeComparison = jsonschema.Object(
		jsonschema.Prop("cmp", UptimeComparisonType),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("cmp"),
	)

	UptimeOp = jsonschema.OneOf(
		UptimeOpStatusCode,
		UptimeOpJsonPath,
		UptimeOpHeaderCheck,
	)

	UptimeOpAnd = jsonschema.Object(
		jsonschema.Prop("id", jsonschema.String()),
		jsonschema.Prop("op", jsonschema.Const("and")),
		jsonschema.Prop("children", jsonschema.Array(jsonschema.Items(UptimeOp))),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("op", "children"),
	)

	UptimeOpStatusCode = jsonschema.Object(
		jsonschema.Prop("id", jsonschema.String()),
		jsonschema.Prop("op", jsonschema.Const("status_code_check")),
		jsonschema.Prop("operator", UptimeComparison),
		jsonschema.Prop("value", jsonschema.Integer()),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("op", "operator", "value"),
	)

	UptimeOpHeaderCheck = jsonschema.Object(
		jsonschema.Prop("id", jsonschema.String()),
		jsonschema.Prop("op", jsonschema.Const("header_check")),
		jsonschema.Prop("key_op", UptimeComparison),
		jsonschema.Prop("key_operand", UptimeOpHeaderOperand),
		jsonschema.Prop("value_op", UptimeComparison),
		jsonschema.Prop("value_operand", UptimeOpHeaderOperand),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("op", "key_op", "value_op"),
	)

	UptimeOpHeaderOperand = jsonschema.OneOf(
		UptimeOpheaderOperandLiteral,
		UptimeOpHeaderOperandGlob,
	)
	UptimeOpheaderOperandLiteral = jsonschema.Object(
		jsonschema.Prop("header_op", jsonschema.Const("literal")),
		jsonschema.Prop("value", jsonschema.String()),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("header_op", "value"),
	)
	UptimeOpHeaderOperandGlob = jsonschema.Object(
		jsonschema.Prop("header_op", jsonschema.Const("glob")),
		jsonschema.Prop("pattern", jsonschema.Object(
			jsonschema.Prop("value", jsonschema.String()),
			jsonschema.AdditionalProps(false),
			jsonschema.Required("value"),
		)),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("header_op", "pattern"),
	)
	UptimeOpJsonPathOperand = jsonschema.OneOf(
		UptimeOpJsonPathOperandLiteral,
		UptimeOpJsonPathOperandGlob,
	)
	UptimeOpJsonPathOperandLiteral = jsonschema.Object(
		jsonschema.Prop("jsonpath_op", jsonschema.Const("literal")),
		jsonschema.Prop("value", jsonschema.String()),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("jsonpath_op", "value"),
	)
	UptimeOpJsonPathOperandGlob = jsonschema.Object(
		jsonschema.Prop("jsonpath_op", jsonschema.Const("glob")),
		jsonschema.Prop("pattern", jsonschema.Object(
			jsonschema.Prop("value", jsonschema.String()),
			jsonschema.AdditionalProps(false),
			jsonschema.Required("value"),
		)),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("jsonpath_op", "pattern"),
	)
	UptimeOpJsonPath = jsonschema.Object(
		jsonschema.Prop("id", jsonschema.String()),
		jsonschema.Prop("op", jsonschema.Const("json_path")),
		jsonschema.Prop("operand", UptimeOpJsonPathOperand),
		jsonschema.Prop("operator", UptimeComparison),
		jsonschema.Prop("value", jsonschema.String()),
		jsonschema.AdditionalProps(false),
		jsonschema.Required("op", "operand", "operator", "value"),
	)

	UptimeAssertionSchema = &jsonschemaori.Schema{
		Definitions: map[string]*jsonschemaori.Schema{
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
				AdditionalProperties: &jsonschemaori.Schema{Not: &jsonschemaori.Schema{}},
				Required:             []string{"cmp"},
				Properties: map[string]*jsonschemaori.Schema{
					"cmp": {
						Ref: "#/definitions/ComparisonType",
					},
				},
			},
			"Op": {
				OneOf: []*jsonschemaori.Schema{
					{
						Ref: "#/definitions/OpStatusCode",
					},
					{
						Ref: "#/definitions/OpJsonPath",
					},
					// {
					// 	Ref: "#/definitions/OpHeaderCheck",
					// },
				},
			},
			"OpAnd": {
				Type:                 "object",
				AdditionalProperties: &jsonschemaori.Schema{Not: &jsonschemaori.Schema{}},
				Required:             []string{"op", "children"},
				Properties: map[string]*jsonschemaori.Schema{
					"op": {
						Const: jsonschemaori.Ptr(any("status_code_check")),
					},
					"children": {
						Type: "array",
						Items: &jsonschemaori.Schema{
							Ref: "#/definitions/Op",
						},
					},
				},
			},
			// "OpOr": {},
			// "OpNot": {},
			"OpStatusCode": {
				Type:                 "object",
				AdditionalProperties: &jsonschemaori.Schema{Not: &jsonschemaori.Schema{}},
				Required:             []string{"op", "operator", "value"},
				Properties: map[string]*jsonschemaori.Schema{
					"id": {
						Type: "string",
					},
					"op": {
						Const: jsonschemaori.Ptr(any("status_code_check")),
					},
					"operator": {
						Ref: "#/definitions/Comparison",
					},
					"value": {
						Type: "integer",
					},
				},
			},
			"OpJsonPathOperand": {
				OneOf: []*jsonschemaori.Schema{
					{
						Anchor: "OpJsonPathOperandLiteral",
						Ref:    "#/definitions/OpJsonPathOperandLiteral",
					},
					{
						Anchor: "OpJsonPathOperandGlob",
						Ref:    "#/definitions/OpJsonPathOperandGlob",
					},
				},
			},
			"OpJsonPathOperandLiteral": {
				Type:                 "object",
				AdditionalProperties: &jsonschemaori.Schema{Not: &jsonschemaori.Schema{}},
				Required:             []string{"jsonpath_op", "value"},
				Properties: map[string]*jsonschemaori.Schema{
					"jsonpath_op": {
						Const: jsonschemaori.Ptr(any("literal")),
					},
					"value": {
						Type: "string",
					},
				},
			},
			"OpJsonPathOperandGlob": {
				Type:                 "object",
				AdditionalProperties: &jsonschemaori.Schema{Not: &jsonschemaori.Schema{}},
				Required:             []string{"jsonpath_op", "pattern"},
				Properties: map[string]*jsonschemaori.Schema{
					"jsonpath_op": {
						Const: jsonschemaori.Ptr(any("glob")),
					},
					"pattern": {
						Type:                 "Object",
						AdditionalProperties: &jsonschemaori.Schema{Not: &jsonschemaori.Schema{}},
						Required:             []string{"value"},
						Properties: map[string]*jsonschemaori.Schema{
							"value": {
								Type: "string",
							},
						},
					},
				},
			},
			"OpJsonPath": {
				Type:                 "object",
				AdditionalProperties: &jsonschemaori.Schema{Not: &jsonschemaori.Schema{}},
				Required:             []string{"op", "operand", "operator", "value"},
				Properties: map[string]*jsonschemaori.Schema{
					"id": {
						Type: "string",
					},
					"op": {
						Const: jsonschemaori.Ptr(any("json_path")),
					},
					"operand": {
						Type: "string",
					},
					"operator": {
						Ref: "#/definitions/Comparison",
					},
					"value": {
						Type: "string",
					},
				},
			},
		},
	}
)

func CollectEvaluationResultErrors(result *jsonschema.EvaluationResult) error {
	localizedErrors := result.ToLocalizeList(nil)
	return errors.New(strings.Join(lo.MapToSlice(localizedErrors.Errors, func(key string, value string) string {
		return fmt.Sprintf("%s: %s", key, value)
	}), "\n"))
}
