#!/usr/bin/env sh
set -eu

root_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp_dir=$(mktemp -d "$root_dir/.research/sqlc-verify.XXXXXX")

cleanup() {
	rm -rf "$tmp_dir"
}
trap cleanup EXIT INT TERM

tmp_config="$tmp_dir/sqlc.yaml"
tmp_output="$tmp_dir/generated"

sed \
	-e "s#schema: internal/storage/migrations/00001_initial_schema.sql#schema: ../../internal/storage/migrations/00001_initial_schema.sql#" \
	-e "s#- internal/storage/migrations/00001_initial_schema.sql#- ../../internal/storage/migrations/00001_initial_schema.sql#" \
	-e "s#- internal/storage/sqlc_schema/river_job.sql#- ../../internal/storage/sqlc_schema/river_job.sql#" \
	-e "s#queries: internal/storage/queries#queries: ../../internal/storage/queries#" \
	-e "s#out: internal/storage/generated#out: generated#" \
	"$root_dir/sqlc.yaml" > "$tmp_config"

(
	cd "$root_dir"
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f "$tmp_config"
)

status=0

if ! diff -ru "$root_dir/internal/storage/generated" "$tmp_output"; then
	printf '%s\n' "sqlc generated storage artifacts are stale: run make sqlc-generate." >&2
	status=1
fi

exit "$status"
