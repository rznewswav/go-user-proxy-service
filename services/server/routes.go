package server

import (
	"service/services/health"
	"service/services/users"
)

func registerRoutes() {
	registerController(health.GetHealthController)
	registerController(users.GetProfileInfo)
	registerController(users.ConcatenateProfileInfo)
	registerController(users.PostProfileWithStruct)
}
