## 目标

* 将当前仓库精简为“仅服务端代码”，移除所有 iOS 客户端相关文件与工程配置；保留并优化 Go 服务端的目录结构与启动方式，保证对外 API 稳定可用。

## 当前结构梳理

* 保留的服务端目录与文件：

  * `server/cmd/server/main.go` 入口，注册路由与中间件

  * `server/internal/*` 业务层、模型、仓储、数据库、鉴权、中间件、上传处理

  * `server/web/admin/*.html` 管理后台模板

  * `server/Dockerfile`、`server/deploy/render.yaml`、`server/go.mod`

* iOS 客户端相关（准备移除）：

  * `App/`、`Models/`、`Services/`、`Views/`、`Resources/`、`Tests/`

  * `PictureBook.xcodeproj/`、`PictureBook.xcworkspace/`

  * `Data/books.json`（仅供 iOS 使用，服务端已有内置 Sample/SQLite 种子）

## 拆除与精简

1. 删除所有 iOS 客户端源代码与资源：`App/`、`Models/`、`Services/`、`Views/`、`Resources/`、`Tests/`
2. 删除 iOS 工程文件：`PictureBook.xcodeproj/` 与 `PictureBook.xcworkspace/`
3. 删除 `Data/books.json`（服务端不依赖该文件；服务端通过 `repo.SeedIfEmpty(sampleBooks())` 自动填充）
4. 保留 `.trae/documents`（若你不需要这类计划文档，可一并移除；默认为保留）

## 目录重构（可选）

* 方案 A（保守）：继续以 `server/` 作为服务端根目录，仓库根仅保留 `server/` 与必要的元文件（如 CI 配置）。

* 方案 B（合并）：将 `server/` 平铺到仓库根，形成经典 Go 项目结构：

  * `cmd/server/main.go`

  * `internal/...`

  * `web/admin/...`

  * `Dockerfile`、`go.mod`

* 两方案均不影响功能；B 方案更贴近常见 Go 项目约定。若选择 B，需要同步调整 `go.mod` 模块路径与 Dockerfile COPY 路径。

## 对外 API 一览（保持不变）

* 健康检查：`GET /healthz`

* 公共接口（用于客户端调用）：

  * `GET /v1/editor-picks`

  * `GET /v1/books?page&page_size&age&sort`

  * `GET /v1/books/:id`

  * `GET /v1/books/:id/recommendations?age&limit`

* 管理与上传接口：

  * `POST /v1/admin/books`（创建书籍）

  * `POST /v1/admin/upload`（通用上传）

  * `POST /v1/admin/books/:id/pages/upload`（批量上传页图）

  * `POST /admin/books/:id/pages/reorder`（页面排序，带鉴权）

  * 管理后台页面路由：`/admin/*`（登录、列表、编辑）

## 运行与部署保持一致

* 环境变量：`CORS_ALLOWED_ORIGINS`（默认 `*`）

* 静态资源：`/assets` 映射至 `uploads/` 目录

* 数据库：`internal/db` 自动初始化 SQLite schema；首次启动空库将使用 `sampleBooks()` 种子

* Docker：沿用 `server/Dockerfile`（Go 构建 → 复制 `web` → 暴露 8080）

* Render 部署：沿用 `server/deploy/render.yaml`（`healthCheckPath: /healthz`）

## 质量与自动化（建议项）

* 增加基础集成测试（Gin 路由层的 HTTP 端到端测试）

* 简单 Makefile：`make build`、`make run`、`make docker-build`

* 可选 CI：Push 时触发 `go test` 与 Docker 构建

## 变更影响与兼容性

* 移除 iOS 客户端后，仓库定位为“纯服务端”；其他客户端（iOS/Android/Web/小程序）均以 `https://<server>/v1/...` 调用。

* 不改动 API 合同与响应结构，现有接口继续可用。

* 若你选择目录平铺（方案 B），会微调 `go.mod` 与 Dockerfile 路径；功能不受影响。

## 交付步骤

1. 执行删除与精简（移除 iOS 相关目录与工程文件）
   2.（可选）将 `server/` 平铺到根并修正 `go.mod`、`Dockerfile` 路径
2. 本地验证：`go build ./server/cmd/server` 或 `docker build -f server/Dockerfile .`
3. 运行检查：本地启动后访问 `GET /healthz`、`GET /v1/editor-picks` 等确认正常

## 需要你确认的选项

* 目录重构选择：保留 `server/` 作为子目录（方案 A）还是平铺到根（方案 B）？

* `.trae/documents` 是否保留？如果不需要，我会一并删除。

执行方案A，这个样就不用更改配置文件
