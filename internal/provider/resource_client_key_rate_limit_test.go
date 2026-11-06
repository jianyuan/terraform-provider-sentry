package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// newClientKeyRateLimitTestServer returns a fake Sentry API that mimics a project whose plan
// does not apply client key rate limits: the key is created, but `rateLimit` is always null.
func newClientKeyRateLimitTestServer(t *testing.T, projectFeatures []string, keysCreated *atomic.Int32) *httptest.Server {
	t.Helper()

	key := map[string]any{
		"id":                      "abc123",
		"name":                    "test",
		"label":                   "test",
		"projectId":               1,
		"public":                  "public",
		"secret":                  "secret",
		"isActive":                true,
		"rateLimit":               nil,
		"dsn":                     map[string]string{"public": "https://public@example.com/1", "secret": "https://public:secret@example.com/1", "csp": "https://example.com/csp"},
		"browserSdkVersion":       "7.x",
		"dynamicSdkLoaderOptions": map[string]bool{"hasDebug": false, "hasPerformance": false, "hasReplay": false},
		"dateCreated":             "2026-01-01T00:00:00Z",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/0/internal/health/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"healthy":{}}`))
	})
	mux.HandleFunc("GET /api/0/projects/test-org/test-project/", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":          "1",
			"slug":        "test-project",
			"name":        "test-project",
			"features":    projectFeatures,
			"dateCreated": "2026-01-01T00:00:00Z",
		})
	})
	mux.HandleFunc("POST /api/0/projects/test-org/test-project/keys/", func(w http.ResponseWriter, r *http.Request) {
		keysCreated.Add(1)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(key)
	})
	mux.HandleFunc("GET /api/0/projects/test-org/test-project/keys/abc123/", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(key)
	})
	mux.HandleFunc("DELETE /api/0/projects/test-org/test-project/keys/abc123/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return server
}

func testClientKeyRateLimitConfig(baseUrl string) string {
	return `
provider "sentry" {
	base_url = "` + baseUrl + `/api/"
	token    = "test"
}

resource "sentry_key" "test" {
	organization = "test-org"
	project      = "test-project"
	name         = "test"

	rate_limit_window = 60
	rate_limit_count  = 1000
}
`
}

// Regression test for https://github.com/jianyuan/terraform-provider-sentry/issues/843
func TestClientKeyResource_rateLimitFeatureNotAvailable(t *testing.T) {
	var keysCreated atomic.Int32
	server := newClientKeyRateLimitTestServer(t, []string{"servicehooks"}, &keysCreated)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testClientKeyRateLimitConfig(server.URL),
				ExpectError: regexp.MustCompile(`Rate limits are not available`),
			},
		},
	})

	if n := keysCreated.Load(); n != 0 {
		t.Fatalf("expected no client key to be created, got %d", n)
	}
}

func TestClientKeyResource_rateLimitNotApplied(t *testing.T) {
	var keysCreated atomic.Int32
	server := newClientKeyRateLimitTestServer(t, []string{"rate-limits"}, &keysCreated)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testClientKeyRateLimitConfig(server.URL),
				ExpectError: regexp.MustCompile(`Rate limit was not applied`),
			},
		},
	})
}
