package server

import (
	"service/services/examples"
	"service/services/health"
)

func registerRoutes() {
	registerController(health.GetHealthController)
	registerController(examples.GetProfileInfo)
	registerController(examples.GetUserNameTranslate)
	registerController(examples.ConcatenateProfileInfo)
	registerController(examples.PostProfileWithStruct)
}
