package sentryclient

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func NewRateLimiterRoundTripper(delegate http.RoundTripper) http.RoundTripper {
	if delegate == nil {
		delegate = http.DefaultTransport
	}

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient = &http.Client{
		Transport: delegate,
	}
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	retryClient.Logger = nil // Disable DEBUG logs
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if rateLimitErr, ok := sentry.CheckResponse(resp).(*sentry.RateLimitError); ok {
			return time.Until(rateLimitErr.Rate.Reset)
		}
		return retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
	}
	retryHTTPClient := retryClient.StandardClient()

	return retryHTTPClient.Transport
}
