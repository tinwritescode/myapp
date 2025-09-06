.PHONY: swagger build run clean

# Generate Swagger documentation
swagger:
	swag init -g ./main.go -o ./docs

# Build the application
build:
	go build -o bin/main main.go

# Run the application
run:
	go run main.go

air:
	air

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run
