package sentrytypes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.StringValuable = (*LossyJson)(nil)
var _ basetypes.StringValuableWithSemanticEquals = (*LossyJson)(nil)

type LossyJson struct {
	basetypes.StringValue
}

func (v LossyJson) Type(_ context.Context) attr.Type {
	return LossyJsonType{}
}

func (v LossyJson) Equal(o attr.Value) bool {
	other, ok := o.(LossyJson)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v LossyJson) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(LossyJson)
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

	result, err := lossyJsonEqual(newValue.ValueString(), v.ValueString())

	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected error occurred while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return false, diags
	}

	return result, diags
}

func lossyJsonEqual(s1, s2 string) (bool, error) {
	v1, err := decodeJson(s1)
	if err != nil {
		return false, err
	}

	v2, err := decodeJson(s2)
	if err != nil {
		return false, err
	}

	return deepLossyEqual(v1, v2), nil
}

func decodeJson(s string) (interface{}, error) {
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()

	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

func deepLossyEqual(v1, v2 interface{}) bool {
	switch v1 := v1.(type) {
	case bool:
		v2, ok := v2.(bool)
		if !ok {
			return false
		}
		return v1 == v2
	case json.Number:
		switch v2 := v2.(type) {
		case json.Number:
			return v1.String() == v2.String()
		case string:
			return v1.String() == v2
		default:
			return false
		}
	case string:
		switch v2 := v2.(type) {
		case json.Number:
			return v1 == v2.String()
		case string:
			return v1 == v2
		default:
			return false
		}
	case []interface{}:
		v2, ok := v2.([]interface{})
		if !ok {
			return false
		}

		if len(v1) != len(v2) {
			return false
		}

		for i := 0; i < len(v1); i++ {
			if !deepLossyEqual(v1[i], v2[i]) {
				return false
			}
		}

		return true
	case map[string]interface{}:
		v2, ok := v2.(map[string]interface{})
		if !ok {
			return false
		}

		if len(v1) > len(v2) {
			return false
		}

		// Check that all keys in v1 are in v2 but not the other way around (lossy)
		for k := range v1 {
			v2Val, ok := v2[k]
			if !ok {
				return false
			}

			if !deepLossyEqual(v1[k], v2Val) {
				return false
			}
		}

		return true
	case nil:
		return v2 == nil
	default:
		panic(fmt.Sprintf("unexpected type %T", v1))
	}
}

func (v LossyJson) Unmarshal(target any) diag.Diagnostics {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("Lossy JSON Unmarshal Error", "json string value is null"))
		return diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("Lossy JSON Unmarshal Error", "json string value is unknown"))
		return diags
	}

	err := json.Unmarshal([]byte(v.ValueString()), target)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("Lossy JSON Unmarshal Error", err.Error()))
	}

	return diags
}

func NewLossyJsonNull() LossyJson {
	return LossyJson{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewLossyJsonUnknown() LossyJson {
	return LossyJson{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewLossyJsonValue(value string) LossyJson {
	return LossyJson{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewLossyJsonPointerValue(value *string) LossyJson {
	return LossyJson{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
