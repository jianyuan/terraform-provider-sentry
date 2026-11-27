package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	intresource "github.com/jianyuan/terraform-provider-sentry/internal/resource"
	"github.com/oapi-codegen/nullable"
)

// externalTeamProviders are the providers accepted by Sentry's external team endpoints.
var externalTeamProviders = []string{
	"github",
	"github_enterprise",
	"jira_server",
	"slack",
	"slack_staging",
	"perforce",
	"gitlab",
	"cursor_origin",
	"msteams",
	"custom_scm",
}

var numericIdRegexp = regexp.MustCompile(`^[0-9]+$`)

// parseNumericId converts a numeric string attribute into the integer the Sentry API expects.
func parseNumericId(attrPath path.Path, value types.String) (int64, diag.Diagnostics) {
	var diags diag.Diagnostics
	v, err := strconv.ParseInt(value.ValueString(), 10, 64)
	if err != nil {
		diags.AddAttributeError(attrPath, "Invalid ID", fmt.Sprintf("Expected a numeric ID, got %q.", value.ValueString()))
	}
	return v, diags
}

// externalIdCreateValue omits external_id from a create request when it is not set.
func externalIdCreateValue(value types.String) nullable.Nullable[string] {
	if value.IsNull() || value.IsUnknown() {
		return nullable.Nullable[string]{}
	}
	return nullable.NewNullableWithValue(value.ValueString())
}

// externalIdUpdateValue sends an explicit null to clear external_id when it is removed from the configuration.
func externalIdUpdateValue(value types.String) nullable.Nullable[string] {
	if value.IsNull() || value.IsUnknown() {
		return nullable.NewNullNullable[string]()
	}
	return nullable.NewNullableWithValue(value.ValueString())
}

// withQueryParam adds a query parameter that the generated client does not expose.
func withQueryParam(key, value string) apiclient.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

type ExternalTeamResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Organization  types.String `tfsdk:"organization"`
	Team          types.String `tfsdk:"team"`
	ProviderName  types.String `tfsdk:"provider_name"`
	ExternalName  types.String `tfsdk:"external_name"`
	IntegrationId types.String `tfsdk:"integration_id"`
	ExternalId    types.String `tfsdk:"external_id"`
}

func (m *ExternalTeamResourceModel) Fill(ctx context.Context, data apiclient.ExternalActor) (diags diag.Diagnostics) {
	m.Id = types.StringValue(data.Id)
	m.ProviderName = types.StringValue(data.Provider)
	m.ExternalName = types.StringValue(data.ExternalName)
	m.IntegrationId = types.StringValue(data.IntegrationId)
	m.ExternalId = types.StringPointerValue(data.ExternalId)

	return
}

func (m ExternalTeamResourceModel) toRequest(update bool) (apiclient.ExternalActorRequest, diag.Diagnostics) {
	integrationId, diags := parseNumericId(path.Root("integration_id"), m.IntegrationId)

	body := apiclient.ExternalActorRequest{
		Provider:      m.ProviderName.ValueString(),
		ExternalName:  m.ExternalName.ValueString(),
		IntegrationId: integrationId,
	}
	if update {
		body.ExternalId = externalIdUpdateValue(m.ExternalId)
	} else {
		body.ExternalId = externalIdCreateValue(m.ExternalId)
	}

	return body, diags
}

var _ resource.Resource = &ExternalTeamResource{}
var _ resource.ResourceWithConfigure = &ExternalTeamResource{}
var _ resource.ResourceWithImportState = &ExternalTeamResource{}

func NewExternalTeamResource() resource.Resource {
	return &ExternalTeamResource{}
}

type ExternalTeamResource struct {
	baseResource
}

func (r *ExternalTeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_team"
}

func (r *ExternalTeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry External Team. Links a Sentry team to a team or channel in an external provider, such as a Slack channel for team notifications, or a GitHub or GitLab team so that CODEOWNERS team handles resolve to the Sentry team. Sentry only allows this on organizations with the `integrations-codeowners` feature enabled; otherwise the API returns `403 Forbidden`.",

		Attributes: map[string]schema.Attribute{
			"id": ResourceIdAttribute(),
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the team belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team": schema.StringAttribute{
				MarkdownDescription: "The slug of the Sentry team to link.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "The external provider. Named `provider_name` because `provider` is a reserved Terraform meta-argument. One of: `" + strings.Join(externalTeamProviders, "`, `") + "`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(externalTeamProviders...),
				},
			},
			"external_name": schema.StringAttribute{
				MarkdownDescription: "The name of the team or channel in the external provider, such as `#backend-alerts` for a Slack channel or `@my-org/backend` for a GitHub team. GitHub, GitHub Enterprise, GitLab and Cursor names must start with `@`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization integration the mapping belongs to, for example `data.sentry_organization_integration.slack.id`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(numericIdRegexp, "must be a numeric ID"),
				},
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the team or channel in the external provider, such as a Slack channel ID (`C0123456789`).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
		},
	}
}

func (r *ExternalTeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ExternalTeamResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := data.toRequest(false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CreateTeamExternalTeamWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Team.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpResp.StatusCode() == http.StatusOK && httpResp.JSON200 != nil {
		// Sentry returns 200 with the existing mapping instead of creating a new one. Its lookup ignores the
		// team, so the existing mapping may belong to a different team. Adopting it silently would let two
		// resources manage the same mapping, so ask the user to import it instead.
		existing := *httpResp.JSON200
		resp.Diagnostics.AddError(
			"External team already exists",
			fmt.Sprintf(
				"Sentry already has an external team mapping for provider %q, integration %q and external name %q "+
					"(mapping ID %q, Sentry team ID %q). An external name can only be linked to one Sentry team. "+
					"Import the existing mapping with ID \"<organization>/<team>/%s\", or delete it before creating a new one.",
				existing.Provider, existing.IntegrationId, existing.ExternalName, existing.Id, stringPointerOrEmpty(existing.TeamId), existing.Id,
			),
		)
		return
	} else if httpResp.StatusCode() != http.StatusCreated || httpResp.JSON201 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON201)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExternalTeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ExternalTeamResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// There is no endpoint to retrieve a single external team, so read the team's mappings through
	// the team details endpoint, which includes them with `?expand=externalTeams`.
	httpResp, err := r.apiClient.GetOrganizationTeamWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Team.ValueString(),
		withQueryParam("expand", "externalTeams"),
	)
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

	var found *apiclient.ExternalActor
	if httpResp.JSON200.ExternalTeams != nil {
		for _, externalTeam := range *httpResp.JSON200.ExternalTeams {
			if externalTeam.Id == data.Id.ValueString() {
				found = &externalTeam
				break
			}
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExternalTeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ExternalTeamResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := data.toRequest(true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UpdateTeamExternalTeamWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Team.ValueString(),
		data.Id.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExternalTeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ExternalTeamResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteTeamExternalTeamWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Team.ValueString(),
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

func (r *ExternalTeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	intresource.ImportState3PartPath("organization", "team", "id")(ctx, req, resp)
}

func stringPointerOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
