GOCACHE ?= $(CURDIR)/.cache/go-build
GOFLAGS ?=
DATABASE_URL ?= postgres://media_manager:media_manager@localhost:15432/media_manager?sslmode=disable
MEDIA_DATA_DIR ?= $(CURDIR)/.data/media
IMAGE ?= project-mema:local
APP_VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 0.0.0-dev)
APP_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
APP_SOURCE_URL ?= Not configured

.PHONY: api-generate api-generate-go api-generate-web build check coverage coverage-backend coverage-web db-clean db-reset db-seed-local dev dev-api dev-api-watch dev-watch dev-web docker-build docs-build docs-dev docs-install docs-preview format river-migrate sqlc-generate test test-api test-deps test-e2e verify-generated verify-sqlc-generated web-install

api-generate: api-generate-go api-generate-web

api-generate-go:
	mkdir -p internal/httpapi
	GOCACHE=$(GOCACHE) go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

api-generate-web:
	cd web && pnpm run api:generate

sqlc-generate:
	GOCACHE=$(GOCACHE) go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

verify-generated:
	sh scripts/verify-openapi-generated.sh

verify-sqlc-generated:
	sh scripts/verify-sqlc-generated.sh

build: api-generate
	cd web && pnpm run build
	GOCACHE=$(GOCACHE) go build $(GOFLAGS) -o bin/server ./cmd/server

docker-build:
	docker build \
		--build-arg APP_VERSION="$(APP_VERSION)" \
		--build-arg APP_COMMIT="$(APP_COMMIT)" \
		--build-arg APP_SOURCE_URL="$(APP_SOURCE_URL)" \
		-t "$(IMAGE)" .

check: verify-generated verify-sqlc-generated
	GOCACHE=$(GOCACHE) go test ./...
	cd web && pnpm run check
	cd web && pnpm run lint
	cd web && pnpm run format:check

coverage: coverage-backend coverage-web

coverage-backend: test-deps
	mkdir -p coverage
	GOCACHE=$(GOCACHE) DATABASE_URL=$(DATABASE_URL) go test ./... -coverpkg=./... -covermode=atomic -coverprofile=coverage/backend.raw.out
	grep -v '/openapi.gen.go:' coverage/backend.raw.out > coverage/backend.out
	go tool cover -func=coverage/backend.out > coverage/backend.txt
	awk '/^total:/ { sub("%", "", $$3); if ($$3 < 60) exit 1 }' coverage/backend.txt

coverage-web:
	cd web && pnpm run test:coverage

db-clean:
	GOCACHE=$(GOCACHE) DATABASE_URL=$(DATABASE_URL) go run ./cmd/devdb clean

db-reset:
	GOCACHE=$(GOCACHE) DATABASE_URL=$(DATABASE_URL) go run ./cmd/devdb reset

db-seed-local:
	GOCACHE=$(GOCACHE) DATABASE_URL=$(DATABASE_URL) go run ./cmd/devdb seed-local

river-migrate:
	GOCACHE=$(GOCACHE) go run github.com/riverqueue/river/cmd/river migrate-up --database-url "$(DATABASE_URL)"

dev:
	$(MAKE) dev-api-watch & $(MAKE) dev-web & wait

dev-api:
	GOCACHE=$(GOCACHE) MEDIA_DATA_DIR=$(MEDIA_DATA_DIR) go run ./cmd/server

dev-api-watch:
	MEDIA_DATA_DIR=$(MEDIA_DATA_DIR) ./scripts/dev-api-watch.sh

dev-watch:
	./scripts/dev-watch.sh

dev-web:
	cd web && pnpm exec vite dev --host 127.0.0.1 --port 15173

docs-install:
	cd docs/website && pnpm install

docs-dev:
	cd docs/website && pnpm exec astro dev --host 0.0.0.0 --port 15174

docs-build:
	cd docs/website && pnpm run build

docs-preview:
	cd docs/website && pnpm exec astro preview --host 0.0.0.0 --port 15174

format:
	gofmt -w cmd internal
	cd web && pnpm run format

test:
	GOCACHE=$(GOCACHE) go test ./...
	cd web && pnpm run test

test-api: test-deps
	GOCACHE=$(GOCACHE) DATABASE_URL=$(DATABASE_URL) go test ./internal/httpapi ./internal/indexers ./internal/metadata

test-deps:
	docker compose up -d postgres

test-e2e: test-deps
	cd web && pnpm run e2e

web-install:
	cd web && pnpm install
