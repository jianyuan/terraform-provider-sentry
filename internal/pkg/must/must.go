package must

func Do(err error) {
	if err != nil {
		panic(err)
	}
}

func Get[T any](v T, err error) T {
	Do(err)
	return v
}
