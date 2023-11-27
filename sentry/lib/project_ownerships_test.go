package sentry

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectOwnershipsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/powerful-abolitionist/ownership/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"raw": "# assign issues to the product team, no matter the area\nurl:https://example.com/areas/*/*/products/* #product-team",
			"fallthrough": false,
			"dateCreated": "2021-11-18T13:09:16.819818Z",
			"lastUpdated": "2022-03-01T14:00:31.317734Z",
			"isActive": true,
			"autoAssignment": true,
			"codeownersAutoSync": null
		}`)
	})

	ctx := context.Background()
	ownership, _, err := client.ProjectOwnerships.Get(ctx, "the-interstellar-jurisdiction", "powerful-abolitionist")
	assert.NoError(t, err)

	expected := &ProjectOwnership{
		Raw:                "# assign issues to the product team, no matter the area\nurl:https://example.com/areas/*/*/products/* #product-team",
		FallThrough:        false,
		IsActive:           true,
		AutoAssignment:     true,
		CodeownersAutoSync: nil,
		DateCreated:        mustParseTime("2021-11-18T13:09:16.819818Z"),
		LastUpdated:        mustParseTime("2022-03-01T14:00:31.317734Z"),
	}

	assert.Equal(t, expected, ownership)
}

func TestProjectOwnershipsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/the-obese-philosophers/ownership/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSON(t, map[string]interface{}{
			"raw": "# assign issues to the product team, no matter the area\nurl:https://example.com/areas/*/*/products/* #product-team",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"raw": "# assign issues to the product team, no matter the area\nurl:https://example.com/areas/*/*/products/* #product-team",
			"fallthrough": false,
			"dateCreated": "2021-11-18T13:09:16.819818Z",
			"lastUpdated": "2022-03-01T14:00:31.317734Z",
			"isActive": true,
			"autoAssignment": true,
			"codeownersAutoSync": null
		}`)
	})

	params := &UpdateProjectOwnershipParams{
		Raw: "# assign issues to the product team, no matter the area\nurl:https://example.com/areas/*/*/products/* #product-team",
	}
	ctx := context.Background()
	ownership, _, err := client.ProjectOwnerships.Update(ctx, "the-interstellar-jurisdiction", "the-obese-philosophers", params)
	assert.NoError(t, err)
	expected := &ProjectOwnership{
		Raw:                "# assign issues to the product team, no matter the area\nurl:https://example.com/areas/*/*/products/* #product-team",
		FallThrough:        false,
		IsActive:           true,
		AutoAssignment:     true,
		CodeownersAutoSync: nil,
		DateCreated:        mustParseTime("2021-11-18T13:09:16.819818Z"),
		LastUpdated:        mustParseTime("2022-03-01T14:00:31.317734Z"),
	}
	assert.Equal(t, expected, ownership)
}
