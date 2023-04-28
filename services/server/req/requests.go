package req

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"service/services/common/structs"
	"service/services/server/constants"
	errors2 "service/services/server/errors"
	t "service/services/translations"
)

// Request A stripped-down version of gin.Context.
//
// All response functions are hidden from

type Request[Body any] struct {
	Body func() Body
	// Deprecated: Request.Context is not mockable.
	// You are encouraged to use the predefined functions
	// or to extend the Request struct instead.
	Context   func() *gin.Context
	Get       func(key string) (value any, exists bool)
	Set       func(key string, value any)
	Header    func(key string) string
	Translate func(key t.TranslationKey, params ...string) string
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
		Translate: func(key t.TranslationKey, params ...string) string {
			lang := GetRequestLanguage(ctx)
			translator := t.GetTranslator(lang)
			translated, translationError := translator.T(key, params...)
			if translationError != nil {
				return string(key)
			} else {
				return translated
			}
		},
	}
	return request
}

func GetRequestLanguage(ctx *gin.Context) (reqLanguage string) {
	reqLanguageMaybeString, _ := ctx.Get(constants.RequestLanguage)
	if castedToString, castable := reqLanguageMaybeString.(string); castable {
		reqLanguage = castedToString
	} else {
		reqLanguage = "en"
	}
	return
}

func WrapRequestBindBody[T any](ctx *gin.Context) (r Request[T], e error) {
	var bindable T
	bindError := ctx.ShouldBind(&bindable)
	if bindError != nil {
		if _, isValidationError := bindError.(validator.ValidationErrors); isValidationError {
			return r, bindError
		} else if jsonTypeError, isJsonFieldTypeError := bindError.(*json.UnmarshalTypeError); isJsonFieldTypeError {
			singleFieldError := errors2.TypeError{
				FieldName:        jsonTypeError.Field,
				CurrentValueType: jsonTypeError.Value,
				TypeOf:           jsonTypeError.Type,
			}
			return r, validator.ValidationErrors([]validator.FieldError{singleFieldError})
		} else {
			return r, errors.Wrap(errors2.ErrCannotProcessReqBody, bindError.Error())
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
	requestHeaders *structs.StringDefaultedMap,
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
		Translate: func(key t.TranslationKey, params ...string) string {
			lang := ctx[constants.RequestLanguage].(string)
			translator := t.GetTranslator(lang)
			translated, translationError := translator.T(key, params...)
			if translationError != nil {
				return string(key)
			} else {
				return translated
			}
		},
	}
	return request
}
