package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &ProjectSymbolSourcesResource{}
var _ resource.ResourceWithConfigure = &ProjectSymbolSourcesResource{}
var _ resource.ResourceWithImportState = &ProjectSymbolSourcesResource{}

func NewProjectSymbolSourcesResource() resource.Resource {
	return &ProjectSymbolSourcesResource{}
}

type ProjectSymbolSourcesResource struct {
	baseResource
}

type ProjectSymbolSourcesResourceModel struct {
	Id                   types.String                             `tfsdk:"id"`
	Organization         types.String                             `tfsdk:"organization"`
	Project              types.String                             `tfsdk:"project"`
	Type                 types.String                             `tfsdk:"type"`
	Name                 types.String                             `tfsdk:"name"`
	Layout               *ProjectSymbolSourcesResourceLayoutModel `tfsdk:"layout"`
	AppConnectIssuer     types.String                             `tfsdk:"app_connect_issuer"`
	AppConnectPrivateKey types.String                             `tfsdk:"app_connect_private_key"`
	AppId                types.String                             `tfsdk:"app_id"`
	Url                  types.String                             `tfsdk:"url"`
	Username             types.String                             `tfsdk:"username"`
	Password             types.String                             `tfsdk:"password"`
	Bucket               types.String                             `tfsdk:"bucket"`
	Region               types.String                             `tfsdk:"region"`
	AccessKey            types.String                             `tfsdk:"access_key"`
	SecretKey            types.String                             `tfsdk:"secret_key"`
	Prefix               types.String                             `tfsdk:"prefix"`
	ClientEmail          types.String                             `tfsdk:"client_email"`
	PrivateKey           types.String                             `tfsdk:"private_key"`
}

func (data *ProjectSymbolSourcesResourceModel) Fill(source sentry.ProjectSymbolSource) error {
	data.Id = types.StringPointerValue(source.ID)
	data.Type = types.StringPointerValue(source.Type)
	data.Name = types.StringPointerValue(source.Name)
	if source.Layout != nil {
		data.Layout = &ProjectSymbolSourcesResourceLayoutModel{}
		if err := data.Layout.Fill(*source.Layout); err != nil {
			return err
		}
	}
	data.AppConnectIssuer = types.StringPointerValue(source.AppConnectIssuer)
	data.AppId = types.StringPointerValue(source.AppId)
	data.Url = types.StringPointerValue(source.Url)
	data.Username = types.StringPointerValue(source.Username)
	data.Bucket = types.StringPointerValue(source.Bucket)
	data.Region = types.StringPointerValue(source.Region)
	data.AccessKey = types.StringPointerValue(source.AccessKey)
	data.Prefix = types.StringPointerValue(source.Prefix)
	data.ClientEmail = types.StringPointerValue(source.ClientEmail)

	return nil
}

type ProjectSymbolSourcesResourceLayoutModel struct {
	Type   types.String `tfsdk:"type"`
	Casing types.String `tfsdk:"casing"`
}

func (m *ProjectSymbolSourcesResourceLayoutModel) Fill(layout sentry.ProjectSymbolSourceLayout) error {
	m.Type = types.StringPointerValue(layout.Type)
	m.Casing = types.StringPointerValue(layout.Casing)

	return nil
}

func (m ProjectSymbolSourcesResourceLayoutModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":   types.StringType,
		"casing": types.StringType,
	}
}

func (r *ProjectSymbolSourcesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_symbol_source"
}

func (r *ProjectSymbolSourcesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Project Symbol Source. See the [Sentry documentation](https://docs.sentry.io/api/projects/add-a-symbol-source-to-a-project/) for more information.",

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
			"type": schema.StringAttribute{
				Description: "The type of symbol source. One of `appStoreConnect` (App Store Connect), `http` (SymbolServer (HTTP)), `gcs` (Google Cloud Storage), `s3` (Amazon S3).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"appStoreConnect",
						"http",
						"gcs",
						"s3",
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "The human-readable name of the source.",
				Required:    true,
			},
			"layout": schema.SingleNestedAttribute{
				Description: "Layout settings for the source. This is required for HTTP, GCS, and S3 sources and invalid for AppStoreConnect sources.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The layout of the folder structure. The options are: `native` - Platform-Specific (SymStore / GDB / LLVM), `symstore` - Microsoft SymStore, `symstore_index2` - Microsoft SymStore (with index2.txt), `ssqp` - Microsoft SSQP, `unified` - Unified Symbol Server Layout, `debuginfod` - debuginfod.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"native",
								"symstore",
								"symstore_index2",
								"ssqp",
								"unified",
								"debuginfod",
							),
						},
					},
					"casing": schema.StringAttribute{
						Description: "The casing of the symbol source layout. The layout of the folder structure. The options are: `default` - Default (mixed case), `uppercase` - Uppercase, `lowercase` - Lowercase.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"default",
								"uppercase",
								"lowercase",
							),
						},
					},
				},
			},
			"app_connect_issuer": schema.StringAttribute{
				Description: "The App Store Connect Issuer ID. Required for AppStoreConnect sources, invalid for all others.",
				Optional:    true,
			},
			"app_connect_private_key": schema.StringAttribute{
				Description: "The App Store Connect API Private Key. Required for AppStoreConnect sources, invalid for all others.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Description: "The App Store Connect App ID. Required for AppStoreConnect sources, invalid for all others.",
				Optional:    true,
			},
			"url": schema.StringAttribute{
				Description: "The source's URL. Optional for HTTP sources, invalid for all others.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The user name for accessing the source. Optional for HTTP sources, invalid for all others.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for accessing the source. Optional for HTTP sources, invalid for all others.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket": schema.StringAttribute{
				Description: "The GCS or S3 bucket where the source resides. Required for GCS and S3 sourcse, invalid for HTTP and AppStoreConnect sources.",
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Description: "The source's S3 region. Required for S3 sources, invalid for all others.",
				Optional:    true,
			},
			"access_key": schema.StringAttribute{
				Description: "The AWS Access Key.Required for S3 sources, invalid for all others.",
				Optional:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "The AWS Secret Access Key.Required for S3 sources, invalid for all others.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prefix": schema.StringAttribute{
				Description: "The GCS or S3 prefix. Optional for GCS and S3 sourcse, invalid for HTTP and AppStoreConnect sources.",
				Optional:    true,
			},
			"client_email": schema.StringAttribute{
				Description: "The GCS email address for authentication. Required for GCS sources, invalid for all others.",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The GCS private key. Required for GCS sources, invalid for all others.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ProjectSymbolSourcesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectSymbolSourcesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.CreateProjectSymbolSourceParams{
		Type:                 data.Type.ValueStringPointer(),
		Name:                 data.Name.ValueStringPointer(),
		AppConnectIssuer:     data.AppConnectIssuer.ValueStringPointer(),
		AppConnectPrivateKey: data.AppConnectPrivateKey.ValueStringPointer(),
		AppId:                data.AppId.ValueStringPointer(),
		Url:                  data.Url.ValueStringPointer(),
		Username:             data.Username.ValueStringPointer(),
		Password:             data.Password.ValueStringPointer(),
		Bucket:               data.Bucket.ValueStringPointer(),
		Region:               data.Region.ValueStringPointer(),
		AccessKey:            data.AccessKey.ValueStringPointer(),
		SecretKey:            data.SecretKey.ValueStringPointer(),
		Prefix:               data.Prefix.ValueStringPointer(),
		ClientEmail:          data.ClientEmail.ValueStringPointer(),
		PrivateKey:           data.PrivateKey.ValueStringPointer(),
	}
	if data.Layout != nil {
		params.Layout = &sentry.ProjectSymbolSourceLayout{
			Type:   data.Layout.Type.ValueStringPointer(),
			Casing: data.Layout.Casing.ValueStringPointer(),
		}
	}

	source, _, err := r.client.ProjectSymbolSources.Create(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating project symbol source: %s", err.Error()))
		return
	}

	if err := data.Fill(*source); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling project symbol source: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectSymbolSourcesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectSymbolSourcesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sources, apiResp, err := r.client.ProjectSymbolSources.List(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		&sentry.ProjectSymbolSourceQueryParams{
			ID: data.Id.ValueStringPointer(),
		},
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

	if len(sources) != 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Symbol source not found: %s", data.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	source := sources[0]

	if err := data.Fill(*source); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling project symbol source: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectSymbolSourcesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectSymbolSourcesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.UpdateProjectSymbolSourceParams{
		ID:                   data.Id.ValueStringPointer(),
		Type:                 data.Type.ValueStringPointer(),
		Name:                 data.Name.ValueStringPointer(),
		AppConnectIssuer:     data.AppConnectIssuer.ValueStringPointer(),
		AppConnectPrivateKey: data.AppConnectPrivateKey.ValueStringPointer(),
		AppId:                data.AppId.ValueStringPointer(),
		Url:                  data.Url.ValueStringPointer(),
		Username:             data.Username.ValueStringPointer(),
		Password:             data.Password.ValueStringPointer(),
		Bucket:               data.Bucket.ValueStringPointer(),
		Region:               data.Region.ValueStringPointer(),
		AccessKey:            data.AccessKey.ValueStringPointer(),
		SecretKey:            data.SecretKey.ValueStringPointer(),
		Prefix:               data.Prefix.ValueStringPointer(),
		ClientEmail:          data.ClientEmail.ValueStringPointer(),
		PrivateKey:           data.PrivateKey.ValueStringPointer(),
	}
	if data.Layout != nil {
		params.Layout = &sentry.ProjectSymbolSourceLayout{
			Type:   data.Layout.Type.ValueStringPointer(),
			Casing: data.Layout.Casing.ValueStringPointer(),
		}
	}

	source, _, err := r.client.ProjectSymbolSources.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating project symbol source: %s", err.Error()))
		return
	}

	if err := data.Fill(*source); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling project symbol source: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectSymbolSourcesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectSymbolSourcesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.ProjectSymbolSources.Delete(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting project symbol source: %s", err.Error()))
		return
	}
}

func (r *ProjectSymbolSourcesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, symbolSourceId, err := splitThreePartID(req.ID, "organization", "project-slug", "symbol-source-id")
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
		ctx, path.Root("id"), symbolSourceId,
	)...)
}
