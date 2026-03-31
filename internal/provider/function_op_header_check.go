package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

var _ function.Function = &OpHeaderCheckFunction{}

func NewOpHeaderCheckFunction() function.Function {
	return &OpHeaderCheckFunction{}
}

type OpHeaderCheckFunction struct {
}

func (f OpHeaderCheckFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_header_check"
}

func (f OpHeaderCheckFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "The HTTP header comparison operation. The operators can be one of `equals`, `not_equal`, `less_than`, `greater_than`, `always`, or `never`. The operands must be constructed using either `op_header_operand_literal` or `op_header_operand_glob`.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "key_operator",
				MarkdownDescription: "The operator to use when matching the header key. Can be one of `equals`, `not_equal`, `less_than`, `greater_than`, `always`, or `never`.",
			},
			function.StringParameter{
				Name:                "key_operand",
				MarkdownDescription: "The header key operand.",
				CustomType:          jsontypes.NormalizedType{},
			},
			function.StringParameter{
				Name:                "value_operator",
				MarkdownDescription: "The operator to use when matching the header value. Can be one of `equals`, `not_equal`, `less_than`, `greater_than`, `always`, or `never`.",
			},
			function.StringParameter{
				Name:                "value_operand",
				MarkdownDescription: "The header value operand.",
				CustomType:          jsontypes.NormalizedType{},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpHeaderCheckFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var keyOperator string
	var keyOperand string
	var valueOperator string
	var valueOperand string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &keyOperator, &keyOperand, &valueOperator, &valueOperand))
	if resp.Error != nil {
		return
	}

	if err := sentrydata.ValidateUptimeAssertionForDefinition("ComparisonType", keyOperator); err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	if err := sentrydata.ValidateJSONUptimeAssertionForDefinition("OpHeaderOperand", []byte(keyOperand)); err != nil {
		resp.Error = function.NewArgumentFuncError(1, err.Error())
		return
	}

	if err := sentrydata.ValidateUptimeAssertionForDefinition("ComparisonType", valueOperator); err != nil {
		resp.Error = function.NewArgumentFuncError(2, err.Error())
		return
	}

	if err := sentrydata.ValidateJSONUptimeAssertionForDefinition("OpHeaderOperand", []byte(valueOperand)); err != nil {
		resp.Error = function.NewArgumentFuncError(3, err.Error())
		return
	}

	var out struct {
		Op          string `json:"op"`
		KeyOperator struct {
			Cmp string `json:"cmp"`
		} `json:"key_op"`
		KeyOperand    json.RawMessage `json:"key_operand"`
		ValueOperator struct {
			Cmp string `json:"cmp"`
		} `json:"value_op"`
		ValueOperand json.RawMessage `json:"value_operand"`
	}
	out.Op = "header_check"
	out.KeyOperator.Cmp = keyOperator
	out.KeyOperand = json.RawMessage(keyOperand)
	out.ValueOperator.Cmp = valueOperator
	out.ValueOperand = json.RawMessage(valueOperand)

	jsonOut, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonOut)))
}
