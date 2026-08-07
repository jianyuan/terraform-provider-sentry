package diagutils

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var (
	ErrEmptyResponse = errors.New("empty response")
)

func NewClientError(action string, err error) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Client error", fmt.Sprintf("Unable to %s, got error: %s", action, err))
}

func NewClientStatusError(action string, status int, body []byte) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Client error", fmt.Sprintf("Unable to %s, got status %d: %s", action, status, string(body)))
}

func NewNotFoundError(resource string) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Not found", fmt.Sprintf("No matching %s found", resource))
}

func NewNotSupportedError(action string) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Not supported", fmt.Sprintf("Action %q is not supported", action))
}

func NewFillError(err error) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Fill error", fmt.Sprintf("Unable to fill model: %s", err))
}

func NewImportError(err error) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Import error", fmt.Sprintf("Unable to import: %s", err))
}

func DiagnosticsError(diags diag.Diagnostics) error {
	errList := diags.Errors()
	if len(errList) == 0 {
		return nil
	}

	errs := make([]error, len(errList))
	for i, d := range errList {
		errs[i] = errors.New(FormatDiagnostic(d))
	}

	return errors.Join(errs...)
}

func FormatDiagnostic(d diag.Diagnostic) string {
	msg := d.Summary()

	if detail := d.Detail(); detail != "" {
		msg += "\n\n" + detail
	}

	if p, ok := d.(diag.DiagnosticWithPath); ok {
		msg += "\n" + p.Path().String()
	}

	return msg
}
