# PLAN-11: 初始化后端 Go 工程架构

## 目标
根据《后端架构书.md》初始化 Go 后端项目 `test-ebook-api`，搭建核心分层目录结构，配置 Gin 框架基础路由与中间件，并完成 SQLite 数据库与 GORM 的初步集成。

## 状态
- [x] 1. 创建项目目录并执行 `go mod init test-ebook-api`。
- [x] 2. 搭建分层目录结构（cmd, internal/handler, service, repo, model...）。
- [x] 3. 初始化核心配置加载逻辑 (`internal/config`) 与统一响应工具 (`internal/pkg/response`)。
- [x] 4. 初始化数据库连接与 GORM 集成。
- [x] 5. 实现基础路由骨架与中间件 (CORS, Logger, JWT)。
- [x] 6. 完成用户登录验证与 Seed 管理员账号机制。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `test-ebook-api/go.mod` | 依赖管理 |
| `test-ebook-api/cmd/server/main.go` | 工程入口 |
| `test-ebook-api/internal/config/config.go` | 配置管理 |
| `test-ebook-api/internal/database/sqlite.go` | 数据库驱动配置 |
| `test-ebook-api/internal/router/router.go` | 路由分发 |
