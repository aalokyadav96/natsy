package infra

import "os"

type Config struct {
	PostgresURL string
	RedisURL    string
	NATSURL     string
	HTTPAddr    string
}

func LoadConfig() Config {
	return Config{
		PostgresURL: getenv("POSTGRES_URL", "postgres://apeman:ningning@localhost:5432/eventdb?sslmode=disable"),
		RedisURL:    getenv("REDIS_URL", "redis://localhost:6379/0"),
		NATSURL:     getenv("NATS_URL", "nats://localhost:4222"),
		HTTPAddr:    getenv("HTTP_ADDR", ":4000"),
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
