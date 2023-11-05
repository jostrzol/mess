package utils

import "fmt"

func SingleEntry[K comparable, V any](m map[K]V) (K, V) {
	if len(m) != 1 {
		panic(fmt.Errorf("expected exactly 1 entry, got %v", len(m)))
	}
	for k, v := range m {
		return k, v
	}
	panic("unreachable")
}

func Entries[K comparable, V any](m map[K]V) (keys []K, values []V) {
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return
}
