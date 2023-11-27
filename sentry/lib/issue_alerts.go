package sentry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// IssueAlert represents an issue alert configured for this project.
// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/serializers/models/rule.py#L131-L155
type IssueAlert struct {
	ID          *string                `json:"id,omitempty"`
	Conditions  []*IssueAlertCondition `json:"conditions,omitempty"`
	Filters     []*IssueAlertFilter    `json:"filters,omitempty"`
	Actions     []*IssueAlertAction    `json:"actions,omitempty"`
	ActionMatch *string                `json:"actionMatch,omitempty"`
	FilterMatch *string                `json:"filterMatch,omitempty"`
	Frequency   *int                   `json:"frequency,omitempty"`
	Name        *string                `json:"name,omitempty"`
	DateCreated *time.Time             `json:"dateCreated,omitempty"`
	Owner       *string                `json:"owner,omitempty"`
	CreatedBy   *IssueAlertCreatedBy   `json:"createdBy,omitempty"`
	Environment *string                `json:"environment,omitempty"`
	Projects    []string               `json:"projects,omitempty"`
	TaskUUID    *string                `json:"uuid,omitempty"` // This is actually the UUID of the async task that can be spawned to create the rule
}

// IssueAlertCreatedBy for defining the rule creator.
type IssueAlertCreatedBy struct {
	ID    *int    `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// IssueAlertCondition for defining conditions.
type IssueAlertCondition map[string]interface{}

// IssueAlertAction for defining actions.
type IssueAlertAction map[string]interface{}

// IssueAlertFilter for defining actions.
type IssueAlertFilter map[string]interface{}

// IssueAlertTaskDetail represents the inline struct Sentry defines for task details
// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/endpoints/project_rule_task_details.py#L29
type IssueAlertTaskDetail struct {
	Status *string     `json:"status,omitempty"`
	Rule   *IssueAlert `json:"rule,omitempty"`
	Error  *string     `json:"error,omitempty"`
}

// IssueAlertsService provides methods for accessing Sentry project
// client key API endpoints.
// https://docs.sentry.io/api/projects/
type IssueAlertsService service

// List issue alerts configured for a project.
func (s *IssueAlertsService) List(ctx context.Context, organizationSlug string, projectSlug string, params *ListCursorParams) ([]*IssueAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/rules/", organizationSlug, projectSlug)
	u, err := addQuery(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alerts := []*IssueAlert{}
	resp, err := s.client.Do(ctx, req, &alerts)
	if err != nil {
		return nil, resp, err
	}
	return alerts, resp, nil
}

// Get details on an issue alert.
func (s *IssueAlertsService) Get(ctx context.Context, organizationSlug string, projectSlug string, id string) (*IssueAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/rules/%v/", organizationSlug, projectSlug, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alert := new(IssueAlert)
	resp, err := s.client.Do(ctx, req, alert)
	if err != nil {
		return nil, resp, err
	}
	return alert, resp, nil
}

// Create a new issue alert bound to a project.
func (s *IssueAlertsService) Create(ctx context.Context, organizationSlug string, projectSlug string, params *IssueAlert) (*IssueAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/rules/", organizationSlug, projectSlug)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	alert := new(IssueAlert)
	resp, err := s.client.Do(ctx, req, alert)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode == 202 {
		if alert.TaskUUID == nil {
			return nil, resp, errors.New("missing task uuid")
		}
		// We just received a reference to an async task, we need to check another endpoint to retrieve the issue alert we created
		return s.getIssueAlertFromTaskDetail(ctx, organizationSlug, projectSlug, *alert.TaskUUID)
	}

	return alert, resp, nil
}

// getIssueAlertFromTaskDetail is called when Sentry offloads the issue alert creation process to an async task and sends us back the task's uuid.
// It usually doesn't happen, but when creating Slack notification rules, it seemed to be sometimes the case. During testing it
// took very long for a task to finish (10+ seconds) which is why this method can take long to return.
func (s *IssueAlertsService) getIssueAlertFromTaskDetail(ctx context.Context, organizationSlug string, projectSlug string, taskUUID string) (*IssueAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/rule-task/%v/", organizationSlug, projectSlug, taskUUID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var resp *Response
	for i := 0; i < 5; i++ {
		// TODO: Read poll interval from context
		time.Sleep(5 * time.Second)

		taskDetail := new(IssueAlertTaskDetail)
		resp, err := s.client.Do(ctx, req, taskDetail)
		if err != nil {
			return nil, resp, err
		}

		if resp.StatusCode == 404 {
			return nil, resp, fmt.Errorf("cannot find issue alert creation task with UUID %v", taskUUID)
		}
		if taskDetail.Status != nil && taskDetail.Rule != nil {
			if *taskDetail.Status == "success" {
				return taskDetail.Rule, resp, err
			} else if *taskDetail.Status == "failed" {
				if taskDetail != nil {
					return taskDetail.Rule, resp, errors.New(*taskDetail.Error)
				}

				return taskDetail.Rule, resp, errors.New("error while running the issue alert creation task")
			}
		}
	}
	return nil, resp, errors.New("getting the status of the issue alert creation from Sentry took too long")
}

// Update an issue alert.
func (s *IssueAlertsService) Update(ctx context.Context, organizationSlug string, projectSlug string, issueAlertID string, params *IssueAlert) (*IssueAlert, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/rules/%v/", organizationSlug, projectSlug, issueAlertID)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	alert := new(IssueAlert)
	resp, err := s.client.Do(ctx, req, alert)
	if err != nil {
		return nil, resp, err
	}
	return alert, resp, nil
}

// Delete an issue alert.
func (s *IssueAlertsService) Delete(ctx context.Context, organizationSlug string, projectSlug string, id string) (*Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/rules/%v/", organizationSlug, projectSlug, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
