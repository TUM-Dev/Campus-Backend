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
In the following, we provide instructions for installing [MariaDB](https://mariadb.org/) as the DB server of choice.

#### Fedora

```bash
sudo dnf install mariadb-server

# Start the MariaDB server
sudo systemctl start mariadb

# Optional: Enable autostart
sudo systemctl enable mariadb
```

More details are available here: https://docs.fedoraproject.org/en-US/quick-docs/installing-mysql-mariadb/

#### Debian/Ubuntu

```bash
sudo apt install mariadb-server

# Start the MariaDB server
sudo systemctl start mariadb

# Optional: Enable autostart
sudo systemctl enable mariadb
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
go run ./main.go
```

#### Environment Variables

There are a few environment variables available:

* [REQUIRED] `DB_DSN`: The [GORM](https://gorm.io/) [DB connection string](https://gorm.io/docs/connecting_to_the_database.html#MySQL) for connecting to the MySQL DB. Example: `gorm@tcp(localhost:3306)/campus_backend`
* [OPTIONAL] `SENTRY_DSN`: The Sentry [Data Source Name](https://sentry-docs-git-patch-1.sentry.dev/product/sentry-basics/dsn-explainer/) for reporting issues and crashes.

## Running the Server (Docker)
```bash
docker compose up -d
```
The docker compose will start the server and a mariadb instance.
The server will be available at `localhost:50051` and the mariadb instance at `localhost:3306`.
Additionally, docker creates the volume `campus-db-data` to persist the data of the mariadb instances.

### Environment Variables
The following environment variables need to be set for the server to work properly:
* [REQUIRED] `DB_NAME`: The name of the database to use.
* [REQUIRED] `DB_ROOT_PASSWORD`: The password of the root user.
* [OPTIONAL] `DB_PORT`: The port of the database server. Defaults to `3306`.
* [OPTIONAL] `SENTRY_DSN`: The Sentry [Data Source Name](https://sentry-docs-git-patch-1.sentry.dev/product/sentry-basics/dsn-explainer/) for reporting issues and crashes.
* **[iOS Push Notification Service [OPTIONAL]](#ios-push-notifications-service)**:
  * [REQUIRED] `APNS_KEY_ID`: The key ID of the APNs key => APNs Key needs to be downloaded from the Apple Developer Portal the name of the file also contains the key ID.
  * [REQUIRED] `APNS_TEAM_ID`: The team ID of the iOS app can be found in AppStoreConnect.
  * [REQUIRED] `APNS_P8_FILE_PATH`: The path to the APNs key file (e.g. `/secrets/AuthKey_XXXX.p8`) in the docker container. The file itself needs to exist in the same directory as the `docker-compose.yml` file and called `apns_auth_key.p8`.
  * [REQUIRED] `CAMPUS_API_TOKEN`: A token used to authenticate with TUMonline (used for example for the grades)

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
