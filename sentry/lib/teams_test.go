package sentry

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/teams/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"id": "3",
				"slug": "ancient-gabelers",
				"name": "Ancient Gabelers",
				"dateCreated": "2017-07-18T19:29:46.305Z",
				"isMember": false,
				"teamRole": "admin",
				"hasAccess": true,
				"isPending": false,
				"memberCount": 1,
				"avatar": {
					"avatarType": "letter_avatar",
					"avatarUuid": null
				},
				"externalTeams": [],
				"projects": []
			},
			{
				"id": "2",
				"slug": "powerful-abolitionist",
				"name": "Powerful Abolitionist",
				"dateCreated": "2017-07-18T19:29:24.743Z",
				"isMember": false,
				"teamRole": "admin",
				"hasAccess": true,
				"isPending": false,
				"memberCount": 1,
				"avatar": {
					"avatarType": "letter_avatar",
					"avatarUuid": null
				},
				"externalTeams": [],
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

	ctx := context.Background()
	teams, _, err := client.Teams.List(ctx, "the-interstellar-jurisdiction")
	assert.NoError(t, err)

	expected := []*Team{
		{
			ID:          String("3"),
			Slug:        String("ancient-gabelers"),
			Name:        String("Ancient Gabelers"),
			DateCreated: Time(mustParseTime("2017-07-18T19:29:46.305Z")),
			IsMember:    Bool(false),
			TeamRole:    String("admin"),
			HasAccess:   Bool(true),
			IsPending:   Bool(false),
			MemberCount: Int(1),
			Avatar: &Avatar{
				Type: "letter_avatar",
			},
		},
		{
			ID:          String("2"),
			Slug:        String("powerful-abolitionist"),
			Name:        String("Powerful Abolitionist"),
			DateCreated: Time(mustParseTime("2017-07-18T19:29:24.743Z")),
			IsMember:    Bool(false),
			TeamRole:    String("admin"),
			HasAccess:   Bool(true),
			IsPending:   Bool(false),
			MemberCount: Int(1),
			Avatar: &Avatar{
				Type: "letter_avatar",
			},
		},
	}
	assert.Equal(t, expected, teams)
}

func TestTeamsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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

	ctx := context.Background()
	team, _, err := client.Teams.Get(ctx, "the-interstellar-jurisdiction", "powerful-abolitionist")
	assert.NoError(t, err)

	expected := &Team{
		ID:          String("2"),
		Slug:        String("powerful-abolitionist"),
		Name:        String("Powerful Abolitionist"),
		DateCreated: Time(mustParseTime("2017-07-18T19:29:24.743Z")),
		HasAccess:   Bool(true),
		IsPending:   Bool(false),
		IsMember:    Bool(false),
	}
	assert.Equal(t, expected, team)
}

func TestTeamsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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

	params := &CreateTeamParams{
		Name: String("Ancient Gabelers"),
	}
	ctx := context.Background()
	team, _, err := client.Teams.Create(ctx, "the-interstellar-jurisdiction", params)
	assert.NoError(t, err)

	expected := &Team{
		ID:          String("3"),
		Slug:        String("ancient-gabelers"),
		Name:        String("Ancient Gabelers"),
		DateCreated: Time(mustParseTime("2017-07-18T19:29:46.305Z")),
		HasAccess:   Bool(true),
		IsPending:   Bool(false),
		IsMember:    Bool(false),
	}
	assert.Equal(t, expected, team)
}

func TestTeamsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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

	params := &UpdateTeamParams{
		Name: String("The Inflated Philosophers"),
	}
	ctx := context.Background()
	team, _, err := client.Teams.Update(ctx, "the-interstellar-jurisdiction", "the-obese-philosophers", params)
	assert.NoError(t, err)
	expected := &Team{
		ID:          String("4"),
		Slug:        String("the-obese-philosophers"),
		Name:        String("The Inflated Philosophers"),
		DateCreated: Time(mustParseTime("2017-07-18T19:30:14.736Z")),
		HasAccess:   Bool(true),
		IsPending:   Bool(false),
		IsMember:    Bool(false),
	}
	assert.Equal(t, expected, team)
}

func TestTeamsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/teams/the-interstellar-jurisdiction/the-obese-philosophers/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	ctx := context.Background()
	_, err := client.Teams.Delete(ctx, "the-interstellar-jurisdiction", "the-obese-philosophers")
	assert.NoError(t, err)

}
