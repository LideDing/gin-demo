## Requirements

### Requirement: docker-compose service definition
The project SHALL provide a `docker-compose.yml` at the repository root that defines a single `app` service for the gin-demo application, loading OIDC configuration from the local `.env` file.

#### Scenario: Service starts with env file
- **WHEN** `docker compose up` is run and a `.env` file exists
- **THEN** the `app` service SHALL start with all variables from `.env` injected as environment variables

#### Scenario: Port is accessible on host
- **WHEN** the `app` service is running
- **THEN** port `8080` on the host SHALL be mapped to port `8080` in the container

#### Scenario: Image is built locally
- **WHEN** `docker compose up --build` is run
- **THEN** the compose file SHALL build the image from the local `Dockerfile` instead of pulling from a registry

### Requirement: Restart policy for resilience
The `app` service SHALL have a restart policy so it recovers from unexpected exits during local development.

#### Scenario: Container restarts on failure
- **WHEN** the `app` container exits with a non-zero code
- **THEN** Docker SHALL automatically restart it (`restart: unless-stopped`)
