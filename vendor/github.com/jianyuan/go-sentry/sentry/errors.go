package sentry

import "fmt"

// APIError represents a Sentry API Error response
type APIError map[string]interface{}

// TODO: use this instead
// type apiError struct {
// 	Detail string `json:"detail"`
// }

func (e APIError) Error() string {
	return fmt.Sprintf("sentry: %v", e)
}

// Empty returns true if empty.
func (e APIError) Empty() bool {
	return len(e) == 0
}

func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if !apiError.Empty() {
		return apiError
	}
	return nil
}
