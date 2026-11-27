package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

func TestProvider_ephemeralResources(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	server, err := providerserver.NewProtocol6WithError(New("test")())()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("unexpected error: %s: %s", d.Summary, d.Detail)
		}
	}

	schema, ok := resp.EphemeralResourceSchemas["sentry_key"]
	if !ok {
		t.Fatal("expected the sentry_key ephemeral resource to be registered")
	}
	for _, attr := range schema.Block.Attributes {
		if attr.Name == "dsn" && !attr.Sensitive {
			t.Fatal("expected dsn to be sensitive")
		}
		if attr.Name == "secret" && !attr.Sensitive {
			t.Fatal("expected secret to be sensitive")
		}
	}
}

func testClientKeyJSON(id, name, dateCreated string) map[string]any {
	return map[string]any{
		"id":                      id,
		"name":                    name,
		"label":                   name,
		"projectId":               1,
		"public":                  "public-" + id,
		"secret":                  "secret-" + id,
		"isActive":                true,
		"rateLimit":               nil,
		"dsn":                     map[string]string{"public": "https://public-" + id + "@example.com/1", "secret": "https://public-" + id + ":secret-" + id + "@example.com/1", "csp": "https://example.com/csp"},
		"browserSdkVersion":       "7.x",
		"dynamicSdkLoaderOptions": map[string]bool{"hasDebug": false, "hasPerformance": true, "hasReplay": true},
		"dateCreated":             dateCreated,
	}
}

// newClientKeyEphemeralTestServer returns a fake Sentry API for a project
// with two client keys, "Default" (the oldest) and "other".
func newClientKeyEphemeralTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	keys := []map[string]any{
		testClientKeyJSON("bbb", "other", "2026-02-01T00:00:00Z"),
		testClientKeyJSON("aaa", "Default", "2026-01-01T00:00:00Z"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/0/internal/health/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"healthy":{}}`))
	})
	mux.HandleFunc("GET /api/0/projects/test-org/test-project/keys/", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(keys)
	})
	mux.HandleFunc("GET /api/0/projects/test-org/test-project/keys/{id}/", func(w http.ResponseWriter, r *http.Request) {
		for _, key := range keys {
			if key["id"] == r.PathValue("id") {
				_ = json.NewEncoder(w).Encode(key)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"The requested resource does not exist"}`))
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return server
}

func openClientKeyEphemeralResourceAgainst(t *testing.T, baseURL string, attrs map[string]tftypes.Value) *ephemeral.OpenResponse {
	t.Helper()
	ctx := context.Background()

	apiClient, err := apiclient.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatal(err)
	}

	e := NewClientKeyEphemeralResource()
	configureResp := &ephemeral.ConfigureResponse{}
	e.(ephemeral.EphemeralResourceWithConfigure).Configure(ctx, ephemeral.ConfigureRequest{
		ProviderData: &providerdata.ProviderData{ApiClient: apiClient},
	}, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configure: %v", configureResp.Diagnostics)
	}

	schemaResp := &ephemeral.SchemaResponse{}
	e.Schema(ctx, ephemeral.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}

	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		if v, ok := attrs[name]; ok {
			values[name] = v
		} else {
			values[name] = tftypes.NewValue(typ, nil)
		}
	}

	config := tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, values)}
	resp := &ephemeral.OpenResponse{Result: tfsdk.EphemeralResultData{Schema: schemaResp.Schema}}
	e.Open(ctx, ephemeral.OpenRequest{Config: config}, resp)
	return resp
}

func TestClientKeyEphemeralResourceOpen(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	server := newClientKeyEphemeralTestServer(t)

	project := map[string]tftypes.Value{
		"organization": tftypes.NewValue(tftypes.String, "test-org"),
		"project":      tftypes.NewValue(tftypes.String, "test-project"),
	}

	for _, tc := range []struct {
		name      string
		attrs     map[string]tftypes.Value
		wantId    string
		wantError string
	}{
		{name: "id", attrs: map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, "bbb")}, wantId: "bbb"},
		{name: "name", attrs: map[string]tftypes.Value{"name": tftypes.NewValue(tftypes.String, "other")}, wantId: "bbb"},
		{name: "first", attrs: map[string]tftypes.Value{"first": tftypes.NewValue(tftypes.Bool, true)}, wantId: "aaa"},
		{name: "multiple keys", attrs: map[string]tftypes.Value{}, wantError: "Multiple keys found"},
		{name: "name not found", attrs: map[string]tftypes.Value{"name": tftypes.NewValue(tftypes.String, "missing")}, wantError: "No matching client key found"},
		{name: "id not found", attrs: map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, "missing")}, wantError: "got status 404"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			attrs := map[string]tftypes.Value{}
			for k, v := range project {
				attrs[k] = v
			}
			for k, v := range tc.attrs {
				attrs[k] = v
			}

			resp := openClientKeyEphemeralResourceAgainst(t, server.URL+"/api/", attrs)

			if tc.wantError != "" {
				if !resp.Diagnostics.HasError() {
					t.Fatal("expected an error")
				}
				if got := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(got, tc.wantError) {
					t.Fatalf("expected the error to contain %q, got: %s", tc.wantError, got)
				}
				return
			}

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}

			var result ClientKeyLookupModel
			if diags := resp.Result.Get(ctx, &result); diags.HasError() {
				t.Fatalf("get result: %v", diags)
			}
			if got := result.Id.ValueString(); got != tc.wantId {
				t.Fatalf("expected key %q, got %q", tc.wantId, got)
			}
			if got := result.Secret.ValueString(); got != "secret-"+tc.wantId {
				t.Fatalf("unexpected secret: %q", got)
			}
			dsn, diags := result.Dsn.Get(ctx)
			if diags.HasError() {
				t.Fatalf("get dsn: %v", diags)
			}
			if got := dsn["public"]; got != "https://public-"+tc.wantId+"@example.com/1" {
				t.Fatalf("unexpected public DSN: %q", got)
			}
		})
	}
}

func TestClientKeyEphemeralResource_echo(t *testing.T) {
	t.Parallel()
	server := newClientKeyEphemeralTestServer(t)

	for _, tc := range []struct {
		name   string
		extras string
		wantId string
	}{
		{name: "id", extras: `id = "bbb"`, wantId: "bbb"},
		{name: "name", extras: `name = "other"`, wantId: "bbb"},
		{name: "first", extras: `first = true`, wantId: "aaa"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_10_0),
				},
				Steps: []resource.TestStep{
					{
						Config: `
provider "sentry" {
	base_url = "` + server.URL + `/api/"
	token    = "test"
}

ephemeral "sentry_key" "test" {
	organization = "test-org"
	project      = "test-project"
	` + tc.extras + `
}

provider "echo" {
	data = ephemeral.sentry_key.test
}

resource "echo" "test" {}
`,
						ConfigStateChecks: []statecheck.StateCheck{
							expectNoEphemeralResourcesInState{},
							statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("id"), knownvalue.StringExact(tc.wantId)),
							statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("secret"), knownvalue.StringExact("secret-"+tc.wantId)),
							statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("project_id"), knownvalue.StringExact("1")),
							statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("dsn").AtMapKey("public"), knownvalue.StringExact("https://public-"+tc.wantId+"@example.com/1")),
						},
					},
				},
			})
		})
	}
}

func TestClientKeyEphemeralResource_conflictingLookups(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "sentry_key" "test" {
	organization = "test-org"
	project      = "test-project"
	id           = "aaa"
	name         = "Default"
}
`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}
