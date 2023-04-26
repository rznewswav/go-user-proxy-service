package controllers

type MockController[T any] struct {
	Controller[T]
}

func Mock[T any](controller Controller[T]) MockController[T] {
	cloned := controller.Clone()
	return MockController[T]{
		Controller: cloned,
	}
}

func (mc MockController[T]) ReplaceMiddleware(
	target *Handler[any],
	replaceWith *Handler[any],
) MockController[T] {
	for index, middleware := range mc.Middlewares {
		if target == middleware {
			mc.Middlewares[index] = replaceWith
		}
	}
	return mc
}
