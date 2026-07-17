package sentry

import (
	"context"
	"errors"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mzglinski/go-sentry/v2/sentry"
	"github.com/mzglinski/terraform-provider-sentry/internal/providerdata"
)

func resourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Organization resource.",

		CreateContext: resourceSentryOrganizationCreate,
		ReadContext:   resourceSentryOrganizationRead,
		UpdateContext: resourceSentryOrganizationUpdate,
		DeleteContext: resourceSentryOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The human readable name for the organization.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The unique URL slug for this organization.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"agree_terms": {
				Description: "You agree to the applicable terms of service and privacy policy. This is only used for creation.",
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
			},
			"internal_id": {
				Description: "The internal ID for this organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_early_adopter": {
				Description: "Opt-in to new features before they're released to the public.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"hide_ai_features": {
				Description: "Hide AI features from the organization.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"codecov_access": {
				Description: "Enable Code Coverage Insights. This feature is only available for organizations on the Team plan and above.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"default_role": {
				Description: "The default role new members will receive. Valid values are `member`, `admin`, `manager`, `owner`.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"open_membership": {
				Description: "Allow organization members to freely join any team.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"events_member_admin": {
				Description: "Allow members to delete events by granting them the `event:admin` scope.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"alerts_member_write": {
				Description: "Allow members to create, edit, and delete alert rules by granting them the `alerts:write` scope.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"attachments_role": {
				Description: "The role required to download event attachments. Valid values are `member`, `admin`, `manager`, `owner`.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"debug_files_role": {
				Description: "The role required to download debug information files. Valid values are `member`, `admin`, `manager`, `owner`.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"avatar_type": {
				Description: "The type of display picture for the organization. Valid values are `letter_avatar`, `upload`.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"avatar": {
				Description: "The image to upload as the organization avatar, in base64. Required if `avatar_type` is `upload`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"require_2fa": {
				Description: "Require and enforce two-factor authentication for all members.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"allow_shared_issues": {
				Description: "Allow sharing of limited details on issues to anonymous users.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enhanced_privacy": {
				Description: "Enable enhanced privacy controls to limit personally identifiable information (PII) as well as source code in things like notifications.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"scrape_javascript": {
				Description: "Allow Sentry to scrape missing JavaScript source context when possible.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"store_crash_reports": {
				Description: "How many native crash reports to store per issue. Valid values are `0`, `1`, `5`, `10`, `20`, `50`, `100`, `-1` (unlimited).",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"allow_join_requests": {
				Description: "Allow users to request to join your organization.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"data_scrubber": {
				Description: "Require server-side data scrubbing for all projects.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"data_scrubber_defaults": {
				Description: "Apply the default scrubbers to prevent things like passwords and credit cards from being stored for all projects.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"sensitive_fields": {
				Description: "A list of additional global field names to match against when scrubbing data for all projects.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"safe_fields": {
				Description: "A list of global field names which data scrubbers should ignore.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scrub_ip_addresses": {
				Description: "Prevent IP addresses from being stored for new events on all projects.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"relay_pii_config": {
				Description: "Advanced data scrubbing rules that can be configured for each project as a JSON string.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"trusted_relays": {
				Description: "A list of local Relays registered for the organization.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the relay.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The public key of the relay.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A description for the relay.",
						},
					},
				},
			},
			"github_pr_bot": {
				Description: "Allow Sentry to comment on recent pull requests suspected of causing issues.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"github_open_pr_bot": {
				Description: "Allow Sentry to comment on open pull requests to show recent error issues for the code being changed.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"github_nudge_invite": {
				Description: "Allow Sentry to detect users committing to your GitHub repositories that are not part of your Sentry organization.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"gitlab_pr_bot": {
				Description: "Allow Sentry to comment on recent pull requests suspected of causing issues.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"allow_member_project_creation": {
				Description: "Allow members to create projects.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client

	params := &sentry.CreateOrganizationParams{
		Name:       sentry.String(d.Get("name").(string)),
		AgreeTerms: sentry.Bool(d.Get("agree_terms").(bool)),
	}
	if slug, ok := d.GetOk("slug"); ok {
		params.Slug = sentry.String(slug.(string))
	}

	tflog.Debug(ctx, "Creating organization", map[string]interface{}{"org": params.Name})
	organization, _, err := client.Organizations.Create(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sentry.StringValue(organization.Slug))
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client
	org := d.Id()

	tflog.Debug(ctx, "Reading organization", map[string]interface{}{"org": org})
	organization, _, err := client.Organizations.Get(ctx, org)
	if err != nil {
		if sErr, ok := err.(*sentry.ErrorResponse); ok {
			if sErr.Response.StatusCode == http.StatusNotFound {
				tflog.Info(ctx, "Removing organization from state because it no longer exists in Sentry", map[string]interface{}{"org": org})
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	err = errors.Join(
		d.Set("name", organization.Name),
		d.Set("slug", organization.Slug),
		d.Set("agree_terms", true),
		d.Set("internal_id", organization.ID),
		d.Set("is_early_adopter", organization.IsEarlyAdopter),
		d.Set("hide_ai_features", organization.HideAIFeatures),
		d.Set("codecov_access", organization.CodecovAccess),
		d.Set("default_role", organization.DefaultRole),
		d.Set("open_membership", organization.OpenMembership),
		d.Set("events_member_admin", organization.EventsMemberAdmin),
		d.Set("alerts_member_write", organization.AlertsMemberWrite),
		d.Set("attachments_role", organization.AttachmentsRole),
		d.Set("debug_files_role", organization.DebugFilesRole),
		d.Set("avatar_type", organization.Avatar.AvatarType),
		d.Set("require_2fa", organization.Require2FA),
		d.Set("allow_shared_issues", organization.AllowSharedIssues),
		d.Set("enhanced_privacy", organization.EnhancedPrivacy),
		d.Set("scrape_javascript", organization.ScrapeJavaScript),
		d.Set("store_crash_reports", organization.StoreCrashReports),
		d.Set("allow_join_requests", organization.AllowJoinRequests),
		d.Set("data_scrubber", organization.DataScrubber),
		d.Set("data_scrubber_defaults", organization.DataScrubberDefaults),
		d.Set("sensitive_fields", organization.SensitiveFields),
		d.Set("safe_fields", organization.SafeFields),
		d.Set("scrub_ip_addresses", organization.ScrubIPAddresses),
		d.Set("relay_pii_config", organization.RelayPiiConfig),
		d.Set("github_pr_bot", organization.GithubPRBot),
		d.Set("github_open_pr_bot", organization.GithubOpenPRBot),
		d.Set("github_nudge_invite", organization.GithubNudgeInvite),
		d.Set("gitlab_pr_bot", organization.GitlabPRBot),
		d.Set("allow_member_project_creation", organization.AllowMemberProjectCreation),
	)

	if organization.TrustedRelays != nil {
		var trustedRelays []map[string]interface{}
		for _, apiRelay := range organization.TrustedRelays {
			relay := make(map[string]interface{})
			relay["name"] = apiRelay.Name
			relay["public_key"] = apiRelay.PublicKey
			relay["description"] = apiRelay.Description
			trustedRelays = append(trustedRelays, relay)
		}
		err = errors.Join(err, d.Set("trusted_relays", trustedRelays))
	}

	return diag.FromErr(err)
}

func resourceSentryOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client
	org := d.Id()
	params := &sentry.UpdateOrganizationParams{
		Name: sentry.String(d.Get("name").(string)),
	}

	if d.HasChange("slug") {
		if slug, ok := d.GetOk("slug"); ok {
			params.Slug = sentry.String(slug.(string))
		}
	}

	if d.HasChange("is_early_adopter") {
		params.IsEarlyAdopter = sentry.Bool(d.Get("is_early_adopter").(bool))
	}
	if d.HasChange("hide_ai_features") {
		params.HideAIFeatures = sentry.Bool(d.Get("hide_ai_features").(bool))
	}
	if d.HasChange("codecov_access") {
		params.CodecovAccess = sentry.Bool(d.Get("codecov_access").(bool))
	}
	if d.HasChange("default_role") {
		params.DefaultRole = sentry.String(d.Get("default_role").(string))
	}
	if d.HasChange("open_membership") {
		params.OpenMembership = sentry.Bool(d.Get("open_membership").(bool))
	}
	if d.HasChange("events_member_admin") {
		params.EventsMemberAdmin = sentry.Bool(d.Get("events_member_admin").(bool))
	}
	if d.HasChange("alerts_member_write") {
		params.AlertsMemberWrite = sentry.Bool(d.Get("alerts_member_write").(bool))
	}
	if d.HasChange("attachments_role") {
		params.AttachmentsRole = sentry.String(d.Get("attachments_role").(string))
	}
	if d.HasChange("debug_files_role") {
		params.DebugFilesRole = sentry.String(d.Get("debug_files_role").(string))
	}
	if d.HasChange("avatar_type") {
		params.AvatarType = sentry.String(d.Get("avatar_type").(string))
	}
	if d.HasChange("avatar") {
		params.Avatar = sentry.String(d.Get("avatar").(string))
	}
	if d.HasChange("require_2fa") {
		params.Require2FA = sentry.Bool(d.Get("require_2fa").(bool))
	}
	if d.HasChange("allow_shared_issues") {
		params.AllowSharedIssues = sentry.Bool(d.Get("allow_shared_issues").(bool))
	}
	if d.HasChange("enhanced_privacy") {
		params.EnhancedPrivacy = sentry.Bool(d.Get("enhanced_privacy").(bool))
	}
	if d.HasChange("scrape_javascript") {
		params.ScrapeJavaScript = sentry.Bool(d.Get("scrape_javascript").(bool))
	}
	if d.HasChange("store_crash_reports") {
		params.StoreCrashReports = sentry.Int(d.Get("store_crash_reports").(int))
	}
	if d.HasChange("allow_join_requests") {
		params.AllowJoinRequests = sentry.Bool(d.Get("allow_join_requests").(bool))
	}
	if d.HasChange("data_scrubber") {
		params.DataScrubber = sentry.Bool(d.Get("data_scrubber").(bool))
	}
	if d.HasChange("data_scrubber_defaults") {
		params.DataScrubberDefaults = sentry.Bool(d.Get("data_scrubber_defaults").(bool))
	}
	if d.HasChange("sensitive_fields") {
		fields := d.Get("sensitive_fields").([]interface{})
		s := make([]string, len(fields))
		for i, v := range fields {
			s[i] = v.(string)
		}
		params.SensitiveFields = s
	}
	if d.HasChange("safe_fields") {
		fields := d.Get("safe_fields").([]interface{})
		s := make([]string, len(fields))
		for i, v := range fields {
			s[i] = v.(string)
		}
		params.SafeFields = s
	}
	if d.HasChange("scrub_ip_addresses") {
		params.ScrubIPAddresses = sentry.Bool(d.Get("scrub_ip_addresses").(bool))
	}
	if d.HasChange("relay_pii_config") {
		params.RelayPiiConfig = sentry.String(d.Get("relay_pii_config").(string))
	}
	if d.HasChange("github_pr_bot") {
		params.GithubPRBot = sentry.Bool(d.Get("github_pr_bot").(bool))
	}
	if d.HasChange("github_open_pr_bot") {
		params.GithubOpenPRBot = sentry.Bool(d.Get("github_open_pr_bot").(bool))
	}
	if d.HasChange("github_nudge_invite") {
		params.GithubNudgeInvite = sentry.Bool(d.Get("github_nudge_invite").(bool))
	}
	if d.HasChange("gitlab_pr_bot") {
		params.GitlabPRBot = sentry.Bool(d.Get("gitlab_pr_bot").(bool))
	}
	if d.HasChange("allow_member_project_creation") {
		params.AllowMemberProjectCreation = sentry.Bool(d.Get("allow_member_project_creation").(bool))
	}

	if d.HasChange("trusted_relays") {
		relaysList := d.Get("trusted_relays").([]interface{})
		apiRelays := make([]sentry.TrustedRelayUpdateParams, len(relaysList))
		for i, r := range relaysList {
			relayMap := r.(map[string]interface{})
			apiRelays[i] = sentry.TrustedRelayUpdateParams{
				Name:        sentry.String(relayMap["name"].(string)),
				PublicKey:   sentry.String(relayMap["public_key"].(string)),
				Description: sentry.String(relayMap["description"].(string)),
			}
		}
		params.TrustedRelays = apiRelays
	}

	tflog.Debug(ctx, "Updating organization", map[string]interface{}{"org": org, "params": params})
	organization, _, err := client.Organizations.Update(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sentry.StringValue(organization.Slug))
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client
	org := d.Id()

	tflog.Debug(ctx, "Deleting organization", map[string]interface{}{"org": org})
	_, err := client.Organizations.Delete(ctx, org)
	return diag.FromErr(err)
}
