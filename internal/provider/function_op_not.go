package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

var _ function.Function = &OpNotFunction{}

func NewOpNotFunction() function.Function {
	return &OpNotFunction{}
}

type OpNotFunction struct {
}

func (f OpNotFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_not"
}

func (f OpNotFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Parameters: []function.Parameter{
			&function.StringParameter{
				Name:       "operand",
				CustomType: jsontypes.NormalizedType{},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpNotFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
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
		"op":      "not",
		"operand": json.RawMessage(operand),
	}

	jsonBytes, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonBytes)))
}
