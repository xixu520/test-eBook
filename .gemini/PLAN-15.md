# PLAN-15: 前后端 Docker 容器化与部署实践

你好 xixu520。本计划旨在将项目全面容器化，实现一键启动与环境一致性。

## 目标
- [ ] 为 Go 后端 `test-ebook-api` 编写多阶段构建 Dockerfile
- [ ] 为 Vue 前端 `test-ebook-web` 编写 Nginx 托管 Dockerfile
- [ ] 编写 `docker-compose.yml` 实现服务编排
- [ ] 处理容器间的网络通信与环境变量注入
- [ ] 验证容器化部署后的全链路功能（包括 OCR 与文件上传）

## 详细步骤

### 1. 后端容器化 (test-ebook-api)
- **Dockerfile**: 使用 `golang:1.21-alpine` 作为构建环境，`alpine:latest` 作为运行环境。
- **配置处理**: 将 `config.yaml` 映射为容器卷或内置默认配置。
- **存储映射**: 将 `uploads/` 和 `data/` (SQLite 数据文件) 映射为宿主机持久化卷。

### 2. 前端容器化 (test-ebook-web)
- **Dockerfile**: 使用 `node:18-alpine` 进行构建，`nginx:alpine` 提供 Web 服务。
- **Nginx 配置**: 处理 SPA 路由（history mode）及反向代理（可选，或由独立网关处理）。
- **BaseURL 注入**: 通过构建参数或运行时环境替换方案解决 `VITE_API_BASE_URL` 的动态指向。

### 3. 服务编排 (docker-compose)
- 定义 `backend` 服务：暴露端口，挂载卷。
- 定义 `frontend` 服务：暴露 3000 或 80 端口，链接至 backend。
- 配置网络：确保 frontend 容器或浏览器能访问至 backend。

### 4. 部署验证
- 执行 `docker-compose up --build`。
- 测试登录、上传、OCR 轮询及分类管理。

## 注意事项
- [!IMPORTANT]
- 需要处理 SQLite 数据库文件权限问题。
- 前端 `VITE_API_BASE_URL` 如果在浏览器端运行，需指向宿主机 IP 或域名，而非容器局域网 IP。
