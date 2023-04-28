package middlewares

import (
	"service/services/server/constants"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
)

var AssignRequestLanguage handlers.Handler[any] = func(
	Request req.Request[any],
) (Response resp.Response) {
	mainLangInHeader := Request.Header(constants.RequestLanguage)
	if len(mainLangInHeader) == 0 {
		mainLangInHeader = "en"
	}

	Request.Set(constants.RequestLanguage, mainLangInHeader)
	return
}
