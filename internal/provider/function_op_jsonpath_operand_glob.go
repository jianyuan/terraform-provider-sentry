package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &OpJsonpathOperandGlobFunction{}

func NewOpJsonpathOperandGlobFunction() function.Function {
	return &OpJsonpathOperandGlobFunction{}
}

type OpJsonpathOperandGlobFunction struct {
}

func (f OpJsonpathOperandGlobFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_jsonpath_operand_glob"
}

func (f OpJsonpathOperandGlobFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "A JSONPath operand that matches a glob pattern. Intended to be used with the `op_jsonpath` operation.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "pattern",
				MarkdownDescription: "The glob pattern to match.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpJsonpathOperandGlobFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var pattern string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &pattern))
	if resp.Error != nil {
		return
	}

	out := map[string]any{
		"jsonpath_op": "glob",
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
