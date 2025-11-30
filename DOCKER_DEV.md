<!-- @format -->

# Docker Development Guide

This guide explains how to develop with Docker where Go 1.24 is available inside the container.

## Quick Start

```bash
# Start the full stack (PostgreSQL + App with Go 1.24)
make docker-up

# View logs
make docker-logs

# Enter the container shell (Go commands available here!)
make docker-exec

# Stop everything
make docker-down
```

## Using Go Commands Inside Docker

Once you run `make docker-exec`, you'll be inside the container with Go 1.24 installed:

```bash
# Inside the container shell
go version              # Check Go version (should show 1.24)
go build ./...          # Build all packages
go test ./...           # Run tests
go run cmd/api/main.go  # Run the app
go mod tidy             # Tidy dependencies
make run                # Make commands work too!
```

## Development Workflow

### Option 1: Local Development (Current)

```bash
make run  # Runs locally, uses local Go installation
```

### Option 2: Docker Development (New)

```bash
# Start everything in Docker
make docker-up

# Your code changes are automatically reflected (volume mounted)
# The app will restart on code changes if you add a hot-reload tool

# To manually restart after changes:
make docker-rebuild
```

### Option 3: Hybrid (Database in Docker, App Local)

```bash
# Start only the database
make dev_infra

# Run app locally (pointing to Docker database)
make run
```

## File Structure

- `Dockerfile` - Production build (multi-stage, minimal final image)
- `Dockerfile.dev` - Development build (includes Go 1.24 toolchain)
- `infra/docker.compose.dev.yml` - Development stack configuration

## Environment Variables

The docker-compose file sets these automatically:

- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=dev_pg_user`
- `DB_PASSWORD=dev_pg_password`
- `DB_NAME=treasure_dev_db`
- `SERVER_PORT=8008`

## Tips

1. **Volume Mounting**: Your entire project is mounted to `/app` in the container, so code changes are immediately available
2. **Go Module Cache**: Go modules are cached in a Docker volume for faster builds
3. **Database Setup**: After `docker-up`, run `make tables` to create database tables
4. **Hot Reload**: Consider adding [Air](https://github.com/cosmtrek/air) for automatic reloading on file changes

## Adding Hot Reload (Optional)

1. Install Air in the Dockerfile.dev:

```dockerfile
RUN go install github.com/cosmtrek/air@latest
```

2. Create `.air.toml` in the project root

3. Update docker-compose command:

```yaml
command: ["air", "-c", ".air.toml"]
```

## Troubleshooting

**Container won't start?**

```bash
docker compose -f infra/docker.compose.dev.yml logs app
```

**Database connection issues?**

```bash
# Check if postgres is ready
docker exec treasure_dev_postgres pg_isready -U dev_pg_user
```

**Need to rebuild from scratch?**

```bash
docker compose -f infra/docker.compose.dev.yml down -v
make docker-up
make tables
```
