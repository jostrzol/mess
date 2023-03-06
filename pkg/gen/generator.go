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
