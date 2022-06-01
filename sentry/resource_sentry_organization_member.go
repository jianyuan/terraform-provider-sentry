package sentry

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func resourceSentryOrganizationMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryOrganizationMemberCreate,
		ReadContext:   resourceSentryOrganizationMemberRead,
		UpdateContext: resourceSentryOrganizationMemberUpdate,
		DeleteContext: resourceSentryOrganizationMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOrganizationMemberImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the team should be created for",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email of the organization member",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the role of the organization member",
			},
			"teams": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The slugs of the teams to add the member too",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"member_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique member identifier",
			},
			"pending": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The invite is pending",
			},
			"expired": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The invite is expired",
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
		teamsInput := v.([]interface{})
		teams := make([]string, len(teamsInput))
		for i, team := range teamsInput {
			teams[i] = fmt.Sprint(team)
		}
		if len(teams) > 0 {
			params.Teams = teams
		}
	}

	tflog.Debug(ctx, "Inviting Sentry organization member", map[string]interface{}{
		"email": params.Email,
		"org":   org,
		"teams": params.Teams,
	})
	member, _, err := client.OrganizationMembers.Create(org, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Invited Sentry organization member", map[string]interface{}{
		"email": params.Email,
		"org":   org,
		"teams": params.Teams,
	})

	d.Set("organization", org)
	d.SetId(member.ID)
	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	memberId := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Reading Sentry organization member", map[string]interface{}{
		"memberId": memberId,
		"org":      org,
	})
	member, resp, err := client.OrganizationMembers.Get(org, memberId)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry organization member", map[string]interface{}{
		"email": member.Email,
		"role":  member.Role,
		"id":    member.ID,
		"teams": member.Teams,
		"org":   org,
	})

	d.SetId(member.ID)
	d.Set("member_id", member.ID)
	d.Set("email", member.Email)
	d.Set("role", member.Role)
	d.Set("teams", member.Teams)
	d.Set("expired", member.Expired)
	d.Set("pending", member.Pending)
	return nil
}

func resourceSentryOrganizationMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	memberId := d.Id()
	org := d.Get("organization").(string)
	params := &sentry.UpdateOrganizationMemberParams{
		Role: d.Get("role").(string),
	}

	if v, ok := d.GetOk("teams"); ok {
		teamsInput := v.([]interface{})
		teams := make([]string, len(teamsInput))
		for i, team := range teamsInput {
			teams[i] = fmt.Sprint(team)
		}
		if len(teams) > 0 {
			params.Teams = teams
		}
	}

	tflog.Debug(ctx, "Updating organization member", map[string]interface{}{
		"email": d.Get("email"),
		"role":  params.Role,
		"id":    memberId,
		"teams": params.Teams,
		"org":   org,
	})

	member, _, err := client.OrganizationMembers.Update(org, memberId, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry organization member", map[string]interface{}{
		"email": d.Get("email"),
		"role":  params.Role,
		"id":    memberId,
		"teams": params.Teams,
		"org":   org,
	})

	d.SetId(member.ID)
	return resourceSentryOrganizationMemberRead(ctx, d, meta)
}

func resourceSentryOrganizationMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	memberId := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Deleting Sentry organization member", map[string]interface{}{
		"memberId": memberId,
		"email":    d.Get("email"),
		"org":      org,
	})
	_, err := client.OrganizationMembers.Delete(org, memberId)
	tflog.Debug(ctx, "Deleted Sentry organization member", map[string]interface{}{
		"memberId": memberId,
		"email":    d.Get("email"),
		"org":      org,
	})

	return diag.FromErr(err)
}
