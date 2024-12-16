package sentrytypes

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/go-utils/sliceutils"
)

var _ basetypes.SetValuable = (*StringSet)(nil)

type StringSet struct {
	basetypes.SetValue
}

func (v StringSet) Type(_ context.Context) attr.Type {
	return StringSetType{
		SetType: basetypes.SetType{
			ElemType: basetypes.StringType{},
		},
	}
}

func (v StringSet) Equal(o attr.Value) bool {
	other, ok := o.(StringSet)

	if !ok {
		return false
	}

	return v.SetValue.Equal(other.SetValue)
}

func (v StringSet) ValueString(ctx context.Context) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() || v.IsUnknown() || len(v.Elements()) == 0 {
		return "", diags
	}

	var items []types.String
	diags.Append(v.ElementsAs(ctx, &items, false)...)
	if diags.HasError() {
		return "", diags
	}

	var parsedItems []string
	for _, item := range items {
		if item.IsNull() || item.IsUnknown() {
			continue
		}

		parsedItem := strings.TrimSpace(item.ValueString())
		if len(parsedItem) > 0 {
			parsedItems = append(parsedItems, parsedItem)
		}
	}

	return strings.Join(parsedItems, ","), diags
}

func (v StringSet) ValueStringPointer(ctx context.Context) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() || v.IsUnknown() || len(v.Elements()) == 0 {
		return nil, diags
	}

	var items []types.String
	diags.Append(v.ElementsAs(ctx, &items, false)...)
	if diags.HasError() {
		return nil, diags
	}

	var parsedItems []string
	for _, item := range items {
		if item.IsNull() || item.IsUnknown() {
			continue
		}

		parsedItem := strings.TrimSpace(item.ValueString())
		if len(parsedItem) > 0 {
			parsedItems = append(parsedItems, parsedItem)
		}
	}

	return ptr.Ptr(strings.Join(parsedItems, ",")), diags
}

func StringSetNull() StringSet {
	return StringSet{SetValue: basetypes.NewSetNull(types.StringType)}
}

func StringSetUnknown() StringSet {
	return StringSet{SetValue: basetypes.NewSetUnknown(types.StringType)}
}

func StringSetPointerValue(value *string) (StringSet, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value == nil || strings.TrimSpace(*value) == "" {
		return StringSetNull(), diags
	}

	items := strings.Split(*value, ",")
	elements := sliceutils.Map(func(item string) attr.Value {
		return types.StringValue(strings.TrimSpace(item))
	}, items)

	setValue, d := types.SetValue(types.StringType, elements)
	diags.Append(d...)

	return StringSet{SetValue: setValue}, diags
}
