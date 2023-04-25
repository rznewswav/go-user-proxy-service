package app

import (
	"service/services/health"
	"service/services/mysql"
	"service/services/redis"
	"service/services/server"
)

func Start() {
	shutdownHandler := InitShutdownHook()
	onFinishShutdown := InitSignal(shutdownHandler)

	mysql.Connect(shutdownHandler)
	redis.Connect(shutdownHandler)
	health.NewPatient("mysql", mysql.Health)
	health.NewPatient("redis", redis.Health)
	server.Start(shutdownHandler)
	shutdownHandler.Start(onFinishShutdown)
}
