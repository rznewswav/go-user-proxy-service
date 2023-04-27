package controllers

import (
	"net/http"
)

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

type MockConfig struct {
	Key   string
	Value string
}

type MockAddHeader MockConfig
type MockBody any

func (mc MockController[T]) SendMockRequest(opt ...any) (
	response any,
	status int,
	header Headers,
) {
	status = http.StatusOK

	var body T
	var requestHeaders Headers

	for _, castable := range opt {
		switch casted := castable.(type) {
		case MockBody:
			body = casted.(T)
		case MockAddHeader:
			requestHeaders.Set(casted.Key, casted.Value)
		}
	}

	ctx := make(map[string]interface{})
	requestForMiddie := WrapRequestMockBody[any](
		body,
		ctx,
		&requestHeaders,
	)
	for _, middleware := range mc.Middlewares {
		response = (*middleware)(
			requestForMiddie,
			func(i int) {
				status = i
			},
			header.SetterFunc(),
		)

		if response != nil {
			return
		}
	}

	request := WrapRequestMockBody[T](
		body,
		ctx,
		&requestHeaders,
	)
	response = mc.Handler(
		request,
		func(i int) {
			status = i
		},
		header.SetterFunc(),
	)
	return
}
