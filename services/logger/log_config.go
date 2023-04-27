package logger

type Config struct {
	LogLevel       string `env:"LOG_LEVEL" default:"info" printDebug:"true"`
	PrettifyLogger bool   `env:"LOG_PRETTIFY"  default:"true" printDebug:"true"`
	LogEvents      string `env:"LOG_EVENTS" default:""`
}
