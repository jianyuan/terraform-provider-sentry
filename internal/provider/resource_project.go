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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jianyuan/go-utils/sliceutils"

	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type ProjectFilterResourceModel struct {
	BlacklistedIps types.Set `tfsdk:"blacklisted_ips"`
	Releases       types.Set `tfsdk:"releases"`
	ErrorMessages  types.Set `tfsdk:"error_messages"`
}

func (m ProjectFilterResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"blacklisted_ips": types.SetType{ElemType: types.StringType},
		"releases":        types.SetType{ElemType: types.StringType},
		"error_messages":  types.SetType{ElemType: types.StringType},
	}
}

func (m *ProjectFilterResourceModel) Fill(ctx context.Context, project apiclient.Project) (diags diag.Diagnostics) {
	if project.Options == nil {
		m.BlacklistedIps = types.SetNull(types.StringType)
		m.Releases = types.SetNull(types.StringType)
		m.ErrorMessages = types.SetNull(types.StringType)

		return
	}

	if values, ok := project.Options["filters:blacklisted_ips"].(string); ok {
		if values == "" {
			m.BlacklistedIps = types.SetValueMust(types.StringType, []attr.Value{})
		} else {
			m.BlacklistedIps = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		diags.Append(diagutils.NewFillError(fmt.Errorf("invalid type for filters:blacklisted_ips: %T", project.Options["filters:blacklisted_ips"])))
	}

	if values, ok := project.Options["filters:releases"].(string); ok {
		if values == "" {
			m.Releases = types.SetValueMust(types.StringType, []attr.Value{})
		} else {
			m.Releases = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		diags.Append(diagutils.NewFillError(fmt.Errorf("invalid type for filters:releases: %T", project.Options["filters:releases"])))
	}

	if values, ok := project.Options["filters:error_messages"].(string); ok {
		if values == "" {
			m.ErrorMessages = types.SetValueMust(types.StringType, []attr.Value{})
		} else {
			m.ErrorMessages = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
				return types.StringValue(v)
			}, strings.Split(values, "\n")))
		}
	} else {
		diags.Append(diagutils.NewFillError(fmt.Errorf("invalid type for filters:error_messages: %T", project.Options["filters:error_messages"])))
	}

	return
}

type ProjectClientSecurityResourceModel struct {
	AllowedDomains      types.Set    `tfsdk:"allowed_domains"`
	ScrapeJavascript    types.Bool   `tfsdk:"scrape_javascript"`
	SecurityToken       types.String `tfsdk:"security_token"`
	SecurityTokenHeader types.String `tfsdk:"security_token_header"`
	VerifyTlsSsl        types.Bool   `tfsdk:"verify_tls_ssl"`
}

func (m ProjectClientSecurityResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_domains":       types.SetType{ElemType: types.StringType},
		"scrape_javascript":     types.BoolType,
		"security_token":        types.StringType,
		"security_token_header": types.StringType,
		"verify_tls_ssl":        types.BoolType,
	}
}

func (m *ProjectClientSecurityResourceModel) Fill(ctx context.Context, project apiclient.Project) (diags diag.Diagnostics) {
	m.AllowedDomains = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
		return types.StringValue(v)
	}, project.AllowedDomains))
	m.ScrapeJavascript = types.BoolValue(project.ScrapeJavaScript)
	m.SecurityToken = types.StringValue(project.SecurityToken)

	if v, err := project.SecurityTokenHeader.Get(); err == nil {
		m.SecurityTokenHeader = types.StringValue(v)
	} else {
		m.SecurityTokenHeader = types.StringValue("")
	}

	m.VerifyTlsSsl = types.BoolValue(project.VerifySSL)

	return
}

type ProjectResourceModel struct {
	Id                   types.String              `tfsdk:"id"`
	Organization         types.String              `tfsdk:"organization"`
	Teams                types.Set                 `tfsdk:"teams"`
	Name                 types.String              `tfsdk:"name"`
	Slug                 types.String              `tfsdk:"slug"`
	Platform             types.String              `tfsdk:"platform"`
	DefaultRules         types.Bool                `tfsdk:"default_rules"`
	DefaultKey           types.Bool                `tfsdk:"default_key"`
	InternalId           types.String              `tfsdk:"internal_id"`
	Features             types.Set                 `tfsdk:"features"`
	DigestsMinDelay      types.Int64               `tfsdk:"digests_min_delay"`
	DigestsMaxDelay      types.Int64               `tfsdk:"digests_max_delay"`
	ResolveAge           types.Int64               `tfsdk:"resolve_age"`
	Filters              types.Object              `tfsdk:"filters"`
	FingerprintingRules  sentrytypes.TrimmedString `tfsdk:"fingerprinting_rules"`
	GroupingEnhancements sentrytypes.TrimmedString `tfsdk:"grouping_enhancements"`
	ClientSecurity       types.Object              `tfsdk:"client_security"`
	HighlightTags        types.Set                 `tfsdk:"highlight_tags"`
}

func (m *ProjectResourceModel) Fill(ctx context.Context, project apiclient.Project) (diags diag.Diagnostics) {
	m.Id = types.StringValue(project.Slug)
	m.Organization = types.StringValue(project.Organization.Slug)
	m.Teams = types.SetValueMust(types.StringType, sliceutils.Map(func(v apiclient.Team) attr.Value {
		return types.StringValue(v.Slug)
	}, project.Teams))
	m.Name = types.StringValue(project.Name)
	m.Slug = types.StringValue(project.Slug)

	if v, err := project.Platform.Get(); err == nil && v != "" {
		m.Platform = types.StringValue(v)
	} else {
		m.Platform = types.StringNull()
	}

	m.InternalId = types.StringValue(project.Id)
	m.Features = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
		return types.StringValue(v)
	}, project.Features))

	m.DigestsMinDelay = types.Int64Value(project.DigestsMinDelay)
	m.DigestsMaxDelay = types.Int64Value(project.DigestsMaxDelay)
	m.ResolveAge = types.Int64Value(project.ResolveAge)

	var filters ProjectFilterResourceModel
	diags.Append(filters.Fill(ctx, project)...)
	m.Filters = tfutils.MergeDiagnostics(types.ObjectValueFrom(ctx, filters.AttributeTypes(), filters))(&diags)

	m.FingerprintingRules = sentrytypes.TrimmedStringValue(project.FingerprintingRules)
	m.GroupingEnhancements = sentrytypes.TrimmedStringValue(project.GroupingEnhancements)

	var clientSecurity ProjectClientSecurityResourceModel
	diags.Append(clientSecurity.Fill(ctx, project)...)
	m.ClientSecurity = tfutils.MergeDiagnostics(types.ObjectValueFrom(ctx, clientSecurity.AttributeTypes(), clientSecurity))(&diags)

	if project.HighlightTags != nil {
		m.HighlightTags = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
			return types.StringValue(v)
		}, *project.HighlightTags))
	} else {
		m.HighlightTags = types.SetNull(types.StringType)
	}

	return
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"platform": tfutils.WithEnumStringAttribute(schema.StringAttribute{
				MarkdownDescription: "The platform for this project. Use `other` for platforms not listed.",
				Optional:            true,
			}, sentrydata.Platforms),
			"default_rules": schema.BoolAttribute{
				Description: "Whether to create a default issue alert. Defaults to true where the behavior is to alert the user on every new issue.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_key": schema.BoolAttribute{
				Description: "Whether to create a default key on project creation. By default, Sentry will create a key for you. If you wish to manage keys manually, set this to false and create keys using the `sentry_key` resource. Note that this only takes effect on project creation, not on project update.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"internal_id": schema.StringAttribute{
				Description: "The internal ID for this project.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"features": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"digests_min_delay": schema.Int64Attribute{
				Description: "The minimum amount of time (in seconds) to wait between scheduling digests for delivery after the initial scheduling.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"digests_max_delay": schema.Int64Attribute{
				Description: "The maximum amount of time (in seconds) to wait between scheduling digests for delivery.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"resolve_age": schema.Int64Attribute{
				Description: "Hours in which an issue is automatically resolve if not seen after this amount of time.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"filters": schema.SingleNestedAttribute{
				Description: "Custom filters for this project.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"blacklisted_ips": schema.SetAttribute{
						Description: "Filter events from these IP addresses. (e.g. 127.0.0.1 or 10.0.0.0/8)",
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
					"releases": schema.SetAttribute{
						MarkdownDescription: "Filter events from these releases. Allows [glob pattern matching](https://en.wikipedia.org/wiki/Glob_(programming)). (e.g. 1.* or [!3].[0-9].*)",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
					"error_messages": schema.SetAttribute{
						MarkdownDescription: "Filter events by error messages. Allows [glob pattern matching](https://en.wikipedia.org/wiki/Glob_(programming)). (e.g. TypeError* or *: integer division or modulo by zero)",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"fingerprinting_rules": schema.StringAttribute{
				MarkdownDescription: "This can be used to modify the fingerprint rules on the server with custom rules. Rules follow the pattern `matcher:glob -> fingerprint, values`. To learn more about fingerprint rules, [read the docs](https://docs.sentry.io/concepts/data-management/event-grouping/fingerprint-rules/).",
				CustomType:          sentrytypes.TrimmedStringType{},
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"grouping_enhancements": schema.StringAttribute{
				MarkdownDescription: "This can be used to enhance the grouping algorithm with custom rules. Rules follow the pattern `matcher:glob [v^]?[+-]flag`. To learn more about stack trace rules, [read the docs](https://docs.sentry.io/concepts/data-management/event-grouping/stack-trace-rules/).",
				CustomType:          sentrytypes.TrimmedStringType{},
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_security": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure origin URLs which Sentry should accept events from. This is used for communication with clients like [sentry-javascript](https://github.com/getsentry/sentry-javascript).",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"allowed_domains": schema.SetAttribute{
						MarkdownDescription: "A list of allowed domains. Examples: https://example.com, *, *.example.com, *:80.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
					"scrape_javascript": schema.BoolAttribute{
						MarkdownDescription: "Enable JavaScript source fetching. Allow Sentry to scrape missing JavaScript source context when possible.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"security_token": schema.StringAttribute{
						MarkdownDescription: "Security Token. Outbound requests matching Allowed Domains will have the header \"{security_token_header}: {security_token}\" appended.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"security_token_header": schema.StringAttribute{
						MarkdownDescription: "Security Token Header. Outbound requests matching Allowed Domains will have the header \"{security_token_header}: {security_token}\" appended.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(20),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"verify_tls_ssl": schema.BoolAttribute{
						MarkdownDescription: "Verify TLS/SSL. Outbound requests will verify TLS (sometimes known as SSL) connections.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"highlight_tags": schema.SetAttribute{
				MarkdownDescription: "A list of strings with tag keys to highlight on this project's issues. E.g. ['release', 'environment']",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
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
	createBody := apiclient.CreateOrganizationTeamProjectJSONRequestBody{
		Name:         data.Name.ValueString(),
		Platform:     data.Platform.ValueStringPointer(),
		DefaultRules: data.DefaultRules.ValueBoolPointer(),
	}

	if !data.Slug.IsUnknown() {
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
		Name: data.Name.ValueStringPointer(),
	}

	if !data.Slug.IsUnknown() {
		updateBody.Slug = data.Slug.ValueStringPointer()
	}

	if !data.Platform.IsUnknown() {
		updateBody.Platform = data.Platform.ValueStringPointer()
	}

	if !data.DigestsMinDelay.IsUnknown() {
		updateBody.DigestsMinDelay = data.DigestsMinDelay.ValueInt64Pointer()
	}

	if !data.DigestsMaxDelay.IsUnknown() {
		updateBody.DigestsMaxDelay = data.DigestsMaxDelay.ValueInt64Pointer()
	}

	if !data.ResolveAge.IsUnknown() {
		updateBody.ResolveAge = data.ResolveAge.ValueInt64Pointer()
	}

	if !data.Filters.IsUnknown() {
		var filters ProjectFilterResourceModel
		resp.Diagnostics.Append(data.Filters.As(ctx, &filters, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		options := make(map[string]interface{})
		if !filters.BlacklistedIps.IsUnknown() {
			var values []string
			resp.Diagnostics.Append(filters.BlacklistedIps.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:blacklisted_ips"] = strings.Join(values, "\n")
		}

		if !filters.Releases.IsUnknown() {
			var values []string
			resp.Diagnostics.Append(filters.Releases.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:releases"] = strings.Join(values, "\n")
		}

		if !filters.ErrorMessages.IsUnknown() {
			var values []string
			resp.Diagnostics.Append(filters.ErrorMessages.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:error_messages"] = strings.Join(values, "\n")
		}

		if len(options) > 0 {
			updateBody.Options = &options
		}
	}

	if !data.FingerprintingRules.IsUnknown() {
		updateBody.FingerprintingRules = data.FingerprintingRules.ValueStringPointer()
	}

	if !data.GroupingEnhancements.IsUnknown() {
		updateBody.GroupingEnhancements = data.GroupingEnhancements.ValueStringPointer()
	}

	if !data.ClientSecurity.IsUnknown() {
		var clientSecurity ProjectClientSecurityResourceModel
		resp.Diagnostics.Append(data.ClientSecurity.As(ctx, &clientSecurity, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !clientSecurity.AllowedDomains.IsUnknown() {
			resp.Diagnostics.Append(clientSecurity.AllowedDomains.ElementsAs(ctx, &updateBody.AllowedDomains, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if !clientSecurity.ScrapeJavascript.IsUnknown() {
			updateBody.ScrapeJavaScript = clientSecurity.ScrapeJavascript.ValueBoolPointer()
		}

		if !clientSecurity.SecurityToken.IsUnknown() {
			updateBody.SecurityToken = clientSecurity.SecurityToken.ValueStringPointer()
		}

		if !clientSecurity.SecurityTokenHeader.IsUnknown() {
			updateBody.SecurityTokenHeader = clientSecurity.SecurityTokenHeader.ValueStringPointer()
		}

		if !clientSecurity.VerifyTlsSsl.IsUnknown() {
			updateBody.VerifySSL = clientSecurity.VerifyTlsSsl.ValueBoolPointer()
		}
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

	resp.Diagnostics.Append(data.Fill(ctx, *httpRespRead.JSON200)...)
	if resp.Diagnostics.HasError() {
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

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateBody := apiclient.UpdateOrganizationProjectJSONRequestBody{}

	if !plan.Name.Equal(state.Name) {
		updateBody.Name = plan.Name.ValueStringPointer()
	}

	if !plan.Slug.Equal(state.Slug) {
		updateBody.Slug = plan.Slug.ValueStringPointer()
	}

	if !plan.Platform.Equal(state.Platform) {
		updateBody.Platform = plan.Platform.ValueStringPointer()
	}

	if !plan.DigestsMinDelay.Equal(state.DigestsMinDelay) {
		updateBody.DigestsMinDelay = plan.DigestsMinDelay.ValueInt64Pointer()
	}

	if !plan.DigestsMaxDelay.Equal(state.DigestsMaxDelay) {
		updateBody.DigestsMaxDelay = plan.DigestsMaxDelay.ValueInt64Pointer()
	}

	if !plan.ResolveAge.Equal(state.ResolveAge) {
		updateBody.ResolveAge = plan.ResolveAge.ValueInt64Pointer()
	}

	if !plan.Filters.Equal(state.Filters) {
		var filtersPlan, filtersState ProjectFilterResourceModel
		resp.Diagnostics.Append(plan.Filters.As(ctx, &filtersPlan, basetypes.ObjectAsOptions{})...)
		resp.Diagnostics.Append(state.Filters.As(ctx, &filtersState, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		options := make(map[string]interface{})
		if !filtersPlan.BlacklistedIps.Equal(filtersState.BlacklistedIps) {
			var values []string
			resp.Diagnostics.Append(filtersPlan.BlacklistedIps.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:blacklisted_ips"] = strings.Join(values, "\n")
		}

		if !filtersPlan.Releases.Equal(filtersState.Releases) {
			var values []string
			resp.Diagnostics.Append(filtersPlan.Releases.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:releases"] = strings.Join(values, "\n")
		}

		if !filtersPlan.ErrorMessages.Equal(filtersState.ErrorMessages) {
			var values []string
			resp.Diagnostics.Append(filtersPlan.ErrorMessages.ElementsAs(ctx, &values, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			options["filters:error_messages"] = strings.Join(values, "\n")
		}

		if len(options) > 0 {
			updateBody.Options = &options
		}
	}

	if !plan.FingerprintingRules.Equal(state.FingerprintingRules) {
		updateBody.FingerprintingRules = plan.FingerprintingRules.ValueStringPointer()
	}

	if !plan.GroupingEnhancements.Equal(state.GroupingEnhancements) {
		updateBody.GroupingEnhancements = plan.GroupingEnhancements.ValueStringPointer()
	}

	if !plan.ClientSecurity.Equal(state.ClientSecurity) {
		var clientSecurityPlan, clientSecurityState ProjectClientSecurityResourceModel
		resp.Diagnostics.Append(plan.ClientSecurity.As(ctx, &clientSecurityPlan, basetypes.ObjectAsOptions{})...)
		resp.Diagnostics.Append(state.ClientSecurity.As(ctx, &clientSecurityState, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !clientSecurityPlan.AllowedDomains.Equal(clientSecurityState.AllowedDomains) {
			resp.Diagnostics.Append(clientSecurityPlan.AllowedDomains.ElementsAs(ctx, &updateBody.AllowedDomains, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if !clientSecurityPlan.ScrapeJavascript.Equal(clientSecurityState.ScrapeJavascript) {
			updateBody.ScrapeJavaScript = clientSecurityPlan.ScrapeJavascript.ValueBoolPointer()
		}

		if !clientSecurityPlan.SecurityToken.Equal(clientSecurityState.SecurityToken) {
			updateBody.SecurityToken = clientSecurityPlan.SecurityToken.ValueStringPointer()
		}

		if !clientSecurityPlan.SecurityTokenHeader.Equal(clientSecurityState.SecurityTokenHeader) {
			updateBody.SecurityTokenHeader = clientSecurityPlan.SecurityTokenHeader.ValueStringPointer()
		}

		if !clientSecurityPlan.VerifyTlsSsl.Equal(clientSecurityState.VerifyTlsSsl) {
			updateBody.VerifySSL = clientSecurityPlan.VerifyTlsSsl.ValueBoolPointer()
		}
	}

	if !plan.HighlightTags.Equal(state.HighlightTags) {
		var highlightTags []string
		resp.Diagnostics.Append(plan.HighlightTags.ElementsAs(ctx, &highlightTags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateBody.HighlightTags = &highlightTags
	}

	httpRespUpdate, err := r.apiClient.UpdateOrganizationProjectWithResponse(
		ctx,
		plan.Organization.ValueString(),
		plan.Id.ValueString(),
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

	resp.Diagnostics.Append(plan.Fill(ctx, *httpRespRead.JSON200)...)
	if resp.Diagnostics.HasError() {
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
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
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
