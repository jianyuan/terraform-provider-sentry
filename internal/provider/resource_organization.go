package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

type OrganizationModel struct {
	Id   types.String `tfsdk:"id"`
	Slug types.String `tfsdk:"slug"`
	Name types.String `tfsdk:"name"`
}

func (m *OrganizationModel) Fill(org sentry.Organization) error {
	m.Id = types.StringPointerValue(org.ID)
	m.Slug = types.StringPointerValue(org.Slug)
	m.Name = types.StringPointerValue(org.Name)

	return nil
}
