package utils

func KeysToSlice[K comparable, V any](set map[K]V) []K {
	var result []K
	for key := range set {
		result = append(result, key)
	}
	return result
}

func ValuesToSlice[K comparable, V any](set map[K]V) []V {
	var result []V
	for _, value := range set {
		result = append(result, value)
	}
	return result
}
