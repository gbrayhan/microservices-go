#!/usr/bin/env bash
set -euo pipefail

trap 'error_handler $LINENO' ERR

# Parse command line arguments
VERBOSE=false
FEATURE_FILE=""
SCENARIO_TAGS=""
while [[ $# -gt 0 ]]; do
  case $1 in
    -v|--verbose)
      VERBOSE=true
      shift
      ;;
    -f|--feature)
      FEATURE_FILE="$2"
      shift 2
      ;;
    -t|--tags)
      SCENARIO_TAGS="$2"
      shift 2
      ;;
    *)
      echo "Usage: $0 [-v|--verbose] [-f|--feature <feature_file>] [-t|--tags <scenario_tags>]"
      echo "  -v, --verbose    Enable verbose output for tests"
      echo "  -f, --feature    Run only the specified feature file (e.g., auth.feature)"
      echo "  -t, --tags       Run only scenarios with specific tags (e.g., @smoke @critical)"
      echo
      echo "Examples:"
      echo "  $0                                    # Run all integration tests"
      echo "  $0 -v                                 # Run all tests with verbose output"
      echo "  $0 -f auth.feature                    # Run only auth.feature tests"
      echo "  $0 -t @smoke                          # Run only scenarios tagged with @smoke"
      echo "  $0 -f auth.feature -t @critical       # Run only critical scenarios in auth.feature"
      echo "  $0 -v -f order.feature -t @smoke      # Run smoke tests in order.feature with verbose output"
      echo
      echo "Tag Examples:"
      echo "  @smoke     - Quick tests for basic functionality"
      echo "  @critical  - Critical path tests"
      echo "  @slow      - Tests that take longer to run"
      echo "  @auth      - Authentication related tests"
      echo "  @api       - API endpoint tests"
      exit 1
      ;;
  esac
done

error_handler() {
  local exit_code=$?
  local line_no=$1
  echo "❌ Error on line $line_no (exit code $exit_code)."
  echo "   ⮡ Check the output of the previous steps to identify the cause."
  exit "$exit_code"
}

# Function to validate required environment variables
validate_required_env_vars() {
  local missing_vars=()
  
  # Database variables
  [[ -z "${DB_HOST:-}" ]] && missing_vars+=("DB_HOST")
  [[ -z "${DB_PORT:-}" ]] && missing_vars+=("DB_PORT")
  [[ -z "${DB_USER:-}" ]] && missing_vars+=("DB_USER")
  [[ -z "${DB_PASSWORD:-}" ]] && missing_vars+=("DB_PASSWORD")
  [[ -z "${DB_NAME:-}" ]] && missing_vars+=("DB_NAME")
  [[ -z "${DB_SSLMODE:-}" ]] && missing_vars+=("DB_SSLMODE")
  
  # JWT variables
  [[ -z "${ACCESS_SECRET_KEY:-}" ]] && missing_vars+=("ACCESS_SECRET_KEY")
  [[ -z "${REFRESH_SECRET_KEY:-}" ]] && missing_vars+=("REFRESH_SECRET_KEY")
  [[ -z "${JWT_ISSUER:-}" ]] && missing_vars+=("JWT_ISSUER")
  
  # External services
  [[ -z "${IMGUR_CLIENT_ID:-}" ]] && missing_vars+=("IMGUR_CLIENT_ID")
  
  
  # Initial user (for migrations)
  [[ -z "${START_USER_EMAIL:-}" ]] && missing_vars+=("START_USER_EMAIL")
  [[ -z "${START_USER_PW:-}" ]] && missing_vars+=("START_USER_PW")
  
  # Optional variables with defaults
  [[ -z "${APP_PORT:-}" ]] && export APP_PORT=8080
  [[ -z "${ACCESS_TOKEN_TTL:-}" ]] && export ACCESS_TOKEN_TTL=15
  [[ -z "${REFRESH_TOKEN_TTL:-}" ]] && export REFRESH_TOKEN_TTL=10080
  
  if [[ ${#missing_vars[@]} -gt 0 ]]; then
    echo "❌ Error: The following required environment variables are not set:"
    printf "   - %s\n" "${missing_vars[@]}"
    echo
    echo "Please set these variables before running the integration tests."
    exit 1
  fi
}

BUILD_NAME="dev-aceso"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🔍 Validating required environment variables..."
validate_required_env_vars
echo "✅ All required environment variables are set"

echo "🛠 Compiling project from '$PROJECT_ROOT'..."
cd "$PROJECT_ROOT"
go build -o "$BUILD_NAME" .

echo "🔍 Ensuring port $APP_PORT is free..."
PIDS=$(lsof -ti tcp:"$APP_PORT" || true)
if [[ -n "$PIDS" ]]; then
  echo "⚠️ Killing stale process(es) on port $APP_PORT: $PIDS"
  kill -9 $PIDS
  while lsof -ti tcp:"$APP_PORT" >/dev/null; do sleep 0.1; done
  echo "✅ Port $APP_PORT is now free."
else
  echo "✅ Port $APP_PORT was already free."
fi

echo "🔧 Environment variables are already set and validated"

echo "▶️ Starting the application (logs suppressed)…"
"$PROJECT_ROOT/$BUILD_NAME" > /dev/null 2>&1 &
APP_PID=$!

until lsof -ti tcp:"$APP_PORT" >/dev/null; do sleep 0.1; done
echo "✅ App listening on port $APP_PORT (PID $APP_PID)"

echo
echo "🧪 Running integration tests…"
trap '' ERR
set +e

# Set integration test environment variable
export INTEGRATION_TEST=true

# Build the test command with optional verbose flag
TEST_CMD="go test -count=1 ./Test/integration -tags=integration"
if [[ "$VERBOSE" == "true" ]]; then
  TEST_CMD="$TEST_CMD -v"
fi

if [[ -n "$FEATURE_FILE" ]]; then
  # Validate that the feature file exists
  FEATURE_PATH="$PROJECT_ROOT/Test/integration/features/$FEATURE_FILE"
  if [[ ! -f "$FEATURE_PATH" ]]; then
    echo "❌ Error: Feature file '$FEATURE_FILE' not found at '$FEATURE_PATH'"
    echo "Available feature files:"
    ls -1 "$PROJECT_ROOT/Test/integration/features/"*.feature | sed 's|.*/||' | sed 's/^/  - /'
    exit 1
  fi
  echo "🎯 Running only feature file: $FEATURE_FILE"
  export INTEGRATION_FEATURE_FILE="$FEATURE_FILE"
else
  echo "🧪 Running all integration tests..."
  unset INTEGRATION_FEATURE_FILE
fi

if [[ -n "$SCENARIO_TAGS" ]]; then
  echo "🎯 Running scenarios with tags: $SCENARIO_TAGS"
  export INTEGRATION_SCENARIO_TAGS="$SCENARIO_TAGS"
else
  unset INTEGRATION_SCENARIO_TAGS
fi

$TEST_CMD
TEST_EXIT=$?
set -e
trap 'error_handler $LINENO' ERR

if [ $TEST_EXIT -eq 0 ]; then
  echo "🎉 Integration tests passed!"
else
  echo "⚠️ Integration tests finished with exit code $TEST_EXIT."
fi

echo "🛑 Stopping the application (PID $APP_PID)…"
kill "$APP_PID" 2>/dev/null || true
echo "✅ Application stopped."

echo "💡 All done."
exit $TEST_EXIT 