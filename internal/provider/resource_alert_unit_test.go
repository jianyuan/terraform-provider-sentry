package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func TestAlertResource_frequencyTriggerConditions(t *testing.T) {
	ctx := context.Background()
	data := AlertResourceModel{
		TriggerConditions: supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelTriggerConditionsItem{
			newAlertTriggerConditionsItem(ctx, func(item *AlertResourceModelTriggerConditionsItem) {
				item.EventFrequencyCount = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemEventFrequencyCount{
					Value:    supertypes.NewInt64Value(0),
					Interval: supertypes.NewStringValue("1m"),
					Filters: supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelTriggerConditionsItemEventFrequencyCountFiltersItem{
						{
							Attribute: supertypes.NewStringValue("message"),
							Match:     supertypes.NewStringValue("co"),
							Value:     supertypes.NewStringValue("error"),
						},
					}),
				})
			}),
			newAlertTriggerConditionsItem(ctx, func(item *AlertResourceModelTriggerConditionsItem) {
				item.EventFrequencyPercent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemEventFrequencyPercent{
					Value:              supertypes.NewInt64Value(25),
					Interval:           supertypes.NewStringValue("1h"),
					ComparisonInterval: supertypes.NewStringValue("1w"),
					Filters:            supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelTriggerConditionsItemEventFrequencyPercentFiltersItem{}),
				})
			}),
			newAlertTriggerConditionsItem(ctx, func(item *AlertResourceModelTriggerConditionsItem) {
				item.EventUniqueUserFrequencyPercent = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemEventUniqueUserFrequencyPercent{
					Value:              supertypes.NewInt64Value(30),
					Interval:           supertypes.NewStringValue("15m"),
					ComparisonInterval: supertypes.NewStringValue("1d"),
					Filters:            supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelTriggerConditionsItemEventUniqueUserFrequencyPercentFiltersItem{}),
				})
			}),
			newAlertTriggerConditionsItem(ctx, func(item *AlertResourceModelTriggerConditionsItem) {
				item.PercentSessionsCount = supertypes.NewSingleNestedObjectValueOf(ctx, &AlertResourceModelTriggerConditionsItemPercentSessionsCount{
					Value:    types.Float64Value(17.2),
					Interval: supertypes.NewStringValue("30m"),
					Filters: supertypes.NewListNestedObjectValueOfValueSlice(ctx, []AlertResourceModelTriggerConditionsItemPercentSessionsCountFiltersItem{
						{
							Key:   supertypes.NewStringValue("environment"),
							Match: supertypes.NewStringValue("is"),
						},
					}),
				})
			}),
		}),
		LegacyTriggerConditions: supertypes.NewListValueOfNull[string](ctx),
	}

	got, diags := (&AlertResource{}).getTriggerConditions(ctx, data)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 trigger conditions, got %d", len(got))
	}

	if got[0].Type != "event_frequency_count" {
		t.Fatalf("expected event_frequency_count, got %q", got[0].Type)
	}
	frequencyCountComparison, err := got[0].Comparison.AsOrganizationWorkflowTriggerConditionComparison1()
	if err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	if frequencyCountComparison["value"] != float64(0) {
		t.Fatalf("expected value 0, got %#v", frequencyCountComparison["value"])
	}
	filters, ok := frequencyCountComparison["filters"].([]any)
	if !ok || len(filters) != 1 {
		t.Fatalf("expected one filter, got %#v", frequencyCountComparison["filters"])
	}

	frequencyPercentComparison, err := got[1].Comparison.AsOrganizationWorkflowTriggerConditionComparison1()
	if err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	if frequencyPercentComparison["comparison_interval"] != "1w" {
		t.Fatalf("expected snake_case comparison_interval, got %#v", frequencyPercentComparison)
	}
	if _, ok := frequencyPercentComparison["comparisonInterval"]; ok {
		t.Fatalf("did not expect camelCase comparisonInterval in %#v", frequencyPercentComparison)
	}

	if got[2].Type != "event_unique_user_frequency_percent" {
		t.Fatalf("expected event_unique_user_frequency_percent, got %q", got[2].Type)
	}
	uniqueUserPercentComparison, err := got[2].Comparison.AsOrganizationWorkflowTriggerConditionComparison1()
	if err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	if uniqueUserPercentComparison["comparison_interval"] != "1d" {
		t.Fatalf("expected unique user comparison_interval, got %#v", uniqueUserPercentComparison)
	}

	if got[3].Type != "percent_sessions_count" {
		t.Fatalf("expected percent_sessions_count, got %q", got[3].Type)
	}
	percentSessionsCountComparison, err := got[3].Comparison.AsOrganizationWorkflowTriggerConditionComparison1()
	if err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	if percentSessionsCountComparison["value"] != 17.2 {
		t.Fatalf("expected fractional percent session value, got %#v", percentSessionsCountComparison["value"])
	}
	filters, ok = percentSessionsCountComparison["filters"].([]any)
	if !ok || len(filters) != 1 {
		t.Fatalf("expected one percent sessions filter, got %#v", percentSessionsCountComparison["filters"])
	}
}

func TestAlertResource_fillFrequencyTriggerConditions(t *testing.T) {
	ctx := context.Background()
	firstId := json.Number("1")
	secondId := json.Number("2")
	var frequencyPercentComparison apiclient.OrganizationWorkflowTriggerCondition_Comparison
	if err := frequencyPercentComparison.FromOrganizationWorkflowTriggerConditionComparison1(map[string]any{
		"interval":            "1h",
		"value":               float64(50),
		"comparison_interval": "1w",
	}); err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	var percentSessionsCountComparison apiclient.OrganizationWorkflowTriggerCondition_Comparison
	if err := percentSessionsCountComparison.FromOrganizationWorkflowTriggerConditionComparison1(map[string]any{
		"interval": "30m",
		"value":    17.2,
		"filters": []map[string]any{
			{
				"key":   "environment",
				"match": "is",
			},
		},
	}); err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	var triggers apiclient.OrganizationWorkflow_Triggers
	if err := triggers.FromOrganizationWorkflowTrigger(apiclient.OrganizationWorkflowTrigger{
		LogicType: apiclient.OrganizationWorkflowTriggerLogicTypeAnyShort,
		Conditions: []apiclient.OrganizationWorkflowTriggerCondition{
			{
				Id:              &firstId,
				Type:            "event_frequency_percent",
				Comparison:      frequencyPercentComparison,
				ConditionResult: true,
			},
			{
				Id:              &secondId,
				Type:            "percent_sessions_count",
				Comparison:      percentSessionsCountComparison,
				ConditionResult: true,
			},
		},
	}); err != nil {
		t.Fatalf("unexpected triggers error: %v", err)
	}
	var actionFilters apiclient.OrganizationWorkflow_ActionFilters
	if err := actionFilters.FromOrganizationWorkflowActionFilters0([]apiclient.OrganizationWorkflowActionFilter{}); err != nil {
		t.Fatalf("unexpected action filters error: %v", err)
	}

	var model AlertResourceModel
	diags := model.Fill(ctx, apiclient.OrganizationWorkflow{
		ActionFilters: actionFilters,
		Config: apiclient.OrganizationWorkflowConfig{
			Frequency: 1440,
		},
		DetectorIds: []string{"1"},
		Enabled:     true,
		Id:          "1",
		Name:        "test",
		Triggers:    triggers,
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	triggerConditions := model.TriggerConditions.DiagsGet(ctx, diags)
	if len(triggerConditions) != 2 {
		t.Fatalf("expected two trigger conditions, got %d", len(triggerConditions))
	}
	condition := triggerConditions[0].EventFrequencyPercent.DiagsGet(ctx, diags)
	if condition.Interval.Get() != "1h" || condition.Value.Get() != 50 || condition.ComparisonInterval.Get() != "1w" {
		t.Fatalf("unexpected trigger condition: %#v", condition)
	}
	if !condition.Filters.IsKnown() {
		t.Fatalf("expected empty filters list to be known")
	}
	if filters := condition.Filters.DiagsGet(ctx, diags); len(filters) != 0 {
		t.Fatalf("expected no filters, got %#v", filters)
	}

	percentSessionsCount := triggerConditions[1].PercentSessionsCount.DiagsGet(ctx, diags)
	if percentSessionsCount.Interval.Get() != "30m" || percentSessionsCount.Value.ValueFloat64() != 17.2 {
		t.Fatalf("unexpected percent sessions trigger condition: %#v", percentSessionsCount)
	}
	if filters := percentSessionsCount.Filters.DiagsGet(ctx, diags); len(filters) != 1 {
		t.Fatalf("expected one percent sessions filter, got %#v", filters)
	}
}

// Sentry's workflow DataConditionSerializer camelCases the comparison blob on
// read (convert_dict_key_case(comparison, snake_to_camel_case)), so a GET
// response returns "comparisonInterval" even though writes must send snake_case
// "comparison_interval". This guards the camelCase fallback in
// freqComparisonComparisonInterval so it is not mistaken for dead code: on the
// real API this fallback is the branch that actually populates state.
func TestAlertResource_fillFrequencyTriggerConditions_camelCaseComparisonInterval(t *testing.T) {
	ctx := context.Background()
	id := json.Number("1")
	var comparison apiclient.OrganizationWorkflowTriggerCondition_Comparison
	if err := comparison.FromOrganizationWorkflowTriggerConditionComparison1(map[string]any{
		"interval":           "1h",
		"value":              float64(25),
		"comparisonInterval": "1w",
	}); err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	var triggers apiclient.OrganizationWorkflow_Triggers
	if err := triggers.FromOrganizationWorkflowTrigger(apiclient.OrganizationWorkflowTrigger{
		LogicType: apiclient.OrganizationWorkflowTriggerLogicTypeAnyShort,
		Conditions: []apiclient.OrganizationWorkflowTriggerCondition{
			{
				Id:              &id,
				Type:            "event_frequency_percent",
				Comparison:      comparison,
				ConditionResult: true,
			},
		},
	}); err != nil {
		t.Fatalf("unexpected triggers error: %v", err)
	}
	var actionFilters apiclient.OrganizationWorkflow_ActionFilters
	if err := actionFilters.FromOrganizationWorkflowActionFilters0([]apiclient.OrganizationWorkflowActionFilter{}); err != nil {
		t.Fatalf("unexpected action filters error: %v", err)
	}

	var model AlertResourceModel
	diags := model.Fill(ctx, apiclient.OrganizationWorkflow{
		ActionFilters: actionFilters,
		Config:        apiclient.OrganizationWorkflowConfig{Frequency: 1440},
		DetectorIds:   []string{"1"},
		Enabled:       true,
		Id:            "1",
		Name:          "test",
		Triggers:      triggers,
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	triggerConditions := model.TriggerConditions.DiagsGet(ctx, diags)
	if len(triggerConditions) != 1 {
		t.Fatalf("expected one trigger condition, got %d", len(triggerConditions))
	}
	condition := triggerConditions[0].EventFrequencyPercent.DiagsGet(ctx, diags)
	if condition.ComparisonInterval.Get() != "1w" {
		t.Fatalf("expected camelCase comparisonInterval to populate comparison_interval, got %q", condition.ComparisonInterval.Get())
	}
	if condition.Interval.Get() != "1h" || condition.Value.Get() != 25 {
		t.Fatalf("unexpected trigger condition: %#v", condition)
	}
}

func newAlertTriggerConditionsItem(ctx context.Context, set func(*AlertResourceModelTriggerConditionsItem)) AlertResourceModelTriggerConditionsItem {
	item := AlertResourceModelTriggerConditionsItem{
		FirstSeenEvent:                  supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemFirstSeenEvent](ctx),
		IssueResolvedTrigger:            supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemIssueResolvedTrigger](ctx),
		ReappearedEvent:                 supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemReappearedEvent](ctx),
		RegressionEvent:                 supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemRegressionEvent](ctx),
		EventFrequencyCount:             supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventFrequencyCount](ctx),
		EventUniqueUserFrequencyCount:   supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventUniqueUserFrequencyCount](ctx),
		EventUniqueUserFrequencyPercent: supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventUniqueUserFrequencyPercent](ctx),
		EventFrequencyPercent:           supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemEventFrequencyPercent](ctx),
		PercentSessionsCount:            supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemPercentSessionsCount](ctx),
		PercentSessionsPercent:          supertypes.NewSingleNestedObjectValueOfNull[AlertResourceModelTriggerConditionsItemPercentSessionsPercent](ctx),
	}
	set(&item)
	return item
}
