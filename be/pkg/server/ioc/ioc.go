package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
)

func MustResolve[T any]() *T {
	var result *T
	container.MustResolve(container.Global, &result)
	return result
}

func MustSingleton[T any]() {
	container.MustSingletonLazy(container.Global, func() *T {
		var result T
		container.MustFill(container.Global, &result)
		return &result
	})
}

func MustSingletonAs[TSrc, TDst any]() {
	container.MustSingletonLazy(container.Global, func() TDst {
		var result TSrc
		container.MustFill(container.Global, &result)
		return (interface{})(&result).(TDst)
	})
}

var HandlerInitializers []func(*gin.Engine)

func MustHandler[T any](addHandlerFuncs ...func(*T, *gin.Engine)) {
	container.MustSingletonLazy(container.Global, func() *T {
		var result T
		container.MustFill(container.Global, &result)
		return &result
	})

	HandlerInitializers = append(HandlerInitializers, func(g *gin.Engine) {
		handler := MustResolve[T]()
		for _, addHandlerFunc := range addHandlerFuncs {
			addHandlerFunc(handler, g)
		}
	})
}
