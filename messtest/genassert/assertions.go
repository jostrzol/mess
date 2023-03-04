package genassert

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Iterable[T any] interface {
	[]T | <-chan T
}

func Len[T any](t *testing.T, channel <-chan T, length int, msgsAndArgs ...interface{}) bool {
	t.Helper()
	slice := generatorToSlice(channel)
	return assert.Len(t, slice, length, msgsAndArgs...)
}

func Empty[T any](t *testing.T, channel <-chan T, msgsAndArgs ...interface{}) bool {
	t.Helper()
	slice := generatorToSlice(channel)
	return assert.Empty(t, slice, msgsAndArgs...)
}

func ElementsMatch[T any, TA Iterable[T], TB Iterable[T]](
	t *testing.T, iterableA TA, iterableB TB, msgsAndArgs ...interface{},
) bool {
	t.Helper()
	listA := generatorToSlice(iterableA)
	listB := generatorToSlice(iterableB)
	return assert.ElementsMatch(t, listA, listB, msgsAndArgs...)
}

func generatorToSlice(value interface{}) interface{} {
	val := reflect.ValueOf(value)
	if val.Type().Kind() == reflect.Chan {
		result := make([]interface{}, 0)
		ok := true
		for ok {
			var element reflect.Value
			_, element, ok = reflect.Select([]reflect.SelectCase{
				{Chan: val, Dir: reflect.SelectRecv},
			})
			result = append(result, element.Interface())
		}
		return result
	}
	return value
}
