package sentry

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SENTRY_TOKEN", nil),
				Description: "The authentication token used to connect to Sentry",
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SENTRY_BASE_URL", "https://sentry.io/api/"),
				Description: "The Sentry Base API URL",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"sentry_organization": resourceSentryOrganization(),
			"sentry_team":         resourceSentryTeam(),
			"sentry_project":      resourceSentryProject(),
			"sentry_key":          resourceSentryKey(),
			"sentry_plugin":       resourceSentryPlugin(),
			"sentry_rule":         resourceSentryRule(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sentry_key":          dataSourceSentryKey(),
			"sentry_organization": dataSourceSentryOrganization(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token:   d.Get("token").(string),
		BaseURL: d.Get("base_url").(string),
	}
	return config.Client()
}
