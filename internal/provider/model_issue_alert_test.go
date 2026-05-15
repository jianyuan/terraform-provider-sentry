package provider

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStripLegacyActionDisplayFields(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    string
		expected map[string]interface{}
	}{
		"removes top-level display fields": {
			input:    `{"id":"x","name":"display","formFields":{"a":1},"hasSchemaFormConfig":true}`,
			expected: map[string]interface{}{"id": "x"},
		},
		"strips label from settings entries": {
			input: `{"id":"sentry.rules.actions.notify_event_sentry_app.NotifyEventSentryAppAction","sentryAppInstallationUuid":"u","settings":[{"label":"[Team] foo","name":"target","value":"Group:1"},{"label":"P3","name":"urgency","value":"abc"}]}`,
			expected: map[string]interface{}{
				"id":                        "sentry.rules.actions.notify_event_sentry_app.NotifyEventSentryAppAction",
				"sentryAppInstallationUuid": "u",
				"settings": []interface{}{
					map[string]interface{}{"name": "target", "value": "Group:1"},
					map[string]interface{}{"name": "urgency", "value": "abc"},
				},
			},
		},
		"leaves settings entries without label untouched": {
			input: `{"id":"x","settings":[{"name":"n","value":"v"}]}`,
			expected: map[string]interface{}{
				"id": "x",
				"settings": []interface{}{
					map[string]interface{}{"name": "n", "value": "v"},
				},
			},
		},
		"no settings key is a no-op": {
			input:    `{"id":"x"}`,
			expected: map[string]interface{}{"id": "x"},
		},
	}
	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := stripLegacyActionDisplayFields([]byte(tc.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var gotMap map[string]interface{}
			if err := json.Unmarshal(got, &gotMap); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			if diff := cmp.Diff(gotMap, tc.expected); diff != "" {
				t.Errorf("unexpected stripped JSON (-got, +expected): %s", diff)
			}
		})
	}
}
