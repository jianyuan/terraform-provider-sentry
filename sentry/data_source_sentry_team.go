package sentry

import (
	"context"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Team data source.",

		ReadContext: dataSourceSentryTeamRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the team should be created for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The unique URL slug for this team.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "The internal ID for this team.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The human readable name for this organization.",
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
		},
	}
}

func dataSourceSentryTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	teamSlug := d.Get("slug").(string)

	tflog.Debug(ctx, "Reading team", map[string]interface{}{
		"org":  org,
		"team": teamSlug,
	})
	team, _, err := client.Teams.Get(ctx, org, teamSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sentry.StringValue(team.Slug))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("slug", team.Slug),
		d.Set("internal_id", team.ID),
		d.Set("name", team.Name),
		d.Set("has_access", team.HasAccess),
		d.Set("is_pending", team.IsPending),
		d.Set("is_member", team.IsMember),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
