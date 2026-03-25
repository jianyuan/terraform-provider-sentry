package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &OpJsonpathOperandLiteralFunction{}

func NewOpJsonpathOperandLiteralFunction() function.Function {
	return &OpJsonpathOperandLiteralFunction{}
}

type OpJsonpathOperandLiteralFunction struct {
}

func (f OpJsonpathOperandLiteralFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_jsonpath_operand_literal"
}

func (f OpJsonpathOperandLiteralFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "A literal value for comparison.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name: "value",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpJsonpathOperandLiteralFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &value))
	if resp.Error != nil {
		return
	}

	out := map[string]any{
		"jsonpath_op": "literal",
		"value":       value,
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonOut)))
}
