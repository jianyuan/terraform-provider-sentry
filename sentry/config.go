package sentry

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"golang.org/x/oauth2"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	Token   string
	BaseURL string
}

// Client to connect to Sentry.
func (c *Config) Client(ctx context.Context) (interface{}, diag.Diagnostics) {
	tflog.Info(ctx, "Instantiating Sentry client...")

	// Rate limit
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil // Disable DEBUG logs
	retryClient.CheckRetry = retryablehttp.ErrorPropagatedRetryPolicy
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if rateLimitErr, ok := sentry.CheckResponse(resp).(*sentry.RateLimitError); ok {
			return time.Until(rateLimitErr.Rate.Reset)
		}
		return retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
	}
	retryHTTPClient := retryClient.StandardClient()

	ctx = context.WithValue(ctx, oauth2.HTTPClient, retryHTTPClient)

	// Auth
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token})
	httpClient := oauth2.NewClient(ctx, ts)

	if c.BaseURL == "" {
		return sentry.NewClient(httpClient), nil
	} else {
		cl, err := sentry.NewOnPremiseClient(c.BaseURL, httpClient)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return cl, nil
	}
}
