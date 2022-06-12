package sentry

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func resourceSentryDashboard() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Dashboard resource.",

		CreateContext: resourceSentryDashboardCreate,
		ReadContext:   resourceSentryDashboardRead,
		UpdateContext: resourceSentryDashboardUpdate,
		DeleteContext: resourceSentryDashboardDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the dashboard belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"title": {
				Description: "Dashboard title.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"widget": {
				Description: "Dashboard widgets.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"display_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fields": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"aggregates": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"columns": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"field_aliases": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"conditions": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"order_by": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// Computed
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"widget_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"x": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"y": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"w": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"h": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"min_h": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"internal_id": {
				Description: "The internal ID for this dashboard.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryDashboardObject(d *schema.ResourceData) *sentry.Dashboard {
	dashboard := &sentry.Dashboard{
		Title: sentry.String(d.Get("title").(string)),
	}

	if widgetList, ok := d.GetOk("widget"); ok {
		widgetList := widgetList.([]interface{})
		dashboard.Widgets = make([]*sentry.DashboardWidget, 0, len(widgetList))
		for _, widgetMap := range widgetList {
			widgetMap := widgetMap.(map[string]interface{})
			widget := new(sentry.DashboardWidget)
			if v, ok := widgetMap["id"].(string); ok && v != "" {
				widget.ID = sentry.String(v)
			}
			if v, ok := widgetMap["title"].(string); ok && v != "" {
				widget.Title = sentry.String(v)
			}
			if v, ok := widgetMap["display_type"].(string); ok && v != "" {
				widget.DisplayType = sentry.String(v)
			}
			if v, ok := widgetMap["interval"].(string); ok && v != "" {
				widget.Interval = sentry.String(v)
			}

			if queryList, ok := widgetMap["query"].([]interface{}); ok {
				widget.Queries = make([]*sentry.DashboardWidgetQuery, 0, len(queryList))
				for _, queryMap := range queryList {
					queryMap := queryMap.(map[string]interface{})
					query := new(sentry.DashboardWidgetQuery)
					if v, ok := queryMap["fields"].([]interface{}); ok && len(v) > 0 {
						query.Fields = expandStringList(v)
					}
					if v, ok := queryMap["aggregates"].([]interface{}); ok && len(v) > 0 {
						query.Aggregates = expandStringList(v)
					}
					if v, ok := queryMap["columns"].([]interface{}); ok && len(v) > 0 {
						query.Columns = expandStringList(v)
					}
					if v, ok := queryMap["field_aliases"].([]interface{}); ok && len(v) > 0 {
						query.FieldAliases = expandStringList(v)
					}
					if v, ok := queryMap["name"].(string); ok && v != "" {
						query.Name = sentry.String(v)
					}
					if v, ok := queryMap["conditions"].(string); ok && v != "" {
						query.Conditions = sentry.String(v)
					}
					if v, ok := queryMap["order_by"].(string); ok && v != "" {
						query.OrderBy = sentry.String(v)
					}
					if v, ok := queryMap["id"].(string); ok && v != "" {
						query.ID = sentry.String(v)
					}
					widget.Queries = append(widget.Queries, query)
				}
			}

			if v, ok := widgetMap["widget_type"].(string); ok && v != "" {
				widget.WidgetType = sentry.String(v)
			}
			if v, ok := widgetMap["limit"].(int); ok && v > 0 {
				widget.Limit = sentry.Int(v)
			}

			if layoutList, ok := widgetMap["layout"].([]interface{}); ok && len(layoutList) == 1 {
				layoutMap := layoutList[0].(map[string]interface{})
				widget.Layout = &sentry.DashboardWidgetLayout{
					X:    sentry.Int(layoutMap["x"].(int)),
					Y:    sentry.Int(layoutMap["y"].(int)),
					W:    sentry.Int(layoutMap["w"].(int)),
					H:    sentry.Int(layoutMap["h"].(int)),
					MinH: sentry.Int(layoutMap["min_h"].(int)),
				}
			}
			dashboard.Widgets = append(dashboard.Widgets, widget)
		}
	}

	return dashboard
}

func resourceSentryDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	dashboardReq := resourceSentryDashboardObject(d)

	tflog.Debug(ctx, "Creating dashboard", map[string]interface{}{
		"org":   org,
		"title": dashboardReq.Title,
	})
	dashboard, _, err := client.Dashboards.Create(ctx, org, dashboardReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(org, sentry.StringValue(dashboard.ID)))
	return resourceSentryDashboardRead(ctx, d, meta)
}

func resourceSentryDashboardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, dashboardID, err := splitSentryDashboardID(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	tflog.Debug(ctx, "Reading dashboard", map[string]interface{}{
		"org":         org,
		"dashboardID": dashboardID,
	})
	dashboard, _, err := client.Dashboards.Get(ctx, org, dashboardID)
	if err != nil {
		if sErr, ok := err.(*sentry.ErrorResponse); ok {
			if sErr.Response.StatusCode == http.StatusNotFound {
				tflog.Info(ctx, "Removing dashboard from state because it no longer exists in Sentry", map[string]interface{}{
					"org":         org,
					"dashboardID": dashboardID,
				})
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(org, sentry.StringValue(dashboard.ID)))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("title", dashboard.Title),
		d.Set("widget", flattenDashboardWidgets(dashboard.Widgets)),
		d.Set("internal_id", dashboard.ID),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryDashboardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, dashboardID, err := splitSentryDashboardID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	dashboardReq := resourceSentryDashboardObject(d)

	tflog.Debug(ctx, "Updating dashboard", map[string]interface{}{
		"org":         org,
		"dashboardID": dashboardID,
	})
	_, _, err = client.Dashboards.Update(ctx, org, dashboardID, dashboardReq)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceSentryDashboardRead(ctx, d, meta)
}

func resourceSentryDashboardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, dashboardID, err := splitSentryDashboardID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting dashboard", map[string]interface{}{
		"org":         org,
		"dashboardID": dashboardID,
	})
	_, err = client.Dashboards.Delete(ctx, org, dashboardID)
	return diag.FromErr(err)
}

func splitSentryDashboardID(id string) (org string, dashboardID string, err error) {
	org, dashboardID, err = splitTwoPartID(id, "organization-slug", "dashboard-id")
	return
}

func flattenDashboardWidgets(widgets []*sentry.DashboardWidget) []interface{} {
	if widgets == nil {
		return []interface{}{}
	}

	l := make([]interface{}, 0, len(widgets))
	for _, widget := range widgets {
		m := make(map[string]interface{})
		m["id"] = widget.ID
		m["title"] = widget.Title
		m["display_type"] = widget.DisplayType
		m["interval"] = widget.Interval
		m["query"] = flattenDashboardWidgetQueries(widget.Queries)
		m["widget_type"] = widget.WidgetType
		m["limit"] = widget.Limit
		m["layout"] = flattenDashboardWidgetLayout(widget.Layout)
		l = append(l, m)
	}
	return l
}

func flattenDashboardWidgetQueries(queries []*sentry.DashboardWidgetQuery) []interface{} {
	if queries == nil {
		return []interface{}{}
	}

	l := make([]interface{}, 0, len(queries))
	for _, query := range queries {
		m := make(map[string]interface{})
		m["id"] = query.ID
		m["fields"] = query.Fields
		m["aggregates"] = query.Aggregates
		m["columns"] = query.Columns
		m["field_aliases"] = query.FieldAliases
		m["name"] = query.Name
		m["conditions"] = query.Conditions
		m["order_by"] = query.OrderBy
		l = append(l, m)
	}
	return l
}

func flattenDashboardWidgetLayout(layout *sentry.DashboardWidgetLayout) []interface{} {
	if layout == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})
	m["x"] = layout.X
	m["y"] = layout.Y
	m["w"] = layout.W
	m["h"] = layout.H
	m["min_h"] = layout.MinH

	return []interface{}{m}
}
