package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
