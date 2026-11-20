package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

func readClientKeyDataSourceAgainst(t *testing.T, baseURL string, attrs map[string]tftypes.Value) *datasource.ReadResponse {
	t.Helper()
	ctx := context.Background()

	apiClient, err := apiclient.NewClientWithResponses(baseURL)
	if err != nil {
		t.Fatal(err)
	}

	d := NewClientKeyDataSource()
	configureResp := &datasource.ConfigureResponse{}
	d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{
		ProviderData: &providerdata.ProviderData{ApiClient: apiClient},
	}, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configure: %v", configureResp.Diagnostics)
	}

	schemaResp := &datasource.SchemaResponse{}
	d.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
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
	resp := &datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	d.Read(ctx, datasource.ReadRequest{Config: config}, resp)
	return resp
}

func TestClientKeyDataSourceRead_projectWithoutKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	t.Cleanup(server.Close)

	project := map[string]tftypes.Value{
		"organization": tftypes.NewValue(tftypes.String, "my-org"),
		"project":      tftypes.NewValue(tftypes.String, "my-project"),
	}

	for _, tc := range []struct {
		name  string
		first *bool
	}{
		{name: "first unset"},
		{name: "first true", first: new(true)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			attrs := map[string]tftypes.Value{}
			for k, v := range project {
				attrs[k] = v
			}
			if tc.first != nil {
				attrs["first"] = tftypes.NewValue(tftypes.Bool, *tc.first)
			}

			resp := readClientKeyDataSourceAgainst(t, server.URL+"/api/", attrs)
			if !resp.Diagnostics.HasError() {
				t.Fatal("expected an error")
			}
			if got := resp.Diagnostics.Errors()[0].Detail(); got != "No matching client key found" {
				t.Fatalf("unexpected error: %s", got)
			}
		})
	}
}

func TestClientKeyDataSourceRead_idNotFoundReportsStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"The requested resource does not exist"}`))
	}))
	t.Cleanup(server.Close)

	resp := readClientKeyDataSourceAgainst(t, server.URL+"/api/", map[string]tftypes.Value{
		"organization": tftypes.NewValue(tftypes.String, "my-org"),
		"project":      tftypes.NewValue(tftypes.String, "my-project"),
		"id":           tftypes.NewValue(tftypes.String, "abc"),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error")
	}
	if got := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(got, "got status 404") {
		t.Fatalf("expected the error to report the 404 status, got: %s", got)
	}
}
