package organization

import (
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
)

func Schema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "Sentry Organization",
		},
		Resource: superschema.SchemaDetails{},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The unique URL slug for this organization.",
					Computed:            true,
				},
			},
			"slug": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The unique URL slug for this organization.",
					Required:            true,
				},
			},
			"internal_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The internal ID for this organization.",
					Computed:            true,
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The human readable name for this organization.",
					Computed:            true,
				},
			},
		},
	}
}
