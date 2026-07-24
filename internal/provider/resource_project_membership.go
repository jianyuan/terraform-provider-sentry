package provider

import (
	"context"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type ProjectMembershipResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	ProjectId    types.String `tfsdk:"project_id"`
	TeamId       types.String `tfsdk:"team_id"`
}

func (m *ProjectMembershipResourceModel) Fill(organization, projectId, teamId string) {
	m.Id = types.StringValue(tfutils.BuildThreePartId(organization, projectId, teamId))
	m.Organization = types.StringValue(organization)
	m.ProjectId = types.StringValue(projectId)
	m.TeamId = types.StringValue(teamId)
}

func NewProjectMembershipResource() resource.Resource {
	return &ProjectMembershipResource{}
}

var _ resource.Resource = &ProjectMembershipResource{}
var _ resource.ResourceWithConfigure = &ProjectMembershipResource{}
var _ resource.ResourceWithImportState = &ProjectMembershipResource{}

type ProjectMembershipResource struct {
	baseResource
}

func (r *ProjectMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_membership"
}

func (r *ProjectMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Linking a Sentry team to a Sentry project.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The slug of the Sentry project.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "The slug of the Sentry team.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ProjectMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AddTeamToProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.ProjectId.ValueString(),
		data.TeamId.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("add team to project", err))
		return
	} else if httpResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("add team to project", httpResp.StatusCode(), httpResp.Body))
		return
	}

	data.Fill(
		data.Organization.ValueString(),
		data.ProjectId.ValueString(),
		data.TeamId.ValueString(),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectMembershipResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetOrganizationProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.ProjectId.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project membership"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	project := httpResp.JSON200
	teamFound := slices.ContainsFunc(project.Teams, func(team apiclient.Team) bool {
		return team.Slug == data.TeamId.ValueString()
	})

	if !teamFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project membership"))
		resp.State.RemoveResource(ctx)
		return
	}

	data.Fill(
		data.Organization.ValueString(),
		data.ProjectId.ValueString(),
		data.TeamId.ValueString(),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Project membership cannot be updated. Changes to project_id or team_id require resource replacement.",
	)
}

func (r *ProjectMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectMembershipResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RemoveTeamFromProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.ProjectId.ValueString(),
		data.TeamId.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *ProjectMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, projectId, teamId, err := tfutils.SplitThreePartId(req.ID, "organization", "project-slug", "team-slug")
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewImportError(err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("project_id"), projectId,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("team_id"), teamId,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}
