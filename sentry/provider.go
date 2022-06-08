package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Description: "The authentication token used to connect to Sentry. The value can be sourced from " +
					"the `SENTRY_AUTH_TOKEN` environment variable.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SENTRY_AUTH_TOKEN", "SENTRY_TOKEN"}, nil),
				Sensitive:   true,
			},
			"base_url": {
				Description: "The target Sentry Base API URL in the format `https://[hostname]/api/`. " +
					"The default value is `https://sentry.io/api/`. The value must be provided when working with " +
					"Sentry On-Premise. The value can be sourced from the `SENTRY_BASE_URL` environment variable.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SENTRY_BASE_URL", "https://sentry.io/api/"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"sentry_alert_rule":   resourceSentryAlertRule(),
			"sentry_issue_alert":  resourceSentryIssueAlert(),
			"sentry_default_key":  resourceSentryDefaultKey(),
			"sentry_key":          resourceSentryKey(),
			"sentry_organization": resourceSentryOrganization(),
			"sentry_plugin":       resourceSentryPlugin(),
			"sentry_project":      resourceSentryProject(),
			"sentry_rule":         resourceSentryRule(),
			"sentry_team":         resourceSentryTeam(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sentry_key":          dataSourceSentryKey(),
			"sentry_organization": dataSourceSentryOrganization(),
			"sentry_alert_rules":  dataSourceSentryAlertRules(),
		},

		ConfigureContextFunc: providerContextConfigure,
	}
}

func providerContextConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Token:   d.Get("token").(string),
		BaseURL: d.Get("base_url").(string),
	}
	return config.Client(ctx)
}
