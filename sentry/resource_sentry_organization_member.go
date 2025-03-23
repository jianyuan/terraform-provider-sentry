package sentry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func resourceSentryOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for managing Sentry organization members. To add a member to a team, use the `sentry_team_member` resource.",

		CreateContext: resourceSentryOrganizationMemberCreate,
		ReadContext:   resourceSentryOrganizationMemberRead,
		UpdateContext: resourceSentryOrganizationMemberUpdate,
		DeleteContext: resourceSentryOrganizationMemberDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the user should be invited to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "The email of the organization member.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"role": {
				Description: "This is the role of the organization member.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"billing",
						"member",
						"admin",
						"manager",
						"owner",
					},
					false,
				),
			},
			"internal_id": {
				Description: "The internal ID for this organization membership.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pending": {
				Description: "The invite is pending.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"expired": {
				Description: "The invite has expired.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*providerdata.ProviderData).ApiClient

	org := d.Get("organization").(string)
	params := apiclient.CreateOrganizationMemberJSONRequestBody{
		Email:   d.Get("email").(string),
		OrgRole: d.Get("role").(string),
	}

	tflog.Debug(ctx, "Inviting organization member", map[string]interface{}{
		"email": params.Email,
		"org":   org,
	})
	httpResp, err := apiClient.CreateOrganizationMemberWithResponse(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	} else if httpResp.StatusCode() != http.StatusCreated || httpResp.JSON201 == nil {
		return diag.FromErr(fmt.Errorf("unexpected status code: %d", httpResp.StatusCode()))
	}

	member := httpResp.JSON201

	d.SetId(tfutils.BuildTwoPartId(org, member.Id))
	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*providerdata.ProviderData).ApiClient

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Reading organization member", map[string]interface{}{
		"org":      org,
		"memberID": memberID,
	})
	httpResp, err := apiClient.GetOrganizationMemberWithResponse(ctx, org, memberID)
	if err != nil {
		return diag.FromErr(err)
	} else if httpResp.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		return diag.FromErr(fmt.Errorf("unexpected status code: %d", httpResp.StatusCode()))
	}

	member := httpResp.JSON200

	d.SetId(tfutils.BuildTwoPartId(org, member.Id))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("internal_id", member.Id),
		d.Set("email", member.Email),
		d.Set("role", member.OrgRole),
		d.Set("expired", member.Expired),
		d.Set("pending", member.Pending),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryOrganizationMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*providerdata.ProviderData).ApiClient

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getHttpResp, err := apiClient.GetOrganizationMemberWithResponse(ctx, org, memberID)
	if err != nil {
		return diag.FromErr(err)
	} else if getHttpResp.StatusCode() != http.StatusOK || getHttpResp.JSON200 == nil {
		return diag.FromErr(fmt.Errorf("unexpected status code: %d", getHttpResp.StatusCode()))
	}
	orgMember := getHttpResp.JSON200

	teamRoles := make([]apiclient.TeamRole, len(orgMember.TeamRoles))
	for i, teamRole := range orgMember.TeamRoles {
		teamRoles[i] = apiclient.TeamRole{
			TeamSlug: teamRole.TeamSlug,
			Role:     teamRole.Role,
		}
	}
	params := apiclient.UpdateOrganizationMemberJSONRequestBody{
		OrgRole:   ptr.Ptr(d.Get("role").(string)),
		TeamRoles: &teamRoles,
	}

	tflog.Debug(ctx, "Updating organization member", map[string]interface{}{
		"email": d.Get("email"),
		"role":  params.OrgRole,
		"id":    memberID,
		"org":   org,
	})

	httpResp, err := apiClient.UpdateOrganizationMemberWithResponse(ctx, org, memberID, params)
	if err != nil {
		return diag.FromErr(err)
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		return diag.FromErr(fmt.Errorf("unexpected status code: %d", httpResp.StatusCode()))
	}

	member := httpResp.JSON200

	d.SetId(tfutils.BuildTwoPartId(org, member.Id))
	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*providerdata.ProviderData).ApiClient

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting organization member", map[string]interface{}{
		"org":      org,
		"memberID": memberID,
	})
	httpResp, err := apiClient.DeleteOrganizationMemberWithResponse(ctx, org, memberID)
	if err != nil {
		return diag.FromErr(err)
	} else if httpResp.StatusCode() != http.StatusNoContent {
		return diag.FromErr(fmt.Errorf("unexpected status code: %d", httpResp.StatusCode()))
	}

	return nil
}

func splitSentryOrganizationMemberID(id string) (org string, memberID string, err error) {
	org, memberID, err = tfutils.SplitTwoPartId(id, "organization-id", "member-id")
	return
}
