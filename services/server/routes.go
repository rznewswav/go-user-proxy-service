package server

import (
	"service/services/health"
	"service/services/users"
)

func registerRoutes() {
	registerController(health.HealthController)
	registerController(users.GetProfileInfo)
	registerController(users.ConcatenateProfileInfo)
}
