package sentry

// Bool returns a pointer to the bool value.
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the int value.
func Int(v int) *int {
	return &v
}
