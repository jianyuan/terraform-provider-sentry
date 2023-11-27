package sentry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
				"avatar": {
					"avatarType": "letter_avatar",
					"avatarUuid": null
				},
				"color": "#bf6e3f",
				"dateCreated": "2018-09-20T15:48:07.592Z",
				"features": [
					"data-forwarding",
					"rate-limits"
				],
				"firstEvent": null,
				"hasAccess": true,
				"id": "4",
				"isBookmarked": false,
				"isInternal": false,
				"isMember": false,
				"isPublic": false,
				"name": "The Spoiled Yoghurt",
				"organization": {
					"avatar": {
						"avatarType": "letter_avatar",
						"avatarUuid": null
					},
					"dateCreated": "2018-09-20T15:47:52.908Z",
					"id": "2",
					"isEarlyAdopter": false,
					"name": "The Interstellar Jurisdiction",
					"require2FA": false,
					"slug": "the-interstellar-jurisdiction",
					"status": {
						"id": "active",
						"name": "active"
					}
				},
				"platform": null,
				"slug": "the-spoiled-yoghurt",
				"status": "active"
			},
			{
				"avatar": {
					"avatarType": "letter_avatar",
					"avatarUuid": null
				},
				"color": "#bf5b3f",
				"dateCreated": "2018-09-20T15:47:56.723Z",
				"features": [
					"data-forwarding",
					"rate-limits",
					"releases"
				],
				"firstEvent": null,
				"hasAccess": true,
				"id": "3",
				"isBookmarked": false,
				"isInternal": false,
				"isMember": false,
				"isPublic": false,
				"name": "Prime Mover",
				"organization": {
					"avatar": {
						"avatarType": "letter_avatar",
						"avatarUuid": null
					},
					"dateCreated": "2018-09-20T15:47:52.908Z",
					"id": "2",
					"isEarlyAdopter": false,
					"name": "The Interstellar Jurisdiction",
					"require2FA": false,
					"slug": "the-interstellar-jurisdiction",
					"status": {
						"id": "active",
						"name": "active"
					}
				},
				"platform": null,
				"slug": "prime-mover",
				"status": "active"
			},
			{
				"avatar": {
					"avatarType": "letter_avatar",
					"avatarUuid": null
				},
				"color": "#3fbf7f",
				"dateCreated": "2018-09-20T15:47:52.926Z",
				"features": [
					"data-forwarding",
					"rate-limits",
					"releases"
				],
				"firstEvent": null,
				"hasAccess": true,
				"id": "2",
				"isBookmarked": false,
				"isInternal": false,
				"isMember": false,
				"isPublic": false,
				"name": "Pump Station",
				"organization": {
					"avatar": {
						"avatarType": "letter_avatar",
						"avatarUuid": null
					},
					"dateCreated": "2018-09-20T15:47:52.908Z",
					"id": "2",
					"isEarlyAdopter": false,
					"name": "The Interstellar Jurisdiction",
					"require2FA": false,
					"slug": "the-interstellar-jurisdiction",
					"status": {
						"id": "active",
						"name": "active"
					}
				},
				"platform": null,
				"slug": "pump-station",
				"fingerprintingRules": "fingerprinting rule",
				"groupingEnhancements": "pump station grouping enhancement rule",
				"status": "active"
			}
		]`)
	})

	ctx := context.Background()
	projects, _, err := client.Projects.List(ctx)
	assert.NoError(t, err)

	expectedOrganization := Organization{
		ID:   String("2"),
		Slug: String("the-interstellar-jurisdiction"),
		Status: &OrganizationStatus{
			ID:   String("active"),
			Name: String("active"),
		},
		Name:           String("The Interstellar Jurisdiction"),
		DateCreated:    Time(mustParseTime("2018-09-20T15:47:52.908Z")),
		IsEarlyAdopter: Bool(false),
		Require2FA:     Bool(false),
		Avatar: &Avatar{
			Type: "letter_avatar",
		},
	}
	expected := []*Project{
		{
			ID:          "4",
			Slug:        "the-spoiled-yoghurt",
			Name:        "The Spoiled Yoghurt",
			Color:       "#bf6e3f",
			DateCreated: mustParseTime("2018-09-20T15:48:07.592Z"),
			Features: []string{
				"data-forwarding",
				"rate-limits",
			},
			Status:    "active",
			HasAccess: true,
			Avatar: Avatar{
				Type: "letter_avatar",
			},
			Organization: expectedOrganization,
		},
		{
			ID:           "3",
			Slug:         "prime-mover",
			Name:         "Prime Mover",
			DateCreated:  mustParseTime("2018-09-20T15:47:56.723Z"),
			IsPublic:     false,
			IsBookmarked: false,
			Color:        "#bf5b3f",
			Features: []string{
				"data-forwarding",
				"rate-limits",
				"releases",
			},
			Status:    "active",
			HasAccess: true,
			Avatar: Avatar{
				Type: "letter_avatar",
			},
			Organization: expectedOrganization,
		},
		{
			ID:           "2",
			Slug:         "pump-station",
			Name:         "Pump Station",
			DateCreated:  mustParseTime("2018-09-20T15:47:52.926Z"),
			IsPublic:     false,
			IsBookmarked: false,
			Color:        "#3fbf7f",
			Features: []string{
				"data-forwarding",
				"rate-limits",
				"releases",
			},
			Status:    "active",
			HasAccess: true,
			Avatar: Avatar{
				Type: "letter_avatar",
			},
			Organization:         expectedOrganization,
			FingerprintingRules:  "fingerprinting rule",
			GroupingEnhancements: "pump station grouping enhancement rule",
		},
	}
	assert.Equal(t, expected, projects)
}

func TestProjectsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"allowedDomains": [
				"*"
			],
			"avatar": {
				"avatarType": "letter_avatar",
				"avatarUuid": null
			},
			"color": "#3fbf7f",
			"dataScrubber": true,
			"dataScrubberDefaults": true,
			"dateCreated": "2018-10-02T14:19:09.864Z",
			"defaultEnvironment": null,
			"digestsMaxDelay": 1800,
			"digestsMinDelay": 300,
			"features": [
				"data-forwarding",
				"rate-limits",
				"releases"
			],
			"firstEvent": null,
			"hasAccess": true,
			"id": "2",
			"isBookmarked": false,
			"isInternal": false,
			"isMember": false,
			"isPublic": false,
			"latestRelease": {
				"authors": [],
				"commitCount": 0,
				"data": {},
				"dateCreated": "2018-10-02T14:19:25.397Z",
				"dateReleased": null,
				"deployCount": 0,
				"firstEvent": null,
				"lastCommit": null,
				"lastDeploy": null,
				"lastEvent": null,
				"newGroups": 0,
				"owner": null,
				"projects": [{
					"name": "Pump Station",
					"slug": "pump-station"
				}],
				"ref": "6ba09a7c53235ee8a8fa5ee4c1ca8ca886e7fdbb",
				"shortVersion": "2.0rc2",
				"url": null,
				"version": "2.0rc2"
			},
			"name": "Pump Station",
			"options": {
				"feedback:branding": true,
				"filters:blacklisted_ips": "",
				"filters:error_messages": "",
				"filters:releases": "",
				"sentry:csp_ignored_sources": "",
				"sentry:csp_ignored_sources_defaults": true,
				"sentry:reprocessing_active": false
			},
			"organization": {
				"avatar": {
					"avatarType": "letter_avatar",
					"avatarUuid": null
				},
				"dateCreated": "2018-10-02T14:19:09.817Z",
				"id": "2",
				"isEarlyAdopter": false,
				"name": "The Interstellar Jurisdiction",
				"require2FA": false,
				"slug": "the-interstellar-jurisdiction",
				"status": {
					"id": "active",
					"name": "active"
				}
			},
			"platform": null,
			"platforms": [],
			"plugins": [{
				"assets": [],
				"author": {
					"name": "Sentry Team",
					"url": "https://github.com/getsentry/sentry"
				},
				"canDisable": true,
				"contexts": [],
				"description": "Integrates web hooks.",
				"doc": "",
				"enabled": false,
				"hasConfiguration": true,
				"id": "webhooks",
				"isTestable": true,
				"metadata": {},
				"name": "WebHooks",
				"resourceLinks": [{
						"title": "Bug Tracker",
						"url": "https://github.com/getsentry/sentry/issues"
					},
					{
						"title": "Source",
						"url": "https://github.com/getsentry/sentry"
					}
				],
				"shortName": "WebHooks",
				"slug": "webhooks",
				"status": "unknown",
				"type": "notification",
				"version": "9.1.0.dev0"
			}],
			"processingIssues": 0,
			"relayPiiConfig": null,
			"resolveAge": 720,
			"safeFields": [],
			"scrapeJavaScript": true,
			"scrubIPAddresses": false,
			"securityToken": "320e3180c64e11e8b61e0242ac110002",
			"securityTokenHeader": null,
			"sensitiveFields": [],
			"slug": "pump-station",
			"status": "active",
			"storeCrashReports": false,
			"subjectPrefix": "[Sentry] ",
			"subjectTemplate": "$shortID - $title",
			"team": {
				"id": "2",
				"name": "Powerful Abolitionist",
				"slug": "powerful-abolitionist"
			},
			"teams": [{
				"id": "2",
				"name": "Powerful Abolitionist",
				"slug": "powerful-abolitionist"
			}],
			"fingerprintingRules": "fingerprinting rule",
			"groupingEnhancements": "pump-station grouping enhancement rule",
			"verifySSL": false
		}`)
	})

	ctx := context.Background()
	project, _, err := client.Projects.Get(ctx, "the-interstellar-jurisdiction", "pump-station")
	assert.NoError(t, err)

	expectedOrganization := Organization{
		ID:   String("2"),
		Slug: String("the-interstellar-jurisdiction"),
		Status: &OrganizationStatus{
			ID:   String("active"),
			Name: String("active"),
		},
		Name:           String("The Interstellar Jurisdiction"),
		DateCreated:    Time(mustParseTime("2018-10-02T14:19:09.817Z")),
		IsEarlyAdopter: Bool(false),
		Require2FA:     Bool(false),
		Avatar: &Avatar{
			Type: "letter_avatar",
		},
	}
	expected := &Project{
		ID:          "2",
		Slug:        "pump-station",
		Name:        "Pump Station",
		Color:       "#3fbf7f",
		DateCreated: mustParseTime("2018-10-02T14:19:09.864Z"),
		Features: []string{
			"data-forwarding",
			"rate-limits",
			"releases",
		},
		Status:    "active",
		HasAccess: true,
		Avatar: Avatar{
			Type: "letter_avatar",
		},
		Options: map[string]interface{}{
			"feedback:branding":                   true,
			"filters:blacklisted_ips":             "",
			"filters:error_messages":              "",
			"filters:releases":                    "",
			"sentry:csp_ignored_sources":          "",
			"sentry:csp_ignored_sources_defaults": true,
			"sentry:reprocessing_active":          false,
		},
		DigestsMinDelay:      300,
		DigestsMaxDelay:      1800,
		ResolveAge:           720,
		SubjectPrefix:        "[Sentry] ",
		AllowedDomains:       []string{"*"},
		DataScrubber:         true,
		DataScrubberDefaults: true,
		FingerprintingRules:  "fingerprinting rule",
		GroupingEnhancements: "pump-station grouping enhancement rule",
		SafeFields:           []string{},
		SensitiveFields:      []string{},
		SubjectTemplate:      "$shortID - $title",
		SecurityToken:        "320e3180c64e11e8b61e0242ac110002",
		ScrapeJavaScript:     true,
		Organization:         expectedOrganization,
		Team: Team{
			ID:   String("2"),
			Slug: String("powerful-abolitionist"),
			Name: String("Powerful Abolitionist"),
		},
		Teams: []Team{
			{
				ID:   String("2"),
				Slug: String("powerful-abolitionist"),
				Name: String("Powerful Abolitionist"),
			},
		},
	}
	assert.Equal(t, expected, project)
}

func TestProjectsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/teams/the-interstellar-jurisdiction/powerful-abolitionist/projects/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "The Spoiled Yoghurt",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
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
		}`)
	})

	params := &CreateProjectParams{
		Name: "The Spoiled Yoghurt",
	}
	ctx := context.Background()
	project, _, err := client.Projects.Create(ctx, "the-interstellar-jurisdiction", "powerful-abolitionist", params)
	assert.NoError(t, err)

	expected := &Project{
		ID:           "4",
		Slug:         "the-spoiled-yoghurt",
		Name:         "The Spoiled Yoghurt",
		DateCreated:  mustParseTime("2017-07-18T19:29:44.996Z"),
		IsPublic:     false,
		IsBookmarked: false,
		Color:        "#bf6e3f",
		Features: []string{
			"data-forwarding",
			"rate-limits",
		},
		Status: "active",
	}
	assert.Equal(t, expected, project)
}

func TestProjectsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/plain-proxy/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "Plane Proxy",
			"slug": "plane-proxy",
			"options": map[string]interface{}{
				"sentry:origins": "http://example.com\nhttp://example.invalid",
			},
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"status": "active",
			"digestsMinDelay": 300,
			"options": {
				"sentry:origins": "http://example.com\nhttp://example.invalid",
				"sentry:resolve_age": 720
			},
			"defaultEnvironment": null,
			"features": [
				"data-forwarding",
				"rate-limits",
				"releases"
			],
			"subjectPrefix": null,
			"color": "#bf803f",
			"slug": "plane-proxy",
			"isPublic": false,
			"dateCreated": "2017-07-18T19:30:09.751Z",
			"platforms": [],
			"callSign": "PLANE-PROXY",
			"firstEvent": null,
			"digestsMaxDelay": 1800,
			"resolveAge": 720,
			"processingIssues": 0,
			"isBookmarked": false,
			"callSignReviewed": false,
			"id": "5",
			"subjectTemplate": "[$project] ${tag:level}: $title",
			"fingerprintingRules": "fingerprinting rule",
			"groupingEnhancements": "Plane Proxy grouping enhancement rule",
			"name": "Plane Proxy"
		}`)
	})

	params := &UpdateProjectParams{
		Name: "Plane Proxy",
		Slug: "plane-proxy",
		Options: map[string]interface{}{
			"sentry:origins": "http://example.com\nhttp://example.invalid",
		},
	}
	ctx := context.Background()
	project, _, err := client.Projects.Update(ctx, "the-interstellar-jurisdiction", "plain-proxy", params)
	assert.NoError(t, err)
	expected := &Project{
		ID:           "5",
		Slug:         "plane-proxy",
		Name:         "Plane Proxy",
		DateCreated:  mustParseTime("2017-07-18T19:30:09.751Z"),
		IsPublic:     false,
		IsBookmarked: false,
		Color:        "#bf803f",
		Features: []string{
			"data-forwarding",
			"rate-limits",
			"releases",
		},
		Status: "active",
		Options: map[string]interface{}{
			"sentry:origins":     "http://example.com\nhttp://example.invalid",
			"sentry:resolve_age": json.Number("720"),
		},
		DigestsMinDelay:      300,
		DigestsMaxDelay:      1800,
		ResolveAge:           720,
		SubjectTemplate:      "[$project] ${tag:level}: $title",
		FingerprintingRules:  "fingerprinting rule",
		GroupingEnhancements: "Plane Proxy grouping enhancement rule",
	}
	assert.Equal(t, expected, project)
}

func TestProjectsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/plain-proxy/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	ctx := context.Background()
	_, err := client.Projects.Delete(ctx, "the-interstellar-jurisdiction", "plain-proxy")
	assert.NoError(t, err)

}

func TestProjectsService_UpdateTeam(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/teams/planet-express/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"slug": "plane-proxy",
			"id": "5",
			"name": "Plane Proxy",
			"team": {
				"id": "420",
				"name": "Planet Express",
				"slug": "planet-express"
			}
		}`)
	})

	ctx := context.Background()
	project, _, err := client.Projects.AddTeam(ctx, "the-interstellar-jurisdiction", "pump-station", "planet-express")
	assert.NoError(t, err)
	expected := &Project{
		ID:   "5",
		Slug: "plane-proxy",
		Name: "Plane Proxy",
		Team: Team{
			ID:   String("420"),
			Slug: String("planet-express"),
			Name: String("Planet Express"),
		},
	}
	assert.Equal(t, expected, project)
}

func TestProjectsService_DeleteTeam(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/teams/powerful-abolitionist/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	ctx := context.Background()
	_, err := client.Projects.RemoveTeam(ctx, "the-interstellar-jurisdiction", "pump-station", "powerful-abolitionist")
	assert.NoError(t, err)
}
