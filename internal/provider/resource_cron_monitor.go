package provider

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
    "github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

type CronMonitorResourceModel struct {
    ID           types.String `tfsdk:"id"`
    Organization types.String `tfsdk:"organization"`
    Project      types.String `tfsdk:"project"`
    Name         types.String `tfsdk:"name"`
    Schedule     types.String `tfsdk:"schedule"`
    ScheduleType types.String `tfsdk:"schedule_type"`
    CheckinMargin types.Int64  `tfsdk:"checkin_margin"`
    MaxRuntime   types.Int64  `tfsdk:"max_runtime"`
    Timezone     types.String `tfsdk:"timezone"`
    Environment  types.String `tfsdk:"environment"`
    Enabled      types.Bool   `tfsdk:"enabled"`
}

type CronMonitorResource struct {
    baseResource
}

func NewCronMonitorResource() resource.Resource {
    return &CronMonitorResource{}
}

func (r *CronMonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cron_monitor"
}

func (r *CronMonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
                Default: stringdefault.StaticString("cron_job"),
            },
            "checkin_margin": schema.Int64Attribute{
                Optional: true,
            },
            "max_runtime": schema.Int64Attribute{
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

func (r *CronMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data CronMonitorResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    createBody := map[string]interface{}{
        "name": data.Name.ValueString(),
        "type": "cron_job",
        "config": map[string]interface{}{
            "schedule":       data.Schedule.ValueString(),
            "timezone":       data.Timezone.ValueString(),
            "checkin_margin": data.CheckinMargin.ValueInt64Pointer(),
            "max_runtime":    data.MaxRuntime.ValueInt64Pointer(),
        },
        "environment": data.Environment.ValueString(),
        "enabled":     data.Enabled.ValueBool(),
    }

    body, err := json.Marshal(createBody)
    if err != nil {
        resp.Diagnostics.AddError("Failed to marshal request body", err.Error())
        return
    }

    httpResp, err := r.apiClient.CreateMonitorWithBodyWithResponse(
        ctx,
        data.Organization.ValueString(),
        data.Project.ValueString(),
        "application/json",
        strings.NewReader(string(body)),
    )
    if err != nil {
        resp.Diagnostics.Append(diagutils.NewClientError("create", err))
        return
    } else if httpResp.StatusCode() != http.StatusCreated {
        resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
        return
    }

    var monitor apiclient.Monitor
    if err := json.Unmarshal(httpResp.Body, &monitor); err != nil {
        resp.Diagnostics.AddError("Failed to parse response body", err.Error())
        return
    }

    data.ID = types.StringValue(*monitor.Id)
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data CronMonitorResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    httpResp, err := r.apiClient.GetMonitorWithResponse(
        ctx,
        data.Organization.ValueString(),
        data.Project.ValueString(),
        data.ID.ValueString(),
    )
    if err != nil {
        resp.Diagnostics.Append(diagutils.NewClientError("read", err))
        return
    } else if httpResp.StatusCode() == http.StatusNotFound {
        resp.Diagnostics.Append(diagutils.NewNotFoundError("cron monitor"))
        resp.State.RemoveResource(ctx)
        return
    } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
        return
    }

    var monitor apiclient.Monitor
    if err := json.Unmarshal(httpResp.Body, &monitor); err != nil {
        resp.Diagnostics.AddError("Failed to parse response body", err.Error())
        return
    }

    data.Name = types.StringValue(*monitor.Name)
    if monitor.Config != nil {
        config := *monitor.Config
        if schedule, ok := config["schedule"].(string); ok {
            data.Schedule = types.StringValue(schedule)
        }
        if timezone, ok := config["timezone"].(string); ok {
            data.Timezone = types.StringValue(timezone)
        }
        if checkinMargin, ok := config["checkin_margin"].(float64); ok {
            data.CheckinMargin = types.Int64Value(int64(checkinMargin))
        }
        if maxRuntime, ok := config["max_runtime"].(float64); ok {
            data.MaxRuntime = types.Int64Value(int64(maxRuntime))
        }
    }
    data.ScheduleType = types.StringValue(*monitor.Type)

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data CronMonitorResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    updateBody := map[string]interface{}{
        "name": data.Name.ValueString(),
        "type": "cron_job",
        "config": map[string]interface{}{
            "schedule":       data.Schedule.ValueString(),
            "timezone":       data.Timezone.ValueString(),
            "checkin_margin": data.CheckinMargin.ValueInt64Pointer(),
            "max_runtime":    data.MaxRuntime.ValueInt64Pointer(),
        },
        "environment": data.Environment.ValueString(),
        "enabled":     data.Enabled.ValueBool(),
    }

    body, err := json.Marshal(updateBody)
    if err != nil {
        resp.Diagnostics.AddError("Failed to marshal request body", err.Error())
        return
    }

    httpResp, err := r.apiClient.UpdateMonitorWithBodyWithResponse(
        ctx,
        data.Organization.ValueString(),
        data.Project.ValueString(),
        data.ID.ValueString(),
        "application/json",
        strings.NewReader(string(body)),
    )
    if err != nil {
        resp.Diagnostics.Append(diagutils.NewClientError("update", err))
        return
    } else if httpResp.StatusCode() != http.StatusOK {
        resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data CronMonitorResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    httpResp, err := r.apiClient.DeleteMonitorWithResponse(
        ctx,
        data.Organization.ValueString(),
        data.Project.ValueString(),
        data.ID.ValueString(),
    )
    if err != nil {
        resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
        return
    } else if httpResp.StatusCode() == http.StatusNotFound {
        return
    } else if httpResp.StatusCode() != http.StatusNoContent {
        resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
        return
    }
}

func (r *CronMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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