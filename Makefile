GO ?= go
BUN ?= bun
IMAGE ?= go-http-template:latest

.PHONY: install dev dev-api dev-web swagger build test check clean docker-build docker-test docker-prod help

install: ## 下载 Go 和前端依赖
	$(GO) mod download
	cd web && $(BUN) install --frozen-lockfile

dev: ## 同时启动后端和前端开发服务
	$(MAKE) -j2 dev-api dev-web

dev-api: ## 启动后端
	$(GO) run ./cmd/server

dev-web: ## 启动前端开发服务
	cd web && $(BUN) run dev

swagger: ## 生成 Swagger 接口文档
	$(GO) run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g main.go -d cmd/server,internal/controller,internal/router,internal/dto,internal/pkg/response --parseInternal -o docs

build: ## 构建前端和后端
	cd web && $(BUN) run build
	mkdir -p bin
	CGO_ENABLED=1 $(GO) build -o bin/server ./cmd/server

test: ## 运行 Go 测试
	$(GO) test ./...

check: ## 执行 Go vet 和前端检查、构建
	$(GO) vet ./...
	cd web && $(BUN) run build

clean: ## 删除构建产物
	rm -rf bin web/dist

docker-build: ## 构建 Docker 镜像（可通过 IMAGE 指定名称和标签）
	docker build -f docker/Dockerfile -t $(IMAGE) .

docker-test: ## 构建并后台启动测试环境
	docker compose -f deploy/test/compose.yaml up --build -d

docker-prod: ## 构建并后台启动生产环境
	docker compose -f deploy/prod/compose.yaml up --build -d

help: ## 显示可用命令
	@printf '用法: make <目标>\n\n'
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
