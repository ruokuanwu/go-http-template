FROM oven/bun:1.3.6 AS frontend
WORKDIR /app/web
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile
COPY web/ ./
RUN bun run build

FROM golang:1.25.5-bookworm AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY docs/ ./docs/
RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 10001 app && useradd --uid 10001 --gid app --no-create-home app \
    && mkdir -p /app/data && chown app:app /app/data
WORKDIR /app
COPY --from=backend /server ./server
COPY --from=frontend /app/web/dist ./web/dist
ENV HTTP_ADDR=:8080 DB_PATH=/app/data/app.db WEB_DIR=/app/web/dist GIN_MODE=release
USER app
EXPOSE 8080
ENTRYPOINT ["./server"]
