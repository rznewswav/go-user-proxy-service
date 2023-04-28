package resp

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"net/http"
	"service/services/common/structs"
	"service/services/server/req"
	t "service/services/translations"
)

type Payload struct {
	Status int
	Header structs.StringDefaultedMap
	Data   gin.H
}
type Response interface {
	Success() bool
	Status(int) Response
	Header(string, string) Response
	Send(ctx *gin.Context)
	Next() bool
	GetResponsePayload(ut.Translator) Payload
}

type response struct {
	data    gin.H
	status  int
	code    string
	title   t.TranslationKey
	message t.TranslationKey
	next    bool
	success bool
	header  structs.StringDefaultedMap
}

func (r response) applyTitleMessageTranslation(translator ut.Translator) (hasTitleMessage bool) {

	// both r.title and r.message must present
	if len(r.title) == 0 || len(r.message) == 0 {
		hasTitleMessage = false
		return
	}
	hasTitleMessage = true
	translatedTitle, translatorTitleError := translator.T(r.title)
	if translatorTitleError != nil {
		translatedTitle = string(r.title)
	}

	translatedMessage, translatorMessageError := translator.T(r.message)
	if translatorMessageError != nil {
		translatedMessage = string(r.message)
	}

	r.data["title"] = translatedTitle
	r.data["message"] = translatedMessage
	return
}

func (r response) GetResponsePayload(translator ut.Translator) Payload {
	if r.applyTitleMessageTranslation(translator) {
		return Payload{
			Status: r.status,
			Header: r.header,
			Data: gin.H{
				"success": r.Success(),
				"title":   r.data["title"],
				"data":    r.data,
			},
		}
	} else {
		return Payload{
			Status: r.status,
			Header: r.header,
			Data: gin.H{
				"success": r.Success(),
				"data":    r.data,
			},
		}
	}
}

func (r response) Next() bool {
	return r.next
}

func (r response) Send(ctx *gin.Context) {
	r.header.ForEach(func(key, value string) {
		ctx.Header(key, value)
	})
	ctx.Status(r.status)

	if r.Next() {
		ctx.Next()
		return
	}

	lang := req.GetRequestLanguage(ctx)
	translator := t.GetTranslator(lang)
	payload := r.GetResponsePayload(translator)

	if r.Success() {
		ctx.JSON(payload.Status, payload.Data)
	} else {
		ctx.AbortWithStatusJSON(payload.Status, payload.Data)
	}
}

func (r response) Status(i int) Response {
	r.status = i
	return r
}

func (r response) Header(key, value string) Response {
	r.header.Set(key, value)
	return r
}

func (r response) Success() bool {
	return r.success
}

// S alias for resp.Success
func S(data ...any) Response {
	return Success(makeGinH(data...))
}

func Success(data gin.H) Response {
	return response{
		data:    smcopy(data),
		status:  http.StatusOK,
		success: true,
	}
}
