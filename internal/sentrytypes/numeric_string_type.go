package sentrytypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = (*NumericStringType)(nil)

// NumericStringType is a string attribute holding a numeric ID, such as an integration ID, that the
// Sentry API takes as an integer in requests but returns as a string.
type NumericStringType struct {
	basetypes.StringType
}

func (t NumericStringType) String() string {
	return "sentrytypes.NumericStringType"
}

func (t NumericStringType) ValueType(_ context.Context) attr.Value {
	return NumericString{}
}

func (t NumericStringType) Equal(o attr.Type) bool {
	other, ok := o.(NumericStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t NumericStringType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Numeric String Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var valueString string

	if err := in.As(&valueString); err != nil {
		diags.AddAttributeError(
			path,
			"Numeric String Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)

		return diags
	}

	if _, err := parseNumericString(valueString); err != nil {
		diags.AddAttributeError(
			path,
			"Invalid Numeric ID",
			"A string value was provided that is not a numeric ID.\n\n"+
				"Given Value: "+valueString+"\n",
		)
	}

	return diags
}

func (t NumericStringType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return NumericString{StringValue: in}, nil
}

func (t NumericStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
