package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

const (
	monitorScheduleTypeCrontab  = "crontab"
	monitorScheduleTypeInterval = "interval"
)

type MonitorIntervalModel struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (m MonitorIntervalModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Int64Type,
		"unit":  types.StringType,
	}
}

type MonitorResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Organization          types.String `tfsdk:"organization"`
	Project               types.String `tfsdk:"project"`
	Name                  types.String `tfsdk:"name"`
	Slug                  types.String `tfsdk:"slug"`
	Owner                 types.String `tfsdk:"owner"`
	ScheduleCrontab       types.String `tfsdk:"schedule_crontab"`
	Timezone              types.String `tfsdk:"timezone"`
	ScheduleInterval      types.Object `tfsdk:"schedule_interval"`
	CheckinMargin         types.Int64  `tfsdk:"checkin_margin"`
	MaxRuntime            types.Int64  `tfsdk:"max_runtime"`
	FailureIssueThreshold types.Int64  `tfsdk:"failure_issue_threshold"`
	RecoveryThreshold     types.Int64  `tfsdk:"recovery_threshold"`
}

func (m *MonitorResourceModel) Fill(ctx context.Context, monitor apiclient.MonitorResponse) (diags diag.Diagnostics) {
	m.Id = types.StringValue(monitor.Id)
	m.Name = types.StringValue(monitor.Name)
	m.Slug = types.StringValue(monitor.Slug)

	if monitor.Owner != nil {
		m.Owner = types.StringValue(fmt.Sprintf("%s:%s", monitor.Owner.Type, monitor.Owner.Id))
	} else {
		m.Owner = types.StringNull()
	}

	scheduleType := string(monitor.Config.ScheduleType)
	if scheduleType == "" {
		scheduleType = monitorScheduleTypeCrontab
	}
	m.Timezone = nullableToString(monitor.Config.Timezone, "timezone", &diags)

	switch scheduleType {
	case monitorScheduleTypeCrontab:
		schedule, err := monitor.Config.Schedule.AsMonitorScheduleInput0()
		if err != nil {
			diags.Append(diagutils.NewFillError(fmt.Errorf("invalid schedule type for crontab: %w", err)))
			return
		}

		m.ScheduleCrontab = types.StringValue(schedule)
		m.ScheduleInterval = types.ObjectNull(MonitorIntervalModel{}.AttributeTypes())
	case monitorScheduleTypeInterval:
		schedule, err := monitor.Config.Schedule.AsMonitorScheduleInput1()
		if err != nil || len(schedule) != 2 {
			diags.Append(diagutils.NewFillError(fmt.Errorf("invalid schedule type for interval")))
			return
		}

		value, err := schedule[0].AsMonitorScheduleInput10()
		if err != nil {
			diags.Append(diagutils.NewFillError(fmt.Errorf("invalid interval value: %w", err)))
			return
		}

		unit, err := schedule[1].AsMonitorScheduleInput11()
		if err != nil {
			diags.Append(diagutils.NewFillError(fmt.Errorf("invalid interval unit: %w", err)))
			return
		}

		interval := MonitorIntervalModel{
			Value: types.Int64Value(value),
			Unit:  types.StringValue(string(unit)),
		}
		m.ScheduleInterval = tfutils.MergeDiagnostics(types.ObjectValueFrom(ctx, interval.AttributeTypes(), interval))(&diags)
		m.ScheduleCrontab = types.StringNull()
	default:
		diags.Append(diagutils.NewFillError(fmt.Errorf("unsupported schedule_type %q", scheduleType)))
		return
	}

	m.CheckinMargin = nullableToInt64(monitor.Config.CheckinMargin, "checkin_margin", &diags)
	m.MaxRuntime = nullableToInt64(monitor.Config.MaxRuntime, "max_runtime", &diags)
	m.FailureIssueThreshold = nullableToInt64(monitor.Config.FailureIssueThreshold, "failure_issue_threshold", &diags)
	m.RecoveryThreshold = nullableToInt64(monitor.Config.RecoveryThreshold, "recovery_threshold", &diags)

	return
}

func (m *MonitorResourceModel) buildCreateRequest(ctx context.Context) (apiclient.MonitorRequest, diag.Diagnostics) {
	return m.buildRequest(ctx)
}

func (m *MonitorResourceModel) buildUpdateRequest(ctx context.Context) (apiclient.MonitorRequest, diag.Diagnostics) {
	return m.buildRequest(ctx)
}

func (m *MonitorResourceModel) buildRequest(ctx context.Context) (apiclient.MonitorRequest, diag.Diagnostics) {
	config, diags := m.buildRequestConfig(ctx)
	if diags.HasError() {
		return apiclient.MonitorRequest{}, diags
	}

	body := apiclient.MonitorRequest{
		Config:  config,
		Name:    m.Name.ValueString(),
		Project: m.Project.ValueString(),
	}

	if !m.Slug.IsNull() && !m.Slug.IsUnknown() {
		body.Slug = m.Slug.ValueStringPointer()
	}

	if !m.Owner.IsUnknown() {
		if m.Owner.IsNull() {
			body.Owner.SetNull()
		} else {
			body.Owner.Set(m.Owner.ValueString())
		}
	}

	return body, diags
}

func (m *MonitorResourceModel) buildRequestConfig(ctx context.Context) (apiclient.MonitorConfigRequest, diag.Diagnostics) {
	schedule, scheduleType, diags := m.buildScheduleInput(ctx)
	if diags.HasError() {
		return apiclient.MonitorConfigRequest{}, diags
	}

	configType := apiclient.MonitorConfigRequestScheduleType(scheduleType)
	config := apiclient.MonitorConfigRequest{
		Schedule:     schedule,
		ScheduleType: configType,
	}
	setNullableInt64(m.CheckinMargin, &config.CheckinMargin)
	setNullableInt64(m.MaxRuntime, &config.MaxRuntime)
	setNullableString(m.Timezone, &config.Timezone)
	setNullableInt64(m.FailureIssueThreshold, &config.FailureIssueThreshold)
	setNullableInt64(m.RecoveryThreshold, &config.RecoveryThreshold)

	return config, diags
}

func (m *MonitorResourceModel) buildScheduleInput(ctx context.Context) (apiclient.MonitorScheduleInput, string, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	hasScheduleCrontab := !m.ScheduleCrontab.IsNull() && !m.ScheduleCrontab.IsUnknown()
	hasScheduleInterval := !m.ScheduleInterval.IsNull() && !m.ScheduleInterval.IsUnknown()

	if hasScheduleCrontab && hasScheduleInterval {
		diags.AddError("Invalid configuration", "Only one of `schedule_crontab` or `schedule_interval` can be configured.")
		return apiclient.MonitorScheduleInput{}, "", diags
	}

	if !hasScheduleCrontab && !hasScheduleInterval {
		diags.AddError("Invalid configuration", "One of `schedule_crontab` or `schedule_interval` must be configured.")
		return apiclient.MonitorScheduleInput{}, "", diags
	}

	scheduleType := monitorScheduleTypeCrontab
	if hasScheduleInterval {
		scheduleType = monitorScheduleTypeInterval
	}

	switch scheduleType {
	case monitorScheduleTypeCrontab:
		if !hasScheduleCrontab {
			diags.AddError("Invalid configuration", "`schedule_crontab` must be configured.")
			return apiclient.MonitorScheduleInput{}, "", diags
		}

		var schedule apiclient.MonitorScheduleInput
		if err := schedule.FromMonitorScheduleInput0(m.ScheduleCrontab.ValueString()); err != nil {
			diags.AddError("Invalid configuration", fmt.Sprintf("Failed to serialize `schedule_crontab`: %s", err))
			return apiclient.MonitorScheduleInput{}, "", diags
		}

		return schedule, scheduleType, diags
	case monitorScheduleTypeInterval:
		if !hasScheduleInterval {
			diags.AddError("Invalid configuration", "`schedule_interval` must be configured.")
			return apiclient.MonitorScheduleInput{}, "", diags
		}

		var interval MonitorIntervalModel
		diags.Append(m.ScheduleInterval.As(ctx, &interval, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return apiclient.MonitorScheduleInput{}, "", diags
		}

		valueItem := apiclient.MonitorScheduleInput_1_Item{}
		if err := valueItem.FromMonitorScheduleInput10(interval.Value.ValueInt64()); err != nil {
			diags.AddError("Invalid configuration", fmt.Sprintf("Failed to serialize `interval.value`: %s", err))
			return apiclient.MonitorScheduleInput{}, "", diags
		}

		unitItem := apiclient.MonitorScheduleInput_1_Item{}
		if err := unitItem.FromMonitorScheduleInput11(apiclient.MonitorScheduleInput11(interval.Unit.ValueString())); err != nil {
			diags.AddError("Invalid configuration", fmt.Sprintf("Failed to serialize `interval.unit`: %s", err))
			return apiclient.MonitorScheduleInput{}, "", diags
		}

		var schedule apiclient.MonitorScheduleInput
		if err := schedule.FromMonitorScheduleInput1(apiclient.MonitorScheduleInput1{valueItem, unitItem}); err != nil {
			diags.AddError("Invalid configuration", fmt.Sprintf("Failed to serialize `schedule_interval`: %s", err))
			return apiclient.MonitorScheduleInput{}, "", diags
		}

		return schedule, scheduleType, diags
	default:
		diags.AddError("Invalid configuration", fmt.Sprintf("Unsupported `schedule_type` value %q.", scheduleType))
		return apiclient.MonitorScheduleInput{}, "", diags
	}
}

var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithConfigure = &MonitorResource{}
var _ resource.ResourceWithConfigValidators = &MonitorResource{}
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

func (r *MonitorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("schedule_crontab"),
			path.MatchRoot("schedule_interval"),
		),
	}
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Cron Monitor resource. This resource manages monitors used by the Crons product.",

		Attributes: map[string]schema.Attribute{
			"id": ResourceIdAttribute(),
			"organization": schema.StringAttribute{
				MarkdownDescription: "Organization slug or ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "Project slug or ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the monitor.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "Unique monitor slug. If omitted, Sentry derives one from `name`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "Owner actor in the format `user:<id>` or `team:<id>`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_crontab": schema.StringAttribute{
				MarkdownDescription: "Crontab schedule expression.",
				Optional:            true,
				Computed:            true,
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "IANA timezone for crontab schedules.",
				Optional:            true,
				Computed:            true,
			},
			"schedule_interval": schema.SingleNestedAttribute{
				MarkdownDescription: "Interval schedule.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{
						MarkdownDescription: "Interval value.",
						Required:            true,
					},
					"unit": schema.StringAttribute{
						MarkdownDescription: "Interval unit. One of `year`, `month`, `week`, `day`, `hour`, `minute`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("year", "month", "week", "day", "hour", "minute"),
						},
					},
				},
			},
			"checkin_margin": schema.Int64Attribute{
				MarkdownDescription: "Minutes after the expected check-in before the run is considered missed.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_runtime": schema.Int64Attribute{
				MarkdownDescription: "Maximum runtime in minutes before an in-progress check-in times out.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"failure_issue_threshold": schema.Int64Attribute{
				MarkdownDescription: "Consecutive failed or missed check-ins before creating an issue.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"recovery_threshold": schema.Int64Attribute{
				MarkdownDescription: "Consecutive successful check-ins before resolving an issue.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, diags := data.buildCreateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CreateOrganizationMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		requestBody,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpResp.StatusCode() != http.StatusCreated || httpResp.JSON201 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON201)...)
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

	httpResp, err := r.apiClient.GetProjectMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("monitor"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MonitorResourceModel
	var state MonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, diags := plan.buildUpdateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UpdateProjectMonitorWithResponse(
		ctx,
		plan.Organization.ValueString(),
		plan.Project.ValueString(),
		state.Id.ValueString(),
		requestBody,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("monitor"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(plan.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteProjectMonitorWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() != http.StatusAccepted {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateThreePartId(ctx, "organization", "project", req, resp)
}
