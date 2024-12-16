package sentrytypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.SetTypable = (*StringSetType)(nil)

type StringSetType struct {
	basetypes.SetType
}

func (t StringSetType) String() string {
	return "sentrytypes.StringSetType"
}

func (t StringSetType) ValueType(_ context.Context) attr.Value {
	return StringSet{}
}

func (t StringSetType) Equal(o attr.Type) bool {
	other, ok := o.(StringSetType)

	if !ok {
		return false
	}

	return t.SetType.Equal(other.SetType)
}

func (t StringSetType) ValueFromSet(_ context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	return StringSet{SetValue: in}, nil
}

func (t StringSetType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	setValue, ok := attrValue.(basetypes.SetValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := t.ValueFromSet(ctx, setValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}
