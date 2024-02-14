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
							// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/models/dashboard_widget.py#L51
							ValidateFunc: validation.StringInSlice(
								[]string{
									"line",
									"area",
									"stacked_area",
									"bar",
									"table",
									"big_number",
									"top_n",
								},
								false,
							),
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
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"aggregates": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"columns": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"field_aliases": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
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
										Computed: true,
									},
									"order_by": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
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
							// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/models/dashboard_widget.py#L39
							ValidateFunc: validation.StringInSlice(
								[]string{
									"discover",
									"issue",
									"metrics",
								},
								false,
							),
						},
						"limit": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
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
			widget.Title = sentry.String(widgetMap["title"].(string))
			widget.DisplayType = sentry.String(widgetMap["display_type"].(string))
			if v := widgetMap["id"].(string); v != "" {
				widget.ID = sentry.String(v)
			}
			if v := widgetMap["interval"].(string); v != "" {
				widget.Interval = sentry.String(v)
			}
			if v := widgetMap["widget_type"].(string); v != "" {
				widget.WidgetType = sentry.String(v)
			}
			if v := widgetMap["limit"].(int); v > 0 {
				widget.Limit = sentry.Int(v)
			}

			if queryList, ok := widgetMap["query"].([]interface{}); ok {
				widget.Queries = make([]*sentry.DashboardWidgetQuery, 0, len(queryList))
				for _, queryMap := range queryList {
					queryMap := queryMap.(map[string]interface{})
					query := new(sentry.DashboardWidgetQuery)
					query.Fields = expandStringList(queryMap["fields"].([]interface{}))
					query.Aggregates = expandStringList(queryMap["aggregates"].(*schema.Set).List())
					query.Columns = expandStringList(queryMap["columns"].(*schema.Set).List())
					query.FieldAliases = expandStringList(queryMap["field_aliases"].([]interface{}))
					query.Name = sentry.String(queryMap["name"].(string))
					query.Conditions = sentry.String(queryMap["conditions"].(string))
					query.OrderBy = sentry.String(queryMap["order_by"].(string))
					if v := queryMap["id"].(string); v != "" {
						query.ID = sentry.String(v)
					}
					widget.Queries = append(widget.Queries, query)
				}
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

	widgetList := make([]interface{}, 0, len(widgets))
	for _, widget := range widgets {
		layoutMap := make(map[string]interface{})
		layoutMap["x"] = widget.Layout.X
		layoutMap["y"] = widget.Layout.Y
		layoutMap["w"] = widget.Layout.W
		layoutMap["h"] = widget.Layout.H
		layoutMap["min_h"] = widget.Layout.MinH

		widgetMap := make(map[string]interface{})
		widgetMap["id"] = widget.ID
		widgetMap["title"] = widget.Title
		widgetMap["display_type"] = widget.DisplayType
		widgetMap["interval"] = widget.Interval
		widgetMap["query"] = flattenDashboardWidgetQueries(widget.Queries)
		widgetMap["widget_type"] = widget.WidgetType
		widgetMap["limit"] = widget.Limit
		widgetMap["layout"] = []interface{}{layoutMap}
		widgetList = append(widgetList, widgetMap)
	}
	return widgetList
}

func flattenDashboardWidgetQueries(queries []*sentry.DashboardWidgetQuery) []interface{} {
	if queries == nil {
		return []interface{}{}
	}

	queryList := make([]interface{}, 0, len(queries))
	for _, query := range queries {
		queryMap := make(map[string]interface{})
		queryMap["id"] = query.ID
		queryMap["fields"] = query.Fields
		queryMap["aggregates"] = flattenStringSet(query.Aggregates)
		queryMap["columns"] = flattenStringSet(query.Columns)
		queryMap["field_aliases"] = query.FieldAliases
		queryMap["name"] = query.Name
		queryMap["conditions"] = query.Conditions
		queryMap["order_by"] = query.OrderBy
		queryList = append(queryList, queryMap)
	}
	return queryList
}
