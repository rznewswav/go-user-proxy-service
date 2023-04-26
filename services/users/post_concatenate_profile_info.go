package users

import "service/services/server/controllers"

var ConcatenateProfileInfo = controllers.C[map[string]interface{}]().
	Post("/api/v1/me").
	UseMiddleware(&AuthMiddleware).
	Handle(func(
		request controllers.Request[map[string]interface{}],
		SetStatus controllers.SetStatus,
		SetHeader controllers.SetHeader,
	) (Response any) {
		profile, _ := request.Context().Get(UserProfileToken)
		return map[string]interface{}{
			"profile": profile,
			"body":    request.Body(),
		}
	})
