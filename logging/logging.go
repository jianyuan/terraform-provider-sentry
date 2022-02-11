package logging

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// LogKeyValuePair is a helper struct to help organise the use of tflog. tflog
// uses pairs to print key-value items/args, so this just to help visualize the
// way the logging works
type LogKeyValuePair struct {
	// Key is what will appear next to the value, use it as a small hint/description
	Key string
	// The value is what you want to log.
	Value interface{}
}

func MakePair(key string, value interface{}) LogKeyValuePair {
	return LogKeyValuePair{
		Key:   key,
		Value: value,
	}
}

func AttachHttpResponse(ctx context.Context, resp *http.Response) context.Context {
	// Makes it more usable than a map
	jsonResponseHeaders, _ := json.Marshal(resp.Header)

	respHeaders := MakePair("responseHeaders", jsonResponseHeaders)
	respStatus := MakePair("responseStatus", resp.Status)
	reqUrl := MakePair("requestUrl", resp.Request.URL)
	pairs := []LogKeyValuePair{reqUrl, respStatus, respHeaders}

	for _, pair := range pairs {
		ctx = tflog.With(ctx, pair.Key, pair.Value)
	}
	return ctx
}
