package sentry

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func resourceSentryProject() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Project resource.",

		CreateContext: resourceSentryProjectCreate,
		ReadContext:   resourceSentryProjectRead,
		UpdateContext: resourceSentryProjectUpdate,
		DeleteContext: resourceSentryProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSentryProjectImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"team": {
				Description: "The slug of the team to create the project for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name for the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The optional slug for this project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"platform": {
				Description: "The optional platform for this project.",
				Type:        schema.TypeString,
				Optional:    true,
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
				Deprecated: "is_bookmarked is no longer used",
				Type:       schema.TypeBool,
				Computed:   true,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"digests_min_delay": {
				Description: "The minimum amount of time (in seconds) to wait between scheduling digests for delivery after the initial scheduling.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"digests_max_delay": {
				Description: "The maximum amount of time (in seconds) to wait between scheduling digests for delivery.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"resolve_age": {
				Description: "Hours in which an issue is automatically resolve if not seen after this amount of time.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},

			// TODO: Project options
		},
	}
}

func resourceSentryProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	team := d.Get("team").(string)
	params := &sentry.CreateProjectParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	tflog.Debug(ctx, "Creating Sentry project", map[string]interface{}{
		"teamName": team,
		"org":      org,
	})
	proj, _, err := client.Projects.Create(ctx, org, team, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry project", map[string]interface{}{
		"projectSlug": proj.Slug,
		"projectID":   proj.ID,
		"team":        team,
		"org":         org,
	})

	d.SetId(proj.Slug)
	return resourceSentryProjectUpdate(ctx, d, meta)
}

func resourceSentryProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Reading Sentry project", map[string]interface{}{
		"projectSlug": slug,
		"org":         org,
	})
	proj, resp, err := client.Projects.Get(ctx, org, slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry project", map[string]interface{}{
		"projectSlug": proj.Slug,
		"projectID":   proj.ID,
		"org":         org,
	})

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

func resourceSentryProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		params.DigestsMinDelay = sentry.Int(v.(int))
	}

	if v, ok := d.GetOk("digests_max_delay"); ok {
		params.DigestsMaxDelay = sentry.Int(v.(int))
	}

	if v, ok := d.GetOk("resolve_age"); ok {
		params.ResolveAge = sentry.Int(v.(int))
	}

	tflog.Debug(ctx, "Updating Sentry project", map[string]interface{}{
		"projectSlug": slug,
		"org":         org,
	})
	proj, _, err := client.Projects.Update(ctx, org, slug, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry project", map[string]interface{}{
		"projectSlug": proj.Slug,
		"projectID":   proj.ID,
		"org":         org,
	})

	d.SetId(proj.Slug)
	return resourceSentryProjectRead(ctx, d, meta)
}

func resourceSentryProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Deleting Sentry project", map[string]interface{}{
		"projectSlug": slug,
		"org":         org,
	})
	_, err := client.Projects.Delete(ctx, org, slug)
	tflog.Debug(ctx, "Deleted Sentry project", map[string]interface{}{
		"projectSlug": slug,
		"org":         org,
	})

	return diag.FromErr(err)
}

func resourceSentryProjectImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	tflog.Debug(ctx, "Importing Sentry project", map[string]interface{}{
		"projectID": addrID,
	})

	parts := strings.Split(addrID, "/")

	if len(parts) != 2 {
		return nil, errors.New("Project import requires an ADDR ID of the following schema org-slug/project-slug")
	}

	d.Set("organization", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
