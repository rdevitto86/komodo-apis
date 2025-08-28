# Default to Docker-first workflow: `make` -> build image + up (dev)
.DEFAULT_GOAL := bootstrap

# === Service metadata ===
APP := komodo-auth-api
ENV ?= dev
ENV := $(strip $(ENV))
TAG ?= $(ENV)

# Paths
BUILD_DIR := build
CONFIG_ENV := config/.env.$(ENV)

# Compose (base + overlay)
BASE_COMPOSE := $(BUILD_DIR)/docker-compose.base.yml
ENV_COMPOSE := $(BUILD_DIR)/docker-compose.$(ENV).yml

.PHONY: local build run bootstrap stop restart cleanup test_all test_unit test_int test_e2e fmt

help:
	@echo "Targets:"
	@echo "  build             Build Docker image for ENV ($(ENV))"
	@echo "  run               Build image and start container via Docker Compose"
	@echo "  local             Run app locally (no Docker, loads $(CONFIG_ENV))"
	@echo "  stop              Stop and remove container stack"
	@echo "  restart           Restart container stack"
	@echo "  cleanup           Prune Docker and remove local build artifacts"
	@echo "  test_all          Run all tests"
	@echo "  test_unit         Run unit tests"
	@echo "  test_int          Run integration tests"
	@echo "  test_e2e          Run end-to-end tests"
	@echo "  fmt               Run go fmt and go vet"
	@echo ""
	@echo "Examples:"
	@echo "  make build ENV=prod"
	@echo "  make run ENV=dev"
	@echo "  make local ENV=dev"
	@echo "  make stop ENV=prod"
	@echo "  make restart ENV=dev"
	@echo "  make cleanup"
	@echo "  make test_all"

local:
	@export $$(grep -v '^\s*#' $(CONFIG_ENV) | grep -v '^\s*$$' | xargs) ; \
	go run ./cmd/$(APP)

build:
	docker build -f $(BUILD_DIR)/Dockerfile -t $(APP):$(TAG) --build-arg API_ENV=$(ENV) .

run:
	docker compose -f $(BASE_COMPOSE) -f $(ENV_COMPOSE) up -d

stop:
	docker compose -f $(BASE_COMPOSE) -f $(ENV_COMPOSE) down -v

bootstrap:
	$(MAKE) build
	$(MAKE) run

restart:
	$(MAKE) stop ENV=$(ENV)
	$(MAKE) run ENV=$(ENV)

cleanup:
	docker container prune -f
	docker image prune -f
	docker network prune -f
	docker volume prune -f
	rm -rf bin

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
