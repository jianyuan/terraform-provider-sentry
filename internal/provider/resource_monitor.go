package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
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

func (r *MonitorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("config").AtName("schedule_crontab"),
			path.MatchRoot("config").AtName("schedule_interval"),
		),
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

	httpResp, err := r.apiClient.CreateMonitorWithResponse(ctx, data.Organization.ValueString(), monitorRequest)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	}

	if httpResp.StatusCode() != http.StatusCreated || httpResp.JSON201 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, data.Organization.ValueString(), *httpResp.JSON201)...)
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

	httpResp, err := r.apiClient.GetOrganizationMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, data.Organization.ValueString(), *httpResp.JSON200)...)
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

	httpResp, err := r.apiClient.UpdateOrganizationMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
		monitorRequest,
	)

	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, data.Organization.ValueString(), *httpResp.JSON200)...)
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

	httpResp, err := r.apiClient.DeleteOrganizationMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		return
	}

	if httpResp.StatusCode() != http.StatusAccepted {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateTwoPartId(ctx, "organization", req, resp)
}
