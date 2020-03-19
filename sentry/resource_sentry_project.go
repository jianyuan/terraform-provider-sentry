package sentry

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
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
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the project belongs to",
			},
			"team": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the team to create the project for",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the project",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional slug for this project",
				Computed:    true,
			},
			"platform": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional platform for this project",
				Computed:    true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_bookmarked": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "is_bookmarked is no longer used",
			},
			"color": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"digests_min_delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"digests_max_delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// TODO: Project options
		},
	}
}

func resourceSentryProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	team := d.Get("team").(string)
	params := &sentry.CreateProjectParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	proj, _, err := client.Projects.Create(org, team, params)
	if err != nil {
		return err
	}

	d.SetId(proj.Slug)
	resourceSentryProjectUpdate(d, meta)
	return resourceSentryProjectRead(d, meta)
}

func resourceSentryProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	proj, _, err := client.Projects.Get(org, slug)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(proj.Slug)
	d.Set("organization", proj.Organization.Slug)
	d.Set("team", proj.Team.Slug)
	d.Set("name", proj.Name)
	d.Set("slug", proj.Slug)
	d.Set("platform", proj.Platform)
	d.Set("project_id", proj.ID)
	d.Set("is_public", proj.IsPublic)
	d.Set("color", proj.Color)
	d.Set("features", proj.Features)
	d.Set("status", proj.Status)
	d.Set("digests_min_delay", proj.DigestsMinDelay)
	d.Set("digests_max_delay", proj.DigestsMaxDelay)

	// TODO: Project options

	return nil
}

func resourceSentryProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	params := &sentry.UpdateProjectParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	platform := d.Get("platform").(string)
	if platform != "" {
		params.Platform = platform
	}

	if v, ok := d.GetOk("digests_min_delay"); ok {
		params.DigestsMinDelay = Int(v.(int))
	}

	if v, ok := d.GetOk("digests_max_delay"); ok {
		params.DigestsMaxDelay = Int(v.(int))
	}

	proj, _, err := client.Projects.Update(org, slug, params)
	if err != nil {
		return err
	}

	d.SetId(proj.Slug)
	return resourceSentryProjectRead(d, meta)
}

func resourceSentryProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	_, err := client.Projects.Delete(org, slug)
	return err
}

func resourceSentryProjectImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
