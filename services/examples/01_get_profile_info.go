package examples

import (
	"service/services/server/constants"
	"service/services/server/controllers"
	"service/services/server/middlewares"
	"service/services/server/req"
	"service/services/server/resp"
)

type GetProfileInfoResponse struct {
	Profile any `json:"profile,omitempty"`
}

var GetProfileInfo = controllers.
	Get("/api/v1/me").
	UseMiddleware(&middlewares.AuthMiddleware).
	Handle(func(
		Request req.Request[any],
	) (Response resp.Response) {
		profile, _ := Request.Get(constants.UserProfile)

		return resp.S(GetProfileInfoResponse{
			profile,
		})
	})
