package infra

import "testing"

func TestLoadConfigUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("POSTGRES_URL", "postgres://override:5432/test")
	t.Setenv("REDIS_URL", "redis://override:6379/1")
	t.Setenv("NATS_URL", "nats://override:4222")
	t.Setenv("HTTP_ADDR", ":8080")

	cfg := LoadConfig()

	if cfg.PostgresURL != "postgres://override:5432/test" {
		t.Fatalf("expected overridden postgres url, got %s", cfg.PostgresURL)
	}
	if cfg.RedisURL != "redis://override:6379/1" {
		t.Fatalf("expected overridden redis url, got %s", cfg.RedisURL)
	}
	if cfg.NATSURL != "nats://override:4222" {
		t.Fatalf("expected overridden nats url, got %s", cfg.NATSURL)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected overridden http addr, got %s", cfg.HTTPAddr)
	}
}
