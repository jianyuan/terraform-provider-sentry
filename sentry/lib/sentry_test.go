package sentry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client = NewClient(nil)
	url, _ := url.Parse(server.URL + "/api/")
	client.BaseURL = url
	return client, mux, server.URL, server.Close
}

// RewriteTransport rewrites https requests to http to avoid TLS cert issues
// during testing.
type RewriteTransport struct {
	Transport http.RoundTripper
}

// RoundTrip rewrites the request scheme to http and calls through to the
// composed RoundTripper or if it is nil, to the http.DefaultTransport.
func (t *RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}

func assertMethod(t *testing.T, expectedMethod string, req *http.Request) {
	assert.Equal(t, expectedMethod, req.Method)
}

// assertQuery tests that the Request has the expected url query key/val pairs
func assertQuery(t *testing.T, expected map[string]string, req *http.Request) {
	queryValues := req.URL.Query()
	expectedValues := url.Values{}
	for key, value := range expected {
		expectedValues.Add(key, value)
	}
	assert.Equal(t, expectedValues, queryValues)
}

// assertPostJSON tests that the Request has the expected JSON in its Body
func assertPostJSON(t *testing.T, expected interface{}, req *http.Request) {
	var actual interface{}

	d := json.NewDecoder(req.Body)
	d.UseNumber()

	err := d.Decode(&actual)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

// assertPostJSON tests that the request has the expected values in its body.
func assertPostJSONValue(t *testing.T, expected interface{}, req *http.Request) {
	var actual interface{}

	d := json.NewDecoder(req.Body)
	d.UseNumber()

	err := d.Decode(&actual)
	assert.NoError(t, err)
	assert.ObjectsAreEqualValues(expected, actual)
}

func mustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(fmt.Sprintf("mustParseTime: %s", err))
	}
	return t
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	assert.Equal(t, "https://sentry.io/api/", c.BaseURL.String())
}

func TestNewOnPremiseClient(t *testing.T) {
	testCases := []struct {
		baseURL string
	}{
		{"https://example.com"},
		{"https://example.com/"},
		{"https://example.com/api"},
		{"https://example.com/api/"},
	}
	for _, tc := range testCases {
		t.Run(tc.baseURL, func(t *testing.T) {
			c, err := NewOnPremiseClient(tc.baseURL, nil)

			assert.NoError(t, err)
			assert.Equal(t, "https://example.com/api/", c.BaseURL.String())
		})
	}

}

func TestResponse_populatePaginationCursor_hasNextResults(t *testing.T) {
	r := &http.Response{
		Header: http.Header{
			"Link": {`<https://sentry.io/api/0/organizations/terraform-provider-sentry/members/?&cursor=100:-1:1>; rel="previous"; results="false"; cursor="100:-1:1", ` +
				`<https://sentry.io/api/0/organizations/terraform-provider-sentry/members/?&cursor=100:1:0>; rel="next"; results="true"; cursor="100:1:0"`,
			},
		},
	}

	response := newResponse(r)
	assert.Equal(t, response.Cursor, "100:1:0")
}

func TestResponse_populatePaginationCursor_noNextResults(t *testing.T) {
	r := &http.Response{
		Header: http.Header{
			"Link": {`<https://sentry.io/api/0/organizations/terraform-provider-sentry/members/?&cursor=100:-1:1>; rel="previous"; results="false"; cursor="100:-1:1", ` +
				`<https://sentry.io/api/0/organizations/terraform-provider-sentry/members/?&cursor=100:1:0>; rel="next"; results="false"; cursor="100:1:0"`,
			},
		},
	}

	response := newResponse(r)
	assert.Equal(t, response.Cursor, "")
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	ctx := context.Background()
	client.Do(ctx, req, body)

	expected := &foo{A: "a"}

	assert.Equal(t, expected, body)
}

func TestDo_rateLimit(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerRateLimit, "40")
		w.Header().Set(headerRateRemaining, "39")
		w.Header().Set(headerRateReset, "1654566542")
		w.Header().Set(headerRateConcurrentLimit, "25")
		w.Header().Set(headerRateConcurrentRemaining, "24")
	})

	req, _ := client.NewRequest("GET", "/", nil)
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.Rate.Limit, 40)
	assert.Equal(t, resp.Rate.Remaining, 39)
	assert.Equal(t, resp.Rate.Reset, time.Date(2022, time.June, 7, 1, 49, 2, 0, time.UTC))
	assert.Equal(t, resp.Rate.ConcurrentLimit, 25)
	assert.Equal(t, resp.Rate.ConcurrentRemaining, 24)
}

func TestDo_nilContext(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(nil, req, nil)

	assert.Equal(t, errNonNilContext, err)
}

func TestDo_httpErrorPlainText(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)

	assert.Equal(t, &ErrorResponse{Response: resp.Response, Detail: "Bad Request"}, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDo_apiError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"detail": "API error message"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)

	assert.Equal(t, &ErrorResponse{Response: resp.Response, Detail: "API error message"}, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDo_apiError_noDetail(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `"API error message"`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)

	assert.Equal(t, &ErrorResponse{Response: resp.Response, Detail: "API error message"}, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCheckResponse(t *testing.T) {
	testcases := []struct {
		description string
		body        string
	}{
		{
			description: "JSON object",
			body:        `{"detail": "Error message"}`,
		},
		{
			description: "JSON string",
			body:        `"Error message"`,
		},
		{
			description: "plain text",
			body:        `Error message`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			res := &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader(tc.body)),
			}

			err := CheckResponse(res)

			expected := &ErrorResponse{
				Response: res,
				Detail:   "Error message",
			}
			assert.ErrorIs(t, err, expected)
		})
	}

}

func TestCheckResponse_rateLimit(t *testing.T) {
	testcases := []struct {
		description string
		addHeaders  func(res *http.Response)
	}{
		{
			description: "headerRateRemaining",
			addHeaders: func(res *http.Response) {
				res.Header.Set(headerRateRemaining, "0")
				res.Header.Set(headerRateReset, "123456")
			},
		},
		{
			description: "headerRateConcurrentRemaining",
			addHeaders: func(res *http.Response) {
				res.Header.Set(headerRateConcurrentRemaining, "0")
				res.Header.Set(headerRateReset, "123456")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			res := &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{},
				Body:       ioutil.NopCloser(strings.NewReader(`{"detail": "Rate limit exceeded"}`)),
			}
			tc.addHeaders(res)

			err := CheckResponse(res)

			expected := &RateLimitError{
				Rate:     ParseRate(res),
				Response: res,
				Detail:   "Rate limit exceeded",
			}
			assert.ErrorIs(t, err, expected)
		})
	}
}
