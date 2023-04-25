package controllers

type Controller[T any] struct {
	Route       string
	Method      Method
	Middlewares []Handler[any]
	Handler     Handler[T]
}

func NewController[T any]() Controller[T] {
	return Controller[T]{}
}

func C[T any]() Controller[T] {
	return NewController[T]()
}

func (c Controller[T]) Get(Route string) Controller[T] {
	c.Method = GET
	c.Route = Route
	return c
}

func (c Controller[T]) Post(Route string) Controller[T] {
	c.Method = POST
	c.Route = Route
	return c
}

func (c Controller[T]) Put(Route string) Controller[T] {
	c.Method = PUT
	c.Route = Route
	return c
}

func (c Controller[T]) Patch(Route string) Controller[T] {
	c.Method = PATCH
	c.Route = Route
	return c
}

func (c Controller[T]) Delete(Route string) Controller[T] {
	c.Method = DELETE
	c.Route = Route
	return c
}

func (c Controller[T]) UseMiddleware(handler Handler[any]) Controller[T] {
	c.Middlewares = append(c.Middlewares, handler)
	return c
}

func (c Controller[T]) ResetMiddleware(handler Handler[any]) Controller[T] {
	c.Middlewares = make([]Handler[any], 0)
	return c
}

func (c Controller[T]) Handle(Handler Handler[T]) Controller[T] {
	c.Handler = Handler
	return c
}

func (c Controller[T]) Clone() (newController Controller[T]) {
	newController.Route = c.Route
	newController.Method = c.Method
	copy(newController.Middlewares, c.Middlewares)
	newController.Handler = c.Handler
	return
}
