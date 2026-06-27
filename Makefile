GOCACHE ?= $(CURDIR)/.cache/go-build
GOFLAGS ?=
DATABASE_URL ?= postgres://media_manager:media_manager@localhost:5432/media_manager?sslmode=disable

.PHONY: api-generate api-generate-go api-generate-web build check db-reset dev dev-api dev-web format river-migrate test web-install

api-generate: api-generate-go api-generate-web

api-generate-go:
	mkdir -p internal/httpapi
	GOCACHE=$(GOCACHE) go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

api-generate-web:
	cd web && pnpm run api:generate

build: api-generate
	cd web && pnpm run build
	GOCACHE=$(GOCACHE) go build $(GOFLAGS) -o bin/server ./cmd/server

check: api-generate
	GOCACHE=$(GOCACHE) go test ./...
	cd web && pnpm run check
	cd web && pnpm run lint
	cd web && pnpm run format:check

db-reset:
	APP_ENV=development ALLOW_DEV_RESET=true GOCACHE=$(GOCACHE) go run ./cmd/server reset-dev

river-migrate:
	GOCACHE=$(GOCACHE) go run github.com/riverqueue/river/cmd/river migrate-up --database-url "$(DATABASE_URL)"

dev:
	$(MAKE) dev-api & $(MAKE) dev-web & wait

dev-api:
	GOCACHE=$(GOCACHE) go run ./cmd/server

dev-web:
	cd web && pnpm run dev -- --host 127.0.0.1

format:
	gofmt -w cmd internal
	cd web && pnpm run format

test:
	GOCACHE=$(GOCACHE) go test ./...
	cd web && pnpm run test

web-install:
	cd web && pnpm install
