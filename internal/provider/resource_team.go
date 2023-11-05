package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &TeamResource{}
var _ resource.ResourceWithImportState = &TeamResource{}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

type TeamResource struct {
	client *sentry.Client
}

type TeamResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	InternalId   types.String `tfsdk:"internal_id"`
	TeamId       types.String `tfsdk:"team_id"`
	HasAccess    types.Bool   `tfsdk:"has_access"`
	IsPending    types.Bool   `tfsdk:"is_pending"`
	IsMember     types.Bool   `tfsdk:"is_member"`
}

func (r *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Team resource.",

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
			"name": schema.StringAttribute{
				Description: "The name of the team.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slug": schema.StringAttribute{
				Description: "The optional slug for this team.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"internal_id": schema.StringAttribute{
				Description: "The internal ID for this team.",
				Computed:    true,
			},
			"team_id": schema.StringAttribute{
				Description:        "The internal ID of this team.",
				Computed:           true,
				DeprecationMessage: "Use `internal_id` instead.",
			},
			"has_access": schema.BoolAttribute{
				Description: "Whether the authenticated user has access to this team.",
				Computed:    true,
			},
			"is_pending": schema.BoolAttribute{
				Description: "Whether the team is pending.",
				Computed:    true,
			},
			"is_member": schema.BoolAttribute{
				Description: "Whether the authenticated user is a member of this team.",
				Computed:    true,
			},
		},
	}
}

func (r *TeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	team, _, err := r.client.Teams.Create(ctx, data.Organization.ValueString(), &sentry.CreateTeamParams{
		Name: data.Name.ValueStringPointer(),
		Slug: data.Slug.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team, got error: %s", err))
		return
	}

	data.Id = types.StringPointerValue(team.Slug)
	data.Name = types.StringPointerValue(team.Name)
	data.Slug = types.StringPointerValue(team.Slug)
	data.InternalId = types.StringPointerValue(team.ID)
	data.TeamId = types.StringPointerValue(team.ID)
	data.HasAccess = types.BoolPointerValue(team.HasAccess)
	data.IsPending = types.BoolPointerValue(team.IsPending)
	data.IsMember = types.BoolPointerValue(team.IsMember)

	tflog.Trace(ctx, "created a team")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	team, _, err := r.client.Teams.Get(ctx, data.Organization.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team, got error: %s", err))
		return
	}

	data.Id = types.StringPointerValue(team.Slug)
	data.Name = types.StringPointerValue(team.Name)
	data.Slug = types.StringPointerValue(team.Slug)
	data.InternalId = types.StringPointerValue(team.ID)
	data.TeamId = types.StringPointerValue(team.ID)
	data.HasAccess = types.BoolPointerValue(team.HasAccess)
	data.IsPending = types.BoolPointerValue(team.IsPending)
	data.IsMember = types.BoolPointerValue(team.IsMember)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	team, _, err := r.client.Teams.Update(ctx, data.Organization.ValueString(), data.Id.ValueString(), &sentry.UpdateTeamParams{
		Name: data.Name.ValueStringPointer(),
		Slug: data.Slug.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team, got error: %s", err))
		return
	}

	data.Id = types.StringPointerValue(team.Slug)
	data.Name = types.StringPointerValue(team.Name)
	data.Slug = types.StringPointerValue(team.Slug)
	data.InternalId = types.StringPointerValue(team.ID)
	data.TeamId = types.StringPointerValue(team.ID)
	data.HasAccess = types.BoolPointerValue(team.HasAccess)
	data.IsPending = types.BoolPointerValue(team.IsPending)
	data.IsMember = types.BoolPointerValue(team.IsMember)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Teams.Delete(ctx, data.Organization.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team, got error: %s", err))
		return
	}
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	org, id, err := splitTwoPartID(req.ID, "organization", "id")
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to import team, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), org,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)...)
}
