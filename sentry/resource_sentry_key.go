package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryKeyCreate,
		ReadContext:   resourceSentryKeyRead,
		UpdateContext: resourceSentryKeyUpdate,
		DeleteContext: resourceSentryKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the key should be created for",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project the key should be created for",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the key",
			},
			"public": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rate_limit_window": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"rate_limit_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dsn_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsn_public": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsn_csp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSentryKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	_, resp, err := client.Projects.Get(org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.Errorf("project not found \"%v\": %v", project, err)
	}

	params := &sentry.CreateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	tflog.Debug(ctx, "Creating Sentry key", map[string]interface{}{
		"keyName": params.Name,
		"org":     org,
		"project": project,
	})
	key, resp, err := client.ProjectKeys.Create(org, project, params)
	tflog.Debug(ctx, "Sentry key create http response data", logging.ExtractHttpResponse(resp))
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry key", map[string]interface{}{
		"keyID":   key.ID,
		"keyName": key.Name,
		"org":     org,
		"project": project,
	})
	d.SetId(key.ID)

	return resourceSentryKeyRead(ctx, d, meta)
}

func resourceSentryKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry key", map[string]interface{}{
		"keyID":   id,
		"org":     org,
		"project": project,
	})
	keys, resp, err := client.ProjectKeys.List(org, project)
	tflog.Debug(ctx, "Sentry key read http response data", logging.ExtractHttpResponse(resp))
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "Read Sentry keys", map[string]interface{}{
		"keyCount": len(keys),
		"keys":     logging.TryJsonify(keys),
	})

	found := false

	for _, key := range keys {
		if key.ID == id {
			tflog.Debug(ctx, "Found Sentry key", map[string]interface{}{
				"keyID":   id,
				"org":     org,
				"project": project,
			})
			d.SetId(key.ID)
			d.Set("name", key.Name)
			d.Set("public", key.Public)
			d.Set("secret", key.Secret)
			d.Set("project_id", key.ProjectID)
			d.Set("is_active", key.IsActive)

			if key.RateLimit != nil {
				d.Set("rate_limit_window", key.RateLimit.Window)
				d.Set("rate_limit_count", key.RateLimit.Count)
			}

			d.Set("dsn_secret", key.DSN.Secret)
			d.Set("dsn_public", key.DSN.Public)
			d.Set("dsn_csp", key.DSN.CSP)

			found = true

			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Sentry key could not be found...", map[string]interface{}{
			"keyID": id,
		})
		d.SetId("")
	}

	return nil
}

func resourceSentryKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	params := &sentry.UpdateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	tflog.Debug(ctx, "Updating Sentry key", map[string]interface{}{
		"keyID": id,
	})
	key, resp, err := client.ProjectKeys.Update(org, project, id, params)
	tflog.Debug(ctx, "Sentry key update http response data", logging.ExtractHttpResponse(resp))
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry key", map[string]interface{}{
		"keyID": key.ID,
	})

	d.SetId(key.ID)
	return resourceSentryKeyRead(ctx, d, meta)
}

func resourceSentryKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Deleting Sentry key", map[string]interface{}{
		"keyID": id,
	})
	resp, err := client.ProjectKeys.Delete(org, project, id)
	tflog.Debug(ctx, "Sentry key delete http response data", logging.ExtractHttpResponse(resp))
	tflog.Debug(ctx, "Deleted Sentry key", map[string]interface{}{
		"keyID": id,
	})
	return diag.FromErr(err)
}
