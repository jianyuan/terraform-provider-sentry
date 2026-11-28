package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	intresource "github.com/jianyuan/terraform-provider-sentry/internal/resource"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

// organizationAuthTokenAllowedScopes are the scopes Sentry lets an org auth token have
// (ORG_AUTH_TOKEN_ALLOWED_SCOPES in src/sentry/models/orgauthtoken.py).
var organizationAuthTokenAllowedScopes = []string{"org:ci"}

// organizationAuthTokenDefaultScopes is what Sentry grants when no scopes are sent.
var organizationAuthTokenDefaultScopes = []string{"org:ci"}

const organizationAuthTokenAuthNote = "Sentry only lets a signed-in session, or an OAuth application token with the `org:tokens` scope, manage org auth tokens. " +
	"Personal auth tokens, internal integration tokens and org auth tokens get `403 Forbidden`. " +
	"The OAuth token also needs `org:read` and `org:ci` to create a token, because it can only grant scopes it holds itself, and `org:write` to revoke one. " +
	"One way to get such a token is the [Sentry CLI](https://cli.sentry.dev/): run `sentry auth login --scope org:read,org:write,org:ci,org:tokens`, adding any other scopes the rest of your configuration needs, " +
	"and pass the output of `sentry auth token` to the provider, for example as `SENTRY_AUTH_TOKEN`."

func newOrganizationAuthTokenForbiddenError(action string, body []byte) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic(
		"Not allowed to manage org auth tokens",
		fmt.Sprintf(
			"Unable to %s the org auth token, got status 403: %s\n\n%s",
			action, string(body), organizationAuthTokenAuthNote,
		),
	)
}

func newOrganizationAuthTokenStatusError(action string, status int, body []byte) diag.ErrorDiagnostic {
	if status == http.StatusForbidden {
		return newOrganizationAuthTokenForbiddenError(action, body)
	}
	return diagutils.NewClientStatusError(action, status, body)
}

type OrganizationAuthTokenResourceModel struct {
	Id                  types.String                  `tfsdk:"id"`
	Organization        types.String                  `tfsdk:"organization"`
	Name                types.String                  `tfsdk:"name"`
	Scopes              supertypes.SetValueOf[string] `tfsdk:"scopes"`
	Token               types.String                  `tfsdk:"token"`
	TokenLastCharacters types.String                  `tfsdk:"token_last_characters"`
	DateCreated         types.String                  `tfsdk:"date_created"`
}

// Fill sets every attribute except token, which Sentry only returns when the token is created.
func (m *OrganizationAuthTokenResourceModel) Fill(ctx context.Context, data apiclient.OrgAuthToken) (diags diag.Diagnostics) {
	m.Id = types.StringValue(data.Id)
	m.Name = types.StringValue(data.Name)
	m.Scopes = supertypes.NewSetValueOfSlice(ctx, data.Scopes)
	if v, err := data.TokenLastCharacters.Get(); err == nil {
		m.TokenLastCharacters = types.StringValue(v)
	} else {
		m.TokenLastCharacters = types.StringNull()
	}
	m.DateCreated = types.StringValue(data.DateCreated.Format(time.RFC3339))

	return
}

var _ resource.Resource = &OrganizationAuthTokenResource{}
var _ resource.ResourceWithConfigure = &OrganizationAuthTokenResource{}
var _ resource.ResourceWithImportState = &OrganizationAuthTokenResource{}

func NewOrganizationAuthTokenResource() resource.Resource {
	return &OrganizationAuthTokenResource{}
}

type OrganizationAuthTokenResource struct {
	baseResource
}

func (r *OrganizationAuthTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_auth_token"
}

func (r *OrganizationAuthTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Organization Auth Token. Creates an org auth token, for tools such as CI pipelines that upload source maps or create releases. " +
			"Sentry has no public endpoint to update a token, so changing any argument revokes the token and creates a new one. Sentry records every new token in the audit log and emails the organization owners about it. " +
			"The token value is only returned when the token is created. It is stored in the Terraform state, so treat the state as sensitive. " +
			"It can't be recovered on import, so `token` is null after an import.\n\n" +
			"~> " + organizationAuthTokenAuthNote,

		Attributes: map[string]schema.Attribute{
			"id": ResourceIdAttribute(),
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the token belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the token, up to 255 characters. Renaming the token in Sentry also shows up as a change here, and applying it replaces the token.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 255),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scopes": schema.SetAttribute{
				MarkdownDescription: "The scopes to grant the token. Sentry currently only allows `org:ci`, which is also the default.",
				CustomType:          supertypes.NewSetTypeOf[string](ctx),
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(supertypes.NewSetValueOfSlice(ctx, organizationAuthTokenDefaultScopes).SetValue),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.OneOf(organizationAuthTokenAllowedScopes...)),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The token. Only set when Terraform creates the token; null after an import.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_last_characters": schema.StringAttribute{
				MarkdownDescription: "The last four characters of the token, as shown in the Sentry UI.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"date_created": schema.StringAttribute{
				MarkdownDescription: "When the token was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *OrganizationAuthTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationAuthTokenResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scopes, diags := data.Scopes.Get(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CreateOrganizationAuthTokenWithResponse(
		ctx,
		data.Organization.ValueString(),
		apiclient.CreateOrganizationAuthTokenJSONRequestBody{
			Name:   data.Name.ValueString(),
			Scopes: scopes,
		},
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpResp.StatusCode() != http.StatusCreated || httpResp.JSON201 == nil {
		resp.Diagnostics.Append(newOrganizationAuthTokenStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	token := *httpResp.JSON201
	resp.Diagnostics.Append(data.Fill(ctx, token)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save the token even if Sentry left out its value, so that it is revoked on destroy.
	data.Token = types.StringPointerValue(token.Token)
	if token.Token == nil || *token.Token == "" {
		resp.Diagnostics.AddError("Client error", "Sentry created the org auth token but did not return its value.")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationAuthTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationAuthTokenResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The endpoint to get a single token is private, so look the token up in the list, which is not paginated.
	httpResp, err := r.apiClient.ListOrganizationAuthTokensWithResponse(ctx, data.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(newOrganizationAuthTokenStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	var found *apiclient.OrgAuthToken
	for _, token := range *httpResp.JSON200 {
		if token.Id == data.Id.ValueString() {
			found = &token
			break
		}
	}
	if found == nil {
		// The token was revoked outside Terraform.
		resp.State.RemoveResource(ctx)
		return
	}

	// Fill leaves token alone, so the value saved on create stays in state.
	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationAuthTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every argument forces replacement.
	resp.Diagnostics.Append(diagutils.NewNotSupportedError("update"))
}

func (r *OrganizationAuthTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationAuthTokenResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteOrganizationAuthTokenWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() != http.StatusNoContent && httpResp.StatusCode() != http.StatusNotFound {
		resp.Diagnostics.Append(newOrganizationAuthTokenStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *OrganizationAuthTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	intresource.ImportState2PartPath("organization", "id")(ctx, req, resp)
}
