package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func cronMonitorDetectorJSON(failureIssueThreshold, recoveryThreshold string) []byte {
	return fmt.Appendf(nil, `{
		"id": "1234567",
		"name": "nightly-backup",
		"type": "monitor_check_in_failure",
		"projectId": "4500000000000001",
		"enabled": true,
		"description": null,
		"owner": null,
		"config": {},
		"conditionGroup": null,
		"dateCreated": "2026-11-01T00:00:00Z",
		"dateUpdated": "2026-11-01T00:00:00Z",
		"dataSources": [{
			"type": "cron_monitor",
			"queryObj": {
				"name": "nightly-backup",
				"config": {
					"schedule_type": "interval",
					"schedule": [60, "minute"],
					"checkin_margin": 3,
					"max_runtime": 30,
					"timezone": "UTC",
					"failure_issue_threshold": %s,
					"recovery_threshold": %s
				}
			}
		}]
	}`, failureIssueThreshold, recoveryThreshold)
}

func TestCronMonitorFillThresholds(t *testing.T) {
	// Sentry stores null for a threshold that was never set, for example on a
	// monitor created by an SDK check-in upsert, and treats null as 1.
	for _, tt := range []struct {
		name                  string
		failureIssueThreshold string
		recoveryThreshold     string
		wantFailure           int64
		wantRecovery          int64
	}{
		{name: "null", failureIssueThreshold: "null", recoveryThreshold: "null", wantFailure: 1, wantRecovery: 1},
		{name: "set", failureIssueThreshold: "3", recoveryThreshold: "4", wantFailure: 3, wantRecovery: 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var monitor apiclient.ProjectMonitor
			if err := json.Unmarshal(cronMonitorDetectorJSON(tt.failureIssueThreshold, tt.recoveryThreshold), &monitor); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			var resourceModel CronMonitorResourceModel
			if diags := resourceModel.Fill(ctx, monitor); diags.HasError() {
				t.Fatalf("resource Fill: %v", diags)
			}
			if got := resourceModel.FailureIssueThreshold.Get(); got != tt.wantFailure {
				t.Errorf("resource failure_issue_threshold = %d, want %d", got, tt.wantFailure)
			}
			if got := resourceModel.RecoveryThreshold.Get(); got != tt.wantRecovery {
				t.Errorf("resource recovery_threshold = %d, want %d", got, tt.wantRecovery)
			}

			var dataSourceModel CronMonitorDataSourceModel
			if diags := dataSourceModel.fill(ctx, monitor); diags.HasError() {
				t.Fatalf("data source fill: %v", diags)
			}
			if got := dataSourceModel.FailureIssueThreshold.Get(); got != tt.wantFailure {
				t.Errorf("data source failure_issue_threshold = %d, want %d", got, tt.wantFailure)
			}
			if got := dataSourceModel.RecoveryThreshold.Get(); got != tt.wantRecovery {
				t.Errorf("data source recovery_threshold = %d, want %d", got, tt.wantRecovery)
			}
		})
	}
}
