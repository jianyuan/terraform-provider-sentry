package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type MonitorResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Type         types.String `tfsdk:"type"`
	Name         types.String `tfsdk:"name"`
	Owner        types.String `tfsdk:"owner"`
	Project      types.String `tfsdk:"project"`
	Slug         types.String `tfsdk:"slug"`
	Config       types.Object `tfsdk:"config"`
	IsMuted      types.Bool   `tfsdk:"is_muted"`
	Status       types.String `tfsdk:"status"`
}

func (m MonitorResourceModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":           ResourceIdAttribute(),
		"organization": ResourceOrganizationAttribute(),
		"project":      ResourceProjectAttribute(),
		"type": schema.StringAttribute{
			Required: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"owner": schema.StringAttribute{
			Optional: true,
		},
		"slug": schema.StringAttribute{
			Required: true,
		},
		"config": MonitorConfigResourceModel{}.SchemaAttribute(true),
		"is_muted": schema.BoolAttribute{
			Optional: true,
		},
		"status": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	}
}

func (m *MonitorResourceModel) Fill(ctx context.Context, organization string, monitor apiclient.Monitor) (diags diag.Diagnostics) {
	path := path.Empty()

	m.Id = types.StringValue(monitor.Id)
	m.Type = types.StringValue(string(monitor.Type))
	m.Name = types.StringValue(monitor.Name)
	m.Slug = types.StringValue(monitor.Slug)
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(monitor.Project.Slug)
	m.Status = types.StringValue(monitor.Status)
	m.IsMuted = types.BoolValue(monitor.IsMuted)

	if monitor.Owner == nil {
		m.Owner = types.StringNull()
	} else {
		m.Owner = types.StringValue(string(monitor.Owner.Type) + ":" + monitor.Owner.Id)
	}

	var config MonitorConfigResourceModel
	m.Config.As(ctx, &config, basetypes.ObjectAsOptions{})
	diags.Append(config.Fill(ctx, path.AtName("config"), monitor.Config)...)
	m.Config = tfutils.MergeDiagnostics(types.ObjectValueFrom(ctx, config.AttributeTypes(), config))(&diags)

	return
}
