package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/mariomac/gostream/stream"
	"net/http"
	"service/services/bugsnag"
	"service/services/logger"
	"service/services/server/resp"
	t "service/services/translations"
)

type SetStatus func(int)
type SetHeader func(string, string)
type TranslateOneFunc = func(string, ...string) string

type Handler[T any] func(
	Request Request[T],
	SetStatus SetStatus,
	SetHeader SetHeader,
	Translate TranslateOneFunc,
) (Response any)

func (h Handler[T]) AsGinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestObject, wrapError := WrapRequestBindBody[T](ctx)
		reqLanguage := getRequestLanguage(ctx)
		translator := t.GetTranslator(reqLanguage)
		var translateFunc = translateOneFunc(translator)

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
			} else if errors.Is(wrapError, ErrCannotProcessReqBody) {
				ctx.JSON(http.StatusUnprocessableEntity, resp.RespRequestBodyMalformed)
				return
			} else {
				ctx.JSON(http.StatusServiceUnavailable, resp.RespServiceUnavailableDefault)
				return
			}
		}

		var status = http.StatusOK
		var headers Headers

		responseBody := h(
			requestObject,
			func(s int) {
				status = s
			},
			headers.SetterFunc(),
			translateFunc,
		)

		reply := ctx.JSON
		emptyReply := ctx.Status
		if !IsOk(status) {
			reply = ctx.AbortWithStatusJSON
			emptyReply = ctx.AbortWithStatus
		}

		if responseBody != nil {
			wrappedResponseBody := validateResponseBody(
				responseBody,
				status,
				translateFunc,
			)
			reply(status, wrappedResponseBody)
		} else {
			emptyReply(status)
		}
	}
}

func getRequestLanguage(ctx *gin.Context) (reqLanguage string) {
	reqLanguageMaybeString, _ := ctx.Get(RequestLanguage)
	if castedToString, castable := reqLanguageMaybeString.(string); castable {
		reqLanguage = castedToString
	} else {
		reqLanguage = "en"
	}
	return
}

func translateOneFunc(ut ut.Translator) TranslateOneFunc {
	return func(s string, others ...string) string {
		onTranslationError := s
		if len(others) > 0 {
			onTranslationError = others[0]
		}

		if translated, err := ut.T(s); err != nil {
			return onTranslationError
		} else {
			return translated
		}
	}
}

func validateResponseBody(body any, status int, translate TranslateOneFunc) any {
	isOk := IsOk(status)
	if isOk {
		return gin.H{
			"status": IsOk(status),
			"data":   body,
		}
	}

	if errResponse, castable := body.(resp.ErrorResponse); castable {
		return errResponse.Data()
	} else {
		l := logger.For("request")
		l.Warn(
			"non-success status code is returned %d but response is not type of ErrorResponse struct",
			status,
			bugsnag.New("Response Type Error"),
		)
		return gin.H{
			"status": false,
		}
	}
}

func IsOk(status int) bool {
	return status >= 200 && status < 300
}

func (h Handler[T]) AsGinMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestObject := WrapRequest[T](ctx)
		reqLanguage := getRequestLanguage(ctx)
		translator := t.GetTranslator(reqLanguage)
		var translateFunc = translateOneFunc(translator)
		var status = http.StatusOK
		var headers Headers

		responseBody := h(
			requestObject,
			func(s int) {
				status = s
			},
			headers.SetterFunc(),
			translateFunc,
		)

		if responseBody != nil {
			ctx.JSON(status, gin.H{
				"status": IsOk(status),
				"data":   responseBody,
			})
		} else {
			ctx.Next()
		}
	}
}
