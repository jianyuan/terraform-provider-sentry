package sentrytypes

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestKnownStringPointer(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		value StringValueLike
		want  *string
	}{
		{name: "null", value: NewSlackChannelNull(), want: nil},
		{name: "unknown", value: NewSlackChannelUnknown(), want: nil},
		{name: "set", value: NewSlackChannelValue("#general"), want: new("#general")},
		{name: "base string", value: basetypes.NewStringValue("hello"), want: new("hello")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := KnownStringPointer(tt.value)
			if (got == nil) != (tt.want == nil) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			if got != nil && *got != *tt.want {
				t.Fatalf("got %q, want %q", *got, *tt.want)
			}
		})
	}
}
