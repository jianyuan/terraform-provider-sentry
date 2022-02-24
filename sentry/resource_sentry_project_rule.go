package sentry

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/mitchellh/mapstructure"
)

const (
	defaultActionMatch = "any"
	defaultFilterMatch = "any"
	defaultFrequency   = 30
)

func resourceSentryRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryRuleCreate,
		ReadContext:   resourceSentryRuleRead,
		UpdateContext: resourceSentryRuleUpdate,
		DeleteContext: resourceSentryRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSentryRuleImporter,
		},

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule name",
			},
			"action_match": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter_match": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Perform actions at most once every X minutes",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Perform rule in a specific environment",
			},
		},
	}
}

func resourceSentryRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	name := d.Get("name").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	environment := d.Get("environment").(string)
	actionMatch := d.Get("action_match").(string)
	filterMatch := d.Get("filter_match").(string)
	inputConditions := d.Get("conditions").([]interface{})
	inputActions := d.Get("actions").([]interface{})
	inputFilters := d.Get("filters").([]interface{})
	frequency := d.Get("frequency").(int)

	if actionMatch == "" {
		actionMatch = defaultActionMatch
	}
	if filterMatch == "" {
		filterMatch = defaultFilterMatch
	}
	if frequency == 0 {
		frequency = defaultFrequency
	}

	conditions := make([]sentry.ConditionType, len(inputConditions))
	for i, ic := range inputConditions {
		var condition sentry.ConditionType
		mapstructure.WeakDecode(ic, &condition)
		conditions[i] = condition
	}
	actions := make([]sentry.ActionType, len(inputActions))
	for i, ia := range inputActions {
		var action sentry.ActionType
		mapstructure.WeakDecode(ia, &action)
		actions[i] = action
	}
	filters := make([]sentry.FilterType, len(inputFilters))
	for i, ia := range inputFilters {
		var filter sentry.FilterType
		mapstructure.WeakDecode(ia, &filter)
		filters[i] = filter
	}

	params := &sentry.CreateRuleParams{
		ActionMatch: actionMatch,
		FilterMatch: filterMatch,
		Environment: environment,
		Frequency:   frequency,
		Name:        name,
		Conditions:  conditions,
		Actions:     actions,
		Filters:     filters,
	}

	if environment != "" {
		params.Environment = environment
	}

	tflog.Debug(ctx, "Creating Sentry rule", "ruleName", name, "org", org, "project", project)
	rule, _, err := client.Rules.Create(org, project, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry rule", "ruleName", rule.Name, "ruleID", rule.ID, "org", org, "project", project)

	d.SetId(rule.ID)

	return resourceSentryRuleRead(ctx, d, meta)
}

func resourceSentryRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	id := d.Id()

	tflog.Debug(ctx, "Reading Sentry rule", "ruleID", id, "org", org, "project", project)
	rules, resp, err := client.Rules.List(org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "Read Sentry rules", "ruleCount", len(rules), "rules", rules)

	var rule *sentry.Rule
	for _, r := range rules {
		if r.ID == id {
			rule = &r
			break
		}
	}

	if rule == nil {
		return diag.Errorf("Could not find rule with ID " + id)
	}
	tflog.Debug(ctx, "Read Sentry rule", "ruleID", rule.ID, "org", org, "project", project)

	// workaround for
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/62
	// as the data sent by Sentry is integer
	for _, f := range rule.Actions {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	for _, f := range rule.Conditions {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	for _, f := range rule.Filters {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	d.SetId(rule.ID)
	d.Set("name", rule.Name)
	d.Set("frequency", rule.Frequency)
	d.Set("environment", rule.Environment)
	d.Set("filters", rule.Filters)
	d.Set("actions", rule.Actions)
	d.Set("conditions", rule.Conditions)
	d.Set("action_match", rule.ActionMatch)
	d.Set("filter_match", rule.FilterMatch)

	return nil
}

func resourceSentryRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	name := d.Get("name").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	environment := d.Get("environment").(string)
	actionMatch := d.Get("action_match").(string)
	filterMatch := d.Get("filter_match").(string)
	inputConditions := d.Get("conditions").([]interface{})
	inputActions := d.Get("actions").([]interface{})
	inputFilters := d.Get("filters").([]interface{})
	frequency := d.Get("frequency").(int)

	if actionMatch == "" {
		actionMatch = defaultActionMatch
	}
	if filterMatch == "" {
		filterMatch = defaultFilterMatch
	}
	if frequency == 0 {
		frequency = defaultFrequency
	}

	conditions := make([]sentry.ConditionType, len(inputConditions))
	for i, ic := range inputConditions {
		var condition sentry.ConditionType
		mapstructure.WeakDecode(ic, &condition)
		conditions[i] = condition
	}
	actions := make([]sentry.ActionType, len(inputActions))
	for i, ia := range inputActions {
		var action sentry.ActionType
		mapstructure.WeakDecode(ia, &action)
		actions[i] = action
	}
	filters := make([]sentry.FilterType, len(inputFilters))
	for i, ia := range inputFilters {
		var filter sentry.FilterType
		mapstructure.WeakDecode(ia, &filter)
		filters[i] = filter
	}

	params := &sentry.Rule{
		ID:          id,
		ActionMatch: actionMatch,
		FilterMatch: filterMatch,
		Frequency:   frequency,
		Name:        name,
		Conditions:  conditions,
		Actions:     actions,
		Filters:     filters,
	}

	if environment != "" {
		params.Environment = &environment
	}

	tflog.Debug(ctx, "Updating Sentry rule", "ruleID", id, "org", org, "project", project)
	rule, _, err := client.Rules.Update(org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry rule", "ruleID", rule.ID, "org", org, "project", project)

	return resourceSentryRuleRead(ctx, d, meta)
}

func resourceSentryRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Deleting Sentry rule", "ruleID", id, "org", org, "project", project)
	_, err := client.Rules.Delete(org, project, id)
	tflog.Debug(ctx, "Deleted Sentry rule", "ruleID", id, "org", org, "project", project)

	return diag.FromErr(err)
}
