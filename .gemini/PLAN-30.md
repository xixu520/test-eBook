# PLAN-30: 启动前后端进行调试

## 目标
启动后端 API 服务和前端 Vue 应用，以便 xixu520 进行调试。

## 计划步骤

### 1. 启动后端 (test-ebook-api)
- [x] 进入后端目录 `test-ebook-api`
- [x] 运行 `go run cmd/server/main.go`
- [x] 确认运行在端口 `8182`

### 2. 启动前端 (test-ebook-web)
- [x] 进入前端目录 `test-ebook-web`
- [x] 运行 `npm run dev`
- [x] 确认运行并在浏览器中访问

## 确认事项
- [x] 后端服务启动并连接数据库
- [x] 前端应用启动并连接后端 API

## 备注
- 如果端口冲突，请告知 xixu520。
- 前端 API 基础路径在 `.env.development` 中配置。
