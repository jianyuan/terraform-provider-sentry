package sentry

import (
	"encoding/json"
	"fmt"
)

// APIError represents a Sentry API Error response.
// Should look like:
//
//	type apiError struct {
//		Detail string `json:"detail"`
//	}
type APIError struct {
	f interface{} // unknown
}

func (e *APIError) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.f); err != nil {
		e.f = string(b)
	}
	return nil
}

func (e *APIError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.f)
}

func (e APIError) Detail() string {
	switch v := e.f.(type) {
	case map[string]interface{}:
		if len(v) == 1 {
			if detail, ok := v["detail"].(string); ok {
				return detail
			}
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (e APIError) Error() string {
	return fmt.Sprintf("sentry: %s", e.Detail())
}

// Empty returns true if empty.
func (e APIError) Empty() bool {
	return e.f == nil
}
