package health

import (
	"net/http"
	"service/services/server/controllers"
)

var wasLastHealthy = false

var HealthController = controllers.C[any]().
	Get("/api/health").
	Handle(func(
		body controllers.Request[any],
		SetStatus controllers.SetStatus,
		SetHeader controllers.SetHeader,
	) (
		Response any,
	) {
		isHealthy := IsHealthy()
		wasLastHealthy = isHealthy.Healthy
		if isHealthy.Healthy {
			SetStatus(http.StatusOK)
		} else {
			SetStatus(http.StatusServiceUnavailable)
		}

		return isHealthy
	})

// Use for endpoints that require other services to be ready before handling requests
var HealthMiddleware controllers.Handler[any] = func(
	body controllers.Request[any],
	SetStatus controllers.SetStatus,
	SetHeader controllers.SetHeader,
) (Response any) {
	Response = nil
	if wasLastHealthy {
		return
	}

	isHealthy := IsHealthy()
	wasLastHealthy = isHealthy.Healthy
	if !isHealthy.Healthy {
		SetStatus(http.StatusServiceUnavailable)
		Response = map[string]string{
			"error": "Service is not available",
		}
		return
	}

	return
}
