package sentry

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAPIErrors(t *testing.T) {
	tests := []struct {
		name    string
		have    string
		want    string
		wantErr bool
	}{{
		name: "detail",
		have: `{"detail": "description"}`,
		want: "sentry: description",
	}, {
		name: "detail+others",
		have: `{"detail": "description", "other": "field"}`,
		want: "sentry: map[detail:description other:field]",
	}, {
		name: "jsonstring",
		have: `"jsonstring"`,
		want: "sentry: jsonstring",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// https://github.com/dghubble/sling/blob/master/sling.go#L412
			decoder := json.NewDecoder(strings.NewReader(tt.have))

			var e APIError
			err := decoder.Decode(&e)
			if err != nil {
				if !tt.wantErr {
					t.Fatal(err)
				}
				return
			}
			got := e.Error()
			if tt.want != got {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
