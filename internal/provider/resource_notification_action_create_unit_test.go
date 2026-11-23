package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/must"
)

func TestNotificationActionResource_createFailureAfterActionExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := r.Method + " " + r.URL.Path
		w.Header().Set("Content-Type", "application/json")

		switch request {
		case "POST /api/0/organizations/my-org/notifications/actions/":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 42, "triggerType": "spike-protection", "serviceType": "sentry_notification", "targetIdentifier": "default", "targetDisplay": "default", "projects": [1234]}`))
		case "GET /api/0/organizations/my-org/projects/":
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"detail": "failed"}`))
		default:
			t.Errorf("unexpected request: %s", request)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	ctx := context.Background()
	r := &NotificationActionResource{baseResource: baseResource{
		client:    must.Get(sentry.NewOnPremiseClient(server.URL+"/api/", http.DefaultClient)),
		apiClient: must.Get(apiclient.NewClientWithResponses(server.URL + "/api/")),
	}}

	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	planValues := make(map[string]tftypes.Value)
	for name, attrType := range objectType.AttributeTypes {
		planValues[name] = tftypes.NewValue(attrType, nil)
	}
	planValues["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	planValues["organization"] = tftypes.NewValue(tftypes.String, "my-org")
	planValues["trigger_type"] = tftypes.NewValue(tftypes.String, "spike-protection")
	planValues["service_type"] = tftypes.NewValue(tftypes.String, "sentry_notification")
	planValues["target_identifier"] = tftypes.NewValue(tftypes.String, "default")
	planValues["target_display"] = tftypes.NewValue(tftypes.String, "default")
	planValues["projects"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
		tftypes.NewValue(tftypes.String, "my-project"),
	})

	resp := fwresource.CreateResponse{
		State: tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, nil)},
	}
	r.Create(ctx, fwresource.CreateRequest{
		Plan: tfsdk.Plan{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, planValues)},
	}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error")
	}
	if got := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(got, "Unable to read projects") {
		t.Fatalf("unexpected error: %s", got)
	}

	var id, organization types.String
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("id"), &id)...)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("organization"), &organization)...)
	if id.ValueString() != "42" || organization.ValueString() != "my-org" {
		t.Fatalf("expected the created action to be in the state, got id=%s organization=%s", id, organization)
	}
}
