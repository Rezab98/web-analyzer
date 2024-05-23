package slicetools

func Filter[T any](slice []T, shouldInclude func(t T) bool) []T {
	var result []T

	for i := range slice {
		if shouldInclude(slice[i]) {
			result = append(result, slice[i])
		}
	}

	return result
}
