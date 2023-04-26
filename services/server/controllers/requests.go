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
	}
	return request
}

func WrapRequestBindBody[T any](ctx *gin.Context) (r Request[T], e error) {
	var bindable T
	bindError := ctx.ShouldBind(&bindable)
	if bindError != nil {
		return r, errors.Wrap(ErrCannotProcessReqBody, bindError.Error())
	}
	request := Request[T]{
		Body: func() T {
			return bindable
		},
		Context: func() *gin.Context {
			return ctx
		},
	}
	return request, nil
}
