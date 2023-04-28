package examples

import (
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
		profile, _ := request.Get(constants.UserProfile)
		return resp.S(
			map[string]interface{}{
				"profile": profile,
				"body":    request.Body(),
			},
		)
	})
