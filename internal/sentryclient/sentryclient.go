package sentryclient

import (
	"context"
	"net/http"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	UserAgent string
	Token     string
}

// Client to connect to Sentry.
func (c *Config) HttpClient(ctx context.Context) *http.Client {
	transport := http.DefaultTransport

	// Handle authentication
	transport = NewBearerTokenRoundTripper(transport, c.Token)

	// Handle user agent
	transport = NewUserAgentRoundTripper(transport, c.UserAgent)

	// Handle concurrency limit
	transport = NewSemaphoreRoundTripper(transport)

	// Handle rate limit
	transport = NewRateLimiterRoundTripper(transport)

	return &http.Client{
		Transport: transport,
	}
}
