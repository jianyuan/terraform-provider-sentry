package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
	"github.com/oapi-codegen/nullable"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/samber/lo"
)

func (r *OrganizationResource) getCreateJSONRequestBody(ctx context.Context, data OrganizationResourceModel) (*apiclient.CreateOrganizationJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := apiclient.CreateOrganizationJSONRequestBody{
		Name:       data.Name.ValueString(),
		AgreeTerms: data.AgreeTerms.ValueBool(),
	}
	if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
		body.Slug = lo.ToPtr(data.Slug.ValueString())
	}

	return &body, diags
}

func (r *OrganizationResource) getUpdateJSONRequestBody(ctx context.Context, data OrganizationResourceModel) (*apiclient.UpdateOrganizationJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := apiclient.UpdateOrganizationJSONRequestBody{
		Name: lo.ToPtr(data.Name.ValueString()),
	}

	if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
		body.Slug = lo.ToPtr(data.Slug.ValueString())
	}
	setBoolPtr := func(dst **bool, v supertypes.BoolValue) {
		if !v.IsNull() && !v.IsUnknown() {
			*dst = lo.ToPtr(v.ValueBool())
		}
	}
	setStringPtr := func(dst **string, v supertypes.StringValue) {
		if !v.IsNull() && !v.IsUnknown() {
			*dst = lo.ToPtr(v.ValueString())
		}
	}

	setBoolPtr(&body.IsEarlyAdopter, data.IsEarlyAdopter)
	setBoolPtr(&body.HideAiFeatures, data.HideAiFeatures)
	setBoolPtr(&body.CodecovAccess, data.CodecovAccess)
	setStringPtr(&body.DefaultRole, data.DefaultRole)
	setBoolPtr(&body.OpenMembership, data.OpenMembership)
	setBoolPtr(&body.EventsMemberAdmin, data.EventsMemberAdmin)
	setBoolPtr(&body.AlertsMemberWrite, data.AlertsMemberWrite)
	setStringPtr(&body.AttachmentsRole, data.AttachmentsRole)
	setStringPtr(&body.DebugFilesRole, data.DebugFilesRole)
	setStringPtr(&body.AvatarType, data.AvatarType)
	setStringPtr(&body.Avatar, data.Avatar)
	setBoolPtr(&body.Require2FA, data.Require2fa)
	setBoolPtr(&body.AllowSharedIssues, data.AllowSharedIssues)
	setBoolPtr(&body.EnhancedPrivacy, data.EnhancedPrivacy)
	setBoolPtr(&body.ScrapeJavaScript, data.ScrapeJavascript)
	if !data.StoreCrashReports.IsNull() && !data.StoreCrashReports.IsUnknown() {
		body.StoreCrashReports = lo.ToPtr(int(data.StoreCrashReports.ValueInt64()))
	}
	setBoolPtr(&body.AllowJoinRequests, data.AllowJoinRequests)
	setBoolPtr(&body.DataScrubber, data.DataScrubber)
	setBoolPtr(&body.DataScrubberDefaults, data.DataScrubberDefaults)
	setBoolPtr(&body.ScrubIPAddresses, data.ScrubIpAddresses)
	setBoolPtr(&body.GithubPRBot, data.GithubPrBot)
	setBoolPtr(&body.GithubOpenPRBot, data.GithubOpenPrBot)
	setBoolPtr(&body.GithubNudgeInvite, data.GithubNudgeInvite)
	setBoolPtr(&body.GitlabPRBot, data.GitlabPrBot)
	setBoolPtr(&body.AllowMemberProjectCreation, data.AllowMemberProjectCreation)

	if !data.RelayPiiConfig.IsUnknown() {
		if data.RelayPiiConfig.IsNull() {
			body.RelayPiiConfig = nullable.NewNullNullable[string]()
		} else {
			body.RelayPiiConfig = nullable.NewNullableWithValue(data.RelayPiiConfig.ValueString())
		}
	}

	if !data.SensitiveFields.IsNull() && !data.SensitiveFields.IsUnknown() {
		fields := tfutils.MergeDiagnostics(data.SensitiveFields.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}
		body.SensitiveFields = &fields
	}
	if !data.SafeFields.IsNull() && !data.SafeFields.IsUnknown() {
		fields := tfutils.MergeDiagnostics(data.SafeFields.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}
		body.SafeFields = &fields
	}

	if !data.TrustedRelays.IsNull() && !data.TrustedRelays.IsUnknown() {
		relays := tfutils.MergeDiagnostics(data.TrustedRelays.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}
		out := make([]apiclient.TrustedRelayUpdate, 0, len(relays))
		for _, relay := range relays {
			item := apiclient.TrustedRelayUpdate{}
			if !relay.Name.IsNull() && !relay.Name.IsUnknown() {
				item.Name = lo.ToPtr(relay.Name.ValueString())
			}
			if !relay.PublicKey.IsNull() && !relay.PublicKey.IsUnknown() {
				item.PublicKey = lo.ToPtr(relay.PublicKey.ValueString())
			}
			if !relay.Description.IsNull() && !relay.Description.IsUnknown() {
				item.Description = lo.ToPtr(relay.Description.ValueString())
			}
			out = append(out, item)
		}
		body.TrustedRelays = &out
	}

	return &body, diags
}

func (m *OrganizationResourceModel) Fill(ctx context.Context, data apiclient.Organization) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.Slug)
	m.Name = supertypes.NewStringValue(data.Name)
	m.Slug = supertypes.NewStringValue(data.Slug)
	m.AgreeTerms = supertypes.NewBoolValue(true)
	m.InternalId = supertypes.NewStringValue(data.Id)

	setBool := func(dst *supertypes.BoolValue, src *bool) {
		if src != nil {
			*dst = supertypes.NewBoolValue(*src)
		} else {
			dst.SetNull()
		}
	}
	setString := func(dst *supertypes.StringValue, src *string) {
		if src != nil {
			*dst = supertypes.NewStringValue(*src)
		} else {
			dst.SetNull()
		}
	}

	setBool(&m.IsEarlyAdopter, data.IsEarlyAdopter)
	setBool(&m.HideAiFeatures, data.HideAiFeatures)
	setBool(&m.CodecovAccess, data.CodecovAccess)
	setString(&m.DefaultRole, data.DefaultRole)
	setBool(&m.OpenMembership, data.OpenMembership)
	setBool(&m.EventsMemberAdmin, data.EventsMemberAdmin)
	setBool(&m.AlertsMemberWrite, data.AlertsMemberWrite)
	setString(&m.AttachmentsRole, data.AttachmentsRole)
	setString(&m.DebugFilesRole, data.DebugFilesRole)
	setBool(&m.Require2fa, data.Require2FA)
	setBool(&m.AllowSharedIssues, data.AllowSharedIssues)
	setBool(&m.EnhancedPrivacy, data.EnhancedPrivacy)
	setBool(&m.ScrapeJavascript, data.ScrapeJavaScript)
	setBool(&m.AllowJoinRequests, data.AllowJoinRequests)
	setBool(&m.DataScrubber, data.DataScrubber)
	setBool(&m.DataScrubberDefaults, data.DataScrubberDefaults)
	setBool(&m.ScrubIpAddresses, data.ScrubIPAddresses)
	setBool(&m.GithubPrBot, data.GithubPRBot)
	setBool(&m.GithubOpenPrBot, data.GithubOpenPRBot)
	setBool(&m.GithubNudgeInvite, data.GithubNudgeInvite)
	setBool(&m.GitlabPrBot, data.GitlabPRBot)
	setBool(&m.AllowMemberProjectCreation, data.AllowMemberProjectCreation)

	if data.StoreCrashReports != nil {
		m.StoreCrashReports = supertypes.NewInt64Value(int64(*data.StoreCrashReports))
	} else {
		m.StoreCrashReports.SetNull()
	}

	if data.Avatar != nil {
		setString(&m.AvatarType, data.Avatar.AvatarType)
	} else {
		m.AvatarType.SetNull()
	}

	if v, err := data.RelayPiiConfig.Get(); err == nil {
		m.RelayPiiConfig = supertypes.NewStringValue(v)
	} else {
		m.RelayPiiConfig.SetNull()
	}

	if data.SensitiveFields != nil {
		m.SensitiveFields = supertypes.NewListValueOfSlice(ctx, *data.SensitiveFields)
	} else {
		m.SensitiveFields = supertypes.NewListValueOfNull[string](ctx)
	}
	if data.SafeFields != nil {
		m.SafeFields = supertypes.NewListValueOfSlice(ctx, *data.SafeFields)
	} else {
		m.SafeFields = supertypes.NewListValueOfNull[string](ctx)
	}

	if data.TrustedRelays != nil {
		items := make([]OrganizationResourceModelTrustedRelaysItem, 0, len(*data.TrustedRelays))
		for _, relay := range *data.TrustedRelays {
			item := OrganizationResourceModelTrustedRelaysItem{}
			if relay.Name != nil {
				item.Name = supertypes.NewStringValue(*relay.Name)
			}
			if relay.PublicKey != nil {
				item.PublicKey = supertypes.NewStringValue(*relay.PublicKey)
			}
			if v, err := relay.Description.Get(); err == nil {
				item.Description = supertypes.NewStringValue(v)
			} else {
				item.Description.SetNull()
			}
			items = append(items, item)
		}
		m.TrustedRelays = supertypes.NewListNestedObjectValueOfValueSlice(ctx, items)
	} else {
		m.TrustedRelays = supertypes.NewListNestedObjectValueOfNull[OrganizationResourceModelTrustedRelaysItem](ctx)
	}

	return diags
}
