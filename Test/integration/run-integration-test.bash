#!/usr/bin/env bash
set -euo pipefail

trap 'error_handler $LINENO' ERR

error_handler() {
  local exit_code=$?
  local line_no=$1
  echo "❌ Error on line $line_no (exit code $exit_code)."
  echo "   ⮡ Check the output of the previous steps to identify the cause."
  exit "$exit_code"
}

BUILD_NAME="app-microservice"
: "${APP_PORT:=8080}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

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

echo "🔧 Validating required environment variables..."
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
  echo "❌ Error: The following required environment variables are not set:"
  echo ""
  for var in "${missing_vars[@]}"; do
    echo "   • $var"
  done
  echo ""
  echo "💡 Please set these variables before running integration tests."
  echo ""
  exit 1
fi

echo "✅ All required environment variables are set"

echo "▶️ Starting the application (showing logs)…"
"$PROJECT_ROOT/$BUILD_NAME" &
APP_PID=$!

until lsof -ti tcp:"$APP_PORT" >/dev/null; do sleep 0.1; done
echo "✅ App listening on port $APP_PORT (PID $APP_PID)"

echo
echo "🧪 Running integration tests…"
trap '' ERR
set +e
go test -count=1 ./Test/integration -tags=integration -v
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
