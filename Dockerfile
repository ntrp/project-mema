FROM node:24-bookworm AS web-build
WORKDIR /src/web
COPY web/package.json web/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY web ./
COPY api /src/api
RUN pnpm run api:generate && pnpm run build

FROM golang:1.26-bookworm AS go-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-build /src/web/build ./web/build
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        ffmpeg \
        mediainfo \
        mkvtoolnix \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=go-build /out/server /app/server
COPY --from=web-build /src/web/build /app/web
ENV ADDR=:18080
ENV WEB_DIR=/app/web
EXPOSE 18080
ENTRYPOINT ["/app/server"]
