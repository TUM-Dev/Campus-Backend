# Campus-Backend

This repository holds the following components:
* `api` - the proto API definition in the `CampusService.proto`.
* `client` - example client for how to connect to the backend.
* `server` - the actual server implementation serving both REST at [api.tum.app](https://api.tum.app)
  and gRPC endpoints at [api-grpc.tum.app](https://api-grpc.tum.app).

The API is publicly available for anyone, but most notably, it's the main backend system for the TUM Campus Apps (Android, iOS, and Windows).

## Running the Server (without Docker)

### Installing Requirements

The backend uses MySQL as its backend for storing data.
While it is possible to install [mysql](https://mysql.com/) natively (instructions are on their website), we recommend the following:

```bash
docker run 
```

### Setting up the DB

To setup the DB, connect to your DB server or use `mysql` to connect to it locally.

```sql
DROP DATABASE campus_backend; -- Drop an eventually existing old DB
CREATE DATABASE campus_backend; -- Create a new DB for the backend

CREATE USER 'gorm'@'localhost' IDENTIFIED BY 'gorm'; -- Create a new user called `gorm`.
GRANT ALL PRIVILEGES ON campus_backend.* TO 'gorm'@'localhost'; -- Garant our `gorm` user access to the `campus_backend` DB.
ALTER USER 'gorm'@'localhost' IDENTIFIED BY 'GORM_USER_PASSWORD'; -- Set a password for the `gorm` user.
FLUSH PRIVILEGES;
```

### Starting

To start the server there are environment variables, as well as command line options available for configuring the server behavior.

```bash
cd  server
export DB_DSN="Your gorm DB connection string for example: gorm:GORM_USER_PASSWORD@tcp(localhost:3306)/campus_backend"
export DB_DSN="The DB-name from above string for example: campus_backend"
go run ./main.go
```

#### Environment Variables

There are a few environment variables available:

* [REQUIRED] `DB_DSN`: The [GORM](https://gorm.io/) [DB connection string](https://gorm.io/docs/connecting_to_the_database.html#MySQL) for connecting to the MySQL DB. Example: `gorm@tcp(localhost:3306)/campus_backend`
* [REQUIRED] `DB_DSN`: The name of the database from above connection string. Example: `campus_backend`
* [OPTIONAL] `SENTRY_DSN`: The Sentry [Data Source Name](https://sentry-docs-git-patch-1.sentry.dev/product/sentry-basics/dsn-explainer/) for reporting issues and crashes.

## Running the Server (Docker)
```bash
docker compose -f docker-compose.local.yml up -d
```
The docker compose will start the server and a mariadb instance (=> without the grpc-web layer and without routing/certificates to worry about)
The server will be available at `localhost:50051` and the mariadb instance at `localhost:3306`.
Additionally, docker creates the volume `campus-db-data` to persist the data of the mariadb instances.

### Environment Variables
The following environment variables need to be set for the server to work properly:
* [REQUIRED] `DB_NAME`: The name of the database to use.
* [REQUIRED] `DB_USER_PASSWORD`: The password of the user.
* [OPTIONAL] `DB_USER_NAME`: Name of the user to connect as. Defaults to `root`.
* [OPTIONAL] `DB_PORT`: The port of the database server. Defaults to `3306`.
* [OPTIONAL] `SENTRY_DSN`: The Sentry [Data Source Name](https://sentry-docs-git-patch-1.sentry.dev/product/sentry-basics/dsn-explainer/) for reporting issues and crashes.
* [OPTIONAL] `OMDB_API_KEY`: The key to get more information for tu-film movies from [omdbapi](https://omdbapi.com/). See [omdbapi](https://omdbapi.com/apikey.aspx) for a key.

## Metrics
Our service uses prometheus to collect metrics to display in grafana.
To see the metrics we aggregate, head over to `http://localhost:50051/metrics`

## iOS Push Notifications Service
The iOS Push Notifications Service can be used to send push notifications to iOS devices.

## Visual Studio Code
There are already predefined Visual Studio Code launch tasks for debugging the client and server.
Take a look at the [`lauch.json`](.vscode/launch.json) file for more details.


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
