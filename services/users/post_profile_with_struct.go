package users

import (
	"github.com/gin-gonic/gin"
	"service/services/server/controllers"
)

type PostProfileWithStructBody struct {
	Age             int    `json:"age" binding:"required"`
	FavouriteColour string `json:"favouriteColour" binding:"required"`
}

var PostProfileWithStruct = controllers.C[PostProfileWithStructBody]().
	UseMiddleware(&AuthMiddleware).
	Post("/api/v1/you").
	Handle(func(
		Request controllers.Request[PostProfileWithStructBody],
		SetStatus controllers.SetStatus,
		SetHeader controllers.SetHeader,
	) (Response any) {
		body := Request.Body()
		profile, _ := Request.Get(UserProfileToken)
		return gin.H{
			"body":    body,
			"profile": profile,
		}
	})
