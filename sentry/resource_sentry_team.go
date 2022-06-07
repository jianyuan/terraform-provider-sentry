package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func resourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Team resource.",

		CreateContext: resourceSentryTeamCreate,
		ReadContext:   resourceSentryTeamRead,
		UpdateContext: resourceSentryTeamUpdate,
		DeleteContext: resourceSentryTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importOrganizationAndID,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the team should be created for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the team.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The optional slug for this team.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_pending": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_member": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSentryTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	params := &sentry.CreateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	tflog.Debug(ctx, "Creating Sentry team", map[string]interface{}{
		"teamName": params.Name,
		"org":      org,
	})
	team, _, err := client.Teams.Create(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry team", map[string]interface{}{
		"teamName": team.Name,
		"org":      org,
	})

	d.SetId(team.Slug)
	return resourceSentryTeamRead(ctx, d, meta)
}

func resourceSentryTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Reading Sentry team", map[string]interface{}{
		"teamSlug": slug,
		"org":      org,
	})
	team, resp, err := client.Teams.Get(ctx, org, slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry team", map[string]interface{}{
		"teamSlug": team.Slug,
		"teamID":   team.ID,
		"org":      org,
	})

	d.SetId(team.Slug)
	d.Set("team_id", team.ID)
	d.Set("name", team.Name)
	d.Set("slug", team.Slug)
	d.Set("organization", org)
	d.Set("has_access", team.HasAccess)
	d.Set("is_pending", team.IsPending)
	d.Set("is_member", team.IsMember)
	return nil
}

func resourceSentryTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	params := &sentry.UpdateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	tflog.Debug(ctx, "Updating Sentry team", map[string]interface{}{
		"teamSlug": slug,
		"org":      org,
	})
	team, _, err := client.Teams.Update(ctx, org, slug, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry team", map[string]interface{}{
		"teamSlug": team.Slug,
		"teamID":   team.ID,
		"org":      org,
	})

	d.SetId(team.Slug)
	return resourceSentryTeamRead(ctx, d, meta)
}

func resourceSentryTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Deleting Sentry team", map[string]interface{}{
		"teamSlug": slug,
		"org":      org,
	})
	_, err := client.Teams.Delete(ctx, org, slug)
	tflog.Debug(ctx, "Deleted Sentry team", map[string]interface{}{
		"teamSlug": slug,
		"org":      org,
	})

	return diag.FromErr(err)
}
