package sentrytypes

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.StringValuable = (*NumericString)(nil)

type NumericString struct {
	basetypes.StringValue
}

func (v NumericString) Type(_ context.Context) attr.Type {
	return NumericStringType{}
}

func (v NumericString) Equal(o attr.Value) bool {
	other, ok := o.(NumericString)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// ValueInt64 returns the ID as the integer the Sentry API expects.
func (v NumericString) ValueInt64() (int64, error) {
	return parseNumericString(v.ValueString())
}

func NumericStringNull() NumericString {
	return NumericString{StringValue: basetypes.NewStringNull()}
}

func NumericStringUnknown() NumericString {
	return NumericString{StringValue: basetypes.NewStringUnknown()}
}

func NumericStringValue(value string) NumericString {
	return NumericString{StringValue: basetypes.NewStringValue(value)}
}

// parseNumericString accepts digits only, so signs such as "+1" or "-1" that strconv allows are rejected.
func parseNumericString(s string) (int64, error) {
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("%q is not a numeric ID", s)
		}
	}
	return strconv.ParseInt(s, 10, 64)
}
