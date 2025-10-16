# Run app locally with SQLite
run:
	go run ./cmd/server/main.go

kill-9001:
	lsof -ti:9001 | xargs kill -9 || true

# Database migration commands
migrate-up:
	go run ./cmd/migrate/main.go -command=up

migrate-down:
	go run ./cmd/migrate/main.go -command=down

migrate-status:
	go run ./cmd/migrate/main.go -command=status

# Run tests with coverage
test:
	go test ./... -coverprofile=tmp/coverage.out && go tool cover -func=tmp/coverage.out

# Run tests with detailed coverage for internal/app packages
test-coverage:
	go test -coverpkg=./internal/app/... ./internal/app -coverprofile=tmp/coverage.out && go tool cover -func=tmp/coverage.out

# Generate HTML coverage report
test-coverage-html:
	go test -coverpkg=./internal/app/... ./internal/app -coverprofile=tmp/coverage.out && go tool cover -html=tmp/coverage.out -o tmp/coverage.html

# Generate Swagger docs
swag:
	swag init -g cmd/server/main.go -o docs

# ===============================
# Docker Commands
# ===============================

# Build the Docker image
docker-build:
	docker build -t user-service .

# Run containers with build (foreground mode)
docker-up:
	docker-compose up --build

# Start containers in background
docker-start:
	docker-compose up -d

# Reload containers (rebuild and restart)
docker-reload:
	docker-compose down && docker-compose up --build -d

# Stop running containers (keep data)
docker-stop:
	docker-compose stop

# Stop and remove containers (cleanup)
docker-down:
	docker-compose down

# ===============================
# Database Commands
# ===============================

# Clean SQLite DB file (if persisted to host)
clean-db:
	rm -f test.db

# Format code
fmt:
	go fmt ./...

# Tidy modules
tidy:
	go mod tidy

# Run static analysis (gosec and golangci-lint)
lint:
	golangci-lint run ./...

gosec:
	gosec ./...

# Clean build artifacts and temporary files
clean:
	rm -rf bin/* tmp/coverage.* tmp/*.out tmp/*.db

# Run full CI checks: linting, security, swagger, formatting
ci: tidy fmt swag lint gosec test