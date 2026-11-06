package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestNotificationActionResourceModel_Fill_projects(t *testing.T) {
	action := sentry.NotificationAction{
		ID:       new(json.Number("42")),
		Projects: []json.Number{"1", "2"},
	}

	var m NotificationActionResourceModel
	if err := m.Fill(action, map[string]string{"1": "project-one", "2": "project-two"}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expected := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("project-one"),
		types.StringValue("project-two"),
	})
	if !m.Projects.Equal(expected) {
		t.Errorf("expected projects %s, got %s", expected, m.Projects)
	}
}

func TestNotificationActionResourceModel_Fill_unknownProject(t *testing.T) {
	action := sentry.NotificationAction{
		ID:       new(json.Number("42")),
		Projects: []json.Number{"1", "2"},
	}

	var m NotificationActionResourceModel
	if err := m.Fill(action, map[string]string{"1": "project-one"}); err == nil {
		t.Fatalf("expected an error for an unresolvable project ID, got projects %s", m.Projects)
	}
}
