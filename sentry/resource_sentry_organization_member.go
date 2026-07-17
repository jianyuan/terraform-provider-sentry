package sentry

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gosentry "github.com/mzglinski/go-sentry/v2/sentry"
	"github.com/mzglinski/terraform-provider-sentry/internal/providerdata"
	"github.com/mzglinski/terraform-provider-sentry/internal/tfutils"
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
			"user_id": {
				Description: "The Sentry User ID of the organization member.",
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
	client := meta.(*providerdata.ProviderData).Client

	org := d.Get("organization").(string)
	params := &gosentry.CreateOrganizationMemberParams{
		Email: d.Get("email").(string),
		Role:  d.Get("role").(string),
	}

	tflog.Debug(ctx, "Inviting organization member", map[string]interface{}{
		"email": params.Email,
		"org":   org,
	})
	member, resp, err := client.OrganizationMembers.Create(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error inviting organization member: %s, status: %d, body: %s", resp.Status, resp.StatusCode, resp.Body)
	}

	d.SetId(tfutils.BuildTwoPartId(org, member.ID))
	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Error parsing resource ID for Sentry Organization Member '%s': %s. Assuming deleted or new.", d.Id(), err))
		d.SetId("")
		return nil
	}

	tflog.Debug(ctx, "Reading organization member", map[string]interface{}{
		"org":          org,
		"membershipID": memberID,
	})
	member, resp, err := client.OrganizationMembers.Get(ctx, org, memberID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			tflog.Info(ctx, "Sentry Organization Member not found, removing from state", map[string]interface{}{"org": org, "membershipID": memberID})
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading Sentry Organization Member: %s, status: %d, body: %s", resp.Status, resp.StatusCode, resp.Body)
	}

	d.SetId(tfutils.BuildTwoPartId(org, member.ID))
	err = errors.Join(
		d.Set("organization", org),
		d.Set("internal_id", member.ID),
		d.Set("user_id", member.User.ID),
		d.Set("email", member.Email),
		d.Set("role", member.OrgRole),
		d.Set("expired", member.Expired),
		d.Set("pending", member.Pending),
	)
	return diag.FromErr(err)
}

func resourceSentryOrganizationMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.HasChange("role") {
		tflog.Debug(ctx, "No attribute changes for organization member.")
		return resourceSentryOrganizationMemberRead(ctx, d, meta)
	}

	params := &gosentry.UpdateOrganizationMemberParams{
		OrganizationRole: d.Get("role").(string),
	}

	tflog.Debug(ctx, "Updating organization member", map[string]interface{}{
		"org":          org,
		"membershipID": memberID,
		"newRole":      params.OrganizationRole,
	})

	updatedMember, resp, err := client.OrganizationMembers.Update(ctx, org, memberID, params)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating Sentry Organization Member: %s, status: %d, body: %s", resp.Status, resp.StatusCode, resp.Body)
	}

	if err := d.Set("internal_id", updatedMember.ID); err != nil {
		return diag.FromErr(err)
	}
	if updatedMember.User.ID != "" {
		if err := d.Set("user_id", updatedMember.User.ID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting organization member", map[string]interface{}{
		"org":          org,
		"membershipID": memberID,
	})
	resp, err := client.OrganizationMembers.Delete(ctx, org, memberID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			tflog.Info(ctx, "Sentry Organization Member not found during delete, removing from state", map[string]interface{}{"org": org, "membershipID": memberID})
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting Sentry Organization Member: %s, status: %d, body: %s", resp.Status, resp.StatusCode, resp.Body)
	}

	return nil
}

func splitSentryOrganizationMemberID(id string) (org string, memberID string, err error) {
	org, memberID, err = tfutils.SplitTwoPartId(id, "organization-slug", "membership-id")
	return
}
