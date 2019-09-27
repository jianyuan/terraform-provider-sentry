package sentry

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func dataSourceSentryKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSentryKeyRead,

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

func dataSourceSentryKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	keys, _, err := client.ProjectKeys.List(org, project)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		for _, key := range keys {
			if key.Name == name {
				return sentryKeyAttributes(d, &key)
			}
		}
		return fmt.Errorf("Can't find Sentry key: %s", v)
	}

	if len(keys) == 1 {
		log.Printf("[DEBUG] sentry_key - single key found: %s", keys[0].ID)
		return sentryKeyAttributes(d, &keys[0])
	}

	first := d.Get("first").(bool)
	log.Printf("[DEBUG] sentry_key - multiple results found and `first` is set to: %t", first)
	if first {
		return sentryKeyAttributes(d, &keys[0])
	}

	return fmt.Errorf("There are %d keys associate to this project. "+
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
