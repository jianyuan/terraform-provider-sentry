package sentry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type MetricAlertsService service

type MetricAlert struct {
	ID               *string               `json:"id,omitempty"`
	Name             *string               `json:"name,omitempty"`
	Environment      *string               `json:"environment,omitempty"`
	DataSet          *string               `json:"dataset,omitempty"`
	EventTypes       []string              `json:"eventTypes,omitempty"`
	Query            *string               `json:"query,omitempty"`
	Aggregate        *string               `json:"aggregate,omitempty"`
	TimeWindow       *float64              `json:"timeWindow,omitempty"`
	ThresholdType    *int                  `json:"thresholdType,omitempty"`
	ResolveThreshold *float64              `json:"resolveThreshold,omitempty"`
	Triggers         []*MetricAlertTrigger `json:"triggers,omitempty"`
	Projects         []string              `json:"projects,omitempty"`
	Owner            *string               `json:"owner,omitempty"`
	DateCreated      *time.Time            `json:"dateCreated,omitempty"`
	TaskUUID         *string               `json:"uuid,omitempty"` // This is actually the UUID of the async task that can be spawned to create the metric
}

// MetricAlertTaskDetail represents the inline struct Sentry defines for task details
// https://github.com/getsentry/sentry/blob/22.12.0/src/sentry/incidents/endpoints/project_alert_rule_task_details.py#L31
type MetricAlertTaskDetail struct {
	Status    *string      `json:"status,omitempty"`
	AlertRule *MetricAlert `json:"alertRule,omitempty"`
	Error     *string      `json:"error,omitempty"`
}

// MetricAlertTrigger represents a metric alert trigger.
// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/serializers/models/alert_rule_trigger.py#L35-L47
type MetricAlertTrigger struct {
	ID               *string                     `json:"id,omitempty"`
	AlertRuleID      *string                     `json:"alertRuleId,omitempty"`
	Label            *string                     `json:"label,omitempty"`
	ThresholdType    *int                        `json:"thresholdType,omitempty"`
	AlertThreshold   *float64                    `json:"alertThreshold,omitempty"`
	ResolveThreshold *float64                    `json:"resolveThreshold,omitempty"`
	DateCreated      *time.Time                  `json:"dateCreated,omitempty"`
	Actions          []*MetricAlertTriggerAction `json:"actions"` // Must always be present.
}

// MetricAlertTriggerAction represents a metric alert trigger action.
// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/serializers/models/alert_rule_trigger_action.py#L42-L66
type MetricAlertTriggerAction struct {
	ID                 *string      `json:"id,omitempty"`
	AlertRuleTriggerID *string      `json:"alertRuleTriggerId,omitempty"`
	Type               *string      `json:"type,omitempty"`
	TargetType         *string      `json:"targetType,omitempty"`
	TargetIdentifier   *interface{} `json:"targetIdentifier,omitempty"`
	InputChannelID     *string      `json:"inputChannelId,omitempty"`
	IntegrationID      *int         `json:"integrationId,omitempty"`
	SentryAppID        *string      `json:"sentryAppId,omitempty"`
	DateCreated        *time.Time   `json:"dateCreated,omitempty"`
	Description        *string      `json:"desc,omitempty"`
}

// List Alert Rules configured for a project
func (s *MetricAlertsService) List(ctx context.Context, organizationSlug string, projectSlug string, params *ListCursorParams) ([]*MetricAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/alert-rules/", organizationSlug, projectSlug)
	u, err := addQuery(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alerts := []*MetricAlert{}
	resp, err := s.client.Do(ctx, req, &alerts)
	if err != nil {
		return nil, resp, err
	}
	return alerts, resp, nil
}

// Get details on an issue alert.
func (s *MetricAlertsService) Get(ctx context.Context, organizationSlug string, projectSlug string, id string) (*MetricAlert, *Response, error) {
	// TODO: Remove projectSlug argument
	u := fmt.Sprintf("0/organizations/%v/alert-rules/%v/", organizationSlug, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alert := new(MetricAlert)
	resp, err := s.client.Do(ctx, req, alert)
	if err != nil {
		return nil, resp, err
	}
	return alert, resp, nil
}

// Create a new Alert Rule bound to a project.
func (s *MetricAlertsService) Create(ctx context.Context, organizationSlug string, projectSlug string, params *MetricAlert) (*MetricAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/alert-rules/", organizationSlug, projectSlug)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	alert := new(MetricAlert)
	resp, err := s.client.Do(ctx, req, alert)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode == 202 {
		if alert.TaskUUID == nil {
			return nil, resp, errors.New("missing task uuid")
		}
		// We just received a reference to an async task, we need to check another endpoint to retrieve the metric alert we created
		return s.getMetricAlertFromMetricAlertTaskDetail(ctx, organizationSlug, projectSlug, *alert.TaskUUID)
	}

	return alert, resp, nil
}

// Update an Alert Rule.
func (s *MetricAlertsService) Update(ctx context.Context, organizationSlug string, projectSlug string, alertRuleID string, params *MetricAlert) (*MetricAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/alert-rules/%v/", organizationSlug, projectSlug, alertRuleID)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	alert := new(MetricAlert)
	resp, err := s.client.Do(ctx, req, alert)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode == 202 {
		if alert.TaskUUID == nil {
			return nil, resp, errors.New("missing task uuid")
		}
		// We just received a reference to an async task, we need to check another endpoint to retrieve the metric alert we created
		return s.getMetricAlertFromMetricAlertTaskDetail(ctx, organizationSlug, projectSlug, *alert.TaskUUID)
	}

	return alert, resp, nil
}

func (s *MetricAlertsService) getMetricAlertFromMetricAlertTaskDetail(ctx context.Context, organizationSlug string, projectSlug string, taskUUID string) (*MetricAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/alert-rule-task/%v/", organizationSlug, projectSlug, taskUUID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var resp *Response
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)

		taskDetail := new(MetricAlertTaskDetail)
		resp, err := s.client.Do(ctx, req, taskDetail)
		if err != nil {
			return nil, resp, err
		}

		if resp.StatusCode == 404 {
			return nil, resp, fmt.Errorf("cannot find metric alert creation task with UUID %v", taskUUID)
		}
		if taskDetail.Status != nil && taskDetail.AlertRule != nil {
			if *taskDetail.Status == "success" {
				return taskDetail.AlertRule, resp, err
			} else if *taskDetail.Status == "failed" {
				if taskDetail != nil {
					return taskDetail.AlertRule, resp, errors.New(*taskDetail.Error)
				}

				return taskDetail.AlertRule, resp, errors.New("error while running the metric alert creation task")
			}
		}
	}
	return nil, resp, errors.New("getting the status of the metric alert creation from Sentry took too long")
}

// Delete an Alert Rule.
func (s *MetricAlertsService) Delete(ctx context.Context, organizationSlug string, projectSlug string, alertRuleID string) (*Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/alert-rules/%v/", organizationSlug, projectSlug, alertRuleID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
