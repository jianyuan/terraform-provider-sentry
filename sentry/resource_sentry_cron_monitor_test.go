
package sentry

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceSentryCronMonitor(t *testing.T) {
    r := resourceSentryCronMonitor()
    if r.CreateContext == nil || r.ReadContext == nil || r.UpdateContext == nil || r.DeleteContext == nil {
        t.Fatal("missing CRUD funcs")
    }

    expectedSchema := []string{"organization", "project", "name", "schedule_type", "schedule"}
    for _, k := range expectedSchema {
        if _, ok := r.Schema[k]; !ok {
            t.Fatalf("missing expected schema field: %s", k)
        }
    }
}
