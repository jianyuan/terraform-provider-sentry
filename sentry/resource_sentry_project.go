package sentry

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-multierror"
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
			StateContext: importOrganizationAndID,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"team": {
				Description: "The slug of the team to create the project for. One of 'team' or 'teams' must be set.",
				Type:        schema.TypeString,
				Deprecated:  "To be replaced by 'teams' in a future release.",
				Optional:    true,
			},
			"teams": {
				Description: "The slugs of the teams to create the project for. One of 'team' or 'teams' must be set.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ExactlyOneOf: []string{"team"},
				Optional:     true,
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
			"internal_id": {
				Description: "The internal ID for this project.",
				Type:        schema.TypeString,
				Computed:    true,
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
			"project_id": {
				Deprecated:  "Use `internal_id` instead.",
				Description: "Use `internal_id` instead.",
				Type:        schema.TypeString,
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
	var teams []interface{}
	if team == "" {
		// Since `Set.List()` produces deterministic ordering, `teams[0]` should always
		// resolve to the same value given the same `teams`.
		teams = d.Get("teams").(*schema.Set).List()
		team = teams[0].(string)
	}
	params := &sentry.CreateProjectParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	tflog.Debug(ctx, "Creating Sentry project", map[string]interface{}{
		"team":  team,
		"teams": teams,
		"org":   org,
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

	setTeams := func() error {
		if len(proj.Teams) <= 1 {
			return multierror.Append(
				d.Set("team", proj.Team.Slug),
				d.Set("teams", nil),
			)
		}

		teams := make([]string, len(proj.Teams))
		for i, team := range proj.Teams {
			teams[i] = *team.Slug
		}

		return multierror.Append(
			d.Set("team", nil),
			d.Set("teams", teams),
		)
	}

	d.SetId(proj.Slug)
	retErr := multierror.Append(
		d.Set("organization", proj.Organization.Slug),
		d.Set("name", proj.Name),
		d.Set("slug", proj.Slug),
		setTeams(),
		d.Set("platform", proj.Platform),
		d.Set("internal_id", proj.ID),
		d.Set("is_public", proj.IsPublic),
		d.Set("color", proj.Color),
		d.Set("features", proj.Features),
		d.Set("status", proj.Status),
		d.Set("digests_min_delay", proj.DigestsMinDelay),
		d.Set("digests_max_delay", proj.DigestsMaxDelay),
		d.Set("resolve_age", proj.ResolveAge),
		d.Set("project_id", proj.ID), // Deprecated
	)

	// TODO: Project options

	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	project := d.Id()
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

	tflog.Debug(ctx, "Updating project", map[string]interface{}{
		"org":     org,
		"project": project,
	})
	proj, _, err := client.Projects.Update(ctx, org, project, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(proj.Slug)

	oldTeams := map[string]struct{}{}
	newTeams := map[string]struct{}{}
	if d.HasChange("team") {
		oldTeam, newTeam := d.GetChange("team")
		if oldTeam.(string) != "" {
			oldTeams[oldTeam.(string)] = struct{}{}
		}
		if newTeam.(string) != "" {
			newTeams[newTeam.(string)] = struct{}{}
		}
	}

	if d.HasChange("teams") {
		o, n := d.GetChange("teams")
		for _, oldTeam := range o.(*schema.Set).List() {
			if oldTeam.(string) != "" {
				oldTeams[oldTeam.(string)] = struct{}{}
			}
		}
		for _, newTeam := range n.(*schema.Set).List() {
			if newTeam.(string) != "" {
				newTeams[newTeam.(string)] = struct{}{}
			}
		}
	}

	// Ensure old teams and new teams do not overlap.
	for newTeam := range newTeams {
		if _, exists := oldTeams[newTeam]; exists {
			delete(oldTeams, newTeam)
		}
	}

	tflog.Debug(ctx, "Adding teams to project", map[string]interface{}{
		"org":        org,
		"project":    project,
		"teamsToAdd": newTeams,
	})

	for newTeam := range newTeams {
		_, _, err = client.Projects.AddTeam(ctx, org, project, newTeam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	tflog.Debug(ctx, "Removing teams from project", map[string]interface{}{
		"org":           org,
		"project":       project,
		"teamsToRemove": oldTeams,
	})

	for oldTeam := range oldTeams {
		resp, err := client.Projects.RemoveTeam(ctx, org, project, oldTeam)
		if err != nil {
			if resp.Response.StatusCode != http.StatusNotFound {
				return diag.FromErr(err)
			}
		}
	}

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
