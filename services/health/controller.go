package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/services/server/controllers"
	"service/services/server/handlers"
	"service/services/server/req"
	"service/services/server/resp"
)

var wasLastHealthy = false

var GetHealthController = controllers.C[any]().
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
				F(gin.H{
					"health": isHealthy,
				}).
				Status(http.StatusServiceUnavailable)
		} else {
			return resp.S(gin.H{
				"health": isHealthy,
			})
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
			F(gin.H{
				"health": isHealthy,
			}).
			Status(http.StatusServiceUnavailable)
	}

	return
}
