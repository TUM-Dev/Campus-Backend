<div align="center">

# 🎓 Campus-Backend

**The API behind the TUM Campus Apps. One Go binary serves gRPC and REST, and also scrapes the campus for data.**

[![CI](https://github.com/TUM-Dev/Campus-Backend/actions/workflows/main.yml/badge.svg)](https://github.com/TUM-Dev/Campus-Backend/actions/workflows/main.yml)
[![Lint](https://github.com/TUM-Dev/Campus-Backend/actions/workflows/lint.yml/badge.svg)](https://github.com/TUM-Dev/Campus-Backend/actions/workflows/lint.yml)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![API](https://img.shields.io/badge/API-api.tum.app-0065BD)](https://api.tum.app)

[**Try the API**](https://api.tum.app) · [**Proto definition**](server/api/tumdev/campus_backend.proto) · [**Quick start**](#-quick-start) · [**Contributing**](#-contributing)

</div>

---

Every time a student at the Technical University of Munich opens the TUM Campus App to check today's menu, rate a dish, find out how crowded the Mensa is, or see what [TU Film](https://www.tu-film.de) is showing this week, the request ends up here.

The API is **public**. If you want to build something for campus, you can use it too. Please use it responsibly.

## ✨ What it does

| | Feature | Where the data comes from |
|---|---|---|
| 🍝 | **Canteens and dishes**: menus, opening info, ratings with tags | [eat-api](https://tum-dev.github.io/eat-api/) + ratings from app users |
| 👥 | **Live canteen head count**: how full is the Mensa right now? | Estimated from LRZ Wi-Fi access point stats |
| 📰 | **News and alerts** from TUM sources | Scraped on a schedule by cron jobs |
| 🎬 | **Movies** from TU Film, with posters and metadata | tu-film.de + [OMDb](https://omdbapi.com/) |
| 🏛️ | **Student clubs** directory | Scraped from the student council site |
| 💬 | **Feedback** with image uploads (gRPC streaming), sent out by email | App users |
| 📱 | **Update notes** and device registration | App users |

## 🏗️ How it works

```
 Campus apps ──gRPC──┐
                     ▼
               ┌── :50051 (cmux) ──┐
               │                   │
         application/grpc        HTTP/1
               │                   ├── /v1/*      → grpc-gateway ─┐
               ▼                   ├── /swagger/  → Swagger UI    │
          gRPC server ◀────────────┼───────────────────────────────┘
               │                   ├── /metrics   → Prometheus
               ▼                   ├── /files/    → /Storage
         CampusServer              └── /health
               │
               ▼
             MySQL ◀── CronService (scrapers, polls `crontab` every 30 s)
```

In production, REST is served at [api.tum.app](https://api.tum.app) and gRPC at `api-grpc.tum.app:443`.

- **One proto, two protocols.** [`campus_backend.proto`](server/api/tumdev/campus_backend.proto) is the single source of truth. gRPC stubs, the REST gateway and the OpenAPI spec are all generated from it.
- **One port.** [`cmux`](https://github.com/soheilhy/cmux) looks at each connection and sends it to gRPC or HTTP.
- **Built-in cron.** Scrapers and downloaders run inside the same binary, scheduled from the `crontab` table.
- **Observability.** Prometheus metrics, Sentry error reporting and structured logs via logrus.

## 🚀 Quick start

You need [Docker](https://docs.docker.com/get-docker/). That's it.

```bash
git clone https://github.com/TUM-Dev/Campus-Backend.git
cd Campus-Backend
docker compose -f docker-compose.local.yml up --watch
```

This starts MySQL and the backend at `localhost:50051`. With `--watch`, the server is rebuilt and restarted every time you save a file in `server/`. Changes to tests, test data and Markdown files are ignored. Use `up -d` instead if you just want it running in the background.

Now try it:

```bash
# REST
curl localhost:50051/v1/canteen/allCanteens
curl localhost:50051/v1/news/sources

# gRPC (server reflection is enabled)
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 api.Campus/ListCanteens
```

| URL | What you get |
|---|---|
| `localhost:50051/v1/…` | REST API |
| `localhost:50051/swagger/` | Swagger UI (see [note](#-gotchas)) |
| `localhost:50051/metrics` | Prometheus metrics |
| `localhost:50051/health` | Health check |
| `localhost:3306` | MySQL (data persists in the `mysql-data` volume) |

### Running without Docker

```bash
cd server
DB_DSN='root:secret_root_password@tcp(localhost:3306)/campus_db?charset=utf8mb4&parseTime=True&loc=Local' go run .
```

## 🛠️ Development

All Go commands run from `server/`.

```bash
go build ./...
go test ./...            # spins up a MySQL testcontainer, so Docker is required
bash api/generate.bash   # regenerate gRPC, gateway and Swagger code from the proto
```

Tests that need a database start one MySQL container per test binary via [testcontainers](https://golang.testcontainers.org/) and give each test its own database with the production schema. To reuse a MySQL server that's already running:

```bash
TEST_MYSQL_DSN='root:secret_root_password@tcp(localhost:3306)/' go test ./...
```

VS Code users get debug configs for free in [`.vscode/launch.json`](.vscode/launch.json).

### 🗺️ Repo map

```
server/
├── main.go                       # entry point: cmux, gRPC, gateway, HTTP routes
├── api/tumdev/                   # campus_backend.proto + generated code
├── backend/                      # one file per domain: cafeteria.go, news.go, movie.go, …
│   ├── cron/                     # scheduled scrapers and downloaders
│   └── migration/                # schema migrations (legacy Go + Atlas-generated sql/)
├── model/                        # GORM models, one per table
└── utils/                        # DB setup, test containers, telemetry, error mapping
client/                           # example gRPC clients (separate Go module)
```

### Recipes

<details>
<summary><b>➕ Add a new endpoint</b></summary>

1. Add an `rpc` with a `google.api.http` option, plus its request and reply messages, to [`campus_backend.proto`](server/api/tumdev/campus_backend.proto).
2. Run `bash api/generate.bash` and commit the generated files.
3. Implement the method on `*CampusServer` in `server/backend/<domain>.go`. Start it with `s.checkDevice(ctx)` and return `status.Error(codes.X, …)` errors.

⚠️ Released apps depend on this API. Never renumber or reuse proto field numbers. A proto change on `main` automatically opens a PR in [Campus-Flutter](https://github.com/TUM-Dev/Campus-Flutter).

</details>

<details>
<summary><b>🗄️ Change the database schema</b></summary>

Migrations are generated with [Atlas](https://atlasgo.io) from the GORM models and applied by the server on startup.

1. Edit the model in `server/model/`. If it's a new model, add it to `migration.Models()`.
2. In `server/`, run `atlas migrate diff --env gorm <short_name>` ([install Atlas](https://atlasgo.io/getting-started#installation); Docker is required). It writes `backend/migration/sql/<version>_<short_name>.sql` and updates `atlas.sum`.
3. Review the generated SQL and commit it together with `atlas.sum`.

For data changes, or for views and triggers (which Atlas' free edition doesn't manage), edit the generated file or create one with `atlas migrate new --env gorm <short_name>`, then run `atlas migrate hash --env gorm`.

The desired schema is the models plus [`legacy.sql`](server/tools/atlasloader/legacy.sql), which declares production tables without a model. The Go migrations in `backend/migration/` are frozen, including the post-baseline registry; a regression test rejects new Go migration IDs. Keep existing IDs for compatibility, but never duplicate an Atlas schema change in Go. CI fails, and comments the missing SQL on your PR, if the models change without a migration.

</details>

<details>
<summary><b>⏰ Add a cron job</b></summary>

1. Add a type constant in `backend/cron/cronjobs.go` and wire it into the query and the `switch`.
2. Add the value to the `type` enum in `model/crontab.go`.
3. Generate the migration with `atlas migrate diff --env gorm <short_name>`, append the `INSERT` for the `crontab` row, then run `atlas migrate hash --env gorm`.

</details>

## ⚙️ Configuration

| Variable | Required | Purpose |
|---|:---:|---|
| `DB_DSN` | ✅ | MySQL DSN. Docker Compose builds it from `DB_NAME`, `DB_USER_NAME`, `DB_USER_PASSWORD` and `DB_PORT` in [`.env`](.env). |
| `ENVIRONMENT` | | `dev` or `prod`. Changes telemetry behavior. |
| `SENTRY_DSN` | | Error reporting via [Sentry](https://sentry.io) |
| `OMDB_API_KEY` | | Movie metadata from [OMDb](https://omdbapi.com/apikey.aspx) |
| `SMTP_URL`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` | | Sending feedback emails |

The committed `.env` contains only local dev defaults. Put real secrets in `.env.local`, which is gitignored.

## 🤝 Contributing

Contributions are very welcome, whether that's a bug fix, a new scraper or a whole new feature for the apps.

1. Install [pre-commit](https://pre-commit.com) so formatting, `goimports`, proto regeneration and `go mod tidy` run automatically:
   ```bash
   pip install pre-commit && pre-commit install
   pre-commit run -a   # optional: check every file once
   ```
2. Make your change, add tests, and open a PR. CI runs the tests, golangci-lint and the migration checks.

Working with an AI coding agent? [`AGENTS.md`](AGENTS.md) has the full conventions.

### 🪤 Gotchas

- `/swagger/` locally has no spec until you copy `server/api/tumdev/campus_backend.swagger.json` to `server/swagger/swagger.json`. The Dockerfile does this for you.
- `docker-compose.yml` is the **production** deployment. Use `docker-compose.local.yml` for local work.
- Some legacy columns have odd names (`lastRun`, `barrierFree_moreInfo`). Production relies on them, so leave them as they are.

## 📜 License

[GPL-3.0](LICENSE). Built with ❤️ by [TUM-Dev](https://github.com/TUM-Dev) for students at TUM.
