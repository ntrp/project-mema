#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOCACHE="${GOCACHE:-$ROOT_DIR/.cache/go-build}"
ADDR="${ADDR:-:8080}"
DATABASE_URL="${DATABASE_URL:-postgres://media_manager:media_manager@localhost:5432/media_manager?sslmode=disable}"
WEB_DIR="${WEB_DIR:-web/build}"
APP_ENV="${APP_ENV:-development}"
MEDIA_DATA_DIR="${MEDIA_DATA_DIR:-$ROOT_DIR/.data/media}"

server_pid=""
last_fingerprint=""

cleanup() {
	if [[ -n "$server_pid" ]] && kill -0 "$server_pid" 2>/dev/null; then
		kill "$server_pid" 2>/dev/null || true
		wait "$server_pid" 2>/dev/null || true
	fi
}

trap cleanup EXIT INT TERM

fingerprint_sources() {
	(
		cd "$ROOT_DIR"
		{
			find api cmd internal \
				-type f \
				-print0 2>/dev/null
			find . \
				-maxdepth 1 \
				-type f \
				\( -name 'go.mod' -o -name 'go.sum' -o -name 'tools.go' \) \
				-print0
		} |
			sort -z |
			xargs -0 shasum
	) | shasum
}

restart_api() {
	echo "==> regenerating API contract bindings and rebuilding Go server"
	(
		cd "$ROOT_DIR"
		GOCACHE="$GOCACHE" make api-generate
		GOCACHE="$GOCACHE" go build -o bin/server ./cmd/server
	)

	if [[ -n "$server_pid" ]] && kill -0 "$server_pid" 2>/dev/null; then
		echo "==> stopping API pid $server_pid"
		kill "$server_pid" 2>/dev/null || true
		wait "$server_pid" 2>/dev/null || true
	fi

	echo "==> starting API on $ADDR"
	(
		cd "$ROOT_DIR"
		ADDR="$ADDR" \
			APP_ENV="$APP_ENV" \
			DATABASE_URL="$DATABASE_URL" \
			MEDIA_DATA_DIR="$MEDIA_DATA_DIR" \
			WEB_DIR="$WEB_DIR" \
			"$ROOT_DIR/bin/server"
	) &
	server_pid="$!"
}

echo "==> watching API sources; API URL is http://127.0.0.1${ADDR/:/:}"
while true; do
	current_fingerprint="$(fingerprint_sources)"
	if [[ "$current_fingerprint" != "$last_fingerprint" ]]; then
		last_fingerprint="$current_fingerprint"
		restart_api
	fi
	sleep 1
done
