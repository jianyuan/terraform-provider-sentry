package sentry

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/javaadsnappcar/go-sentry/sentry"
	"github.com/mitchellh/mapstructure"
)

func resourceSentryAPMRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryAPMRuleCreate,
		ReadContext:   resourceSentryAPMRuleRead,
		UpdateContext: resourceSentryAPMRuleUpdate,
		DeleteContext: resourceSentryAPMRuleDelete,

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
				Description: "The APM rule name",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Perform APM rule in a specific environment",
			},
			"dataset": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Sentry APM category",
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
				Description: "The period to evaluate the APM rule in minutes",
			},
			"threshold_type": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The type of threshold",
			},
			"resolve_threshold": &schema.Schema{
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The value at which the APM rule resolves",
			},
			// "triggers": {
			// 	Type:     schema.TypeList,
			// 	Required: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
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
							// Default: [],
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
				Description: "Specifies the owner id of this APM rule",
			},
		},
	}
}

func resourceSentryAPMRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	tflog.Info(ctx, "trying to set projects", d.Get("projects"))
	projects_input := d.Get("projects").([]interface{})
	projects := make([]string, len(projects_input))
	for i, v := range projects_input {
		projects[i] = fmt.Sprint(v)
	}

	//using Type.SchemaSet
	inputTriggers := d.Get("triggers").(*schema.Set)
	triggers := mapTriggersCreate(inputTriggers)
	// inputTriggers := d.Get("triggers").(*schema.Set)
	// inputTriggersList := inputTriggers.List()
	// triggers := make([]sentry.Trigger, len(inputTriggersList))
	// for i, ia := range inputTriggersList {
	// 	var trigger sentry.Trigger
	// 	mapstructure.WeakDecode(ia, &trigger)

	// 	//replace with uppercasing
	// 	trigger["alertThreshold"] = trigger["alert_threshold"]
	// 	trigger["resolveThreshold"] = trigger["resolve_threshold"]
	// 	trigger["thresholdType"] = trigger["threshold_type"]
	// 	delete(trigger, "alert_threshold")
	// 	delete(trigger, "resolve_threshold")
	// 	delete(trigger, "threshold_type")

	// 	//test delete alert and id
	// 	delete(trigger, "alert_rule_id")
	// 	delete(trigger, "id")

	// 	triggers[i] = trigger
	// }
	// //swop trigger elements so critical is first
	// if triggers[0]["label"] != "critical" {
	// 	var criticalTriggerIndex int
	// 	for i, trigger := range triggers {
	// 		if trigger["label"] == "critical" {
	// 			criticalTriggerIndex = i
	// 		}
	// 	}

	// 	temp := triggers[criticalTriggerIndex]
	// 	triggers[criticalTriggerIndex] = triggers[0]
	// 	triggers[0] = temp
	// }

	tflog.Info(ctx, "triggers", triggers)

	//using Type.SchemaList (doesn't work with list because of actions needing to be)
	// inputTriggers := d.Get("triggers").([]interface{})
	// // inputTriggersList := inputTriggers.List()
	// triggers := make([]sentry.Trigger, len(inputTriggers))
	// for i, ia := range inputTriggers {
	// 	var trigger sentry.Trigger
	// 	mapstructure.WeakDecode(ia, &trigger)
	// 	triggers[i] = trigger
	// }
	// tflog.Info(ctx, "triggers", triggers)

	//create hardcoded triggers (works)
	// triggers := make([]sentry.Trigger, 0)
	// triggers = append(triggers, sentry.Trigger{
	// 	// "actions":       [],
	// 	"alertThreshold":   10000,
	// 	"label":            "critical",
	// 	"resolveThreshold": 100.0,
	// 	"thresholdType":    0,
	// })

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

	tflog.Info(ctx, "Creating Sentry APM rule", "ruleName", name, "org", org, "project", project, "params", params)
	apmRule, _, err := client.APMRules.Create(org, project, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "Created Sentry APM rule", "ruleName", apmRule.Name, "ruleID", apmRule.ID, "org", org, "project", project)

	d.SetId(apmRule.ID)

	return resourceSentryAPMRuleRead(ctx, d, meta)
}

func resourceSentryAPMRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	id := d.Id()

	tflog.Debug(ctx, "Reading Sentry APM rule", "apmRuleID", id, "org", org, "project", project)
	apmRules, resp, err := client.APMRules.List(org, project)

	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "Read Sentry APM rules", "ruleCount", len(apmRules), "APM rules", apmRules)

	var apmRule *sentry.APMRule
	for _, r := range apmRules {
		if r.ID == id {
			apmRule = &r
			break
		}
	}

	if apmRule == nil {
		return diag.Errorf("Could not find apmRule with ID" + id)
	}
	tflog.Debug(ctx, "Read Sentry APM rule", "ruleID", apmRule.ID, "org", org, "project", project)

	triggers := mapTriggers(ctx, &apmRule.Triggers)
	if err := d.Set("triggers", triggers); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(apmRule.ID)
	d.Set("name", apmRule.Name)
	// d.Set("organization", apmRule.Owner
	// d.Set("project").(string)
	d.Set("environment", apmRule.Environment)
	d.Set("dataset", apmRule.DataSet)
	d.Set("query", apmRule.Query)
	d.Set("aggregate", apmRule.Aggregate)
	d.Set("time_window", apmRule.TimeWindow)
	d.Set("threshold_type", apmRule.ThresholdType)
	d.Set("resolve_threshold", apmRule.ResolveThreshold)
	// s := mapSchemaTriggers(ctx, d.Get("triggers").([]map[string]*schema.Schema))
	tflog.Debug(ctx, "trying to set projects", apmRule.Projects)
	d.Set("projects", apmRule.Projects)
	tflog.Debug(ctx, "succeeded to set projects", apmRule.Projects)
	d.Set("owner", apmRule.Owner)

	return nil
}

func resourceSentryAPMRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	tflog.Info(ctx, "trying to set projects", d.Get("projects"))
	projects_input := d.Get("projects").([]interface{})
	projects := make([]string, len(projects_input))
	for i, v := range projects_input {
		projects[i] = fmt.Sprint(v)
	}

	//using Type.SchemaSet
	inputTriggers := d.Get("triggers").(*schema.Set)
	triggers := mapTriggersCreate(inputTriggers)
	// inputTriggersList := inputTriggers.List()
	// triggers := make([]sentry.Trigger, len(inputTriggersList))
	// for i, ia := range inputTriggersList {
	// 	var trigger sentry.Trigger
	// 	mapstructure.WeakDecode(ia, &trigger)

	// 	//replace with uppercasing
	// 	trigger["alertThreshold"] = trigger["alert_threshold"]
	// 	trigger["resolveThreshold"] = trigger["resolve_threshold"]
	// 	trigger["thresholdType"] = trigger["threshold_type"]
	// 	delete(trigger, "alert_threshold")
	// 	delete(trigger, "resolve_threshold")
	// 	delete(trigger, "threshold_type")

	// 	//test delete alert and id
	// 	delete(trigger, "alert_rule_id")
	// 	delete(trigger, "id")

	// 	triggers[i] = trigger
	// }

	// //swop trigger elements so critical is first
	// if triggers[0]["label"] != "critical" {
	// 	var criticalTriggerIndex int
	// 	for i, trigger := range triggers {
	// 		if trigger["label"] == "critical" {
	// 			criticalTriggerIndex = i
	// 		}
	// 	}

	// 	temp := triggers[criticalTriggerIndex]
	// 	triggers[criticalTriggerIndex] = triggers[0]
	// 	triggers[0] = temp
	// }

	tflog.Info(ctx, "triggers", triggers)

	params := &sentry.APMRule{
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

	tflog.Debug(ctx, "Updating Sentry APM rule", "ruleName", name, "org", org, "project", project)
	apmRule, _, err := client.APMRules.Update(org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry APM rule", "ruleName", apmRule.Name, "ruleID", apmRule.ID, "org", org, "project", project)

	return resourceSentryAPMRuleRead(ctx, d, meta)
}

func resourceSentryAPMRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Deleting Sentry APM rule", "ruleID", id, "org", org, "project", project)
	_, err := client.APMRules.Delete(org, project, id)
	tflog.Debug(ctx, "Deleted Sentry APM rule", "ruleID", id, "org", org, "project", project)

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

		//test delete alert and id
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

// func mapSchemaTriggers(ctx context.Context, triggers []map[string]*schema.Schema) []sentry.Trigger {
// 	tflog.Debug(ctx, "Mapping triggers")
// 	if triggers != nil {
// 		tflog.Debug(ctx, "Triggers found", triggers)
// 		trs := make([]sentry.Trigger, len(triggers), len(triggers))

// 		for i, trigger := range triggers {
// 			tflog.Debug(ctx, "Reading trigger", trigger)
// 			tr := make(map[string]interface{})

// 			tr["id"] = trigger["id"]
// 			tr["alert_rule_id"] = trigger["alertRuleId"]
// 			tr["label"] = trigger["label"]
// 			tr["threshold_type"] = trigger["thresholdType"]
// 			tr["alert_threshold"] = trigger["alertThreshold"]
// 			tr["resolve_threshold"] = trigger["resolveThreshold"]
// 			// tr["actions"] = mapActions(ctx, trigger["actions"]) //TODO: map later

// 			trs[i] = tr
// 		}

// 		tflog.Debug(ctx, "Mapped triggers", trs)
// 		return trs
// 	}

// 	tflog.Debug(ctx, "No triggers found", triggers)
// 	return make([]sentry.Trigger, 0)
// }
