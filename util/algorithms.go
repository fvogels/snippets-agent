package util

func Map[T any, R any](xs []T, transformer func(t T) R) []R {
	result := make([]R, len(xs))

	for index, x := range xs {
		result[index] = transformer(x)
	}

	return result
}
