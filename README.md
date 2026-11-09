# Campus-Backend

Backend for the TUM Campus Apps. One Go binary serves gRPC and REST and runs scheduled data imports into MySQL.

Production: [REST API](https://api.tum.app), [Swagger](https://api.tum.app/swagger/), gRPC at `api-grpc.tum.app:443`.
The API is defined in [`campus_backend.proto`](server/api/tumdev/campus_backend.proto).

## Local setup

Requires Docker.

```bash
git clone https://github.com/TUM-Dev/Campus-Backend.git
cd Campus-Backend
docker compose -f docker-compose.local.yml up --watch
```

This starts MySQL on port 3306 and the backend on port 50051, rebuilding on source changes.
Use `up -d` instead of `up --watch` to run in the background.
`docker-compose.yml` is for production; use `docker-compose.local.yml` locally.

```bash
curl localhost:50051/v1/canteen/allCanteens
grpcurl -plaintext localhost:50051 list
```

| Path | Purpose |
|---|---|
| `/v1/` | REST API |
| `/swagger/` | Swagger UI |
| `/metrics` | Prometheus metrics |
| `/health` | Health check |

## Development

Run Go commands from `server/`. To run outside Compose, set `DB_DSN` and use `go run .`.

```bash
go build ./...
go test ./...
bash api/generate.bash
```

Database tests require Docker or an existing MySQL server. Each test gets an isolated database.

```bash
TEST_MYSQL_DSN='user:password@tcp(localhost:3306)/' go test ./...
```

API changes start in the proto; regenerate rather than editing generated files. Never renumber or reuse field numbers.
`client/` is a separate Go module and is not covered by server tests.

For schema changes, edit the models, register new models in `migration.Models()`, then generate and review the SQL:

```bash
atlas migrate diff --env gorm <name>
```

Commit the SQL and `atlas.sum`. For manual SQL edits, run `atlas migrate hash --env gorm`.
Existing Go migration registries are frozen; new migrations belong in `backend/migration/sql/`.

Install [pre-commit](https://pre-commit.com) before contributing:

```bash
pip install pre-commit
pre-commit install
pre-commit run -a
```

CI runs tests, lint and migration checks. See [AGENTS.md](AGENTS.md) for the repository layout and change conventions,
and [`.vscode/launch.json`](.vscode/launch.json) for debug configurations.

## Admin dashboard

`/admin/` manages news feeds and shows change history, system/cron status, student clubs and councils with images,
and cafeteria ratings. Directory and rating views are read-only.

The dashboard activates when these settings are supplied. With all four unset, routes return 404; partial configuration fails startup.

```dotenv
ADMIN_BASE_URL=https://api.tum.app
ADMIN_OIDC_CLIENT_ID=<registered-client-id>
ADMIN_OIDC_CLIENT_SECRET=<registered-client-secret>
ADMIN_ALLOWED_SUBJECTS=<exact-subject-1>,<exact-subject-2>
```

Register a TUM Login OIDC client with authorization code flow, PKCE S256, query-mode callbacks,
`client_secret_basic`, `openid` scope and RS256 ID tokens. The callback is `${ADMIN_BASE_URL}/admin/callback`.
Confirm the issuer and endpoints with TUM support: discovery metadata can be inaccurate.
[Defaults](server/admin/config.go) can be overridden with `ADMIN_OIDC_ISSUER`, `ADMIN_OIDC_AUTH_URL`,
`ADMIN_OIDC_TOKEN_URL` and `ADMIN_OIDC_JWKS_URL`.

HTTPS is required. Allowlisted subjects are exact, case-sensitive `sub` claims, not email addresses.
Sessions expire within one hour; changes require CSRF tokens and are audited transactionally.
Logout ends only the admin session. Login state is encrypted in a short-lived browser cookie; replicas with
identical admin configuration do not need sticky routing. Rotating the OIDC client secret invalidates pending logins.
Feed saves advance a shared database cache version in the same transaction, so every replica refreshes public
news-source and article metadata on its next request. News API requests perform a version lookup even on cache hits.

News feeds require HTTPS and reject non-public download destinations. Pausing a feed does not stop an active import.
Feed icons can be reused from existing sources or from downloaded files in `news/sources/`.
Cron outcomes cover only the current process lifetime; `lastRun` records a start, not success.

## Configuration

| Variable | Purpose |
|---|---|
| `DB_DSN` | Required MySQL DSN; supplied by Compose for local use |
| `ENVIRONMENT` | `dev` or `prod` telemetry behavior |
| `SENTRY_DSN` | Error reporting |
| `OMDB_API_KEY` | Movie metadata |
| `SMTP_URL`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` | Feedback emails |

The committed `.env` contains local defaults only. Keep secrets in the ignored `.env.local` and load it explicitly:

```bash
docker compose --env-file .env.local -f docker-compose.yml up -d
```

## Notes

- Outside Docker, copy `server/api/tumdev/campus_backend.swagger.json` to `server/swagger/swagger.json` for Swagger UI.
- Legacy column names such as `lastRun` and `barrierFree_moreInfo` must match production.

## License

[GPL-3.0](LICENSE).
