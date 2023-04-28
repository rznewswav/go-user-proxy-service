package users

import (
	"net/http"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
	t "service/services/translations"
)

const UserProfileToken = "user"

var AuthMiddleware handlers.Handler[any] = func(
	body req.Request[any],
) (Response resp.Response) {
	nwToken := body.Header("nwtoken")

	success, profile := GetUserProfile(nwToken)
	if !success {
		return resp.F().
			Title(t.NotAuthorizedTitle).
			Message(t.NotAuthorized).
			Code("NOT_AUTHORIZED").
			Status(http.StatusForbidden)
	}

	body.Set(UserProfileToken, profile)
	return
}
