
package sentry

import (
    "context"
    "fmt"
    "strings"

    sentryapi "github.com/getsentry/sentry-go"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryCronMonitor() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceSentryCronMonitorCreate,
        ReadContext:   resourceSentryCronMonitorRead,
        UpdateContext: resourceSentryCronMonitorUpdate,
        DeleteContext: resourceSentryCronMonitorDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: map[string]*schema.Schema{
            "organization": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "project": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
            "schedule_type": {
                Type:     schema.TypeString,
                Required: true,
            },
            "schedule": {
                Type:     schema.TypeString,
                Required: true,
            },
            "check_in_margin": {
                Type:     schema.TypeInt,
                Optional: true,
            },
            "max_runtime": {
                Type:     schema.TypeInt,
                Optional: true,
            },
            "timezone": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "alert_rule_enabled": {
                Type:     schema.TypeBool,
                Optional: true,
            },
            "alert_threshold": {
                Type:     schema.TypeInt,
                Optional: true,
            },
        },
    }
}

func resourceSentryCronMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*sentryapi.Client)

    org := d.Get("organization").(string)
    proj := d.Get("project").(string)

    monitorData := map[string]interface{}{
        "name":          d.Get("name").(string),
        "type":          "cron_job",
        "config": map[string]interface{}{
            "schedule_type":     d.Get("schedule_type").(string),
            "schedule":          d.Get("schedule").(string),
            "checkin_margin":    d.Get("check_in_margin").(int),
            "max_runtime":       d.Get("max_runtime").(int),
            "timezone":          d.Get("timezone").(string),
        },
    }

    if d.Get("alert_rule_enabled").(bool) {
        monitorData["alert_rule"] = map[string]interface{}{
            "time_window": d.Get("alert_threshold").(int),
        }
    }

    var result map[string]interface{}
    path := fmt.Sprintf("/api/0/projects/%s/%s/monitors/", org, proj)
    if err := client.Request(ctx, "POST", path, monitorData, &result); err != nil {
        return diag.FromErr(err)
    }

    slug := result["slug"].(string)
    d.SetId(fmt.Sprintf("%s/%s/%s", org, proj, slug))
    return resourceSentryCronMonitorRead(ctx, d, meta)
}

func resourceSentryCronMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*sentryapi.Client)

    org, proj, slug := splitCronMonitorID(d.Id())

    var monitor map[string]interface{}
    path := fmt.Sprintf("/api/0/projects/%s/%s/monitors/%s/", org, proj, slug)
    if err := client.Request(ctx, "GET", path, nil, &monitor); err != nil {
        return diag.FromErr(err)
    }

    d.Set("organization", org)
    d.Set("project", proj)
    d.Set("name", monitor["name"])
    config := monitor["config"].(map[string]interface{})
    d.Set("schedule_type", config["schedule_type"])
    d.Set("schedule", config["schedule"])
    d.Set("check_in_margin", config["checkin_margin"])
    d.Set("max_runtime", config["max_runtime"])
    d.Set("timezone", config["timezone"])

    if ar, ok := monitor["alert_rule"].(map[string]interface{}); ok {
        d.Set("alert_rule_enabled", true)
        d.Set("alert_threshold", ar["time_window"])
    } else {
        d.Set("alert_rule_enabled", false)
    }

    return nil
}

func resourceSentryCronMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*sentryapi.Client)
    org, proj, slug := splitCronMonitorID(d.Id())

    monitorData := map[string]interface{}{
        "name": d.Get("name").(string),
        "config": map[string]interface{}{
            "schedule_type":  d.Get("schedule_type").(string),
            "schedule":       d.Get("schedule").(string),
            "checkin_margin": d.Get("check_in_margin").(int),
            "max_runtime":    d.Get("max_runtime").(int),
            "timezone":       d.Get("timezone").(string),
        },
    }

    if d.Get("alert_rule_enabled").(bool) {
        monitorData["alert_rule"] = map[string]interface{}{
            "time_window": d.Get("alert_threshold").(int),
        }
    }

    path := fmt.Sprintf("/api/0/projects/%s/%s/monitors/%s/", org, proj, slug)
    if err := client.Request(ctx, "PUT", path, monitorData, nil); err != nil {
        return diag.FromErr(err)
    }

    return resourceSentryCronMonitorRead(ctx, d, meta)
}

func resourceSentryCronMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*sentryapi.Client)
    org, proj, slug := splitCronMonitorID(d.Id())

    path := fmt.Sprintf("/api/0/projects/%s/%s/monitors/%s/", org, proj, slug)
    if err := client.Request(ctx, "DELETE", path, nil, nil); err != nil {
        return diag.FromErr(err)
    }

    d.SetId("")
    return nil
}

func splitCronMonitorID(id string) (string, string, string) {
    parts := strings.Split(id, "/")
    return parts[0], parts[1], parts[2]
}
