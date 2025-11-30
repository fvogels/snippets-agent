package util

func Map[T any, R any](xs []T, transformer func(t T) R) []R {
	result := make([]R, len(xs))

	for index, x := range xs {
		result[index] = transformer(x)
	}

	return result
}

func Filter[T any](xs []T, predicate func(t T) bool) []T {
	result := []T{}

	for _, x := range xs {
		if predicate(x) {
			result = append(result, x)
		}
	}

	return result
}
