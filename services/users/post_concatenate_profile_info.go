package users

import (
	"service/services/server/controllers"
	"service/services/server/req"
	"service/services/server/resp"
)

var ConcatenateProfileInfo = controllers.C[map[string]interface{}]().
	Post("/api/v1/me").
	UseMiddleware(&AuthMiddleware).
	Handle(func(
		request req.Request[map[string]interface{}],
	) (Response resp.Response) {
		profile, _ := request.Get(UserProfileToken)
		return resp.S(
			map[string]interface{}{
				"profile": profile,
				"body":    request.Body(),
			},
		)
	})
