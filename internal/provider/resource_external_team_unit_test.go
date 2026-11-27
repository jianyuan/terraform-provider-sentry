package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// fakeExternalActorStore mimics how Sentry stores external teams: it keys mappings by ID, returns
// 200 with the existing mapping when a create matches one, and omits `externalId` when it is empty.
type fakeExternalActorStore struct {
	mu      sync.Mutex
	nextId  int
	actors  map[string]map[string]any
	creates int
}

func newFakeExternalActorStore() *fakeExternalActorStore {
	return &fakeExternalActorStore{nextId: 42, actors: map[string]map[string]any{}}
}

// create stores a mapping from a decoded create request body. Like Sentry, it matches an existing
// mapping regardless of the team it belongs to.
func (s *fakeExternalActorStore) create(w http.ResponseWriter, body map[string]any, teamId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range s.actors {
		if a["provider"] == body["provider"] && a["externalName"] == body["external_name"] {
			_ = json.NewEncoder(w).Encode(a)
			return
		}
	}

	id := strconv.Itoa(s.nextId)
	s.nextId++
	s.creates++
	a := map[string]any{
		"id":            id,
		"provider":      body["provider"],
		"externalName":  body["external_name"],
		"integrationId": fmt.Sprint(body["integration_id"]),
		"teamId":        teamId,
	}
	if v, ok := body["external_id"].(string); ok && v != "" {
		a["externalId"] = v
	}
	s.actors[id] = a

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(a)
}

// update applies a partial update. Like Sentry, the response has no `teamId`.
func (s *fakeExternalActorStore) update(w http.ResponseWriter, r *http.Request, id string) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.actors[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if v, ok := body["provider"]; ok {
		a["provider"] = v
	}
	if v, ok := body["external_name"]; ok {
		a["externalName"] = v
	}
	if v, ok := body["integration_id"]; ok {
		a["integrationId"] = fmt.Sprint(v)
	}
	if v, ok := body["external_id"]; ok {
		if s, ok := v.(string); ok && s != "" {
			a["externalId"] = s
		} else {
			delete(a, "externalId")
		}
	}

	resp := map[string]any{}
	for k, v := range a {
		if k != "teamId" {
			resp[k] = v
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *fakeExternalActorStore) delete(w http.ResponseWriter, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.actors[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	delete(s.actors, id)
	w.WriteHeader(http.StatusNoContent)
}

// list returns the mappings that belong to teamId.
func (s *fakeExternalActorStore) list(teamId string) []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := []map[string]any{}
	for _, a := range s.actors {
		if a["teamId"] == teamId {
			out = append(out, a)
		}
	}
	return out
}

func newFakeSentryServer(t *testing.T, register func(mux *http.ServeMux)) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/0/internal/health/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"healthy":{}}`))
	})
	register(mux)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return server
}

func newExternalTeamTestServer(t *testing.T, store *fakeExternalActorStore) *httptest.Server {
	return newFakeSentryServer(t, func(mux *http.ServeMux) {
		mux.HandleFunc("GET /api/0/teams/test-org/test-team/", func(w http.ResponseWriter, r *http.Request) {
			team := map[string]any{"id": "1", "slug": "test-team", "name": "test-team"}
			if r.URL.Query().Get("expand") == "externalTeams" {
				team["externalTeams"] = store.list("1")
			}
			_ = json.NewEncoder(w).Encode(team)
		})
		mux.HandleFunc("POST /api/0/teams/test-org/test-team/external-teams/", func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			store.create(w, body, "1")
		})
		mux.HandleFunc("PUT /api/0/teams/test-org/test-team/external-teams/{id}/", func(w http.ResponseWriter, r *http.Request) {
			store.update(w, r, r.PathValue("id"))
		})
		mux.HandleFunc("DELETE /api/0/teams/test-org/test-team/external-teams/{id}/", func(w http.ResponseWriter, r *http.Request) {
			store.delete(w, r.PathValue("id"))
		})
	})
}

func testExternalTeamUnitConfig(baseUrl string, externalName string, externalId string) string {
	externalIdLine := ""
	if externalId != "" {
		externalIdLine = fmt.Sprintf("external_id    = %q", externalId)
	}
	return fmt.Sprintf(`
provider "sentry" {
	base_url = "%[1]s/api/"
	token    = "test"
}

resource "sentry_external_team" "test" {
	organization   = "test-org"
	team           = "test-team"
	provider_name  = "slack"
	external_name  = %[2]q
	integration_id = "123"
	%[3]s
}
`, baseUrl, externalName, externalIdLine)
}

func TestExternalTeamResource_lifecycle(t *testing.T) {
	t.Parallel()

	rn := "sentry_external_team.test"
	store := newFakeExternalActorStore()
	server := newExternalTeamTestServer(t, store)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testExternalTeamUnitConfig(server.URL, "#backend-alerts", "C0123456789"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact("42")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("provider_name"), knownvalue.StringExact("slack")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_name"), knownvalue.StringExact("#backend-alerts")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.StringExact("123")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_id"), knownvalue.StringExact("C0123456789")),
				},
			},
			{
				// Rename in place and clear external_id, which must be sent as an explicit null.
				Config: testExternalTeamUnitConfig(server.URL, "#backend-oncall", ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact("42")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_name"), knownvalue.StringExact("#backend-oncall")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_id"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateId:     "test-org/test-team/42",
				ImportStateVerify: true,
			},
		},
	})

	if store.creates != 1 {
		t.Fatalf("expected exactly one create, got %d", store.creates)
	}
	if n := len(store.list("1")); n != 0 {
		t.Fatalf("expected the mapping to be deleted, %d left", n)
	}
}

func TestExternalTeamResource_alreadyExists(t *testing.T) {
	t.Parallel()

	store := newFakeExternalActorStore()
	// A mapping for the same Slack channel that belongs to another team.
	store.actors["7"] = map[string]any{
		"id":            "7",
		"provider":      "slack",
		"externalName":  "#backend-alerts",
		"integrationId": "123",
		"teamId":        "99",
	}
	server := newExternalTeamTestServer(t, store)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testExternalTeamUnitConfig(server.URL, "#backend-alerts", ""),
				ExpectError: regexp.MustCompile(`(?s)External team already exists.*mapping\s+ID\s+"7",\s+Sentry\s+team\s+ID\s+"99"`),
			},
		},
	})

	if _, ok := store.actors["7"]; !ok {
		t.Fatal("the existing mapping must not be deleted")
	}
}
