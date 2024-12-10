package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

type ProjectSpikeProtectionResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
	Enabled      types.Bool   `tfsdk:"enabled"`
}

func (data *ProjectSpikeProtectionResourceModel) Fill(project apiclient.Project) error {
	data.Id = types.StringValue(buildTwoPartID(project.Organization.Slug, project.Slug))
	data.Organization = types.StringValue(project.Organization.Slug)
	data.Project = types.StringValue(project.Slug)
	if disabled, ok := project.Options["quotas:spike-protection-disabled"].(bool); ok {
		data.Enabled = types.BoolValue(!disabled)
	} else {
		data.Enabled = types.BoolNull()
	}

	return nil
}

var _ resource.Resource = &ProjectSpikeProtectionResource{}
var _ resource.ResourceWithConfigure = &ProjectSpikeProtectionResource{}
var _ resource.ResourceWithImportState = &ProjectSpikeProtectionResource{}

func NewProjectSpikeProtectionResource() resource.Resource {
	return &ProjectSpikeProtectionResource{}
}

type ProjectSpikeProtectionResource struct {
	baseResource
}

func (r *ProjectSpikeProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_spike_protection"
}

func (r *ProjectSpikeProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project Spike Protection resource. This resource is used to create and manage spike protection for a project.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Toggle the browser-extensions, localhost, filtered-transaction, or web-crawlers filter on or off.",
				Required:            true,
			},
		},
	}
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
			resp.Diagnostics.Append(diagutils.NewClientError("enable", err))
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
			resp.Diagnostics.Append(diagutils.NewClientError("disable", err))
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

	httpResp, err := r.apiClient.GetOrganizationProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
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
			resp.Diagnostics.Append(diagutils.NewClientError("enable", err))
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
			resp.Diagnostics.Append(diagutils.NewClientError("disable", err))
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
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	}
}

func (r *ProjectSpikeProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, err := splitTwoPartID(req.ID, "organization", "project-slug")
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
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
