package sentry

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSentryProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryProjectCreate,
		Read:   resourceSentryProjectRead,
		Update: resourceSentryProjectUpdate,
		Delete: resourceSentryProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSentryProjectImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the project belongs to",
			},
			"team": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the team to create the project for",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the project",
			},
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional slug for this project",
				Computed:    true,
			},
			"resolve_age": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of hours after issues are automatically resolved",
			},
		},
	}
}

func resourceSentryProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	org := d.Get("organization").(string)
	team := d.Get("team").(string)
	params := &CreateProjectParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	proj, _, err := client.CreateProject(org, team, params)
	if err != nil {
		return err
	}

	d.SetId(proj.Slug)
	return resourceSentryProjectRead(d, meta)
}

func resourceSentryProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	proj, _, err := client.GetProject(org, slug)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(proj.Slug)
	d.Set("internal_id", proj.ID)
	d.Set("name", proj.Name)
	d.Set("slug", proj.Slug)
	d.Set("organization", proj.Organization.Slug)
	d.Set("team", proj.Team.Slug)
	d.Set("resolve_age", proj.Options.ResolveAge)
	return nil
}

func resourceSentryProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	params := &UpdateProjectParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
		Options: UpdateProjectOptionsParams{
			ResolveAge: d.Get("resolve_age").(int),
		},
	}

	proj, _, err := client.UpdateProject(org, slug, params)
	if err != nil {
		return err
	}

	d.SetId(proj.Slug)
	return resourceSentryProjectRead(d, meta)
}

func resourceSentryProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	_, err := client.DeleteProject(org, slug)
	return err
}

func resourceSentryProjectImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	log.Printf("[DEBUG] Importing key using ADDR ID %s", addrID)

	parts := strings.Split(addrID, "/")

	if len(parts) != 2 {
		return nil, errors.New("Project import requires an ADDR ID of the following schema org-slug/team-slug")
	}

	d.Set("organization", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
