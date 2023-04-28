package examples

import (
	"github.com/gin-gonic/gin"
	"service/services/server/constants"
	"service/services/server/controllers"
	"service/services/server/middlewares"
	"service/services/server/req"
	"service/services/server/resp"
)

type ConcateRequestBody struct {
	Age int `json:"age" binding:"required"`
}

var ConcatenateProfileInfo = controllers.
	Post[ConcateRequestBody]("/api/v1/me").
	UseMiddleware(&middlewares.AuthMiddleware).
	Handle(func(
		request req.Request[ConcateRequestBody],
	) (Response resp.Response) {
		body := request.Body()
		profile, _ := request.Get(constants.UserProfile)
		return resp.S(
			gin.H{
				"profile": profile,
				"body":    body,
			},
		)
	})
