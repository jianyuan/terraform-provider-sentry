package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
