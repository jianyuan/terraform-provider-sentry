package provider

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &TeamMemberResource{}
var _ resource.ResourceWithConfigure = &TeamMemberResource{}
var _ resource.ResourceWithImportState = &TeamMemberResource{}

func NewTeamMemberResource() resource.Resource {
	return &TeamMemberResource{}
}

type TeamMemberResource struct {
	baseResource

	roleMu sync.Mutex
}

type TeamMemberResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Organization  types.String `tfsdk:"organization"`
	MemberId      types.String `tfsdk:"member_id"`
	Team          types.String `tfsdk:"team"`
	Role          types.String `tfsdk:"role"`
	EffectiveRole types.String `tfsdk:"effective_role"`
}

func (data *TeamMemberResourceModel) Fill(organization string, team string, memberId string, role *string, effectiveRole string) error {
	data.Id = types.StringValue(buildThreePartID(organization, team, memberId))
	data.Organization = types.StringValue(organization)
	data.MemberId = types.StringValue(memberId)
	data.Team = types.StringValue(team)
	data.Role = types.StringPointerValue(role)
	data.EffectiveRole = types.StringValue(effectiveRole)

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
				Validators: []validator.String{
					stringvalidator.OneOf("contributor", "admin"),
				},
			},
			"effective_role": schema.StringAttribute{
				Description: "The effective role of the member in the team. This represents the highest role, determined by comparing the lower role assigned by the member's organizational role with the role assigned by the member's team role.",
				Computed:    true,
			},
		},
	}
}

func getEffectiveOrgRole(memberOrgRoles []string, orgRoleList []sentry.OrganizationRoleListItem) *sentry.OrganizationRoleListItem {
	orgRoleMap := make(map[string]struct {
		index int
		role  sentry.OrganizationRoleListItem
	}, len(orgRoleList))
	for i, role := range orgRoleList {
		orgRoleMap[role.ID] = struct {
			index int
			role  sentry.OrganizationRoleListItem
		}{
			index: i,
			role:  role,
		}
	}
	memberOrgRolesCopy := make([]string, len(memberOrgRoles))
	copy(memberOrgRolesCopy, memberOrgRoles)

	slices.SortFunc(memberOrgRolesCopy, func(i, j string) int {
		return cmp.Compare(orgRoleMap[j].index, orgRoleMap[i].index)
	})

	if len(memberOrgRolesCopy) > 0 {
		if orgRoleMap, ok := orgRoleMap[memberOrgRolesCopy[0]]; ok {
			return &orgRoleMap.role
		}
	}

	return nil
}

func hasOrgRoleOverwrite(orgRole *sentry.OrganizationRoleListItem, orgRoleList []sentry.OrganizationRoleListItem, teamRoleList []sentry.TeamRoleListItem) bool {
	if orgRole == nil {
		return false
	}

	teamRoleIndex := slices.IndexFunc(teamRoleList, func(teamRole sentry.TeamRoleListItem) bool {
		return teamRole.ID == orgRole.MinimumTeamRole
	})

	return teamRoleIndex > 0
}

// Adapted from https://github.com/getsentry/sentry/blob/23.12.1/static/app/components/teamRoleSelect.tsx#L30-L69
func (r *TeamMemberResource) getEffectiveTeamRole(ctx context.Context, organization string, memberId string, teamSlug string) (*string, error) {
	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	org, _, err := r.client.Organizations.Get(ctx, organization)
	if err != nil {
		return nil, fmt.Errorf("unable to read organization, got error: %s", err)
	}

	team, _, err := r.client.Teams.Get(ctx, organization, teamSlug)
	if err != nil {
		return nil, fmt.Errorf("unable to read team, got error: %s", err)
	}

	member, _, err := r.client.OrganizationMembers.Get(ctx, organization, memberId)
	if err != nil {
		return nil, fmt.Errorf("unable to read organization member, got error: %s", err)
	}

	possibleOrgRoles := []string{member.OrgRole}
	if team.OrgRole != nil {
		possibleOrgRoles = append(possibleOrgRoles, sentry.StringValue(team.OrgRole))
	}

	effectiveOrgRole := getEffectiveOrgRole(possibleOrgRoles, org.OrgRoleList)

	if hasOrgRoleOverwrite(effectiveOrgRole, org.OrgRoleList, org.TeamRoleList) {
		teamRoleIndex := slices.IndexFunc(org.TeamRoleList, func(teamRole sentry.TeamRoleListItem) bool {
			return teamRole.ID == effectiveOrgRole.MinimumTeamRole
		})
		if teamRoleIndex != -1 {
			teamRole := org.TeamRoleList[teamRoleIndex]
			return &teamRole.ID, nil
		}
	}

	teamRoleIndex := slices.IndexFunc(member.TeamRoles, func(teamRole sentry.TeamRole) bool {
		return teamRole.TeamSlug == teamSlug
	})
	if teamRoleIndex != -1 {
		teamRole := member.TeamRoles[teamRoleIndex]
		if teamRole.Role != nil {
			return teamRole.Role, nil
		}
	}

	teamRole := member.TeamRoleList[0]
	return &teamRole.ID, nil
}

func (r *TeamMemberResource) updateRole(ctx context.Context, organization string, memberId string, team string, role string) (*string, error) {
	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	member, _, err := r.client.TeamMembers.Update(ctx, organization, memberId, team, &sentry.UpdateTeamMemberParams{
		TeamRole: sentry.String(role),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to read organization member, got error: %s", err)
	}

	if !sentry.BoolValue(member.IsActive) {
		return nil, fmt.Errorf("team member is not active")
	}

	return member.TeamRole, nil
}

func (r *TeamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamMemberResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.client.TeamMembers.Create(
		ctx,
		data.Organization.ValueString(),
		data.MemberId.ValueString(),
		data.Team.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add member to team, got error: %s", err))
		return
	}

	if !data.Role.IsNull() {
		_, err = r.updateRole(ctx, data.Organization.ValueString(), data.MemberId.ValueString(), data.Team.ValueString(), data.Role.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team member, got error: %s", err))
			return
		}
	}

	effectiveRole, err := r.getEffectiveTeamRole(ctx, data.Organization.ValueString(), data.MemberId.ValueString(), data.Team.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team member role, got error: %s", err))
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(
		data.Organization.ValueString(),
		data.Team.ValueString(),
		data.MemberId.ValueString(),
		data.Role.ValueStringPointer(),
		*effectiveRole,
	); err != nil {
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

	effectiveRole, err := r.getEffectiveTeamRole(ctx, data.Organization.ValueString(), data.MemberId.ValueString(), data.Team.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404 The requested resource does not exist") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Client Error", err.Error())
		}
		return
	}

	if err := data.Fill(
		data.Organization.ValueString(),
		data.Team.ValueString(),
		data.MemberId.ValueString(),
		data.Role.ValueStringPointer(),
		*effectiveRole,
	); err != nil {
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
		_, err := r.updateRole(ctx, plan.Organization.ValueString(), plan.MemberId.ValueString(), plan.Team.ValueString(), plan.Role.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}

		effectiveRole, err := r.getEffectiveTeamRole(ctx, plan.Organization.ValueString(), plan.MemberId.ValueString(), plan.Team.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			resp.State.RemoveResource(ctx)
			return
		}

		if err := state.Fill(
			plan.Organization.ValueString(),
			plan.Team.ValueString(),
			plan.MemberId.ValueString(),
			plan.Role.ValueStringPointer(),
			*effectiveRole,
		); err != nil {
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
