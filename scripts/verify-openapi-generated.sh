#!/usr/bin/env sh
set -eu

root_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp_dir=$(mktemp -d "$root_dir/web/.openapi-verify.XXXXXX")

cleanup() {
	rm -rf "$tmp_dir"
}
trap cleanup EXIT INT TERM

go_config="$tmp_dir/oapi-codegen.yaml"
go_output="$tmp_dir/openapi.gen.go"
web_output="$tmp_dir/schema.d.ts"
query_config="$tmp_dir/orval.config.ts"
query_output="$tmp_dir/tanstack.ts"

sed "s#^output: .*#output: $go_output#" "$root_dir/api/oapi-codegen.yaml" > "$go_config"
sed \
	-e "s#target: '../api/openapi.yaml'#target: '$root_dir/api/openapi.yaml'#" \
	-e "s#target: 'src/lib/api/generated/tanstack.ts'#target: '$query_output'#" \
	"$root_dir/web/orval.config.ts" > "$query_config"

GOCACHE="${GOCACHE:-"$root_dir/.cache/go-build"}" \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
	-config "$go_config" \
	"$root_dir/api/openapi.yaml"

(
	cd "$root_dir/web"
	pnpm exec openapi-typescript ../api/openapi.yaml -o "$web_output"
	pnpm exec orval --config "$query_config"
	pnpm exec prettier --write "$web_output" >/dev/null
	pnpm exec prettier --write "$query_output" >/dev/null
)

status=0

if ! diff -u "$root_dir/internal/httpapi/openapi.gen.go" "$go_output"; then
	printf '%s\n' "OpenAPI Go server artifact is stale: run make api-generate." >&2
	status=1
fi

if ! diff -u "$root_dir/web/src/lib/api/generated/schema.d.ts" "$web_output"; then
	printf '%s\n' "OpenAPI TypeScript schema artifact is stale: run make api-generate." >&2
	status=1
fi

if ! diff -u "$root_dir/web/src/lib/api/generated/tanstack.ts" "$query_output"; then
	printf '%s\n' "OpenAPI TanStack query artifact is stale: run make api-generate." >&2
	status=1
fi

exit "$status"
