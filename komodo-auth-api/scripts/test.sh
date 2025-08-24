#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -u  # Treat unset variables as an error

# Variables
TEST_DIR="./tests"  # Directory containing test files
UNIT_TEST_DIR="$TEST_DIR/unit"
E2E_TEST_DIR="$TEST_DIR/e2e"
INTEGRATION_TEST_DIR="$TEST_DIR/integration"
COVERAGE_FILE="coverage.out"

# Functions
function run_unit_tests() {
    echo "Running unit tests..."
    go test "$UNIT_TEST_DIR/..." -coverprofile="$COVERAGE_FILE" -covermode=atomic
    echo "Unit tests completed."
}

function run_integration_tests() {
    echo "Running integration tests..."
    if [ -d "$INTEGRATION_TEST_DIR" ]; then
        go test "$INTEGRATION_TEST_DIR/..." -v
        echo "Integration tests completed."
    else
        echo "Integration test directory not found: $INTEGRATION_TEST_DIR"
    fi
}

function run_e2e_tests() {
    echo "Running end-to-end (e2e) tests..."
    if [ -d "$E2E_TEST_DIR" ]; then
        go test "$E2E_TEST_DIR/..." -v
        echo "E2E tests completed."
    else
        echo "E2E test directory not found: $E2E_TEST_DIR"
    fi
}

function usage() {
    echo "Usage: $0 [all|unit|integration|e2e]"
    echo "  all: Run all tests (unit, integration, and e2e)."
    echo "  unit: Run only unit tests."
    echo "  integration: Run only integration tests."
    echo "  e2e: Run only end-to-end tests."
    exit 1
}

# Main Script
if [[ $# -ne 1 ]]; then
    usage
fi

case "$1" in
    all)
        run_unit_tests
        run_integration_tests
        run_e2e_tests
        ;;
    unit)
        run_unit_tests
        ;;
    integration)
        run_integration_tests
        ;;
    e2e)
        run_e2e_tests
        ;;
    *)
        usage
        ;;
esac