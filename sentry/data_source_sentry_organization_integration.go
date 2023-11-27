package sentry

import (
	"context"

	"github.com/deste-org/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSentryOrganizationIntegration() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Organization Integration data source.",

		ReadContext: dataSourceSentryOrganizationIntegrationRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the integration belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provider_key": {
				Description: "The key of the organization integration provider.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the organization integration.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "The internal ID for this organization integration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSentryOrganizationIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	providerKey := d.Get("provider_key").(string)
	integrationName := d.Get("name").(string)

	tflog.Debug(ctx, "Reading organization integration", map[string]interface{}{"org": org, "provider_key": providerKey, "name": integrationName})

	// get all paginated integrations with the provider key
	var orgIntegrations []*sentry.OrganizationIntegration
	params := &sentry.ListOrganizationIntegrationsParams{
		ListCursorParams: sentry.ListCursorParams{},
		ProviderKey:      providerKey,
	}
	for {
		keys, resp, err := client.OrganizationIntegrations.List(ctx, org, params)
		if err != nil {
			return diag.FromErr(err)
		}
		orgIntegrations = append(orgIntegrations, keys...)

		tflog.Debug(ctx, "Requested organization integration list cursor", map[string]interface{}{"cursor": resp.Cursor})
		if resp.Cursor == "" {
			break
		}
		params.ListCursorParams.Cursor = resp.Cursor
	}

	// filter for first matching name
	for _, orgIntegration := range orgIntegrations {
		if orgIntegration.Name == integrationName {
			d.SetId(orgIntegration.ID)
			retErr := multierror.Append(
				d.Set("internal_id", orgIntegration.ID),
			)
			return diag.FromErr(retErr.ErrorOrNil())
		}
	}

	return diag.Errorf("Can't find Sentry Organization Integration: %s", integrationName)
}
