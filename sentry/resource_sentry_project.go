package sentry

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

var defaultKey = sentry.ProjectKey{
	Name:  "Default",
	Label: "Default",
}

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
			"resolve_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Hours in which an issue is automatically resolve if not seen after this amount of time.",
				Computed:    true,
			},
			"remove_default_key": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Remove the default project key that Sentry creates for new projects.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
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

	if d.Get("remove_default_key").(bool) {
		keys, _, err := client.ProjectKeys.List(org, proj.Slug)
		if err != nil {
			return err
		}

		for _, key := range keys {
			if key.Name == defaultKey.Name && key.Label == defaultKey.Label {
				_, err := client.ProjectKeys.Delete(org, proj.Slug, key.ID)
				if err != nil {
					return err
				}

				break
			}
		}
	}

	return resourceSentryProjectUpdate(d, meta)
}

func resourceSentryProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	proj, resp, err := client.Projects.Get(org, slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
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
	d.Set("resolve_age", proj.ResolveAge)

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

	if v, ok := d.GetOk("resolve_age"); ok {
		params.ResolveAge = Int(v.(int))
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
