package server

type Config struct {
	AppPort      string `env:"APP_PORT" default:"3000"`
	AppEnv       string `env:"APP_ENV" default:"development"`
	AllowOrigins string `env:"CORS_DOMAINS" default:""`
}
