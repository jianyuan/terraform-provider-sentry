package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &ProjectInboundDataFilterResource{}
var _ resource.ResourceWithImportState = &ProjectInboundDataFilterResource{}

func NewProjectInboundDataFilterResource() resource.Resource {
	return &ProjectInboundDataFilterResource{}
}

type ProjectInboundDataFilterResource struct {
	client *sentry.Client
}

type ProjectInboundDataFilterResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
	FilterId     types.String `tfsdk:"filter_id"`
	Active       types.Bool   `tfsdk:"active"`
	Subfilters   types.List   `tfsdk:"subfilters"`
}

func (r *ProjectInboundDataFilterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_inbound_data_filter"
}

func (r *ProjectInboundDataFilterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project Inbound Data Filter resource. This resource is used to create and manage inbound data filters for a project. For more information on what filters are available, see the [Sentry documentation](https://docs.sentry.io/api/projects/update-an-inbound-data-filter/).",

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
			"project": schema.StringAttribute{
				Description: "The slug of the project to create the filter for.",
				Required:    true,
			},
			"filter_id": schema.StringAttribute{
				Description: "The type of filter toggle to update. See the [Sentry documentation](https://docs.sentry.io/api/projects/update-an-inbound-data-filter/) for a list of available filters.",
				Required:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Toggle the browser-extensions, localhost, filtered-transaction, or web-crawlers filter on or off.",
				Optional:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("subfilters"),
					),
				},
			},
			"subfilters": schema.ListAttribute{
				Description: "Specifies which legacy browser filters should be active. Anything excluded from the list will be disabled. See the [Sentry documentation](https://docs.sentry.io/api/projects/update-an-inbound-data-filter/) for a list of available subfilters.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("active"),
					),
				},
			},
		},
	}
}

func (r *ProjectInboundDataFilterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sentry.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sentry.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ProjectInboundDataFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectInboundDataFilterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	subfilters := []string{}
	if !data.Subfilters.IsNull() {
		resp.Diagnostics.Append(data.Subfilters.ElementsAs(ctx, &subfilters, false)...)
	}

	_, err := r.client.ProjectInboundDataFilters.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.FilterId.ValueString(),
		&sentry.UpdateProjectInboundDataFilterParams{
			Active:     data.Active.ValueBoolPointer(),
			Subfilters: subfilters,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating project inbound data filter: %s", err.Error()))
		return
	}

	data.Id = types.StringValue(buildThreePartID(data.Organization.ValueString(), data.Project.ValueString(), data.FilterId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectInboundDataFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectInboundDataFilterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filters, _, err := r.client.ProjectInboundDataFilters.List(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading project inbound data filters: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	for _, filter := range filters {
		if filter.ID == data.FilterId.ValueString() {
			data.Id = types.StringValue(buildThreePartID(data.Organization.ValueString(), data.Project.ValueString(), data.FilterId.ValueString()))

			if filter.Active.IsBool {
				data.Active = types.BoolValue(filter.Active.BoolVal)
			} else {
				listValue, diags := types.ListValueFrom(ctx, types.StringType, filter.Active.SliceVal)
				data.Subfilters = listValue
				resp.Diagnostics.Append(diags...)
			}
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading project inbound data filters: %s", "Filter not found"))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectInboundDataFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectInboundDataFilterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	subfilters := []string{}
	if !plan.Subfilters.IsNull() {
		resp.Diagnostics.Append(plan.Subfilters.ElementsAs(ctx, &subfilters, false)...)
	}

	_, err := r.client.ProjectInboundDataFilters.Update(
		ctx,
		plan.Organization.ValueString(),
		plan.Project.ValueString(),
		plan.FilterId.ValueString(),
		&sentry.UpdateProjectInboundDataFilterParams{
			Active:     plan.Active.ValueBoolPointer(),
			Subfilters: subfilters,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating project inbound data filter: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectInboundDataFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectInboundDataFilterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.ProjectInboundDataFilters.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.FilterId.ValueString(),
		&sentry.UpdateProjectInboundDataFilterParams{
			Active: sentry.Bool(false),
		},
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting project inbound data filter: %s", err.Error()))
		return
	}
}

func (r *ProjectInboundDataFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, filterID, err := splitThreePartID(req.ID, "organization", "project-slug", "filter-id")
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to import team, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("project"), project,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("filter_id"), filterID,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}
