package providerdata

import (
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

type ProviderData struct {
	Client    *sentry.Client
	ApiClient *apiclient.ClientWithResponses
}
