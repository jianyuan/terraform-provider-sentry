package sentrytypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.StringValuable = (*MsTeamsTeamId)(nil)
var _ basetypes.StringValuableWithSemanticEquals = (*MsTeamsTeamId)(nil)

// MsTeamsTeamId represents the configured team_id of a Microsoft Teams
// notify action. Sentry resolves this value server-side into the
// underlying Microsoft Teams channel's thread ID (e.g.
// "19:xxxxxxxx@thread.tacv2") and returns that instead of echoing back the
// value that was submitted. There is no way to derive one representation
// from the other, so unlike SlackChannel we can't compare them for a
// meaningful difference: any two known values are treated as semantically
// equal so the configured value is retained in state instead of the
// provider raising "produced inconsistent result after apply" or showing a
// permanent diff on every plan.
type MsTeamsTeamId struct {
	basetypes.StringValue
}

func (v MsTeamsTeamId) Type(_ context.Context) attr.Type {
	return MsTeamsTeamIdType{}
}

func (v MsTeamsTeamId) Equal(o attr.Value) bool {
	other, ok := o.(MsTeamsTeamId)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v MsTeamsTeamId) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(MsTeamsTeamId)
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

	if v.IsNull() || v.IsUnknown() || newValue.IsNull() || newValue.IsUnknown() {
		return v.StringValue.Equal(newValue.StringValue), diags
	}

	return true, diags
}

func NewMsTeamsTeamIdNull() MsTeamsTeamId {
	return MsTeamsTeamId{StringValue: basetypes.NewStringNull()}
}

func NewMsTeamsTeamIdUnknown() MsTeamsTeamId {
	return MsTeamsTeamId{StringValue: basetypes.NewStringUnknown()}
}

func NewMsTeamsTeamIdValue(value string) MsTeamsTeamId {
	return MsTeamsTeamId{StringValue: basetypes.NewStringValue(value)}
}

func NewMsTeamsTeamIdPointerValue(value *string) MsTeamsTeamId {
	return MsTeamsTeamId{StringValue: basetypes.NewStringPointerValue(value)}
}
