package app

import (
	"service/services/health"
	// "service/services/mysql"
	"service/services/redis"
	"service/services/server"
)

func Start() {
	shutdownHandler := InitShutdownHook()
	onFinishShutdown := InitSignal(shutdownHandler)

	// mysql.Connect(shutdownHandler)
	redis.Connect(shutdownHandler)
	// health.NewPatient("db", mysql.Health)
	health.NewPatient("cache", redis.Health)
	server.Start(shutdownHandler)
	shutdownHandler.Start(onFinishShutdown)
}
