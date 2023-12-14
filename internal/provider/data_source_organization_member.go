package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &OrganizationMemberDataSource{}

func NewOrganizationMemberDataSource() datasource.DataSource {
	return &OrganizationMemberDataSource{}
}

type OrganizationMemberDataSource struct {
	client *sentry.Client
}

type OrganizationMemberDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Email        types.String `tfsdk:"email"`
	Role         types.String `tfsdk:"role"`
}

func (m *OrganizationMemberDataSourceModel) Fill(organization string, d sentry.OrganizationMember) error {
	m.Id = types.StringValue(d.ID)
	m.Organization = types.StringValue(organization)
	m.Email = types.StringValue(d.Email)
	m.Role = types.StringValue(d.OrganizationRole)

	return nil
}

func (d *OrganizationMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_member"
}

func (d *OrganizationMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve an organization member by email.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization.",
				Required:            true,
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
	}
}

func (d *OrganizationMemberDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sentry.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sentry.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *OrganizationMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationMemberDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var foundMember *sentry.OrganizationMember
	params := &sentry.ListCursorParams{}

out:
	for {
		members, apiResp, err := d.client.OrganizationMembers.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list organization members, got error: %s", err))
			return
		}

		for _, member := range members {
			if member.Email == data.Email.ValueString() {
				foundMember = member
				break out
			}
		}

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	if foundMember == nil {
		resp.Diagnostics.AddError("Not Found", "No matching organization member found")
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *foundMember); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill organization member, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
