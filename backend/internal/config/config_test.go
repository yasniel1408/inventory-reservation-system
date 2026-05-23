package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("expected default database url")
	}
}

func TestLoadUsesEnvironment(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("DATABASE_URL", "postgres://custom")

	cfg := Load()

	if cfg.Port != "9000" {
		t.Fatalf("expected port from env, got %q", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://custom" {
		t.Fatalf("expected database url from env, got %q", cfg.DatabaseURL)
	}
}
