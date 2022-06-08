package sentry

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/mitchellh/mapstructure"
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
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the issue alert belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"conditions": {
				Description: "List of conditions.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				DiffSuppressFunc: SuppressEquivalentJSONDiffs,
			},
			"filters": {
				Description: "List of filters.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				DiffSuppressFunc: SuppressEquivalentJSONDiffs,
			},
			"actions": {
				Description: "List of actions.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				DiffSuppressFunc: SuppressEquivalentJSONDiffs,
			},
			"action_match": {
				Description:  "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
			},
			"filter_match": {
				Description:  "Trigger actions if `all`, `any`, or `none` of the specified filters match.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "any", "none"}, false),
			},
			"frequency": {
				Description: "Perform actions at most once every `X` minutes for this issue. Defaults to `30`.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description: "The issue alert name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"environment": {
				Description: "Perform issue alert in a specific environment.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"project": {
				Description: "The slug of the project to create the issue alert for.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceSentryIssueAlertObject(d *schema.ResourceData) *sentry.IssueAlert {
	alert := &sentry.IssueAlert{
		ActionMatch: sentry.String(d.Get("action_match").(string)),
		FilterMatch: sentry.String(d.Get("filter_match").(string)),
		Frequency:   sentry.Int(d.Get("frequency").(int)),
		Name:        sentry.String(d.Get("name").(string)),
	}

	conditionsIn := d.Get("conditions").([]interface{})
	filtersIn := d.Get("filters").([]interface{})
	actionsIn := d.Get("actions").([]interface{})

	alert.Conditions = make([]*sentry.IssueAlertCondition, 0, len(conditionsIn))
	for _, ic := range conditionsIn {
		condition := new(sentry.IssueAlertCondition)
		mapstructure.WeakDecode(ic, condition)
		alert.Conditions = append(alert.Conditions, condition)
	}
	alert.Filters = make([]*sentry.IssueAlertFilter, 0, len(filtersIn))
	for _, ia := range filtersIn {
		filter := new(sentry.IssueAlertFilter)
		mapstructure.WeakDecode(ia, filter)
		alert.Filters = append(alert.Filters, filter)
	}
	alert.Actions = make([]*sentry.IssueAlertAction, 0, len(actionsIn))
	for _, ia := range actionsIn {
		action := new(sentry.IssueAlertAction)
		mapstructure.WeakDecode(ia, action)
		alert.Actions = append(alert.Actions, action)
	}

	if v, ok := d.GetOk("environment"); ok {
		alert.Environment = sentry.String(v.(string))
	}

	if v, ok := d.GetOk("project"); ok {
		alert.Projects = []string{v.(string)}
	}

	return alert
}

func resourceSentryIssueAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	alertReq := resourceSentryIssueAlertObject(d)

	tflog.Debug(ctx, "Creating issue alert", map[string]interface{}{
		"ruleName": alertReq.Name,
		"org":      org,
		"project":  project,
	})
	alert, _, err := client.IssueAlerts.Create(ctx, org, project, alertReq)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created issue alert", map[string]interface{}{
		"ruleName": alert.Name,
		"ruleID":   alert.ID,
		"org":      org,
		"project":  project,
	})

	d.SetId(buildThreePartID(org, project, *alert.ID))
	return resourceSentryIssueAlertRead(ctx, d, meta)
}

func resourceSentryIssueAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org, project, id, err := splitThreePartID(d.Id(), "organization-slug", "project-slug", "id")
	if err != nil {
		diag.FromErr(err)
	}

	tflog.Debug(ctx, "Reading issue alert", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})

	rule, resp, err := client.IssueAlerts.Get(ctx, org, project, id)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}

	if rule == nil {
		d.SetId("")
		return diag.Errorf("Cannot find issue alert with ID " + id)
	}

	tflog.Debug(ctx, "Read issue alert", map[string]interface{}{
		"ruleID":  rule.ID,
		"org":     org,
		"project": project,
	})

	conditions := make([]interface{}, 0, len(rule.Conditions))
	for _, condition := range rule.Conditions {
		conditions = append(conditions, *condition)
	}
	filters := make([]interface{}, 0, len(rule.Filters))
	for _, filter := range rule.Filters {
		filters = append(filters, *filter)
	}
	actions := make([]interface{}, 0, len(rule.Actions))
	for _, action := range rule.Actions {
		actions = append(actions, *action)
	}

	d.SetId(buildThreePartID(org, project, *rule.ID))
	d.Set("organization", org)
	d.Set("conditions", conditions)
	d.Set("filters", filters)
	d.Set("actions", actions)
	d.Set("action_match", rule.ActionMatch)
	d.Set("filter_match", rule.FilterMatch)
	d.Set("frequency", rule.Frequency)
	d.Set("name", rule.Name)
	d.Set("environment", rule.Environment)
	d.Set("project", project)

	return nil
}

func resourceSentryIssueAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, project, id, err := splitThreePartID(d.Id(), "organization-slug", "project-slug", "id")
	alertReq := resourceSentryIssueAlertObject(d)

	tflog.Debug(ctx, "Updating issue alert", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})
	alert, _, err := client.IssueAlerts.Update(ctx, org, project, id, alertReq)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated issue alert", map[string]interface{}{
		"ruleID":  alert.ID,
		"org":     org,
		"project": project,
	})

	return resourceSentryIssueAlertRead(ctx, d, meta)
}

func resourceSentryIssueAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	org, project, id, err := splitThreePartID(d.Id(), "organization-slug", "project-slug", "id")
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*sentry.Client)

	tflog.Debug(ctx, "Deleting issue rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})
	_, err = client.IssueAlerts.Delete(ctx, org, project, id)
	tflog.Debug(ctx, "Deleted issue rule", map[string]interface{}{
		"ruleID":  id,
		"org":     org,
		"project": project,
	})

	return diag.FromErr(err)
}
