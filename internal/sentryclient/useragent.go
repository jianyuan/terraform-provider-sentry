package sentryclient

import "net/http"

func NewUserAgentRoundTripper(delegate http.RoundTripper, userAgent string) http.RoundTripper {
	if delegate == nil {
		delegate = http.DefaultTransport
	}

	return &UserAgentRoundTripper{
		delegate:  delegate,
		userAgent: userAgent,
	}
}

type UserAgentRoundTripper struct {
	delegate  http.RoundTripper
	userAgent string
}

func (t *UserAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.delegate.RoundTrip(req)
}
