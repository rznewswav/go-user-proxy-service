package redis

type RedisConfig struct {
	RedisHost     string `env:"REDIS_HOST" default:"127.0.0.1"`
	RedisPort     string `env:"REDIS_PORT" default:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" default:""`
}
