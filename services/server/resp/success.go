package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/services/common/structs"
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
	GetResponsePayload() Payload
}

type s struct {
	data   gin.H
	status int
	header structs.StringDefaultedMap
}

func (s s) GetResponsePayload() Payload {
	return Payload{
		Status: s.status,
		Header: s.header,
		Data: gin.H{
			"success": s.Success(),
			"data":    s.data,
		},
	}
}

func (s s) Next() bool {
	return false
}

func (s s) Send(ctx *gin.Context) {
	payload := s.GetResponsePayload()
	payload.Header.ForEach(func(key, value string) {
		ctx.Header(key, value)
	})
	ctx.JSON(payload.Status, payload.Data)
}

func (s s) Status(i int) Response {
	s.status = i
	return s
}

func (s s) Header(key, value string) Response {
	s.header.Set(key, value)
	return s
}

func (s s) Success() bool {
	return true
}

// S alias for resp.Success
func S(data ...gin.H) Response {
	var datum gin.H
	if len(data) > 0 {
		datum = data[0]
	}
	return Success(datum)
}

func Success(data gin.H) Response {
	return s{
		data:   smcopy(data),
		status: http.StatusOK,
	}
}
