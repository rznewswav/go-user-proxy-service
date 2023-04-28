package middlewares

import (
	"net/http"
	"service/services/server/constants"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
	t "service/services/translations"
	"service/services/users"
)

var AuthMiddleware handlers.Handler[any] = func(
	body req.Request[any],
) (Response resp.Response) {
	nwToken := body.Header("nwtoken")

	success, profile := users.GetUserProfile(nwToken)
	if !success {
		return resp.F(
			"NOT_AUTHORIZED",
			t.NotAuthorizedTitle,
			t.NotAuthorized).
			Status(http.StatusForbidden)
	}

	body.Set(constants.UserProfile, profile)
	return
}
