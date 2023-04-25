package mysql

type MysqlConfig struct {
	DatabaseURL string `env:"MYSQL_URL"`
}
