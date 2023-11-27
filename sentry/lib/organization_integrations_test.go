package sentry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationIntegrationsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/integrations/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		assertQuery(t, map[string]string{"cursor": "100:-1:1", "provider_key": "github"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"id": "123456",
				"name": "octocat",
				"icon": "https://avatars.githubusercontent.com/u/583231?v=4",
				"domainName": "github.com/octocat",
				"accountType": "Organization",
				"scopes": ["read", "write"],
				"status": "active",
				"provider": {
					"key": "github",
					"slug": "github",
					"name": "GitHub",
					"canAdd": true,
					"canDisable": false,
					"features": [
						"codeowners",
						"commits",
						"issue-basic",
						"stacktrace-link"
					],
					"aspects": {}
				},
				"configOrganization": [],
				"configData": {},
				"externalId": "87654321",
				"organizationId": 2,
				"organizationIntegrationStatus": "active",
				"gracePeriodEnd": null
			}
		]`)
	})

	ctx := context.Background()
	integrations, _, err := client.OrganizationIntegrations.List(ctx, "the-interstellar-jurisdiction", &ListOrganizationIntegrationsParams{
		ListCursorParams: ListCursorParams{
			Cursor: "100:-1:1",
		},
		ProviderKey: "github",
	})
	assert.NoError(t, err)
	expected := []*OrganizationIntegration{
		{
			ID:          "123456",
			Name:        "octocat",
			Icon:        String("https://avatars.githubusercontent.com/u/583231?v=4"),
			DomainName:  "github.com/octocat",
			AccountType: String("Organization"),
			Scopes:      []string{"read", "write"},
			Status:      "active",
			Provider: OrganizationIntegrationProvider{
				Key:        "github",
				Slug:       "github",
				Name:       "GitHub",
				CanAdd:     true,
				CanDisable: false,
				Features: []string{
					"codeowners",
					"commits",
					"issue-basic",
					"stacktrace-link",
				},
			},
			ConfigData:                    &IntegrationConfigData{},
			ExternalId:                    "87654321",
			OrganizationId:                2,
			OrganizationIntegrationStatus: "active",
			GracePeriodEnd:                nil,
		},
	}
	assert.Equal(t, expected, integrations)
}

func TestOrganizationIntegrationsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/integrations/456789/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
		  "id": "456789",
		  "name": "Interstellar PagerDuty",
		  "icon": null,
		  "domainName": "the-interstellar-jurisdiction",
		  "accountType": null,
		  "scopes": null,
		  "status": "active",
		  "provider": {
			"key": "pagerduty",
			"slug": "pagerduty",
			"name": "PagerDuty",
			"canAdd": true,
			"canDisable": false,
			"features": [
			  "alert-rule",
			  "incident-management"
			],
			"aspects": {
			  "alerts": [
				{
				  "type": "info",
				  "text": "The PagerDuty integration adds a new Alert Rule action to all projects. To enable automatic notifications sent to PagerDuty you must create a rule using the PagerDuty action in your project settings."
				}
			  ]
			}
		  },
		  "configOrganization": [
			{
			  "name": "service_table",
			  "type": "table",
			  "label": "PagerDuty services with the Sentry integration enabled",
			  "help": "If services need to be updated, deleted, or added manually please do so here. Alert rules will need to be individually updated for any additions or deletions of services.",
			  "addButtonText": "",
			  "columnLabels": {
				"service": "Service",
				"integration_key": "Integration Key"
			  },
			  "columnKeys": [
				"service",
				"integration_key"
			  ],
			  "confirmDeleteMessage": "Any alert rules associated with this service will stop working. The rules will still exist but will show a removed service."
			}
		  ],
		  "configData": {
			"service_table": [
			  {
				"service": "testing123",
				"integration_key": "abc123xyz",
				"id": 22222
			  }
			]
		  },
		  "externalId": "999999",
		  "organizationId": 2,
		  "organizationIntegrationStatus": "active",
		  "gracePeriodEnd": null
		}`)
	})

	ctx := context.Background()
	integration, _, err := client.OrganizationIntegrations.Get(ctx, "the-interstellar-jurisdiction", "456789")
	assert.NoError(t, err)
	expected := OrganizationIntegration{
		ID:          "456789",
		Name:        "Interstellar PagerDuty",
		Icon:        nil,
		DomainName:  "the-interstellar-jurisdiction",
		AccountType: nil,
		Scopes:      nil,
		Status:      "active",
		Provider: OrganizationIntegrationProvider{
			Key:        "pagerduty",
			Slug:       "pagerduty",
			Name:       "PagerDuty",
			CanAdd:     true,
			CanDisable: false,
			Features: []string{
				"alert-rule",
				"incident-management",
			},
		},
		ConfigData: &IntegrationConfigData{
			"service_table": []interface{}{
				map[string]interface{}{
					"service":         "testing123",
					"integration_key": "abc123xyz",
					"id":              json.Number("22222"),
				},
			},
		},
		ExternalId:                    "999999",
		OrganizationId:                2,
		OrganizationIntegrationStatus: "active",
		GracePeriodEnd:                nil,
	}
	assert.Equal(t, &expected, integration)
}

func TestOrganizationIntegrationsService_UpdateConfig(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/integrations/456789/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		w.Header().Set("Content-Type", "application/json")
	})

	updateConfigOrganizationIntegrationsParams := UpdateConfigOrganizationIntegrationsParams{
		"service_table": []interface{}{
			map[string]interface{}{
				"service":         "testing123",
				"integration_key": "abc123xyz",
				"id":              json.Number("22222"),
			},
			map[string]interface{}{
				"service":         "testing456",
				"integration_key": "efg456lmn",
				"id":              "",
			},
		},
	}
	ctx := context.Background()
	resp, err := client.OrganizationIntegrations.UpdateConfig(ctx, "the-interstellar-jurisdiction", "456789", &updateConfigOrganizationIntegrationsParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.ContentLength)
}
