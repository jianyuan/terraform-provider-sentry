package sentry

import (
	"errors"
	"fmt"

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
		Create: resourceSentryRuleCreate,
		Read:   resourceSentryRuleRead,
		Update: resourceSentryRuleUpdate,
		Delete: resourceSentryRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSentryRuleImporter,
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

func resourceSentryRuleCreate(d *schema.ResourceData, meta interface{}) error {
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

	rule, _, err := client.Rules.Create(org, project, params)
	if err != nil {
		return err
	}

	d.SetId(rule.ID)

	return resourceSentryRuleRead(d, meta)
}

func resourceSentryRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	id := d.Id()

	rules, resp, err := client.Rules.List(org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
	}

	var rule *sentry.Rule
	for _, r := range rules {
		if r.ID == id {
			rule = &r
			break
		}
	}

	if rule == nil {
		return errors.New("Could not find rule with ID " + id)
	}

	// workaround for
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/62
	// as the data sent by Sentry is integer
	for _, f := range rule.Filters {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				// unparseable so forcing this to be int
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	d.SetId(rule.ID)
	d.Set("name", rule.Name)
	d.Set("frequency", rule.Frequency)
	d.Set("environment", rule.Environment)

	return nil
}

func resourceSentryRuleUpdate(d *schema.ResourceData, meta interface{}) error {
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

	_, _, err := client.Rules.Update(org, project, id, params)
	if err != nil {
		return err
	}

	return resourceSentryRuleRead(d, meta)
}

func resourceSentryRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	_, err := client.Rules.Delete(org, project, id)
	return err
}
