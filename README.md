# cs-stalker

Look up a Counter-Strike player by their Steam profile and pull their Steam and Leetify data.

## Running

Needs Go 1.24+. PostgreSQL 13+ is optional: without database settings the app runs without one, which is enough for lookups since nothing is stored yet.

```sh
cp .env.example .env
go run ./cmd/cs-stalker
```

`make run` does the same, and `task watch` restarts the server whenever a Go, HTML or CSS file changes. The server listens on http://localhost:8085 and, with a database configured, creates or updates its tables on startup.

## Environment

| Variable | Purpose |
| --- | --- |
| `STEAMAPIKEY` | Required. Steam Web API key, from https://steamcommunity.com/dev/apikey |
| `LEETIFY_AUTH_TOKEN` | Optional. Bearer token sent to `api.leetify.com` |
| `DATABASE_URL` | Optional. Postgres connection string, e.g. `postgres://user:pass@host/db?sslmode=require` |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Optional. Postgres connection in parts (all or none), used when `DATABASE_URL` is unset |
| `PORT` | Optional. Port to listen on, default 8085 |

Variables already set in the environment take precedence over `.env`, and the file is optional.

## Deploying

`render.yaml` describes a free [Render](https://render.com) web service in Frankfurt, built with Go from `main` and redeployed on every push.

1. Open https://render.com/deploy?repo=https://github.com/shakerrdev/cs-stalker and sign in with GitHub.
2. Enter your `STEAMAPIKEY` when asked and apply the blueprint.

The service runs without a database. To add one, create a free Postgres project on [Neon](https://neon.com), copy its connection string and add it as `DATABASE_URL` in the service's Environment settings; the tables are created on the next start. `LEETIFY_AUTH_TOKEN` can be added there too.

Free Render services go to sleep after a while without visits, so the first request after a break takes noticeably longer.

## Tests

```sh
go test ./...
```

The database tests are skipped unless `TEST_DATABASE_URL` points at a Postgres database. They run inside a transaction that is rolled back:

```sh
TEST_DATABASE_URL="host=localhost user=postgres dbname=cs_stalker_test sslmode=disable" go test ./internal/db
```
