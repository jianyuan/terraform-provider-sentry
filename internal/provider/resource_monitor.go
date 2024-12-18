package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithConfigure = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

type MonitorResource struct {
	baseResource
}

func (r *MonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Return a client monitor bound to a project.",
		Attributes:          MonitorResourceModel{}.Attributes(),
	}
}

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitorRequest, monitorRequestDiags := data.ToMonitorRequest(ctx)
	resp.Diagnostics.Append(monitorRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.apiClient.CreateMonitorWithResponse(ctx, data.Organization.ValueString(), monitorRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create error: %s", err.Error()))
		return
	}

	if response.JSON201 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create error: %s", response.HTTPResponse.Status))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, data.Organization.ValueString(), *response.JSON201)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.apiClient.GetOrganizationMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	if response.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", "Not found")
		// resp.State.RemoveResource(ctx)
		return
	}

	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", response.HTTPResponse.Status))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, data.Organization.ValueString(), *response.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitorRequest, monitorRequestDiags := data.ToMonitorRequest(ctx)
	resp.Diagnostics.Append(monitorRequestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.apiClient.UpdateOrganizationMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
		monitorRequest,
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Update error: %s", err.Error()))
		return
	}

	if response.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", "Not found")
		// resp.State.RemoveResource(ctx)
		return
	}

	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", response.HTTPResponse.Status))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, data.Organization.ValueString(), *response.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.apiClient.DeleteOrganizationMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Delete error: %s", err.Error()))
		return
	}

	if apiResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.AddWarning("Monitor not found", "Monitor may have been deleted already")
		return
	}
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, id, err := tfutils.SplitTwoPartId(req.ID, "organization", "monitor-id")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)...)
}
