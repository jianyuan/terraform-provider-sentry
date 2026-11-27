package provider

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	intresource "github.com/jianyuan/terraform-provider-sentry/internal/resource"
)

type ProjectSymbolSourcesResourceModel struct {
	Id                            types.String                             `tfsdk:"id"`
	Organization                  types.String                             `tfsdk:"organization"`
	Project                       types.String                             `tfsdk:"project"`
	Type                          types.String                             `tfsdk:"type"`
	Name                          types.String                             `tfsdk:"name"`
	Layout                        *ProjectSymbolSourcesResourceLayoutModel `tfsdk:"layout"`
	AppConnectIssuer              types.String                             `tfsdk:"app_connect_issuer"`
	AppConnectPrivateKey          types.String                             `tfsdk:"app_connect_private_key"`
	AppConnectPrivateKeyWo        types.String                             `tfsdk:"app_connect_private_key_wo"`
	AppConnectPrivateKeyWoVersion types.Int64                              `tfsdk:"app_connect_private_key_wo_version"`
	AppId                         types.String                             `tfsdk:"app_id"`
	Url                           types.String                             `tfsdk:"url"`
	Username                      types.String                             `tfsdk:"username"`
	Password                      types.String                             `tfsdk:"password"`
	PasswordWo                    types.String                             `tfsdk:"password_wo"`
	PasswordWoVersion             types.Int64                              `tfsdk:"password_wo_version"`
	Bucket                        types.String                             `tfsdk:"bucket"`
	Region                        types.String                             `tfsdk:"region"`
	AccessKey                     types.String                             `tfsdk:"access_key"`
	SecretKey                     types.String                             `tfsdk:"secret_key"`
	SecretKeyWo                   types.String                             `tfsdk:"secret_key_wo"`
	SecretKeyWoVersion            types.Int64                              `tfsdk:"secret_key_wo_version"`
	Prefix                        types.String                             `tfsdk:"prefix"`
	ClientEmail                   types.String                             `tfsdk:"client_email"`
	PrivateKey                    types.String                             `tfsdk:"private_key"`
	PrivateKeyWo                  types.String                             `tfsdk:"private_key_wo"`
	PrivateKeyWoVersion           types.Int64                              `tfsdk:"private_key_wo_version"`
}

// projectSymbolSourceSecret is a secret of a symbol source that can be set
// through either a regular attribute or a write-only attribute.
type projectSymbolSourceSecret struct {
	// attr is the name of the regular attribute.
	attr string
	// fields returns the value of the regular attribute and the version of the
	// write-only value.
	fields func(m *ProjectSymbolSourcesResourceModel) (plain types.String, version types.Int64)
	// set sets the secret in an update request body.
	set func(b *projectSymbolSourceUpdateBody, v any)
}

var projectSymbolSourceSecrets = []projectSymbolSourceSecret{
	{
		attr: "app_connect_private_key",
		fields: func(m *ProjectSymbolSourcesResourceModel) (types.String, types.Int64) {
			return m.AppConnectPrivateKey, m.AppConnectPrivateKeyWoVersion
		},
		set: func(b *projectSymbolSourceUpdateBody, v any) { b.AppConnectPrivateKey = v },
	},
	{
		attr: "password",
		fields: func(m *ProjectSymbolSourcesResourceModel) (types.String, types.Int64) {
			return m.Password, m.PasswordWoVersion
		},
		set: func(b *projectSymbolSourceUpdateBody, v any) { b.Password = v },
	},
	{
		attr: "secret_key",
		fields: func(m *ProjectSymbolSourcesResourceModel) (types.String, types.Int64) {
			return m.SecretKey, m.SecretKeyWoVersion
		},
		set: func(b *projectSymbolSourceUpdateBody, v any) { b.SecretKey = v },
	},
	{
		attr: "private_key",
		fields: func(m *ProjectSymbolSourcesResourceModel) (types.String, types.Int64) {
			return m.PrivateKey, m.PrivateKeyWoVersion
		},
		set: func(b *projectSymbolSourceUpdateBody, v any) { b.PrivateKey = v },
	},
}

// projectSymbolSourceHiddenSecret is how Sentry returns a secret, and tells
// Sentry to keep the stored value when sent in an update. An update replaces
// the whole source, so a secret left out of it is removed or, if the source
// type requires it, rejected.
var projectSymbolSourceHiddenSecret = map[string]bool{"hidden-secret": true}

// projectSymbolSourceUpdateBody is the update request body. Its secret fields
// shadow the ones in the embedded params, so they can hold either a string or
// projectSymbolSourceHiddenSecret.
type projectSymbolSourceUpdateBody struct {
	*sentry.UpdateProjectSymbolSourceParams
	AppConnectPrivateKey any `json:"appconnectPrivateKey,omitempty"`
	Password             any `json:"password,omitempty"`
	SecretKey            any `json:"secret_key,omitempty"`
	PrivateKey           any `json:"private_key,omitempty"`
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

var _ resource.Resource = &ProjectSymbolSourcesResource{}
var _ resource.ResourceWithConfigure = &ProjectSymbolSourcesResource{}
var _ resource.ResourceWithImportState = &ProjectSymbolSourcesResource{}
var _ resource.ResourceWithConfigValidators = &ProjectSymbolSourcesResource{}

func NewProjectSymbolSourcesResource() resource.Resource {
	return &ProjectSymbolSourcesResource{}
}

type ProjectSymbolSourcesResource struct {
	baseResource
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
		MarkdownDescription: "Sentry Project Symbol Source. See the [Sentry documentation](https://docs.sentry.io/api/projects/add-a-symbol-source-to-a-project/) for more information.\n\n" +
			"Each secret (`app_connect_private_key`, `password`, `secret_key` and `private_key`) can instead be set through a write-only `<name>_wo` attribute, which is not stored in the Terraform state. " +
			"Write-only attributes require Terraform 1.11 or later. Change `<name>_wo_version` to send a new value.\n\n" +
			"Sentry does not return these secrets, so neither the secrets nor their write-only attributes are populated on import.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
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
				Description: "The App Store Connect API Private Key. Required for AppStoreConnect sources, invalid for all others. Conflicts with `app_connect_private_key_wo`.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_connect_private_key_wo":         writeOnlyStringAttribute("app_connect_private_key", "The App Store Connect API Private Key. Required for AppStoreConnect sources unless `app_connect_private_key` is set, invalid for all others. Conflicts with `app_connect_private_key`."),
			"app_connect_private_key_wo_version": writeOnlyVersionAttribute("app_connect_private_key"),
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
				Description: "The password for accessing the source. Optional for HTTP sources, invalid for all others. Conflicts with `password_wo`.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_wo":         writeOnlyStringAttribute("password", "The password for accessing the source. Optional for HTTP sources, invalid for all others. Conflicts with `password`."),
			"password_wo_version": writeOnlyVersionAttribute("password"),
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
				Description: "The AWS Secret Access Key.Required for S3 sources, invalid for all others. Conflicts with `secret_key_wo`.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_key_wo":         writeOnlyStringAttribute("secret_key", "The AWS Secret Access Key. Required for S3 sources unless `secret_key` is set, invalid for all others. Conflicts with `secret_key`."),
			"secret_key_wo_version": writeOnlyVersionAttribute("secret_key"),
			"prefix": schema.StringAttribute{
				Description: "The GCS or S3 prefix. Optional for GCS and S3 sourcse, invalid for HTTP and AppStoreConnect sources.",
				Optional:    true,
			},
			"client_email": schema.StringAttribute{
				Description: "The GCS email address for authentication. Required for GCS sources, invalid for all others.",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The GCS private key. Required for GCS sources, invalid for all others. Conflicts with `private_key_wo`.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key_wo":         writeOnlyStringAttribute("private_key", "The GCS private key. Required for GCS sources unless `private_key` is set, invalid for all others. Conflicts with `private_key`."),
			"private_key_wo_version": writeOnlyVersionAttribute("private_key"),
		},
	}
}

func (r *ProjectSymbolSourcesResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := make([]resource.ConfigValidator, 0, len(projectSymbolSourceSecrets))
	for _, secret := range projectSymbolSourceSecrets {
		validators = append(validators, resourcevalidator.Conflicting(
			path.MatchRoot(secret.attr),
			path.MatchRoot(secret.attr+"_wo"),
		))
	}
	return validators
}

// resolveSecrets works out what to send for each secret. data is the plan,
// prior is the state before an update, or nil on create.
func (r *ProjectSymbolSourcesResource) resolveSecrets(
	ctx context.Context,
	config tfsdk.Config,
	data *ProjectSymbolSourcesResourceModel,
	prior *ProjectSymbolSourcesResourceModel,
	private privateState,
) (secrets map[string]writeOnlySecret, diags diag.Diagnostics) {
	secrets = make(map[string]writeOnlySecret, len(projectSymbolSourceSecrets))
	for _, secret := range projectSymbolSourceSecrets {
		plain, version := secret.fields(data)

		var priorVersion *types.Int64
		if prior != nil {
			_, v := secret.fields(prior)
			priorVersion = &v
		}

		resolved, d := resolveWriteOnlySecretFromConfig(ctx, config, private, secret.attr, plain, version, priorVersion)
		diags.Append(d...)
		secrets[secret.attr] = resolved
	}
	return
}

// saveSecretsInUse records in private state which write-only secrets are in use.
func saveSecretsInUse(ctx context.Context, private privateState, secrets map[string]writeOnlySecret) (diags diag.Diagnostics) {
	for attr, secret := range secrets {
		diags.Append(setWriteOnlyInUse(ctx, private, attr, secret.InUse)...)
	}
	return
}

func (r *ProjectSymbolSourcesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectSymbolSourcesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secrets, diags := r.resolveSecrets(ctx, req.Config, &data, nil, resp.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.CreateProjectSymbolSourceParams{
		Type:                 data.Type.ValueStringPointer(),
		Name:                 data.Name.ValueStringPointer(),
		AppConnectIssuer:     data.AppConnectIssuer.ValueStringPointer(),
		AppConnectPrivateKey: secrets["app_connect_private_key"].Value,
		AppId:                data.AppId.ValueStringPointer(),
		Url:                  data.Url.ValueStringPointer(),
		Username:             data.Username.ValueStringPointer(),
		Password:             secrets["password"].Value,
		Bucket:               data.Bucket.ValueStringPointer(),
		Region:               data.Region.ValueStringPointer(),
		AccessKey:            data.AccessKey.ValueStringPointer(),
		SecretKey:            secrets["secret_key"].Value,
		Prefix:               data.Prefix.ValueStringPointer(),
		ClientEmail:          data.ClientEmail.ValueStringPointer(),
		PrivateKey:           secrets["private_key"].Value,
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
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	}

	if err := data.Fill(*source); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.Append(saveSecretsInUse(ctx, resp.Private, secrets)...)
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
	if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if len(sources) != 1 {
		resp.State.RemoveResource(ctx)
		return
	}

	source := sources[0]

	if err := data.Fill(*source); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
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

	var state ProjectSymbolSourcesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secrets, diags := r.resolveSecrets(ctx, req.Config, &data, &state, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := &projectSymbolSourceUpdateBody{
		UpdateProjectSymbolSourceParams: &sentry.UpdateProjectSymbolSourceParams{
			ID:               data.Id.ValueStringPointer(),
			Type:             data.Type.ValueStringPointer(),
			Name:             data.Name.ValueStringPointer(),
			AppConnectIssuer: data.AppConnectIssuer.ValueStringPointer(),
			AppId:            data.AppId.ValueStringPointer(),
			Url:              data.Url.ValueStringPointer(),
			Username:         data.Username.ValueStringPointer(),
			Bucket:           data.Bucket.ValueStringPointer(),
			Region:           data.Region.ValueStringPointer(),
			AccessKey:        data.AccessKey.ValueStringPointer(),
			Prefix:           data.Prefix.ValueStringPointer(),
			ClientEmail:      data.ClientEmail.ValueStringPointer(),
		},
	}
	if data.Layout != nil {
		body.Layout = &sentry.ProjectSymbolSourceLayout{
			Type:   data.Layout.Type.ValueStringPointer(),
			Casing: data.Layout.Casing.ValueStringPointer(),
		}
	}
	for _, secret := range projectSymbolSourceSecrets {
		switch s := secrets[secret.attr]; {
		case s.Value != nil:
			secret.set(body, *s.Value)
		case s.Keep:
			secret.set(body, projectSymbolSourceHiddenSecret)
		}
	}

	// The go-sentry update params only take strings for secrets, so build the
	// request here to be able to send projectSymbolSourceHiddenSecret.
	httpReq, err := r.client.NewRequest(
		http.MethodPut,
		"0/projects/"+data.Organization.ValueString()+"/"+data.Project.ValueString()+"/symbol-sources/?id="+url.QueryEscape(data.Id.ValueString()),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	}

	var source sentry.ProjectSymbolSource
	if _, err := r.client.Do(ctx, httpReq, &source); err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	}

	if err := data.Fill(source); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.Append(saveSecretsInUse(ctx, resp.Private, secrets)...)
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
	if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	}
}

func (r *ProjectSymbolSourcesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	intresource.ImportState3PartPath("organization", "project", "id")(ctx, req, resp)
}
