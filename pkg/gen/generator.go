package gen

func Generator[T any](slice []T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, element := range slice {
			ch <- element
		}
	}()
	return ch
}
