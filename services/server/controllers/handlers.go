package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetStatus func(int)
type SetHeader func(string, string)

type Handler[T any] func(
	Request Request[T],
	SetStatus SetStatus,
	SetHeader SetHeader,
) (Response any)

func (s Handler[T]) AsGinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bindable T
		bindError := ctx.ShouldBind(&bindable)
		if bindError != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": false,
				"error": gin.H{
					"code": "BAD_REQUEST_BODY",
				},
			})
			return
		}

		requestObject := RequestStruct[T]{
			body:    bindable,
			context: ctx,
		}

		var status = http.StatusOK
		var headers Headers

		responseBody := s(
			requestObject,
			func(s int) {
				status = s
			},
			headers.SetterFunc(),
		)

		if responseBody != nil {
			ctx.JSON(status, responseBody)
		} else {
			ctx.Status(status)
		}
	}
}

func (s Handler[T]) AsGinMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var nilable T
		requestObject := RequestStruct[T]{
			body:    nilable,
			context: ctx,
		}

		var status = http.StatusOK
		var headers Headers

		responseBody := s(
			requestObject,
			func(s int) {
				status = s
			},
			headers.SetterFunc(),
		)

		if responseBody != nil {
			ctx.JSON(status, responseBody)
		} else {
			ctx.Next()
		}
	}
}
