package sentry

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/mitchellh/mapstructure"
)

const (
	defaultActionMatch = "any"
	defaultFilterMatch = "any"
	defaultFrequency   = 30
)

func resourceSentryIssueAlert() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Issue Alert resource. Note that there's no public documentation for the " +
			"values of conditions, filters, and actions. You can either inspect the request " +
			"payload sent when creating or editing an alert rule on Sentry or inspect " +
			"[Sentry's rules registry in the source code](https://github.com/getsentry/sentry/tree/master/src/sentry/rules).",

		CreateContext: resourceSentryIssueAlertCreate,
		ReadContext:   resourceSentryIssueAlertRead,
		UpdateContext: resourceSentryIssueAlertUpdate,
		DeleteContext: resourceSentryIssueAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importOrganizationProjectAndID,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The slug of the project to create the plugin for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The rule name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action_match": {
				Description:  "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen. Defaults to `any`.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
			},
			"filter_match": {
				Description:  "Trigger actions if `all`, `any`, or `none` of the specified filters match. Defaults to `any`.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "any", "none"}, false),
			},
			"actions": {
				Description: "List of actions.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"conditions": {
				Description: "List of conditions.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"filters": {
				Description: "List of filters.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"frequency": {
				Description: "Perform actions at most once every `X` minutes for this issue. Defaults to `30`.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"environment": {
				Description: "Perform rule in a specific environment.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceSentryIssueAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	params := &sentry.CreateIssueAlertParams{
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

	tflog.Debug(ctx, "Creating Sentry rule", map[string]interface{}{
		"ruleName": name,
		"org":      org,
		"project":  project,
	})
	rule, _, err := client.IssueAlerts.Create(ctx, org, project, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry rule", map[string]interface{}{
		"ruleName": rule.Name,
		"ruleID":   rule.ID,
		"org":      org,
		"project":  project,
	})

	d.SetId(rule.ID)

	return resourceSentryIssueAlertRead(ctx, d, meta)
}

func resourceSentryIssueAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	id := d.Id()

	tflog.Debug(ctx, "Reading Sentry rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})
	rules, resp, err := client.IssueAlerts.List(ctx, org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, "Read Sentry rules", map[string]interface{}{
		"ruleCount": len(rules),
		"rules":     rules,
	})

	var rule *sentry.IssueAlert
	for _, r := range rules {
		if r.ID == id {
			rule = r
			break
		}
	}

	if rule == nil {
		return diag.Errorf("Could not find rule with ID " + id)
	}
	tflog.Debug(ctx, "Read Sentry rule", map[string]interface{}{
		"ruleID":  rule.ID,
		"org":     org,
		"project": project,
	})

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

func resourceSentryIssueAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	params := &sentry.IssueAlert{
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

	tflog.Debug(ctx, "Updating Sentry rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})
	rule, _, err := client.IssueAlerts.Update(ctx, org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry rule", map[string]interface{}{
		"ruleID":  rule.ID,
		"org":     org,
		"project": project,
	})

	return resourceSentryIssueAlertRead(ctx, d, meta)
}

func resourceSentryIssueAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Deleting Sentry rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})
	_, err := client.IssueAlerts.Delete(ctx, org, project, id)
	tflog.Debug(ctx, "Deleted Sentry rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})

	return diag.FromErr(err)
}
