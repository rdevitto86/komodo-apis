#!/bin/bash
# filepath: /Users/rad/komodo-apis/komodo-address-api/deploy/deploy_docker_qa.sh

set -e

# Set environment
ENV=qa
COMPOSE_BASE="build/docker-compose.base.yml"
COMPOSE_QA="build/docker-compose.qa.yml"
ENV_FILE="config/.env.qa"

echo "Deploying komodo-address-api to QA Docker cluster..."

# Build the Docker image with QA settings
docker build -f build/Dockerfile -t komodo-address-api:${ENV} --build-arg API_ENV=${ENV} .

# Start the QA stack using Docker Compose overlays
docker compose -f ${COMPOSE_BASE} -f ${COMPOSE_QA} --env-file ${ENV_FILE} up -d --remove-orphans

echo "QA deployment complete."
docker compose -f ${COMPOSE_BASE} -f ${COMPOSE_QA} ps
