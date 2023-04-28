package resp

import (
	"github.com/gin-gonic/gin"
	"service/services/server/req"
	"service/services/translations"
)

type ErrorResponse interface {
	Response
	Code(string) ErrorResponse
	Title(t.TranslationKey) ErrorResponse
	Message(t.TranslationKey) ErrorResponse
}
type f struct {
	code    string
	title   t.TranslationKey
	message t.TranslationKey
	s
}

func (f f) Code(code string) ErrorResponse {
	f.code = code
	return f
}

func (f f) Title(title t.TranslationKey) ErrorResponse {
	f.title = title
	return f
}

func (f f) Message(message t.TranslationKey) ErrorResponse {
	f.message = message
	return f
}

func (f f) Success() bool {
	return false
}

func (f f) applyTitleMessageTranslation(ctx *gin.Context) {
	lang := req.GetRequestLanguage(ctx)
	translator := t.GetTranslator(lang)
	translatedTitle, translatorTitleError := translator.T(f.title)
	if translatorTitleError != nil {
		translatedTitle = string(f.title)
	}

	translatedMessage, translatorMessageError := translator.T(f.message)
	if translatorMessageError != nil {
		translatedMessage = string(f.message)
	}

	f.data["title"] = translatedTitle
	data, hasData := f.data["data"]
	if !hasData {
		data := make(gin.H)
		f.data["data"] = data
	}

	data.(gin.H)["title"] = translatedTitle
	data.(gin.H)["message"] = translatedMessage
}

func (f f) Send(ctx *gin.Context) {
	f.applyTitleMessageTranslation(ctx)
	payload := f.GetResponsePayload()
	payload.Header.ForEach(func(key, value string) {
		ctx.Header(key, value)
	})
	ctx.AbortWithStatusJSON(payload.Status, payload.Data)
}

func Fail(data gin.H) ErrorResponse {
	return f{
		code:    "GENERIC_ERROR",
		title:   t.GenericErrorTitle,
		message: t.GenericErrorMessage,
		s:       S(data).(s),
	}
}

// F alias for resp.Fail
func F(data ...gin.H) ErrorResponse {
	var datum gin.H
	if len(data) > 0 {
		datum = data[0]
	}
	return Fail(datum)
}
