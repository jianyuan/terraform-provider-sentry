package sentry

import (
	"reflect"
	"testing"
)

func TestFollowShape(t *testing.T) {
	testCases := []struct {
		name  string
		shape interface{}
		value interface{}
		want  interface{}
	}{
		{
			name:  "shape is nil",
			shape: nil,
			value: []interface{}{
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
		},
		{
			name: "complex shape",
			shape: []interface{}{
				map[string]interface{}{
					"a": "",
				},
				map[string]interface{}{
					"b": "",
				},
				map[string]interface{}{
					"c": 0,
				},
			},
			value: []interface{}{
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": 3,
				},
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": 3,
				},
				map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": 3,
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"a": "a",
				},
				map[string]interface{}{
					"b": "b",
				},
				map[string]interface{}{
					"c": 3,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := followShape(tc.shape, tc.value)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
