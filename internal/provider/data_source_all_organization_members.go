package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &AllOrganizationMembersDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationMemberDataSource{}

func NewAllOrganizationMembersDataSource() datasource.DataSource {
	return &AllOrganizationMembersDataSource{}
}

type AllOrganizationMembersDataSource struct {
	baseDataSource
}

type AllOrganizationMembersDataSourceMemberModel struct {
	Id    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func (m *AllOrganizationMembersDataSourceMemberModel) Fill(organization string, member sentry.OrganizationMember) error {
	m.Id = types.StringValue(member.ID)
	m.Email = types.StringValue(member.Email)
	m.Role = types.StringValue(member.OrgRole)

	return nil
}

type AllOrganizationMembersDataSourceModel struct {
	Organization types.String                                  `tfsdk:"organization"`
	Members      []AllOrganizationMembersDataSourceMemberModel `tfsdk:"members"`
}

func (m *AllOrganizationMembersDataSourceModel) Fill(organization string, members []sentry.OrganizationMember) error {
	m.Organization = types.StringValue(organization)

	for _, member := range members {
		mm := AllOrganizationMembersDataSourceMemberModel{}
		if err := mm.Fill(organization, member); err != nil {
			return err
		}
		m.Members = append(m.Members, mm)
	}

	return nil
}

func (d *AllOrganizationMembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_organization_members"
}

func (d *AllOrganizationMembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve all organization members.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization.",
				Required:            true,
			},
			"members": schema.SetNestedAttribute{
				MarkdownDescription: "The list of members.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of of the organization member.",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email of the organization member.",
							Required:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "This is the role of the organization member.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AllOrganizationMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllOrganizationMembersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allMembers []sentry.OrganizationMember
	params := &sentry.ListCursorParams{}

	for {
		members, apiResp, err := d.client.OrganizationMembers.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list organization members, got error: %s", err))
			return
		}

		for _, member := range members {
			allMembers = append(allMembers, *member)
		}

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	if err := data.Fill(data.Organization.ValueString(), allMembers); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
