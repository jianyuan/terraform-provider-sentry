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
		ReadContext: dataSourceSentryAlertRuleRead,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the project belongs to",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project to create the plugin for",
			},
			"alert_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alert rule name",
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
				},
			},
		},
	}
}

func dataSourceSentryAlertRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry Alert rules", map[string]interface{}{
		"org":     org,
		"project": project,
	})
	alertRules, resp, err := client.MetricAlerts.List(ctx, org, project)
	if err != nil {
		return diag.FromErr(err)
	}
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "Read Sentry Alert rules", map[string]interface{}{
		"ruleCount":  len(alertRules),
		"alertRules": alertRules,
	})

	alert_rules := mapAlertRulesData(ctx, alertRules)
	if err := d.Set("alert_rules", alert_rules); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(org)

	return nil
}

func mapAlertRulesData(ctx context.Context, alertRules []*sentry.MetricAlert) []interface{} {
	if alertRules != nil {
		ars := make([]interface{}, len(alertRules), len(alertRules))

		for i, alertRule := range alertRules {
			ar := make(map[string]interface{})

			ar["id"] = alertRule.ID
			ar["name"] = alertRule.Name
			ar["environment"] = alertRule.Environment
			ar["dataset"] = alertRule.DataSet
			ar["query"] = alertRule.Query
			ar["aggregate"] = alertRule.Aggregate
			ar["time_window"] = alertRule.TimeWindow
			ar["threshold_type"] = alertRule.ThresholdType
			ar["resolve_threshold"] = alertRule.ResolveThreshold
			ar["projects"] = alertRule.Projects
			ar["owner"] = alertRule.Owner
			ar["triggers"] = mapTriggers(ctx, alertRule.Triggers)

			ars[i] = ar
		}

		return ars
	}

	return make([]interface{}, 0)
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
