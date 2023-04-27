package app

import "service/services/shutdown"

func InitShutdownHook() *shutdown.Handler {
	return shutdown.NewHandler()
}
