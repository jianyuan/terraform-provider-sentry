package diagutils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func AddClientError(diags diag.Diagnostics, action string, err error) {
	diags.AddError("Client error", fmt.Sprintf("Unable to %s, got error: %s", action, err))
}

func AddNotFoundError(diags diag.Diagnostics, resource string) {
	diags.AddError("Not found", fmt.Sprintf("No matching %s found", resource))
}

func AddNotSupportedError(diags diag.Diagnostics, action string) {
	diags.AddError("Not supported", fmt.Sprintf("Action %q is not supported", action))
}

func AddFillError(diags diag.Diagnostics, err error) {
	diags.AddError("Fill error", fmt.Sprintf("Unable to fill model: %s", err))
}

func AddImportError(diags diag.Diagnostics, err error) {
	diags.AddError("Import error", fmt.Sprintf("Unable to import: %s", err))
}
