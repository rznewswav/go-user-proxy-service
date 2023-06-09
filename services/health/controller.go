package health

import (
	"net/http"
	"service/services/server/controllers"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
	t "service/services/translations"
)

var wasLastHealthy = false

var GetHealthController = controllers.
	Get("/api/health").
	Handle(func(
		body req.Request[any],
	) (
		Response resp.Response,
	) {
		isHealthy := IsHealthy()
		wasLastHealthy = isHealthy.Healthy

		if !isHealthy.Healthy {
			return resp.
				F(
					"SERVICE_UNAVAILABLE",
					t.GenericErrorTitle,
					t.GenericErrorMessage,
					isHealthy,
				).
				Status(http.StatusServiceUnavailable)
		} else {
			return resp.S(isHealthy)
		}
	})

// GetHealthMiddleware Use for endpoints that require other services to be ready before handling requests
//
//goland:noinspection GoUnusedGlobalVariable
var GetHealthMiddleware handlers.Handler[any] = func(
	body req.Request[any],
) (Response resp.Response) {
	Response = nil
	if wasLastHealthy {
		return
	}

	isHealthy := IsHealthy()
	wasLastHealthy = isHealthy.Healthy
	if !isHealthy.Healthy {
		return resp.
			F(
				"SERVICE_UNAVAILABLE",
				t.GenericErrorTitle,
				t.GenericErrorMessage,
				isHealthy,
			).
			Status(http.StatusServiceUnavailable)
	}

	return
}
