package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

type AllProjectsSpikeProtectionResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Projects     types.Set    `tfsdk:"projects"`
}

func (m *AllProjectsSpikeProtectionResourceModel) Fill(organization string, enabled bool, projects []sentry.Project) error {
	m.Organization = types.StringValue(organization)
	m.Enabled = types.BoolValue(enabled)

	projectElements := []attr.Value{}
	for _, project := range projects {
		projectElements = append(projectElements, types.StringValue(project.Slug))
	}
	m.Projects = types.SetValueMust(types.StringType, projectElements)

	return nil
}

var _ resource.Resource = &AllProjectsSpikeProtectionResource{}
var _ resource.ResourceWithConfigure = &AllProjectsSpikeProtectionResource{}

func NewAllProjectsSpikeProtectionResource() resource.Resource {
	return &AllProjectsSpikeProtectionResource{}
}

type AllProjectsSpikeProtectionResource struct {
	baseResource
}

func (r *AllProjectsSpikeProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_projects_spike_protection"
}

func (r *AllProjectsSpikeProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Enable spike protection for all projects in an organization.",

		Attributes: map[string]schema.Attribute{
			"organization": ResourceOrganizationAttribute(),
			"projects": schema.SetAttribute{
				MarkdownDescription: "The slugs of the projects to enable or disable spike protection for.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Toggle the browser-extensions, localhost, filtered-transaction, or web-crawlers filter on or off for all projects.",
				Required:            true,
			},
		},
	}
}

func (r *AllProjectsSpikeProtectionResource) readProjects(ctx context.Context, organization string, enabled bool, projectSlugs []string) ([]sentry.Project, error) {
	var allProjects []sentry.Project
	params := &sentry.ListOrganizationProjectsParams{
		Options: "quotas:spike-protection-disabled",
	}

	for {
		projects, apiResp, err := r.client.OrganizationProjects.List(ctx, organization, params)
		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			for _, projectSlug := range projectSlugs {
				if projectSlug == project.Slug {
					if projectDisabled, ok := project.Options["quotas:spike-protection-disabled"].(bool); ok && projectDisabled != enabled {
						allProjects = append(allProjects, *project)
					}

					break
				}
			}
		}

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	return allProjects, nil
}

func (r *AllProjectsSpikeProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AllProjectsSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	if !data.Projects.IsNull() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	}

	if data.Enabled.ValueBool() {
		_, err := r.client.SpikeProtections.Enable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: projects,
			},
		)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "enable", err)
			return
		}
	} else {
		_, err := r.client.SpikeProtections.Disable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: projects,
			},
		)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "disable", err)
			return
		}
	}

	allProjects, err := r.readProjects(ctx, data.Organization.ValueString(), data.Enabled.ValueBool(), projects)
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "create", err)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Enabled.ValueBool(), allProjects); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AllProjectsSpikeProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AllProjectsSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	if !data.Projects.IsNull() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	}

	allProjects, err := r.readProjects(ctx, data.Organization.ValueString(), data.Enabled.ValueBool(), projects)
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "read", err)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Enabled.ValueBool(), allProjects); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AllProjectsSpikeProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AllProjectsSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	if !data.Projects.IsNull() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	}

	if data.Enabled.ValueBool() {
		_, err := r.client.SpikeProtections.Enable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: projects,
			},
		)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "enable", err)
			return
		}
	} else {
		_, err := r.client.SpikeProtections.Disable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: projects,
			},
		)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "disable", err)
			return
		}
	}

	allProjects, err := r.readProjects(ctx, data.Organization.ValueString(), data.Enabled.ValueBool(), projects)
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "update", err)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Enabled.ValueBool(), allProjects); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AllProjectsSpikeProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AllProjectsSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	if !data.Projects.IsNull() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	}

	if data.Enabled.ValueBool() {
		// We need to disable the spike protection if it was enabled.
		_, err := r.client.SpikeProtections.Disable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: projects,
			},
		)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "disable", err)
			return
		}
	} else {
		// We need to enable the spike protection if it was disabled.
		_, err := r.client.SpikeProtections.Enable(
			ctx,
			data.Organization.ValueString(),
			&sentry.SpikeProtectionParams{
				Projects: projects,
			},
		)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "enable", err)
			return
		}
	}
}
