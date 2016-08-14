package main

import "github.com/hashicorp/terraform/helper/schema"

func resourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryOrganizationCreate,
		Read:   resourceSentryOrganizationRead,
		Update: resourceSentryOrganizationUpdate,
		Delete: resourceSentryOrganizationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The human readable name for the new organization",
			},
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique URL slug for this organization",
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	params := &CreateOrganizationParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	org, _, err := client.CreateOrganization(params)
	if err != nil {
		return err
	}

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(d, meta)
}

func resourceSentryOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()

	org, _, err := client.GetOrganization(slug)
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
	client := meta.(*Client)

	slug := d.Id()

	params := &UpdateOrganizationParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	org, _, err := client.UpdateOrganization(slug, params)
	if err != nil {
		return err
	}

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(d, meta)
}

func resourceSentryOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()

	_, err := client.DeleteOrganization(slug)
	return err
}
