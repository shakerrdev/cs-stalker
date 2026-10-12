# cs-stalker

Look up a Counter-Strike player by their Steam profile and pull their Steam and Leetify data.

## Running

Needs Go 1.24+ and PostgreSQL 13+.

```sh
cp .env.example .env
go run ./cmd/cs-stalker
```

`make run` does the same, and `task watch` restarts the server whenever a Go, HTML or CSS file changes. The server listens on http://localhost:8085 and creates or updates its tables on startup.

## Environment

| Variable | Purpose |
| --- | --- |
| `STEAMAPIKEY` | Steam Web API key, from https://steamcommunity.com/dev/apikey |
| `LEETIFY_AUTH_TOKEN` | Bearer token sent to `api.leetify.com` |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Postgres connection |

Variables already set in the environment take precedence over `.env`, and the file is optional.

## Tests

```sh
go test ./...
```

The database tests are skipped unless `TEST_DATABASE_URL` points at a Postgres database. They run inside a transaction that is rolled back:

```sh
TEST_DATABASE_URL="host=localhost user=postgres dbname=cs_stalker_test sslmode=disable" go test ./internal/db
```
