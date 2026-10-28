package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func TestNewOptionalTicketString(t *testing.T) {
	// Sentry reports an unset ticket field either by omitting the key or by
	// returning an empty string. Both must normalise to null, otherwise an
	// attribute the user never set shows permanent drift after apply.
	for _, tt := range []struct {
		name     string
		value    *string
		wantNull bool
		wantStr  string
	}{
		{name: "absent key", value: nil, wantNull: true},
		{name: "empty string", value: ptr(""), wantNull: true},
		{name: "set value", value: ptr("oncall,triage"), wantStr: "oncall,triage"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := newOptionalTicketString(tt.value)
			if tt.wantNull {
				if !got.IsNull() {
					t.Fatalf("expected null, got %q", got.Get())
				}
				return
			}
			if got.IsNull() {
				t.Fatal("expected a value, got null")
			}
			if got.Get() != tt.wantStr {
				t.Fatalf("expected %q, got %q", tt.wantStr, got.Get())
			}
		})
	}
}

func TestNewOptionalTicketStringSet(t *testing.T) {
	ctx := context.Background()

	if got := newOptionalTicketStringSet(ctx, nil); !got.IsNull() {
		t.Fatal("expected null for an absent components key")
	}

	if got := newOptionalTicketStringSet(ctx, &[]string{}); !got.IsNull() {
		t.Fatal("expected null for an empty components list")
	}

	got := newOptionalTicketStringSet(ctx, &[]string{"10001", "10002"})
	if got.IsNull() {
		t.Fatal("expected a value for a populated components list")
	}
	elems, diags := got.Get(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(elems) != 2 {
		t.Fatalf("expected 2 components, got %d", len(elems))
	}
}

func ptr(s string) *string { return &s }

func TestCamelToSnakeCase(t *testing.T) {
	// Mirrors Django's re_camel_case, which Sentry applies to every incoming
	// payload key. Used to tell the user what their key would be turned into.
	for in, want := range map[string]string{
		"fixVersions":       "fix_versions",
		"dueDate":           "due_date",
		"lastViewed":        "last_viewed",
		"FixVersions":       "fix_versions",
		"FIXVERSIONS":       "fixversions",
		"labels":            "labels",
		"customfield_10101": "customfield_10101",
	} {
		if got := camelToSnakeCase(in); got != want {
			t.Errorf("camelToSnakeCase(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTicketAdditionalFieldKeyValidator(t *testing.T) {
	ctx := context.Background()
	v := ticketAdditionalFieldKeyValidator()

	for _, tt := range []struct {
		key     string
		wantErr bool
	}{
		// Keys that survive Sentry's conversion unchanged.
		{key: "customfield_10101"},
		{key: "fixversions"},
		{key: "duedate"},
		// Pure case changes are harmless: Sentry stores "upper", and Jira
		// lowercases both sides when matching, so it resolves to the same
		// field. Verified against the live API.
		{key: "UPPER"},
		{key: "Environment"},
		// An inserted underscore does break the field.
		{key: "fixVersions", wantErr: true},
		{key: "dueDate", wantErr: true},
		{key: "MixedCase", wantErr: true},
		// Reserved: has a dedicated attribute, or is set by Sentry itself.
		{key: "labels", wantErr: true},
		{key: "project", wantErr: true},
		{key: "description", wantErr: true},
		{key: "dynamic_form_fields", wantErr: true},
		// Reserved keys must be caught in their post-conversion form too --
		// "Labels" is stored as "labels" and would collide with the dedicated
		// attribute just the same.
		{key: "Labels", wantErr: true},
		{key: "PRIORITY", wantErr: true},
	} {
		t.Run(tt.key, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("additional_fields"),
				ConfigValue: types.StringValue(tt.key),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(ctx, req, resp)

			if got := resp.Diagnostics.HasError(); got != tt.wantErr {
				t.Fatalf("key %q: HasError() = %t, want %t (%v)", tt.key, got, tt.wantErr, resp.Diagnostics)
			}
		})
	}
}

// TestTicketAdditionalFieldsOverwriteTypedAttribute documents why the reserved
// key check is needed independently of Sentry's camelCase bug: the generated
// marshaller writes AdditionalProperties after the typed fields, into the same
// JSON object, so a passthrough key silently wins.
func TestTicketAdditionalFieldsOverwriteTypedAttribute(t *testing.T) {
	var action apiclient.OrganizationWorkflowActionFilterActionJira
	action.Data.AdditionalFields.Project = "10000"
	action.Data.AdditionalFields.Issuetype = "10001"
	labels := "oncall,triage"
	action.Data.AdditionalFields.Labels = &labels

	action.Data.AdditionalFields.Set("labels", "clobbered")

	encoded, err := json.Marshal(action.Data.AdditionalFields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded["labels"] != "clobbered" {
		t.Fatalf("expected the passthrough value to win (documenting the collision), got %v", decoded["labels"])
	}

	// The validator is what stops a user from reaching that state.
	resp := &validator.StringResponse{}
	ticketAdditionalFieldKeyValidator().ValidateString(
		context.Background(),
		validator.StringRequest{
			Path:        path.Root("additional_fields"),
			ConfigValue: types.StringValue("labels"),
		},
		resp,
	)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the validator to reject a key that collides with a dedicated attribute")
	}
}

func TestSeerActivityTriggerSendsStages(t *testing.T) {
	// Sentry validates seer_activity_trigger's comparison as a non-empty list of stages and
	// answers `true` with a 400, so the stages must reach the request as a JSON array.
	ctx := context.Background()

	r := &AlertResource{}
	data := AlertResourceModel{
		TriggerConditions: supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelTriggerConditionsItem{
			{
				FirstSeenEvent:                supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemFirstSeenEvent](ctx),
				IssueResolvedTrigger:          supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemIssueResolvedTrigger](ctx),
				ReappearedEvent:               supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemReappearedEvent](ctx),
				RegressionEvent:               supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemRegressionEvent](ctx),
				EventFrequencyCount:           supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventFrequencyCount](ctx),
				EventFrequencyPercent:         supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventFrequencyPercent](ctx),
				EventUniqueUserFrequencyCount: supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventUniqueUserFrequencyCount](ctx),
				PercentSessionsCount:          supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemPercentSessionsCount](ctx),
				PercentSessionsPercent:        supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemPercentSessionsPercent](ctx),
				SeerActivityTrigger: supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemSeerActivityTrigger{
					Stages: supertypes.NewSetValueOfSlice(ctx, []string{"pr_ready_for_review", "coding_completed"}),
				}),
			},
		}),
		LegacyTriggerConditions: supertypes.NewListValueOfNull[string](ctx),
	}

	got, diags := r.getTriggerConditions(ctx, data)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 trigger condition, got %d", len(got))
	}
	if got[0].Type != "seer_activity_trigger" {
		t.Fatalf("expected type seer_activity_trigger, got %q", got[0].Type)
	}
	b, err := json.Marshal(got[0].Comparison)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `["coding_completed","pr_ready_for_review"]` {
		t.Fatalf("expected the sorted stages as the comparison, got %s", b)
	}
}

func TestSeerActivityTriggerRejectedAsLegacy(t *testing.T) {
	// The passthrough can only send `true`, which Sentry rejects for this type. Fail at plan
	// time with a pointer to the native attribute instead of a 400 at apply.
	ctx := context.Background()

	r := &AlertResource{}
	data := AlertResourceModel{
		TriggerConditions:       supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelTriggerConditionsItem{}),
		LegacyTriggerConditions: supertypes.NewListValueOfSlice(ctx, []string{"seer_activity_trigger"}),
	}

	if _, diags := r.getTriggerConditions(ctx, data); !diags.HasError() {
		t.Fatal("expected an error for seer_activity_trigger in legacy_trigger_conditions")
	}
}

func TestSeerActivityTriggerComparisonDecoding(t *testing.T) {
	// Fill maps a list comparison to the native attribute, and anything else to the
	// passthrough. That split rests on AsOrganizationWorkflowTriggerConditionComparison2
	// failing for the boolean form.
	var stages apiclient.OrganizationWorkflowTriggerCondition_Comparison
	if err := json.Unmarshal([]byte(`["pr_ready_for_review"]`), &stages); err != nil {
		t.Fatal(err)
	}
	got, err := stages.AsOrganizationWorkflowTriggerConditionComparison2()
	if err != nil {
		t.Fatalf("expected a list of stages to decode, got %v", err)
	}
	if len(got) != 1 || got[0] != "pr_ready_for_review" {
		t.Fatalf("unexpected stages %v", got)
	}

	var boolean apiclient.OrganizationWorkflowTriggerCondition_Comparison
	if err := json.Unmarshal([]byte(`true`), &boolean); err != nil {
		t.Fatal(err)
	}
	if _, err := boolean.AsOrganizationWorkflowTriggerConditionComparison2(); err == nil {
		t.Fatal("expected a boolean comparison not to decode as stages")
	}
}
