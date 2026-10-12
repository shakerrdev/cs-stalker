package db

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func setDBEnv(t *testing.T) {
	t.Helper()
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

	if _, err := DSNFromEnv(); err == nil {
		t.Fatal("expected an error for missing DB_NAME")
	}
}
