package resp

import (
	"github.com/gin-gonic/gin"
)

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
