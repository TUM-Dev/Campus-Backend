# Campus-Backend

This repository holds the following components:
* `server/api` - the proto API definition in [`tumdev/campus_backend.proto`](server/api/tumdev/campus_backend.proto).
* `client` - example client for how to connect to the backend.
* `server` - the actual server implementation serving both REST at [api.tum.app](https://api.tum.app)
  and gRPC endpoints at [api-grpc.tum.app](https://api-grpc.tum.app).

The API is publicly available for anyone, but most notably, it's the main backend system for the TUM Campus Apps (Android, iOS, and Windows).

## Running the Server
```bash
docker compose -f docker-compose.local.yml up -d
```
The docker compose will start the server and a mysql instance (=> without routing/certificates to worry about)
The server will be available at `localhost:50051` and the mysql instance at `localhost:3306`.
Additionally, docker creates the volume `mysql-data` to persist the data of the mysql instance.

To rebuild and restart the server automatically whenever you change code in `server/`, run it with [watch](https://docs.docker.com/compose/how-tos/file-watch/) instead:
```bash
docker compose -f docker-compose.local.yml up --watch
```
Changes to tests, test data and markdown files are ignored.

### Environment Variables
The following environment variables need to be set for the server to work properly:
* [REQUIRED] `DB_NAME`: The name of the database to use.
* [REQUIRED] `DB_USER_PASSWORD`: The password of the user.
* [OPTIONAL] `DB_USER_NAME`: Name of the user to connect as. Defaults to `root`.
* [OPTIONAL] `DB_PORT`: The port of the database server. Defaults to `3306`.
* [OPTIONAL] `SENTRY_DSN`: The Sentry [Data Source Name](https://sentry-docs-git-patch-1.sentry.dev/product/sentry-basics/dsn-explainer/) for reporting issues and crashes.
* [OPTIONAL] `OMDB_API_KEY`: The key to get more information for tu-film movies from [omdbapi](https://omdbapi.com/). See [omdbapi](https://omdbapi.com/apikey.aspx) for a key.

## Running Tests

```bash
cd server && go test ./...
```

Tests that need a database start a single MySQL container via [testcontainers](https://golang.testcontainers.org/) per test binary (Docker required) and give every test its own freshly migrated database.
To reuse an already running MySQL server instead, set `TEST_MYSQL_DSN`, e.g. `TEST_MYSQL_DSN='root:password@tcp(localhost:3306)/'`.

## Metrics
Our service uses prometheus to collect metrics to display in grafana.
To see the metrics we aggregate, head over to `http://localhost:50051/metrics`

## iOS Push Notifications Service
The iOS Push Notifications Service can be used to send push notifications to iOS devices.

## Visual Studio Code
There are already predefined Visual Studio Code launch tasks for debugging the client and server.
Take a look at the [`launch.json`](.vscode/launch.json) file for more details.


Please be respectful with its usage!

## pre-commit

To ensure that that common pitfalls which can be automated are not done, we recommend you to install `pre-commit`.
You can do so via

```bash
python -m venv venv
source venv/bin/activate
pip install pre-commit
pre-commit install
```

Certain `pre-commit` hooks will now be run on every commit where you change specific files.
If you want to run all files instead, run `pre-commit run -a`
