.PHONY: build run clean test lint deps dev

# Build the game
build:
	go build -o bin/game cmd/game/main.go

# Run the game
run:
	go run cmd/game/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Development mode (with hot reload)
dev:
	air
