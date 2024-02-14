package sentry

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"golang.org/x/oauth2"
	"golang.org/x/sync/semaphore"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	UserAgent string
	Token     string
	BaseURL   string
}

// Client to connect to Sentry.
func (c *Config) Client(ctx context.Context) (interface{}, diag.Diagnostics) {
	tflog.Info(ctx, "Instantiating Sentry client...")

	// Authentication
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token})
	oauth2HTTPClient := oauth2.NewClient(ctx, ts)

	// Handle rate limit
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient = oauth2HTTPClient
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	retryClient.Logger = nil // Disable DEBUG logs
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if rateLimitErr, ok := sentry.CheckResponse(resp).(*sentry.RateLimitError); ok {
			return time.Until(rateLimitErr.Rate.Reset)
		}
		return retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
	}
	retryHTTPClient := retryClient.StandardClient()

	// Handle concurrency limit
	semaphoreHTTPClient := &http.Client{
		Transport: &semaphoreTransport{
			Delegate: retryHTTPClient.Transport,
		},
	}

	// Initialize client
	var cl *sentry.Client
	var err error
	if c.BaseURL == "" {
		cl = sentry.NewClient(semaphoreHTTPClient)
	} else {
		cl, err = sentry.NewOnPremiseClient(c.BaseURL, semaphoreHTTPClient)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	// Set user agent
	cl.UserAgent = c.UserAgent

	return cl, nil
}

type semaphoreTransport struct {
	Delegate http.RoundTripper

	mu sync.RWMutex
	w  *semaphore.Weighted
}

func (t *semaphoreTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.RLock()
	if t.w == nil {
		t.mu.RUnlock()
		t.mu.Lock()
		resp, err := t.Delegate.RoundTrip(req)
		if resp != nil {
			rate := sentry.ParseRate(resp)
			if rate.ConcurrentLimit > 0 {
				t.w = semaphore.NewWeighted(int64(rate.ConcurrentLimit))
			}
		}
		t.mu.Unlock()
		return resp, err
	}
	t.mu.RUnlock()

	ctx := req.Context()
	if err := t.w.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer t.w.Release(1)

	return t.Delegate.RoundTrip(req)
}
