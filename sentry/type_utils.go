package sentry

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the int value passed in.
func Int(v int) *int {
	return &v
}
