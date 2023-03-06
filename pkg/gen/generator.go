package gen

func FromSlice[T any](slice []T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, element := range slice {
			ch <- element
		}
	}()
	return ch
}

func FromValues[K comparable, V any](mab map[K]V) <-chan V {
	ch := make(chan V)
	go func() {
		defer close(ch)
		for _, value := range mab {
			ch <- value
		}
	}()
	return ch
}

func FromKeys[K comparable, V any](mab map[K]V) <-chan K {
	ch := make(chan K)
	go func() {
		defer close(ch)
		for key := range mab {
			ch <- key
		}
	}()
	return ch
}

func ToSlice[T any](generator <-chan T) []T {
	result := make([]T, 0)
	for element := range generator {
		result = append(result, element)
	}
	return result
}
