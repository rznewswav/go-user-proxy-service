package logger

import (
	"service/services/bugsnag"
	"service/services/config"
	"strings"
)

func init() {
	c := config.QuietBuild(Config{})
	SetLogLevel(c.LogLevel)
	if c.PrettifyLogger {
		UsePrettyTransformer()
	} else {
		UseJsonTransformer()
	}
	logEvents := strings.Split(c.LogEvents, ",")
	LogEvents(logEvents)
	For("logger").Info("Log level is set to: %s", c.LogLevel)
	bugsnag.Start(For("bugsnag").Info)
}
