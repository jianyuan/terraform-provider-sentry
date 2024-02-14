package sentry

import (
	"context"
	"sort"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSentryKey() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Key data source.",

		ReadContext: dataSourceSentryKeyRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the key should be created for.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"project": {
				Description: "The slug of the project the key should be created for.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"first": {
				Description:   "Boolean flag indicating that we want the first key of the returned keys.",
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Description:   "The name of the key to retrieve.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"first"},
			},
			// Computed values.
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
				Computed:    true,
			},
			"rate_limit_count": {
				Description: "Number of events that can be reported within the rate limit window.",
				Type:        schema.TypeInt,
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

func dataSourceSentryKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry project keys", map[string]interface{}{
		"org":     org,
		"project": project,
	})

	listParams := &sentry.ListCursorParams{}
	var allKeys []*sentry.ProjectKey
	for {
		keys, resp, err := client.ProjectKeys.List(ctx, org, project, listParams)
		if err != nil {
			return diag.FromErr(err)
		}
		allKeys = append(allKeys, keys...)
		if resp.Cursor == "" {
			break
		}
		listParams.Cursor = resp.Cursor
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		for _, key := range allKeys {
			if key.Name == name {
				return diag.FromErr(sentryKeyAttributes(d, key))
			}
		}
		return diag.Errorf("Can't find Sentry key: %s", v)
	}

	if len(allKeys) == 1 {
		tflog.Debug(ctx, "sentry_key - single key", map[string]interface{}{
			"keyName": allKeys[0].Name,
			"keyID":   allKeys[0].ID,
		})
		return diag.FromErr(sentryKeyAttributes(d, allKeys[0]))
	}

	first := d.Get("first").(bool)
	tflog.Debug(ctx, "sentry_key - multiple results found", map[string]interface{}{
		"first": first,
	})
	if first {
		// Sort keys by date created
		sort.Slice(allKeys, func(i, j int) bool {
			return allKeys[i].DateCreated.Before(allKeys[j].DateCreated)
		})

		tflog.Debug(ctx, "sentry_key - Found more than one key. Returning the oldest (`first`) key.", map[string]interface{}{
			"keyName": allKeys[0].Name,
			"keyID":   allKeys[0].ID,
		})
		return diag.FromErr(sentryKeyAttributes(d, allKeys[0]))
	}

	return diag.Errorf("There are %d keys associate to this project. "+
		"To avoid ambiguity, please set `first` to true or filter the keys by specifying a `name`.",
		len(allKeys))
}

func sentryKeyAttributes(d *schema.ResourceData, key *sentry.ProjectKey) error {
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
	return retErr.ErrorOrNil()
}
