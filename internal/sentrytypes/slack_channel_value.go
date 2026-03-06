package sentrytypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.StringValuable = (*SlackChannel)(nil)
var _ basetypes.StringValuableWithSemanticEquals = (*SlackChannel)(nil)

type SlackChannel struct {
	basetypes.StringValue
}

func (v SlackChannel) Type(_ context.Context) attr.Type {
	return SlackChannelType{}
}

func (v SlackChannel) Equal(o attr.Value) bool {
	other, ok := o.(SlackChannel)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v SlackChannel) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(SlackChannel)
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

	return strings.TrimPrefix(v.ValueString(), "#") == strings.TrimPrefix(newValue.ValueString(), "#"), diags
}

func SlackChannelNull() SlackChannel {
	return SlackChannel{StringValue: basetypes.NewStringNull()}
}

func SlackChannelUnknown() SlackChannel {
	return SlackChannel{StringValue: basetypes.NewStringUnknown()}
}

func SlackChannelValue(value string) SlackChannel {
	return SlackChannel{StringValue: basetypes.NewStringValue(value)}
}

func SlackChannelPointerValue(value *string) SlackChannel {
	return SlackChannel{StringValue: basetypes.NewStringPointerValue(value)}
}
