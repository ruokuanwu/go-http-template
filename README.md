# Go HTTP 全栈模板

最小的前后端同仓库模板，通过任务列表演示添加、查询、完成状态切换与删除，数据持久化到 SQLite。开发时 Vite 代理 API；生产时 Go 同时提供 API 和前端构建产物。

## 技术栈

- 后端：Go 1.25.5、Gin、GORM、SQLite、Zap、samber/oops。
- 前端：Bun 1.3.6、React、TypeScript、Vite、Tailwind CSS v4、shadcn/ui。
- shadcn/ui 组件源码保存在前端工程中，可直接修改。

## 环境要求

- Go 1.25.5 或兼容版本。
- Bun 1.3.6 或兼容版本。
- GCC 与 CGO：GORM SQLite 驱动使用 `go-sqlite3`，构建需要 C 编译器；不需要独立 SQLite 服务。
- Make；容器部署另需 Docker Engine 与 Compose 插件。

确保 `go` 和 `bun` 在 `PATH` 中。使用 Go 调试配置需要 VS Code Go 扩展与 Delve。

## 快速开始

```bash
make install
```

分别在两个终端启动：

```bash
make dev-api
```

```bash
make dev-web
```

打开 http://localhost:5173。前端 `/api` 请求代理到 `http://127.0.0.1:8080`，不需要配置 CORS。

生产构建与本地运行：

```bash
make build
GIN_MODE=release ./bin/server
```

打开 http://localhost:8080。请在项目根目录运行；分发时一并提供 `bin/server` 和 `web/dist`，或通过 `WEB_DIR` 指定构建产物位置。

## 配置

后端从进程环境变量读取配置，**不会自动加载 `.env`**；配置示例见 [.env.example](./.env.example)。

| 变量 | 默认值 | 用途 |
|---|---|---|
| `HTTP_ADDR` | `:8080` | HTTP 监听地址 |
| `DB_PATH` | `data/app.db` | SQLite 文件路径，父目录自动创建 |
| `WEB_DIR` | `web/dist` | 前端构建产物目录 |
| `GIN_MODE` | Gin 默认模式 | 生产运行设置为 `release` |
| `DEBUG` | `false` | 仅为 `true` 时注册 Swagger UI 与文档路由以及其他debug项 |

首次启动自动迁移任务表。SQLite 文件不应提交到版本控制；停止服务后再复制数据库文件作为简单备份。修改后端监听端口时，也需调整前端 Vite 代理目标。

## 常用命令

| 命令 | 说明 |
|---|---|
| `make install` | 下载 Go 依赖并按 Bun 锁文件安装前端依赖 |
| `make dev-api` | 启动后端 |
| `make dev-web` | 启动 Vite 开发服务 |
| `make swagger` | 使用固定版本 Swaggo 生成接口文档 |
| `make build` | TypeScript 检查、前端生产构建、Go 构建 |
| `make test` | 执行 Go 测试命令；模板当前没有测试文件 |
| `make check` | Go vet 与前端类型检查、构建 |
| `make clean` | 删除构建产物，不删除数据库 |
| `make docker-test` | 构建并后台启动测试环境 |
| `make docker-prod` | 构建并后台启动生产环境 |

也可通过 `make GO=/path/to/go BUN=/path/to/bun build` 指定工具路径。

## API

| 方法 | 路径 | 请求体 | 成功状态 |
|---|---|---|---|
| GET | `/api/todos` | — | 200 |
| POST | `/api/todos` | `{"title":"学习 Go"}` | 201 |
| PATCH | `/api/todos/:id` | `{"completed":true}` | 200 |
| DELETE | `/api/todos/:id` | — | 200 |

任务结构：`{"id":1,"title":"学习 Go","completed":false}`。
所有 JSON API 成功响应固定包含 `data`、`code: 0` 和 `message: ""`；失败响应仅包含非 0 的 `code` 和公开的 `message`，不包含 `data`。业务错误码不替代 HTTP 状态码。

- 列表：`{"data":{"items":[]},"code":0,"message":""}`，空列表始终为 `[]`。
- 单条任务：`{"data":{"id":1,"title":"学习 Go","completed":false},"code":0,"message":""}`。
- 删除：HTTP 200，`{"data":{},"code":0,"message":""}`。
- 健康检查：`{"data":{"status":"ok"},"code":0,"message":""}`。
- 错误示例：HTTP 404，`{"code":20001,"message":"任务不存在"}`。

| 错误码 | 含义 | HTTP 状态 |
|---|---|---|
| `10000` | 内部错误 | 500 |
| `10001` | 请求参数错误 | 400 |
| `10004` | 路由不存在 | 404 |
| `20001` | 任务不存在 | 404 |

通用协议与输出函数位于 `internal/pkg/http`（包名 `httpx`）；业务请求、路径参数和响应 DTO 位于 `internal/dto`，不直接返回 GORM 模型。前端统一通过 `web/src/lib/api.ts` 解包 `data`，同时识别 HTTP 状态及业务错误码。HTML、静态资源不包装此协议。
标题去除首尾空白后必须为 1–200 个字符。参数错误返回 400，记录不存在返回 404，内部错误返回 500 并记录日志。

## Swagger 文档

使用 Swaggo 注释生成 Swagger 2.0 文档（不是 OpenAPI 3），由 gin-swagger 提供在线查看与调试。

```bash
make swagger
DEBUG=true make dev-api
```

访问 http://localhost:8080/swagger/index.html；原始文档为 http://localhost:8080/swagger/doc.json。页面直接请求当前站点的 `/api`，不写死主机地址；Swagger UI 不经过 Vite 代理，请从后端端口访问。

接口注释位于 controller 和健康检查处理器，全局元信息位于 `cmd/server/main.go`。修改接口、DTO 或注释后执行 `make swagger`，将 `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml` 一并提交。生成命令固定使用 `swag v1.16.6`，普通构建和 Docker 构建直接使用已有生成文件，无需安装生成工具。

文档路由默认关闭，关闭时返回现有的 404 JSON 响应。生产环境保持关闭；如需开放，应通过网关限制访问。开关只控制路由注册，不提供认证。

## 目录结构

```text
.
├── cmd/server/           # 依赖组装、启动与优雅退出
├── docs/                 # Swaggo 生成的接口文档
├── internal/
│   ├── config/           # 环境变量配置
│   ├── controller/       # HTTP 入参、响应与错误映射
│   ├── db/               # SQLite 连接和自动迁移
│   ├── middleware/       # Zap 请求日志与 panic 恢复
│   ├── dto/              # 业务请求与响应
│   ├── pkg/http/         # 通用 JSON 响应协议
│   ├── model/            # 任务模型
│   ├── repository/       # GORM 数据访问
│   ├── router/           # API 路由与静态资源
│   └── service/          # 业务规则及数据访问接口
├── web/                  # React 前端、Vite 与 Bun 锁文件
├── docker/               # 前后端多阶段镜像构建
├── deploy/test/          # 测试环境 Compose
├── deploy/prod/          # 生产环境 Compose
├── .vscode/launch.json   # Go 与浏览器调试配置
├── Makefile
└── go.mod
```

调用关系：HTTP → controller → service → repository → SQLite。数据访问接口定义在使用方 service 中；未添加无用途的通用工具层。

## 请求日志

基础 Zap Logger 由 `internal/pkg/logger.New()` 显式初始化，进程结束时统一 `Sync()`。

中间件顺序为 `RequestID → Logging → Recovery`：RequestID 只负责生成随机 ID，通过 `internal/pkg/requestid` 写入标准请求 context，并设置 `X-Request-ID` 响应头；不依赖 Logger，也不沿用客户端传入的 ID。Logging 从 context 读取 request ID，派生带 `request_id` 的子 Logger 并写回请求 context，然后在请求结束时记录访问日志。

controller、service、repository 继续传递同一个 `context.Context`，需要日志时调用 `logger.FromContext(ctx)`；不向下传递 Gin Context。没有注入 Logger 的上下文返回无输出 Logger。内部错误在 controller 统一记录，repository 使用 oops 包装返回，避免重复日志。Recovery 的结构化日志也使用请求 Logger。

`db.WithContext(ctx)` 保留上下文，但不会自动将 GORM SQL 日志接入 Zap；当前不增加 SQL 日志适配。

## 容器部署

```bash
make docker-test
# 查看日志
docker compose -f deploy/test/compose.yaml logs -f
# 停止，保留数据库卷
docker compose -f deploy/test/compose.yaml down
```

镜像包含前后端产物，以非 root 用户运行。测试和生产使用各自命名的 SQLite 数据卷；不要用 `down -v` 停止需要保留数据的环境。

测试环境暴露 `8080`；生产环境仅绑定宿主机 `127.0.0.1:8080`，需要自行配置 HTTPS 反向代理。两个环境默认使用相同宿主机端口，不能同时启动。

本项目没有认证和权限功能，只用于最小模板演示；不要未经扩展直接作为公开的多用户应用部署。生产数据库结构演进时，应以版本化迁移替代简单的自动迁移。
