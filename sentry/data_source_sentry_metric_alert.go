package sentry

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func dataSourceSentryMetricAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSentryMetricAlertRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the metric alert belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The slug of the project the metric alert belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "The internal ID for this metric alert.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The metric alert name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_types": {
				Description: "The events type of dataset.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aggregate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_window": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"threshold_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resolve_threshold": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_identifier": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"integration_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"input_channel_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threshold_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"alert_threshold": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"resolve_threshold": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSentryMetricAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	alertID := d.Get("internal_id").(string)

	tflog.Debug(ctx, "Reading metric alert", map[string]interface{}{"org": org, "project": project, "alertID": alertID})
	alert, _, err := client.MetricAlerts.Get(ctx, org, project, alertID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(org, project, sentry.StringValue(alert.ID)))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("project", project),
		d.Set("internal_id", alertID),
		d.Set("name", alert.Name),
		d.Set("environment", alert.Environment),
		d.Set("dataset", alert.DataSet),
		d.Set("event_types", alert.EventTypes),
		d.Set("query", alert.Query),
		d.Set("aggregate", alert.Aggregate),
		d.Set("time_window", alert.TimeWindow),
		d.Set("threshold_type", alert.ThresholdType),
		d.Set("resolve_threshold", alert.ResolveThreshold),
		d.Set("owner", alert.Owner),
		d.Set("trigger", flattenMetricAlertTriggers(alert.Triggers)),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
