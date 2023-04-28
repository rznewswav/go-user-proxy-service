package controllers

import (
	"fmt"
	"service/services/common/utils"
	"service/services/server/handlers"
	"service/services/stack"
)

type Controller[T any] struct {
	Route        string
	Method       Method
	Middlewares  []*handlers.Handler[any]
	Handler      handlers.Handler[T]
	HandlerTrace string
}

func Get(Route string) (c Controller[any]) {
	c.Method = GET
	c.Route = Route
	return c
}

func Post[T any](Route string) (c Controller[T]) {
	c.Method = POST
	c.Route = Route
	return c
}

//goland:noinspection GoUnusedExportedFunction
func Put[T any](Route string) (c Controller[T]) {
	c.Method = PUT
	c.Route = Route
	return c
}

//goland:noinspection GoUnusedExportedFunction
func Patch[T any](Route string) (c Controller[T]) {
	c.Method = PATCH
	c.Route = Route
	return c
}

//goland:noinspection GoUnusedExportedFunction
func Delete[T any](Route string) (c Controller[T]) {
	c.Method = DELETE
	c.Route = Route
	return c
}

func (c Controller[T]) UseMiddleware(handler *handlers.Handler[any]) Controller[T] {
	c.Middlewares = append(c.Middlewares, handler)
	return c
}

func (c Controller[T]) ResetMiddleware() Controller[T] {
	c.Middlewares = make([]*handlers.Handler[any], 0)
	return c
}

func (c Controller[T]) Handle(Handler handlers.Handler[T]) Controller[T] {
	trace := stack.GetStackTrace()
	firstOfTrace := utils.ArrayGetOrNil(trace, 2)
	if firstOfTrace != nil {
		c.HandlerTrace = fmt.Sprintf("%s:%d", firstOfTrace.File, firstOfTrace.LineNumber)
	} else {
		c.HandlerTrace = ""
	}
	c.Handler = Handler
	return c
}

func (c Controller[T]) Clone() (newController Controller[T]) {
	newController.Route = c.Route
	newController.Method = c.Method
	newController.Middlewares = make([]*handlers.Handler[any], len(c.Middlewares))
	copy(newController.Middlewares, c.Middlewares)
	newController.Handler = c.Handler
	return
}
