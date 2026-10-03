package sentryclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func newTestApiClient(t *testing.T, handler http.HandlerFunc) *apiclient.ClientWithResponses {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	apiClient, err := apiclient.NewClientWithResponses(server.URL + "/api/")
	if err != nil {
		t.Fatalf("failed to create API client: %s", err)
	}
	return apiClient
}

func TestGetProjectSlug(t *testing.T) {
	apiClient := newTestApiClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/0/projects/my-org/1234/" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": "1234", "slug": "my-project"}`))
	})

	slug, err := GetProjectSlug(context.Background(), apiClient, "my-org", "1234")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if slug != "my-project" {
		t.Errorf("expected slug %q, got %q", "my-project", slug)
	}
}

func TestGetProjectSlug_notFound(t *testing.T) {
	apiClient := newTestApiClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail": "The requested resource does not exist"}`))
	})

	if _, err := GetProjectSlug(context.Background(), apiClient, "my-org", "1234"); err == nil {
		t.Fatal("expected an error")
	}
}
