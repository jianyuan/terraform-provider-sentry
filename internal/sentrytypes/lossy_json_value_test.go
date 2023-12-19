package sentrytypes

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestLossyJsonStringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentJson   LossyJson
		givenJson     basetypes.StringValuable
		expectedMatch bool
		expectedDiags diag.Diagnostics
	}{
		"not equal - mismatched field values": {
			currentJson:   NewLossyJsonValue(`{"hello": "dlrow", "nums": [3, 2, 1], "nested": {"test-bool": false}}`),
			givenJson:     NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			expectedMatch: false,
		},
		"not equal - mismatched field names": {
			currentJson:   NewLossyJsonValue(`{"Hello": "world", "Nums": [1, 2, 3], "Nested": {"Test-bool": true}}`),
			givenJson:     NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			expectedMatch: false,
		},
		"not equal - array item order difference": {
			currentJson:   NewLossyJsonValue(`[{"nums":[1, 2, 3]}, {"hello": "world"}, {"nested": {"test-bool": true}}]`),
			givenJson:     NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true}}]`),
			expectedMatch: false,
		},
		"semantically equal - null match": {
			currentJson:   NewLossyJsonValue(`null`),
			givenJson:     NewLossyJsonValue(`null`),
			expectedMatch: true,
		},
		"semantically equal - object byte-for-byte match": {
			currentJson:   NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			givenJson:     NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			expectedMatch: true,
		},
		"semantically equal - array byte-for-byte match": {
			currentJson:   NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true}}]`),
			givenJson:     NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true}}]`),
			expectedMatch: true,
		},
		"semantically equal - object field order difference": {
			currentJson:   NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			givenJson:     NewLossyJsonValue(`{"nums": [1, 2, 3], "nested": {"test-bool": true}, "hello": "world"}`),
			expectedMatch: true,
		},
		"semantically equal - object whitespace difference": {
			currentJson: NewLossyJsonValue(`{
				"hello": "world",
				"nums": [1, 2, 3],
				"nested": {
					"test-bool": true
				}
			}`),
			givenJson:     NewLossyJsonValue(`{"hello":"world","nums":[1,2,3],"nested":{"test-bool":true}}`),
			expectedMatch: true,
		},
		"semantically equal - array whitespace difference": {
			currentJson: NewLossyJsonValue(`[
				{
				  "hello": "world"
				},
				{
				  "nums": [
					1,
					2,
					3
				  ]
				},
				{
				  "nested": {
					"test-bool": true
				  }
				}
			  ]`),
			givenJson:     NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true}}]`),
			expectedMatch: true,
		},
		"error - invalid json": {
			currentJson:   NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			givenJson:     NewLossyJsonValue(`&#$^"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			expectedMatch: false,
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Semantic Equality Check Error",
					"An unexpected error occurred while performing semantic equality checks. "+
						"Please report this to the provider developers.\n\n"+
						"Error: invalid character '&' looking for beginning of value",
				),
			},
		},
		"error - not given lossy json value": {
			currentJson:   NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			givenJson:     basetypes.NewStringValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			expectedMatch: false,
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Semantic Equality Check Error",
					"An unexpected value type was received while performing semantic equality checks. "+
						"Please report this to the provider developers.\n\n"+
						"Expected Value Type: sentrytypes.LossyJson\n"+
						"Got Value Type: basetypes.StringValue",
				),
			},
		},
		// JSON Semantic equality uses (decoder).UseNumber to avoid Go parsing JSON numbers into float64. This ensures that Go
		// won't normalize the JSON number representation or impose limits on numeric range.
		"not equal - different JSON number representations": {
			currentJson:   NewLossyJsonValue(`{"large": 12423434}`),
			givenJson:     NewLossyJsonValue(`{"large": 1.2423434e+07}`),
			expectedMatch: false,
		},
		"semantically equal - larger than max float64 values": {
			currentJson:   NewLossyJsonValue(`{"large": 1.79769313486231570814527423731704356798070e+309}`),
			givenJson:     NewLossyJsonValue(`{"large": 1.79769313486231570814527423731704356798070e+309}`),
			expectedMatch: true,
		},
		// JSON Semantic equality uses Go's encoding/json library, which replaces some characters to escape codes
		"semantically equal - HTML escape characters are equal": {
			currentJson:   NewLossyJsonValue(`{"url_ampersand": "http://example.com?foo=bar&hello=world", "left-caret": "<", "right-caret": ">"}`),
			givenJson:     NewLossyJsonValue(`{"url_ampersand": "http://example.com?foo=bar\u0026hello=world", "left-caret": "\u003c", "right-caret": "\u003e"}`),
			expectedMatch: true,
		},
		"semantically equal (lossy) - string number representation": {
			currentJson:   NewLossyJsonValue(`{"large": "12423434"}`),
			givenJson:     NewLossyJsonValue(`{"large": 12423434}`),
			expectedMatch: true,
		},
		"semantically equal (lossy) - object additional field": {
			currentJson:   NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}, "new-field": null}`),
			givenJson:     NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "nested": {"test-bool": true}}`),
			expectedMatch: true,
		},
		"semantically equal (lossy) - array additional field": {
			currentJson:   NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true, "new-field": null}}]`),
			givenJson:     NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true}}]`),
			expectedMatch: true,
		},
		"not equal (lossy) - object missing field": {
			currentJson:   NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true}}]`),
			givenJson:     NewLossyJsonValue(`[{"hello": "world"}, {"nums":[1, 2, 3]}, {"nested": {"test-bool": true, "new-field": null}}]`),
			expectedMatch: false,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentJson.StringSemanticEquals(context.Background(), testCase.givenJson)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestLossyJsonUnmarshal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		json          LossyJson
		target        any
		output        any
		expectedDiags diag.Diagnostics
	}{
		"lossy value is null ": {
			json: NewLossyJsonNull(),
			target: struct {
				Hello string `json:"hello"`
			}{},
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Lossy JSON Unmarshal Error",
					"json string value is null",
				),
			},
		},
		"lossy value is unknown ": {
			json: NewLossyJsonUnknown(),
			target: struct {
				Hello string `json:"hello"`
			}{},
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Lossy JSON Unmarshal Error",
					"json string value is unknown",
				),
			},
		},
		"invalid target - not a pointer ": {
			json: NewLossyJsonValue(`{"hello": "world"}`),
			target: struct {
				Hello string `json:"hello"`
			}{},
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Lossy JSON Unmarshal Error",
					"json: Unmarshal(non-pointer struct { Hello string \"json:\\\"hello\\\"\" })",
				),
			},
		},
		"valid target ": {
			json: NewLossyJsonValue(`{"hello": "world", "nums": [1, 2, 3], "test-bool": true}`),
			target: &struct {
				Hello   string `json:"hello"`
				Numbers []int  `json:"nums"`
				Test    bool   `json:"test-bool"`
			}{},
			output: &struct {
				Hello   string `json:"hello"`
				Numbers []int  `json:"nums"`
				Test    bool   `json:"test-bool"`
			}{
				Hello:   "world",
				Numbers: []int{1, 2, 3},
				Test:    true,
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := testCase.json.Unmarshal(testCase.target)

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
