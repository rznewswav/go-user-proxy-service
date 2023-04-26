package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Request A stripped down version of gin.Context.
//
// All response functions are hidden from

type Request[Body any] struct {
	Body    func() Body
	Context func() *gin.Context
	Get     func(key string) (value any, exists bool)
	Set     func(key string, value any)
}

func WrapRequest[T any](ctx *gin.Context) Request[T] {
	var bindable T
	request := Request[T]{
		Body: func() T {
			return bindable
		},
		Context: func() *gin.Context {
			return ctx
		},
		Get: func(key string) (value any, exists bool) {
			value, exists = ctx.Get(key)
			return
		},
		Set: func(key string, value any) {
			ctx.Set(key, value)
		},
	}
	return request
}

func WrapRequestBindBody[T any](ctx *gin.Context) (r Request[T], e error) {
	var bindable T
	bindError := ctx.ShouldBind(&bindable)
	if bindError != nil {
		return r, errors.Wrap(ErrCannotProcessReqBody, bindError.Error())
	}
	request := WrapRequest[T](ctx)
	request.Body = func() T {
		return bindable
	}
	return request, nil
}

func WrapRequestMockBody[T any](mockBody T) Request[T] {
	ctx := make(map[string]interface{})
	request := Request[T]{
		Body: func() T {
			return mockBody
		},
		Context: func() *gin.Context {
			return nil
		},
		Get: func(key string) (value any, exists bool) {
			value, exists = ctx[key]
			return
		},
		Set: func(key string, value any) {
			ctx[key] = value
		},
	}
	return request
}
