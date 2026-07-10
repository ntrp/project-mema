#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOCACHE="${GOCACHE:-$ROOT_DIR/.cache/go-build}"
ADDR="${ADDR:-:18080}"
DATABASE_URL="${DATABASE_URL:-postgres://media_manager:media_manager@localhost:15432/media_manager?sslmode=disable}"
WEB_DIR="${WEB_DIR:-web/build}"
APP_ENV="${APP_ENV:-development}"
MEDIA_DATA_DIR="${MEDIA_DATA_DIR:-$ROOT_DIR/.data/media}"
OPENAPI_FINGERPRINT_FILE="${OPENAPI_FINGERPRINT_FILE:-$ROOT_DIR/.cache/dev-api-openapi.sha}"

server_pid=""
last_fingerprint=""

openapi_fingerprint() {
	(
		cd "$ROOT_DIR"
		shasum api/openapi.yaml
	) | shasum
}

ensure_api_generated() {
	local current_fingerprint
	local last_openapi_fingerprint

	mkdir -p "$(dirname "$OPENAPI_FINGERPRINT_FILE")"
	current_fingerprint="$(openapi_fingerprint)"
	last_openapi_fingerprint="$(cat "$OPENAPI_FINGERPRINT_FILE" 2>/dev/null || true)"
	if [[ "$current_fingerprint" == "$last_openapi_fingerprint" ]]; then
		return
	fi

	echo "==> regenerating API contract bindings"
	(
		cd "$ROOT_DIR"
		GOCACHE="$GOCACHE" make api-generate
	)
	printf "%s\n" "$current_fingerprint" >"$OPENAPI_FINGERPRINT_FILE"
}

run_api() {
	ensure_api_generated
	echo "==> rebuilding Go server"
	(
		cd "$ROOT_DIR"
		GOCACHE="$GOCACHE" go build -o bin/server ./cmd/server
	)

	echo "==> starting API on $ADDR"
	cd "$ROOT_DIR"
	exec env \
		ADDR="$ADDR" \
		APP_ENV="$APP_ENV" \
		DATABASE_URL="$DATABASE_URL" \
		MEDIA_DATA_DIR="$MEDIA_DATA_DIR" \
		WEB_DIR="$WEB_DIR" \
		"$ROOT_DIR/bin/server"
}

cleanup() {
	if [[ -n "$server_pid" ]] && kill -0 "$server_pid" 2>/dev/null; then
		kill "$server_pid" 2>/dev/null || true
		wait "$server_pid" 2>/dev/null || true
	fi
}

trap cleanup EXIT INT TERM

if [[ "${1:-}" == "--run-api" ]]; then
	trap - EXIT INT TERM
	run_api
fi

fingerprint_sources() {
	(
		cd "$ROOT_DIR"
		{
			find api cmd internal \
				-type f \
				-not -path 'internal/httpapi/openapi.gen.go' \
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
	if [[ -n "$server_pid" ]] && kill -0 "$server_pid" 2>/dev/null; then
		echo "==> stopping API pid $server_pid"
		kill "$server_pid" 2>/dev/null || true
		wait "$server_pid" 2>/dev/null || true
	fi

	"$ROOT_DIR/scripts/dev-api-watch.sh" --run-api &
	server_pid="$!"
}

echo "==> watching API sources; API URL is http://127.0.0.1${ADDR/:/:}"
if command -v watchexec >/dev/null 2>&1; then
	watch_args=(
		--restart
		--debounce 250ms
		--watch "$ROOT_DIR/api"
		--watch "$ROOT_DIR/cmd"
		--watch "$ROOT_DIR/internal"
		--ignore "$ROOT_DIR/bin/**"
		--ignore "$ROOT_DIR/.cache/**"
		--ignore "$ROOT_DIR/internal/httpapi/openapi.gen.go"
		--exts go,yaml,yml,json,sql,mod,sum
	)
	for file in go.mod go.sum tools.go; do
		if [[ -e "$ROOT_DIR/$file" ]]; then
			watch_args+=(--watch "$ROOT_DIR/$file")
		fi
	done
	echo "==> using watchexec for file watching"
	exec watchexec "${watch_args[@]}" -- "$ROOT_DIR/scripts/dev-api-watch.sh" --run-api
fi

echo "==> watchexec not found; falling back to polling"
while true; do
	current_fingerprint="$(fingerprint_sources)"
	if [[ "$current_fingerprint" != "$last_fingerprint" ]]; then
		last_fingerprint="$current_fingerprint"
		restart_api
	fi
	sleep 1
done
