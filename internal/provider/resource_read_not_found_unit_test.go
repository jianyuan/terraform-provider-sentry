package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

// readResourceAgainst configures r with clients pointing at baseURL, then calls
// Read with a prior state built from attrs (all other attributes are null).
func readResourceAgainst(t *testing.T, r resource.Resource, baseURL string, attrs map[string]string) *resource.ReadResponse {
	t.Helper()
	ctx := context.Background()

	client, err := sentry.NewOnPremiseClient(baseURL, http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	apiClient, err := apiclient.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatal(err)
	}

	configureResp := &resource.ConfigureResponse{}
	r.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{
		ProviderData: &providerdata.ProviderData{Client: client, ApiClient: apiClient},
	}, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configure: %v", configureResp.Diagnostics)
	}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}

	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		if v, ok := attrs[name]; ok {
			values[name] = tftypes.NewValue(tftypes.String, v)
		} else {
			values[name] = tftypes.NewValue(typ, nil)
		}
	}

	state := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, values)}
	resp := &resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, resp)
	return resp
}

func TestResourceRead_removesResourceDeletedOutsideTerraform(t *testing.T) {
	notFound := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"The requested resource does not exist"}`))
	}))
	t.Cleanup(notFound.Close)

	emptyList := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	t.Cleanup(emptyList.Close)

	project := map[string]string{"organization": "my-org", "project": "my-project"}
	with := func(extra map[string]string) map[string]string {
		m := map[string]string{}
		for k, v := range project {
			m[k] = v
		}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}

	for _, tc := range []struct {
		name     string
		resource func() resource.Resource
		server   *httptest.Server
		attrs    map[string]string
	}{
		{"sentry_key/404", NewClientKeyResource, notFound, with(map[string]string{"id": "abc"})},
		{"sentry_project_spike_protection/404", NewProjectSpikeProtectionResource, notFound, project},
		{"sentry_project_inbound_data_filter/404", NewProjectInboundDataFilterResource, notFound, with(map[string]string{"filter_id": "browser-extensions"})},
		{"sentry_project_inbound_data_filter/missing", NewProjectInboundDataFilterResource, emptyList, with(map[string]string{"filter_id": "browser-extensions"})},
		{"sentry_project_ownership/404", NewProjectOwnershipResource, notFound, project},
		{"sentry_project_symbol_source/404", NewProjectSymbolSourcesResource, notFound, with(map[string]string{"id": "abc"})},
		{"sentry_project_symbol_source/missing", NewProjectSymbolSourcesResource, emptyList, with(map[string]string{"id": "abc"})},
		{"sentry_integration_pagerduty/404", NewIntegrationPagerDuty, notFound, map[string]string{"organization": "my-org", "integration_id": "123", "id": "1"}},
		{"sentry_integration_opsgenie/404", NewIntegrationOpsgenie, notFound, map[string]string{"organization": "my-org", "integration_id": "123", "id": "1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := readResourceAgainst(t, tc.resource(), tc.server.URL+"/api/", tc.attrs)
			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no error, got: %v", resp.Diagnostics)
			}
			if !resp.State.Raw.IsNull() {
				t.Fatal("expected resource to be removed from state")
			}
		})
	}
}

func TestResourceRead_transportErrorDoesNotPanic(t *testing.T) {
	closed := httptest.NewServer(http.NotFoundHandler())
	baseURL := closed.URL + "/api/"
	closed.Close()

	project := map[string]string{"organization": "my-org", "project": "my-project"}

	for _, tc := range []struct {
		name     string
		resource func() resource.Resource
		attrs    map[string]string
	}{
		{"sentry_project_inbound_data_filter", NewProjectInboundDataFilterResource, map[string]string{"organization": "my-org", "project": "my-project", "filter_id": "browser-extensions"}},
		{"sentry_project_ownership", NewProjectOwnershipResource, project},
		{"sentry_project_symbol_source", NewProjectSymbolSourcesResource, map[string]string{"organization": "my-org", "project": "my-project", "id": "abc"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := readResourceAgainst(t, tc.resource(), baseURL, tc.attrs)
			if !resp.Diagnostics.HasError() {
				t.Fatal("expected a client error")
			}
		})
	}
}
