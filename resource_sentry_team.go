package main

import "github.com/hashicorp/terraform/helper/schema"

func resourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryTeamCreate,
		Read:   resourceSentryTeamRead,
		Update: resourceSentryTeamUpdate,
		Delete: resourceSentryTeamDelete,

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the team should be created for",
			},
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

func resourceSentryTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	org := d.Get("organization").(string)
	params := &CreateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	team, _, err := client.CreateTeam(org, params)
	if err != nil {
		return err
	}

	d.SetId(team.Slug)
	return resourceSentryTeamRead(d, meta)
}

func resourceSentryTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	team, _, err := client.GetTeam(org, slug)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(team.Slug)
	d.Set("internal_id", team.ID)
	d.Set("name", team.Name)
	d.Set("slug", team.Slug)
	return nil
}

func resourceSentryTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	params := &UpdateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	team, _, err := client.UpdateTeam(org, slug, params)
	if err != nil {
		return err
	}

	d.SetId(team.Slug)
	return resourceSentryTeamRead(d, meta)
}

func resourceSentryTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	_, err := client.DeleteTeam(org, slug)
	return err
}
