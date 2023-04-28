package users

import (
	"github.com/gin-gonic/gin"
	"service/services/server/controllers"
	"service/services/server/req"
	"service/services/server/resp"
)

var GetProfileInfo = controllers.C[any]().
	Get("/api/v1/me").
	UseMiddleware(&AuthMiddleware).
	Handle(func(
		body req.Request[any],
	) (Response resp.Response) {
		profile, _ := body.Get(UserProfileToken)
		return resp.S(gin.H{
			"profile": profile,
		})

	})
