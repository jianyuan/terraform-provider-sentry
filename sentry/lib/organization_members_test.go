package sentry

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationMembersService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/members/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		assertQuery(t, map[string]string{"cursor": "100:-1:1"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"inviteStatus": "approved",
				"dateCreated": "2020-01-04T00:00:00.000000Z",
				"user": {
					"username": "test@example.com",
					"lastLogin": "2020-01-02T00:00:00.000000Z",
					"isSuperuser": false,
					"emails": [
						{
							"is_verified": true,
							"id": "1",
							"email": "test@example.com"
						}
					],
					"isManaged": false,
					"experiments": {},
					"lastActive": "2020-01-03T00:00:00.000000Z",
					"isStaff": false,
					"identities": [],
					"id": "1",
					"isActive": true,
					"has2fa": false,
					"name": "John Doe",
					"avatarUrl": "https://secure.gravatar.com/avatar/55502f40dc8b7c769880b10874abc9d0?s=32&d=mm",
					"dateJoined": "2020-01-01T00:00:00.000000Z",
					"options": {
						"timezone": "UTC",
						"stacktraceOrder": -1,
						"language": "en",
						"clock24Hours": false
					},
					"flags": {
						"newsletter_consent_prompt": false
					},
					"avatar": {
						"avatarUuid": null,
						"avatarType": "letter_avatar"
					},
					"hasPasswordAuth": true,
					"email": "test@example.com"
				},
				"roleName": "Owner",
				"expired": false,
				"id": "1",
				"inviterName": null,
				"name": "John Doe",
				"role": "owner",
				"flags": {
					"sso:linked": false,
					"sso:invalid": false
				},
				"email": "test@example.com",
				"pending": false
			}
		]`)
	})

	ctx := context.Background()
	members, _, err := client.OrganizationMembers.List(ctx, "the-interstellar-jurisdiction", &ListCursorParams{
		Cursor: "100:-1:1",
	})
	assert.NoError(t, err)
	expected := []*OrganizationMember{
		{
			ID:    "1",
			Email: "test@example.com",
			Name:  "John Doe",
			User: User{
				ID:              "1",
				Name:            "John Doe",
				Username:        "test@example.com",
				Email:           "test@example.com",
				AvatarURL:       "https://secure.gravatar.com/avatar/55502f40dc8b7c769880b10874abc9d0?s=32&d=mm",
				IsActive:        true,
				HasPasswordAuth: true,
				IsManaged:       false,
				DateJoined:      mustParseTime("2020-01-01T00:00:00.000000Z"),
				LastLogin:       mustParseTime("2020-01-02T00:00:00.000000Z"),
				Has2FA:          false,
				LastActive:      mustParseTime("2020-01-03T00:00:00.000000Z"),
				IsSuperuser:     false,
				IsStaff:         false,
				Avatar: Avatar{
					Type: "letter_avatar",
					UUID: nil,
				},
				Emails: []UserEmail{
					{
						ID:         "1",
						Email:      "test@example.com",
						IsVerified: true,
					},
				},
			},
			Role:     "owner",
			RoleName: "Owner",
			Pending:  false,
			Expired:  false,
			Flags: map[string]bool{
				"sso:invalid": false,
				"sso:linked":  false,
			},
			DateCreated:  mustParseTime("2020-01-04T00:00:00.000000Z"),
			InviteStatus: "approved",
			InviterName:  nil,
		},
	}
	assert.Equal(t, expected, members)
}

func TestOrganizationMembersService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/members/1/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
				"inviteStatus": "approved",
				"dateCreated": "2020-01-04T00:00:00.000000Z",
				"user": {
					"username": "test@example.com",
					"lastLogin": "2020-01-02T00:00:00.000000Z",
					"isSuperuser": false,
					"emails": [
						{
							"is_verified": true,
							"id": "1",
							"email": "test@example.com"
						}
					],
					"isManaged": false,
					"experiments": {},
					"lastActive": "2020-01-03T00:00:00.000000Z",
					"isStaff": false,
					"identities": [],
					"id": "1",
					"isActive": true,
					"has2fa": false,
					"name": "John Doe",
					"avatarUrl": "https://secure.gravatar.com/avatar/55502f40dc8b7c769880b10874abc9d0?s=32&d=mm",
					"dateJoined": "2020-01-01T00:00:00.000000Z",
					"options": {
						"timezone": "UTC",
						"stacktraceOrder": -1,
						"language": "en",
						"clock24Hours": false
					},
					"flags": {
						"newsletter_consent_prompt": false
					},
					"avatar": {
						"avatarUuid": null,
						"avatarType": "letter_avatar"
					},
					"hasPasswordAuth": true,
					"email": "test@example.com"
				},
				"roleName": "Owner",
				"expired": false,
				"id": "1",
				"inviterName": null,
				"name": "John Doe",
				"role": "owner",
				"flags": {
					"sso:linked": false,
					"sso:invalid": false
				},
				"teams": [],
				"email": "test@example.com",
				"pending": false
			}`)
	})

	ctx := context.Background()
	members, _, err := client.OrganizationMembers.Get(ctx, "the-interstellar-jurisdiction", "1")
	assert.NoError(t, err)
	expected := OrganizationMember{
		ID:    "1",
		Email: "test@example.com",
		Name:  "John Doe",
		User: User{
			ID:              "1",
			Name:            "John Doe",
			Username:        "test@example.com",
			Email:           "test@example.com",
			AvatarURL:       "https://secure.gravatar.com/avatar/55502f40dc8b7c769880b10874abc9d0?s=32&d=mm",
			IsActive:        true,
			HasPasswordAuth: true,
			IsManaged:       false,
			DateJoined:      mustParseTime("2020-01-01T00:00:00.000000Z"),
			LastLogin:       mustParseTime("2020-01-02T00:00:00.000000Z"),
			Has2FA:          false,
			LastActive:      mustParseTime("2020-01-03T00:00:00.000000Z"),
			IsSuperuser:     false,
			IsStaff:         false,
			Avatar: Avatar{
				Type: "letter_avatar",
				UUID: nil,
			},
			Emails: []UserEmail{
				{
					ID:         "1",
					Email:      "test@example.com",
					IsVerified: true,
				},
			},
		},
		Role:     "owner",
		RoleName: "Owner",
		Pending:  false,
		Expired:  false,
		Flags: map[string]bool{
			"sso:invalid": false,
			"sso:linked":  false,
		},
		Teams:        []string{},
		DateCreated:  mustParseTime("2020-01-04T00:00:00.000000Z"),
		InviteStatus: "approved",
		InviterName:  nil,
	}
	assert.Equal(t, &expected, members)
}

func TestOrganizationMembersService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/members/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "1",
			"email": "test@example.com",
			"name": "test@example.com",
			"user": null,
			"role": "member",
			"roleName": "Member",
			"pending": true,
			"expired": false,
			"flags": {
				"sso:linked": false,
				"sso:invalid": false,
				"member-limit:restricted": false
			},
			"teams": [],
			"dateCreated": "2020-01-01T00:00:00.000000Z",
			"inviteStatus": "approved",
			"inviterName": "John Doe"
		}`)
	})

	createOrganizationMemberParams := CreateOrganizationMemberParams{
		Email: "test@example.com",
		Role:  RoleMember,
	}
	ctx := context.Background()
	member, _, err := client.OrganizationMembers.Create(ctx, "the-interstellar-jurisdiction", &createOrganizationMemberParams)
	assert.NoError(t, err)

	inviterName := "John Doe"
	expected := OrganizationMember{
		ID:       "1",
		Email:    "test@example.com",
		Name:     "test@example.com",
		User:     User{},
		Role:     "member",
		RoleName: "Member",
		Pending:  true,
		Expired:  false,
		Flags: map[string]bool{
			"sso:linked":              false,
			"sso:invalid":             false,
			"member-limit:restricted": false,
		},
		Teams:        []string{},
		DateCreated:  mustParseTime("2020-01-01T00:00:00.000000Z"),
		InviteStatus: "approved",
		InviterName:  &inviterName,
	}

	assert.Equal(t, &expected, member)
}

func TestOrganizationMembersService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/members/1/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "1",
			"email": "test@example.com",
			"name": "test@example.com",
			"user": null,
			"role": "manager",
			"roleName": "Manager",
			"pending": true,
			"expired": false,
			"flags": {
				"sso:linked": false,
				"sso:invalid": false,
				"member-limit:restricted": false
			},
			"teams": [],
			"dateCreated": "2020-01-01T00:00:00.000000Z",
			"inviteStatus": "approved",
			"inviterName": "John Doe"
		}`)
	})

	updateOrganizationMemberParams := UpdateOrganizationMemberParams{
		Role: RoleMember,
	}
	ctx := context.Background()
	member, _, err := client.OrganizationMembers.Update(ctx, "the-interstellar-jurisdiction", "1", &updateOrganizationMemberParams)
	assert.NoError(t, err)

	inviterName := "John Doe"
	expected := OrganizationMember{
		ID:       "1",
		Email:    "test@example.com",
		Name:     "test@example.com",
		User:     User{},
		Role:     "manager",
		RoleName: "Manager",
		Pending:  true,
		Expired:  false,
		Flags: map[string]bool{
			"sso:linked":              false,
			"sso:invalid":             false,
			"member-limit:restricted": false,
		},
		Teams:        []string{},
		DateCreated:  mustParseTime("2020-01-01T00:00:00.000000Z"),
		InviteStatus: "approved",
		InviterName:  &inviterName,
	}

	assert.Equal(t, &expected, member)
}

func TestOrganizationMembersService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/members/1/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.OrganizationMembers.Delete(ctx, "the-interstellar-jurisdiction", "1")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.ContentLength)
}
