# Default to Docker-first workflow
.DEFAULT_GOAL := bootstrap

# === Service metadata ===
APP := komodo-auth-user-api

# Map environments to compose files
# local, dev, staging -> nonprod compose file
# prod -> prod compose file
ENV ?= local
ENV := $(strip $(ENV))
TAG ?= $(ENV)

# Paths
BUILD_DIR := build

# Determine which compose file to use
ifeq ($(ENV),prod)
	COMPOSE_FILE := $(BUILD_DIR)/docker-compose.prod.yaml
else ifeq ($(ENV),staging)
	COMPOSE_FILE := $(BUILD_DIR)/docker-compose.staging.yaml
else ifeq ($(ENV),dev)
	COMPOSE_FILE := $(BUILD_DIR)/docker-compose.dev.yaml
else ifeq ($(ENV),local)
	COMPOSE_FILE := $(BUILD_DIR)/docker-compose.dev.yaml
else
	$(error Invalid environment: $(ENV). Valid options are: local, dev, staging, prod)
endif

.PHONY: help build run bootstrap stop restart clean test_all test_unit test_int test_e2e fmt

help:
	@echo "Targets:"
	@echo "  build             Build Docker image for ENV ($(ENV))"
	@echo "  run               Start container via Docker Compose"
	@echo "  bootstrap         Build + run in one command"
	@echo "  stop              Stop and remove container stack"
	@echo "  restart           Restart container stack"
	@echo "  clean             Prune Docker and remove local build artifacts"
	@echo "  test_all          Run all tests"
	@echo "  test_unit         Run unit tests"
	@echo "  test_int          Run integration tests"
	@echo "  test_e2e          Run end-to-end tests"
	@echo "  fmt               Run go fmt and go vet"
	@echo ""
	@echo "Supported Environments:"
	@echo "  local    - Local development (no AWS, uses .env.local)"
	@echo "  dev      - Development (AWS enabled, uses .env.dev)"
	@echo "  staging  - Staging (AWS enabled, uses .env.staging)"
	@echo "  prod     - Production (AWS enabled, uses .env.prod)"
	@echo ""
	@echo "Examples:"
	@echo "  make bootstrap               # Default: local"
	@echo "  make bootstrap ENV=local     # Explicit local"
	@echo "  make bootstrap ENV=dev       # Dev environment"
	@echo "  make bootstrap ENV=staging   # Staging environment"
	@echo "  make bootstrap ENV=prod      # Production environment"
	@echo "  make stop ENV=dev"
	@echo "  make restart ENV=staging"

# Build Docker image
build:
	@echo "Building $(APP):$(TAG) for ENV=$(ENV)"
	@docker build \
		-f $(BUILD_DIR)/Dockerfile \
		-t $(APP):$(TAG) \
		--build-arg ENV=$(ENV) \
		..

# Run container with docker-compose
run:
	@echo "Starting $(APP) for ENV=$(ENV) using $(COMPOSE_FILE)"
	@ENV=$(ENV) docker compose -f $(COMPOSE_FILE) up -d

# Stop container
stop:
	@echo "Stopping $(APP) for ENV=$(ENV)"
	ENV=$(ENV) docker compose -f $(COMPOSE_FILE) down --remove-orphans

# Bootstrap: build + run
bootstrap:
	$(MAKE) build ENV=$(ENV)
	$(MAKE) run ENV=$(ENV)

# Restart container
restart:
	$(MAKE) stop ENV=$(ENV)
	$(MAKE) run ENV=$(ENV)

# Clean up Docker resources
clean:
	docker container prune -f
	docker image prune -f
	docker network prune -f
	docker volume prune -f
	rm -rf bin

# Deployment shortcuts
deploy-local:
	$(MAKE) bootstrap ENV=local

deploy-dev:
	$(MAKE) bootstrap ENV=dev

deploy-staging:
	$(MAKE) bootstrap ENV=staging

deploy-prod:
	$(MAKE) bootstrap ENV=prod

# Test Commands
test_all:
	go test ./... -race

test_unit:
	go test -short ./...

test_int:
	go test -tags=integration ./...

test_e2e:
	go test -tags=e2e ./...

fmt:
	go fmt ./... && go vet ./...
