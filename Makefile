# Run app locally with SQLite
run:
	go run ./cmd/server/main.go

# Run tests with coverage
test:
	go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

# Generate Swagger docs
swag:
	swag init -g cmd/server/main.go -o docs

# Build the Docker image
docker-build:
	docker build -t user-service .

# Run using Docker Compose
docker-up:
	docker-compose up --build

# Stop containers
docker-down:
	docker-compose down

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

# Run full CI checks: linting, security, swagger, formatting
ci: tidy fmt swag lint gosec test