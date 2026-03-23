# 建筑标准文件管理系统 (test-eBook)

一个基于现代化前后端分离架构的建筑标准文件（PDF）管理系统。提供标准文件的分类管理、安全上传、全员预览、下载控制以及核心特性的**自动化 OCR 识别与标准信息提取**等功能。

## ✨ 核心特性

### 📂 大规模文件与分类管理
- **多层级分类树**：灵活自定义标准文件目录树，实现按细分业务场景组织。
- **文件全生命周期**：支持高速上传、属性编辑、无缝预览（内置 PDF.js）、授权下载。具有软删除（回收站机制）、彻底删除与智能检索机制。
- **批量管控支持**：支持批量更改分类、批量安全抹除等效率操作。

### 🤖 智能 OCR 扫描提取
- **自动解析联动**：借助 `pdfcpu` 截取长 PDF 首页进行精准分离。
- **云端 OCR 接入**：串联百度云高精度 OCR 接口，全自动识别大批上传的标准文件，准确分离并提取：标准号、名称及发布日期。
- **全异步轮询架构**：后端通过 Goroutine 与 gocron 执行后台兜底与防阻碍处理，前端即时轮询查收识别结果；同时提供 OCR 手动重试挂载点。

### 👨‍💼 严谨的后台及安全体系
- **角色体系与细粒度权限**：划设 Admin 和 User，实施下载限制、操作审计与编辑控制。
- **配置化面板**：具备站点动态配置注入功能，支持百度网盘 OCR Key 与 Token 在控制台进行热更新及调试测试。
- **审计数据链路**：高度敏感操作（文件移除、属性伪装、登录频次）强制登记溯源；自带多维度仪表盘进行运营统计。

---

## 🛠 现代化技术选型

### 前端开发 (test-ebook-web)
- **核心构建**：Vue 3 + Vite + TypeScript (多阶段并行构建优化)
- **基座周边**：Pinia + Vue Router
- **UI 生态系**：Element Plus 搭配深度自定义暗黑/明亮主题切换
- **文件与构建处理**：PDF.js (内嵌预览)、Gzip 预打包、手动 Chunk 分离与资源加速

### 后端开发 (test-ebook-api)
- **性能基座**：Go (1.21+) / Gin Web Framework (带有完备的优雅停机扩展)
- **数据库范式**：纯 Go SQLite (`glebarez/sqlite`) + GORM（开启 WAL，优化单写多读并发上限）
- **安全与环境隔离**：JWT (含 bcrypt 盐值认证)、Viper 环变配置覆盖机制
- **工具链集成**：Zap 结构日志, gocron (超时扫描任务)

---

## 🚀 生产环境部署 (Docker)

本项目全架构经由 Docker 大幅压缩、整合多阶段构建逻辑生成。部署操作异常轻量：

### 1. 基础配置预制
确认服务器已具备 `Docker` 以及 `Docker Compose` 环境。推荐调整及核验当前目录结构内的配置或环境变量设置。如设定特定的 JWT 令牌或 OCR 通道。

### 2. 执行 Docker Compose
系统内自带 `docker-compose.yml` 组装结构，含有一个 Node+Nginx (处理前端) 容器及一个 Go-Alpine (处理后端 API) 高性能容器环境：

```bash
docker-compose up -d --build
```

### 3. 数据层持久化管理说明
在启动后，将在根目录下挂载自动生成的核心数据存储文件夹：
- `backend-data:/app/data`：存储了极其重要且需日常化备份的 SQLite 核心库及 WAL 日志。
- `backend-uploads:/app/uploads`：存放用户上传的所有原生标准 PDF 文件。

---

## 💻 本地开发协作指南

### 启动后端 (API Server)
```bash
cd test-ebook-api
go mod download
go run cmd/server/main.go
```
*API 服务将通过 Gin 在本地初始化，并在 `test-ebook-api/data` 夹下初始化持久化 SQLite 缓存。*

### 启动前端 (Web Client)
```bash
cd test-ebook-web
npm install
npm run dev
```
*开发服务器将在 `3000` 端口开启，请求会自动根据 `vite.config` 转发跨域至本地后端代理，即改即用。*