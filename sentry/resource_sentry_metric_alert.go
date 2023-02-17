package sentry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func resourceSentryMetricAlert() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Metric Alert resource.",

		CreateContext: resourceSentryMetricAlertCreate,
		ReadContext:   resourceSentryMetricAlertRead,
		UpdateContext: resourceSentryMetricAlertUpdate,
		DeleteContext: resourceSentryMetricAlertDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the metric alert belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The slug of the project to create the metric alert for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The metric alert name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Perform Alert rule in a specific environment",
			},
			"dataset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Sentry Alert category",
			},
			"event_types": {
				Description: "The events type of dataset.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The query filter to apply",
			},
			"aggregate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The aggregation criteria to apply",
			},
			"time_window": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "The period to evaluate the Alert rule in minutes",
			},
			"threshold_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The type of threshold",
			},
			"resolve_threshold": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The value at which the Alert rule resolves",
			},
			"trigger": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_identifier": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"integration_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"input_channel_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Slack channel ID to avoid rate-limiting, see [here](https://docs.sentry.io/product/integrations/notification-incidents/slack/#rate-limiting-error)",
									},
								},
							},
						},
						"label": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold_type": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"alert_threshold": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"resolve_threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the owner id of this Alert rule",
			},
			"internal_id": {
				Description: "The internal ID for this metric alert.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryMetricAlertObject(d *schema.ResourceData) *sentry.MetricAlert {
	alert := &sentry.MetricAlert{
		Name:          sentry.String(d.Get("name").(string)),
		DataSet:       sentry.String(d.Get("dataset").(string)),
		Query:         sentry.String(d.Get("query").(string)),
		Aggregate:     sentry.String(d.Get("aggregate").(string)),
		TimeWindow:    sentry.Float64(d.Get("time_window").(float64)),
		ThresholdType: sentry.Int(d.Get("threshold_type").(int)),
	}
	if v, ok := d.GetOk("internal_id"); ok {
		alert.ID = sentry.String(v.(string))
	}
	if v, ok := d.GetOk("environment"); ok {
		alert.Environment = sentry.String(v.(string))
	}
	if v, ok := d.GetOk("dataset"); ok {
		alert.DataSet = sentry.String(v.(string))
	}
	if v, ok := d.GetOk("event_types"); ok {
		eventTypes := expandStringList(v.([]interface{}))
		if len(eventTypes) > 0 {
			alert.EventTypes = eventTypes
		}
	}
	if v, ok := d.GetOk("resolve_threshold"); ok {
		alert.ResolveThreshold = sentry.Float64(v.(float64))
	}
	if v, ok := d.GetOk("owner"); ok {
		alert.Owner = sentry.String(v.(string))
	}
	if v, ok := d.GetOk("project"); ok {
		alert.Projects = []string{v.(string)}
	}

	triggersIn := d.Get("trigger").([]interface{})
	alert.Triggers = expandMetricAlertTriggers(triggersIn)

	return alert
}

func resourceSentryMetricAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	alertReq := resourceSentryMetricAlertObject(d)

	tflog.Info(ctx, "Creating metric alert", map[string]interface{}{
		"org":      org,
		"project":  project,
		"ruleName": alertReq.Name,
		"params":   fmt.Sprintf("%+v", alertReq),
	})
	alert, _, err := client.MetricAlerts.Create(ctx, org, project, alertReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(org, project, sentry.StringValue(alert.ID)))
	return resourceSentryMetricAlertRead(ctx, d, meta)
}

func resourceSentryMetricAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, project, alertID, err := splitSentryAlertID(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	tflog.Debug(ctx, "Reading metric alert", map[string]interface{}{
		"org":     org,
		"project": project,
		"alertID": alertID,
	})
	alert, _, err := client.MetricAlerts.Get(ctx, org, project, alertID)
	if err != nil {
		if sErr, ok := err.(*sentry.ErrorResponse); ok {
			if sErr.Response.StatusCode == http.StatusNotFound {
				tflog.Info(ctx, "Removing metric alert from state because it no longer exists in Sentry", map[string]interface{}{"org": org})
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read metric alert", map[string]interface{}{
		"alert": fmt.Sprintf("%+v", alert),
	})

	d.SetId(buildThreePartID(org, project, sentry.StringValue(alert.ID)))
	retError := multierror.Append(
		d.Set("organization", org),
		d.Set("name", alert.Name),
		d.Set("environment", alert.Environment),
		d.Set("dataset", alert.DataSet),
		d.Set("event_types", alert.EventTypes),
		d.Set("query", alert.Query),
		d.Set("aggregate", alert.Aggregate),
		d.Set("time_window", alert.TimeWindow),
		d.Set("threshold_type", alert.ThresholdType),
		d.Set("resolve_threshold", alert.ResolveThreshold),
		d.Set("trigger", flattenMetricAlertTriggers(alert.Triggers)),
		d.Set("owner", alert.Owner),
		d.Set("internal_id", alert.ID),
	)
	if len(alert.Projects) == 1 {
		retError = multierror.Append(
			retError,
			d.Set("project", alert.Projects[0]),
		)
	}
	return diag.FromErr(retError.ErrorOrNil())
}

func resourceSentryMetricAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, project, alertID, err := splitSentryAlertID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	alertReq := resourceSentryMetricAlertObject(d)

	tflog.Debug(ctx, "Updating metric alert", map[string]interface{}{
		"org":     org,
		"project": project,
		"alertID": alertID,
	})
	alert, _, err := client.MetricAlerts.Update(ctx, org, project, alertID, alertReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(org, project, sentry.StringValue(alert.ID)))
	return resourceSentryMetricAlertRead(ctx, d, meta)
}

func resourceSentryMetricAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, project, alertID, err := splitSentryAlertID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting metric alert", map[string]interface{}{
		"org":     org,
		"project": project,
		"alertID": alertID,
	})
	_, err = client.MetricAlerts.Delete(ctx, org, project, alertID)
	return diag.FromErr(err)
}

func expandMetricAlertTriggers(triggerList []interface{}) []*sentry.MetricAlertTrigger {
	triggers := make([]*sentry.MetricAlertTrigger, 0, len(triggerList))
	for _, triggerMap := range triggerList {
		triggerMap := triggerMap.(map[string]interface{})
		trigger := &sentry.MetricAlertTrigger{
			Label:            sentry.String(triggerMap["label"].(string)),
			ThresholdType:    sentry.Int(triggerMap["threshold_type"].(int)),
			AlertThreshold:   sentry.Float64(triggerMap["alert_threshold"].(float64)),
			ResolveThreshold: sentry.Float64(triggerMap["resolve_threshold"].(float64)),
			Actions:          expandMetricAlertTriggerActions(triggerMap["action"].([]interface{})),
		}
		if v, ok := triggerMap["id"].(string); ok {
			if v != "" {
				trigger.ID = sentry.String(v)
			}
		}
		triggers = append(triggers, trigger)
	}
	return triggers
}

func expandMetricAlertTriggerActions(actionList []interface{}) []*sentry.MetricAlertTriggerAction {
	actions := make([]*sentry.MetricAlertTriggerAction, 0, len(actionList))
	for _, actionMap := range actionList {
		actionMap := actionMap.(map[string]interface{})
		action := &sentry.MetricAlertTriggerAction{
			Type:       sentry.String(actionMap["type"].(string)),
			TargetType: sentry.String(actionMap["target_type"].(string)),
		}
		if v, ok := actionMap["id"].(string); ok {
			if v != "" {
				action.ID = sentry.String(v)
			}
		}
		if v, ok := actionMap["target_identifier"].(string); ok {
			if v != "" {
				action.TargetIdentifier = sentry.String(v)
			}
		}
		if v, ok := actionMap["integration_id"].(int); ok {
			if v != 0 {
				action.IntegrationID = sentry.Int(v)
			}
		}
		if v, ok := actionMap["input_channel_id"].(string); ok {
			if v != "" {
				action.InputChannelID = sentry.String(v)
			}
		}
		actions = append(actions, action)
	}
	return actions
}

func flattenMetricAlertTriggers(triggers []*sentry.MetricAlertTrigger) []interface{} {
	if triggers == nil {
		return []interface{}{}
	}

	triggerList := make([]interface{}, 0, len(triggers))
	for _, trigger := range triggers {
		triggerMap := make(map[string]interface{})
		triggerMap["id"] = trigger.ID
		triggerMap["label"] = trigger.Label
		triggerMap["threshold_type"] = trigger.ThresholdType
		triggerMap["alert_threshold"] = trigger.AlertThreshold
		triggerMap["resolve_threshold"] = trigger.ResolveThreshold
		triggerMap["action"] = flattenMetricAlertTriggerActions(trigger.Actions)
		triggerList = append(triggerList, triggerMap)
	}
	return triggerList
}

func flattenMetricAlertTriggerActions(actions []*sentry.MetricAlertTriggerAction) []interface{} {
	if actions == nil {
		return []interface{}{}
	}

	actionList := make([]interface{}, 0, len(actions))
	for _, action := range actions {
		actionMap := make(map[string]interface{})
		actionMap["id"] = action.ID
		actionMap["type"] = action.Type
		actionMap["target_type"] = action.TargetType
		actionMap["target_identifier"] = action.TargetIdentifier
		actionMap["integration_id"] = action.IntegrationID
		actionMap["input_channel_id"] = action.InputChannelID

		actionList = append(actionList, actionMap)
	}

	return actionList
}
