.PHONY: formatting install build run clean dev_infra test tables help docker-up docker-down docker-logs docker-exec docker-rebuild

# Setup code formatting
formatting:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/segmentio/golines@latest
	brew install pre-commit
	pre-commit install
	brew install make

# Install required libraries
install:
	go mod tidy
	go mod download

# Build packages
build:
	make install
	go build ./...

# Run the application
run:
	make build
	go run cmd/api/main.go

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Clean built cache
clean:
	go clean -modcache

# Start infra
dev_infra:
	docker compose -f infra/docker.compose.dev.yml build
	docker compose -f infra/docker.compose.dev.yml up -d

# Create database tables
tables:
	@echo "Creating database tables..."
	docker exec -i treasure_dev_postgres psql -U dev_pg_user -d treasure_dev_db < db-scripts/players.sql
	docker exec -i treasure_dev_postgres psql -U dev_pg_user -d treasure_dev_db < db-scripts/treasure_chest.sql
	docker exec -i treasure_dev_postgres psql -U dev_pg_user -d treasure_dev_db < db-scripts/treasure_owner.sql
	docker exec -i treasure_dev_postgres psql -U dev_pg_user -d treasure_dev_db < db-scripts/treasure_explorer.sql
	@echo "Database tables created successfully!"

# Docker commands - Full stack with app
docker-up:
	docker compose -f infra/docker.compose.dev.yml up -d --build
	@echo "🐳 Docker environment is up! App running on http://localhost:8008"

docker-down:
	docker compose -f infra/docker.compose.dev.yml down

docker-logs:
	docker compose -f infra/docker.compose.dev.yml logs -f app

docker-exec:
	@echo "🐚 Entering app container with Go 1.24 available..."
	docker exec -it treasure_dev_app /bin/sh

docker-rebuild:
	docker compose -f infra/docker.compose.dev.yml up -d --build --force-recreate app

# Show help
help:
	@echo "Local Development:"
	@echo "  formatting    - Setup code formatting tools"
	@echo "  install       - Install all Go dependencies"
	@echo "  build         - Build all Go packages"
	@echo "  run           - Build and run the application"
	@echo "  test          - Run all tests with race detection and coverage"
	@echo "  clean         - Clean Go module cache"
	@echo ""
	@echo "Database:"
	@echo "  dev_infra     - Start database only (postgres)"
	@echo "  tables        - Create all database tables from SQL scripts"
	@echo ""
	@echo "Docker Development (with Go 1.24):"
	@echo "  docker-up     - Start full stack (postgres + app) in Docker"
	@echo "  docker-down   - Stop all Docker containers"
	@echo "  docker-logs   - View app container logs"
	@echo "  docker-exec   - Enter app container shell (Go commands available)"
	@echo "  docker-rebuild- Rebuild and restart app container"

