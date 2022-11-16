# Campus-Backend

This repository holds the following components:
* `api` - the proto API definition in the `CampusService.proto`.
* `client` - example clients on how to connect to the backend.
* `server` - the actual server implementation serving both REST at [api.tum.app](https://api.tum.app)
   and gRPC endpoints at at [api-grpc.tum.app](https://api-grpc.tum.app).

The API is publicly available for use by anyone, but most notably its the main backend system for the
TUM Campus Apps (Android, iOS and Windows).

### Running the Backend
Optional Commandline Parameter: `-MensaCron 0` deactivates the Mensa Rating cronjobs if not needed in a local setup. The Cronjobs are activated if the parameter is not explicitly added.



Please be respectful with its usage!
