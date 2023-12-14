package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &TeamMemberResource{}
var _ resource.ResourceWithImportState = &TeamMemberResource{}

func NewTeamMemberResource() resource.Resource {
	return &TeamMemberResource{}
}

type TeamMemberResource struct {
	client *sentry.Client

	roleMu sync.Mutex
}

type TeamMemberResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	MemberId     types.String `tfsdk:"member_id"`
	Team         types.String `tfsdk:"team"`
	Role         types.String `tfsdk:"role"`
}

func (data *TeamMemberResourceModel) Fill(organization string, team string, memberId string, role string) error {
	data.Id = types.StringValue(buildThreePartID(organization, team, memberId))
	data.Organization = types.StringValue(organization)
	data.MemberId = types.StringValue(memberId)
	data.Team = types.StringValue(team)
	data.Role = types.StringValue(role)

	return nil
}

func (r *TeamMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_member"
}

func (r *TeamMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Team Member resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Description: "The slug of the organization the team should be created for.",
				Required:    true,
			},
			"member_id": schema.StringAttribute{
				Description: "The ID of the member to add to the team.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team": schema.StringAttribute{
				Description: "The slug of the team to add the member to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Description: "The role of the member in the team. When not set, resolve to the minimum team role given by this member's organization role.",
				Optional:    true,
			},
		},
	}
}

func (r *TeamMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sentry.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sentry.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TeamMemberResource) readRole(ctx context.Context, organization string, memberId string, team string) (*string, error) {
	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	member, _, err := r.client.OrganizationMembers.Get(ctx, organization, memberId)
	if err != nil {
		return nil, fmt.Errorf("unable to read organization member, got error: %s", err)
	}

	teamRole := r.readTeamRole(member.TeamRoles, team)
	if teamRole == nil {
		return nil, fmt.Errorf("unable to find team member")
	}

	return teamRole, nil
}

func (r *TeamMemberResource) readTeamRole(teamRoles []sentry.TeamRole, team string) *string {
	for _, teamRole := range teamRoles {
		if teamRole.TeamSlug == team {
			return &teamRole.Role
		}
	}

	return nil
}

func (r *TeamMemberResource) updateRole(ctx context.Context, organization string, memberId string, team string, role string) (*string, error) {
	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	orgMember, _, err := r.client.OrganizationMembers.Get(ctx, organization, memberId)
	if err != nil {
		return nil, fmt.Errorf("unable to read organization member, got error: %s", err)
	}

	teamRoles := make([]sentry.TeamRole, 0, len(orgMember.TeamRoles))
	for _, teamRole := range orgMember.TeamRoles {
		if teamRole.TeamSlug == team {
			teamRole.Role = role
		}
		teamRoles = append(teamRoles, teamRole)
	}

	orgMember, _, err = r.client.OrganizationMembers.Update(
		ctx,
		organization,
		memberId,
		&sentry.UpdateOrganizationMemberParams{
			OrganizationRole: orgMember.OrganizationRole,
			TeamRoles:        teamRoles,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update organization member's team role, got error: %s", err)
	}

	for _, teamRole := range orgMember.TeamRoles {
		if teamRole.TeamSlug == team {
			return &teamRole.Role, nil
		}
	}

	return nil, fmt.Errorf("unable to find team member")
}

func (r *TeamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamMemberResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	member, _, err := r.client.TeamMembers.Create(
		ctx,
		data.Organization.ValueString(),
		data.MemberId.ValueString(),
		data.Team.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add member to team, got error: %s", err))
		return
	}

	var teamRole *string
	if data.Role.IsNull() {
		teamRole = member.TeamRole
	} else {
		teamRole, err = r.updateRole(ctx, data.Organization.ValueString(), data.MemberId.ValueString(), data.Team.ValueString(), data.Role.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team member, got error: %s", err))
			return
		}
	}

	if teamRole == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to find team member")
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Team.ValueString(), data.MemberId.ValueString(), *teamRole); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill team member, got error: %s", err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamMemberResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	role, err := r.readRole(ctx, data.Organization.ValueString(), data.MemberId.ValueString(), data.Team.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Team.ValueString(), data.MemberId.ValueString(), *role); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill team member, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state TeamMemberResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the role if it has changed
	if !plan.Role.Equal(state.Role) {
		role, err := r.updateRole(ctx, plan.Organization.ValueString(), plan.MemberId.ValueString(), plan.Team.ValueString(), plan.Role.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}

		if err := state.Fill(plan.Organization.ValueString(), plan.Team.ValueString(), plan.MemberId.ValueString(), *role); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fill team member, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TeamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamMemberResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.client.TeamMembers.Delete(
		ctx,
		data.Organization.ValueString(),
		data.MemberId.ValueString(),
		data.Team.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team member, got error: %s", err))
		return
	}
}

func (r *TeamMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, team, memberId, err := splitThreePartID(req.ID, "organization", "team-slug", "member-id")
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to import team, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("team"), team,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("member_id"), memberId,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}
