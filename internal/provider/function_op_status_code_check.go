package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

var _ function.Function = &OpStatusCodeCheckFunction{}

func NewOpStatusCodeCheckFunction() function.Function {
	return &OpStatusCodeCheckFunction{}
}

type OpStatusCodeCheckFunction struct {
}

func (f OpStatusCodeCheckFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "op_status_code_check"
}

func (f OpStatusCodeCheckFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "The HTTP status code comparison operation. The `operator` parameter can be one of `equals`, `not_equal`, `less_than`, `greater_than`, `always`, or `never`.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "operator",
				MarkdownDescription: "The comparison operator. Can be one of `equals`, `not_equal`, `less_than`, `greater_than`, `always`, or `never`.",
			},
			function.Int64Parameter{
				Name:                "value",
				MarkdownDescription: "The HTTP status code to compare against.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f OpStatusCodeCheckFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var operator string
	var value int64

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &operator, &value))
	if resp.Error != nil {
		return
	}

	if err := sentrydata.ValidateUptimeAssertionForDefinition("ComparisonType", operator); err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	out := map[string]any{
		"op":       "status_code_check",
		"operator": map[string]any{"cmp": operator},
		"value":    value,
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal assertion: %s", err.Error())))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(jsonOut)))
}
