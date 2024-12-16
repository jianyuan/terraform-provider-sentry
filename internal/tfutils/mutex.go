package tfutils

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func WithMutuallyExclusiveValidator(attributes map[string]schema.SingleNestedAttribute) map[string]schema.Attribute {
	var names []string
	for name := range attributes {
		names = append(names, name)
	}

	conditionFor := func(name string) []validator.Object {
		var paths []path.Expression

		for _, thisName := range names {
			if thisName != name {
				paths = append(paths, path.MatchRelative().AtParent().AtName(thisName))
			}
		}

		return []validator.Object{objectvalidator.ConflictsWith(paths...)}
	}

	result := make(map[string]schema.Attribute, len(attributes))
	for name, attribute := range attributes {
		attribute.Validators = append(attribute.Validators, conditionFor(name)...)
		result[name] = attribute
	}

	return result
}
