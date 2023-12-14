package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &ProjectSpikeProtectionResource{}
var _ resource.ResourceWithImportState = &ProjectSpikeProtectionResource{}

func NewProjectSpikeProtectionResource() resource.Resource {
	return &ProjectSpikeProtectionResource{}
}

type ProjectSpikeProtectionResource struct {
	client *sentry.Client
}

type ProjectSpikeProtectionResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
	Enabled      types.Bool   `tfsdk:"enabled"`
}

func (data *ProjectSpikeProtectionResourceModel) Fill(organization string, project sentry.Project) error {
	data.Id = types.StringValue(buildTwoPartID(organization, project.Slug))
	data.Organization = types.StringPointerValue(project.Organization.Slug)
	data.Project = types.StringValue(project.Slug)
	if disabled, ok := project.Options["quotas:spike-protection-disabled"].(bool); ok {
		data.Enabled = types.BoolValue(!disabled)
	}

	return nil
}

func (r *ProjectSpikeProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_spike_protection"
}

func (r *ProjectSpikeProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project Spike Protection resource. This resource is used to create and manage spike protection for a project.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Description: "The slug of the organization the project belongs to.",
				Required:    true,
			},
			"project": schema.StringAttribute{
				Description: "The slug of the project to create the filter for.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Toggle the browser-extensions, localhost, filtered-transaction, or web-crawlers filter on or off.",
				Required:    true,
			},
		},
	}
}

func (r *ProjectSpikeProtectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectSpikeProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Enabled.ValueBool() {
		_, err := r.client.SpikeProtections.Enable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: []string{data.Project.ValueString()},
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error enabling spike protection: %s", err.Error()))
			return
		}
	} else {
		_, err := r.client.SpikeProtections.Disable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: []string{data.Project.ValueString()},
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error disabling spike protection: %s", err.Error()))
			return
		}
	}

	data.Id = types.StringValue(buildTwoPartID(data.Organization.ValueString(), data.Project.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectSpikeProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	project, apiResp, err := r.client.Projects.Get(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Project not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading project: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *project); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling project spike protection: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectSpikeProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Enabled.ValueBool() {
		_, err := r.client.SpikeProtections.Enable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: []string{data.Project.ValueString()},
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error enabling spike protection: %s", err.Error()))
			return
		}
	} else {
		_, err := r.client.SpikeProtections.Disable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: []string{data.Project.ValueString()},
			},
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error disabling spike protection: %s", err.Error()))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectSpikeProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.SpikeProtections.Disable(
		ctx,
		data.Organization.ValueString(),
		&sentry.SpikeProtectionParams{
			Projects: []string{data.Project.ValueString()},
		},
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error disabling spike protection: %s", err.Error()))
		return
	}
}

func (r *ProjectSpikeProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, err := splitTwoPartID(req.ID, "organization", "project-slug")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("project"), project,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}
