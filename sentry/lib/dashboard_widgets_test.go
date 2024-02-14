package sentry

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDashboardWidgetsService_Validate_pass(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/dashboards/widgets/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{}`)
	})

	widget := &DashboardWidget{
		Title:       String("Number of Errors"),
		DisplayType: String("big_number"),
		Interval:    String("5m"),
		Queries: []*DashboardWidgetQuery{
			{
				ID:           String("115037"),
				Fields:       []string{"count()"},
				Aggregates:   []string{"count()"},
				Columns:      []string{},
				FieldAliases: []string{},
				Name:         String(""),
				Conditions:   String("!event.type:transaction"),
				OrderBy:      String(""),
			},
		},
		WidgetType: String("discover"),
		Layout: &DashboardWidgetLayout{
			X:    Int(0),
			Y:    Int(0),
			W:    Int(2),
			H:    Int(1),
			MinH: Int(1),
		},
	}
	ctx := context.Background()
	widgetErrors, _, err := client.DashboardWidgets.Validate(ctx, "the-interstellar-jurisdiction", widget)
	assert.Nil(t, widgetErrors)
	assert.NoError(t, err)
}

func TestDashboardWidgetsService_Validate_fail(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/0/organizations/the-interstellar-jurisdiction/dashboards/widgets/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"widgetType":["\"discover-invalid\" is not a valid choice."]}`)
	})

	widget := &DashboardWidget{
		Title:       String("Number of Errors"),
		DisplayType: String("big_number"),
		Interval:    String("5m"),
		Queries: []*DashboardWidgetQuery{
			{
				ID:           String("115037"),
				Fields:       []string{"count()"},
				Aggregates:   []string{"count()"},
				Columns:      []string{},
				FieldAliases: []string{},
				Name:         String(""),
				Conditions:   String("!event.type:transaction"),
				OrderBy:      String(""),
			},
		},
		WidgetType: String("discover-invalid"),
		Layout: &DashboardWidgetLayout{
			X:    Int(0),
			Y:    Int(0),
			W:    Int(2),
			H:    Int(1),
			MinH: Int(1),
		},
	}
	ctx := context.Background()
	widgetErrors, _, err := client.DashboardWidgets.Validate(ctx, "the-interstellar-jurisdiction", widget)
	expected := DashboardWidgetErrors{
		"widgetType": []string{`"discover-invalid" is not a valid choice.`},
	}
	assert.Equal(t, expected, widgetErrors)
	assert.NoError(t, err)
}
