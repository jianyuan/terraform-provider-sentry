package logging

import (
	"encoding/json"
	"net/http"
)

// TryJsonify tries to marshall the object into a json string.
// If a failure should occur, it will return the untouched object.
func TryJsonify(object interface{}) interface{} {
	jsonObject, err := json.Marshal(object)
	if err == nil {
		return string(jsonObject)
	}
	return object
}

// ExtractHttpResponse extracts key-value pairs from the http.Response object.
// This is to match the args signature of the tflog package.
func ExtractHttpResponse(resp *http.Response) map[string]interface{} {
	return map[string]interface{}{
		"responseContentType":   resp.Header.Get("content-type"),
		"responseContentLength": resp.Header.Get("content-length"),
		"responseStatus":        resp.Status,
		"requestUrl":            resp.Request.URL.String(),
	}
}
