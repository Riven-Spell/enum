package internal

// FirstOrZero returns the first value of a list, if any, else the zero of T.
func FirstOrZero[T any](values []T) T {
	var out T
	if len(values) > 0 {
		out = values[0]
	}

	return out
}

// Ptr returns &in.
func Ptr[T any](in T) *T {
	return &in
}
