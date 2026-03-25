package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

var _ function.Function = &OpAndFunction{}

func NewOpAndFunction() function.Function {
	return &OpAndFunction{}
}

type OpAndFunction struct {
}

func (f OpAndFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_and"
}

func (f OpAndFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		VariadicParameter: &function.StringParameter{
			Name:       "children",
			CustomType: jsontypes.NormalizedType{},
		},
		Return: function.StringReturn{},
	}
}

func (f OpAndFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var children []string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &children))
	if resp.Error != nil {
		return
	}

	var out struct {
		Op       string            `json:"op"`
		Children []json.RawMessage `json:"children"`
	}
	out.Op = "and"

	var argErrors []*function.FuncError

	for i, child := range children {
		if result := sentrydata.UptimeOp.ValidateJSON([]byte(child)); !result.IsValid() {
			argErrors = append(argErrors, function.NewArgumentFuncError(int64(i), sentrydata.CollectEvaluationResultErrors(result).Error()))
		} else {
			out.Children = append(out.Children, json.RawMessage(child))
		}
	}

	if len(argErrors) > 0 {
		resp.Error = function.ConcatFuncErrors(argErrors...)
		return
	}

	jsonBytes, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonBytes)))
}
