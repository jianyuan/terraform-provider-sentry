package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &ClientKeyResource{}
var _ resource.ResourceWithConfigure = &ClientKeyResource{}
var _ resource.ResourceWithConfigValidators = &ClientKeyResource{}
var _ resource.ResourceWithImportState = &ClientKeyResource{}

func NewClientKeyResource() resource.Resource {
	return &ClientKeyResource{}
}

type ClientKeyResource struct {
	baseResource
}

type ClientKeyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Organization    types.String `tfsdk:"organization"`
	Project         types.String `tfsdk:"project"`
	ProjectId       types.String `tfsdk:"project_id"`
	Name            types.String `tfsdk:"name"`
	RateLimitWindow types.Int64  `tfsdk:"rate_limit_window"`
	RateLimitCount  types.Int64  `tfsdk:"rate_limit_count"`
	Public          types.String `tfsdk:"public"`
	Secret          types.String `tfsdk:"secret"`
	DsnPublic       types.String `tfsdk:"dsn_public"`
	DsnSecret       types.String `tfsdk:"dsn_secret"`
	DsnCsp          types.String `tfsdk:"dsn_csp"`
}

func (m *ClientKeyResourceModel) Fill(organization string, project string, key sentry.ProjectKey) error {
	m.Id = types.StringValue(key.ID)
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(project)

	if key.RateLimit == nil {
		m.RateLimitWindow = types.Int64Null()
		m.RateLimitCount = types.Int64Null()
	} else {
		m.RateLimitWindow = types.Int64Value(int64(key.RateLimit.Window))
		m.RateLimitCount = types.Int64Value(int64(key.RateLimit.Count))
	}

	m.ProjectId = types.StringValue(key.ProjectID.String())
	m.Name = types.StringValue(key.Name)
	m.Public = types.StringValue(key.Public)
	m.Secret = types.StringValue(key.Secret)
	m.DsnPublic = types.StringValue(key.DSN.Public)
	m.DsnSecret = types.StringValue(key.DSN.Secret)
	m.DsnCsp = types.StringValue(key.DSN.CSP)

	return nil
}

func (r *ClientKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *ClientKeyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.RequiredTogether(
			path.MatchRoot("rate_limit_window"),
			path.MatchRoot("rate_limit_count"),
		),
	}
}

func (r *ClientKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Return a client key bound to a project.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the resource belongs to.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The slug of the project the resource belongs to.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the client key.",
				Required:            true,
			},
			"rate_limit_window": schema.Int64Attribute{
				MarkdownDescription: "Length of time that will be considered when checking the rate limit.",
				Optional:            true,
			},
			"rate_limit_count": schema.Int64Attribute{
				MarkdownDescription: "Number of events that can be reported within the rate limit window.",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project that the key belongs to.",
				Computed:            true,
			},
			"public": schema.StringAttribute{
				MarkdownDescription: "The public key.",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "The secret key.",
				Computed:            true,
			},
			"dsn_public": schema.StringAttribute{
				MarkdownDescription: "The DSN tells the SDK where to send the events to.",
				Computed:            true,
			},
			"dsn_secret": schema.StringAttribute{
				MarkdownDescription: "Deprecated DSN includes a secret which is no longer required by newer SDK versions. If you are unsure which to use, follow installation instructions for your language.",
				Computed:            true,
			},
			"dsn_csp": schema.StringAttribute{
				MarkdownDescription: "Security header endpoint for features like CSP and Expect-CT reports.",
				Computed:            true,
			},
		},
	}
}

func (r *ClientKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClientKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.CreateProjectKeyParams{
		Name: data.Name.ValueString(),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: int(data.RateLimitWindow.ValueInt64()),
			Count:  int(data.RateLimitCount.ValueInt64()),
		},
	}

	key, _, err := r.client.ProjectKeys.Create(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create error: %s", err.Error()))
		return
	}
	if err := data.Fill(data.Organization.ValueString(), data.Project.ValueString(), *key); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClientKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, apiResp, err := r.client.ProjectKeys.Get(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Project.ValueString(), *key); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClientKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.UpdateProjectKeyParams{
		Name: data.Name.ValueString(),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: int(data.RateLimitWindow.ValueInt64()),
			Count:  int(data.RateLimitCount.ValueInt64()),
		},
	}

	key, apiResp, err := r.client.ProjectKeys.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
		params,
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Update error: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Project.ValueString(), *key); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClientKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.ProjectKeys.Delete(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Delete error: %s", err.Error()))
		return
	}
}

func (r *ClientKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, id, err := splitThreePartID(req.ID, "organization", "project-slug", "key-id")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("project"), project,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)...)
}
