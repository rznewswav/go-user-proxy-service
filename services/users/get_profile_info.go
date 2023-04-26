package users

import "service/services/server/controllers"

var GetProfileInfo = controllers.C[any]().
	Get("/api/v1/me").
	UseMiddleware(&AuthMiddleware).
	Handle(func(
		body controllers.Request[any],
		SetStatus controllers.SetStatus,
		SetHeader controllers.SetHeader,
	) (Response any) {
		profile, _ := body.Context().Get(UserProfileToken)
		return profile
	})
