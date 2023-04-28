package examples

import (
	"github.com/gin-gonic/gin"
	"service/services/server/constants"
	"service/services/server/controllers"
	"service/services/server/middlewares"
	"service/services/server/req"
	"service/services/server/resp"
)

type PostProfileWithStructBody struct {
	Age             int    `json:"age" binding:"required"`
	FavouriteColour string `json:"favouriteColour" binding:"required"`
}

var PostProfileWithStruct = controllers.
	Post[PostProfileWithStructBody]("/api/v1/you").
	UseMiddleware(&middlewares.AuthMiddleware).
	Handle(func(
		Request req.Request[PostProfileWithStructBody],
	) (Response resp.Response) {
		body := Request.Body()
		profile, _ := Request.Get(constants.UserProfile)
		return resp.S(gin.H{
			"body":    body,
			"profile": profile,
		})
	})
