package utils

func FirstEntry[K comparable, V any](m map[K]V) (K, V) {
	for k, v := range m {
		return k, v
	}
	panic("map empty")
}

func Entries[K comparable, V any](m map[K]V) (keys []K, values []V) {
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return
}
