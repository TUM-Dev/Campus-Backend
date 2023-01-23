# Campus-Backend

This repository holds the following components:
* `api` - the proto API definition in the `CampusService.proto`.
* `client` - example client for how to connect to the backend.
* `server` - the actual server implementation serving both REST at [api.tum.app](https://api.tum.app)
   and gRPC endpoints at [api-grpc.tum.app](https://api-grpc.tum.app).

The API is publicly available for anyone, but most notably, it's the main backend system for the TUM Campus Apps (Android, iOS, and Windows).

## Running the Server

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
go run ./main.go [-MensaCron 0]
```

#### Environment Variables

There are a few environment variables available:

* [REQUIRED] `DB_DSN`: The [GORM](https://gorm.io/) [DB connection string](https://gorm.io/docs/connecting_to_the_database.html#MySQL) for connecting to the MySQL DB. Example: `gorm@tcp(localhost:3306)/campus_backend`
* [OPTIONAL] `SENTRY_DSN`: The Sentry [Data Source Name](https://sentry-docs-git-patch-1.sentry.dev/product/sentry-basics/dsn-explainer/) for reporting issues and crashes.

#### Command Line Arguments

* [OPTIONAL] `-MensaCron 0`: Providing this argument deactivates the Mensa Rating cronjobs if not needed in a local setup. Be aware, this option will change in a future version ([#117](https://github.com/TUM-Dev/Campus-Backend/issues/117) and [#115](https://github.com/TUM-Dev/Campus-Backend/issues/115)).

### Visual Studio Code

There are already predefined Visual Studio Code launch tasks for debugging the client and server.
Take a look at the [`lauch.json`](.vscode/launch.json) file for more details.


Please be respectful with its usage!
