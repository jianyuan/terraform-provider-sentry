package sentry

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func resourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryOrganizationCreate,
		Read:   resourceSentryOrganizationRead,
		Update: resourceSentryOrganizationUpdate,
		Delete: resourceSentryOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The human readable name for the organization",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique URL slug for this organization",
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	params := &sentry.CreateOrganizationParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	org, _, err := client.Organizations.Create(params)
	if err != nil {
		return err
	}

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(d, meta)
}

func resourceSentryOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()

	org, _, err := client.Organizations.Get(slug)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(org.Slug)
	d.Set("internal_id", org.ID)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)
	return nil
}

func resourceSentryOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	params := &sentry.UpdateOrganizationParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	org, _, err := client.Organizations.Update(slug, params)
	if err != nil {
		return err
	}

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(d, meta)
}

func resourceSentryOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()

	_, err := client.Organizations.Delete(slug)
	return err
}
