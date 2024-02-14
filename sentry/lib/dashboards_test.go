package sentry

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDashboardsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/dashboards/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{
				"id": "11833",
				"title": "General",
				"dateCreated": "2022-06-07T16:48:26.255520Z"
			},
			{
				"id": "11832",
				"title": "Mobile Template",
				"dateCreated": "2022-06-07T16:43:40.456607Z"
			}
		]`)
	})

	ctx := context.Background()
	widgetErrors, _, err := client.Dashboards.List(ctx, "the-interstellar-jurisdiction", nil)

	expected := []*Dashboard{
		{
			ID:          String("11833"),
			Title:       String("General"),
			DateCreated: Time(mustParseTime("2022-06-07T16:48:26.255520Z")),
		},
		{
			ID:          String("11832"),
			Title:       String("Mobile Template"),
			DateCreated: Time(mustParseTime("2022-06-07T16:43:40.456607Z")),
		},
	}
	assert.Equal(t, expected, widgetErrors)
	assert.NoError(t, err)
}

func TestDashboardsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/dashboards/12072/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "12072",
			"title": "General",
			"dateCreated": "2022-06-07T16:48:26.255520Z",
			"widgets": [
				{
					"id": "105567",
					"title": "Custom Widget",
					"displayType": "world_map",
					"interval": "5m",
					"dateCreated": "2022-06-12T15:37:19.886736Z",
					"dashboardId": "12072",
					"queries": [
						{
							"id": "117838",
							"name": "",
							"fields": [
								"count()"
							],
							"aggregates": [
								"count()"
							],
							"columns": [],
							"fieldAliases": [],
							"conditions": "",
							"orderby": "",
							"widgetId": "105567"
						}
					],
					"limit": null,
					"widgetType": "discover",
					"layout": {
						"y": 0,
						"x": 0,
						"h": 2,
						"minH": 2,
						"w": 2
					}
				}
			]
		}`)
	})

	ctx := context.Background()
	dashboard, _, err := client.Dashboards.Get(ctx, "the-interstellar-jurisdiction", "12072")

	expected := &Dashboard{
		ID:          String("12072"),
		Title:       String("General"),
		DateCreated: Time(mustParseTime("2022-06-07T16:48:26.255520Z")),
		Widgets: []*DashboardWidget{
			{
				ID:          String("105567"),
				Title:       String("Custom Widget"),
				DisplayType: String("world_map"),
				Interval:    String("5m"),
				Queries: []*DashboardWidgetQuery{
					{
						ID:           String("117838"),
						Fields:       []string{"count()"},
						Aggregates:   []string{"count()"},
						Columns:      []string{},
						FieldAliases: []string{},
						Name:         String(""),
						Conditions:   String(""),
						OrderBy:      String(""),
					},
				},
				WidgetType: String("discover"),
				Layout: &DashboardWidgetLayout{
					X:    Int(0),
					Y:    Int(0),
					W:    Int(2),
					H:    Int(2),
					MinH: Int(2),
				},
			},
		},
	}
	assert.Equal(t, expected, dashboard)
	assert.NoError(t, err)
}

func TestDashboardsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/dashboards/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		assertPostJSONValue(t, map[string]interface{}{
			"title":   "General",
			"widgets": map[string]interface{}{},
		}, r)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "12072",
			"title": "General",
			"dateCreated": "2022-06-07T16:48:26.255520Z",
			"widgets": []
		}`)
	})

	params := &Dashboard{
		Title:   String("General"),
		Widgets: []*DashboardWidget{},
	}
	ctx := context.Background()
	dashboard, _, err := client.Dashboards.Create(ctx, "the-interstellar-jurisdiction", params)

	expected := &Dashboard{
		ID:          String("12072"),
		Title:       String("General"),
		DateCreated: Time(mustParseTime("2022-06-07T16:48:26.255520Z")),
		Widgets:     []*DashboardWidget{},
	}
	assert.Equal(t, expected, dashboard)
	assert.NoError(t, err)
}

func TestDashboardsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/dashboards/12072/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PUT", r)
		assertPostJSONValue(t, map[string]interface{}{
			"id":      "12072",
			"title":   "General",
			"widgets": map[string]interface{}{},
		}, r)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "12072",
			"title": "General",
			"dateCreated": "2022-06-07T16:48:26.255520Z",
			"widgets": []
		}`)
	})

	params := &Dashboard{
		ID:      String("12072"),
		Title:   String("General"),
		Widgets: []*DashboardWidget{},
	}
	ctx := context.Background()
	dashboard, _, err := client.Dashboards.Update(ctx, "the-interstellar-jurisdiction", "12072", params)

	expected := &Dashboard{
		ID:          String("12072"),
		Title:       String("General"),
		DateCreated: Time(mustParseTime("2022-06-07T16:48:26.255520Z")),
		Widgets:     []*DashboardWidget{},
	}
	assert.Equal(t, expected, dashboard)
	assert.NoError(t, err)
}

func TestDashboardsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/dashboards/12072/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
	})

	ctx := context.Background()
	_, err := client.Dashboards.Delete(ctx, "the-interstellar-jurisdiction", "12072")
	assert.NoError(t, err)
}
