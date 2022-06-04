package sentry

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/mitchellh/mapstructure"
)

func resourceSentryAlertRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryAlertRuleCreate,
		ReadContext:   resourceSentryAlertRuleRead,
		UpdateContext: resourceSentryAlertRuleUpdate,
		DeleteContext: resourceSentryAlertRuleDelete,

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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Alert rule name",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Perform Alert rule in a specific environment",
			},
			"dataset": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Sentry Alert category",
			},
			"query": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The query filter to apply",
			},
			"aggregate": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The aggregation criteria to apply",
			},
			"time_window": &schema.Schema{
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "The period to evaluate the Alert rule in minutes",
			},
			"threshold_type": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The type of threshold",
			},
			"resolve_threshold": &schema.Schema{
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The value at which the Alert rule resolves",
			},
			"triggers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
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
							Required: true,
						},
						"threshold_type": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"alert_threshold": &schema.Schema{
							Type:     schema.TypeFloat,
							Required: true,
						},
						"resolve_threshold": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"projects": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the owner id of this Alert rule",
			},
		},
	}
}

func resourceSentryAlertRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	name := d.Get("name").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	environment := d.Get("environment").(string)
	dataset := d.Get("dataset").(string)
	query := d.Get("query").(string)
	aggregate := d.Get("aggregate").(string)
	time_window := d.Get("time_window").(float64)
	threshold_type := d.Get("threshold_type").(int)
	resolve_threshold := d.Get("resolve_threshold").(float64)
	owner := d.Get("owner").(string)

	inputProjects := d.Get("projects").([]interface{})
	projects := make([]string, len(inputProjects))
	for i, v := range inputProjects {
		projects[i] = fmt.Sprint(v)
	}

	inputTriggers := d.Get("triggers").(*schema.Set)
	triggers := mapTriggersCreate(inputTriggers)

	params := &sentry.CreateAPMRuleParams{
		Name:             name,
		DataSet:          dataset,
		Query:            query,
		Aggregate:        aggregate,
		TimeWindow:       time_window,
		ThresholdType:    threshold_type,
		ResolveThreshold: resolve_threshold,
		Triggers:         triggers,
		Projects:         projects,
		Owner:            owner,
	}

	if environment != "" {
		params.Environment = &environment
	}

	tflog.Info(ctx, "Creating Sentry Alert rule", map[string]interface{}{
		"ruleName": name,
		"org":      org,
		"project":  project,
		"params":   params,
	})
	alertRule, _, err := client.APMRules.Create(org, project, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "Created Sentry Alert rule", map[string]interface{}{
		"ruleName": alertRule.Name,
		"ruleID":   alertRule.ID,
		"org":      org,
		"project":  project,
	})

	d.SetId(alertRule.ID)

	return resourceSentryAlertRuleRead(ctx, d, meta)
}

func resourceSentryAlertRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	id := d.Id()

	tflog.Debug(ctx, "Reading Sentry Alert rule", map[string]interface{}{
		"alertRuleID": id,
		"org":         org,
		"project":     project,
	})
	alertRules, resp, err := client.APMRules.List(org, project)

	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "Read Sentry Alert rules", map[string]interface{}{
		"ruleCount":  len(alertRules),
		"alertRules": alertRules,
	})

	var alertRule *sentry.APMRule
	for _, r := range alertRules {
		if r.ID == id {
			alertRule = &r
			break
		}
	}

	if alertRule == nil {
		return diag.Errorf("Could not find alertRule with ID" + id)
	}
	tflog.Debug(ctx, "Read Sentry Alert rule", map[string]interface{}{
		"ruleID":  alertRule.ID,
		"org":     org,
		"project": project,
	})

	triggers := mapResourceTriggersRead(ctx, &alertRule.Triggers)
	if err := d.Set("triggers", triggers); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(alertRule.ID)
	d.Set("name", alertRule.Name)
	d.Set("environment", alertRule.Environment)
	d.Set("dataset", alertRule.DataSet)
	d.Set("query", alertRule.Query)
	d.Set("aggregate", alertRule.Aggregate)
	d.Set("time_window", alertRule.TimeWindow)
	d.Set("threshold_type", alertRule.ThresholdType)
	d.Set("resolve_threshold", alertRule.ResolveThreshold)
	d.Set("projects", alertRule.Projects)
	d.Set("owner", alertRule.Owner)

	return nil
}

func resourceSentryAlertRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	name := d.Get("name").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	environment := d.Get("environment").(string)
	dataset := d.Get("dataset").(string)
	query := d.Get("query").(string)
	aggregate := d.Get("aggregate").(string)
	time_window := d.Get("time_window").(float64)
	threshold_type := d.Get("threshold_type").(int)
	resolve_threshold := d.Get("resolve_threshold").(float64)
	owner := d.Get("owner").(string)

	inputProjects := d.Get("projects").([]interface{})
	projects := make([]string, len(inputProjects))
	for i, v := range inputProjects {
		projects[i] = fmt.Sprint(v)
	}

	inputTriggers := d.Get("triggers").(*schema.Set)
	triggers := mapTriggersCreate(inputTriggers)

	params := &sentry.APMRule{
		ID:               id,
		Name:             name,
		Environment:      &environment,
		DataSet:          dataset,
		Query:            query,
		Aggregate:        aggregate,
		TimeWindow:       time_window,
		ThresholdType:    threshold_type,
		ResolveThreshold: resolve_threshold,
		Triggers:         triggers,
		Projects:         projects,
		Owner:            owner,
	}

	tflog.Debug(ctx, "Updating Sentry Alert rule", map[string]interface{}{
		"ruleName": name,
		"org":      org,
		"project":  project,
	})
	alertRule, _, err := client.APMRules.Update(org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry Alert rule", map[string]interface{}{
		"ruleName": alertRule.Name,
		"ruleID":   alertRule.ID,
		"org":      org,
		"project":  project,
	})

	return resourceSentryAlertRuleRead(ctx, d, meta)
}

func resourceSentryAlertRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Deleting Sentry Alert rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})
	_, err := client.APMRules.Delete(org, project, id)
	tflog.Debug(ctx, "Deleted Sentry Alert rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})

	return diag.FromErr(err)
}

func mapTriggersCreate(inputTriggers *schema.Set) []sentry.Trigger {
	inputTriggersList := inputTriggers.List()
	triggers := make([]sentry.Trigger, len(inputTriggersList))
	for i, ia := range inputTriggersList {
		var trigger sentry.Trigger
		mapstructure.WeakDecode(ia, &trigger)

		//replace with uppercasing
		trigger["alertThreshold"] = trigger["alert_threshold"]
		trigger["resolveThreshold"] = trigger["resolve_threshold"]
		trigger["thresholdType"] = trigger["threshold_type"]
		delete(trigger, "alert_threshold")
		delete(trigger, "resolve_threshold")
		delete(trigger, "threshold_type")

		//delete id and alert_rule_id as they are not required in POST&PUT requests
		delete(trigger, "alert_rule_id")
		delete(trigger, "id")

		triggers[i] = trigger
	}

	//swop trigger elements so critical is first
	if triggers[0]["label"] != "critical" {
		var criticalTriggerIndex int
		for i, trigger := range triggers {
			if trigger["label"] == "critical" {
				criticalTriggerIndex = i
			}
		}

		temp := triggers[criticalTriggerIndex]
		triggers[criticalTriggerIndex] = triggers[0]
		triggers[0] = temp
	}

	return triggers
}

func mapResourceTriggersRead(ctx context.Context, triggers *[]sentry.Trigger) []interface{} {
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
			tr["actions"] = mapResourceActionsRead(ctx, trigger["actions"])

			trs[i] = tr
		}

		return trs
	}

	return make([]interface{}, 0)
}

func mapResourceActionsRead(ctx context.Context, a interface{}) interface{} {
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
