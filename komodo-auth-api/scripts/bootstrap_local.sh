#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -u  # Treat unset variables as an error

# Variables
APP_NAME="komodo-address-api"
DOCKERFILE_PATH="./build/Dockerfile"
LOCAL_TAG="local/$APP_NAME:latest"
CONTAINER_NAME="komodo-address-api"
PORT=7010

# Functions
function build_image() {
  echo "Building Docker image..."
  docker build -t "$LOCAL_TAG" -f "$DOCKERFILE_PATH" .
  echo "Docker image built: $LOCAL_TAG"
}

function run_container() {
  echo "Running Docker container..."
  docker run -d \
    -p "$PORT:$PORT" \
    --name "$CONTAINER_NAME" \
    "$LOCAL_TAG"
  echo "Docker container running: $CONTAINER_NAME"
}

function build_and_run() {
  build_image
  stop_container
  run_container
}

function stop_container() {
  if docker ps -q -f name="$CONTAINER_NAME" > /dev/null; then
    echo "Stopping existing container: $CONTAINER_NAME..."
    docker stop "$CONTAINER_NAME" > /dev/null || true
    docker rm "$CONTAINER_NAME" > /dev/null || true
    echo "Stopped and removed container: $CONTAINER_NAME"
  else
    echo "No running container found with name: $CONTAINER_NAME"
  fi
}

function restart_container() {
  echo "Restarting Docker container..."
  stop_container
  run_container
  echo "Docker container restarted: $CONTAINER_NAME"
}

function cleanup() {
  echo "Cleaning up unused Docker resources..."
  # Remove stopped containers
  docker container prune -f
  # Remove dangling images
  docker image prune -f
  # Remove unused networks
  docker network prune -f
  # Remove unused volumes (optional, be cautious)
  docker volume prune -f
  echo "Docker cleanup completed."
}

function usage() {
  echo "Usage: $0 [build|run|bootstrap|stop|restart|cleanup]"
  echo "  build: Build the Docker image."
  echo "  run: Run the Docker container."
  echo "  bootstrap: Build the Docker image and run the container."
  echo "  stop: Stop and remove the running container."
  echo "  restart: Restart the Docker container without rebuilding the image."
  echo "  cleanup: Remove unused Docker resources."
  exit 1
}

# Main Script
if [[ $# -ne 1 ]]; then
  usage
fi

case "$1" in
  build)
    build_image
    ;;
  run)
    run_container
    ;;
  bootstrap)
    build_and_run
    ;;
  stop)
    stop_container
    ;;
  restart)
    restart_container
    ;;
  cleanup)
    cleanup
    ;;
  *)
    usage
    ;;
esac