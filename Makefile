GO ?= go
BUN ?= bun

.PHONY: install dev-api dev-web build test check clean docker-test docker-prod

install:
	$(GO) mod download
	cd web && $(BUN) install --frozen-lockfile

dev-api:
	$(GO) run ./cmd/server

dev-web:
	cd web && $(BUN) run dev

build:
	cd web && $(BUN) run build
	mkdir -p bin
	CGO_ENABLED=1 $(GO) build -o bin/server ./cmd/server

test:
	$(GO) test ./...

check:
	$(GO) vet ./...
	cd web && $(BUN) run build

clean:
	rm -rf bin web/dist

docker-test:
	docker compose -f deploy/test/compose.yaml up --build -d

docker-prod:
	docker compose -f deploy/prod/compose.yaml up --build -d
