package sentry

import (
	"context"
	"sort"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryOrganizationMember() *schema.Resource {
	return &schema.Resource{
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
			},
			"teams": {
				Description: "The teams the organization member should be added to.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	params := &sentry.CreateOrganizationMemberParams{
		Email: d.Get("email").(string),
		Role:  d.Get("role").(string),
	}

	if v, ok := d.GetOk("teams"); ok {
		teams := expandStringList(v.([]interface{}))
		if len(teams) > 0 {
			params.Teams = teams
		}
	}

	tflog.Debug(ctx, "Inviting organization member", map[string]interface{}{
		"email": params.Email,
		"org":   org,
		"teams": params.Teams,
	})
	member, _, err := client.OrganizationMembers.Create(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(org, member.ID))
	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())

	tflog.Debug(ctx, "Reading organization member", map[string]interface{}{
		"org":      org,
		"memberID": memberID,
	})
	member, resp, err := client.OrganizationMembers.Get(ctx, org, memberID)
	if found, err := checkClientGet(resp, err, d); !found {
		tflog.Info(ctx, "Removed organization membership from state because it no longer exists in Sentry", map[string]interface{}{
			"org":      org,
			"memberID": memberID,
		})
		return diag.FromErr(err)
	}

	sort.Strings(member.Teams)

	d.SetId(buildTwoPartID(org, member.ID))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("internal_id", member.ID),
		d.Set("email", member.Email),
		d.Set("role", member.Role),
		d.Set("teams", member.Teams),
		d.Set("expired", member.Expired),
		d.Set("pending", member.Pending),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryOrganizationMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	params := &sentry.UpdateOrganizationMemberParams{
		Role: d.Get("role").(string),
	}

	if v, ok := d.GetOk("teams"); ok {
		teams := expandStringList(v.([]interface{}))
		if len(teams) > 0 {
			params.Teams = teams
		}
	}

	tflog.Debug(ctx, "Updating organization member", map[string]interface{}{
		"email": d.Get("email"),
		"role":  params.Role,
		"id":    memberID,
		"teams": params.Teams,
		"org":   org,
	})

	member, _, err := client.OrganizationMembers.Update(ctx, org, memberID, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(org, member.ID))
	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, memberID, err := splitSentryOrganizationMemberID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting organization member", map[string]interface{}{
		"org":      org,
		"memberID": memberID,
	})
	_, err = client.OrganizationMembers.Delete(ctx, org, memberID)
	return diag.FromErr(err)
}

func splitSentryOrganizationMemberID(id string) (org string, memberID string, err error) {
	org, memberID, err = splitTwoPartID(id, "organization-id", "member-id")
	return
}
