package tfutils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/jianyuan/go-utils/sliceutils"
)

func WithEnumInt64Attribute(base schema.Int64Attribute, choices []int64) schema.Int64Attribute {
	// Add a markdown description that lists the valid values
	if base.MarkdownDescription != "" {
		base.MarkdownDescription += " "
	}
	validValues := sliceutils.Map(func(v int64) string {
		return fmt.Sprintf("`%d`", v)
	}, choices)
	if len(validValues) > 1 {
		base.MarkdownDescription += "Valid values are: " + strings.Join(validValues[:len(validValues)-1], ", ") + ", and " + validValues[len(validValues)-1] + "."
	} else {
		base.MarkdownDescription += "Valid values are: " + validValues[0] + "."
	}

	// Add a validator that checks the value is one of the valid values
	base.Validators = append(base.Validators, int64validator.OneOf(choices...))

	return base
}
