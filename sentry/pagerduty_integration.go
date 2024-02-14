package sentry

import (
	"context"
	sentry "github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func dataSourcePagerdutyIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerdutyIntegrationRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"pagerduty_integration": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerdutyIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	integrationId := d.Get("integration_id").(int)

	pagerDutyIntegration, _, err := client.Pagerduty.Get(ctx, org, integrationId)
	if err != nil {
		return diag.FromErr(err)
	}

	var pagerdutyIntegrationMap = make(map[string]interface{})
	for _, svc := range pagerDutyIntegration.ConfigData.ServiceTable {
		pagerdutyIntegrationMap[svc.Service] = strconv.Itoa(svc.Id)
	}

	d.SetId(strconv.Itoa(integrationId))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("integration_id", integrationId),
		d.Set("pagerduty_integration", pagerdutyIntegrationMap),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
