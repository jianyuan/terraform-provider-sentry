package sentry

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func dataSourceSentryProject() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Project data source.",

		ReadContext: dataSourceSentryProjectRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The unique URL slug for this project.",
				Type:        schema.TypeString,
				Required:    true,
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
		},
	}
}

func dataSourceSentryProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	projectSlug := d.Get("slug").(string)

	tflog.Debug(ctx, "Reading project", map[string]interface{}{
		"org":     org,
		"project": projectSlug,
	})
	project, _, err := client.Projects.Get(ctx, org, projectSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.Slug)
	retErr := multierror.Append(
		d.Set("organization", project.Organization.Slug),
		d.Set("slug", project.Slug),
		d.Set("internal_id", project.ID),
		d.Set("is_public", project.IsPublic),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
