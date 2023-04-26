package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Request A stripped down version of gin.Context.
//
// All response functions are hidden from
type Request[Body any] interface {
	Body() Body
	Context() *gin.Context
}

type RequestStruct[Body any] struct {
	Request[Body]
	Body    func() Body
	Context func() *gin.Context
}

func WrapRequest[T any](ctx *gin.Context) Request[T] {
	var bindable T
	request := RequestStruct[T]{
		Body: func() T {
			return bindable
		},
		Context: func() *gin.Context {
			return ctx
		},
	}
	request.Request = request
	return request
}

func WrapRequestBindBody[T any](ctx *gin.Context) (Request[T], error) {
	var bindable T
	bindError := ctx.ShouldBind(&bindable)
	if bindError != nil {
		return nil, errors.Wrap(ErrCannotProcessReqBody, bindError.Error())
	}
	request := RequestStruct[T]{
		Body: func() T {
			return bindable
		},
		Context: func() *gin.Context {
			return ctx
		},
	}
	request.Request = request
	return request, nil
}
