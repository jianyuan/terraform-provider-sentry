package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

var _ function.Function = &AssertionFunction{}

func NewAssertionFunction() function.Function {
	return &AssertionFunction{}
}

type AssertionFunction struct {
}

func (f AssertionFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "assertion"
}

func (f AssertionFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "Creates an uptime assertion. The root operation can be any of the other operation functions that starts with `op_`.",
		Parameters: []function.Parameter{
			&function.StringParameter{
				Name:                "root",
				MarkdownDescription: "The root operation of the assertion.",
				CustomType:          jsontypes.NormalizedType{},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f AssertionFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var operand string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &operand))
	if resp.Error != nil {
		return
	}

	if err := sentrydata.ValidateJSONUptimeAssertionForDefinition("Op", []byte(operand)); err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	out := map[string]any{
		"root": json.RawMessage(operand),
	}

	jsonBytes, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonBytes)))
}
