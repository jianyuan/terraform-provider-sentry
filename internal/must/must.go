// Package must provides a simple way to unwrap errors.
package must

// Do panics if err is not nil.
func Do(err error) {
	if err != nil {
		panic(err)
	}
}

// Get returns v if err is nil, otherwise panics.
func Get[T any](v T, err error) T {
	Do(err)
	return v
}
