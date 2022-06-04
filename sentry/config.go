package sentry

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/jianyuan/go-sentry/sentry"
	"golang.org/x/time/rate"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	Token   string
	BaseURL string
}

// Client to connect to Sentry.
func (c *Config) Client(ctx context.Context) (interface{}, diag.Diagnostics) {
	var baseURL *url.URL
	var err error

	if c.BaseURL != "" {
		tflog.Debug(ctx, "Parsing base url", map[string]interface{}{
			"BaseUrl": c.BaseURL,
		})
		baseURL, err = url.Parse(c.BaseURL)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	} else {
		tflog.Warn(ctx, "No base URL was set for the Sentry client")
	}

	tflog.Info(ctx, "Instantiating Sentry client...")
	client := &http.Client{
		Transport: &transport{
			// 40 requests every second.
			limiter: rate.NewLimiter(40, 1),
		},
	}

	cl := sentry.NewClient(client, baseURL, c.Token)

	return cl, nil
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
