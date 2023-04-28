package examples

import (
	"service/services/server/constants"
	"service/services/server/controllers"
	"service/services/server/middlewares"
	"service/services/server/req"
	"service/services/server/resp"
	t "service/services/translations"
	"service/services/users"
)

// This controller demonstrates how to use translation keys
// in the response body.
//
// The function Request.Translate will use the request header
// "mainlanguage" to get the translation language of the request.
//
// Remember to register translations.

const (
	YourNameIs t.TranslationKey = "YourNameIs"
)

func init() {
	t.AddTranslations(
		YourNameIs,
		"Your name is {0}",
		"Nama anda ialah {0}",
		"ni de ming zi shi {0}",
	)
}

type GetUserNameTranslateResponse struct {
	Message string `json:"message,omitempty"`
}

var GetUserNameTranslate = controllers.
	Get("/api/v1/me/name").
	UseMiddleware(&middlewares.AuthMiddleware).
	Handle(
		func(Request req.Request[any]) (Response resp.Response) {
			profileUncasted, _ := Request.Get(constants.UserProfile)
			profile := profileUncasted.(users.AppUserType)
			return resp.S(GetUserNameTranslateResponse{
				Request.Translate(YourNameIs, profile.LoginDisplayName),
			})
		})
