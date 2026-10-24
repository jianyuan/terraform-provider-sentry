package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-plugin-framework-utils/fwdiag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	intresource "github.com/jianyuan/terraform-provider-sentry/internal/resource"
	"github.com/samber/lo"
)

type OrganizationDataScrubbingResourceModel struct {
	Id                   types.String         `tfsdk:"id"`
	Organization         types.String         `tfsdk:"organization"`
	DataScrubber         types.Bool           `tfsdk:"data_scrubber"`
	DataScrubberDefaults types.Bool           `tfsdk:"data_scrubber_defaults"`
	ScrubIpAddresses     types.Bool           `tfsdk:"scrub_ip_addresses"`
	SensitiveFields      types.Set            `tfsdk:"sensitive_fields"`
	SafeFields           types.Set            `tfsdk:"safe_fields"`
	RelayPiiConfig       jsontypes.Normalized `tfsdk:"relay_pii_config"`
}

// Fill keeps the configured organization, which may be the internal ID.
func (m *OrganizationDataScrubbingResourceModel) Fill(ctx context.Context, organization apiclient.Organization) (diags diag.Diagnostics) {
	m.Id = types.StringValue(organization.Slug)
	m.DataScrubber = types.BoolValue(lo.FromPtr(organization.DataScrubber))
	m.DataScrubberDefaults = types.BoolValue(lo.FromPtr(organization.DataScrubberDefaults))
	m.ScrubIpAddresses = types.BoolValue(lo.FromPtr(organization.ScrubIPAddresses))
	m.SensitiveFields = fwdiag.Merge(stringSetValue(ctx, organization.SensitiveFields))(&diags)
	m.SafeFields = fwdiag.Merge(stringSetValue(ctx, organization.SafeFields))(&diags)
	m.RelayPiiConfig = relayPiiConfigValue(organization.RelayPiiConfig)
	return
}

// UpdateBody writes all six settings. A null relay_pii_config clears the rules.
func (m OrganizationDataScrubbingResourceModel) UpdateBody(ctx context.Context) (body apiclient.UpdateOrganizationJSONRequestBody, diags diag.Diagnostics) {
	sensitiveFields := fwdiag.Merge(stringSetElements(ctx, m.SensitiveFields))(&diags)
	safeFields := fwdiag.Merge(stringSetElements(ctx, m.SafeFields))(&diags)

	body = apiclient.UpdateOrganizationJSONRequestBody{
		DataScrubber:         m.DataScrubber.ValueBoolPointer(),
		DataScrubberDefaults: m.DataScrubberDefaults.ValueBoolPointer(),
		ScrubIPAddresses:     m.ScrubIpAddresses.ValueBoolPointer(),
		SensitiveFields:      &sensitiveFields,
		SafeFields:           &safeFields,
		RelayPiiConfig:       new(m.RelayPiiConfig.ValueString()),
	}
	return
}

// organizationDataScrubbingDefaults holds the Sentry defaults: the schema defaults and what destroy restores.
func organizationDataScrubbingDefaults() OrganizationDataScrubbingResourceModel {
	emptyStringSet := types.SetValueMust(types.StringType, []attr.Value{})
	return OrganizationDataScrubbingResourceModel{
		DataScrubber:         types.BoolValue(false),
		DataScrubberDefaults: types.BoolValue(false),
		ScrubIpAddresses:     types.BoolValue(false),
		SensitiveFields:      emptyStringSet,
		SafeFields:           emptyStringSet,
		RelayPiiConfig:       jsontypes.NewNormalizedNull(),
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
		MarkdownDescription: "Sentry Organization Data Scrubbing resource. Manages the data scrubbing settings that apply to every project of an existing organization, including Advanced Data Scrubbing rules. Omitted attributes are reset to the Sentry defaults, and destroying the resource resets every setting to the Sentry defaults. Project-level settings are attributes of `sentry_project`.",

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
				MarkdownDescription: "Require the default scrubbers, which prevent things like passwords and credit cards from being stored, for all projects. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(defaults.DataScrubberDefaults.ValueBool()),
			},
			"scrub_ip_addresses": schema.BoolAttribute{
				MarkdownDescription: "Prevent IP addresses from being stored for new events in all projects. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(defaults.ScrubIpAddresses.ValueBool()),
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
			"relay_pii_config": schema.StringAttribute{
				MarkdownDescription: relayPiiConfigMarkdownDescription + " These rules apply to all projects. Omit it to remove all rules.",
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
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

	resp.Diagnostics.Append(r.write(ctx, &data, "create")...)
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

func (r *OrganizationDataScrubbingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OrganizationDataScrubbingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.write(ctx, &data, "update")...)
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

	body := fwdiag.Merge(organizationDataScrubbingDefaults().UpdateBody(ctx))(&resp.Diagnostics)
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
}

// write is shared by Create and Update: the organization already exists, so both send every setting and read them back.
func (r *OrganizationDataScrubbingResource) write(ctx context.Context, data *OrganizationDataScrubbingResourceModel, action string) (diags diag.Diagnostics) {
	body := fwdiag.Merge(data.UpdateBody(ctx))(&diags)
	if diags.HasError() {
		return
	}

	httpRespUpdate, err := r.apiClient.UpdateOrganizationWithResponse(ctx, data.Organization.ValueString(), body)
	if err != nil {
		diags.Append(diagutils.NewClientError(action, err))
		return
	} else if httpRespUpdate.StatusCode() != http.StatusOK {
		diags.Append(diagutils.NewClientStatusError(action, httpRespUpdate.StatusCode(), httpRespUpdate.Body))
		return
	}

	httpRespRead, err := r.apiClient.GetOrganizationWithResponse(ctx, data.Organization.ValueString())
	if err != nil {
		diags.Append(diagutils.NewClientError("read", err))
		return
	} else if httpRespRead.StatusCode() != http.StatusOK || httpRespRead.JSON200 == nil {
		diags.Append(diagutils.NewClientStatusError("read", httpRespRead.StatusCode(), httpRespRead.Body))
		return
	}

	diags.Append(data.Fill(ctx, *httpRespRead.JSON200)...)
	return
}
