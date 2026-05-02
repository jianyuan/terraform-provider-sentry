package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
)

func TestIssueAlertResource_elemTypesInitialized(t *testing.T) {
	if len(issueAlertConditionV2ElemType.AttrTypes) == 0 {
		t.Fatal("issueAlertConditionV2ElemType was not initialized")
	}
	if len(issueAlertFilterV2ElemType.AttrTypes) == 0 {
		t.Fatal("issueAlertFilterV2ElemType was not initialized")
	}
	if len(issueAlertActionV2ElemType.AttrTypes) == 0 {
		t.Fatal("issueAlertActionV2ElemType was not initialized")
	}
}

func TestIssueAlertModel_v2FieldsHandleUnknown(t *testing.T) {
	model := IssueAlertModel{
		ConditionsV2: types.ListUnknown(issueAlertConditionV2ElemType),
		FiltersV2:    types.ListUnknown(issueAlertFilterV2ElemType),
		ActionsV2:    types.ListUnknown(issueAlertActionV2ElemType),
	}

	if !model.ConditionsV2.IsUnknown() {
		t.Error("expected ConditionsV2 to be unknown")
	}
	if !model.FiltersV2.IsUnknown() {
		t.Error("expected FiltersV2 to be unknown")
	}
	if !model.ActionsV2.IsUnknown() {
		t.Error("expected ActionsV2 to be unknown")
	}
}

func TestIssueAlertModel_v2FieldsHandleNull(t *testing.T) {
	model := IssueAlertModel{
		ConditionsV2: types.ListNull(issueAlertConditionV2ElemType),
		FiltersV2:    types.ListNull(issueAlertFilterV2ElemType),
		ActionsV2:    types.ListNull(issueAlertActionV2ElemType),
	}

	if !model.ConditionsV2.IsNull() {
		t.Error("expected ConditionsV2 to be null")
	}
	if !model.FiltersV2.IsNull() {
		t.Error("expected FiltersV2 to be null")
	}
	if !model.ActionsV2.IsNull() {
		t.Error("expected ActionsV2 to be null")
	}
}

func TestIssueAlertConditionModel_ToApi_emptyElement(t *testing.T) {
	model := IssueAlertConditionModel{}
	_, diags := model.ToApi(context.Background())
	if !diags.HasError() {
		t.Fatal("expected error for empty condition element")
	}
	found := false
	for _, d := range diags.Errors() {
		if d.Summary() == "Exactly one condition must be set" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Exactly one condition must be set' error, got: %s", diags)
	}
}

func TestIssueAlertFilterModel_ToApi_emptyElement(t *testing.T) {
	model := IssueAlertFilterModel{}
	_, diags := model.ToApi(context.Background())
	if !diags.HasError() {
		t.Fatal("expected error for empty filter element")
	}
	found := false
	for _, d := range diags.Errors() {
		if d.Summary() == "Exactly one filter must be set" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Exactly one filter must be set' error, got: %s", diags)
	}
}

func TestIssueAlertActionModel_ToApi_emptyElement(t *testing.T) {
	model := IssueAlertActionModel{}
	_, diags := model.ToApi(context.Background())
	if !diags.HasError() {
		t.Fatal("expected error for empty action element")
	}
	found := false
	for _, d := range diags.Errors() {
		if d.Summary() == "Exactly one action must be set" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Exactly one action must be set' error, got: %s", diags)
	}
}

func TestIssueAlertConditionModel_ToApi_validCondition(t *testing.T) {
	model := IssueAlertConditionModel{
		FirstSeenEvent: &IssueAlertConditionFirstSeenEventModel{},
	}
	result, diags := model.ToApi(context.Background())
	if diags.HasError() {
		t.Fatalf("unexpected error: %s", diags)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestIssueAlertFilterModel_ToApi_validFilter(t *testing.T) {
	model := IssueAlertFilterModel{
		LatestRelease: &IssueAlertFilterLatestReleaseModel{},
	}
	result, diags := model.ToApi(context.Background())
	if diags.HasError() {
		t.Fatalf("unexpected error: %s", diags)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestIssueAlertActionModel_ToApi_validAction(t *testing.T) {
	model := IssueAlertActionModel{
		NotifyEvent: &IssueAlertActionNotifyEventModel{},
	}
	result, diags := model.ToApi(context.Background())
	if diags.HasError() {
		t.Fatalf("unexpected error: %s", diags)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestIssueAlertActionSlackNotifyServiceModel_ToApi_withChannelId(t *testing.T) {
	channelId := "C1234567890"
	model := IssueAlertActionSlackNotifyServiceModel{
		Workspace: types.StringValue("ws123"),
		Channel:   sentrytypes.SlackChannelValue("#general"),
		ChannelId: types.StringValue(channelId),
	}
	result, diags := model.ToApi(context.Background())
	if diags.HasError() {
		t.Fatalf("unexpected error: %s", diags)
	}
	slack, err := result.AsProjectRuleActionSlackNotifyService()
	if err != nil {
		t.Fatalf("failed to decode result: %s", err)
	}
	if slack.ChannelId == nil || *slack.ChannelId != channelId {
		t.Errorf("expected channel_id %q, got %v", channelId, slack.ChannelId)
	}
}

func TestIssueAlertActionSlackNotifyServiceModel_ToApi_withoutChannelId(t *testing.T) {
	model := IssueAlertActionSlackNotifyServiceModel{
		Workspace: types.StringValue("ws123"),
		Channel:   sentrytypes.SlackChannelValue("#general"),
		ChannelId: types.StringNull(),
	}
	result, diags := model.ToApi(context.Background())
	if diags.HasError() {
		t.Fatalf("unexpected error: %s", diags)
	}
	slack, err := result.AsProjectRuleActionSlackNotifyService()
	if err != nil {
		t.Fatalf("failed to decode result: %s", err)
	}
	if slack.ChannelId != nil {
		t.Errorf("expected channel_id to be nil, got %q", *slack.ChannelId)
	}
}

func TestIssueAlertActionSlackNotifyServiceModel_Fill_setsChannelId(t *testing.T) {
	channelId := "C9876543210"
	action := apiclient.ProjectRuleActionSlackNotifyService{
		Workspace: "ws123",
		Channel:   "#general",
		ChannelId: &channelId,
	}
	var model IssueAlertActionSlackNotifyServiceModel
	diags := model.Fill(context.Background(), action)
	if diags.HasError() {
		t.Fatalf("unexpected error: %s", diags)
	}
	if model.ChannelId.IsNull() || model.ChannelId.ValueString() != channelId {
		t.Errorf("expected channel_id %q, got %q", channelId, model.ChannelId.ValueString())
	}
}

func TestSlackChannel_SemanticEquals_ignoresHashPrefix(t *testing.T) {
	tests := []struct {
		name  string
		a, b  string
		equal bool
	}{
		{"both with hash", "#general", "#general", true},
		{"both without hash", "general", "general", true},
		{"first with hash", "#general", "general", true},
		{"second with hash", "general", "#general", true},
		{"different channels", "#general", "#random", false},
		{"different without hash", "general", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := sentrytypes.SlackChannelValue(tt.a)
			b := sentrytypes.SlackChannelValue(tt.b)
			result, diags := a.StringSemanticEquals(context.Background(), b)
			if diags.HasError() {
				t.Fatalf("unexpected error: %s", diags)
			}
			if result != tt.equal {
				t.Errorf("SlackChannel(%q).SemanticEquals(%q) = %v, want %v", tt.a, tt.b, result, tt.equal)
			}
		})
	}
}

func TestReorderToMatchPrior_preservesPriorOrder(t *testing.T) {
	identity := func(s string) string { return s }

	prior := []string{"a", "b", "c"}
	incoming := []string{"c", "a", "b"}
	got := reorderToMatchPrior(prior, incoming, identity)
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestReorderToMatchPrior_appendsNewItems(t *testing.T) {
	identity := func(s string) string { return s }

	prior := []string{"a", "b"}
	incoming := []string{"c", "b", "a"}
	got := reorderToMatchPrior(prior, incoming, identity)
	// prior items first in order, then unmatched "c"
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestReorderToMatchPrior_handlesDuplicateTypes(t *testing.T) {
	identity := func(s string) string { return s }

	prior := []string{"x", "x", "y"}
	incoming := []string{"y", "x", "x"}
	got := reorderToMatchPrior(prior, incoming, identity)
	want := []string{"x", "x", "y"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestReorderToMatchPrior_emptyPriorReturnsincoming(t *testing.T) {
	identity := func(s string) string { return s }

	prior := []string{}
	incoming := []string{"c", "b", "a"}
	got := reorderToMatchPrior(prior, incoming, identity)
	if len(got) != len(incoming) {
		t.Fatalf("len mismatch: got %v, want %v", got, incoming)
	}
	for i := range incoming {
		if got[i] != incoming[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], incoming[i])
		}
	}
}
