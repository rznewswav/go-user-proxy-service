package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mariomac/gostream/stream"
	"net/http"
	"service/services/logger"
	errors2 "service/services/server/errors"
	"service/services/server/req"
	"service/services/server/resp"
	t "service/services/translations"
)

type Handler[T any] func(Request req.Request[T]) (Response resp.Response)

func (h Handler[T]) AsGinHandler() gin.HandlerFunc {
	l := logger.For("handler")
	return func(ctx *gin.Context) {
		requestObject, wrapError := req.WrapRequestBindBody[T](ctx)
		reqLanguage := req.GetRequestLanguage(ctx)
		translator := t.GetTranslator(reqLanguage)

		if wrapError != nil {
			if validationError, isValidationError := wrapError.(validator.ValidationErrors); isValidationError && len(validationError) > 0 {
				validationErrorStream := stream.OfSlice(validationError)
				validationErrorStrings := stream.Map(validationErrorStream, func(it validator.FieldError) string {
					switch it.Tag() {
					case "required":
						translated, translationError := translator.T("field-required", it.Field())
						if translationError != nil {
							return it.Translate(translator)
						} else {
							return translated
						}
					default:
						return it.Translate(translator)

					}
				}).ToSlice()
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status": false,
					"data": gin.H{
						"code":         "BAD_REQUEST_BODY",
						"message":      validationErrorStrings[0],
						"verboseError": validationErrorStrings,
					},
				})
				return
			} else if errors.Is(wrapError, errors2.ErrCannotProcessReqBody) {
				ctx.JSON(http.StatusUnprocessableEntity, resp.RespRequestBodyMalformed)
				return
			} else {
				ctx.JSON(http.StatusServiceUnavailable, resp.RespServiceUnavailableDefault)
				return
			}
		}

		responseBody := h(
			requestObject,
		)

		if responseBody != nil {
			if response, castable := responseBody.(resp.Response); castable {
				response.Send(ctx)
			} else {
				l.Warn(
					"response body for route %s %s is not returning in resp.Response interface",
					ctx.Request.Method,
					ctx.FullPath(),
				)
				ctx.JSON(http.StatusOK, gin.H{
					"status": true,
					"data":   response,
				})
			}
		} else {
			ctx.Status(http.StatusNoContent)
		}
	}
}

func (h Handler[T]) AsGinMiddleware() gin.HandlerFunc {
	l := logger.For("middie")

	return func(ctx *gin.Context) {
		requestObject := req.WrapRequest[T](ctx)
		responseBody := h(requestObject)

		if responseBody != nil {
			if response, castable := responseBody.(resp.Response); castable {
				response.Send(ctx)
			} else {
				l.Warn(
					"response body for route %s %s is not returning in resp.Response interface",
					ctx.Request.Method,
					ctx.FullPath(),
				)
				ctx.JSON(http.StatusOK, gin.H{
					"status": true,
					"data":   response,
				})
			}
		} else {
			ctx.Next()
		}
	}
}
