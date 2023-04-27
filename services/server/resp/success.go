package resp

import "github.com/gin-gonic/gin"

type Response interface {
	Success() bool
}

type S gin.H

func (s S) Success() bool {
	return true
}

func (s F) Success() bool {
	return false
}

func Success(data gin.H) Response {
	return S{
		"success": true,
		"data":    data,
	}
}
