package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func TestProjectResourceModelFillRelayAutomaticJsonExpansion(t *testing.T) {
	t.Parallel()

	baseOptions := map[string]interface{}{
		"filters:blacklisted_ips": "",
		"filters:releases":        "",
		"filters:error_messages":  "",
	}

	tests := []struct {
		name    string
		options map[string]interface{}
		want    bool
		wantErr bool
	}{
		{
			name:    "missing option",
			options: baseOptions,
			want:    false,
		},
		{
			name: "enabled",
			options: mergeOptions(baseOptions, map[string]interface{}{
				"sentry:relay_automatic_json_expansion": true,
			}),
			want: true,
		},
		{
			name: "disabled",
			options: mergeOptions(baseOptions, map[string]interface{}{
				"sentry:relay_automatic_json_expansion": false,
			}),
			want: false,
		},
		{
			name: "invalid type",
			options: mergeOptions(baseOptions, map[string]interface{}{
				"sentry:relay_automatic_json_expansion": "true",
			}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var model ProjectResourceModel
			diags := model.Fill(context.Background(), apiclient.Project{
				Slug: "web",
				Name: "Web",
				Organization: apiclient.Organization{
					Slug: "acme",
				},
				Options: tt.options,
			})

			if tt.wantErr {
				if !diags.HasError() {
					t.Fatal("expected a fill error")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if model.RelayAutomaticJsonExpansion.IsNull() || model.RelayAutomaticJsonExpansion.IsUnknown() {
				t.Fatal("expected a known boolean")
			}

			if got := model.RelayAutomaticJsonExpansion.ValueBool(); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssignRelayAutomaticJsonExpansion(t *testing.T) {
	t.Parallel()

	t.Run("skips unset", func(t *testing.T) {
		t.Parallel()

		var options *map[string]interface{}
		assignRelayAutomaticJsonExpansion(&options, types.BoolNull())
		if options != nil {
			t.Fatalf("expected no options map, got %#v", *options)
		}
	})

	t.Run("merges with existing options", func(t *testing.T) {
		t.Parallel()

		existing := map[string]interface{}{
			"filters:releases": "1.*",
		}
		options := &existing
		assignRelayAutomaticJsonExpansion(&options, types.BoolValue(true))

		if options != &existing {
			t.Fatal("expected the existing options map to be reused")
		}
		if existing["filters:releases"] != "1.*" {
			t.Fatalf("existing option was overwritten: %#v", existing)
		}
		if existing["sentry:relay_automatic_json_expansion"] != true {
			t.Fatalf("option was not set: %#v", existing)
		}
	})
}

func mergeOptions(base map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(base)+len(extra))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range extra {
		merged[key] = value
	}
	return merged
}
