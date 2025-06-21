#!/usr/bin/env bash
set -euo pipefail

# Integration Test Runner Script
# This script runs comprehensive integration tests for the microservices application

trap 'error_handler $LINENO' ERR

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

error_handler() {
  local exit_code=$?
  local line_no=$1
  print_error "Error on line $line_no (exit code $exit_code)."
  print_error "Check the output of the previous steps to identify the cause."
  exit "$exit_code"
}

# Configuration
BUILD_NAME="app-microservice"
: "${APP_PORT:=8080}"
: "${TEST_TIMEOUT:=300}"  # 5 minutes timeout

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

print_status "Starting integration test suite..."

print_status "Compiling project from '$PROJECT_ROOT'..."
cd "$PROJECT_ROOT"

# Clean previous build
if [[ -f "$BUILD_NAME" ]]; then
    rm "$BUILD_NAME"
fi

# Build the application
go build -o "$BUILD_NAME" .
print_success "Application compiled successfully"

print_status "Ensuring port $APP_PORT is free..."
PIDS=$(lsof -ti tcp:"$APP_PORT" || true)
if [[ -n "$PIDS" ]]; then
  print_warning "Killing stale process(es) on port $APP_PORT: $PIDS"
  kill -9 $PIDS
  while lsof -ti tcp:"$APP_PORT" >/dev/null; do sleep 0.1; done
  print_success "Port $APP_PORT is now free."
else
  print_success "Port $APP_PORT was already free."
fi

print_status "Validating required environment variables..."
# Array of required environment variables
required_vars=(
  "DB_HOST"
  "DB_NAME"
  "DB_PASSWORD"
  "DB_PORT"
  "DB_SSLMODE"
  "DB_USER"
  "JWT_ACCESS_SECRET"
  "JWT_ACCESS_TIME_MINUTE"
  "JWT_REFRESH_SECRET"
  "JWT_REFRESH_TIME_HOUR"
  "START_USER_EMAIL"
  "START_USER_PW"
)

# Check all required variables and collect missing ones
missing_vars=()
for var in "${required_vars[@]}"; do
  if [[ -z "${!var:-}" ]]; then
    missing_vars+=("$var")
  fi
done

# If there are missing variables, show them all and exit
if [[ ${#missing_vars[@]} -gt 0 ]]; then
  print_error "The following required environment variables are not set:"
  echo ""
  for var in "${missing_vars[@]}"; do
    echo "   • $var"
  done
  echo ""
  print_status "Please set these variables before running integration tests."
  echo ""
  exit 1
fi

print_success "All required environment variables are set"

# Check if database is accessible
print_status "Testing database connectivity..."
if command -v psql &> /dev/null; then
    if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" &> /dev/null; then
        print_success "Database connection successful"
    else
        print_error "Cannot connect to database. Please check your database configuration."
        exit 1
    fi
else
    print_warning "psql not found, skipping database connectivity test"
fi

print_status "Starting the application..."
"$PROJECT_ROOT/$BUILD_NAME" &
APP_PID=$!

# Wait for application to start
print_status "Waiting for application to start on port $APP_PORT..."
timeout_counter=0
while ! lsof -ti tcp:"$APP_PORT" >/dev/null 2>&1; do
    sleep 1
    timeout_counter=$((timeout_counter + 1))
    if [[ $timeout_counter -ge $TEST_TIMEOUT ]]; then
        print_error "Application failed to start within $TEST_TIMEOUT seconds"
        kill "$APP_PID" 2>/dev/null || true
        exit 1
    fi
done

print_success "Application started successfully (PID $APP_PID)"

# Wait a bit more for the application to fully initialize
sleep 3

# Test health endpoint
print_status "Testing health endpoint..."
if curl -f -s "http://localhost:$APP_PORT/v1/health" >/dev/null; then
    print_success "Health endpoint is responding"
else
    print_error "Health endpoint is not responding"
    kill "$APP_PID" 2>/dev/null || true
    exit 1
fi

echo
print_status "Running integration tests..."
trap '' ERR
set +e

# Run integration tests
go test -count=1 ./Test/integration -tags=integration -v
TEST_EXIT=$?

set -e
trap 'error_handler $LINENO' ERR

echo
if [ $TEST_EXIT -eq 0 ]; then
    print_success "🎉 Integration tests passed!"
else
    print_warning "⚠️ Integration tests finished with exit code $TEST_EXIT."
fi

print_status "Stopping the application (PID $APP_PID)..."
kill "$APP_PID" 2>/dev/null || true

# Wait for application to stop
timeout_counter=0
while lsof -ti tcp:"$APP_PORT" >/dev/null 2>&1; do
    sleep 0.1
    timeout_counter=$((timeout_counter + 1))
    if [[ $timeout_counter -ge 30 ]]; then
        print_warning "Application did not stop gracefully, forcing termination"
        kill -9 "$APP_PID" 2>/dev/null || true
        break
    fi
done

print_success "Application stopped successfully"

# Clean up build artifact
if [[ -f "$BUILD_NAME" ]]; then
    rm "$BUILD_NAME"
    print_status "Cleaned up build artifact"
fi

echo
print_status "Integration test suite completed"
print_status "Exit code: $TEST_EXIT"

if [ $TEST_EXIT -eq 0 ]; then
    print_success "✅ All tests passed successfully!"
else
    print_error "❌ Some tests failed. Please check the test output above."
fi

exit $TEST_EXIT
