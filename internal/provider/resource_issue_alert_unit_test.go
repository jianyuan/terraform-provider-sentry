package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

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

func TestIssueAlertActionSlackNotifyServiceModel_ToApi_withChannelId(t *testing.T) {
	channelId := "C1234567890"
	model := IssueAlertActionSlackNotifyServiceModel{
		Workspace: types.StringValue("ws123"),
		Channel:   sentrytypes.NewSlackChannelValue("#general"),
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
		Channel:   sentrytypes.NewSlackChannelValue("#general"),
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
			a := sentrytypes.NewSlackChannelValue(tt.a)
			b := sentrytypes.NewSlackChannelValue(tt.b)
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

// nullActionModel returns an IssueAlertActionModel with every action variant
// set to null, so a caller can populate exactly one variant.
func nullActionModel(ctx context.Context) IssueAlertActionModel {
	return IssueAlertActionModel{
		NotifyEmail:                  supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEmailModel](ctx),
		NotifyEvent:                  supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEventModel](ctx),
		NotifyEventService:           supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEventServiceModel](ctx),
		NotifyEventSentryApp:         supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionNotifyEventSentryAppModel](ctx),
		OpsgenieNotifyTeam:           supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionOpsgenieNotifyTeam](ctx),
		PagerDutyNotifyService:       supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionPagerDutyNotifyServiceModel](ctx),
		SlackNotifyService:           supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionSlackNotifyServiceModel](ctx),
		MsTeamsNotifyService:         supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionMsTeamsNotifyServiceModel](ctx),
		DiscordNotifyService:         supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionDiscordNotifyServiceModel](ctx),
		JiraCreateTicket:             supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionJiraCreateTicketModel](ctx),
		JiraServerCreateTicket:       supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionJiraServerCreateTicketModel](ctx),
		GitHubCreateTicket:           supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionGitHubCreateTicketModel](ctx),
		GitHubEnterpriseCreateTicket: supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionGitHubEnterpriseCreateTicketModel](ctx),
		AzureDevopsCreateTicket:      supertypes.NewSingleNestedObjectValueOfNull[IssueAlertActionAzureDevopsCreateTicketModel](ctx),
	}
}

func slackAction(ctx context.Context, workspace, channel string) IssueAlertActionModel {
	m := nullActionModel(ctx)
	m.SlackNotifyService = supertypes.NewSingleNestedObjectValueOf(ctx, &IssueAlertActionSlackNotifyServiceModel{
		Name:      types.StringNull(),
		Workspace: types.StringValue(workspace),
		Channel:   sentrytypes.NewSlackChannelValue(channel),
		ChannelId: types.StringValue("C123"),
		Tags:      sentrytypes.StringSetNull(),
		Notes:     types.StringNull(),
	})
	return m
}

// A slack_notify_service action key must ignore a leading "#" on the channel,
// matching sentrytypes.SlackChannel's StringSemanticEquals. Sentry returns the
// channel without the "#", so without this the prior and incoming keys diverge
// and reorderToMatchPrior cannot pair them.
func TestIssueAlertActionModelKey_slackChannelIgnoresHashPrefix(t *testing.T) {
	ctx := context.Background()
	key := issueAlertActionModelKey(ctx)

	withHash := key(slackAction(ctx, "ws123", "#errors-prod"))
	withoutHash := key(slackAction(ctx, "ws123", "errors-prod"))

	if withHash != withoutHash {
		t.Errorf("slack action key should ignore leading '#': %q != %q", withHash, withoutHash)
	}
}

// End-to-end: when Sentry reorders actions AND returns channels without the
// "#", reorderToMatchPrior must still restore the prior order.
func TestReorderToMatchPrior_slackActionsWithNormalizedChannel(t *testing.T) {
	ctx := context.Background()
	keyFn := issueAlertActionModelKey(ctx)

	prior := []IssueAlertActionModel{
		slackAction(ctx, "A", "#errors-prod"),
		slackAction(ctx, "B", "#errors-staging"),
	}
	// Sentry: reordered, and "#" stripped.
	incoming := []IssueAlertActionModel{
		slackAction(ctx, "B", "errors-staging"),
		slackAction(ctx, "A", "errors-prod"),
	}

	got := reorderToMatchPrior(prior, incoming, keyFn)
	if len(got) != 2 {
		t.Fatalf("len mismatch: got %d, want 2", len(got))
	}
	if ws := got[0].SlackNotifyService.MustGet(ctx).Workspace.ValueString(); ws != "A" {
		t.Errorf("index 0: got workspace %q, want A", ws)
	}
	if ws := got[1].SlackNotifyService.MustGet(ctx).Workspace.ValueString(); ws != "B" {
		t.Errorf("index 1: got workspace %q, want B", ws)
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
