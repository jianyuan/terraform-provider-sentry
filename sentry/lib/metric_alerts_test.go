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

func TestMetricAlertService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/alert-rules/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `[
			{
				"id": "12345",
				"name": "pump-station-alert",
				"environment": "production",
				"dataset": "transactions",
				"eventTypes": ["transaction"],
				"query": "http.url:http://service/unreadmessages",
				"aggregate": "p50(transaction.duration)",
				"thresholdType": 0,
				"resolveThreshold": 100.0,
				"timeWindow": 5.0,
				"triggers": [
					{
						"id": "6789",
						"alertRuleId": "12345",
						"label": "critical",
						"thresholdType": 0,
						"alertThreshold": 55501.0,
						"resolveThreshold": 100.0,
						"dateCreated": "2022-04-07T16:46:48.607583Z",
						"actions": [
							{
								"id": "12345",
								"alertRuleTriggerId": "12345",
								"type": "slack",
								"targetType": "specific",
								"targetIdentifier": "#alert-rule-alerts",
								"inputChannelId": "C038NF00X4F",
								"integrationId": 123,
								"sentryAppId": null,
								"dateCreated": "2022-04-07T16:46:49.154638Z",
								"desc": "Send a Slack notification to #alert-rule-alerts"
							}
						]
					}
				],
				"projects": [
					"pump-station"
				],
				"owner": "pump-station:12345",
				"dateCreated": "2022-04-07T16:46:48.569571Z"
			}
		]`)
	})

	ctx := context.Background()
	alertRules, _, err := client.MetricAlerts.List(ctx, "the-interstellar-jurisdiction", "pump-station", nil)
	require.NoError(t, err)

	expected := []*MetricAlert{
		{
			ID:               String("12345"),
			Name:             String("pump-station-alert"),
			Environment:      String("production"),
			DataSet:          String("transactions"),
			Query:            String("http.url:http://service/unreadmessages"),
			Aggregate:        String("p50(transaction.duration)"),
			EventTypes:       []string{"transaction"},
			ThresholdType:    Int(0),
			ResolveThreshold: Float64(100.0),
			TimeWindow:       Float64(5.0),
			Triggers: []*MetricAlertTrigger{
				{
					ID:               String("6789"),
					AlertRuleID:      String("12345"),
					Label:            String("critical"),
					ThresholdType:    Int(0),
					AlertThreshold:   Float64(55501.0),
					ResolveThreshold: Float64(100.0),
					DateCreated:      Time(mustParseTime("2022-04-07T16:46:48.607583Z")),
					Actions: []*MetricAlertTriggerAction{
						{
							ID:                 String("12345"),
							AlertRuleTriggerID: String("12345"),
							Type:               String("slack"),
							TargetType:         String("specific"),
							TargetIdentifier:   InterfaceString("#alert-rule-alerts"),
							InputChannelID:     String("C038NF00X4F"),
							IntegrationID:      Int(123),
							DateCreated:        Time(mustParseTime("2022-04-07T16:46:49.154638Z")),
							Description:        String("Send a Slack notification to #alert-rule-alerts"),
						},
					},
				},
			},
			Projects:    []string{"pump-station"},
			Owner:       String("pump-station:12345"),
			DateCreated: Time(mustParseTime("2022-04-07T16:46:48.569571Z")),
		},
	}
	require.Equal(t, expected, alertRules)
}

func TestMetricAlertService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/alert-rules/12345/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"id": "12345",
				"name": "pump-station-alert",
				"environment": "production",
				"dataset": "transactions",
				"eventTypes": ["transaction"],
				"query": "http.url:http://service/unreadmessages",
				"aggregate": "p50(transaction.duration)",
				"timeWindow": 10,
				"thresholdType": 0,
				"resolveThreshold": 0,
				"triggers": [
				  {
					"actions": [
					  {
						"alertRuleTriggerId": "56789",
						"dateCreated": "2022-04-15T15:06:01.087054Z",
						"desc": "Send a Slack notification to #alert-rule-alerts",
						"id": "12389",
						"inputChannelId": "C0XXXFKLXXX",
						"integrationId": 111,
						"sentryAppId": null,
						"targetIdentifier": 123456,
						"targetType": "specific",
						"type": "slack"
					  }
					],
					"alertRuleId": "12345",
					"alertThreshold": 10000,
					"dateCreated": "2022-04-15T15:06:01.079598Z",
					"id": "56789",
					"label": "critical",
					"resolveThreshold": 0,
					"thresholdType": 0
				  }
				],
				"projects": [
				  "pump-station"
				],
				"owner": "pump-station:12345",
				"dateCreated": "2022-04-15T15:06:01.05618Z"
			}
		`)
	})

	ctx := context.Background()
	alert, _, err := client.MetricAlerts.Get(ctx, "the-interstellar-jurisdiction", "pump-station", "12345")
	require.NoError(t, err)

	expected := &MetricAlert{
		ID:               String("12345"),
		Name:             String("pump-station-alert"),
		Environment:      String("production"),
		DataSet:          String("transactions"),
		EventTypes:       []string{"transaction"},
		Query:            String("http.url:http://service/unreadmessages"),
		Aggregate:        String("p50(transaction.duration)"),
		TimeWindow:       Float64(10),
		ThresholdType:    Int(0),
		ResolveThreshold: Float64(0),
		Triggers: []*MetricAlertTrigger{
			{
				ID:               String("56789"),
				AlertRuleID:      String("12345"),
				Label:            String("critical"),
				ThresholdType:    Int(0),
				AlertThreshold:   Float64(10000.0),
				ResolveThreshold: Float64(0.0),
				DateCreated:      Time(mustParseTime("2022-04-15T15:06:01.079598Z")),
				Actions: []*MetricAlertTriggerAction{
					{
						ID:                 String("12389"),
						AlertRuleTriggerID: String("56789"),
						Type:               String("slack"),
						TargetType:         String("specific"),
						TargetIdentifier:   InterfaceNumber("123456"),
						InputChannelID:     String("C0XXXFKLXXX"),
						IntegrationID:      Int(111),
						DateCreated:        Time(mustParseTime("2022-04-15T15:06:01.087054Z")),
						Description:        String("Send a Slack notification to #alert-rule-alerts"),
					},
				},
			},
		},
		Projects:    []string{"pump-station"},
		Owner:       String("pump-station:12345"),
		DateCreated: Time(mustParseTime("2022-04-15T15:06:01.05618Z")),
	}
	require.Equal(t, expected, alert)
}

func TestMetricAlertsService_CreateWithAsyncTask(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/alert-rule-task/fakeuuid/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"status": "success",
				"error": null,
				"alertRule": {
					"id": "12345",
					"name": "pump-station-alert",
					"environment": "production",
					"dataset": "transactions",
					"eventTypes": ["transaction"],
					"query": "http.url:http://service/unreadmessages",
					"aggregate": "p50(transaction.duration)",
					"timeWindow": 10,
					"thresholdType": 0,
					"resolveThreshold": 0,
					"triggers": [
					  {
						"actions": [
						  {
							"alertRuleTriggerId": "56789",
							"dateCreated": "2022-04-15T15:06:01.087054Z",
							"desc": "Send a Slack notification to #alert-rule-alerts",
							"id": "12389",
							"inputChannelId": "C0XXXFKLXXX",
							"integrationId": 111,
							"sentryAppId": null,
							"targetIdentifier": "#alert-rule-alerts",
							"targetType": "specific",
							"type": "slack"
						  }
						],
						"alertRuleId": "12345",
						"alertThreshold": 10000,
						"dateCreated": "2022-04-15T15:06:01.079598Z",
						"id": "56789",
						"label": "critical",
						"resolveThreshold": 0,
						"thresholdType": 0
					  }
					],
					"projects": [
					  "pump-station"
					],
					"owner": "pump-station:12345",
					"dateCreated": "2022-04-15T15:06:01.05618Z"
				}
			}
		`)
	})

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/alert-rules/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSONValue(t, map[string]interface{}{
			"id":               "12345",
			"name":             "pump-station-alert",
			"environment":      "production",
			"dataset":          "transactions",
			"eventTypes":       []string{"transaction"},
			"query":            "http.url:http://service/unreadmessages",
			"aggregate":        "p50(transaction.duration)",
			"timeWindow":       10,
			"thresholdType":    0,
			"resolveThreshold": 0,
			"triggers": []map[string]interface{}{
				{
					"actions": []map[string]interface{}{
						{
							"alertRuleTriggerId": "56789",
							"dateCreated":        "2022-04-15T15:06:01.087054Z",
							"desc":               "Send a Slack notification to #alert-rule-alerts",
							"id":                 "12389",
							"inputChannelId":     "C0XXXFKLXXX",
							"integrationId":      111,
							"sentryAppId":        nil,
							"targetIdentifier":   "#alert-rule-alerts",
							"targetType":         "specific",
							"type":               "slack",
						},
					},
					"alertRuleId":      "12345",
					"alertThreshold":   10000,
					"dateCreated":      "2022-04-15T15:06:01.079598Z",
					"id":               "56789",
					"label":            "critical",
					"resolveThreshold": 0,
					"thresholdType":    0,
				},
			},
			"projects":    []string{"pump-station"},
			"owner":       "pump-station:12345",
			"dateCreated": "2022-04-15T15:06:01.05618Z",
		}, r)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"uuid": "fakeuuid"}`)
	})

	params := &MetricAlert{
		Name:             String("pump-station-alert"),
		Environment:      String("production"),
		DataSet:          String("transactions"),
		Query:            String("http.url:http://service/unreadmessages"),
		Aggregate:        String("p50(transaction.duration)"),
		TimeWindow:       Float64(10.0),
		ThresholdType:    Int(0),
		ResolveThreshold: Float64(0),
		Triggers: []*MetricAlertTrigger{
			{
				ID:               String("56789"),
				AlertRuleID:      String("12345"),
				Label:            String("critical"),
				ThresholdType:    Int(0),
				AlertThreshold:   Float64(55501.0),
				ResolveThreshold: Float64(100.0),
				DateCreated:      Time(mustParseTime("2022-04-15T15:06:01.079598Z")),
				Actions: []*MetricAlertTriggerAction{
					{
						ID:                 String("12389"),
						AlertRuleTriggerID: String("56789"),
						Type:               String("slack"),
						TargetType:         String("specific"),
						TargetIdentifier:   InterfaceString("#alert-rule-alerts"),
						InputChannelID:     String("C0XXXFKLXXX"),
						IntegrationID:      Int(123),
						DateCreated:        Time(mustParseTime("2022-04-15T15:06:01.087054Z")),
						Description:        String("Send a Slack notification to #alert-rule-alerts"),
					},
				},
			},
		},
		Projects: []string{"pump-station"},
		Owner:    String("pump-station:12345"),
	}
	ctx := context.Background()
	alertRule, _, err := client.MetricAlerts.Create(ctx, "the-interstellar-jurisdiction", "pump-station", params)
	require.NoError(t, err)

	expected := &MetricAlert{
		ID:               String("12345"),
		Name:             String("pump-station-alert"),
		Environment:      String("production"),
		DataSet:          String("transactions"),
		EventTypes:       []string{"transaction"},
		Query:            String("http.url:http://service/unreadmessages"),
		Aggregate:        String("p50(transaction.duration)"),
		ThresholdType:    Int(0),
		ResolveThreshold: Float64(0),
		TimeWindow:       Float64(10.0),
		Triggers: []*MetricAlertTrigger{
			{
				ID:               String("56789"),
				AlertRuleID:      String("12345"),
				Label:            String("critical"),
				ThresholdType:    Int(0),
				AlertThreshold:   Float64(10000.0),
				ResolveThreshold: Float64(0.0),
				DateCreated:      Time(mustParseTime("2022-04-15T15:06:01.079598Z")),
				Actions: []*MetricAlertTriggerAction{
					{
						ID:                 String("12389"),
						AlertRuleTriggerID: String("56789"),
						Type:               String("slack"),
						TargetType:         String("specific"),
						TargetIdentifier:   InterfaceString("#alert-rule-alerts"),
						InputChannelID:     String("C0XXXFKLXXX"),
						IntegrationID:      Int(111),
						DateCreated:        Time(mustParseTime("2022-04-15T15:06:01.087054Z")),
						Description:        String("Send a Slack notification to #alert-rule-alerts"),
					},
				},
			},
		},
		Projects:    []string{"pump-station"},
		Owner:       String("pump-station:12345"),
		DateCreated: Time(mustParseTime("2022-04-15T15:06:01.05618Z")),
	}

	require.Equal(t, expected, alertRule)
}

func TestMetricAlertService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/alert-rules/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"id": "12345",
				"name": "pump-station-alert",
				"environment": "production",
				"dataset": "transactions",
				"eventTypes": ["transaction"],
				"query": "http.url:http://service/unreadmessages",
				"aggregate": "p50(transaction.duration)",
				"timeWindow": 10,
				"thresholdType": 0,
				"resolveThreshold": 0,
				"triggers": [
				  {
					"actions": [
					  {
						"alertRuleTriggerId": "56789",
						"dateCreated": "2022-04-15T15:06:01.087054Z",
						"desc": "Send a Slack notification to #alert-rule-alerts",
						"id": "12389",
						"inputChannelId": "C0XXXFKLXXX",
						"integrationId": 111,
						"sentryAppId": null,
						"targetIdentifier": "#alert-rule-alerts",
						"targetType": "specific",
						"type": "slack"
					  }
					],
					"alertRuleId": "12345",
					"alertThreshold": 10000,
					"dateCreated": "2022-04-15T15:06:01.079598Z",
					"id": "56789",
					"label": "critical",
					"resolveThreshold": 0,
					"thresholdType": 0
				  }
				],
				"projects": [
				  "pump-station"
				],
				"owner": "pump-station:12345",
				"dateCreated": "2022-04-15T15:06:01.05618Z"
			}
		`)
	})

	params := &MetricAlert{
		Name:             String("pump-station-alert"),
		Environment:      String("production"),
		DataSet:          String("transactions"),
		Query:            String("http.url:http://service/unreadmessages"),
		Aggregate:        String("p50(transaction.duration)"),
		TimeWindow:       Float64(10.0),
		ThresholdType:    Int(0),
		ResolveThreshold: Float64(0),
		Triggers: []*MetricAlertTrigger{
			{
				ID:               String("56789"),
				AlertRuleID:      String("12345"),
				Label:            String("critical"),
				ThresholdType:    Int(0),
				AlertThreshold:   Float64(55501.0),
				ResolveThreshold: Float64(100.0),
				DateCreated:      Time(mustParseTime("2022-04-15T15:06:01.079598Z")),
				Actions: []*MetricAlertTriggerAction{
					{
						ID:                 String("12389"),
						AlertRuleTriggerID: String("56789"),
						Type:               String("slack"),
						TargetType:         String("specific"),
						TargetIdentifier:   InterfaceString("#alert-rule-alerts"),
						InputChannelID:     String("C0XXXFKLXXX"),
						IntegrationID:      Int(123),
						DateCreated:        Time(mustParseTime("2022-04-15T15:06:01.087054Z")),
						Description:        String("Send a Slack notification to #alert-rule-alerts"),
					},
				},
			},
		},
		Projects: []string{"pump-station"},
		Owner:    String("pump-station:12345"),
	}
	ctx := context.Background()
	alertRule, _, err := client.MetricAlerts.Create(ctx, "the-interstellar-jurisdiction", "pump-station", params)
	require.NoError(t, err)

	expected := &MetricAlert{
		ID:               String("12345"),
		Name:             String("pump-station-alert"),
		Environment:      String("production"),
		DataSet:          String("transactions"),
		EventTypes:       []string{"transaction"},
		Query:            String("http.url:http://service/unreadmessages"),
		Aggregate:        String("p50(transaction.duration)"),
		ThresholdType:    Int(0),
		ResolveThreshold: Float64(0),
		TimeWindow:       Float64(10.0),
		Triggers: []*MetricAlertTrigger{
			{
				ID:               String("56789"),
				AlertRuleID:      String("12345"),
				Label:            String("critical"),
				ThresholdType:    Int(0),
				AlertThreshold:   Float64(10000.0),
				ResolveThreshold: Float64(0.0),
				DateCreated:      Time(mustParseTime("2022-04-15T15:06:01.079598Z")),
				Actions: []*MetricAlertTriggerAction{
					{
						ID:                 String("12389"),
						AlertRuleTriggerID: String("56789"),
						Type:               String("slack"),
						TargetType:         String("specific"),
						TargetIdentifier:   InterfaceString("#alert-rule-alerts"),
						InputChannelID:     String("C0XXXFKLXXX"),
						IntegrationID:      Int(111),
						DateCreated:        Time(mustParseTime("2022-04-15T15:06:01.087054Z")),
						Description:        String("Send a Slack notification to #alert-rule-alerts"),
					},
				},
			},
		},
		Projects:    []string{"pump-station"},
		Owner:       String("pump-station:12345"),
		DateCreated: Time(mustParseTime("2022-04-15T15:06:01.05618Z")),
	}

	require.Equal(t, expected, alertRule)
}

func TestMetricAlertService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	params := &MetricAlert{
		ID:               String("12345"),
		Name:             String("pump-station-alert"),
		Environment:      String("production"),
		DataSet:          String("transactions"),
		Query:            String("http.url:http://service/unreadmessages"),
		Aggregate:        String("p50(transaction.duration)"),
		TimeWindow:       Float64(10),
		ThresholdType:    Int(0),
		ResolveThreshold: Float64(0),
		Triggers: []*MetricAlertTrigger{
			{
				ID:               String("6789"),
				AlertRuleID:      String("12345"),
				Label:            String("critical"),
				ThresholdType:    Int(0),
				AlertThreshold:   Float64(55501.0),
				ResolveThreshold: Float64(100.0),
				DateCreated:      Time(mustParseTime("2022-04-07T16:46:48.607583Z")),
				Actions:          []*MetricAlertTriggerAction{},
			},
		},
		Owner:       String("pump-station:12345"),
		DateCreated: Time(mustParseTime("2022-04-15T15:06:01.079598Z")),
	}

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/alert-rules/12345/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSON(t, map[string]interface{}{
			"id":               "12345",
			"name":             "pump-station-alert",
			"environment":      "production",
			"dataset":          "transactions",
			"query":            "http.url:http://service/unreadmessages",
			"aggregate":        "p50(transaction.duration)",
			"timeWindow":       json.Number("10"),
			"thresholdType":    json.Number("0"),
			"resolveThreshold": json.Number("0"),
			"triggers": []interface{}{
				map[string]interface{}{
					"id":               "6789",
					"alertRuleId":      "12345",
					"label":            "critical",
					"thresholdType":    json.Number("0"),
					"alertThreshold":   json.Number("55501"),
					"resolveThreshold": json.Number("100"),
					"dateCreated":      "2022-04-07T16:46:48.607583Z",
					"actions":          []interface{}{},
				},
			},
			"owner":       "pump-station:12345",
			"dateCreated": "2022-04-15T15:06:01.079598Z",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"id": "12345",
				"name": "pump-station-alert",
				"environment": "production",
				"dataset": "transactions",
				"eventTypes": ["transaction"],
				"query": "http.url:http://service/unreadmessages",
				"aggregate": "p50(transaction.duration)",
				"timeWindow": 10,
				"thresholdType": 0,
				"resolveThreshold": 0,
				"triggers": [
				  {
					"actions": [
						{
							"id":                 "12389",
							"alertRuleTriggerId": "56789",
							"type":               "slack",
							"targetType":         "specific",
							"targetIdentifier":   "#alert-rule-alerts",
							"inputChannelId":     "C0XXXFKLXXX",
							"integrationId":      111,
							"sentryAppId":        null,
							"dateCreated":        "2022-04-15T15:06:01.087054Z",
							"desc":               "Send a Slack notification to #alert-rule-alerts"
						}
					],
					"alertRuleId": "12345",
					"alertThreshold": 10000,
					"dateCreated": "2022-04-15T15:06:01.079598Z",
					"id": "56789",
					"label": "critical",
					"resolveThreshold": 0,
					"thresholdType": 0
				  }
				],
				"projects": [
				  "pump-station"
				],
				"owner": "pump-station:12345",
				"dateCreated": "2022-04-15T15:06:01.05618Z"
			}
		`)
	})

	ctx := context.Background()
	alertRule, _, err := client.MetricAlerts.Update(ctx, "the-interstellar-jurisdiction", "pump-station", "12345", params)
	assert.NoError(t, err)

	expected := &MetricAlert{
		ID:               String("12345"),
		Name:             String("pump-station-alert"),
		Environment:      String("production"),
		DataSet:          String("transactions"),
		EventTypes:       []string{"transaction"},
		Query:            String("http.url:http://service/unreadmessages"),
		Aggregate:        String("p50(transaction.duration)"),
		ThresholdType:    Int(0),
		ResolveThreshold: Float64(0),
		TimeWindow:       Float64(10.0),
		Triggers: []*MetricAlertTrigger{
			{
				ID:               String("56789"),
				AlertRuleID:      String("12345"),
				Label:            String("critical"),
				ThresholdType:    Int(0),
				AlertThreshold:   Float64(10000.0),
				ResolveThreshold: Float64(0.0),
				DateCreated:      Time(mustParseTime("2022-04-15T15:06:01.079598Z")),
				Actions: []*MetricAlertTriggerAction{
					{
						ID:                 String("12389"),
						AlertRuleTriggerID: String("56789"),
						Type:               String("slack"),
						TargetType:         String("specific"),
						TargetIdentifier:   InterfaceString("#alert-rule-alerts"),
						InputChannelID:     String("C0XXXFKLXXX"),
						IntegrationID:      Int(111),
						DateCreated:        Time(mustParseTime("2022-04-15T15:06:01.087054Z")),
						Description:        String("Send a Slack notification to #alert-rule-alerts"),
					},
				},
			},
		},
		Projects:    []string{"pump-station"},
		Owner:       String("pump-station:12345"),
		DateCreated: Time(mustParseTime("2022-04-15T15:06:01.05618Z")),
	}

	require.Equal(t, expected, alertRule)
}

func TestMetricAlertService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/alert-rules/12345/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	ctx := context.Background()
	_, err := client.MetricAlerts.Delete(ctx, "the-interstellar-jurisdiction", "pump-station", "12345")
	require.NoError(t, err)
}
