package provider

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryplatforms"
)

var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithConfigure = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	baseResource
}

type ProjectResourceFilterModel struct {
	BlacklistedIps types.Set `tfsdk:"blacklisted_ips"`
	Releases       types.Set `tfsdk:"releases"`
	ErrorMessages  types.Set `tfsdk:"error_messages"`
}

func (data *ProjectResourceFilterModel) Fill(project sentry.Project) error {
	if project.Options == nil {
		data.BlacklistedIps = types.SetNull(types.StringType)
		data.Releases = types.SetNull(types.StringType)
		data.ErrorMessages = types.SetNull(types.StringType)
		return nil
	}

	if values, ok := project.Options["filters:blacklisted_ips"].(string); ok {
		if values == "" {
			data.BlacklistedIps = types.SetNull(types.StringType)
		} else {
			var elements []attr.Value
			for _, value := range strings.Split(values, "\n") {
				elements = append(elements, types.StringValue(value))
			}
			data.BlacklistedIps = types.SetValueMust(types.StringType, elements)
		}
	} else {
		return fmt.Errorf("invalid type for filters:blacklisted_ips: %T", project.Options["filters:blacklisted_ips"])
	}

	if values, ok := project.Options["filters:releases"].(string); ok {
		if values == "" {
			data.Releases = types.SetNull(types.StringType)
		} else {
			var elements []attr.Value
			for _, value := range strings.Split(values, "\n") {
				elements = append(elements, types.StringValue(value))
			}
			data.Releases = types.SetValueMust(types.StringType, elements)
		}
	} else {
		return fmt.Errorf("invalid type for filters:releases: %T", project.Options["filters:releases"])
	}

	if values, ok := project.Options["filters:error_messages"].(string); ok {
		if values == "" {
			data.ErrorMessages = types.SetNull(types.StringType)
		} else {
			var elements []attr.Value
			for _, value := range strings.Split(values, "\n") {
				elements = append(elements, types.StringValue(value))
			}
			data.ErrorMessages = types.SetValueMust(types.StringType, elements)
		}
	} else {
		return fmt.Errorf("invalid type for filters:error_messages: %T", project.Options["filters:error_messages"])
	}

	return nil
}

type ProjectResourceModel struct {
	Id              types.String                `tfsdk:"id"`
	Organization    types.String                `tfsdk:"organization"`
	Teams           types.Set                   `tfsdk:"teams"`
	Name            types.String                `tfsdk:"name"`
	Slug            types.String                `tfsdk:"slug"`
	Platform        types.String                `tfsdk:"platform"`
	DefaultRules    types.Bool                  `tfsdk:"default_rules"`
	DefaultKey      types.Bool                  `tfsdk:"default_key"`
	InternalId      types.String                `tfsdk:"internal_id"`
	Features        types.Set                   `tfsdk:"features"`
	DigestsMinDelay types.Int64                 `tfsdk:"digests_min_delay"`
	DigestsMaxDelay types.Int64                 `tfsdk:"digests_max_delay"`
	ResolveAge      types.Int64                 `tfsdk:"resolve_age"`
	Filters         *ProjectResourceFilterModel `tfsdk:"filters"`
}

func (data *ProjectResourceModel) Fill(organization string, project sentry.Project) error {
	data.Id = types.StringValue(project.Slug)
	data.Organization = types.StringValue(organization)
	data.Name = types.StringValue(project.Name)
	data.Slug = types.StringValue(project.Slug)

	if project.Platform == "" {
		data.Platform = types.StringNull()
	} else {
		data.Platform = types.StringValue(project.Platform)
	}

	data.InternalId = types.StringValue(project.ID)

	if data.DigestsMinDelay.IsNull() {
		data.DigestsMinDelay = types.Int64Null()
	} else {
		data.DigestsMinDelay = types.Int64Value(int64(project.DigestsMinDelay))
	}
	if data.DigestsMaxDelay.IsNull() {
		data.DigestsMaxDelay = types.Int64Null()
	} else {
		data.DigestsMaxDelay = types.Int64Value(int64(project.DigestsMaxDelay))
	}
	if data.ResolveAge.IsNull() {
		data.ResolveAge = types.Int64Null()
	} else {
		data.ResolveAge = types.Int64Value(int64(project.ResolveAge))
	}

	var teamElements []attr.Value
	for _, team := range project.Teams {
		teamElements = append(teamElements, types.StringPointerValue(team.Slug))
	}
	data.Teams = types.SetValueMust(types.StringType, teamElements)

	var featureElements []attr.Value
	for _, feature := range project.Features {
		featureElements = append(featureElements, types.StringValue(feature))
	}
	data.Features = types.SetValueMust(types.StringType, featureElements)

	if data.Filters != nil {
		data.Filters.Fill(project)
	}

	return nil
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Description: "The slug of the organization the project belongs to.",
				Required:    true,
			},
			"teams": schema.SetAttribute{
				Description: "The slugs of the teams to create the project for.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name for the project.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The optional slug for this project.",
				Optional:    true,
				Computed:    true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "The platform for this project. For a list of valid values, [see this page](https://github.com/jianyuan/terraform-provider-sentry/blob/main/internal/sentryplatforms/platforms.txt). Use `other` for platforms not listed.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(sentryplatforms.Platforms...),
				},
			},
			"default_rules": schema.BoolAttribute{
				Description: "Whether to create a default issue alert. Defaults to true where the behavior is to alert the user on every new issue.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_key": schema.BoolAttribute{
				Description: "Whether to create a default key. By default, Sentry will create a key for you. If you wish to manage keys manually, set this to false and create keys using the `sentry_key` resource.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"internal_id": schema.StringAttribute{
				Description: "The internal ID for this project.",
				Computed:    true,
			},
			"features": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"digests_min_delay": schema.Int64Attribute{
				Description: "The minimum amount of time (in seconds) to wait between scheduling digests for delivery after the initial scheduling.",
				Optional:    true,
			},
			"digests_max_delay": schema.Int64Attribute{
				Description: "The maximum amount of time (in seconds) to wait between scheduling digests for delivery.",
				Optional:    true,
			},
			"resolve_age": schema.Int64Attribute{
				Description: "Hours in which an issue is automatically resolve if not seen after this amount of time.",
				Optional:    true,
			},
			"filters": schema.SingleNestedAttribute{
				Description: "Custom filters for this project.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"blacklisted_ips": schema.SetAttribute{
						Description: "Filter events from these IP addresses. (e.g. 127.0.0.1 or 10.0.0.0/8)",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"releases": schema.SetAttribute{
						MarkdownDescription: "Filter events from these releases. Allows [glob pattern matching](https://en.wikipedia.org/wiki/Glob_(programming)). (e.g. 1.* or [!3].[0-9].*)",
						ElementType:         types.StringType,
						Optional:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"error_messages": schema.SetAttribute{
						MarkdownDescription: "Filter events by error messages. Allows [glob pattern matching](https://en.wikipedia.org/wiki/Glob_(programming)). (e.g. TypeError* or *: integer division or modulo by zero)",
						ElementType:         types.StringType,
						Optional:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
		},
	}
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var teams []string
	resp.Diagnostics.Append(data.Teams.ElementsAs(ctx, &teams, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(teams) == 0 {
		resp.Diagnostics.AddError("Client Error", "At least one team is required")
		return
	}

	// Create the project
	project, _, err := r.client.Projects.Create(
		ctx,
		data.Organization.ValueString(),
		teams[0],
		&sentry.CreateProjectParams{
			Name:         data.Name.ValueString(),
			Slug:         data.Slug.ValueString(),
			Platform:     data.Platform.ValueString(),
			DefaultRules: data.DefaultRules.ValueBoolPointer(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating project: %s", err.Error()))
		return
	}

	// Update the project
	updateParams := &sentry.UpdateProjectParams{
		Name:     data.Name.ValueString(),
		Slug:     data.Slug.ValueString(),
		Platform: data.Platform.ValueString(),
	}
	if !data.DigestsMinDelay.IsNull() {
		updateParams.DigestsMinDelay = sentry.Int(int(data.DigestsMinDelay.ValueInt64()))
	}
	if !data.DigestsMaxDelay.IsNull() {
		updateParams.DigestsMaxDelay = sentry.Int(int(data.DigestsMaxDelay.ValueInt64()))
	}
	if !data.ResolveAge.IsNull() {
		updateParams.ResolveAge = sentry.Int(int(data.ResolveAge.ValueInt64()))
	}
	if data.Filters != nil {
		updateParams.Options = make(map[string]interface{})
		if data.Filters.BlacklistedIps.IsNull() {
			updateParams.Options["filters:blacklisted_ips"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(data.Filters.BlacklistedIps.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			updateParams.Options["filters:blacklisted_ips"] = strings.Join(values, "\n")
		}

		if data.Filters.Releases.IsNull() {
			updateParams.Options["filters:releases"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(data.Filters.Releases.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			updateParams.Options["filters:releases"] = strings.Join(values, "\n")
		}

		if data.Filters.ErrorMessages.IsNull() {
			updateParams.Options["filters:error_messages"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(data.Filters.ErrorMessages.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			updateParams.Options["filters:error_messages"] = strings.Join(values, "\n")
		}
	}

	project, apiResp, err := r.client.Projects.Update(
		ctx,
		data.Organization.ValueString(),
		project.ID,
		updateParams,
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Project not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating project: %s", err.Error()))
		return
	}

	// If the default key is set to false, remove the default key
	if !data.DefaultKey.IsNull() && !data.DefaultKey.ValueBool() {
		if err := r.removeDefaultKey(ctx, data.Organization.ValueString(), project); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing default key: %s", err.Error()))
			return
		}
	}

	// Add additional teams
	if len(teams) > 1 {
		for _, team := range teams[1:] {
			_, _, err := r.client.Projects.AddTeam(ctx, data.Organization.ValueString(), project.Slug, team)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error adding team to project: %s", err.Error()))
				return
			}
		}
	}

	project, apiResp, err = r.client.Projects.Get(ctx, data.Organization.ValueString(), project.Slug)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Project not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading project: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *project); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling project: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, apiResp, err := r.client.Projects.Get(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Project not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading project: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *project); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling project: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.UpdateProjectParams{
		Name:     plan.Name.ValueString(),
		Slug:     plan.Slug.ValueString(),
		Platform: plan.Platform.ValueString(),
	}
	if !plan.DigestsMinDelay.IsNull() {
		params.DigestsMinDelay = sentry.Int(int(plan.DigestsMinDelay.ValueInt64()))
	}
	if !plan.DigestsMaxDelay.IsNull() {
		params.DigestsMaxDelay = sentry.Int(int(plan.DigestsMaxDelay.ValueInt64()))
	}
	if !plan.ResolveAge.IsNull() {
		params.ResolveAge = sentry.Int(int(plan.ResolveAge.ValueInt64()))
	}
	if plan.Filters != nil {
		params.Options = make(map[string]interface{})
		if plan.Filters.BlacklistedIps.IsNull() {
			params.Options["filters:blacklisted_ips"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(plan.Filters.BlacklistedIps.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			params.Options["filters:blacklisted_ips"] = strings.Join(values, "\n")
		}

		if plan.Filters.Releases.IsNull() {
			params.Options["filters:releases"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(plan.Filters.Releases.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			params.Options["filters:releases"] = strings.Join(values, "\n")
		}

		if plan.Filters.ErrorMessages.IsNull() {
			params.Options["filters:error_messages"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(plan.Filters.ErrorMessages.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			params.Options["filters:error_messages"] = strings.Join(values, "\n")
		}
	}

	project, apiResp, err := r.client.Projects.Update(
		ctx,
		plan.Organization.ValueString(),
		plan.Id.ValueString(),
		params,
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Project not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating project: %s", err.Error()))
		return
	}

	// If the default key is set to false, remove the default key
	if !plan.DefaultKey.IsNull() && !plan.DefaultKey.ValueBool() {
		if err := r.removeDefaultKey(ctx, plan.Organization.ValueString(), project); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing default key: %s", err.Error()))
			return
		}
	}

	// Update teams
	if !plan.Teams.Equal(state.Teams) {
		var planTeams, stateTeams []string

		resp.Diagnostics.Append(plan.Teams.ElementsAs(ctx, &planTeams, false)...)
		resp.Diagnostics.Append(state.Teams.ElementsAs(ctx, &stateTeams, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Add teams
		for _, team := range planTeams {
			if !slices.Contains(stateTeams, team) {
				_, _, err := r.client.Projects.AddTeam(
					ctx,
					plan.Organization.ValueString(),
					project.Slug,
					team,
				)
				if err != nil {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error adding team to project: %s", err.Error()))
					return
				}
			}
		}

		// Remove teams
		for _, team := range stateTeams {
			if !slices.Contains(planTeams, team) {
				_, err := r.client.Projects.RemoveTeam(
					ctx,
					plan.Organization.ValueString(),
					project.Slug,
					team,
				)
				if err != nil {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing team from project: %s", err.Error()))
					return
				}
			}
		}
	}

	project, apiResp, err = r.client.Projects.Get(
		ctx,
		plan.Organization.ValueString(),
		plan.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Project not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading project: %s", err.Error()))
		return
	}

	if err := plan.Fill(plan.Organization.ValueString(), *project); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling project: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectResource) removeDefaultKey(ctx context.Context, organization string, project *sentry.Project) error {
	params := &sentry.ListProjectKeysParams{}

	for {
		keys, apiResp, err := r.client.ProjectKeys.List(ctx, organization, project.Slug, params)
		if err != nil {
			return err
		}

		for _, key := range keys {
			if key.Name == "Default" {
				_, err := r.client.ProjectKeys.Delete(ctx, organization, project.ID, key.ID)
				if err != nil {
					return err
				}

				return nil
			}
		}

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	return nil
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.Projects.Delete(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting project: %s", err.Error()))
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, err := splitTwoPartID(req.ID, "organization", "project-slug")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), project,
	)...)
}
