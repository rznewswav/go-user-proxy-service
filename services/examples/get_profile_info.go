package examples

import (
	"github.com/gin-gonic/gin"
	"service/services/server/constants"
	"service/services/server/controllers"
	"service/services/server/middlewares"
	"service/services/server/req"
	"service/services/server/resp"
)

var GetProfileInfo = controllers.
	Get("/api/v1/me").
	UseMiddleware(&middlewares.AuthMiddleware).
	Handle(func(
		body req.Request[any],
	) (Response resp.Response) {
		profile, _ := body.Get(constants.UserProfile)
		return resp.S(gin.H{
			"profile": profile,
		})

	})
