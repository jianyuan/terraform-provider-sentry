package sentry

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func dataSourceSentryAPMRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSentryAPMRuleRead,
		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the project belongs to",
			},
			"project": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project to create the plugin for",
			},
			"apm_rules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The APM rule name",
						},
						"environment": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataset": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregate": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_window": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"threshold_type": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resolve_threshold": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"projects": &schema.Schema{
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"owner": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"triggers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
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
									"alert_rule_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"label": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"threshold_type": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"alert_threshold": &schema.Schema{
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"resolve_threshold": &schema.Schema{
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

func dataSourceSentryAPMRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry APM rules", "org", org, "project", project)
	apmRules, resp, err := client.APMRules.List(org, project)
	if err != nil {
		return diag.FromErr(err)
	}
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "Read Sentry APM rules", "ruleCount", len(apmRules), "APM rules", apmRules)

	apm_rules := mapApmRulesData(ctx, &apmRules)
	if err := d.Set("apm_rules", apm_rules); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(org)

	return nil
}

func mapApmRulesData(ctx context.Context, apmRules *[]sentry.APMRule) []interface{} {
	if apmRules != nil {
		ars := make([]interface{}, len(*apmRules), len(*apmRules))

		for i, apmRule := range *apmRules {
			ar := make(map[string]interface{})

			ar["id"] = apmRule.ID
			ar["name"] = apmRule.Name
			ar["environment"] = apmRule.Environment
			ar["dataset"] = apmRule.DataSet
			ar["query"] = apmRule.Query
			ar["aggregate"] = apmRule.Aggregate
			ar["time_window"] = apmRule.TimeWindow
			ar["threshold_type"] = apmRule.ThresholdType
			ar["resolve_threshold"] = apmRule.ResolveThreshold
			ar["projects"] = apmRule.Projects
			ar["owner"] = apmRule.Owner
			ar["triggers"] = mapTriggers(ctx, &apmRule.Triggers)
			// ar["created"] = apmRule.Created //TODO: map later

			ars[i] = ar
		}

		return ars
	}

	return make([]interface{}, 0)
}

func mapTriggers(ctx context.Context, triggers *[]sentry.Trigger) []interface{} {
	if triggers != nil {
		trs := make([]interface{}, len(*triggers), len(*triggers))

		for i, trigger := range *triggers {
			tflog.Debug(ctx, "Reading trigger", trigger)
			tr := make(map[string]interface{})

			tr["id"] = trigger["id"]
			tr["alert_rule_id"] = trigger["alertRuleId"]
			tr["label"] = trigger["label"]
			tr["threshold_type"] = trigger["thresholdType"]
			tr["alert_threshold"] = trigger["alertThreshold"]
			tr["resolve_threshold"] = trigger["resolveThreshold"]
			tr["actions"] = mapActions(ctx, trigger["actions"])

			trs[i] = tr
		}

		return trs
	}

	return make([]interface{}, 0)
}

func mapActions(ctx context.Context, a interface{}) interface{} {
	//convert actions which appears as interface{} but is actually []interface{}
	// var actions []interface{}
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
