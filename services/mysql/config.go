package mysql

type Config struct {
	DatabaseURL string `env:"MYSQL_URL"`
}
