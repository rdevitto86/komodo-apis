#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -u  # Treat unset variables as an error

# Variables
AUTH_API_IMAGE="local/komodo-auth-api:latest"
AUTH_API_CONTAINER="komodo-auth-api"
AUTH_API_PORT=7001
AUTH_API_URL="http://localhost:$AUTH_API_PORT"

ADDRESS_API_IMAGE="local/komodo-address-api:latest"
ADDRESS_API_CONTAINER="komodo-address-api"
ADDRESS_API_PORT=7010

# Functions
function start_api() {
  local image=$1
  local container=$2
  local port=$3
  local env_vars=$4

  echo "Starting $container on port $port..."

  # Stop and remove the container if it already exists
  if docker ps -q -f name="$container" > /dev/null; then
    echo "Stopping existing container: $container..."
    docker stop "$container" > /dev/null || true
    docker rm "$container" > /dev/null || true
  fi

  # Run the container
  docker run -d \
    -p "$port:$port" \
    --name "$container" \
    $env_vars \
    "$image"

  echo "$container is running on port $port."
}

function start_auth_api() {
  start_api "$AUTH_API_IMAGE" "$AUTH_API_CONTAINER" \
}

function start_address_api() {
  start_api "$ADDRESS_API_IMAGE" "$ADDRESS_API_CONTAINER" "$ADDRESS_API_PORT" \
    "-e GEOCODER=mock -e AUTH_API_URL=$AUTH_API_URL"
}

function stop_all() {
  echo "Stopping all API containers..."
  docker stop "$AUTH_API_CONTAINER" "$ADDRESS_API_CONTAINER" > /dev/null || true
  docker rm "$AUTH_API_CONTAINER" "$ADDRESS_API_CONTAINER" > /dev/null || true
  echo "All API containers stopped and removed."
}

function usage() {
  echo "Usage: $0 [start|stop|restart]"
  echo "  start: Start all APIs."
  echo "  stop: Stop all APIs."
  echo "  restart: Restart all APIs."
  exit 1
}

# Main Script
if [[ $# -ne 1 ]]; then
  usage
fi

case "$1" in
  start)
    start_auth_api
    start_address_api
    start_other_api
    ;;
  stop)
    stop_all
    ;;
  restart)
    stop_all
    start_auth_api
    start_address_api
    start_other_api
    ;;
  *)
    usage
    ;;
esac