package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	intresource "github.com/jianyuan/terraform-provider-sentry/internal/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
)

type ProjectCodeOwnersResourceModel struct {
	Id            types.String              `tfsdk:"id"`
	Organization  types.String              `tfsdk:"organization"`
	Project       types.String              `tfsdk:"project"`
	Raw           sentrytypes.TrimmedString `tfsdk:"raw"`
	CodeMappingId types.String              `tfsdk:"code_mapping_id"`
}

func (m *ProjectCodeOwnersResourceModel) Fill(ctx context.Context, codeOwners apiclient.ProjectCodeOwners) (diags diag.Diagnostics) {
	m.Id = types.StringValue(codeOwners.Id)
	m.Raw = sentrytypes.TrimmedStringValue(codeOwners.Raw)
	m.CodeMappingId = types.StringValue(codeOwners.CodeMappingId)

	return
}

var _ resource.Resource = &ProjectCodeOwnersResource{}
var _ resource.ResourceWithConfigure = &ProjectCodeOwnersResource{}
var _ resource.ResourceWithImportState = &ProjectCodeOwnersResource{}

func NewProjectCodeOwnersResource() resource.Resource {
	return &ProjectCodeOwnersResource{}
}

type ProjectCodeOwnersResource struct {
	baseResource
}

func (r *ProjectCodeOwnersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_codeowners"
}

func (r *ProjectCodeOwnersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project CODEOWNERS. Imports a CODEOWNERS file for a project using an existing code mapping. A project can have more than one CODEOWNERS configuration, one per code mapping.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
			"raw": schema.StringAttribute{
				MarkdownDescription: "The raw contents of the CODEOWNERS file.",
				CustomType:          sentrytypes.TrimmedStringType{},
				Required:            true,
			},
			"code_mapping_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the code mapping used to translate repository paths to stack trace paths.",
				Required:            true,
			},
		},
	}
}

func (r *ProjectCodeOwnersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectCodeOwnersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CreateProjectCodeOwnersWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		apiclient.CreateProjectCodeOwnersJSONRequestBody{
			Raw:           data.Raw.ValueString(),
			CodeMappingId: data.CodeMappingId.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON201)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectCodeOwnersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectCodeOwnersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetProjectCodeOwnersWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
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

func (r *ProjectCodeOwnersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectCodeOwnersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UpdateProjectCodeOwnersWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
		apiclient.UpdateProjectCodeOwnersJSONRequestBody{
			Raw:           data.Raw.ValueStringPointer(),
			CodeMappingId: data.CodeMappingId.ValueStringPointer(),
		},
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectCodeOwnersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectCodeOwnersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteProjectCodeOwnersWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() != http.StatusNoContent && httpResp.StatusCode() != http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *ProjectCodeOwnersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	intresource.ImportState3PartPath("organization", "project", "id")(ctx, req, resp)
}
