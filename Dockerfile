# ======================
# Builder Stage
# ======================
FROM golang:1.23-alpine AS builder

WORKDIR /go/src/app

# Install required packages for CGO
RUN apk add --no-cache git gcc musl-dev

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire source code
COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/server

# ======================
# Runtime Stage
# ======================
FROM alpine:3.19

# Install SQLite CLI (optional: remove if not needed)
# RUN apk add --no-cache sqlite

# Copy binary from builder
COPY --from=builder /go/bin/app /app

ENV PORT=8080
EXPOSE 8080

# Default entrypoint
CMD ["/app"]
