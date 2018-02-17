package sentry

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationService_List(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		assertQuery(t, map[string]string{"cursor": "1500300636142:0:1"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"name": "The Interstellar Jurisdiction",
				"slug": "the-interstellar-jurisdiction",
				"avatar": {
					"avatarUuid": null,
					"avatarType": "letter_avatar"
				},
				"dateCreated": "2017-07-17T14:10:36.141Z",
				"id": "2",
				"isEarlyAdopter": false
			}
		]`)
	})

	client := NewClient(httpClient, nil, "")
	organizations, _, err := client.Organizations.List(&ListOrganizationParams{
		Cursor: "1500300636142:0:1",
	})
	assert.NoError(t, err)
	expected := []Organization{
		{
			ID:             "2",
			Slug:           "the-interstellar-jurisdiction",
			Name:           "The Interstellar Jurisdiction",
			DateCreated:    mustParseTime("2017-07-17T14:10:36.141Z"),
			IsEarlyAdopter: false,
			Avatar: OrganizationAvatar{
				UUID: nil,
				Type: "letter_avatar",
			},
		},
	}
	assert.Equal(t, expected, organizations)
}

func TestOrganizationService_Get(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"defaultRole": "member",
			"features": [
				"sso",
				"callsigns",
				"api-keys",
				"open-membership",
				"shared-issues"
			],
			"safeFields": [],
			"id": "2",
			"isEarlyAdopter": false,
			"scrubIPAddresses": false,
			"access": [],
			"allowSharedIssues": true,
			"isDefault": false,
			"sensitiveFields": [],
			"quota": {
				"maxRateInterval": 60,
				"projectLimit": 100,
				"accountLimit": 0,
				"maxRate": 0
			},
			"dateCreated": "2017-07-18T19:29:24.565Z",
			"slug": "the-interstellar-jurisdiction",
			"openMembership": true,
			"availableRoles": [
				{
					"id": "member",
					"name": "Member"
				},
				{
					"id": "admin",
					"name": "Admin"
				},
				{
					"id": "manager",
					"name": "Manager"
				},
				{
					"id": "owner",
					"name": "Owner"
				}
			],
			"name": "The Interstellar Jurisdiction",
			"enhancedPrivacy": false,
			"teams": [
				{
					"slug": "ancient-gabelers",
					"name": "Ancient Gabelers",
					"hasAccess": true,
					"isPending": false,
					"dateCreated": "2017-07-18T19:29:46.305Z",
					"isMember": false,
					"id": "3",
					"projects": []
				},
				{
					"slug": "powerful-abolitionist",
					"name": "Powerful Abolitionist",
					"hasAccess": true,
					"isPending": false,
					"dateCreated": "2017-07-18T19:29:24.743Z",
					"isMember": false,
					"id": "2",
					"projects": [
						{
							"status": "active",
							"slug": "prime-mover",
							"defaultEnvironment": null,
							"features": [
								"data-forwarding",
								"rate-limits",
								"releases"
							],
							"color": "#bf5b3f",
							"isPublic": false,
							"dateCreated": "2017-07-18T19:29:30.063Z",
							"platforms": [],
							"callSign": "PRIME-MOVER",
							"firstEvent": null,
							"processingIssues": 0,
							"isBookmarked": false,
							"callSignReviewed": false,
							"id": "3",
							"name": "Prime Mover"
						},
						{
							"status": "active",
							"slug": "pump-station",
							"defaultEnvironment": null,
							"features": [
								"data-forwarding",
								"rate-limits",
								"releases"
							],
							"color": "#3fbf7f",
							"isPublic": false,
							"dateCreated": "2017-07-18T19:29:24.793Z",
							"platforms": [],
							"callSign": "PUMP-STATION",
							"firstEvent": null,
							"processingIssues": 0,
							"isBookmarked": false,
							"callSignReviewed": false,
							"id": "2",
							"name": "Pump Station"
						},
						{
							"status": "active",
							"slug": "the-spoiled-yoghurt",
							"defaultEnvironment": null,
							"features": [
								"data-forwarding",
								"rate-limits"
							],
							"color": "#bf6e3f",
							"isPublic": false,
							"dateCreated": "2017-07-18T19:29:44.996Z",
							"platforms": [],
							"callSign": "THE-SPOILED-YOGHURT",
							"firstEvent": null,
							"processingIssues": 0,
							"isBookmarked": false,
							"callSignReviewed": false,
							"id": "4",
							"name": "The Spoiled Yoghurt"
						}
					]
				}
			],
			"pendingAccessRequests": 0,
			"dataScrubberDefaults": false,
			"dataScrubber": false,
			"avatar": {
				"avatarUuid": null,
				"avatarType": "letter_avatar"
			},
			"onboardingTasks": [
				{
					"status": "complete",
					"dateCompleted": "2017-07-18T19:29:45.084Z",
					"task": 1,
					"data": {},
					"user": null
				}
			]
		}`)
	})

	client := NewClient(httpClient, nil, "")
	organization, _, err := client.Organizations.Get("the-interstellar-jurisdiction")
	assert.NoError(t, err)
	expected := &Organization{
		ID:          "2",
		Slug:        "the-interstellar-jurisdiction",
		Name:        "The Interstellar Jurisdiction",
		DateCreated: mustParseTime("2017-07-18T19:29:24.565Z"),
		Quota: OrganizationQuota{
			MaxRate:         0,
			MaxRateInterval: 60,
			AccountLimit:    0,
			ProjectLimit:    100,
		},
		Access: []string{},
		Features: []string{
			"sso",
			"callsigns",
			"api-keys",
			"open-membership",
			"shared-issues",
		},
		PendingAccessRequests: 0,
		IsDefault:             false,
		DefaultRole:           "member",
		AvailableRoles: []OrganizationAvailableRole{
			{
				ID:   "member",
				Name: "Member",
			},
			{
				ID:   "admin",
				Name: "Admin",
			},
			{
				ID:   "manager",
				Name: "Manager",
			},
			{
				ID:   "owner",
				Name: "Owner",
			},
		},
		AccountRateLimit: 0,
		ProjectRateLimit: 0,
		Avatar: OrganizationAvatar{
			Type: "letter_avatar",
		},
		OpenMembership:       true,
		AllowSharedIssues:    true,
		EnhancedPrivacy:      false,
		DataScrubber:         false,
		DataScrubberDefaults: false,
		SensitiveFields:      []string{},
		SafeFields:           []string{},
		ScrubIPAddresses:     false,
		IsEarlyAdopter:       false,
	}
	assert.Equal(t, expected, organization)
}

func TestOrganizationService_Create(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "The Interstellar Jurisdiction",
			"slug": "the-interstellar-jurisdiction",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"name": "The Interstellar Jurisdiction",
			"slug": "the-interstellar-jurisdiction",
			"id": "2"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &CreateOrganizationParams{
		Name: "The Interstellar Jurisdiction",
		Slug: "the-interstellar-jurisdiction",
	}
	organization, _, err := client.Organizations.Create(params)
	assert.NoError(t, err)
	expected := &Organization{
		ID:   "2",
		Name: "The Interstellar Jurisdiction",
		Slug: "the-interstellar-jurisdiction",
	}
	assert.Equal(t, expected, organization)
}

func TestOrganizationService_Create_AgreeTerms(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSON(t, map[string]interface{}{
			"name":       "The Interstellar Jurisdiction",
			"slug":       "the-interstellar-jurisdiction",
			"agreeTerms": true,
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"name": "The Interstellar Jurisdiction",
			"slug": "the-interstellar-jurisdiction",
			"id": "2"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &CreateOrganizationParams{
		Name:       "The Interstellar Jurisdiction",
		Slug:       "the-interstellar-jurisdiction",
		AgreeTerms: Bool(true),
	}
	organization, _, err := client.Organizations.Create(params)
	assert.NoError(t, err)
	expected := &Organization{
		ID:   "2",
		Name: "The Interstellar Jurisdiction",
		Slug: "the-interstellar-jurisdiction",
	}
	assert.Equal(t, expected, organization)
}

func TestOrganizationService_Update(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/badly-misnamed/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "Impeccably Designated",
			"slug": "impeccably-designated",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"name": "Impeccably Designated",
			"slug": "impeccably-designated",
			"id": "2"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &UpdateOrganizationParams{
		Name: "Impeccably Designated",
		Slug: "impeccably-designated",
	}
	organization, _, err := client.Organizations.Update("badly-misnamed", params)
	assert.NoError(t, err)
	expected := &Organization{
		ID:   "2",
		Name: "Impeccably Designated",
		Slug: "impeccably-designated",
	}
	assert.Equal(t, expected, organization)
}

func TestOrganizationService_Delete(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	client := NewClient(httpClient, nil, "")
	_, err := client.Organizations.Delete("the-interstellar-jurisdiction")
	assert.NoError(t, err)
}
