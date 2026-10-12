package db

import (
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func setDBEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "cs_stalker")
}

func TestDSNFromEnv(t *testing.T) {
	setDBEnv(t)

	for _, password := range []string{"", "with space", `it's`, `back\slash`} {
		t.Setenv("DB_PASSWORD", password)

		dsn, err := DSNFromEnv()
		if err != nil {
			t.Fatal(err)
		}
		config, err := pgconn.ParseConfig(dsn)
		if err != nil {
			t.Fatalf("parse %q: %v", dsn, err)
		}
		if config.Password != password || config.Database != "cs_stalker" {
			t.Errorf("password %q: parsed password = %q, database = %q", password, config.Password, config.Database)
		}
	}
}

func TestDSNFromEnv_Missing(t *testing.T) {
	setDBEnv(t)
	os.Unsetenv("DB_NAME")

	if _, err := DSNFromEnv(); err == nil || errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected a missing DB_NAME error, got %v", err)
	}
}

func TestDSNFromEnv_NotConfigured(t *testing.T) {
	setDBEnv(t)
	for _, key := range []string{"DATABASE_URL", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"} {
		os.Unsetenv(key)
	}

	if _, err := DSNFromEnv(); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("err = %v, want ErrNotConfigured", err)
	}
}

func TestDSNFromEnv_DatabaseURLWins(t *testing.T) {
	setDBEnv(t)
	url := "postgres://app:secret@db.example.com/cs?sslmode=require"
	t.Setenv("DATABASE_URL", url)

	dsn, err := DSNFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if dsn != url {
		t.Errorf("dsn = %q, want DATABASE_URL", dsn)
	}
}
