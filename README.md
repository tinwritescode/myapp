# MyApp

A Go web application built with Gin framework and Swagger documentation.

## Prerequisites

- Go 1.25.1 or higher
- PostgreSQL (for database)
- Swag CLI for generating Swagger docs

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

## Available Commands

This project uses a Makefile for convenient command execution:

### Development
```bash
make run        # Run the application
make build      # Build the application binary
make swagger    # Generate Swagger documentation
```

### Code Quality
```bash
make fmt        # Format Go code
make lint       # Run linter
make test       # Run tests
```

### Dependencies
```bash
make deps       # Download and tidy dependencies
```

### Cleanup
```bash
make clean      # Remove build artifacts and generated docs
```

## API Documentation

After running `make swagger`, you can access the Swagger UI at:
- http://localhost:8080/swagger/index.html

## Project Structure

```
├── cmd/           # Application entry points
├── internal/      # Private application code
│   ├── config/    # Configuration
│   ├── database/  # Database connection and migrations
│   ├── handlers/  # HTTP handlers
│   ├── models/    # Data models
│   ├── routes/    # Route definitions
│   └── service/   # Business logic
├── docs/          # Generated Swagger documentation
├── migrations/    # Database migrations
├── pkg/           # Public library code
└── tests/         # Test files
```

## Getting Started

1. Start the application:
   ```bash
   make run
   ```

2. Generate API documentation:
   ```bash
   make swagger
   ```

3. Visit http://localhost:8080/ping to test the API
4. Visit http://localhost:8080/swagger/index.html for API documentation
