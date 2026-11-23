package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

func TestProjectResourceCreate_sendsHighlightTags(t *testing.T) {
	ctx := context.Background()

	const project = `{"id":"1","slug":"my-project","name":"My Project","platform":"go","organization":{"id":"1","slug":"my-org"},"teams":[{"id":"1","slug":"my-team"}]}`

	var mu sync.Mutex
	var updateBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			_ = json.Unmarshal(body, &updateBody)
			mu.Unlock()
		}
		_, _ = w.Write([]byte(project))
	}))
	t.Cleanup(server.Close)

	apiClient, err := apiclient.NewClientWithResponses(server.URL + "/api/")
	if err != nil {
		t.Fatal(err)
	}

	r := NewProjectResource()
	configureResp := &resource.ConfigureResponse{}
	r.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{
		ProviderData: &providerdata.ProviderData{ApiClient: apiClient},
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
		values[name] = tftypes.NewValue(typ, tftypes.UnknownValue)
	}
	values["organization"] = tftypes.NewValue(tftypes.String, "my-org")
	values["name"] = tftypes.NewValue(tftypes.String, "My Project")
	values["platform"] = tftypes.NewValue(tftypes.String, "go")
	values["default_rules"] = tftypes.NewValue(tftypes.Bool, nil)
	values["default_key"] = tftypes.NewValue(tftypes.Bool, nil)
	values["teams"] = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
		tftypes.NewValue(tftypes.String, "my-team"),
	})
	values["highlight_tags"] = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
		tftypes.NewValue(tftypes.String, "foo"),
		tftypes.NewValue(tftypes.String, "bar"),
	})

	plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, values)}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, nil)}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	mu.Lock()
	defer mu.Unlock()
	if updateBody == nil {
		t.Fatalf("project was not updated after creation: %v", resp.Diagnostics)
	}
	got, ok := updateBody["highlightTags"].([]any)
	if !ok {
		t.Fatalf("expected highlightTags in the update request, got body: %v", updateBody)
	}
	tags := map[any]bool{}
	for _, tag := range got {
		tags[tag] = true
	}
	if len(got) != 2 || !tags["foo"] || !tags["bar"] {
		t.Fatalf("expected highlightTags [foo bar], got %v", got)
	}
}
