package logging

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// logKeyValuePair is a helper struct to help organise the use of tflog. tflog
// uses pairs to print key-value items/args, so this just to help visualize the
// way the logging works
type logKeyValuePair struct {
	// Key is what will appear next to the value, e.g. keyID=...
	Key string
	// The value is what you want to log. e.g. Key=Value
	Value interface{}
}

func makePair(key string, value interface{}) logKeyValuePair {
	return logKeyValuePair{
		Key:   key,
		Value: value,
	}
}

// TryJsonify tries to marshall the object into a json string.
// If a failure should occur, it will return the untouched object.
func TryJsonify(object interface{}) interface{} {
	jsonObject, err := json.Marshal(object)
	if err == nil {
		return string(jsonObject)
	}
	return object
}

func AttachHttpResponse(ctx context.Context, resp *http.Response) context.Context {
	// Makes resp.Header (a bit) more usable than a map
	respHeaders := makePair("responseHeaders", TryJsonify(resp.Header))
	respStatus := makePair("responseStatus", resp.Status)
	reqUrl := makePair("requestUrl", resp.Request.URL.String())
	pairs := []logKeyValuePair{reqUrl, respStatus, respHeaders}

	for _, pair := range pairs {
		ctx = tflog.With(ctx, pair.Key, pair.Value)
	}
	return ctx
}

func extractHttpResponseElements(resp *http.Response) []logKeyValuePair {
	return []logKeyValuePair{
		makePair("responseHeaders", TryJsonify(resp.Header)),
		makePair("responseStatus", resp.Status),
		makePair("requestUrl", resp.Request.URL.String()),
	}
}

func ExtractHttpResponse(resp *http.Response) []interface{} {
	elements := extractHttpResponseElements(resp)
	args := make([]interface{}, 2*len(elements))
	for _, pair := range elements {
		args = append(args, pair.Key, pair.Value)
	}
	return args
}
