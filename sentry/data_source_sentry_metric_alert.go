package sentry

import (
	"context"
	"fmt"
	"reflect"

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
			"projects": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"triggers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
							},
						},
						"alert_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
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
	d.Set("organization", org)
	d.Set("project", project)
	d.Set("internal_id", alertID)
	d.Set("name", alert.Name)
	d.Set("environment", alert.Environment)
	d.Set("dataset", alert.DataSet)
	d.Set("query", alert.Query)
	d.Set("aggregate", alert.Aggregate)
	d.Set("time_window", alert.TimeWindow)
	d.Set("threshold_type", alert.ThresholdType)
	d.Set("resolve_threshold", alert.ResolveThreshold)
	d.Set("projects", alert.Projects)
	d.Set("owner", alert.Owner)
	d.Set("triggers", mapTriggers(ctx, alert.Triggers))
	return nil
}

func mapTriggers(ctx context.Context, triggers []*sentry.MetricAlertTrigger) []interface{} {
	if triggers != nil {
		trs := make([]interface{}, 0, len(triggers))

		for _, trigger := range triggers {
			tflog.Debug(ctx, "Reading trigger", *trigger)
			tr := make(map[string]interface{})

			tr["id"] = (*trigger)["id"]
			tr["alert_rule_id"] = (*trigger)["alertRuleId"]
			tr["label"] = (*trigger)["label"]
			tr["threshold_type"] = (*trigger)["thresholdType"]
			tr["alert_threshold"] = (*trigger)["alertThreshold"]
			tr["resolve_threshold"] = (*trigger)["resolveThreshold"]
			tr["actions"] = mapActions(ctx, (*trigger)["actions"])

			trs = append(trs, tr)
		}

		return trs
	}

	return make([]interface{}, 0)
}

func mapActions(ctx context.Context, a interface{}) interface{} {
	//convert actions which appears as interface{} but is actually []interface{}
	var actions []map[string]interface{}
	rv := reflect.ValueOf(a)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			t := rv.Index(i).Interface().(map[string]interface{})
			actions = append(actions, t)
		}
	}

	for _, f := range actions {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	return actions
}
