package sentry

import (
	"context"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Organization data source.",

		ReadContext: dataSourceSentryOrganizationRead,

		Schema: map[string]*schema.Schema{
			"slug": {
				Description: "The unique URL slug for this organization.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "The internal ID for this organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The human readable name for this organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSentryOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("slug").(string)

	tflog.Debug(ctx, "Reading organization", map[string]interface{}{"org": org})
	organization, _, err := client.Organizations.Get(ctx, org)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sentry.StringValue(organization.Slug))
	retErr := multierror.Append(
		d.Set("name", organization.Name),
		d.Set("internal_id", organization.ID),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
