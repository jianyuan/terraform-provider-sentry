package sentrytypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.StringValuable = (*TrimmedString)(nil)
var _ basetypes.StringValuableWithSemanticEquals = (*TrimmedString)(nil)

type TrimmedString struct {
	basetypes.StringValue
}

func (v TrimmedString) Type(_ context.Context) attr.Type {
	return TrimmedStringType{}
}

func (v TrimmedString) Equal(o attr.Value) bool {
	other, ok := o.(TrimmedString)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v TrimmedString) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(TrimmedString)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return strings.TrimSpace(v.ValueString()) == strings.TrimSpace(newValue.ValueString()), diags
}

func TrimmedStringNull() TrimmedString {
	return TrimmedString{StringValue: basetypes.NewStringNull()}
}

func TrimmedStringUnknown() TrimmedString {
	return TrimmedString{StringValue: basetypes.NewStringUnknown()}
}

func TrimmedStringValue(value string) TrimmedString {
	return TrimmedString{StringValue: basetypes.NewStringValue(value)}
}

func TrimmedStringPointerValue(value *string) TrimmedString {
	return TrimmedString{StringValue: basetypes.NewStringPointerValue(value)}
}
