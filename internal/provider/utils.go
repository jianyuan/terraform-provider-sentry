package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/oapi-codegen/nullable"
)

func ResourceIdAttribute() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func ResourceOrganizationAttribute() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: "The organization of this resource.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func ResourceProjectAttribute() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: "The project of this resource.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func DataSourceOrganizationAttribute() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: "The organization the resource belongs to.",
		Required:            true,
	}
}

func DataSourceProjectAttribute() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: "The project the resource belongs to.",
		Required:            true,
	}
}

func nullableFromPtr[T any](v *T) nullable.Nullable[T] {
	if v == nil {
		return nullable.NewNullNullable[T]()
	}
	return nullable.NewNullableWithValue(*v)
}
