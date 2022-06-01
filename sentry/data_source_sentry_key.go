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

func dataSourceSentryKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSentryKeyRead,

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
			"first": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"first"},
			},
			// Computed values.
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
				Computed: true,
			},
			"rate_limit_count": {
				Type:     schema.TypeInt,
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

func dataSourceSentryKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry project keys", map[string]interface{}{
		"org":     org,
		"project": project,
	})
	keys, resp, err := client.ProjectKeys.List(org, project)
	tflog.Debug(ctx, "Sentry key read http response data", logging.ExtractHttpResponse(resp))
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
