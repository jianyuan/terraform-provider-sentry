package sentry

import (
	"context"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryKey() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Key resource.",

		CreateContext: resourceSentryKeyCreate,
		ReadContext:   resourceSentryKeyRead,
		UpdateContext: resourceSentryKeyUpdate,
		DeleteContext: resourceSentryKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importOrganizationProjectAndID,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the key should be created for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The slug of the project the key should be created for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"public": {
				Description: "Public key portion of the client key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"secret": {
				Description: "Secret key portion of the client key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"project_id": {
				Description: "The ID of the project that the key belongs to.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"is_active": {
				Description: "Flag indicating the key is active.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"rate_limit_window": {
				Description: "Length of time that will be considered when checking the rate limit.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"rate_limit_count": {
				Description: "Number of events that can be reported within the rate limit window.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"dsn_secret": {
				Deprecated: "DSN (Deprecated) for the key.",
				Type:       schema.TypeString,
				Computed:   true,
			},
			"dsn_public": {
				Description: "DSN for the key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dsn_csp": {
				Description: "DSN for the Content Security Policy (CSP) for the key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	_, resp, err := client.Projects.Get(ctx, org, project)
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
	key, _, err := client.ProjectKeys.Create(ctx, org, project, params)
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

	listParams := &sentry.ListCursorParams{}
	var allKeys []*sentry.ProjectKey
	for {
		keys, resp, err := client.ProjectKeys.List(ctx, org, project, listParams)
		if found, err := checkClientGet(resp, err, d); !found {
			return diag.FromErr(err)
		}
		allKeys = append(allKeys, keys...)
		if resp.Cursor == "" {
			break
		}
		listParams.Cursor = resp.Cursor
	}
	tflog.Trace(ctx, "Read Sentry keys", map[string]interface{}{
		"keyCount": len(allKeys),
		"keys":     allKeys,
	})

	found := false

	for _, key := range allKeys {
		if key.ID == id {
			tflog.Debug(ctx, "Found Sentry key", map[string]interface{}{
				"keyID":   id,
				"org":     org,
				"project": project,
			})
			d.SetId(key.ID)
			retErr := multierror.Append(
				d.Set("name", key.Name),
				d.Set("public", key.Public),
				d.Set("secret", key.Secret),
				d.Set("project_id", key.ProjectID),
				d.Set("is_active", key.IsActive),
				d.Set("dsn_secret", key.DSN.Secret),
				d.Set("dsn_public", key.DSN.Public),
				d.Set("dsn_csp", key.DSN.CSP),
			)
			if key.RateLimit != nil {
				retErr = multierror.Append(
					retErr,
					d.Set("rate_limit_window", key.RateLimit.Window),
					d.Set("rate_limit_count", key.RateLimit.Count),
				)
			}
			if err := retErr.ErrorOrNil(); err != nil {
				return diag.FromErr(err)
			}

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
	key, _, err := client.ProjectKeys.Update(ctx, org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry key", map[string]interface{}{
		"keyID": id,
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
	_, err := client.ProjectKeys.Delete(ctx, org, project, id)
	tflog.Debug(ctx, "Deleted Sentry key", map[string]interface{}{
		"keyID": id,
	})
	return diag.FromErr(err)
}
