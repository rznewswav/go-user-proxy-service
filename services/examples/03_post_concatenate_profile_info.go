package examples

import (
	"github.com/gin-gonic/gin"
	"service/services/server/constants"
	"service/services/server/controllers"
	"service/services/server/middlewares"
	"service/services/server/req"
	"service/services/server/resp"
)

var ConcatenateProfileInfo = controllers.
	Post[map[string]interface{}]("/api/v1/me").
	UseMiddleware(&middlewares.AuthMiddleware).
	Handle(func(
		request req.Request[map[string]interface{}],
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
