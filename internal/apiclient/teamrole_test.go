package apiclient

import (
	"encoding/json"
	"testing"

	"github.com/oapi-codegen/nullable"
)

func TestTeamRoleMarshalRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  string
	}{
		{"explicit null role", `{"teamSlug":"test","role":null}`, `{"role":null,"teamSlug":"test"}`},
		{"valued role", `{"teamSlug":"test","role":"admin"}`, `{"role":"admin","teamSlug":"test"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var tr TeamRole
			if err := json.Unmarshal([]byte(tc.in), &tr); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			b, err := json.Marshal(tr)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if string(b) != tc.out {
				t.Errorf("round-trip mismatch:\n  in:  %s\n  got: %s\n  want:%s", tc.in, b, tc.out)
			}
		})
	}
}

func TestTeamRoleMarshalConstructed(t *testing.T) {
	// Null roles should marshal as an explicit null
	tr := TeamRole{TeamSlug: "test", Role: nullable.NewNullNullable[string]()}
	b, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	want := `{"role":null,"teamSlug":"test"}`
	if string(b) != want {
		t.Errorf("got %s, want %s", b, want)
	}
}
