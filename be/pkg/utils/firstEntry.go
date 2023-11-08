package utils

import "fmt"

func Single[V any](slice []V) V {
	if len(slice) != 1 {
		panic(fmt.Errorf("expected exactly 1 element, got %v", len(slice)))
	}
	return slice[0]
}

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

func Transpose[V any](slice [][]V) (result [][]V) {
	height := len(slice)
	if height == 0 {
		return nil
	}
	width := len(slice[0])

	result = make([][]V, width)
	for _, row := range slice {
		for x, value := range row {
			result[x] = append(result[x], value)
		}
	}
	return
}
