package logger

import (
	"service/services/config"
	logger_structs "service/services/logger/structs"
)

func init() {
	config := config.QuietBuild(logger_structs.LoggerConfig{})
	SetLogLevel(config.LogLevel)
	if config.PrettifyLogger {
		UsePrettyTransformer()
	} else {
		UseJsonTransformer()
	}
	WithContext("logger").Info("Log level is set to: %s", config.LogLevel)
}
