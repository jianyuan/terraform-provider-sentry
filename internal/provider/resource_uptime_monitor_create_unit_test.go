package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func TestUptimeMonitorResourceCreate_preservesPlannedEnabled(t *testing.T) {
	ctx := context.Background()

	var config apiclient.ProjectMonitorConfig
	if err := config.FromProjectMonitorConfigUptimeDomainFailure(apiclient.ProjectMonitorConfigUptimeDomainFailure{
		Mode:              apiclient.ProjectMonitorConfigUptimeDomainFailureMode(1),
		Environment:       "production",
		RecoveryThreshold: 1,
		DowntimeThreshold: 3,
	}); err != nil {
		t.Fatal(err)
	}

	var dataSource apiclient.ProjectMonitorDataSourceWrapper
	if err := dataSource.FromProjectMonitorDataSourceWrapperUptimeSubscription(apiclient.ProjectMonitorDataSourceWrapperUptimeSubscription{
		QueryObj: apiclient.ProjectMonitorDataSourceUptimeDomainFailure{
			Url:             "https://example.com",
			Method:          "GET",
			Headers:         [][]string{},
			IntervalSeconds: 60,
			TimeoutMs:       5000,
		},
		Type: apiclient.ProjectMonitorDataSourceWrapperUptimeSubscriptionType("uptime_subscription"),
	}); err != nil {
		t.Fatal(err)
	}

	responseBody, err := json.Marshal(apiclient.ProjectMonitor{
		Config:      config,
		DataSources: []apiclient.ProjectMonitorDataSourceWrapper{dataSource},
		DateCreated: time.Date(2026, time.November, 27, 0, 0, 0, 0, time.UTC),
		DateUpdated: time.Date(2026, time.November, 27, 0, 0, 0, 0, time.UTC),
		Enabled:     false,
		Id:          "123",
		Name:        "uptime monitor",
		ProjectId:   "my-project",
		Type:        "uptime_domain_failure",
	})
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			var requestBody struct {
				Enabled bool `json:"enabled"`
			}
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				t.Errorf("decode create request: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !requestBody.Enabled {
				t.Error("expected create request to enable the monitor")
			}
			w.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
		default:
			t.Errorf("unexpected request method %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		_, _ = w.Write(responseBody)
	}))
	t.Cleanup(server.Close)

	apiClient, err := apiclient.NewClientWithResponses(server.URL + "/api/")
	if err != nil {
		t.Fatal(err)
	}

	r := &UptimeMonitorResource{baseResource: baseResource{apiClient: apiClient}}
	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}

	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	planValues := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		planValues[name] = tftypes.NewValue(attributeType, tftypes.UnknownValue)
	}
	planValues["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	planValues["organization"] = tftypes.NewValue(tftypes.String, "my-org")
	planValues["project"] = tftypes.NewValue(tftypes.String, "my-project")
	planValues["enabled"] = tftypes.NewValue(tftypes.Bool, true)
	planValues["name"] = tftypes.NewValue(tftypes.String, "uptime monitor")
	planValues["url"] = tftypes.NewValue(tftypes.String, "https://example.com")
	planValues["method"] = tftypes.NewValue(tftypes.String, "GET")
	planValues["interval_seconds"] = tftypes.NewValue(tftypes.Number, 60)
	planValues["timeout_ms"] = tftypes.NewValue(tftypes.Number, 5000)
	planValues["environment"] = tftypes.NewValue(tftypes.String, "production")
	planValues["recovery_threshold"] = tftypes.NewValue(tftypes.Number, 1)
	planValues["downtime_threshold"] = tftypes.NewValue(tftypes.Number, 3)

	plan := tfsdk.Plan{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(objectType, planValues),
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: schemaResp.Schema,
			Raw:    tftypes.NewValue(objectType, nil),
		},
	}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("create: %v", resp.Diagnostics)
	}
	if got := resp.Diagnostics.WarningsCount(); got != 1 {
		t.Fatalf("expected one warning for the mismatched enabled response, got %d", got)
	}

	var state UptimeMonitorResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		t.Fatalf("read created state: %v", resp.Diagnostics)
	}
	if got := state.Enabled.Get(); !got {
		t.Fatalf("expected enabled state to preserve the plan's true value, got %t", got)
	}

	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("read: %v", readResp.Diagnostics)
	}

	var refreshedState UptimeMonitorResourceModel
	readResp.Diagnostics.Append(readResp.State.Get(ctx, &refreshedState)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("read refreshed state: %v", readResp.Diagnostics)
	}
	if got := refreshedState.Enabled.Get(); got {
		t.Fatalf("expected refresh to report Sentry's disabled value, got %t", got)
	}
}
