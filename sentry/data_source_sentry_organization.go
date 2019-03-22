package sentry

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func dataSourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSentryOrganizationRead,
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

func dataSourceSentryOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Get("slug").(string)

	org, _, err := client.Organizations.Get(slug)
	if err != nil {
		return err
	}

	d.SetId(org.Slug)
	d.Set("internal_id", org.ID)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)

	return nil
}
