package sentry

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectKeyService_List(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/keys/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"name": "Fabulous Key",
				"projectId": 2,
				"secret": "e3e7c32c93f44a19b144e4e85940d3a6",
				"label": "Fabulous Key",
				"dsn": {
					"secret": "https://befdbf32724c4ae0a3d286717b1f8127:e3e7c32c93f44a19b144e4e85940d3a6@sentry.io/2",
					"csp": "https://sentry.io/api/2/csp-report/?sentry_key=befdbf32724c4ae0a3d286717b1f8127",
					"public": "https://befdbf32724c4ae0a3d286717b1f8127@sentry.io/2"
				},
				"public": "befdbf32724c4ae0a3d286717b1f8127",
				"rateLimit": null,
				"dateCreated": "2017-07-18T19:29:44.791Z",
				"id": "befdbf32724c4ae0a3d286717b1f8127",
				"isActive": true
			}
		]`)
	})

	client := NewClient(httpClient, nil, "")
	projectKeys, _, err := client.ProjectKeys.List("the-interstellar-jurisdiction", "pump-station")
	assert.NoError(t, err)

	expected := []ProjectKey{
		{
			ID:        "befdbf32724c4ae0a3d286717b1f8127",
			Name:      "Fabulous Key",
			Public:    "befdbf32724c4ae0a3d286717b1f8127",
			Secret:    "e3e7c32c93f44a19b144e4e85940d3a6",
			ProjectID: 2,
			IsActive:  true,
			DSN: ProjectKeyDSN{
				Secret: "https://befdbf32724c4ae0a3d286717b1f8127:e3e7c32c93f44a19b144e4e85940d3a6@sentry.io/2",
				Public: "https://befdbf32724c4ae0a3d286717b1f8127@sentry.io/2",
				CSP:    "https://sentry.io/api/2/csp-report/?sentry_key=befdbf32724c4ae0a3d286717b1f8127",
			},
			DateCreated: mustParseTime("2017-07-18T19:29:44.791Z"),
		},
	}
	assert.Equal(t, expected, projectKeys)
}

func TestProjectKeyService_Create(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/keys/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "Fabulous Key",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"name": "Fabulous Key",
			"projectId": 2,
			"secret": "e3e7c32c93f44a19b144e4e85940d3a6",
			"label": "Fabulous Key",
			"dsn": {
				"secret": "https://befdbf32724c4ae0a3d286717b1f8127:e3e7c32c93f44a19b144e4e85940d3a6@sentry.io/2",
				"csp": "https://sentry.io/api/2/csp-report/?sentry_key=befdbf32724c4ae0a3d286717b1f8127",
				"public": "https://befdbf32724c4ae0a3d286717b1f8127@sentry.io/2"
			},
			"public": "befdbf32724c4ae0a3d286717b1f8127",
			"rateLimit": null,
			"dateCreated": "2017-07-18T19:29:44.791Z",
			"id": "befdbf32724c4ae0a3d286717b1f8127",
			"isActive": true
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &CreateProjectKeyParams{
		Name: "Fabulous Key",
	}
	projectKey, _, err := client.ProjectKeys.Create("the-interstellar-jurisdiction", "pump-station", params)
	assert.NoError(t, err)
	expected := &ProjectKey{
		ID:        "befdbf32724c4ae0a3d286717b1f8127",
		Name:      "Fabulous Key",
		Public:    "befdbf32724c4ae0a3d286717b1f8127",
		Secret:    "e3e7c32c93f44a19b144e4e85940d3a6",
		ProjectID: 2,
		IsActive:  true,
		DSN: ProjectKeyDSN{
			Secret: "https://befdbf32724c4ae0a3d286717b1f8127:e3e7c32c93f44a19b144e4e85940d3a6@sentry.io/2",
			Public: "https://befdbf32724c4ae0a3d286717b1f8127@sentry.io/2",
			CSP:    "https://sentry.io/api/2/csp-report/?sentry_key=befdbf32724c4ae0a3d286717b1f8127",
		},
		DateCreated: mustParseTime("2017-07-18T19:29:44.791Z"),
	}
	assert.Equal(t, expected, projectKey)
}

func TestProjectKeyService_Update(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/keys/befdbf32724c4ae0a3d286717b1f8127/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSON(t, map[string]interface{}{
			"name": "Fabulous Key",
		}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"name": "Fabulous Key",
			"projectId": 2,
			"secret": "e3e7c32c93f44a19b144e4e85940d3a6",
			"label": "Fabulous Key",
			"dsn": {
				"secret": "https://befdbf32724c4ae0a3d286717b1f8127:e3e7c32c93f44a19b144e4e85940d3a6@sentry.io/2",
				"csp": "https://sentry.io/api/2/csp-report/?sentry_key=befdbf32724c4ae0a3d286717b1f8127",
				"public": "https://befdbf32724c4ae0a3d286717b1f8127@sentry.io/2"
			},
			"public": "befdbf32724c4ae0a3d286717b1f8127",
			"rateLimit": null,
			"dateCreated": "2017-07-18T19:29:44.791Z",
			"id": "befdbf32724c4ae0a3d286717b1f8127",
			"isActive": true
		}`)
	})

	client := NewClient(httpClient, nil, "")
	params := &UpdateProjectKeyParams{
		Name: "Fabulous Key",
	}
	projectKey, _, err := client.ProjectKeys.Update("the-interstellar-jurisdiction", "pump-station", "befdbf32724c4ae0a3d286717b1f8127", params)
	assert.NoError(t, err)
	expected := &ProjectKey{
		ID:        "befdbf32724c4ae0a3d286717b1f8127",
		Name:      "Fabulous Key",
		Public:    "befdbf32724c4ae0a3d286717b1f8127",
		Secret:    "e3e7c32c93f44a19b144e4e85940d3a6",
		ProjectID: 2,
		IsActive:  true,
		DSN: ProjectKeyDSN{
			Secret: "https://befdbf32724c4ae0a3d286717b1f8127:e3e7c32c93f44a19b144e4e85940d3a6@sentry.io/2",
			Public: "https://befdbf32724c4ae0a3d286717b1f8127@sentry.io/2",
			CSP:    "https://sentry.io/api/2/csp-report/?sentry_key=befdbf32724c4ae0a3d286717b1f8127",
		},
		DateCreated: mustParseTime("2017-07-18T19:29:44.791Z"),
	}
	assert.Equal(t, expected, projectKey)
}

func TestProjectKeyService_Delete(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/0/projects/the-interstellar-jurisdiction/pump-station/keys/befdbf32724c4ae0a3d286717b1f8127/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	client := NewClient(httpClient, nil, "")
	_, err := client.ProjectKeys.Delete("the-interstellar-jurisdiction", "pump-station", "befdbf32724c4ae0a3d286717b1f8127")
	assert.NoError(t, err)

}
