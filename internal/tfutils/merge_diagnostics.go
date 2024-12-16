package tfutils

import "github.com/hashicorp/terraform-plugin-framework/diag"

// MergeDiagnostics is a utility function that merges the given diagnostics into the
// provided diagnostics and returns the original value.
func MergeDiagnostics[T any](v T, diagsOut diag.Diagnostics) func(diags *diag.Diagnostics) T {
	return func(diags *diag.Diagnostics) T {
		diags.Append(diagsOut...)
		return v
	}
}
