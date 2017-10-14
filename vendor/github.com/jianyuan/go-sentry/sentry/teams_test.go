package sentry

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamService_List(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/teams/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
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
		]`)
	})

	client := NewClient(httpClient, nil, "")
	teams, _, err := client.Teams.List("the-interstellar-jurisdiction")
	assert.NoError(t, err)

	expected := []Team{
		{
			ID:          "3",
			Slug:        "ancient-gabelers",
			Name:        "Ancient Gabelers",
			DateCreated: mustParseTime("2017-07-18T19:29:46.305Z"),
			HasAccess:   true,
			IsPending:   false,
			IsMember:    false,
		},
		{
			ID:          "2",
			Slug:        "powerful-abolitionist",
			Name:        "Powerful Abolitionist",
			DateCreated: mustParseTime("2017-07-18T19:29:24.743Z"),
			HasAccess:   true,
			IsPending:   false,
			IsMember:    false,
		},
	}
	assert.Equal(t, expected, teams)
}

func TestTeamService_Get(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/teams/the-interstellar-jurisdiction/powerful-abolitionist/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"slug": "powerful-abolitionist",
			"name": "Powerful Abolitionist",
			"hasAccess": true,
			"isPending": false,
			"dateCreated": "2017-07-18T19:29:24.743Z",
			"isMember": false,
			"organization": {
				"name": "The Interstellar Jurisdiction",
				"slug": "the-interstellar-jurisdiction",
				"avatar": {
					"avatarUuid": null,
					"avatarType": "letter_avatar"
				},
				"dateCreated": "2017-07-18T19:29:24.565Z",
				"id": "2",
				"isEarlyAdopter": false
			},
			"id": "2"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	team, _, err := client.Teams.Get("the-interstellar-jurisdiction", "powerful-abolitionist")
	assert.NoError(t, err)

	expected := &Team{
		ID:          "2",
		Slug:        "powerful-abolitionist",
		Name:        "Powerful Abolitionist",
		DateCreated: mustParseTime("2017-07-18T19:29:24.743Z"),
		HasAccess:   true,
		IsPending:   false,
		IsMember:    false,
	}
	assert.Equal(t, expected, team)
}

func TestTeamService_Create(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/teams/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "Ancient Gabelers",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"slug": "ancient-gabelers",
			"name": "Ancient Gabelers",
			"hasAccess": true,
			"isPending": false,
			"dateCreated": "2017-07-18T19:29:46.305Z",
			"isMember": false,
			"id": "3"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &CreateTeamParams{
		Name: "Ancient Gabelers",
	}
	team, _, err := client.Teams.Create("the-interstellar-jurisdiction", params)
	assert.NoError(t, err)

	expected := &Team{
		ID:          "3",
		Slug:        "ancient-gabelers",
		Name:        "Ancient Gabelers",
		DateCreated: mustParseTime("2017-07-18T19:29:46.305Z"),
		HasAccess:   true,
		IsPending:   false,
		IsMember:    false,
	}
	assert.Equal(t, expected, team)
}

func TestTeamService_Update(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/teams/the-interstellar-jurisdiction/the-obese-philosophers/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "The Inflated Philosophers",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"slug": "the-obese-philosophers",
			"name": "The Inflated Philosophers",
			"hasAccess": true,
			"isPending": false,
			"dateCreated": "2017-07-18T19:30:14.736Z",
			"isMember": false,
			"id": "4"
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &UpdateTeamParams{
		Name: "The Inflated Philosophers",
	}
	team, _, err := client.Teams.Update("the-interstellar-jurisdiction", "the-obese-philosophers", params)
	assert.NoError(t, err)
	expected := &Team{
		ID:          "4",
		Slug:        "the-obese-philosophers",
		Name:        "The Inflated Philosophers",
		DateCreated: mustParseTime("2017-07-18T19:30:14.736Z"),
		HasAccess:   true,
		IsPending:   false,
		IsMember:    false,
	}
	assert.Equal(t, expected, team)
}

func TestTeamService_Delete(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/teams/the-interstellar-jurisdiction/the-obese-philosophers/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	client := NewClient(httpClient, nil, "")
	_, err := client.Teams.Delete("the-interstellar-jurisdiction", "the-obese-philosophers")
	assert.NoError(t, err)

}
