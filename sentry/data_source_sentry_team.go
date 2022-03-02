package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func dataSourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	slug := d.Get("slug").(string)

	tflog.Debug(ctx, "Reading Sentry Team", "org", org, "teamSlug", slug)
	team, _, err := client.Teams.Get(org, slug)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry Team", "teamName", team.Name, "teamSlug", team.Slug, "teamID", team.ID)

	d.SetId(team.Slug)
	d.Set("slug", team.Slug)
	d.Set("name", team.Name)
	d.Set("team_id", team.ID)

	return nil
}
