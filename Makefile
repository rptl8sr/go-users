.PHONY: help generate test-local test-all dev dev-down prod prod-down

help:
	@echo "Available commands:"
	@echo "  make generate     - generate code"
	@echo "  make test-local   - run unit tests locally"
	@echo "  make test-all     - run all tests with docker-compose"
	@echo "  make dev          - run in development mode with docker-compose"
	@echo "  make dev-down     - stop development environment"
	@echo "  make prod         - run in production mode with docker-compose"
	@echo "  make prod-down    - stop production environment"

generate:
	go generate ./...

test-local:
	APP_MODE=test go test -v ./internal/...

test-all:
	go vet ./...
	APP_MODE=test docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.test.yml --env-file .env.test up -d postgres --remove-orphans
	APP_MODE=test docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.test.yml --env-file .env.test run --rm migrations
	APP_MODE=test docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.test.yml --env-file .env.test up --build test
	APP_MODE=test docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.test.yml --env-file .env.test up -d app
	sleep 5
	go test -v ./tests/integration/...
	APP_MODE=test docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.test.yml --env-file .env.test down -v

dev:
	APP_MODE=dev docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.dev.yml --env-file .env.dev up --build

dev-down:
	APP_MODE=dev docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.dev.yml --env-file .env.dev down -v

prod:
	APP_MODE=prod docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.prod.yml --env-file .env.prod up --build -d

prod-down:
	APP_MODE=prod docker-compose -f docker/compose/docker-compose.base.yml -f docker/compose/docker-compose.prod.yml --env-file .env.prod down