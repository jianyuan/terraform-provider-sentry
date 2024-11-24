package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

type ProjectInboundDataFilterResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
	FilterId     types.String `tfsdk:"filter_id"`
	Active       types.Bool   `tfsdk:"active"`
	Subfilters   types.Set    `tfsdk:"subfilters"`
}

func (m *ProjectInboundDataFilterResourceModel) Fill(organization string, project string, filterId string, filter sentry.ProjectInboundDataFilter) error {
	m.Id = types.StringValue(buildThreePartID(organization, project, filterId))
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(project)
	m.FilterId = types.StringValue(filterId)

	if filter.Active.IsBool {
		m.Active = types.BoolValue(filter.Active.BoolVal)
	} else {
		subfilterElements := []attr.Value{}
		for _, subfilter := range filter.Active.StringSliceVal {
			subfilterElements = append(subfilterElements, types.StringValue(subfilter))
		}

		m.Subfilters = types.SetValueMust(types.StringType, subfilterElements)
	}

	return nil
}

var _ resource.Resource = &ProjectInboundDataFilterResource{}
var _ resource.ResourceWithConfigure = &ProjectInboundDataFilterResource{}
var _ resource.ResourceWithImportState = &ProjectInboundDataFilterResource{}

func NewProjectInboundDataFilterResource() resource.Resource {
	return &ProjectInboundDataFilterResource{}
}

type ProjectInboundDataFilterResource struct {
	baseResource
}

func (r *ProjectInboundDataFilterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_inbound_data_filter"
}

func (r *ProjectInboundDataFilterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project Inbound Data Filter resource. This resource is used to create and manage inbound data filters for a project. For more information on what filters are available, see the [Sentry documentation](https://docs.sentry.io/api/projects/update-an-inbound-data-filter/).",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
			"filter_id": schema.StringAttribute{
				Description: "The type of filter toggle to update. See the [Sentry documentation](https://docs.sentry.io/api/projects/update-an-inbound-data-filter/) for a list of available filters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
			"subfilters": schema.SetAttribute{
				Description: "Specifies which legacy browser filters should be active. Anything excluded from the list will be disabled. See the [Sentry documentation](https://docs.sentry.io/api/projects/update-an-inbound-data-filter/) for a list of available subfilters.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("active"),
					),
				},
			},
		},
	}
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
		diagutils.AddClientError(resp.Diagnostics, "create", err)
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

	filters, apiResp, err := r.client.ProjectInboundDataFilters.List(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		diagutils.AddNotFoundError(resp.Diagnostics, "project inbound data filter")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "read", err)
		return
	}

	var foundFilter *sentry.ProjectInboundDataFilter

	for _, filter := range filters {
		if filter.ID == data.FilterId.ValueString() {
			foundFilter = filter
			break
		}
	}

	if foundFilter == nil {
		diagutils.AddNotFoundError(resp.Diagnostics, "project inbound data filter")
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.Project.ValueString(), data.FilterId.ValueString(), *foundFilter); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
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
		diagutils.AddClientError(resp.Diagnostics, "update", err)
		return
	}

	plan.Id = types.StringValue(buildThreePartID(plan.Organization.ValueString(), plan.Project.ValueString(), plan.FilterId.ValueString()))

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
		diagutils.AddClientError(resp.Diagnostics, "delete", err)
		return
	}
}

func (r *ProjectInboundDataFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, filterID, err := splitThreePartID(req.ID, "organization", "project-slug", "filter-id")
	if err != nil {
		diagutils.AddImportError(resp.Diagnostics, err)
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
