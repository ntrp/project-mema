GOCACHE ?= $(CURDIR)/.cache/go-build
GOFLAGS ?=
DATABASE_URL ?= postgres://media_manager:media_manager@localhost:15432/media_manager?sslmode=disable
MEDIA_DATA_DIR ?= $(CURDIR)/.data/media

.PHONY: api-generate api-generate-go api-generate-web build check coverage coverage-backend coverage-web db-reset dev dev-api dev-api-watch dev-watch dev-web format river-migrate storage-generate test test-api test-deps test-e2e web-install

api-generate: api-generate-go api-generate-web

api-generate-go:
	mkdir -p internal/httpapi
	GOCACHE=$(GOCACHE) go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

api-generate-web:
	cd web && pnpm run api:generate

storage-generate:
	GOCACHE=$(GOCACHE) go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

build: api-generate
	cd web && pnpm run build
	GOCACHE=$(GOCACHE) go build $(GOFLAGS) -o bin/server ./cmd/server

check: api-generate
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

db-reset:
	APP_ENV=development ALLOW_DEV_RESET=true GOCACHE=$(GOCACHE) go run ./cmd/server reset-dev

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
