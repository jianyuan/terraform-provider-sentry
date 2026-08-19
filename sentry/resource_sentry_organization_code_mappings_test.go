package sentry

import (
	"testing"
)

func TestOrganizationCodeMappingsListPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		org           string
		integrationId string
		projectId     string
		cursor        string
		want          string
	}{
		{
			name:          "includes project filter so Sentry returns a stable page for one project",
			org:           "my-org",
			integrationId: "1",
			projectId:     "39",
			want:          "0/organizations/my-org/code-mappings/?integrationId=1&per_page=100&project=39",
		},
		{
			name:      "omits empty optional filters",
			org:       "my-org",
			projectId: "39",
			want:      "0/organizations/my-org/code-mappings/?per_page=100&project=39",
		},
		{
			name:   "keeps org-wide list when project is unknown (import)",
			org:    "my-org",
			cursor: "100:1:0",
			want:   "0/organizations/my-org/code-mappings/?cursor=100%3A1%3A0&per_page=100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := organizationCodeMappingsListPath(tt.org, tt.integrationId, tt.projectId, tt.cursor)
			if got != tt.want {
				t.Fatalf("organizationCodeMappingsListPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
