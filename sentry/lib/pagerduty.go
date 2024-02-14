package sentry

import (
	"context"
	"fmt"
)

type PagerdutyIntegration struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ConfigData struct {
		ServiceTable []struct {
			Service        string `json:"service"`
			IntegrationKey string `json:"integration_key"`
			Id             int    `json:"id"`
		} `json:"service_table"`
	} `json:"configData"`
	ExternalId                    string `json:"externalId"`
	OrganizationId                int    `json:"organizationId"`
	OrganizationIntegrationStatus string `json:"organizationIntegrationStatus"`
}
type PagerdutyService service

func (s *PagerdutyService) Get(ctx context.Context, organization string, integrationId int) (*PagerdutyIntegration, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/integrations/%d/", organization, integrationId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pagerdutyIntegration := new(PagerdutyIntegration)
	resp, err := s.client.Do(ctx, req, pagerdutyIntegration)
	if err != nil {
		return nil, resp, err
	}
	return pagerdutyIntegration, resp, nil
}
