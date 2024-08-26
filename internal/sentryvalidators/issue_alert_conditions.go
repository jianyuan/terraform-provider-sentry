package sentryvalidators

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = issueAlertConditionsValidator{}

type issueAlertConditionsValidator struct{}

func (issueAlertConditionsValidator) Description(_ context.Context) string {
	return "issue alert conditions must be valid"
}

func (validator issueAlertConditionsValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v issueAlertConditionsValidator) ValidateJSON(value string) (diagnostics diag.Diagnostics) {
	var conditions []map[string]interface{}
	err := json.Unmarshal([]byte(value), &conditions)
	if err != nil {
		diagnostics.AddError(
			"Invalid Issue Alert Conditions",
			fmt.Sprintf("Issue alert conditions must be a valid JSON object, got: %s", err),
		)
		return
	}

	for i, condition := range conditions {
		if _, ok := condition["id"]; !ok {
			diagnostics.AddError(
				"Invalid Issue Alert Conditions",
				fmt.Sprintf("Condition %d must contain an 'id' key", i+1),
			)
		}

		if _, ok := condition["name"]; ok {
			diagnostics.AddError(
				"Invalid Issue Alert Conditions",
				fmt.Sprintf("Condition %d must not contain a 'name' key", i+1),
			)
		}
	}

	return
}

func (v issueAlertConditionsValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	diagnostics := v.ValidateJSON(request.ConfigValue.ValueString())
	for _, diagnostic := range diagnostics {
		response.Diagnostics.Append(diag.WithPath(request.Path, diagnostic))
	}

}

func IssueAlertConditions() validator.String {
	return issueAlertConditionsValidator{}
}
