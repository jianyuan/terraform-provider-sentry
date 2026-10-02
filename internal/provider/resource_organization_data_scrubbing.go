package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	intresource "github.com/jianyuan/terraform-provider-sentry/internal/resource"
	"github.com/samber/lo"
)

type OrganizationDataScrubbingResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Organization         types.String `tfsdk:"organization"`
	DataScrubber         types.Bool   `tfsdk:"data_scrubber"`
	DataScrubberDefaults types.Bool   `tfsdk:"data_scrubber_defaults"`
	SensitiveFields      types.Set    `tfsdk:"sensitive_fields"`
	SafeFields           types.Set    `tfsdk:"safe_fields"`
	ScrubIpAddresses     types.Bool   `tfsdk:"scrub_ip_addresses"`
}

func (data *OrganizationDataScrubbingResourceModel) Fill(ctx context.Context, organization apiclient.Organization) (diags diag.Diagnostics) {
	// organization keeps the configured value, which may be the internal ID, so the apply result matches the plan.
	data.Id = types.StringValue(organization.Slug)
	data.DataScrubber = types.BoolValue(lo.FromPtr(organization.DataScrubber))
	data.DataScrubberDefaults = types.BoolValue(lo.FromPtr(organization.DataScrubberDefaults))
	data.ScrubIpAddresses = types.BoolValue(lo.FromPtr(organization.ScrubIPAddresses))

	var setDiags diag.Diagnostics
	data.SensitiveFields, setDiags = stringSetValue(ctx, organization.SensitiveFields)
	diags.Append(setDiags...)
	data.SafeFields, setDiags = stringSetValue(ctx, organization.SafeFields)
	diags.Append(setDiags...)

	return
}

// stringSetValue returns an empty set for a missing list, because a null set would differ from the empty default.
func stringSetValue(ctx context.Context, values *[]string) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, append([]string{}, lo.FromPtr(values)...))
}

func (data OrganizationDataScrubbingResourceModel) UpdateBody(ctx context.Context) (body apiclient.UpdateOrganizationJSONRequestBody, diags diag.Diagnostics) {
	// Sentry rejects null for list fields, so the slices must be non-nil even when the set is empty.
	sensitiveFields := []string{}
	diags.Append(data.SensitiveFields.ElementsAs(ctx, &sensitiveFields, false)...)
	safeFields := []string{}
	diags.Append(data.SafeFields.ElementsAs(ctx, &safeFields, false)...)

	body = apiclient.UpdateOrganizationJSONRequestBody{
		DataScrubber:         data.DataScrubber.ValueBoolPointer(),
		DataScrubberDefaults: data.DataScrubberDefaults.ValueBoolPointer(),
		SensitiveFields:      &sensitiveFields,
		SafeFields:           &safeFields,
		ScrubIPAddresses:     data.ScrubIpAddresses.ValueBoolPointer(),
	}
	return
}

// organizationDataScrubbingDefaults holds the Sentry defaults, which are both the schema defaults and the settings that destroy restores.
func organizationDataScrubbingDefaults() OrganizationDataScrubbingResourceModel {
	emptyStringSet := types.SetValueMust(types.StringType, []attr.Value{})
	return OrganizationDataScrubbingResourceModel{
		DataScrubber:         types.BoolValue(false),
		DataScrubberDefaults: types.BoolValue(false),
		SensitiveFields:      emptyStringSet,
		SafeFields:           emptyStringSet,
		ScrubIpAddresses:     types.BoolValue(false),
	}
}

var _ resource.Resource = &OrganizationDataScrubbingResource{}
var _ resource.ResourceWithConfigure = &OrganizationDataScrubbingResource{}
var _ resource.ResourceWithImportState = &OrganizationDataScrubbingResource{}

func NewOrganizationDataScrubbingResource() resource.Resource {
	return &OrganizationDataScrubbingResource{}
}

type OrganizationDataScrubbingResource struct {
	baseResource
}

func (r *OrganizationDataScrubbingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_data_scrubbing"
}

func (r *OrganizationDataScrubbingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	defaults := organizationDataScrubbingDefaults()
	nonEmptyStrings := []validator.Set{
		setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Organization Data Scrubbing resource. This resource manages the server-side data scrubbing settings of an existing organization. Omitted attributes are reset to the Sentry defaults, and destroying the resource resets all settings to the Sentry defaults.",

		Attributes: map[string]schema.Attribute{
			"id": ResourceIdAttribute(),
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug or internal ID of the organization.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_scrubber": schema.BoolAttribute{
				MarkdownDescription: "Require server-side data scrubbing for all projects. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(defaults.DataScrubber.ValueBool()),
			},
			"data_scrubber_defaults": schema.BoolAttribute{
				MarkdownDescription: "Apply the default scrubbers to prevent things like passwords and credit cards from being stored for all projects. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(defaults.DataScrubberDefaults.ValueBool()),
			},
			"sensitive_fields": schema.SetAttribute{
				MarkdownDescription: "Additional field names to match against when scrubbing data for all projects. Defaults to an empty set.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(defaults.SensitiveFields),
				Validators:          nonEmptyStrings,
			},
			"safe_fields": schema.SetAttribute{
				MarkdownDescription: "Field names which data scrubbers should ignore for all projects. Defaults to an empty set.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(defaults.SafeFields),
				Validators:          nonEmptyStrings,
			},
			"scrub_ip_addresses": schema.BoolAttribute{
				MarkdownDescription: "Prevent IP addresses from being stored for new events in all projects. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(defaults.ScrubIpAddresses.ValueBool()),
			},
		},
	}
}

func (r *OrganizationDataScrubbingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationDataScrubbingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(ctx, &data, "create", &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationDataScrubbingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationDataScrubbingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetOrganizationWithResponse(ctx, data.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("organization"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationDataScrubbingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OrganizationDataScrubbingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(ctx, &data, "update", &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationDataScrubbingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationDataScrubbingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, bodyDiags := organizationDataScrubbingDefaults().UpdateBody(ctx)
	resp.Diagnostics.Append(bodyDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UpdateOrganizationWithResponse(ctx, data.Organization.ValueString(), body)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *OrganizationDataScrubbingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	intresource.ImportState1PartPassthrough("organization")(ctx, req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// update is shared by Create and Update because the organization already exists, so both write the settings with the same request.
func (r *OrganizationDataScrubbingResource) update(ctx context.Context, data *OrganizationDataScrubbingResourceModel, action string, diags *diag.Diagnostics) {
	body, bodyDiags := data.UpdateBody(ctx)
	diags.Append(bodyDiags...)
	if diags.HasError() {
		return
	}

	httpResp, err := r.apiClient.UpdateOrganizationWithResponse(ctx, data.Organization.ValueString(), body)
	if err != nil {
		diags.Append(diagutils.NewClientError(action, err))
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		diags.Append(diagutils.NewClientStatusError(action, httpResp.StatusCode(), httpResp.Body))
		return
	}

	diags.Append(data.Fill(ctx, *httpResp.JSON200)...)
}
