package sentry

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func dataSourceSentryKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSentryKeyRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The slug of the organization the key should be created for.",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The slug of the project the key should be created for.",
			},
			"first": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				Description:   "Boolean flag indicating that we want the first key of the returned keys.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"first"},
				Description:   "The name of the key to retrieve.",
			},
			// Computed values.
			"public": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public key portion of the client key.",
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secret key portion of the client key.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the project that the key belongs to.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating the key is active.",
			},
			"rate_limit_window": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Length of time that will be considered when checking the rate limit.",
			},
			"rate_limit_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of events that can be reported within the rate limit window.",
			},
			"dsn_secret": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "DSN (Deprecated) for the key.",
			},
			"dsn_public": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DSN for the key.",
			},
			"dsn_csp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DSN for the Content Security Policy (CSP) for the key.",
			},
		},
	}
}

func dataSourceSentryKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry project keys", map[string]interface{}{
		"org":     org,
		"project": project,
	})
	keys, _, err := client.ProjectKeys.List(org, project)
	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		for _, key := range keys {
			if key.Name == name {
				return diag.FromErr(sentryKeyAttributes(d, &key))
			}
		}
		return diag.Errorf("Can't find Sentry key: %s", v)
	}

	if len(keys) == 1 {
		tflog.Debug(ctx, "sentry_key - single key", map[string]interface{}{
			"keyName": keys[0].Name,
			"keyID":   keys[0].ID,
		})
		return diag.FromErr(sentryKeyAttributes(d, &keys[0]))
	}

	first := d.Get("first").(bool)
	tflog.Debug(ctx, "sentry_key - multiple results found", map[string]interface{}{
		"first": first,
	})
	if first {
		// Sort keys by date created
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].DateCreated.Before(keys[j].DateCreated)
		})

		tflog.Debug(ctx, "sentry_key - Found more than one key. Returning the oldest (`first`) key.", map[string]interface{}{
			"keyName": keys[0].Name,
			"keyID":   keys[0].ID,
		})
		return diag.FromErr(sentryKeyAttributes(d, &keys[0]))
	}

	return diag.Errorf("There are %d keys associate to this project. "+
		"To avoid ambiguity, please set `first` to true or filter the keys by specifying a `name`.",
		len(keys))
}

func sentryKeyAttributes(d *schema.ResourceData, key *sentry.ProjectKey) error {
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

	return nil
}
