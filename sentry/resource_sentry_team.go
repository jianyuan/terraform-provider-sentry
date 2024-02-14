package sentry

import (
	"context"
	"net/http"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"internal_id": {
				Description: "The internal ID for this team.",
				Type:        schema.TypeString,
				Computed:    true,
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
			"team_id": {
				Deprecated:  "Use `internal_id` instead.",
				Description: "Use `internal_id` instead.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	params := &sentry.CreateTeamParams{
		Name: sentry.String(d.Get("name").(string)),
	}
	if slug, ok := d.GetOk("slug"); ok {
		params.Slug = sentry.String(slug.(string))
	}

	tflog.Debug(ctx, "Creating team", map[string]interface{}{"org": org, "teamName": params.Name})
	team, _, err := client.Teams.Create(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sentry.StringValue(team.Slug))
	return resourceSentryTeamRead(ctx, d, meta)
}

func resourceSentryTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	teamSlug := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Reading team", map[string]interface{}{"org": org, "team": teamSlug})
	team, _, err := client.Teams.Get(ctx, org, teamSlug)
	if err != nil {
		if sErr, ok := err.(*sentry.ErrorResponse); ok {
			if sErr.Response.StatusCode == http.StatusNotFound {
				tflog.Info(ctx, "Removing team from state because it no longer exists in Sentry", map[string]interface{}{"team": teamSlug})
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("name", team.Name),
		d.Set("slug", team.Slug),
		d.Set("internal_id", team.ID),
		d.Set("has_access", team.HasAccess),
		d.Set("is_pending", team.IsPending),
		d.Set("is_member", team.IsMember),
		d.Set("team_id", team.ID), // Deprecated
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	teamSlug := d.Id()
	org := d.Get("organization").(string)
	params := &sentry.UpdateTeamParams{
		Name: sentry.String(d.Get("name").(string)),
		Slug: sentry.String(d.Get("slug").(string)),
	}

	tflog.Debug(ctx, "Updating team", map[string]interface{}{"org": org, "team": teamSlug})
	team, _, err := client.Teams.Update(ctx, org, teamSlug, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sentry.StringValue(team.Slug))
	return resourceSentryTeamRead(ctx, d, meta)
}

func resourceSentryTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	teamSlug := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Deleting team", map[string]interface{}{"org": org, "team": teamSlug})
	_, err := client.Teams.Delete(ctx, org, teamSlug)
	return diag.FromErr(err)
}
