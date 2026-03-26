package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &OpHeaderOperandGlobFunction{}

func NewOpHeaderOperandGlobFunction() function.Function {
	return &OpHeaderOperandGlobFunction{}
}

type OpHeaderOperandGlobFunction struct {
}

func (f OpHeaderOperandGlobFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_header_operand_glob"
}

func (f OpHeaderOperandGlobFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "An HTTP header operand that matches a glob pattern. Intended to be used with the `op_header_check` operation.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "pattern",
				MarkdownDescription: "The glob pattern to match.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpHeaderOperandGlobFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var pattern string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &pattern))
	if resp.Error != nil {
		return
	}

	out := map[string]any{
		"header_op": "glob",
		"pattern": map[string]any{
			"value": pattern,
		},
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonOut)))
}
