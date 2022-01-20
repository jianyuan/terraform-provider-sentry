package sentry

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryKeyCreate,
		Read:   resourceSentryKeyRead,
		Update: resourceSentryKeyUpdate,
		Delete: resourceSentryKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeyImport,
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

func resourceSentryKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	_, resp, err := client.Projects.Get(org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return fmt.Errorf("project not found \"%v\": %w", project, err)
	}

	params := &sentry.CreateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	key, _, err := client.ProjectKeys.Create(org, project, params)
	if err != nil {
		return err
	}
	logging.Debugf("Created Sentry key with id %s", key.ID)
	d.SetId(key.ID)

	return resourceSentryKeyRead(d, meta)
}

func resourceSentryKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	keys, resp, err := client.ProjectKeys.List(org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
	}

	found := false

	for _, key := range keys {
		if key.ID == id {
			logging.Debugf("Found the sentry key with id %s", id)
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
		logging.Debugf("Sentry key with id %s could not be found...", id)
		d.SetId("")
	}

	return nil
}

func resourceSentryKeyUpdate(d *schema.ResourceData, meta interface{}) error {
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

	key, _, err := client.ProjectKeys.Update(org, project, id, params)
	if err != nil {
		return err
	}

	logging.Debugf("Updating current Sentry key id to %s", key.ID)
	d.SetId(key.ID)
	return resourceSentryKeyRead(d, meta)
}

func resourceSentryKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	_, err := client.ProjectKeys.Delete(org, project, id)
	logging.Debugf("Deleted Sentry key with id %s", id)
	return err
}
