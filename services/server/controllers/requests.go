package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// Request A stripped-down version of gin.Context.
//
// All response functions are hidden from

type Request[Body any] struct {
	Body func() Body
	// Deprecated: Request.Context is not mockable.
	// You are encouraged to use the predefined functions
	// or to extend the Request struct instead.
	Context func() *gin.Context
	Get     func(key string) (value any, exists bool)
	Set     func(key string, value any)
	Header  func(key string) string
}

func WrapRequest[T any](ctx *gin.Context) Request[T] {
	var bindable T
	//goland:noinspection GoDeprecation
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
		Header: func(key string) string {
			return ctx.GetHeader(key)
		},
	}
	return request
}

func WrapRequestBindBody[T any](ctx *gin.Context) (r Request[T], e error) {
	var bindable T
	bindError := ctx.ShouldBind(&bindable)
	if bindError != nil {
		if _, isValidationError := bindError.(validator.ValidationErrors); isValidationError {
			return r, bindError
		} else if jsonTypeError, isJsonFieldTypeError := bindError.(*json.UnmarshalTypeError); isJsonFieldTypeError {
			singleFieldError := TypeError{
				FieldName:        jsonTypeError.Field,
				CurrentValueType: jsonTypeError.Value,
				TypeOf:           jsonTypeError.Type,
			}
			return r, validator.ValidationErrors([]validator.FieldError{singleFieldError})
		} else {
			return r, errors.Wrap(ErrCannotProcessReqBody, bindError.Error())
		}
	}
	request := WrapRequest[T](ctx)
	request.Body = func() T {
		return bindable
	}
	return request, nil
}

func WrapRequestMockBody[T any](
	mockBody T,
	ctx map[string]interface{},
	requestHeaders *Headers,
) Request[T] {
	//goland:noinspection GoDeprecation
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
		Header: func(key string) string {
			return requestHeaders.Get(key)
		},
	}
	return request
}
