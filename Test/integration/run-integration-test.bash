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

BUILD_NAME="dev-aceso"
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

echo "🔧 Exporting environment variables (if not already set)..."
export ACCESS_SECRET_KEY="${ACCESS_SECRET_KEY:-yourAccessSecretKey}"
export ACCESS_TOKEN_TTL="${ACCESS_TOKEN_TTL:-15}"
export APP_PORT
export DB_HOST="${DB_HOST:-localhost}"
export DB_NAME="${DB_NAME:-aceso}"
export DB_PASSWORD="${DB_PASSWORD:-yourpassword}"
export DB_PORT="${DB_PORT:-5432}"
export DB_SSLMODE="${DB_SSLMODE:-disable}"
export DB_USER="${DB_USER:-app_user}"
export IMGUR_CLIENT_ID="${IMGUR_CLIENT_ID:-yourImgurClientId}"
export JWT_ISSUER="${JWT_ISSUER:-aceso}"
export REFRESH_SECRET_KEY="${REFRESH_SECRET_KEY:-yourRefreshSecretKey}"
export REFRESH_TOKEN_TTL="${REFRESH_TOKEN_TTL:-10080}"
export WKHTMLTOPDF_BIN="${WKHTMLTOPDF_BIN:-/usr/local/bin/wkhtmltopdf}"

echo "▶️ Starting the application (logs suppressed)…"
"$PROJECT_ROOT/$BUILD_NAME" > /dev/null 2>&1 &
APP_PID=$!

until lsof -ti tcp:"$APP_PORT" >/dev/null; do sleep 0.1; done
echo "✅ App listening on port $APP_PORT (PID $APP_PID)"

echo
echo "🧪 Running integration tests…"
trap '' ERR
set +e
go test -count=1 ./Test/integration -tags=integration
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
