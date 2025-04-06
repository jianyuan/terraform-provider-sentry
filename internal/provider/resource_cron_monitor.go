package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/jianyuan/terraform-provider-sentry/internal/api"
)

type cronMonitorResource struct {
	client *api.Client
}

func NewCronMonitorResource() resource.Resource {
	return &cronMonitorResource{}
}

func (r *cronMonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cron_monitor"
}

func (r *cronMonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization": schema.StringAttribute{
				Required: true,
			},
			"project": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"schedule": schema.StringAttribute{
				Required: true,
			},
			"schedule_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default: stringdefault.StaticString("crontab"),
			},
			"checkin_margin": schema.StringAttribute{
				Optional: true,
			},
			"max_runtime": schema.StringAttribute{
				Optional: true,
			},
			"timezone": schema.StringAttribute{
				Optional: true,
				Default: stringdefault.StaticString("UTC"),
			},
			"environment": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default: booldefault.StaticBool(true),
			},
		},
	}
}

func (r *cronMonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if client, ok := req.ProviderData.(*api.Client); ok {
		r.client = client
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *api.Client, got something else.",
		)
	}
}

type cronMonitorResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Organization    types.String `tfsdk:"organization"`
	Project         types.String `tfsdk:"project"`
	Name            types.String `tfsdk:"name"`
	Schedule        types.String `tfsdk:"schedule"`
	ScheduleType    types.String `tfsdk:"schedule_type"`
	CheckinMargin   types.String `tfsdk:"checkin_margin"`
	MaxRuntime      types.String `tfsdk:"max_runtime"`
	Timezone        types.String `tfsdk:"timezone"`
	Environment     types.String `tfsdk:"environment"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

func (r *cronMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cronMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.client.CreateCronMonitor(ctx, data.Organization.ValueString(), data.Project.ValueString(), api.CreateCronMonitorParams{
		Name:          data.Name.ValueString(),
		Schedule:      data.Schedule.ValueString(),
		ScheduleType:  data.ScheduleType.ValueString(),
		CheckinMargin: data.CheckinMargin.ValueString(),
		MaxRuntime:    data.MaxRuntime.ValueString(),
		Timezone:      data.Timezone.ValueString(),
		Environment:   data.Environment.ValueString(),
		Enabled:       data.Enabled.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to create cron monitor: "+err.Error())
		return
	}

	data.ID = types.StringValue(monitor.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cronMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cronMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) 
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.client.GetCronMonitor(ctx, data.Organization.ValueString(), data.Project.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read cron monitor: "+err.Error())
		return
	}

	data.Name = types.StringValue(monitor.Name)
	data.Schedule = types.StringValue(monitor.Schedule)
	data.ScheduleType = types.StringValue(monitor.ScheduleType)
	data.CheckinMargin = types.StringValue(monitor.CheckinMargin)
	data.MaxRuntime = types.StringValue(monitor.MaxRuntime)
	data.Timezone = types.StringValue(monitor.Timezone)
	data.Environment = types.StringValue(monitor.Environment)
	data.Enabled = types.BoolValue(monitor.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...) 
}

func (r *cronMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data cronMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateCronMonitor(ctx, data.Organization.ValueString(), data.Project.ValueString(), data.ID.ValueString(), api.UpdateCronMonitorParams{
		Name:          data.Name.ValueString(),
		Schedule:      data.Schedule.ValueString(),
		ScheduleType:  data.ScheduleType.ValueString(),
		CheckinMargin: data.CheckinMargin.ValueString(),
		MaxRuntime:    data.MaxRuntime.ValueString(),
		Timezone:      data.Timezone.ValueString(),
		Environment:   data.Environment.ValueString(),
		Enabled:       data.Enabled.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to update cron monitor: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...) 
}

func (r *cronMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data cronMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCronMonitor(ctx, data.Organization.ValueString(), data.Project.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to delete cron monitor: "+err.Error())
		return
	}
}

func (r *cronMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect format: org/project/id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected format: organization/project/id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
}