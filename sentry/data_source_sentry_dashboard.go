package sentry

import (
	"context"

	"github.com/deste-org/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSentryDashboard() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSentryDashboardRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the dashboard belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "The internal ID for this dashboard.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"title": {
				Description: "Dashboard title.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"widget": {
				Description: "Dashboard widgets.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fields": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"aggregates": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"columns": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"field_aliases": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"conditions": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"order_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"widget_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"x": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"y": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"w": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"h": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_h": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSentryDashboardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	dashboardID := d.Get("internal_id").(string)

	tflog.Debug(ctx, "Reading dashboard", map[string]interface{}{
		"org":         org,
		"dashboardID": dashboardID,
	})
	dashboard, _, err := client.Dashboards.Get(ctx, org, dashboardID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(org, sentry.StringValue(dashboard.ID)))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("internal_id", dashboard.ID),
		d.Set("title", dashboard.Title),
		d.Set("widget", flattenDashboardWidgets(dashboard.Widgets)),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
