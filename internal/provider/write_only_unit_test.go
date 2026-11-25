package provider

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

func TestResolveWriteOnlySecret(t *testing.T) {
	plain := types.StringValue("plain")
	writeOnly := types.StringValue("write-only")
	null := types.StringNull()

	for _, tt := range []struct {
		name           string
		plain          types.String
		writeOnly      types.String
		versionChanged bool
		wasInUse       bool
		want           writeOnlySecret
	}{
		{"plain", plain, null, false, false, writeOnlySecret{Value: new("plain")}},
		{"plain after write-only", plain, null, false, true, writeOnlySecret{Value: new("plain")}},
		{"neither", null, null, true, false, writeOnlySecret{}},
		{"write-only on create", null, writeOnly, true, false, writeOnlySecret{Value: new("write-only"), InUse: true}},
		{"write-only version changed", null, writeOnly, true, true, writeOnlySecret{Value: new("write-only"), InUse: true}},
		{"write-only unchanged", null, writeOnly, false, true, writeOnlySecret{Keep: true, InUse: true}},
		{"write-only after plain or import", null, writeOnly, false, false, writeOnlySecret{Value: new("write-only"), InUse: true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveWriteOnlySecret(tt.plain, tt.writeOnly, tt.versionChanged, tt.wasInUse)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

const writeOnlyTestProviderConfig = `
provider "sentry" {
	base_url = "%s/api/"
	token    = "test"
}
`

// fakeSymbolSourceAPI mimics the Sentry symbol source endpoints for a single S3
// source: secrets are never returned, `{"hidden-secret": true}` on update keeps
// the stored secret, and an S3 source without a secret is rejected.
type fakeSymbolSourceAPI struct {
	mu sync.Mutex
	// source is the stored source, including secrets.
	source map[string]any
	// lastSecretKey is the `secret_key` value of the last create or update request.
	lastSecretKey any
}

var fakeSymbolSourceSecretFields = []string{"secret_key", "private_key", "password", "appconnectPrivateKey"}

func (f *fakeSymbolSourceAPI) redacted() map[string]any {
	out := maps.Clone(f.source)
	for _, k := range fakeSymbolSourceSecretFields {
		if _, ok := out[k]; ok {
			out[k] = map[string]bool{"hidden-secret": true}
		}
	}
	return out
}

func (f *fakeSymbolSourceAPI) storedSecretKey() any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.source["secret_key"]
}

func (f *fakeSymbolSourceAPI) sentSecretKey() any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.lastSecretKey
}

func newFakeSymbolSourceAPI(t *testing.T) (*fakeSymbolSourceAPI, *httptest.Server) {
	t.Helper()
	f := &fakeSymbolSourceAPI{}

	const sourcesPath = "/api/0/projects/my-org/my-project/symbol-sources/"
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/0/internal/health/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"healthy":{}}`))
	})
	mux.HandleFunc("POST "+sourcesPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.lastSecretKey = body["secret_key"]
		if _, ok := body["secret_key"].(string); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body["id"] = "source-1"
		f.source = body
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(f.redacted())
	})
	mux.HandleFunc("GET "+sourcesPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		if f.source == nil || r.URL.Query().Get("id") != "source-1" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		_ = json.NewEncoder(w).Encode([]any{f.redacted()})
	})
	mux.HandleFunc("PUT "+sourcesPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.lastSecretKey = body["secret_key"]
		if r.URL.Query().Get("id") != "source-1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if body["id"] != "source-1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Backfill hidden secrets like Sentry does.
		for _, k := range fakeSymbolSourceSecretFields {
			if reflect.DeepEqual(body[k], map[string]any{"hidden-secret": true}) {
				body[k] = f.source[k]
			}
		}
		if _, ok := body["secret_key"].(string); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f.source = body
		_ = json.NewEncoder(w).Encode(f.redacted())
	})
	mux.HandleFunc("DELETE "+sourcesPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		if r.URL.Query().Get("id") != "source-1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		f.source = nil
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return f, server
}

func testProjectSymbolSourceWriteOnlyFakeConfig(baseUrl, name, secretAttrs string) string {
	return fmt.Sprintf(writeOnlyTestProviderConfig, baseUrl) + fmt.Sprintf(`
resource "sentry_project_symbol_source" "test" {
	organization = "my-org"
	project      = "my-project"
	type         = "s3"
	name         = "%[1]s"
	layout = {
		type   = "native"
		casing = "default"
	}
	bucket     = "bucket"
	region     = "us-east-1"
	access_key = "access_key"
	%[2]s
}
`, name, secretAttrs)
}

func expectSecretKey(f *fakeSymbolSourceAPI, sent any, stored string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := f.sentSecretKey(); !reflect.DeepEqual(got, sent) {
			return fmt.Errorf("sent secret_key %#v, want %#v", got, sent)
		}
		if got := f.storedSecretKey(); got != stored {
			return fmt.Errorf("stored secret_key %#v, want %q", got, stored)
		}
		return nil
	}
}

func TestProjectSymbolSourceResource_writeOnlyFakeApi(t *testing.T) {
	rn := "sentry_project_symbol_source.test"
	f, server := newFakeSymbolSourceAPI(t)
	hidden := map[string]any{"hidden-secret": true}

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig(server.URL, "s3", `secret_key_wo = "secret-1"
	secret_key_wo_version = 1`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key_wo_version"), knownvalue.Int64Exact(1)),
				},
				Check: expectSecretKey(f, "secret-1", "secret-1"),
			},
			// An update that does not bump the version keeps the stored secret.
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig(server.URL, "s3-renamed", `secret_key_wo = "secret-1"
	secret_key_wo_version = 1`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("s3-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.Null()),
				},
				Check: expectSecretKey(f, hidden, "secret-1"),
			},
			// Bumping the version sends the new value.
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig(server.URL, "s3-renamed", `secret_key_wo = "secret-2"
	secret_key_wo_version = 2`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key_wo_version"), knownvalue.Int64Exact(2)),
				},
				Check: expectSecretKey(f, "secret-2", "secret-2"),
			},
			// Changing only the write-only value is not detected.
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig(server.URL, "s3-renamed", `secret_key_wo = "secret-3"
	secret_key_wo_version = 2`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: expectSecretKey(f, "secret-2", "secret-2"),
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateIdFunc:       resourceid.ImportState3PartIDFunc(rn, "organization", "project", "id"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key_wo", "secret_key_wo_version"},
			},
			// Switching back to the regular attribute stores the secret in state.
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig(server.URL, "s3-renamed", `secret_key = "plain"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.StringExact("plain")),
				},
				Check: expectSecretKey(f, "plain", "plain"),
			},
			// Switching to the write-only attribute sends it even without a version.
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig(server.URL, "s3-renamed", `secret_key_wo = "secret-4"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.Null()),
				},
				Check: expectSecretKey(f, "secret-4", "secret-4"),
			},
		},
	})
}

func TestProjectSymbolSourceResource_writeOnlyValidation(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig("http://127.0.0.1:1", "s3", `secret_key = "plain"
	secret_key_wo = "write-only"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{
				Config: testProjectSymbolSourceWriteOnlyFakeConfig("http://127.0.0.1:1", "s3", `secret_key = "plain"
	secret_key_wo_version = 1`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

// fakeIntegrationAPI mimics a PagerDuty or Opsgenie organization integration.
// Sentry returns the integration keys of all rows in full.
type fakeIntegrationAPI struct {
	mu          sync.Mutex
	providerKey string
	tableKey    string
	nameKey     string
	rows        []map[string]any
	nextId      int
}

func (f *fakeIntegrationAPI) keyOf(name string) any {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, row := range f.rows {
		if row[f.nameKey] == name {
			return row["integration_key"]
		}
	}
	return nil
}

func (f *fakeIntegrationAPI) newId() any {
	f.nextId++
	if f.providerKey == "pagerduty" {
		return f.nextId
	}
	return strconv.Itoa(f.nextId)
}

func newFakeIntegrationAPI(t *testing.T, providerKey, tableKey, nameKey string) (*fakeIntegrationAPI, *httptest.Server) {
	t.Helper()
	f := &fakeIntegrationAPI{providerKey: providerKey, tableKey: tableKey, nameKey: nameKey}
	f.rows = []map[string]any{{nameKey: "other", "integration_key": "other-key", "id": f.newId()}}

	const integrationPath = "/api/0/organizations/my-org/integrations/123/"
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/0/internal/health/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"healthy":{}}`))
	})
	mux.HandleFunc("GET "+integrationPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":                            "123",
			"name":                          providerKey,
			"provider":                      map[string]any{"key": providerKey},
			"configData":                    map[string]any{tableKey: f.rows},
			"organizationId":                1,
			"organizationIntegrationStatus": "active",
			"status":                        "active",
		})
	})
	mux.HandleFunc("POST "+integrationPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		var body map[string][]map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		rows := body[tableKey]
		for _, row := range rows {
			if key, _ := row["integration_key"].(string); key == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if id := fmt.Sprint(row["id"]); id == "0" || id == "" {
				row["id"] = f.newId()
			}
		}
		f.rows = rows
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return f, server
}

var writeOnlyIntegrationTestCases = []struct {
	resourceType string
	providerKey  string
	tableKey     string
	nameKey      string
}{
	{"sentry_integration_pagerduty", "pagerduty", "service_table", "service"},
	{"sentry_integration_opsgenie", "opsgenie", "team_table", "team"},
}

func testIntegrationWriteOnlyFakeConfig(baseUrl, resourceType, nameKey, name, keyAttrs string) string {
	return fmt.Sprintf(writeOnlyTestProviderConfig, baseUrl) + fmt.Sprintf(`
resource "%[1]s" "test" {
	organization   = "my-org"
	integration_id = "123"
	%[2]s = "%[3]s"
	%[4]s
}
`, resourceType, nameKey, name, keyAttrs)
}

func expectIntegrationKey(f *fakeIntegrationAPI, name, key string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := f.keyOf(name); got != key {
			return fmt.Errorf("integration key of %q is %#v, want %q", name, got, key)
		}
		if got := f.keyOf("other"); got != "other-key" {
			return fmt.Errorf("integration key of the other row is %#v, want %q", got, "other-key")
		}
		return nil
	}
}

func TestIntegrationResources_writeOnlyFakeApi(t *testing.T) {
	for _, tc := range writeOnlyIntegrationTestCases {
		t.Run(tc.resourceType, func(t *testing.T) {
			rn := tc.resourceType + ".test"
			f, server := newFakeIntegrationAPI(t, tc.providerKey, tc.tableKey, tc.nameKey)
			config := func(name, keyAttrs string) string {
				return testIntegrationWriteOnlyFakeConfig(server.URL, tc.resourceType, tc.nameKey, name, keyAttrs)
			}

			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_11_0),
				},
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: config("svc", `integration_key_wo = "key-1"
	integration_key_wo_version = 1`),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.Null()),
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key_wo"), knownvalue.Null()),
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key_wo_version"), knownvalue.Int64Exact(1)),
						},
						Check: expectIntegrationKey(f, "svc", "key-1"),
					},
					// An update that does not bump the version keeps the key.
					{
						Config: config("svc-renamed", `integration_key_wo = "key-1"
	integration_key_wo_version = 1`),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.Null()),
						},
						Check: expectIntegrationKey(f, "svc-renamed", "key-1"),
					},
					// Bumping the version sends the new key.
					{
						Config: config("svc-renamed", `integration_key_wo = "key-2"
	integration_key_wo_version = 2`),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
							},
						},
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.Null()),
						},
						Check: expectIntegrationKey(f, "svc-renamed", "key-2"),
					},
					// Changing only the write-only value is not detected.
					{
						Config: config("svc-renamed", `integration_key_wo = "key-3"
	integration_key_wo_version = 2`),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectEmptyPlan(),
							},
						},
						Check: expectIntegrationKey(f, "svc-renamed", "key-2"),
					},
					// Import reads the key from Sentry, as before.
					{
						ResourceName:            rn,
						ImportState:             true,
						ImportStateIdFunc:       resourceid.ImportState3PartIDFunc(rn, "organization", "integration_id", "id"),
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"integration_key", "integration_key_wo", "integration_key_wo_version"},
					},
					// Switching back to the regular attribute stores the key in state.
					{
						Config: config("svc-renamed", `integration_key = "plain"`),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.StringExact("plain")),
						},
						Check: expectIntegrationKey(f, "svc-renamed", "plain"),
					},
					// Switching to the write-only attribute sends it even without a version.
					{
						Config: config("svc-renamed", `integration_key_wo = "key-4"`),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.Null()),
						},
						Check: expectIntegrationKey(f, "svc-renamed", "key-4"),
					},
				},
			})
		})
	}
}

// The regular integration_key attribute keeps working as before, including
// being read back from Sentry on refresh and import.
func TestIntegrationResources_integrationKeyFakeApi(t *testing.T) {
	for _, tc := range writeOnlyIntegrationTestCases {
		t.Run(tc.resourceType, func(t *testing.T) {
			rn := tc.resourceType + ".test"
			f, server := newFakeIntegrationAPI(t, tc.providerKey, tc.tableKey, tc.nameKey)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: testIntegrationWriteOnlyFakeConfig(server.URL, tc.resourceType, tc.nameKey, "svc", `integration_key = "key-1"`),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.StringExact("key-1")),
						},
						Check: expectIntegrationKey(f, "svc", "key-1"),
					},
					{
						ResourceName:      rn,
						ImportState:       true,
						ImportStateIdFunc: resourceid.ImportState3PartIDFunc(rn, "organization", "integration_id", "id"),
						ImportStateVerify: true,
					},
				},
			})
		})
	}
}

func TestIntegrationResources_writeOnlyValidation(t *testing.T) {
	for _, tc := range writeOnlyIntegrationTestCases {
		t.Run(tc.resourceType, func(t *testing.T) {
			config := func(keyAttrs string) string {
				return testIntegrationWriteOnlyFakeConfig("http://127.0.0.1:1", tc.resourceType, tc.nameKey, "svc", keyAttrs)
			}

			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_11_0),
				},
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      config(``),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured`),
					},
					{
						Config: config(`integration_key = "plain"
	integration_key_wo = "write-only"`),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured`),
					},
					{
						Config: config(`integration_key = "plain"
	integration_key_wo_version = 1`),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
					},
				},
			})
		})
	}
}

// The update body must put each secret under the key Sentry uses for it, as a
// string or as the hidden secret marker, and leave unset secrets out.
func TestProjectSymbolSourceUpdateBody_secrets(t *testing.T) {
	jsonKeys := map[string]string{
		"app_connect_private_key": "appconnectPrivateKey",
		"password":                "password",
		"secret_key":              "secret_key",
		"private_key":             "private_key",
	}

	for _, secret := range projectSymbolSourceSecrets {
		for _, value := range []any{"value", projectSymbolSourceHiddenSecret} {
			t.Run(fmt.Sprintf("%s/%v", secret.attr, value), func(t *testing.T) {
				body := &projectSymbolSourceUpdateBody{
					UpdateProjectSymbolSourceParams: &sentry.UpdateProjectSymbolSourceParams{ID: new("source-1")},
				}
				secret.set(body, value)

				b, err := json.Marshal(body)
				if err != nil {
					t.Fatal(err)
				}
				var got map[string]any
				if err := json.Unmarshal(b, &got); err != nil {
					t.Fatal(err)
				}

				want := value
				if _, ok := value.(string); !ok {
					want = map[string]any{"hidden-secret": true}
				}
				if !reflect.DeepEqual(got[jsonKeys[secret.attr]], want) {
					t.Fatalf("got %s, want %q set to %v", b, jsonKeys[secret.attr], want)
				}
				for attr, key := range jsonKeys {
					if _, ok := got[key]; ok && attr != secret.attr {
						t.Fatalf("got %s, want %q left out", b, key)
					}
				}
			})
		}
	}
}
