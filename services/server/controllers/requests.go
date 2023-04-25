package controllers

import "github.com/gin-gonic/gin"

// A stripped down version of gin.Context.
//
// All response functions are hidden from
type Request[Body any] interface {
	Body() Body
	Context() *gin.Context
}

type RequestStruct[Body any] struct {
	Request[Body]
	body    Body
	context *gin.Context
}

func (r RequestStruct[Body]) Body() Body {
	return r.body
}

func (r RequestStruct[Body]) Context() *gin.Context {
	return r.context
}
