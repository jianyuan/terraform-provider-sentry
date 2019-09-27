package sentry

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SENTRY_TOKEN", nil),
				Description: "The authentication token used to connect to Sentry",
			},
			"base_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SENTRY_BASE_URL", "https://app.getsentry.com/api/"),
				Description: "The Sentry Base API URL",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"sentry_organization": resourceSentryOrganization(),
			"sentry_team":         resourceSentryTeam(),
			"sentry_project":      resourceSentryProject(),
			"sentry_key":          resourceSentryKey(),
			"sentry_plugin":       resourceSentryPlugin(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sentry_key": dataSourceSentryKey(),
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
