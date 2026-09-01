package sentrytypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = (*MsTeamsTeamIdType)(nil)

type MsTeamsTeamIdType struct {
	basetypes.StringType
}

func (t MsTeamsTeamIdType) String() string {
	return "sentrytypes.MsTeamsTeamIdType"
}

func (t MsTeamsTeamIdType) ValueType(_ context.Context) attr.Value {
	return MsTeamsTeamId{}
}

func (t MsTeamsTeamIdType) Equal(o attr.Type) bool {
	other, ok := o.(MsTeamsTeamIdType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t MsTeamsTeamIdType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return MsTeamsTeamId{StringValue: in}, nil
}

func (t MsTeamsTeamIdType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
