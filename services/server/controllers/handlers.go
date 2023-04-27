package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ms"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	"github.com/mariomac/gostream/stream"
	"net/http"
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

var localeEn = en.New()
var localeZh = zh.New()
var localeMs = ms.New()
var uni = ut.New(localeEn, localeEn, localeZh, localeMs)

func init() {
	if enTranslator, found := uni.GetTranslator("en"); found {
		RegisterEnTranslations(enTranslator)
	}

	if msTranslator, found := uni.GetTranslator("ms"); found {
		RegisterMsTranslations(msTranslator)
	}

	if zhTranslator, found := uni.GetTranslator("zh"); found {
		RegisterZhTranslations(zhTranslator)
	}
}

func (h Handler[T]) AsGinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestObject, wrapError := WrapRequestBindBody[T](ctx)
		reqLanguageMaybeString, _ := ctx.Get(RequestLanguage)
		var reqLanguage string
		if castedToString, castable := reqLanguageMaybeString.(string); castable {
			reqLanguage = castedToString
		} else {
			reqLanguage = "en"
		}

		translator, _ := uni.GetTranslator(reqLanguage)

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
				ctx.JSON(http.StatusUnprocessableEntity, RespRequestBodyMalformed)
				return
			} else {
				ctx.JSON(http.StatusServiceUnavailable, RespServiceUnavailableDefault)
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
		)

		if responseBody != nil {
			ctx.JSON(status, responseBody)
		} else {
			ctx.Status(status)
		}
	}
}

func (h Handler[T]) AsGinMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestObject := WrapRequest[T](ctx)

		var status = http.StatusOK
		var headers Headers

		responseBody := h(
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
