package app

import "service/services/shutdown"

func InitShutdownHook() *shutdown.ShutdownHandler {
	return shutdown.NewHandler()
}
