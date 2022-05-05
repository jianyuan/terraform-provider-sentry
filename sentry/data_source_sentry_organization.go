package sentry

import (
	"context"

	"github.com/canva/go-sentry/sentry"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSentryOrganizationRead,
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},

			"internal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSentryOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Get("slug").(string)

	tflog.Debug(ctx, "Reading Sentry org", map[string]interface{}{"orgSlug": slug})
	org, _, err := client.Organizations.Get(slug)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry org", map[string]interface{}{"orgName": org.Name, "orgSlug": org.Slug, "orgID": org.ID})

	d.SetId(org.Slug)
	d.Set("internal_id", org.ID)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)

	return nil
}
