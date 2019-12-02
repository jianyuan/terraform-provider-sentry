package sentry

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/mitchellh/mapstructure"
)

const (
	defaultActionMatch = "any"
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
	inputConditions := d.Get("conditions").([]interface{})
	inputActions := d.Get("actions").([]interface{})
	frequency := d.Get("frequency").(int)

	if actionMatch == "" {
		actionMatch = defaultActionMatch
	}
	if frequency == 0 {
		frequency = defaultFrequency
	}

	conditions := make([]*sentry.CreateRuleConditionParams, len(inputConditions))
	for i, ic := range inputConditions {
		var condition sentry.CreateRuleConditionParams
		mapstructure.Decode(ic, &condition)
		conditions[i] = &condition
	}
	actions := make([]*sentry.CreateRuleActionParams, len(inputActions))
	for i, ia := range inputActions {
		var action sentry.CreateRuleActionParams
		mapstructure.Decode(ia, &action)
		actions[i] = &action
	}

	log.Printf("%v, %v, %v", name, org, project)

	params := &sentry.CreateRuleParams{
		ActionMatch: actionMatch,
		Environment: environment,
		Frequency:   frequency,
		Name:        name,
		Conditions:  conditions,
		Actions:     actions,
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

	rules, _, err := client.Rules.List(org, project)
	if err != nil {
		d.SetId("")
		return nil
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

	d.SetId(rule.ID)
	d.Set("name", rule.Name)
	d.Set("actions", rule.Actions)
	d.Set("conditions", rule.Conditions)
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
	inputConditions := d.Get("conditions").([]interface{})
	inputActions := d.Get("actions").([]interface{})
	frequency := d.Get("frequency").(int)

	if actionMatch == "" {
		actionMatch = defaultActionMatch
	}
	if frequency == 0 {
		frequency = defaultFrequency
	}

	conditions := make([]sentry.RuleCondition, len(inputConditions))
	for i, ic := range inputConditions {
		var condition sentry.RuleCondition
		mapstructure.Decode(ic, &condition)
		conditions[i] = condition
	}
	actions := make([]sentry.RuleAction, len(inputActions))
	for i, ia := range inputActions {
		var action sentry.RuleAction
		mapstructure.Decode(ia, &action)
		actions[i] = action
	}

	log.Printf("%v, %v, %v, %v", id, name, org, project)

	params := &sentry.Rule{
		ID:          id,
		ActionMatch: actionMatch,
		Environment: environment,
		Frequency:   frequency,
		Name:        name,
		Conditions:  conditions,
		Actions:     actions,
	}

	if environment != "" {
		params.Environment = environment
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
