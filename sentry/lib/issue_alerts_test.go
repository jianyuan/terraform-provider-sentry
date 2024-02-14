package sentry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueAlertsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/rules/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
			  "environment": "production",
			  "actionMatch": "any",
			  "frequency": 30,
			  "name": "Notify errors",
			  "conditions": [
				{
				  "id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
				  "name": "An issue is first seen",
				  "value": 500,
				  "interval": "1h"
				}
			  ],
			  "id": "12345",
			  "actions": [
				{
				  "name": "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
				  "tags": "environment",
				  "channel_id": "XX00X0X0X",
				  "workspace": "1234",
				  "id": "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
				  "channel": "#dummy-channel"
				}
			  ],
			  "dateCreated": "2019-08-24T18:12:16.321Z"
			}
		]`)
	})

	ctx := context.Background()
	alerts, _, err := client.IssueAlerts.List(ctx, "the-interstellar-jurisdiction", "pump-station", nil)
	require.NoError(t, err)

	expected := []*IssueAlert{
		{
			ID:          String("12345"),
			ActionMatch: String("any"),
			Environment: String("production"),
			Frequency:   Int(30),
			Name:        String("Notify errors"),
			Conditions: []*IssueAlertCondition{
				{
					"id":       "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
					"name":     "An issue is first seen",
					"value":    json.Number("500"),
					"interval": "1h",
				},
			},
			Actions: []*IssueAlertAction{
				{
					"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
					"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
					"tags":       "environment",
					"channel_id": "XX00X0X0X",
					"channel":    "#dummy-channel",
					"workspace":  "1234",
				},
			},
			DateCreated: Time(mustParseTime("2019-08-24T18:12:16.321Z")),
		},
	}
	require.Equal(t, expected, alerts)

}

func TestIssueAlertsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/rules/11185158/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "11185158",
			"conditions": [
				{
					"id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
					"name": "A new issue is created"
				},
				{
					"id": "sentry.rules.conditions.regression_event.RegressionEventCondition",
					"name": "The issue changes state from resolved to unresolved"
				},
				{
					"id": "sentry.rules.conditions.reappeared_event.ReappearedEventCondition",
					"name": "The issue changes state from ignored to unresolved"
				},
				{
					"interval": "1h",
					"id": "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
					"comparisonType": "count",
					"value": 100,
					"name": "The issue is seen more than 100 times in 1h"
				},
				{
					"interval": "1h",
					"id": "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
					"comparisonType": "count",
					"value": 100,
					"name": "The issue is seen by more than 100 users in 1h"
				},
				{
					"interval": "1h",
					"id": "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
					"comparisonType": "count",
					"value": 100,
					"name": "The issue affects more than 100.0 percent of sessions in 1h"
				}
			],
			"filters": [
				{
					"comparison_type": "older",
					"time": "minute",
					"id": "sentry.rules.filters.age_comparison.AgeComparisonFilter",
					"value": 10,
					"name": "The issue is older than 10 minute"
				},
				{
					"id": "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
					"value": 10,
					"name": "The issue has happened at least 10 times"
				},
				{
					"targetType": "Team",
					"id": "sentry.rules.filters.assigned_to.AssignedToFilter",
					"targetIdentifier": 1322366,
					"name": "The issue is assigned to Team"
				},
				{
					"id": "sentry.rules.filters.latest_release.LatestReleaseFilter",
					"name": "The event is from the latest release"
				},
				{
					"attribute": "message",
					"match": "co",
					"id": "sentry.rules.filters.event_attribute.EventAttributeFilter",
					"value": "test",
					"name": "The event's message value contains test"
				},
				{
					"match": "co",
					"id": "sentry.rules.filters.tagged_event.TaggedEventFilter",
					"key": "test",
					"value": "test",
					"name": "The event's tags match test contains test"
				},
				{
					"level": "50",
					"match": "eq",
					"id": "sentry.rules.filters.level.LevelFilter",
					"name": "The event's level is equal to fatal"
				}
			],
			"actions": [
				{
					"targetType": "IssueOwners",
					"id": "sentry.mail.actions.NotifyEmailAction",
					"targetIdentifier": "",
					"name": "Send a notification to IssueOwners"
				},
				{
					"targetType": "Team",
					"id": "sentry.mail.actions.NotifyEmailAction",
					"targetIdentifier": 1322366,
					"name": "Send a notification to Team"
				},
				{
					"targetType": "Member",
					"id": "sentry.mail.actions.NotifyEmailAction",
					"targetIdentifier": 94401,
					"name": "Send a notification to Member"
				},
				{
					"id": "sentry.rules.actions.notify_event.NotifyEventAction",
					"name": "Send a notification (for all legacy integrations)"
				}
			],
			"actionMatch": "any",
			"filterMatch": "any",
			"frequency": 30,
			"name": "My Rule Name",
			"dateCreated": "2022-05-23T19:54:30.860115Z",
			"owner": "team:1322366",
			"createdBy": {
				"id": 94401,
				"name": "John Doe",
				"email": "test@example.com"
			},
			"environment": null,
			"projects": [
				"python"
			]
		}`)
	})

	ctx := context.Background()
	alerts, _, err := client.IssueAlerts.Get(ctx, "the-interstellar-jurisdiction", "pump-station", "11185158")
	require.NoError(t, err)

	expected := &IssueAlert{
		ID: String("11185158"),
		Conditions: []*IssueAlertCondition{
			{
				"id":   "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
				"name": "A new issue is created",
			},
			{
				"id":   "sentry.rules.conditions.regression_event.RegressionEventCondition",
				"name": "The issue changes state from resolved to unresolved",
			},
			{
				"id":   "sentry.rules.conditions.reappeared_event.ReappearedEventCondition",
				"name": "The issue changes state from ignored to unresolved",
			},
			{
				"interval":       "1h",
				"id":             "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
				"comparisonType": "count",
				"value":          json.Number("100"),
				"name":           "The issue is seen more than 100 times in 1h",
			},
			{
				"interval":       "1h",
				"id":             "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
				"comparisonType": "count",
				"value":          json.Number("100"),
				"name":           "The issue is seen by more than 100 users in 1h",
			},
			{
				"interval":       "1h",
				"id":             "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
				"comparisonType": "count",
				"value":          json.Number("100"),
				"name":           "The issue affects more than 100.0 percent of sessions in 1h",
			},
		},
		Filters: []*IssueAlertFilter{
			{
				"comparison_type": "older",
				"time":            "minute",
				"id":              "sentry.rules.filters.age_comparison.AgeComparisonFilter",
				"value":           json.Number("10"),
				"name":            "The issue is older than 10 minute",
			},
			{
				"id":    "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
				"value": json.Number("10"),
				"name":  "The issue has happened at least 10 times",
			},
			{
				"targetType":       "Team",
				"id":               "sentry.rules.filters.assigned_to.AssignedToFilter",
				"targetIdentifier": json.Number("1322366"),
				"name":             "The issue is assigned to Team",
			},
			{
				"id":   "sentry.rules.filters.latest_release.LatestReleaseFilter",
				"name": "The event is from the latest release",
			},
			{
				"attribute": "message",
				"match":     "co",
				"id":        "sentry.rules.filters.event_attribute.EventAttributeFilter",
				"value":     "test",
				"name":      "The event's message value contains test",
			},
			{
				"match": "co",
				"id":    "sentry.rules.filters.tagged_event.TaggedEventFilter",
				"key":   "test",
				"value": "test",
				"name":  "The event's tags match test contains test",
			},
			{
				"level": "50",
				"match": "eq",
				"id":    "sentry.rules.filters.level.LevelFilter",
				"name":  "The event's level is equal to fatal",
			},
		},
		Actions: []*IssueAlertAction{
			{
				"targetType":       "IssueOwners",
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"targetIdentifier": "",
				"name":             "Send a notification to IssueOwners",
			},
			{
				"targetType":       "Team",
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"targetIdentifier": json.Number("1322366"),
				"name":             "Send a notification to Team",
			},
			{
				"targetType":       "Member",
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"targetIdentifier": json.Number("94401"),
				"name":             "Send a notification to Member",
			},
			{
				"id":   "sentry.rules.actions.notify_event.NotifyEventAction",
				"name": "Send a notification (for all legacy integrations)",
			},
		},
		ActionMatch: String("any"),
		FilterMatch: String("any"),
		Frequency:   Int(30),
		Name:        String("My Rule Name"),
		DateCreated: Time(mustParseTime("2022-05-23T19:54:30.860115Z")),
		Owner:       String("team:1322366"),
		CreatedBy: &IssueAlertCreatedBy{
			ID:    Int(94401),
			Name:  String("John Doe"),
			Email: String("test@example.com"),
		},
		Projects: []string{"python"},
	}
	require.Equal(t, expected, alerts)

}

func TestIssueAlertsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/rules/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSONValue(t, map[string]interface{}{
			"actionMatch": "all",
			"environment": "production",
			"frequency":   30,
			"name":        "Notify errors",
			"conditions": []map[string]interface{}{
				{
					"interval": "1h",
					"name":     "The issue is seen more than 10 times in 1h",
					"value":    10,
					"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
				},
			},
			"actions": []map[string]interface{}{
				{
					"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
					"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
					"tags":       "environment",
					"channel":    "#dummy-channel",
					"channel_id": "XX00X0X0X",
					"workspace":  "1234",
				},
			},
		}, r)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "123456",
			"actionMatch": "all",
			"environment": "production",
			"frequency": 30,
			"name": "Notify errors",
			"conditions": [
				{
					"interval": "1h",
					"name": "The issue is seen more than 10 times in 1h",
					"value": 10,
					"id": "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
				}
			],
			"actions": [
				{
					"id": "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
					"name": "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
					"tags": "environment",
					"channel_id": "XX00X0X0X",
					"workspace": "1234",
					"channel": "#dummy-channel"
				}
			],
			"dateCreated": "2019-08-24T18:12:16.321Z"
		}`)
	})

	params := &IssueAlert{
		ActionMatch: String("all"),
		Environment: String("production"),
		Frequency:   Int(30),
		Name:        String("Notify errors"),
		Conditions: []*IssueAlertCondition{
			{
				"interval": "1h",
				"name":     "The issue is seen more than 10 times in 1h",
				"value":    json.Number("10"),
				"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
			},
		},
		Actions: []*IssueAlertAction{
			{
				"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
				"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
				"tags":       "environment",
				"channel_id": "XX00X0X0X",
				"workspace":  "1234",
				"channel":    "#dummy-channel",
			},
		},
	}
	ctx := context.Background()
	alerts, _, err := client.IssueAlerts.Create(ctx, "the-interstellar-jurisdiction", "pump-station", params)
	require.NoError(t, err)

	expected := &IssueAlert{
		ID:          String("123456"),
		ActionMatch: String("all"),
		Environment: String("production"),
		Frequency:   Int(30),
		Name:        String("Notify errors"),
		Conditions: []*IssueAlertCondition{
			{
				"interval": "1h",
				"name":     "The issue is seen more than 10 times in 1h",
				"value":    json.Number("10"),
				"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
			},
		},
		Actions: []*IssueAlertAction{
			{
				"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
				"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
				"tags":       "environment",
				"channel_id": "XX00X0X0X",
				"channel":    "#dummy-channel",
				"workspace":  "1234",
			},
		},
		DateCreated: Time(mustParseTime("2019-08-24T18:12:16.321Z")),
	}
	require.Equal(t, expected, alerts)

}

func TestIssueAlertsService_CreateWithAsyncTask(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/rule-task/fakeuuid/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"status": "success",
			"error": null,
			"rule": {
				"id": "123456",
				"actionMatch": "all",
				"environment": "production",
				"frequency": 30,
				"name": "Notify errors",
				"conditions": [
					{
						"interval": "1h",
						"name": "The issue is seen more than 10 times in 1h",
						"value": 10,
						"id": "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
					}
				],
				"actions": [
					{
						"id": "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
						"name": "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
						"tags": "environment",
						"channel_id": "XX00X0X0X",
						"workspace": "1234",
						"channel": "#dummy-channel"
					}
				],
				"dateCreated": "2019-08-24T18:12:16.321Z"
			}
		}`)
	})
	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/rules/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSONValue(t, map[string]interface{}{
			"actionMatch": "all",
			"environment": "production",
			"frequency":   30,
			"name":        "Notify errors",
			"conditions": []map[string]interface{}{
				{
					"interval": "1h",
					"name":     "The issue is seen more than 10 times in 1h",
					"value":    10,
					"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
				},
			},
			"actions": []map[string]interface{}{
				{
					"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
					"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
					"tags":       "environment",
					"channel":    "#dummy-channel",
					"channel_id": "XX00X0X0X",
					"workspace":  "1234",
				},
			},
		}, r)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"uuid": "fakeuuid"}`)

	})

	params := &IssueAlert{
		ActionMatch: String("all"),
		Environment: String("production"),
		Frequency:   Int(30),
		Name:        String("Notify errors"),
		Conditions: []*IssueAlertCondition{
			{
				"interval": "1h",
				"name":     "The issue is seen more than 10 times in 1h",
				"value":    json.Number("10"),
				"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
			},
		},
		Actions: []*IssueAlertAction{
			{
				"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
				"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
				"tags":       "environment",
				"channel_id": "XX00X0X0X",
				"workspace":  "1234",
				"channel":    "#dummy-channel",
			},
		},
	}
	ctx := context.Background()
	alert, _, err := client.IssueAlerts.Create(ctx, "the-interstellar-jurisdiction", "pump-station", params)
	require.NoError(t, err)

	expected := &IssueAlert{
		ID:          String("123456"),
		ActionMatch: String("all"),
		Environment: String("production"),
		Frequency:   Int(30),
		Name:        String("Notify errors"),
		Conditions: []*IssueAlertCondition{
			{
				"interval": "1h",
				"name":     "The issue is seen more than 10 times in 1h",
				"value":    json.Number("10"),
				"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
			},
		},
		Actions: []*IssueAlertAction{
			{
				"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
				"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
				"tags":       "environment",
				"channel_id": "XX00X0X0X",
				"channel":    "#dummy-channel",
				"workspace":  "1234",
			},
		},
		DateCreated: Time(mustParseTime("2019-08-24T18:12:16.321Z")),
	}
	require.Equal(t, expected, alert)

}

func TestIssueAlertsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	params := &IssueAlert{
		ID:          String("12345"),
		ActionMatch: String("all"),
		FilterMatch: String("any"),
		Environment: String("staging"),
		Frequency:   Int(30),
		Name:        String("Notify errors"),
		Conditions: []*IssueAlertCondition{
			{
				"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
				"value":    500,
				"interval": "1h",
			},
		},
		Actions: []*IssueAlertAction{
			{
				"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
				"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
				"tags":       "environment",
				"channel_id": "XX00X0X0X",
				"channel":    "#dummy-channel",
				"workspace":  "1234",
			},
		},
		Filters: []*IssueAlertFilter{
			{
				"id":    "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
				"name":  "The issue has happened at least 4 times",
				"value": 4,
			},
			{
				"attribute": "message",
				"id":        "sentry.rules.filters.event_attribute.EventAttributeFilter",
				"match":     "eq",
				"name":      "The event's message value equals test",
				"value":     "test",
			},
		},
		DateCreated: Time(mustParseTime("2019-08-24T18:12:16.321Z")),
	}

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/rules/12345/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSONValue(t, map[string]interface{}{
			"id":          "12345",
			"actionMatch": "all",
			"filterMatch": "any",
			"environment": "staging",
			"frequency":   json.Number("30"),
			"name":        "Notify errors",
			"dateCreated": "2019-08-24T18:12:16.321Z",
			"conditions": []map[string]interface{}{
				{
					"id":       "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
					"value":    json.Number("500"),
					"interval": "1h",
				},
			},
			"actions": []map[string]interface{}{
				{
					"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
					"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
					"tags":       "environment",
					"channel":    "#dummy-channel",
					"channel_id": "XX00X0X0X",
					"workspace":  "1234",
				},
			},
			"filters": []map[string]interface{}{
				{
					"id":    "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
					"name":  "The issue has happened at least 4 times",
					"value": json.Number("4"),
				},
				{
					"attribute": "message",
					"id":        "sentry.rules.filters.event_attribute.EventAttributeFilter",
					"match":     "eq",
					"name":      "The event's message value equals test",
					"value":     "test",
				},
			},
		}, r)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"environment": "staging",
			"actionMatch": "any",
			"frequency": 30,
			"name": "Notify errors",
			"conditions": [
				{
					"id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
					"name": "An issue is first seen"
				}
			],
			"id": "12345",
			"actions": [
				{
					"name": "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
					"tags": "environment",
					"channel_id": "XX00X0X0X",
					"workspace": "1234",
					"id": "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
					"channel": "#dummy-channel"
				}
			],
			"dateCreated": "2019-08-24T18:12:16.321Z"
		}`)
	})
	ctx := context.Background()
	alerts, _, err := client.IssueAlerts.Update(ctx, "the-interstellar-jurisdiction", "pump-station", "12345", params)
	assert.NoError(t, err)

	expected := &IssueAlert{
		ID:          String("12345"),
		ActionMatch: String("any"),
		Environment: String("staging"),
		Frequency:   Int(30),
		Name:        String("Notify errors"),
		Conditions: []*IssueAlertCondition{
			{
				"id":   "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
				"name": "An issue is first seen",
			},
		},
		Actions: []*IssueAlertAction{
			{
				"id":         "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
				"name":       "Send a notification to the Dummy Slack workspace to #dummy-channel and show tags [environment] in notification",
				"tags":       "environment",
				"channel_id": "XX00X0X0X",
				"channel":    "#dummy-channel",
				"workspace":  "1234",
			},
		},
		DateCreated: Time(mustParseTime("2019-08-24T18:12:16.321Z")),
	}
	require.Equal(t, expected, alerts)

}

func TestIssueAlertsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/rules/12345/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	ctx := context.Background()
	_, err := client.IssueAlerts.Delete(ctx, "the-interstellar-jurisdiction", "pump-station", "12345")
	require.NoError(t, err)
}
