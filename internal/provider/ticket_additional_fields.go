package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// reservedTicketAdditionalFieldKeys are keys that must not appear in
// `additional_fields`, either because the action exposes a dedicated attribute
// for them or because Sentry populates them itself when the ticket is created.
var reservedTicketAdditionalFieldKeys = []string{
	// Have a dedicated attribute on the ticket action.
	"project",
	"issuetype",
	"labels",
	"components",
	"priority",
	"reporter",
	"integration",
	// Set by Sentry when the ticket is created.
	"title",
	"description",
	"summary",
	"dynamic_form_fields",
}

// ticketAdditionalFieldKeyValidator rejects `additional_fields` keys that would
// not survive a round trip. It enforces two independent rules; see
// validateTicketFieldKeyCollision and validateTicketFieldKeyCasing.
func ticketAdditionalFieldKeyValidator() validator.String {
	return ticketAdditionalFieldKeyValidatorImpl{}
}

type ticketAdditionalFieldKeyValidatorImpl struct{}

func (v ticketAdditionalFieldKeyValidatorImpl) Description(ctx context.Context) string {
	return "key must survive Sentry's snake_case conversion unchanged and must not duplicate a dedicated attribute"
}

func (v ticketAdditionalFieldKeyValidatorImpl) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ticketAdditionalFieldKeyValidatorImpl) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	key := req.ConfigValue.ValueString()

	// Sentry stores the converted form of the key, so both checks compare
	// against that rather than against what the user typed.
	stored := camelToSnakeCase(key)

	if validateTicketFieldKeyCasing(req, resp, key, stored); resp.Diagnostics.HasError() {
		return
	}

	validateTicketFieldKeyCollision(req, resp, key, stored)
}

// validateTicketFieldKeyCollision rejects keys that collide with a dedicated
// attribute or with a value Sentry sets itself.
//
// This guards a provider-side bug, not an API one, and is required
// indefinitely. The generated `additionalFields` marshaller writes
// AdditionalProperties into the same JSON object *after* the typed fields, so a
// passthrough key silently overwrites the dedicated attribute. In practice that
// means `additional_fields = { labels = "x" }` discards whatever the `labels`
// attribute was set to, Terraform reports "Provider produced inconsistent
// result after apply", and the resource is left orphaned in Sentry.
func validateTicketFieldKeyCollision(req validator.StringRequest, resp *validator.StringResponse, key, stored string) {
	for _, reserved := range reservedTicketAdditionalFieldKeys {
		if stored != reserved {
			continue
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid additional field key",
			fmt.Sprintf(
				"The key %q is set by Sentry or has a dedicated attribute on this action. "+
					"Setting it here would silently override that attribute. Set it there instead.",
				key,
			),
		)
		return
	}
}

// validateTicketFieldKeyCasing rejects keys that Sentry's API would rewrite.
//
// TODO: remove this check, camelToSnakeCase, and their tests once
// getsentry/sentry#122398 has shipped and the provider's minimum supported
// Sentry version includes it. It exists solely to work around an upstream
// asymmetry: `TicketingActionHandler.serialize_data` exempts
// `additional_fields` from case conversion when reading, but the write path
// (`BaseActionValidator`, via `CamelSnakeSerializer.__init__`) runs
// camel_to_snake_case recursively with no such exemption, so third-party field
// IDs are mangled on the way in. That PR exempts nested keys under
// `additional_fields`, after which camelCase field IDs round-trip verbatim and
// rejecting them here would be wrong.
//
// Note this only affects Sentry SaaS and self-hosted versions predating the
// fix; the collision check above is unrelated and stays.
//
// Only an inserted underscore actually breaks the field. Jira Cloud lowercases
// both the payload key and the Jira field ID before matching, so a pure case
// change (`UPPER` -> `upper`) still resolves to the same field; `fixVersions`
// -> `fix_versions` does not, and Jira Server matches its field ID exactly so
// it drops the value with no error at all. Hence the test is whether conversion
// changes the key by more than lowercasing it.
func validateTicketFieldKeyCasing(req validator.StringRequest, resp *validator.StringResponse, key, stored string) {
	if stored == strings.ToLower(key) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid additional field key",
		fmt.Sprintf(
			"Sentry's API converts camelCase keys to snake_case when storing them, which would change "+
				"%[1]q to %[2]q and drop the field when the ticket is created. Use the all-lowercase "+
				"spelling %[3]q instead -- Jira matches field IDs case-insensitively, so it resolves to "+
				"the same field.",
			key, stored, strings.ToLower(key),
		),
	)
}

// camelToSnakeCase mirrors Django's re_camel_case substitution, which Sentry
// applies to incoming payload keys. It is used to determine the key Sentry will
// actually store, and to show the user what their key would be turned into.
func camelToSnakeCase(value string) string {
	var b strings.Builder
	runes := []rune(value)
	for i, r := range runes {
		isUpper := r >= 'A' && r <= 'Z'
		if isUpper && i > 0 {
			prev := runes[i-1]
			prevIsLower := prev >= 'a' && prev <= 'z'
			// Mirrors Django's `[A-Z](?![A-Z]|$)`: a run-final uppercase only splits
			// when followed by a lowercase character, not at end of string.
			nextIsNotUpper := i+1 < len(runes) && (runes[i+1] < 'A' || runes[i+1] > 'Z')
			if prevIsLower || nextIsNotUpper {
				b.WriteRune('_')
			}
		}
		b.WriteRune(r)
	}
	return strings.Trim(strings.ToLower(b.String()), "_")
}
