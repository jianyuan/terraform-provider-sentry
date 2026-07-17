package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type ProjectOwnershipResourceModel struct {
	Organization       types.String              `tfsdk:"organization"`
	Project            types.String              `tfsdk:"project"`
	Raw                sentrytypes.TrimmedString `tfsdk:"raw"`
	Fallthrough        types.Bool                `tfsdk:"fallthrough"`
	AutoAssignment     types.String              `tfsdk:"auto_assignment"`
	CodeownersAutoSync types.Bool                `tfsdk:"codeowners_auto_sync"`
}

func (data *ProjectOwnershipResourceModel) Fill(ownership sentry.ProjectOwnership) error {
	data.Raw = sentrytypes.TrimmedStringValue(ownership.Raw)
	data.Fallthrough = types.BoolValue(ownership.FallThrough)
	data.AutoAssignment = types.StringValue(ownership.AutoAssignment)

	if ownership.CodeownersAutoSync == nil {
		data.CodeownersAutoSync = types.BoolValue(true)
	} else {
		data.CodeownersAutoSync = types.BoolValue(*ownership.CodeownersAutoSync)
	}

	return nil
}

var _ resource.Resource = &ProjectOwnershipResource{}
var _ resource.ResourceWithConfigure = &ProjectOwnershipResource{}
var _ resource.ResourceWithImportState = &ProjectOwnershipResource{}

func NewProjectOwnershipResource() resource.Resource {
	return &ProjectOwnershipResource{}
}

type ProjectOwnershipResource struct {
	baseResource
}

func (r *ProjectOwnershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_ownership"
}

func (r *ProjectOwnershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project Ownership. See the [Sentry documentation](https://docs.sentry.io/api/projects/update-ownership-configuration-for-a-project/) for more information.",

		Attributes: map[string]schema.Attribute{
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
			"raw": schema.StringAttribute{
				Description: "Raw input for ownership configuration.",
				CustomType:  sentrytypes.TrimmedStringType{},
				Required:    true,
			},
			"fallthrough": schema.BoolAttribute{
				Description: "Whether to fall through to the default ownership rules.",
				Required:    true,
			},
			"auto_assignment": schema.StringAttribute{
				Description: "The auto-assignment mode. The options are: `Auto Assign to Issue Owner`, `Auto Assign to Suspect Commits`, and `Turn off Auto-Assignment`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Auto Assign to Issue Owner",
						"Auto Assign to Suspect Commits",
						"Turn off Auto-Assignment",
					),
				},
			},
			"codeowners_auto_sync": schema.BoolAttribute{
				Description: "Whether to automatically sync codeowners.",
				Required:    true,
			},
		},
	}
}

func (r *ProjectOwnershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectOwnershipResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.UpdateProjectOwnershipParams{
		Raw:                data.Raw.ValueString(),
		FallThrough:        data.Fallthrough.ValueBoolPointer(),
		AutoAssignment:     data.AutoAssignment.ValueStringPointer(),
		CodeownersAutoSync: data.CodeownersAutoSync.ValueBoolPointer(),
	}

	source, _, err := r.client.ProjectOwnerships.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	}

	if err := data.Fill(*source); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectOwnershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectOwnershipResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownership, apiResp, err := r.client.ProjectOwnerships.Get(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project ownership"))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if ownership == nil {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project ownership"))
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(*ownership); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectOwnershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectOwnershipResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.UpdateProjectOwnershipParams{
		Raw:                data.Raw.ValueString(),
		FallThrough:        data.Fallthrough.ValueBoolPointer(),
		AutoAssignment:     data.AutoAssignment.ValueStringPointer(),
		CodeownersAutoSync: data.CodeownersAutoSync.ValueBoolPointer(),
	}

	source, _, err := r.client.ProjectOwnerships.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	}

	if err := data.Fill(*source); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectOwnershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectOwnershipResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	trueValue := true
	autoAssign := "Auto Assign to Issue Owner"
	_, apiResp, err := r.client.ProjectOwnerships.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		&sentry.UpdateProjectOwnershipParams{
			Raw:                "",
			FallThrough:        &trueValue,
			AutoAssignment:     &autoAssign,
			CodeownersAutoSync: &trueValue,
		},
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	}
}

func (r *ProjectOwnershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateTwoPart(ctx, "organization", "project", req, resp)
}
