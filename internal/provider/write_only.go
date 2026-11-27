package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-plugin-framework-utils/fwtypes"
)

// Some secrets can be set either through a regular attribute `<attr>`, which is
// stored in the Terraform state, or through a write-only attribute `<attr>_wo`,
// which is never stored. Write-only values are only available from the
// configuration, so `<attr>_wo_version` lets the user signal that the value has
// changed and needs to be sent again.

// privateState is implemented by the private state data in resource requests
// and responses.
type privateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// writeOnlyPrivateStateKey is the private state key that records that the
// write-only form of attr was sent to Sentry.
func writeOnlyPrivateStateKey(attr string) string {
	return attr + "_wo"
}

// writeOnlyInUse reports whether the write-only form of attr was sent to Sentry
// by a previous create or update.
func writeOnlyInUse(ctx context.Context, private privateState, attr string) (bool, diag.Diagnostics) {
	v, diags := private.GetKey(ctx, writeOnlyPrivateStateKey(attr))
	return len(v) > 0, diags
}

// setWriteOnlyInUse records in private state whether the write-only form of
// attr is in use.
func setWriteOnlyInUse(ctx context.Context, private privateState, attr string, inUse bool) diag.Diagnostics {
	var value []byte
	if inUse {
		value = []byte("true")
	}
	return private.SetKey(ctx, writeOnlyPrivateStateKey(attr), value)
}

// writeOnlySecret is the outcome of resolveWriteOnlySecret.
type writeOnlySecret struct {
	// Value is the secret to send to Sentry, or nil if there is nothing new to send.
	Value *string
	// Keep is true when the write-only form is in use but its value has not
	// changed, so the value Sentry already has should be kept.
	Keep bool
	// InUse is true when the write-only form is in use after this operation.
	InUse bool
}

// resolveWriteOnlySecret decides what to send to Sentry for a secret that can
// be set through either plain or writeOnly (from the configuration).
//
// A write-only value is sent when it was not sent before (wasInUse is false,
// as on create, after import, or when switching from the plain attribute) or
// when its version changed. Otherwise the value stored in Sentry is kept.
func resolveWriteOnlySecret(plain, writeOnly types.String, versionChanged, wasInUse bool) writeOnlySecret {
	switch {
	case fwtypes.IsKnown(plain):
		return writeOnlySecret{Value: plain.ValueStringPointer()}
	case !fwtypes.IsKnown(writeOnly):
		return writeOnlySecret{}
	case versionChanged || !wasInUse:
		return writeOnlySecret{Value: writeOnly.ValueStringPointer(), InUse: true}
	default:
		return writeOnlySecret{Keep: true, InUse: true}
	}
}

// writeOnlyStringAttribute returns the schema for `<attr>_wo`.
func writeOnlyStringAttribute(attr string, description string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: fmt.Sprintf(
			"%[1]s This value is write-only and is not stored in the Terraform state; it requires Terraform 1.11 or later. "+
				"It is sent on create, when switching from `%[2]s`, and whenever `%[2]s_wo_version` changes.",
			description, attr,
		),
		Optional:  true,
		Sensitive: true,
		WriteOnly: true,
	}
}

// writeOnlyVersionAttribute returns the schema for `<attr>_wo_version`.
func writeOnlyVersionAttribute(attr string) schema.Int64Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: fmt.Sprintf(
			"Version of `%[1]s_wo`. Terraform cannot detect changes to a write-only value, so change this number to send a new `%[1]s_wo` value to Sentry.",
			attr,
		),
		Optional: true,
		Validators: []validator.Int64{
			int64validator.AlsoRequires(path.MatchRoot(attr + "_wo")),
		},
	}
}

// resolveWriteOnlySecretFromConfig reads `<attr>_wo` from the configuration and
// resolves what to send for attr. plain and version are the planned values of
// `<attr>` and `<attr>_wo_version`. priorVersion is the value of
// `<attr>_wo_version` in the prior state, or nil on create.
func resolveWriteOnlySecretFromConfig(
	ctx context.Context,
	config tfsdk.Config,
	private privateState,
	attr string,
	plain types.String,
	version types.Int64,
	priorVersion *types.Int64,
) (writeOnlySecret, diag.Diagnostics) {
	var writeOnly types.String
	diags := config.GetAttribute(ctx, path.Root(attr+"_wo"), &writeOnly)
	if diags.HasError() {
		return writeOnlySecret{}, diags
	}

	versionChanged, wasInUse := true, false
	if priorVersion != nil {
		versionChanged = !version.Equal(*priorVersion)

		var d diag.Diagnostics
		wasInUse, d = writeOnlyInUse(ctx, private, attr)
		diags.Append(d...)
	}

	return resolveWriteOnlySecret(plain, writeOnly, versionChanged, wasInUse), diags
}

// integrationKeyWriteOnlyNote documents integration_key_wo on the PagerDuty and
// Opsgenie integration resources.
const integrationKeyWriteOnlyNote = "The integration key can be set through the write-only `integration_key_wo` attribute instead of `integration_key`, " +
	"so that it is not stored in the Terraform state. Write-only attributes require Terraform 1.11 or later. Change `integration_key_wo_version` to send a new key.\n\n" +
	"Sentry returns the integration key, so `integration_key` is populated on import while `integration_key_wo` and `integration_key_wo_version` are not. " +
	"If the configuration uses `integration_key_wo`, the next apply sends it and removes `integration_key` from the state."
