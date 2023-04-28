package users

import (
	"github.com/gin-gonic/gin"
	"service/services/server/controllers"
	"service/services/server/req"
	"service/services/server/resp"
)

type PostProfileWithStructBody struct {
	Age             int    `json:"age" binding:"required"`
	FavouriteColour string `json:"favouriteColour" binding:"required"`
}

var PostProfileWithStruct = controllers.C[PostProfileWithStructBody]().
	UseMiddleware(&AuthMiddleware).
	Post("/api/v1/you").
	Handle(func(
		Request req.Request[PostProfileWithStructBody],
	) (Response resp.Response) {
		body := Request.Body()
		profile, _ := Request.Get(UserProfileToken)
		return resp.S(gin.H{
			"body":    body,
			"profile": profile,
		})
	})
