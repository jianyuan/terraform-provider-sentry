package sentryplatforms

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		platform string
		want     bool
	}{
		{"android", true},
		{"python", true},
		{"javascript", true},
		{"other", true},
		{"bogus", false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("platform=%s", tc.platform), func(t *testing.T) {
			got := Validate(tc.platform)
			if got != tc.want {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
