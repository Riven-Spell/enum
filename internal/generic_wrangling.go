package internal

func FirstOrZero[T any](values []T) T {
	var out T
	if len(values) > 0 {
		out = values[0]
	}

	return out
}

func Ptr[T any](in T) *T {
	return &in
}
