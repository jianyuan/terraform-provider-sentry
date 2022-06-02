package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryTeamCreate,
		ReadContext:   resourceSentryTeamRead,
		UpdateContext: resourceSentryTeamUpdate,
		DeleteContext: resourceSentryTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSentryTeamImport,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the team should be created for",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional slug for this team",
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
	team, resp, err := client.Teams.Create(org, params)
	tflog.Debug(ctx, "Sentry team create http response data", logging.ExtractHttpResponse(resp))
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
	team, resp, err := client.Teams.Get(org, slug)
	tflog.Debug(ctx, "Sentry team read http response data", logging.ExtractHttpResponse(resp))
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
	team, resp, err := client.Teams.Update(org, slug, params)
	tflog.Debug(ctx, "Sentry team update http response data", logging.ExtractHttpResponse(resp))
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
	resp, err := client.Teams.Delete(org, slug)
	tflog.Debug(ctx, "Sentry team delete http response data", logging.ExtractHttpResponse(resp))
	tflog.Debug(ctx, "Deleted Sentry team", map[string]interface{}{
		"teamSlug": slug,
		"org":      org,
	})

	return diag.FromErr(err)
}
