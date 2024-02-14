package sentry

import (
	"context"
	"net/http"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/mapstructure"
)

func resourceSentryIssueAlert() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Issue Alert resource. Note that there's no public documentation for the " +
			"values of conditions, filters, and actions. You can either inspect the request " +
			"payload sent when creating or editing an issue alert on Sentry or inspect " +
			"[Sentry's rules registry in the source code](https://github.com/getsentry/sentry/tree/master/src/sentry/rules). " +
			"Since v0.11.2, you should also omit the name property of each condition, filter, and action.",

		CreateContext: resourceSentryIssueAlertCreate,
		ReadContext:   resourceSentryIssueAlertRead,
		UpdateContext: resourceSentryIssueAlertUpdate,
		DeleteContext: resourceSentryIssueAlertDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema:        resourceSentryIssueAlertSchema(),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceSentryIssueAlertResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceSentryIssueAlertStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceSentryIssueAlertSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"organization": {
			Description: "The slug of the organization the issue alert belongs to.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"project": {
			Description: "The slug of the project to create the issue alert for.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"name": {
			Description:  "The issue alert name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
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
		"actions": {
			Description: "List of actions.",
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Schema{
				Type: schema.TypeMap,
			},
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
		"environment": {
			Description: "Perform issue alert in a specific environment.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		"projects": {
			Deprecated:  "Use `project` (singular) instead.",
			Description: "Use `project` (singular) instead.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"internal_id": {
			Description: "The internal ID for this issue alert.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceSentryIssueAlertResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSentryIssueAlertSchema(),
	}
}

func resourceSentryIssueAlertStateUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	id := rawState["id"].(string)
	org := rawState["organization"].(string)
	project := rawState["project"].(string)
	rawState["id"] = buildThreePartID(org, project, id)
	return rawState, nil
}

func resourceSentryIssueAlertObject(d *schema.ResourceData) *sentry.IssueAlert {
	alert := &sentry.IssueAlert{
		Name:        sentry.String(d.Get("name").(string)),
		ActionMatch: sentry.String(d.Get("action_match").(string)),
		FilterMatch: sentry.String(d.Get("filter_match").(string)),
		Frequency:   sentry.Int(d.Get("frequency").(int)),
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
		"org":       org,
		"project":   project,
		"alertName": alertReq.Name,
	})
	alert, _, err := client.IssueAlerts.Create(ctx, org, project, alertReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(org, project, sentry.StringValue(alert.ID)))
	return resourceSentryIssueAlertRead(ctx, d, meta)
}

func resourceSentryIssueAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, project, alertID, err := splitSentryAlertID(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	tflog.Debug(ctx, "Reading issue alert", map[string]interface{}{"org": org, "project": project, "alertID": alertID})
	alert, _, err := client.IssueAlerts.Get(ctx, org, project, alertID)
	if err != nil {
		if sErr, ok := err.(*sentry.ErrorResponse); ok {
			if sErr.Response.StatusCode == http.StatusNotFound {
				tflog.Info(ctx, "Removing issue alert from state because it no longer exists in Sentry", map[string]interface{}{"org": org, "project": project, "alertID": alertID})
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	conditions := followShape(d.Get("conditions"), normalizeSentryIssueAlertProperty(alert.Conditions))
	filters := followShape(d.Get("filters"), normalizeSentryIssueAlertProperty(alert.Filters))
	actions := followShape(d.Get("actions"), normalizeSentryIssueAlertProperty(alert.Actions))

	d.SetId(buildThreePartID(org, project, sentry.StringValue(alert.ID)))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("projects", alert.Projects),
		d.Set("name", alert.Name),
		d.Set("conditions", conditions),
		d.Set("filters", filters),
		d.Set("actions", actions),
		d.Set("action_match", alert.ActionMatch),
		d.Set("filter_match", alert.FilterMatch),
		d.Set("frequency", alert.Frequency),
		d.Set("environment", alert.Environment),
		d.Set("internal_id", alert.ID),
	)
	if len(alert.Projects) == 1 {
		retErr = multierror.Append(
			retErr,
			d.Set("project", alert.Projects[0]),
		)
	}
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryIssueAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, project, alertID, err := splitSentryAlertID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	alertReq := resourceSentryIssueAlertObject(d)

	tflog.Debug(ctx, "Updating issue alert", map[string]interface{}{
		"org":     org,
		"project": project,
		"alertID": alertID,
	})
	_, _, err = client.IssueAlerts.Update(ctx, org, project, alertID, alertReq)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceSentryIssueAlertRead(ctx, d, meta)
}

func resourceSentryIssueAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, project, alertID, err := splitSentryAlertID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting issue alert", map[string]interface{}{
		"org":     org,
		"project": project,
		"alertID": alertID,
	})
	_, err = client.IssueAlerts.Delete(ctx, org, project, alertID)
	return diag.FromErr(err)
}

func normalizeSentryIssueAlertProperty[T interface{ ~map[string]interface{} }](v []*T) []interface{} {
	out := make([]interface{}, 0, len(v))
	for _, c := range v {
		m := make(map[string]interface{})
		for k, v := range *c {
			m[k] = v
		}
		out = append(out, m)
	}
	return out
}
