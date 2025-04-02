# Go Users API

A simple REST API service for user management built with Go.

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- Make

## Environment Setup

The project requires different environment files for different modes:

- `.env.test` - for running tests
- `.env.dev` - for development
- `.env.prod` - for production

Each environment file should contain necessary configuration variables. You can find an example environment file at `example.env`. Copy it to create your environment files and adjust the values as needed.

## Available Commands

```bash
# Generate code
make generate

# Run unit tests locally
make test-local

# Run all tests with docker-compose
make test-all   # Will be stop automatically after tests ends

# Development mode
make dev        # Start development containers
make dev-down   # Stop development containers

# Production mode
make prod       # Start production containers
make prod-down  # Stop production containers
```

## Project Structure

```
.
├── cmd/                # Application entry points
├── internal/           # Internal packages
│   ├── api/            # Generated API code and handlers
│   ├── config/         # Configuration
│   ├── database/       # Database models and migrations
│   ├── ownErrors/      # Custom error types
│   └── router/         # Router setup
├── openapi/            # OpenAPI specification
├── tests/              # Test files
│   └── integration/    # Integration tests
├── docker/             # Docker configuration
│   ├── build/          # Dockerfileы
│   └── compose/        # Docker Compose files
├── tools/              # Development tools and scripts for API generation
└── .env.MODE           # .env files with needed mode suffix [ test | dev | prod ]
```

## API Documentation

The API is generated using [oapi-codegen](https://github.com/deepmap/oapi-codegen) from OpenAPI specification located in `openapi/openapi.yaml`. When running in development mode, you can access the Swagger UI documentation at http://localhost:8080/swagger/.

## API Endpoints

### Create User
```bash
  curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com"
  }'
```

### Get User
```bash
  curl http://localhost:8080/api/users/1
```

### Update User
```bash
  curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Smith",
    "email": "john.smith@example.com"
  }'
```

### Health Check
```bash
  curl http://localhost:8080/api/health
``` 