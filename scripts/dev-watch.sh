#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOCACHE="${GOCACHE:-$ROOT_DIR/.cache/go-build}"
ADDR="${ADDR:-:8080}"
DATABASE_URL="${DATABASE_URL:-postgres://media_manager:media_manager@localhost:5432/media_manager?sslmode=disable}"
WEB_DIR="${WEB_DIR:-web/build}"
APP_ENV="${APP_ENV:-development}"

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
			find api cmd internal web/src \
				-type f \
				-not -path '*/node_modules/*' \
				-not -path '*/.svelte-kit/*' \
				-print0 2>/dev/null
			find web \
				-maxdepth 1 \
				-type f \
				\( -name 'package.json' -o -name 'pnpm-lock.yaml' -o -name 'vite.config.ts' -o -name 'svelte.config.*' \) \
				-print0
		} |
			sort -z |
			xargs -0 shasum
	) | shasum
}

restart_app() {
	echo "==> rebuilding contract, web assets, and Go server"
	(
		cd "$ROOT_DIR"
		GOCACHE="$GOCACHE" make build
	)

	if [[ -n "$server_pid" ]] && kill -0 "$server_pid" 2>/dev/null; then
		echo "==> stopping server pid $server_pid"
		kill "$server_pid" 2>/dev/null || true
		wait "$server_pid" 2>/dev/null || true
	fi

	echo "==> starting app on $ADDR"
	(
		cd "$ROOT_DIR"
		ADDR="$ADDR" \
			APP_ENV="$APP_ENV" \
			DATABASE_URL="$DATABASE_URL" \
			WEB_DIR="$WEB_DIR" \
			"$ROOT_DIR/bin/server"
	) &
	server_pid="$!"
}

echo "==> watching source changes; app default URL is http://127.0.0.1${ADDR/:/:}"
while true; do
	current_fingerprint="$(fingerprint_sources)"
	if [[ "$current_fingerprint" != "$last_fingerprint" ]]; then
		last_fingerprint="$current_fingerprint"
		restart_app
	fi
	sleep 1
done
