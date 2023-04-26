package controllers

import (
	"errors"
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

var RespServiceUnavailableDefault = gin.H{
	"status": false,
	"data":   gin.H{},
}

var RespRequestBodyMalformed = gin.H{
	"status": false,
	"data": gin.H{
		"message": "Your request body is not a valid JSON body.",
		"code":    "INVALID_JSON_BODY",
	},
}

func (s Handler[T]) AsGinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestObject, wrapError := WrapRequestBindBody[T](ctx)
		if wrapError != nil {
			if errors.Is(wrapError, ErrCannotProcessReqBody) {
				ctx.JSON(http.StatusUnprocessableEntity, RespRequestBodyMalformed)
				return
			} else {
				ctx.JSON(http.StatusServiceUnavailable, RespServiceUnavailableDefault)
				return
			}
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
		requestObject := WrapRequest[T](ctx)

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
