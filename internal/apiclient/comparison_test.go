package apiclient

import (
	"encoding/json"
	"testing"
)

func TestProjectMonitorConditionGroupCondition_FloatComparison(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected float64
	}{
		{
			name:     "integer comparison",
			jsonData: `{"type":"gt","comparison":100,"conditionResult":75}`,
			expected: 100,
		},
		{
			name:     "float comparison - 0.25",
			jsonData: `{"type":"gt","comparison":0.25,"conditionResult":75}`,
			expected: 0.25,
		},
		{
			name:     "float comparison - 0.1",
			jsonData: `{"type":"gt","comparison":0.1,"conditionResult":50}`,
			expected: 0.1,
		},
		{
			name:     "small float comparison - 0.07",
			jsonData: `{"type":"gt","comparison":0.07,"conditionResult":75}`,
			expected: 0.07,
		},
		{
			name:     "negative float comparison",
			jsonData: `{"type":"gt","comparison":-0.5,"conditionResult":75}`,
			expected: -0.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var condition ProjectMonitorConditionGroupCondition
			err := json.Unmarshal([]byte(tc.jsonData), &condition)
			if err != nil {
				t.Fatalf("failed to unmarshal JSON: %v", err)
			}

			comparison, err := condition.Comparison.AsProjectMonitorConditionGroupConditionComparison1()
			if err != nil {
				t.Fatalf("failed to get comparison as number: %v", err)
			}

			floatValue, err := comparison.Float64()
			if err != nil {
				t.Fatalf("failed to convert to float64: %v", err)
			}

			if floatValue != tc.expected {
				t.Errorf("expected %f, got %f", tc.expected, floatValue)
			}
		})
	}
}

func TestProjectMonitorConditionGroupCondition_AnomalyDetection(t *testing.T) {
	jsonData := `{"type":"anomaly_detection","comparison":{"seasonality":"auto","sensitivity":"high","thresholdType":0},"conditionResult":75}`

	var condition ProjectMonitorConditionGroupCondition
	err := json.Unmarshal([]byte(jsonData), &condition)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// This should fail for number comparison
	_, err = condition.Comparison.AsProjectMonitorConditionGroupConditionComparison1()
	if err == nil {
		t.Error("expected error when trying to get anomaly detection as number")
	}

	// This should succeed for object comparison
	comparison, err := condition.Comparison.AsProjectMonitorConditionGroupConditionComparison2()
	if err != nil {
		t.Fatalf("failed to get comparison as object: %v", err)
	}

	if comparison.Sensitivity != "high" {
		t.Errorf("expected sensitivity 'high', got %q", comparison.Sensitivity)
	}
	if comparison.Seasonality != "auto" {
		t.Errorf("expected seasonality 'auto', got %q", comparison.Seasonality)
	}
}
