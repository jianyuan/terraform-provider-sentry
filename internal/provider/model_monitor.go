package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
	"github.com/oapi-codegen/nullable"
)

type MonitorResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Owner        types.String `tfsdk:"owner"`
	Config       types.Object `tfsdk:"config"`
	Status       types.String `tfsdk:"status"`
}

func (m MonitorResourceModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":           ResourceIdAttribute(),
		"organization": ResourceOrganizationAttribute(),
		"project": schema.StringAttribute{
			MarkdownDescription: "The project of this resource.",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"slug": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"owner": schema.StringAttribute{
			Optional: true,
		},
		"config": MonitorConfigResourceModel{}.SchemaAttribute(true),
		"status": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	}
}

func (m MonitorResourceModel) ToMonitorRequest(ctx context.Context) (apiclient.MonitorRequest, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	var monitorConfigResourceModel MonitorConfigResourceModel
	diags.Append(m.Config.As(ctx, &monitorConfigResourceModel, basetypes.ObjectAsOptions{})...)

	monitorConfig, monitorConfigDiags := monitorConfigResourceModel.ToMonitorRequest(ctx, path.Root("config"))
	diags.Append(monitorConfigDiags...)

	owner := nullable.Nullable[string]{}
	if !m.Owner.IsNull() && !m.Owner.IsUnknown() {
		owner.Set(m.Owner.ValueString())
	}

	var slug *string
	if !m.Slug.IsNull() && !m.Slug.IsUnknown() {
		value := m.Slug.ValueString()
		if value != "" {
			slug = &value
		}
	}

	request := apiclient.MonitorRequest{
		Name:    m.Name.ValueString(),
		Slug:    slug,
		Project: m.Project.ValueString(),
		Owner:   owner,
		Config:  monitorConfig,
		Status:  (*apiclient.MonitorRequestStatus)(m.Status.ValueStringPointer()),
	}

	return request, diags
}

func (m *MonitorResourceModel) Fill(ctx context.Context, organization string, monitor apiclient.Monitor) (diags diag.Diagnostics) {
	path := path.Empty()

	var config MonitorConfigResourceModel
	diags.Append(config.Fill(ctx, path.AtName("config"), monitor.Config)...)

	m.Organization = types.StringValue(organization)
	m.Id = types.StringValue(monitor.Id)
	m.Name = types.StringValue(monitor.Name)
	m.Slug = types.StringValue(monitor.Slug)
	m.Project = types.StringValue(monitor.Project.Slug)
	m.Owner = types.StringPointerValue(formatMonitorOwner(monitor.Owner))
	m.Config = tfutils.MergeDiagnostics(types.ObjectValueFrom(ctx, config.AttributeTypes(), config))(&diags)

	m.Status = types.StringValue(monitor.Status)

	return
}

func formatMonitorOwner(owner nullable.Nullable[struct {
	Id   string                     `json:"id"`
	Type apiclient.MonitorOwnerType `json:"type"`
}]) *string {
	if owner.IsNull() || !owner.IsSpecified() {
		return nil
	}

	parsed, err := owner.Get()
	if err != nil || parsed.Id == "" || parsed.Type == "" {
		return nil
	}

	formatted := string(parsed.Type) + ":" + parsed.Id

	return &formatted
}
