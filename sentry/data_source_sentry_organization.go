package sentry

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
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

	logging.Debugf("Received org named %s with ID: %s and internal_id: %s", org.Name, org.Slug, org.ID)
	d.SetId(org.Slug)
	d.Set("internal_id", org.ID)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)

	return nil
}
