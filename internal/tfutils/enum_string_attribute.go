package tfutils

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/jianyuan/go-utils/sliceutils"
)

func WithEnumStringAttribute(base schema.StringAttribute, choices []string) schema.StringAttribute {
	// Add a markdown description that lists the valid values
	if base.MarkdownDescription != "" {
		base.MarkdownDescription += " "
	}
	validValues := sliceutils.Map(func(v string) string {
		return "`" + v + "`"
	}, choices)
	if len(validValues) > 1 {
		base.MarkdownDescription += "Valid values are: " + strings.Join(validValues[:len(validValues)-1], ", ") + ", and " + validValues[len(validValues)-1] + "."
	} else {
		base.MarkdownDescription += "Valid values are: " + validValues[0] + "."
	}

	// Add a validator that checks the value is one of the valid values
	base.Validators = append(base.Validators, stringvalidator.OneOf(choices...))

	return base
}
