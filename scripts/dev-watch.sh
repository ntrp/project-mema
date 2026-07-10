#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MEDIA_DATA_DIR="${MEDIA_DATA_DIR:-$ROOT_DIR/.data/media}"
ADDR="${ADDR:-:18080}"
WEB_PORT="${WEB_PORT:-15173}"

api_pid=""
cleaned_up=0

kill_tree() {
	local pid="$1"
	local child

	if [[ -z "$pid" ]] || ! kill -0 "$pid" 2>/dev/null; then
		return
	fi

	for child in $(pgrep -P "$pid" 2>/dev/null || true); do
		kill_tree "$child"
	done
	kill "$pid" 2>/dev/null || true
}

cleanup() {
	if [[ "$cleaned_up" -eq 1 ]]; then
		return
	fi
	cleaned_up=1

	kill_tree "$api_pid"
	if [[ -n "$api_pid" ]]; then
		wait "$api_pid" 2>/dev/null || true
	fi
}

trap cleanup EXIT
trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM

start_api() {
	(
		cd "$ROOT_DIR"
		MEDIA_DATA_DIR="$MEDIA_DATA_DIR" ADDR="$ADDR" ./scripts/dev-api-watch.sh
	) &
	api_pid="$!"
}

start_web() {
	(
		cd "$ROOT_DIR/web"
		NVIM_LISTEN_ADDRESS=/tmp/project-mema.nvim \
			LAUNCH_EDITOR="$ROOT_DIR/scripts/open-in-nvim.sh" \
			pnpm exec vite dev --host 0.0.0.0 --port "$WEB_PORT" --clearScreen=false
	)
}

echo "==> starting API watcher on http://127.0.0.1${ADDR/:/:}"
start_api

echo "==> starting frontend dev server on http://127.0.0.1:$WEB_PORT"
start_web
