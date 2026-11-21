package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

// unreachableBaseURL returns the URL of a server that has already been
// closed, so every request fails before an HTTP response exists.
func unreachableBaseURL(t *testing.T) string {
	t.Helper()
	srv := httptest.NewServer(http.NotFoundHandler())
	baseURL := srv.URL + "/api/"
	srv.Close()
	return baseURL
}

// configuredResourceState configures r against baseURL and returns a state
// built from attrs (all other attributes are null).
func configuredResourceState(t *testing.T, r resource.Resource, baseURL string, attrs map[string]string) (schema.Schema, tftypes.Value) {
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

	return schemaResp.Schema, tftypes.NewValue(objType, values)
}

func TestResource_transportErrorIsReportedNotPanicked(t *testing.T) {
	baseURL := unreachableBaseURL(t)

	project := map[string]string{"organization": "my-org", "project": "my-project", "id": "1"}

	resources := []struct {
		name     string
		resource func() resource.Resource
		attrs    map[string]string
		read     bool
		update   bool
	}{
		{"notification_action", NewNotificationActionResource, map[string]string{"organization": "my-org", "id": "1"}, true, true},
		{"organization_repository", NewOrganizationRepositoryResource, map[string]string{"organization": "my-org", "id": "1"}, false, false},
		{"project_inbound_data_filter", NewProjectInboundDataFilterResource, project, false, false},
		{"project_ownership", NewProjectOwnershipResource, project, false, false},
		{"project_symbol_source", NewProjectSymbolSourcesResource, project, false, false},
	}

	for _, tc := range resources {
		t.Run(tc.name+"/delete", func(t *testing.T) {
			ctx := context.Background()
			r := tc.resource()
			sch, raw := configuredResourceState(t, r, baseURL, tc.attrs)
			state := tfsdk.State{Schema: sch, Raw: raw}

			resp := &resource.DeleteResponse{State: state}
			r.Delete(ctx, resource.DeleteRequest{State: state}, resp)

			if !resp.Diagnostics.HasError() {
				t.Fatal("expected an error diagnostic for an unreachable API")
			}
		})

		if tc.read {
			t.Run(tc.name+"/read", func(t *testing.T) {
				ctx := context.Background()
				r := tc.resource()
				sch, raw := configuredResourceState(t, r, baseURL, tc.attrs)
				state := tfsdk.State{Schema: sch, Raw: raw}

				resp := &resource.ReadResponse{State: state}
				r.Read(ctx, resource.ReadRequest{State: state}, resp)

				if !resp.Diagnostics.HasError() {
					t.Fatal("expected an error diagnostic for an unreachable API")
				}
				if resp.State.Raw.IsNull() {
					t.Fatal("the resource must stay in state when the API cannot be reached")
				}
			})
		}

		if tc.update {
			t.Run(tc.name+"/update", func(t *testing.T) {
				ctx := context.Background()
				r := tc.resource()
				sch, raw := configuredResourceState(t, r, baseURL, tc.attrs)
				state := tfsdk.State{Schema: sch, Raw: raw}
				plan := tfsdk.Plan{Schema: sch, Raw: raw}

				resp := &resource.UpdateResponse{State: state}
				r.Update(ctx, resource.UpdateRequest{State: state, Plan: plan}, resp)

				if !resp.Diagnostics.HasError() {
					t.Fatal("expected an error diagnostic for an unreachable API")
				}
			})
		}
	}
}
