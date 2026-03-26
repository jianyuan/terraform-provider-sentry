package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &OpHeaderOperandLiteralFunction{}

func NewOpHeaderOperandLiteralFunction() function.Function {
	return &OpHeaderOperandLiteralFunction{}
}

type OpHeaderOperandLiteralFunction struct {
}

func (f OpHeaderOperandLiteralFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_header_operand_literal"
}

func (f OpHeaderOperandLiteralFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "An HTTP header operand that matches a literal value. Intended to be used with the `op_header_check` operation.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				MarkdownDescription: "The literal value to match.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpHeaderOperandLiteralFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &value))
	if resp.Error != nil {
		return
	}

	out := map[string]any{
		"header_op": "literal",
		"value":     value,
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonOut)))
}
