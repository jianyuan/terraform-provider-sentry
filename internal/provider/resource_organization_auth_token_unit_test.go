package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// fakeOrgAuthTokenAPI mimics the Sentry org auth token endpoints: the token value is only
// returned by create, the list is not paginated and never includes the value, and revoking a
// token that is already revoked returns 404.
type fakeOrgAuthTokenAPI struct {
	mu     sync.Mutex
	nextId int
	tokens map[string]map[string]any
	// creates counts successful creates.
	creates int
	// lastScopes holds the `scopes` of the last create request.
	lastScopes any
	// forbidden makes every org auth token endpoint return 403, as Sentry does for personal tokens.
	forbidden bool
}

func (f *fakeOrgAuthTokenAPI) revoke(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.tokens, id)
}

func (f *fakeOrgAuthTokenAPI) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.tokens)
}

func newFakeOrgAuthTokenAPI(t *testing.T) (*fakeOrgAuthTokenAPI, *httptest.Server) {
	t.Helper()
	f := &fakeOrgAuthTokenAPI{nextId: 100, tokens: map[string]map[string]any{}}

	const tokensPath = "/api/0/organizations/my-org/org-auth-tokens/"
	forbidden := func(w http.ResponseWriter) bool {
		if !f.forbidden {
			return false
		}
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"detail":"You do not have permission to perform this action."}`))
		return true
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/0/internal/health/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"healthy":{}}`))
	})
	mux.HandleFunc("POST "+tokensPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		if forbidden(w) {
			return
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.lastScopes = body["scopes"]

		id := strconv.Itoa(f.nextId)
		f.nextId++
		f.creates++
		value := "sntrys_secret" + id
		token := map[string]any{
			"id":                  id,
			"name":                body["name"],
			"scopes":              body["scopes"],
			"tokenLastCharacters": value[len(value)-4:],
			"dateCreated":         "2026-11-26T12:00:00.123456Z",
			"dateLastUsed":        nil,
			"projectLastUsedId":   nil,
		}
		f.tokens[id] = token

		resp := map[string]any{"token": value}
		for k, v := range token {
			resp[k] = v
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("GET "+tokensPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		if forbidden(w) {
			return
		}
		list := []map[string]any{}
		for _, token := range f.tokens {
			// Like Sentry, the last use changes on every request the token makes.
			token["dateLastUsed"] = "2026-11-27T08:00:00Z"
			list = append(list, token)
		}
		_ = json.NewEncoder(w).Encode(list)
	})
	mux.HandleFunc("DELETE "+tokensPath+"{id}/", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		if forbidden(w) {
			return
		}
		if _, ok := f.tokens[r.PathValue("id")]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		delete(f.tokens, r.PathValue("id"))
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return f, server
}

func testOrganizationAuthTokenUnitConfig(baseUrl string, attrs string) string {
	return fmt.Sprintf(`
provider "sentry" {
	base_url = "%s/api/"
	token    = "test"
}

resource "sentry_organization_auth_token" "test" {
	organization = "my-org"
	%s
}
`, baseUrl, attrs)
}

// captureAttr stores the value of a state attribute so a later step can compare against it.
func captureAttr(rn, attr string, dst *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource %s not found in state", rn)
		}
		*dst = rs.Primary.Attributes[attr]
		return nil
	}
}

func TestOrganizationAuthTokenResource_lifecycle(t *testing.T) {
	t.Parallel()

	rn := "sentry_organization_auth_token.test"
	f, server := newFakeOrgAuthTokenAPI(t)
	var firstId string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_6),
		},
		Steps: []resource.TestStep{
			{
				Config: testOrganizationAuthTokenUnitConfig(server.URL, `name = "ci"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact("100")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("ci")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("org:ci")})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token"), knownvalue.StringExact("sntrys_secret100")),
					statecheck.ExpectSensitiveValue(rn, tfjsonpath.New("token")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token_last_characters"), knownvalue.StringExact("t100")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("date_created"), knownvalue.StringExact("2026-11-26T12:00:00Z")),
				},
				Check: captureAttr(rn, "id", &firstId),
			},
			{
				// A refresh reads the list, which has no token value and a new last use: the token
				// stays in state and nothing changes.
				Config: testOrganizationAuthTokenUnitConfig(server.URL, `name = "ci"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token"), knownvalue.StringExact("sntrys_secret100")),
				},
			},
			{
				// Setting the default scope explicitly is not a change.
				Config: testOrganizationAuthTokenUnitConfig(server.URL, `
	name   = "ci"
	scopes = ["org:ci"]
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				// Renaming replaces the token.
				Config: testOrganizationAuthTokenUnitConfig(server.URL, `name = "ci-renamed"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionDestroyBeforeCreate),
						plancheck.ExpectUnknownValue(rn, tfjsonpath.New("token")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact("101")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("ci-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token"), knownvalue.StringExact("sntrys_secret101")),
				},
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateId:           "my-org/101",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported resource, got %d", len(states))
					}
					if v, ok := states[0].Attributes["token"]; ok && v != "" {
						return fmt.Errorf("expected token to be null after import, got %q", v)
					}
					return nil
				},
			},
		},
	})

	if firstId != "100" {
		t.Fatalf("expected the first token to have ID 100, got %q", firstId)
	}
	if f.creates != 2 {
		t.Fatalf("expected 2 creates, got %d", f.creates)
	}
	if !reflect.DeepEqual(f.lastScopes, []any{"org:ci"}) {
		t.Fatalf("expected scopes to be sent as [org:ci], got %#v", f.lastScopes)
	}
	if n := f.count(); n != 0 {
		t.Fatalf("expected every token to be revoked, %d left", n)
	}
}

func TestOrganizationAuthTokenResource_revokedOutsideTerraform(t *testing.T) {
	t.Parallel()

	rn := "sentry_organization_auth_token.test"
	f, server := newFakeOrgAuthTokenAPI(t)
	config := testOrganizationAuthTokenUnitConfig(server.URL, `name = "ci"`)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact("100")),
				},
			},
			{
				// The revoked token is removed from state on refresh, and a new one is created.
				PreConfig: func() { f.revoke("100") },
				Config:    config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact("101")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token"), knownvalue.StringExact("sntrys_secret101")),
				},
			},
		},
	})
}

func TestOrganizationAuthTokenResource_forbidden(t *testing.T) {
	t.Parallel()

	f, server := newFakeOrgAuthTokenAPI(t)
	f.forbidden = true

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testOrganizationAuthTokenUnitConfig(server.URL, `name = "ci"`),
				ExpectError: regexp.MustCompile(`(?s)Not\s+allowed\s+to\s+manage\s+org\s+auth\s+tokens.*got\s+status\s+403.*OAuth\s+application\s+token\s+with\s+the\s+` + "`org:tokens`" + `\s+scope.*Personal\s+auth\s+tokens,\s+internal\s+integration\s+tokens`),
			},
		},
	})

	if f.creates != 0 {
		t.Fatalf("expected no creates, got %d", f.creates)
	}
}

func TestOrganizationAuthTokenResource_validation(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		attrs     string
		wantError string
	}{
		{"scope not allowed", "name = \"ci\"\n\tscopes = [\"org:ci\", \"project:releases\"]", `(?s)value\s+must\s+be\s+one\s+of:\s+\["org:ci"\]`},
		{"no scopes", "name = \"ci\"\n\tscopes = []", `(?s)set\s+must\s+contain\s+at\s+least\s+1\s+elements`},
		{"empty name", `name = ""`, `(?s)character\s+count\s+must\s+be\s+between\s+1\s+and\s+255`},
		{"name too long", fmt.Sprintf("name = %q", strings.Repeat("a", 256)), `(?s)character\s+count\s+must\s+be\s+between\s+1\s+and\s+255`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, server := newFakeOrgAuthTokenAPI(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testOrganizationAuthTokenUnitConfig(server.URL, tc.attrs),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(tc.wantError),
					},
				},
			})
		})
	}

	t.Run("255 characters", func(t *testing.T) {
		t.Parallel()
		_, server := newFakeOrgAuthTokenAPI(t)
		name := strings.Repeat("é", 255)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testOrganizationAuthTokenUnitConfig(server.URL, fmt.Sprintf("name = %q", name)),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("sentry_organization_auth_token.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					},
				},
			},
		})
	})
}
