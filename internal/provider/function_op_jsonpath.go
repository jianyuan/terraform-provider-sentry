package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

var _ function.Function = &OpJsonpathFunction{}

func NewOpJsonpathFunction() function.Function {
	return &OpJsonpathFunction{}
}

type OpJsonpathFunction struct {
}

func (f OpJsonpathFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_jsonpath"
}

func (f OpJsonpathFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "The JSONPath query operation.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:       "operand",
				CustomType: jsontypes.NormalizedType{},
			},
			function.StringParameter{
				Name: "operator",
			},
			function.StringParameter{
				Name: "value",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpJsonpathFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var operand string
	var operator string
	var value string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &operand, &operator, &value))
	if resp.Error != nil {
		return
	}

	if result := sentrydata.UptimeOpJsonPathOperand.ValidateJSON([]byte(operand)); !result.IsValid() {
		resp.Error = function.NewArgumentFuncError(0, sentrydata.CollectEvaluationResultErrors(result).Error())
		return
	}

	if result := sentrydata.UptimeComparisonType.Validate(operator); !result.IsValid() {
		resp.Error = function.NewArgumentFuncError(1, sentrydata.CollectEvaluationResultErrors(result).Error())
		return
	}

	out := map[string]any{
		"op":      "json_path",
		"operand": json.RawMessage(operand),
		"operator": map[string]string{
			"cmp": operator,
		},
		"value": value,
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonOut)))
}
