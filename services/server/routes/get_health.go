package server_routes

import (
	"net/http"
	"service/services/health"

	"github.com/gin-gonic/gin"
)

func GetHealth(c *gin.Context) {
	healthStatus := health.IsHealthy()
	var healthStatusCode int
	if healthStatus.Healthy {
		healthStatusCode = http.StatusOK
	} else {
		healthStatusCode = http.StatusInternalServerError
	}
	c.JSON(healthStatusCode, healthStatus)
}
