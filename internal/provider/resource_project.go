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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/go-utils/sliceutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryplatforms"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type ProjectFilterResourceModel struct {
	BlacklistedIps types.Set `tfsdk:"blacklisted_ips"`
	Releases       types.Set `tfsdk:"releases"`
	ErrorMessages  types.Set `tfsdk:"error_messages"`
}

func (m *ProjectFilterResourceModel) Fill(project apiclient.Project) error {
	if project.Options == nil {
		m.BlacklistedIps = types.SetNull(types.StringType)
		m.Releases = types.SetNull(types.StringType)
		m.ErrorMessages = types.SetNull(types.StringType)
		return nil
	}

	if values, ok := project.Options["filters:blacklisted_ips"].(string); ok {
		if values == "" {
			m.BlacklistedIps = types.SetNull(types.StringType)
		} else {
			m.BlacklistedIps = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		return fmt.Errorf("invalid type for filters:blacklisted_ips: %T", project.Options["filters:blacklisted_ips"])
	}

	if values, ok := project.Options["filters:releases"].(string); ok {
		if values == "" {
			m.Releases = types.SetNull(types.StringType)
		} else {
			m.Releases = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		return fmt.Errorf("invalid type for filters:releases: %T", project.Options["filters:releases"])
	}

	if values, ok := project.Options["filters:error_messages"].(string); ok {
		if values == "" {
			m.ErrorMessages = types.SetNull(types.StringType)
		} else {
			m.ErrorMessages = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		return fmt.Errorf("invalid type for filters:error_messages: %T", project.Options["filters:error_messages"])
	}

	return nil
}

func (m *ProjectFilterResourceModel) FillOld(project sentry.Project) error {
	if project.Options == nil {
		m.BlacklistedIps = types.SetNull(types.StringType)
		m.Releases = types.SetNull(types.StringType)
		m.ErrorMessages = types.SetNull(types.StringType)
		return nil
	}

	if values, ok := project.Options["filters:blacklisted_ips"].(string); ok {
		if values == "" {
			m.BlacklistedIps = types.SetNull(types.StringType)
		} else {
			m.BlacklistedIps = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		return fmt.Errorf("invalid type for filters:blacklisted_ips: %T", project.Options["filters:blacklisted_ips"])
	}

	if values, ok := project.Options["filters:releases"].(string); ok {
		if values == "" {
			m.Releases = types.SetNull(types.StringType)
		} else {
			m.Releases = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		return fmt.Errorf("invalid type for filters:releases: %T", project.Options["filters:releases"])
	}

	if values, ok := project.Options["filters:error_messages"].(string); ok {
		if values == "" {
			m.ErrorMessages = types.SetNull(types.StringType)
		} else {
			m.ErrorMessages = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		return fmt.Errorf("invalid type for filters:error_messages: %T", project.Options["filters:error_messages"])
	}

	return nil
}

type ProjectResourceModel struct {
	Id                   types.String                `tfsdk:"id"`
	Organization         types.String                `tfsdk:"organization"`
	Teams                types.Set                   `tfsdk:"teams"`
	Name                 types.String                `tfsdk:"name"`
	Slug                 types.String                `tfsdk:"slug"`
	Platform             types.String                `tfsdk:"platform"`
	DefaultRules         types.Bool                  `tfsdk:"default_rules"`
	DefaultKey           types.Bool                  `tfsdk:"default_key"`
	InternalId           types.String                `tfsdk:"internal_id"`
	Features             types.Set                   `tfsdk:"features"`
	DigestsMinDelay      types.Int64                 `tfsdk:"digests_min_delay"`
	DigestsMaxDelay      types.Int64                 `tfsdk:"digests_max_delay"`
	ResolveAge           types.Int64                 `tfsdk:"resolve_age"`
	Filters              *ProjectFilterResourceModel `tfsdk:"filters"`
	FingerprintingRules  types.String                `tfsdk:"fingerprinting_rules"`
	GroupingEnhancements types.String                `tfsdk:"grouping_enhancements"`
}

func (m *ProjectResourceModel) Fill(project apiclient.Project) error {
	m.Id = types.StringValue(project.Slug)
	m.Organization = types.StringValue(project.Organization.Slug)
	m.Teams = types.SetValueMust(types.StringType, sliceutils.Map(func(v apiclient.Team) attr.Value {
		return types.StringValue(v.Slug)
	}, project.Teams))
	m.Name = types.StringValue(project.Name)
	m.Slug = types.StringValue(project.Slug)

	if project.Platform == nil || *project.Platform == "" {
		m.Platform = types.StringNull()
	} else {
		m.Platform = types.StringPointerValue(project.Platform)
	}

	m.InternalId = types.StringValue(project.Id)
	m.Features = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
		return types.StringValue(v)
	}, project.Features))

	if !m.DigestsMinDelay.IsNull() {
		m.DigestsMinDelay = types.Int64Value(project.DigestsMinDelay)
	}
	if !m.DigestsMaxDelay.IsNull() {
		m.DigestsMaxDelay = types.Int64Value(project.DigestsMaxDelay)
	}
	if !m.ResolveAge.IsNull() {
		m.ResolveAge = types.Int64Value(project.ResolveAge)
	}

	if m.Filters != nil {
		m.Filters.Fill(project)
	}

	if !m.FingerprintingRules.IsNull() {
		if strings.TrimRight(m.FingerprintingRules.ValueString(), "\n") != project.FingerprintingRules {
			m.FingerprintingRules = types.StringValue(project.FingerprintingRules)
		}
	}

	if !m.GroupingEnhancements.IsNull() {
		if strings.TrimRight(m.GroupingEnhancements.ValueString(), "\n") != project.GroupingEnhancements {
			m.GroupingEnhancements = types.StringValue(project.GroupingEnhancements)
		}
	}

	return nil
}

func (m *ProjectResourceModel) FillOld(organization string, project sentry.Project) error {
	m.Id = types.StringValue(project.Slug)
	m.Organization = types.StringValue(organization)
	m.Name = types.StringValue(project.Name)
	m.Slug = types.StringValue(project.Slug)

	if project.Platform == "" {
		m.Platform = types.StringNull()
	} else {
		m.Platform = types.StringValue(project.Platform)
	}

	m.InternalId = types.StringValue(project.ID)

	if !m.DigestsMinDelay.IsNull() {
		m.DigestsMinDelay = types.Int64Value(int64(project.DigestsMinDelay))
	}
	if !m.DigestsMaxDelay.IsNull() {
		m.DigestsMaxDelay = types.Int64Value(int64(project.DigestsMaxDelay))
	}
	if !m.ResolveAge.IsNull() {
		m.ResolveAge = types.Int64Value(int64(project.ResolveAge))
	}

	m.Teams = types.SetValueMust(types.StringType, sliceutils.Map(func(v sentry.Team) attr.Value {
		return types.StringPointerValue(v.Slug)
	}, project.Teams))

	m.Features = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
		return types.StringValue(v)
	}, project.Features))

	if m.Filters != nil {
		m.Filters.FillOld(project)
	}

	if project.FingerprintingRules == "" {
		if !m.FingerprintingRules.IsNull() {
			m.FingerprintingRules = types.StringNull()
		}
	} else {
		if m.FingerprintingRules.IsNull() || strings.TrimRight(m.FingerprintingRules.ValueString(), "\n") != project.FingerprintingRules {
			m.FingerprintingRules = types.StringValue(project.FingerprintingRules)
		}
	}

	if project.GroupingEnhancements == "" {
		if !m.GroupingEnhancements.IsNull() {
			m.GroupingEnhancements = types.StringNull()
		}
	} else {
		if m.GroupingEnhancements.IsNull() || strings.TrimRight(m.GroupingEnhancements.ValueString(), "\n") != project.GroupingEnhancements {
			m.GroupingEnhancements = types.StringValue(project.GroupingEnhancements)
		}
	}

	return nil
}

var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithConfigure = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	baseResource
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project resource.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
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
					stringvalidator.OneOf(append(sentryplatforms.Platforms, "other")...),
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
			"fingerprinting_rules": schema.StringAttribute{
				MarkdownDescription: "This can be used to modify the fingerprint rules on the server with custom rules. Rules follow the pattern `matcher:glob -> fingerprint, values`. To learn more about fingerprint rules, [read the docs](https://docs.sentry.io/concepts/data-management/event-grouping/fingerprint-rules/).",
				Optional:            true,
			},
			"grouping_enhancements": schema.StringAttribute{
				MarkdownDescription: "This can be used to enhance the grouping algorithm with custom rules. Rules follow the pattern `matcher:glob [v^]?[+-]flag`. To learn more about stack trace rules, [read the docs](https://docs.sentry.io/concepts/data-management/event-grouping/stack-trace-rules/).",
				Optional:            true,
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
	createBody := apiclient.CreateOrganizationTeamProjectJSONRequestBody{
		Name:         data.Name.ValueString(),
		Platform:     data.Platform.ValueStringPointer(),
		DefaultRules: data.DefaultRules.ValueBoolPointer(),
	}
	if data.Slug.ValueString() != "" {
		createBody.Slug = data.Slug.ValueStringPointer()
	}

	httpRespCreate, err := r.apiClient.CreateOrganizationTeamProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		teams[0],
		createBody,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpRespCreate.StatusCode() != http.StatusCreated || httpRespCreate.JSON201 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpRespCreate.StatusCode(), httpRespCreate.Body))
		return
	}

	// Update the project
	updateBody := apiclient.UpdateOrganizationProjectJSONRequestBody{
		Name:            data.Name.ValueStringPointer(),
		Platform:        data.Platform.ValueStringPointer(),
		DigestsMinDelay: data.DigestsMinDelay.ValueInt64Pointer(),
		DigestsMaxDelay: data.DigestsMaxDelay.ValueInt64Pointer(),
		ResolveAge:      data.ResolveAge.ValueInt64Pointer(),
	}

	if data.Slug.ValueString() != "" {
		updateBody.Slug = data.Slug.ValueStringPointer()
	}

	if data.Filters != nil {
		options := make(map[string]interface{})
		if data.Filters.BlacklistedIps.IsNull() {
			options["filters:blacklisted_ips"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(data.Filters.BlacklistedIps.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:blacklisted_ips"] = strings.Join(values, "\n")
		}

		if data.Filters.Releases.IsNull() {
			options["filters:releases"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(data.Filters.Releases.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:releases"] = strings.Join(values, "\n")
		}

		if data.Filters.ErrorMessages.IsNull() {
			options["filters:error_messages"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(data.Filters.ErrorMessages.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:error_messages"] = strings.Join(values, "\n")
		}

		if len(options) > 0 {
			updateBody.Options = &options
		}
	}

	if data.FingerprintingRules.IsNull() {
		updateBody.FingerprintingRules = sentry.String("")
	} else {
		updateBody.FingerprintingRules = sentry.String(data.FingerprintingRules.ValueString())
	}

	if data.GroupingEnhancements.IsNull() {
		updateBody.GroupingEnhancements = sentry.String("")
	} else {
		updateBody.GroupingEnhancements = sentry.String(data.GroupingEnhancements.ValueString())
	}

	httpRespUpdate, err := r.apiClient.UpdateOrganizationProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		httpRespCreate.JSON201.Slug,
		updateBody,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpRespUpdate.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpRespUpdate.StatusCode() != http.StatusOK || httpRespUpdate.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpRespUpdate.StatusCode(), httpRespUpdate.Body))
		return
	}

	// If the default key is set to false, remove the default key
	if !data.DefaultKey.IsNull() && !data.DefaultKey.ValueBool() {
		if err := r.removeDefaultKey(ctx, data.Organization.ValueString(), *httpRespUpdate.JSON200); err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("remove default key", err))
			return
		}
	}

	// Add additional teams
	if len(teams) > 1 {
		for _, team := range teams[1:] {
			httpResp, err := r.apiClient.AddTeamToProjectWithResponse(
				ctx,
				data.Organization.ValueString(),
				httpRespUpdate.JSON200.Slug,
				team,
			)
			if err != nil {
				resp.Diagnostics.Append(diagutils.NewClientError("add team to project", err))
				return
			} else if httpResp.StatusCode() != http.StatusCreated {
				resp.Diagnostics.Append(diagutils.NewClientStatusError("add team to project", httpResp.StatusCode(), httpResp.Body))
				return
			}
		}
	}

	httpRespRead, err := r.apiClient.GetOrganizationProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		httpRespUpdate.JSON200.Slug,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpRespRead.StatusCode() != http.StatusOK || httpRespRead.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project"))
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(*httpRespRead.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
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

	httpResp, err := r.apiClient.GetOrganizationProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
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

	params := apiclient.UpdateOrganizationProjectJSONRequestBody{
		Name:            plan.Name.ValueStringPointer(),
		Platform:        plan.Platform.ValueStringPointer(),
		DigestsMinDelay: plan.DigestsMinDelay.ValueInt64Pointer(),
		DigestsMaxDelay: plan.DigestsMaxDelay.ValueInt64Pointer(),
		ResolveAge:      plan.ResolveAge.ValueInt64Pointer(),
	}

	if plan.Slug.ValueString() != "" {
		params.Slug = plan.Slug.ValueStringPointer()
	}

	if plan.Filters != nil {
		options := make(map[string]interface{})
		if plan.Filters.BlacklistedIps.IsNull() {
			options["filters:blacklisted_ips"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(plan.Filters.BlacklistedIps.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:blacklisted_ips"] = strings.Join(values, "\n")
		}

		if plan.Filters.Releases.IsNull() {
			options["filters:releases"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(plan.Filters.Releases.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:releases"] = strings.Join(values, "\n")
		}

		if plan.Filters.ErrorMessages.IsNull() {
			options["filters:error_messages"] = ""
		} else {
			values := []string{}
			resp.Diagnostics.Append(plan.Filters.ErrorMessages.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:error_messages"] = strings.Join(values, "\n")
		}

		if len(options) > 0 {
			params.Options = &options
		}
	}

	if plan.FingerprintingRules.IsNull() {
		params.FingerprintingRules = sentry.String("")
	} else {
		params.FingerprintingRules = sentry.String(plan.FingerprintingRules.ValueString())
	}

	if plan.GroupingEnhancements.IsNull() {
		params.GroupingEnhancements = sentry.String("")
	} else {
		params.GroupingEnhancements = sentry.String(plan.GroupingEnhancements.ValueString())
	}

	httpRespUpdate, err := r.apiClient.UpdateOrganizationProjectWithResponse(
		ctx,
		plan.Organization.ValueString(),
		plan.Id.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpRespUpdate.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpRespUpdate.StatusCode() != http.StatusOK || httpRespUpdate.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpRespUpdate.StatusCode(), httpRespUpdate.Body))
		return
	}

	// If the default key is set to false, remove the default key
	if !plan.DefaultKey.IsNull() && !plan.DefaultKey.ValueBool() {
		if err := r.removeDefaultKey(ctx, plan.Organization.ValueString(), *httpRespUpdate.JSON200); err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("remove default key", err))
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
				httpResp, err := r.apiClient.AddTeamToProjectWithResponse(
					ctx,
					plan.Organization.ValueString(),
					httpRespUpdate.JSON200.Slug,
					team,
				)
				if err != nil {
					resp.Diagnostics.Append(diagutils.NewClientError("add team to project", err))
					return
				} else if httpResp.StatusCode() != http.StatusCreated {
					resp.Diagnostics.Append(diagutils.NewClientStatusError("add team to project", httpResp.StatusCode(), httpResp.Body))
					return
				}
			}
		}

		// Remove teams
		for _, team := range stateTeams {
			if !slices.Contains(planTeams, team) {
				httpResp, err := r.apiClient.RemoveTeamFromProjectWithResponse(
					ctx,
					plan.Organization.ValueString(),
					httpRespUpdate.JSON200.Slug,
					team,
				)
				if err != nil {
					resp.Diagnostics.Append(diagutils.NewClientError("remove team from project", err))
					return
				} else if httpResp.StatusCode() != http.StatusOK {
					resp.Diagnostics.Append(diagutils.NewClientStatusError("remove team from project", httpResp.StatusCode(), httpResp.Body))
					return
				}
			}
		}
	}

	httpRespRead, err := r.apiClient.GetOrganizationProjectWithResponse(
		ctx,
		plan.Organization.ValueString(),
		plan.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpRespRead.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpRespRead.StatusCode() != http.StatusOK || httpRespRead.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpRespRead.StatusCode(), httpRespRead.Body))
		return
	}

	if err := plan.Fill(*httpRespRead.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectResource) removeDefaultKey(ctx context.Context, organization string, project apiclient.Project) error {
	params := &apiclient.ListProjectClientKeysParams{}

	for {
		httpResp, err := r.apiClient.ListProjectClientKeysWithResponse(
			ctx,
			organization,
			project.Slug,
			params,
		)
		if err != nil {
			return err
		} else if httpResp.StatusCode() != http.StatusOK {
			return fmt.Errorf("unable to list project client keys, got status %d: %s", httpResp.StatusCode(), string(httpResp.Body))
		}

		for _, key := range *httpResp.JSON200 {
			if key.Name == "Default" {
				httpResp, err := r.apiClient.DeleteProjectClientKeyWithResponse(
					ctx,
					organization,
					project.Id,
					key.Id,
				)
				if err != nil {
					return err
				} else if httpResp.StatusCode() == http.StatusNotFound {
					return nil
				} else if httpResp.StatusCode() != http.StatusNoContent {
					return fmt.Errorf("unable to delete project client key, got status %d: %s", httpResp.StatusCode(), string(httpResp.Body))
				}

				return nil
			}
		}

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	return nil
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteOrganizationProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
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

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateTwoPartId(ctx, "organization", req, resp)
}
