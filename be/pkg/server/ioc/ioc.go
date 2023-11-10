package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/pkg/server/core/event"
	"golang.org/x/exp/maps"
)

func MustResolve[T any]() T {
	var result T
	container.MustResolve(container.Global, &result)
	return result
}

func MustSingleton[T any](instance T) {
	container.MustSingletonLazy(container.Global, func() T {
		return instance
	})
}

func MustSingletonFill[T any]() {
	container.MustSingletonLazy(container.Global, func() *T {
		var result T
		container.MustFill(container.Global, &result)
		return &result
	})
}

func MustSingletonFillAs[TSrc, TDst any]() {
	container.MustSingletonLazy(container.Global, func() TDst {
		var result TSrc
		container.MustFill(container.Global, &result)
		return (interface{})(&result).(TDst)
	})
}

var HandlerInitializers []func(*gin.Engine)

func MustHandlerFill[T any](addHandlerFuncs ...func(*T, *gin.Engine)) {
	container.MustSingletonLazy(container.Global, func() *T {
		var result T
		container.MustFill(container.Global, &result)
		return &result
	})

	HandlerInitializers = append(HandlerInitializers, func(g *gin.Engine) {
		handler := MustResolve[*T]()
		for _, addHandlerFunc := range addHandlerFuncs {
			addHandlerFunc(handler, g)
		}
	})
}

func MustSingletonObserverFill[T any, PT interface {
	event.Observer
	*T
}]() {
	container.MustSingletonLazy(container.Global, func(broker *event.Broker) PT {
		var result T
		container.MustFill(container.Global, &result)
		ptr := PT(&result)
		broker.Observe(ptr)
		return ptr
	})
}

func MakeChild(parent container.Container) container.Container {
	child := container.New()
	for ty, bindMap := range parent {
		child[ty] = maps.Clone(bindMap)
	}
	return child
}
