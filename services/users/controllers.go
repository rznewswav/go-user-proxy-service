package users

import (
	"net/http"
	"service/services/server/controllers"
)

const UserProfileToken = "user"

var AuthMiddleware controllers.Handler[any] = func(
	body controllers.Request[any],
	SetStatus controllers.SetStatus,
	SetHeader controllers.SetHeader,
) (Response any) {
	nwToken := body.Context().GetHeader("nwtoken")

	success, profile := GetUserProfile(nwToken)
	if !success {
		SetStatus(http.StatusForbidden)
		return false
	}

	body.Context().Set(UserProfileToken, profile)
	return
}

var GetProfileInfo = controllers.C[any]().
	Get("/api/v1/me").
	UseMiddleware(AuthMiddleware).
	Handle(func(
		body controllers.Request[any],
		SetStatus controllers.SetStatus,
		SetHeader controllers.SetHeader,
	) (Response any) {
		profile, _ := body.Context().Get(UserProfileToken)
		return profile
	})
