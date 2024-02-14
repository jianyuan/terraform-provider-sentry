package sentry

import (
	"context"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSentryIssueAlertSentryIssueAlert() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Issue Alert data source. As the object structure of `conditions`, `filters`, and " + "" +
			"`actions` are undocumented, a tip is to set up an Issue Alert via the Web UI, and use this data source " +
			"to copy its object structure to your resources.",

		ReadContext: dataSourceSentryIssueAlertRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the issue alert belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The slug of the project the issue alert belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "The internal ID for this issue alert.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"conditions": {
				Description: "List of conditions.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"filters": {
				Description: "List of filters.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"actions": {
				Description: "List of actions.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"action_match": {
				Description: "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"filter_match": {
				Description: "Trigger actions if `all`, `any`, or `none` of the specified filters match.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"frequency": {
				Description: "Perform actions at most once every `X` minutes for this issue. Defaults to `30`.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "The issue alert name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environment": {
				Description: "Perform issue alert in a specific environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSentryIssueAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	alertID := d.Get("internal_id").(string)

	tflog.Debug(ctx, "Reading issue alert", map[string]interface{}{"org": org, "project": project, "alertID": alertID})
	alert, _, err := client.IssueAlerts.Get(ctx, org, project, alertID)
	if err != nil {
		return diag.FromErr(err)
	}

	conditions := make([]interface{}, 0, len(alert.Conditions))
	for _, condition := range alert.Conditions {
		conditions = append(conditions, *condition)
	}
	filters := make([]interface{}, 0, len(alert.Filters))
	for _, filter := range alert.Filters {
		filters = append(filters, *filter)
	}
	actions := make([]interface{}, 0, len(alert.Actions))
	for _, action := range alert.Actions {
		actions = append(actions, *action)
	}

	d.SetId(buildThreePartID(org, project, sentry.StringValue(alert.ID)))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("project", project),
		d.Set("internal_id", alert.ID),
		d.Set("conditions", conditions),
		d.Set("filters", filters),
		d.Set("actions", actions),
		d.Set("action_match", alert.ActionMatch),
		d.Set("filter_match", alert.FilterMatch),
		d.Set("frequency", alert.Frequency),
		d.Set("name", alert.Name),
		d.Set("environment", alert.Environment),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
