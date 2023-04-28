package resp

import (
	"github.com/gin-gonic/gin"
	"service/services/translations"
)

func (r response) Code(code string) Response {
	r.code = code
	return r
}

func (r response) Title(title t.TranslationKey) Response {
	r.title = title
	return r
}

func (r response) Message(message t.TranslationKey) Response {
	r.message = message
	return r
}

func Fail(code string, title, message t.TranslationKey, data gin.H) Response {
	r := S(data).(response)
	r.success = false
	r.code = code
	r.title = title
	r.message = message
	return r
}

// F alias for resp.Fail
func F(code string, title, message t.TranslationKey, data ...gin.H) Response {
	var datum gin.H
	if len(data) > 0 {
		datum = data[0]
	}
	return Fail(code, title, message, datum)
}
