package sentrytypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestNumericStringTypeValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in        tftypes.Value
		expectErr bool
	}{
		"empty-struct": {in: tftypes.Value{}},
		"null":         {in: tftypes.NewValue(tftypes.String, nil)},
		"unknown":      {in: tftypes.NewValue(tftypes.String, tftypes.UnknownValue)},
		"numeric":      {in: tftypes.NewValue(tftypes.String, "123")},
		"empty":        {in: tftypes.NewValue(tftypes.String, ""), expectErr: true},
		"letters":      {in: tftypes.NewValue(tftypes.String, "12a"), expectErr: true},
		"plus sign":    {in: tftypes.NewValue(tftypes.String, "+1"), expectErr: true},
		"minus sign":   {in: tftypes.NewValue(tftypes.String, "-1"), expectErr: true},
		"overflow":     {in: tftypes.NewValue(tftypes.String, "99999999999999999999"), expectErr: true},
		"wrong-type":   {in: tftypes.NewValue(tftypes.Number, 123), expectErr: true},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := NumericStringType{}.Validate(context.Background(), testCase.in, path.Root("test"))
			if diags.HasError() != testCase.expectErr {
				t.Fatalf("expected error %v, got %v", testCase.expectErr, diags)
			}
		})
	}
}

func TestNumericStringValueInt64(t *testing.T) {
	t.Parallel()

	got, err := NumericStringValue("4242").ValueInt64()
	if err != nil || got != 4242 {
		t.Fatalf("expected 4242, got %d, %v", got, err)
	}
}
