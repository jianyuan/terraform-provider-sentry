package sentry

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryTeamCreate,
		Read:   resourceSentryTeamRead,
		Update: resourceSentryTeamUpdate,
		Delete: resourceSentryTeamDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSentryTeamImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the team should be created for",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team",
			},
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional slug for this team",
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
	d.Set("organization", team.Organization.Slug)
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

func resourceSentryTeamImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	log.Printf("[DEBUG] Importing key using ADDR ID %s", addrID)

	parts := strings.Split(addrID, "/")

	if len(parts) != 2 {
		return nil, errors.New("Project import requires an ADDR ID of the following schema org-slug/project-slug")
	}

	d.Set("organization", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
