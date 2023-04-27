package resp

import (
	"github.com/gin-gonic/gin"
	"service/services/translations"
)

type ErrorResponse interface {
	Response
	Code(string) ErrorResponse
	Title(t.TranslationKey) ErrorResponse
	Message(t.TranslationKey) ErrorResponse
	Data() gin.H
}
type F struct {
	code    string
	title   t.TranslationKey
	message t.TranslationKey
	data    gin.H
}

func (s F) Data() gin.H {
	return s.data
}

func (s F) Code(code string) ErrorResponse {
	s.code = code
	return s
}

func (s F) Title(title t.TranslationKey) ErrorResponse {
	s.title = title
	return s
}

func (s F) Message(message t.TranslationKey) ErrorResponse {
	s.message = message
	return s
}

func Fail(data gin.H) ErrorResponse {
	populateMap := make(gin.H)
	smcopy(data, populateMap)
	return F{
		code:    "GENERIC_ERROR",
		title:   t.GenericErrorTitle,
		message: t.GenericErrorMessage,
		data:    data,
	}
}

// smcopy shallow copy map
func smcopy(data gin.H, populateMap gin.H) {
	for k, v := range data {
		populateMap[k] = v
	}
}
