package bugsnag

type Config struct {
	AppEnv        string `env:"APP_ENV" default:"development"`
	BugsnagApiKey string `env:"BUGSNAG_API_KEY"`
	NoReport      bool   `env:"BUGSNAG_NO_REPORT" default:"false"`
	Debug         bool   `env:"BUGSNAG_DEBUG" default:"false"`
}
