package sentrytypes

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestLossyJsonTypeValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            tftypes.Value
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			in: tftypes.Value{},
		},
		"null": {
			in: tftypes.NewValue(tftypes.String, nil),
		},
		"unknown": {
			in: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"valid json object": {
			in: tftypes.NewValue(tftypes.String, `{"hello":"world", "array": [1, 2, 3]}`),
		},
		"valid json array": {
			in: tftypes.NewValue(tftypes.String, `["hello", "world"]`),
		},
		"invalid json - bracket mismatch": {
			in: tftypes.NewValue(tftypes.String, `{"hello":"world"`),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid JSON String Value",
					"A string value was provided that is not valid JSON string format (RFC 7159).\n\n"+
						"Given Value: {\"hello\":\"world\"\n",
				),
			},
		},
		"invalid json - normal string": {
			in: tftypes.NewValue(tftypes.String, "notvalidjson123"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid JSON String Value",
					"A string value was provided that is not valid JSON string format (RFC 7159).\n\n"+
						"Given Value: notvalidjson123\n",
				),
			},
		},
		"wrong-value-type": {
			in: tftypes.NewValue(tftypes.Number, 123),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Lossy JSON Type Validation Error",
					"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
						"expected String value, received tftypes.Value with value: tftypes.Number<\"123\">",
				),
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := LossyJsonType{}.Validate(context.Background(), testCase.in, path.Root("test"))

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestLossyJsonTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in          tftypes.Value
		expectation attr.Value
		expectedErr string
	}{
		"true": {
			in:          tftypes.NewValue(tftypes.String, `{"hello":"world"}`),
			expectation: NewLossyJsonValue(`{"hello":"world"}`),
		},
		"unknown": {
			in:          tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expectation: NewLossyJsonUnknown(),
		},
		"null": {
			in:          tftypes.NewValue(tftypes.String, nil),
			expectation: NewLossyJsonNull(),
		},
		"wrongType": {
			in:          tftypes.NewValue(tftypes.Number, 123),
			expectedErr: "can't unmarshal tftypes.Number into *string, expected string",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			got, err := LossyJsonType{}.ValueFromTerraform(ctx, testCase.in)
			if err != nil {
				if testCase.expectedErr == "" {
					t.Fatalf("Unexpected error: %s", err)
				}
				if testCase.expectedErr != err.Error() {
					t.Fatalf("Expected error to be %q, got %q", testCase.expectedErr, err.Error())
				}
				return
			}
			if err == nil && testCase.expectedErr != "" {
				t.Fatalf("Expected error to be %q, didn't get an error", testCase.expectedErr)
			}
			if !got.Equal(testCase.expectation) {
				t.Errorf("Expected %+v, got %+v", testCase.expectation, got)
			}
			if testCase.expectation.IsNull() != testCase.in.IsNull() {
				t.Errorf("Expected null-ness match: expected %t, got %t", testCase.expectation.IsNull(), testCase.in.IsNull())
			}
			if testCase.expectation.IsUnknown() != !testCase.in.IsKnown() {
				t.Errorf("Expected unknown-ness match: expected %t, got %t", testCase.expectation.IsUnknown(), !testCase.in.IsKnown())
			}
		})
	}
}
