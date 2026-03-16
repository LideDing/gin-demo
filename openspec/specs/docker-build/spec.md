## Requirements

### Requirement: Multi-stage Dockerfile
The project SHALL provide a `Dockerfile` at the repository root that builds the application using a multi-stage build: a builder stage compiles the Go binary and a runtime stage packages it into a minimal Alpine-based image.

#### Scenario: Builder stage compiles binary
- **WHEN** `docker build` is executed
- **THEN** the builder stage SHALL use `golang:1.25-alpine`, copy `go.mod`/`go.sum`, run `go mod download`, copy source, and produce a statically-linked binary named `gin-demo`

#### Scenario: Runtime stage is minimal
- **WHEN** the final image is produced
- **THEN** it SHALL be based on `alpine:3.21`, include `ca-certificates`, expose port `8080`, and contain only the compiled `gin-demo` binary

#### Scenario: Image runs the application
- **WHEN** a container is started from the image without arguments
- **THEN** the `gin-demo` binary SHALL be executed as the container entrypoint

### Requirement: .dockerignore excludes non-essential files
The project SHALL provide a `.dockerignore` file that prevents source-code artifacts, documentation, and local config files from being included in the Docker build context.

#### Scenario: Build context is lean
- **WHEN** `docker build` is invoked
- **THEN** the `.dockerignore` SHALL exclude `.git`, `openspec/`, `*.md`, `.env`, `.env.*`, and the local binary `gin-demo`
