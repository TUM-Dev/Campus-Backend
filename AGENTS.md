# AGENTS.md

Backend for the TUM Campus Apps (Android, iOS, Flutter). It is one Go binary that serves gRPC and REST on a single port, and it also runs cron jobs that scrape and download data into MySQL.
It runs in production at `api.tum.app`.

## Layout

| Path | What |
|---|---|
| `server/` | The Go module `github.com/TUM-Dev/Campus-Backend/server`. Run all `go` commands from here. |
| `server/main.go` | Entry point. Uses `cmux` to split port `:50051` into gRPC (HTTP/2, `content-type: application/grpc`) and HTTP/1, and runs the grpc-gateway under `/v1/`. Also serves `/health`, `/metrics` (Prometheus), `/swagger/` and `/files/` (from `/Storage`). |
| `server/api/tumdev/campus_backend.proto` | **The single API definition.** The `*.pb.go`, `*.pb.gw.go`, `*_grpc.pb.go` and `*.swagger.json` files next to it are generated. Don't edit them by hand. |
| `server/backend/` | `CampusServer` (see `rpcserver.go`) implements `pb.CampusServer`, with one file per domain (`news.go`, `cafeteria.go`, `movie.go`, …). |
| `server/backend/cron/` | `CronService`. Every 30 s it polls the `crontab` table and runs due jobs in goroutines. |
| `server/backend/migration/` | Schema migrations (see below). |
| `server/model/` | GORM models, one file per table, snake_case file names. |
| `server/utils/` | DB setup (`SetupDB`), test DB helper (`SetupTestContainer`), Sentry/logrus telemetry, gRPC→HTTP error mapping. |
| `server/swagger/` | Swagger UI assets, embedded via `go:embed`. `swagger.json` is copied in only by the Dockerfile. |
| `client/` | A separate Go module with example gRPC clients (`local/`, `public/`). It pins an old server version and is rarely touched. |

## Commands (from `server/`)

```bash
go build ./...
go test ./...                        # Docker required: one MySQL testcontainer per test binary
TEST_MYSQL_DSN='root:pw@tcp(localhost:3306)/' go test ./...   # reuse an existing MySQL instead
DB_DSN='user:pw@tcp(localhost:3306)/campus_db?charset=utf8mb4&parseTime=True&loc=Local' go run .
bash api/generate.bash               # regenerate protobuf/gateway/swagger (needs buf, goimports, protoc plugins in $(go env GOPATH)/bin)
```

For a full local stack, run `docker compose -f docker-compose.local.yml up -d` from the repo root.

## Making changes

### Adding or changing an endpoint
1. Edit `server/api/tumdev/campus_backend.proto`. Add the `rpc` to `service Campus` with a `google.api.http` option, plus the request and reply messages.
2. Run `bash api/generate.bash` and commit the generated files.
3. Implement the method on `*CampusServer` in `server/backend/<domain>.go`.

Notes:
- Most handlers start with `if err := s.checkDevice(ctx); err != nil { return nil, err }`. It records the `x-device-id`, `x-os-version` and `x-app-version` metadata in a buffer that is flushed to the `devices` table every minute.
- Return gRPC status errors (`status.Error(codes.X, "...")`). `utils.GrpcErrorToWebError` maps them to HTTP for REST clients. Don't leak internal error text in `codes.Internal` messages. Log the error with logrus instead.
- Pushing a proto change to `main` triggers `client_update.yml`, which opens a PR in `TUM-Dev/Campus-Flutter`. Keep changes backward compatible: never renumber or reuse field numbers, because released apps depend on them.

### Changing the database schema
- `migration.Migrate` runs at startup and records applied migrations in the `migrations` table:
  1. The legacy gormigrate migrations in `server/backend/migration/<YYYYMMDDhhmmss>.go`. These are frozen, so don't add new ones.
  2. The SQL migrations in `server/backend/migration/sql/`, which [Atlas](https://atlasgo.io) generates from the models.
- The desired schema is the GORM models (`migration.Models()`) plus `server/tools/atlasloader/legacy.sql`. `legacy.sql` declares the production tables that have no model, and the foreign keys whose production names GORM can't generate.
- To change the schema:
  1. Edit the model. If you added a new model, add it to `migration.Models()`.
  2. In `server/`, run `atlas migrate diff --env gorm <short_name>` (needs Docker). It writes `sql/<version>_<short_name>.sql` and updates `atlas.sum`.
  3. Review the generated SQL and commit it together with `atlas.sum`.
- For data changes, or anything Atlas doesn't manage (views, triggers), edit the generated file or create one with `atlas migrate new --env gorm <short_name>`, then run `atlas migrate hash --env gorm`.
- `test_migration.yml` runs the migrations like production and checks that the result matches the models. If the models changed without a migration, it fails and comments the missing SQL on the PR.
- Model conventions:
  - Use explicit `gorm:"column:…;type:…"` tags.
  - Foreign keys are declared on the relation field, e.g. `constraint:OnUpdate:CASCADE,OnDelete:CASCADE`. Relations tagged `constraint:-` have their foreign key in `legacy.sql` instead.
  - Use `guregu/null` types for nullable columns.
  - Tables are created with `ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`, set in `migration.TableOptions`.

### Adding a cron job
1. Add a type constant in `backend/cron/cronjobs.go`.
2. Add it to the `IN (...)` query and to the `switch`.
3. Add the value to the `type` enum in `model/crontab.go`.
4. Generate the migration with `atlas migrate diff --env gorm <short_name>`, append the `INSERT` for the `crontab` row to the generated file, then run `atlas migrate hash --env gorm`.

## Testing conventions
- Tests use `testify` (`require`, `suite`).
- Unit tests in `backend/` mostly use `go-sqlmock` with a GORM mysql dialector, and match queries with `regexp.QuoteMeta`. Copy the setup from an existing `*_test.go` (e.g. `news_test.go`).
- Tests that need a real DB call `utils.SetupTestContainer(ctx, t)`. Each test gets its own database, created from the migration baseline (`migration.MigrateFromBaseline`), so it has the production schema.
- Cron jobs write to `cron.StorageDir`, which defaults to `/Storage/`. Override it in tests.

## Lint / CI
- CI (`main.yml`): `go test -v ./...` in `server/`, and a Docker build that pushes to `ghcr.io/tum-dev/campus-backend/backend-server`.
- `lint.yml`:
  - golangci-lint v2.11.4 (default config) on `server` and `client`.
  - pre-commit: trailing whitespace, EOF newline, json/yaml checks, `go fmt`, `goimports`, protobuf regen, and `go mod tidy` in `server/`.
  - ls-lint: directories lowercase, `.go` files snake_case.
- `staticcheck.conf` turns off ST1000/ST1003.
- Go version is 1.27 everywhere (CI, Dockerfile, `go.mod`).

## Configuration (env vars)
- `DB_DSN` is required.
- Optional:
  - `ENVIRONMENT` (`prod` changes telemetry behavior)
  - `SENTRY_DSN`
  - `OMDB_API_KEY` (movie metadata)
  - `SMTP_URL`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` (feedback emails)
  - `CI_EXIT_AFTER_MIGRATION=true` (exit after migrating, used by CI)
- The root `.env` is committed and holds only local dev defaults. Never put real secrets in it. Use `.env.local`, which is gitignored.

## Gotchas
- When running locally, `/swagger/` has no `swagger.json` unless you copy `api/tumdev/campus_backend.swagger.json` to `swagger/swagger.json` (the Dockerfile does this).
- `docker-compose.yml` is the **production** deployment (Cloud SQL proxy + Traefik). Use `docker-compose.local.yml` for local work.
- `client/` is a separate Go module. `go test ./...` in `server/` does not cover it.
- Several legacy tables and columns use camelCase or odd names (`lastRun`, `barrierFree_moreInfo`). Match the existing column names exactly and don't "fix" them, because production depends on them.
