package sentry

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	Token   string
	BaseURL string

	RateLimitPerSecond int
}

// Client to connect to Sentry.
func (c *Config) Client(ctx context.Context) (interface{}, diag.Diagnostics) {
	tflog.Info(ctx, "Instantiating Sentry client...")

	// Rate limit
	rateLimitHTTPClient := &http.Client{
		Transport: &transport{
			// 40 requests every second.
			limiter: rate.NewLimiter(rate.Limit(c.RateLimitPerSecond), 1),
		},
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, rateLimitHTTPClient)

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

type transport struct {
	limiter *rate.Limiter
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if err := t.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return http.DefaultTransport.RoundTrip(req)
}
