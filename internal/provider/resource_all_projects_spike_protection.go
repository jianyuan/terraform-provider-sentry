package provider

import (
	"context"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/go-utils/sliceutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type AllProjectsSpikeProtectionResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Projects     types.Set    `tfsdk:"projects"`
}

func (m *AllProjectsSpikeProtectionResourceModel) Fill(ctx context.Context, projects []apiclient.Project) (diags diag.Diagnostics) {
	m.Projects = types.SetValueMust(types.StringType, sliceutils.Map(func(project apiclient.Project) attr.Value {
		return types.StringValue(project.Slug)
	}, projects))
	return
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Toggle the browser-extensions, localhost, filtered-transaction, or web-crawlers filter on or off for all projects.",
				Required:            true,
			},
		},
	}
}

func (r *AllProjectsSpikeProtectionResource) readProjects(ctx context.Context, organization string, enabled bool, projectSlugs []string) ([]apiclient.Project, diag.Diagnostics) {
	var diags diag.Diagnostics

	var allProjects []apiclient.Project
	params := &apiclient.ListOrganizationProjectsParams{
		Options: ptr.Ptr([]string{"quotas:spike-protection-disabled"}),
	}
	for {
		httpResp, err := r.apiClient.ListOrganizationProjectsWithResponse(
			ctx,
			organization,
			params,
		)
		if err != nil {
			diags.Append(diagutils.NewClientError("read", err))
			return nil, diags
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			diags.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return nil, diags
		}

		for _, project := range *httpResp.JSON200 {
			if slices.Contains(projectSlugs, project.Slug) {
				if projectDisabled, ok := project.Options["quotas:spike-protection-disabled"].(bool); ok && projectDisabled != enabled {
					allProjects = append(allProjects, project)
				}
			}
		}

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	return allProjects, nil
}

func (r *AllProjectsSpikeProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AllProjectsSpikeProtectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var projects []string
	if !data.Projects.IsNull() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	}

	if data.Enabled.ValueBool() {
		httpResp, err := r.apiClient.EnableSpikeProtectionWithResponse(
			ctx,
			data.Organization.ValueString(),
			apiclient.EnableSpikeProtectionJSONRequestBody{
				Projects: projects,
			},
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("enable", err))
			return
		} else if httpResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("enable", httpResp.StatusCode(), httpResp.Body))
			return
		}
	} else {
		httpResp, err := r.apiClient.DisableSpikeProtectionWithResponse(
			ctx,
			data.Organization.ValueString(),
			apiclient.DisableSpikeProtectionJSONRequestBody{
				Projects: projects,
			},
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("disable", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("disable", httpResp.StatusCode(), httpResp.Body))
			return
		}
	}

	allProjects := tfutils.MergeDiagnostics(r.readProjects(ctx, data.Organization.ValueString(), data.Enabled.ValueBool(), projects))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, allProjects)...)
	if resp.Diagnostics.HasError() {
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

	allProjects := tfutils.MergeDiagnostics(r.readProjects(ctx, data.Organization.ValueString(), data.Enabled.ValueBool(), projects))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, allProjects)...)
	if resp.Diagnostics.HasError() {
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
		httpResp, err := r.apiClient.EnableSpikeProtectionWithResponse(
			ctx,
			data.Organization.ValueString(),
			apiclient.EnableSpikeProtectionJSONRequestBody{
				Projects: projects,
			},
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("enable", err))
			return
		} else if httpResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("enable", httpResp.StatusCode(), httpResp.Body))
			return
		}
	} else {
		httpResp, err := r.apiClient.DisableSpikeProtectionWithResponse(
			ctx,
			data.Organization.ValueString(),
			apiclient.DisableSpikeProtectionJSONRequestBody{
				Projects: projects,
			},
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("disable", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("disable", httpResp.StatusCode(), httpResp.Body))
			return
		}
	}

	allProjects := tfutils.MergeDiagnostics(r.readProjects(ctx, data.Organization.ValueString(), data.Enabled.ValueBool(), projects))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, allProjects)...)
	if resp.Diagnostics.HasError() {
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
		httpResp, err := r.apiClient.DisableSpikeProtectionWithResponse(
			ctx,
			data.Organization.ValueString(),
			apiclient.DisableSpikeProtectionJSONRequestBody{
				Projects: projects,
			},
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("disable", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("disable", httpResp.StatusCode(), httpResp.Body))
			return
		}
	} else {
		// We need to enable the spike protection if it was disabled.
		httpResp, err := r.apiClient.EnableSpikeProtectionWithResponse(
			ctx,
			data.Organization.ValueString(),
			apiclient.EnableSpikeProtectionJSONRequestBody{
				Projects: projects,
			},
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("enable", err))
			return
		} else if httpResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("enable", httpResp.StatusCode(), httpResp.Body))
			return
		}
	}
}
