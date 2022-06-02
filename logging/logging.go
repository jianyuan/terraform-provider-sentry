package logging

import (
	"encoding/json"
	"net/http"
)

// logKeyValuePair is a helper struct to help organise the use of tflog. tflog
// uses pairs to print key-value items/args, so this is just to help visualize the
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

// extractHttpResponseElements is an indirection for if it is ever needed to add more elements.
func extractHttpResponseElements(resp *http.Response) []logKeyValuePair {
	return []logKeyValuePair{
		makePair("responseContentType", resp.Header.Get("content-type")),
		makePair("responseContentLength", resp.Header.Get("content-length")),
		makePair("responseStatus", resp.Status),
		makePair("requestUrl", resp.Request.URL.String()),
	}
}

// ExtractHttpResponse extracts key-value pairs from the http.Response object.
// This is to match the args signature of the tflog package.
func ExtractHttpResponse(resp *http.Response) map[string]interface{} {
	elements := extractHttpResponseElements(resp)
	args := make(map[string]interface{})
	for _, pair := range elements {
		args[pair.Key] = pair.Value
	}
	return args
}
