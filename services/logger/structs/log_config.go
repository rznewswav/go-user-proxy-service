package logger_structs

type LoggerConfig struct {
	LogLevel       string `env:"LOG_LEVEL" default:"info" printDebug:"true"`
	PrettifyLogger bool   `env:"LOG_PRETTIFY"  default:"true" printDebug:"true"`
}
