package sentry

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryDefaultKey() *schema.Resource {
	// reuse read and update operations
	dKey := resourceSentryKey()
	dKey.CreateContext = resourceSentryDefaultKeyCreate
	dKey.DeleteContext = resourceAwsDefaultVpcDelete

	// Key name is a computed resource for default key
	dKey.Schema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    true,
		Description: "The name of the key",
	}

	return dKey
}

func resourceSentryDefaultKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	keys, resp, err := client.ProjectKeys.List(org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}

	if len(keys) < 1 {
		return diag.Errorf("Default key not found on the project")
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].DateCreated.Before(keys[j].DateCreated)
	})

	id := keys[0].ID
	params := &sentry.UpdateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	tflog.Debug(ctx, "Creating Sentry default key", map[string]interface{}{
		"org":     org,
		"project": project,
		"keyID":   id,
	})
	key, resp, err := client.ProjectKeys.Update(org, project, id, params)
	tflog.Debug(ctx, "Sentry default key create http response data", logging.ExtractHttpResponse(resp))
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry default key", map[string]interface{}{
		"org":     org,
		"project": project,
		"keyID":   key.ID,
	})

	d.SetId(id)
	return resourceSentryKeyRead(ctx, d, meta)
}

func resourceAwsDefaultVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Warn(ctx, "Cannot destroy Default Key. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
