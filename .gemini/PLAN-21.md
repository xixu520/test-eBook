# PLAN-21: 系统前后端接口规范全量对齐与实现方案

基于 `PLAN-19` 发现的系统性脱节问题，以及 `PLAN-20`（已完成分类模块对齐）的经验，现针对剩余模块制定精确到代码层面的修改与实现方案。

## 整改目标
彻底消除前后端及《API 文档》之间的接口路径、参数和缺失功能的矛盾。以《API 文档》和合理的业务需求为最终基准。

---

## 模块修改明细表

### 1. 标准文件模块 (Documents)
**目标**：将所有文件管理路径对齐到 `/documents`，并补全回收站和历史记录功能。
**后端实现方案**：
- **Router (`router.go`)**: 
  - 删除原有的 `v1.Group("/standards")` 中涉及 `/files` 的部分。
  - 新增 `documents := v1.Group("/documents")`，注册：`GET ""` (列表), `GET "/:id"` (详情), `DELETE "/:id"` (软删除), `POST "/upload"`, `GET "/history"`。
- **Repository (`standard_repo.go`)**:
  - `DeleteFile` 改为依赖 GORM的软删除。
  - 新增 `GetFileHistory(number string)` 获取指定标准号的所有历史版本。
  - 新增 `GetRecycleBinFiles()` 和 `RestoreFiles(ids []uint)`, `HardDeleteFiles(ids []uint)`。
- **Handler (`standard.go`)**: 
  - 增加对应的历史记录和回收站的处理函数。
**前端修改方案**：
- **API (`api/document.ts`)**: 
  - 将 `/standards/files` 相关路径全部替换为 `/documents` 及 `/documents/:id`。对齐回收站相关的路径为 `/recycle-bin/documents` 等。

### 2. OCR 与上传模块 (Upload & OCR)
**目标**：对齐文档上传与 OCR 异步轮询的流程。
**后端实现方案**：
- **Router & Handler**:
  - `POST /documents/upload`：代替原 `/standards/files` 的上传，保存文件并返回 `task_id`（将其记录到 Redis 或 DB 的 task 表）。
  - 新增 `GET /tasks/:task_id/status`：实现轮询接口，实时去数据库或缓存中查询指定 OCR 任务进度。
  - 新增 `POST /documents/:id/ocr/retry`：提交后重新将其推入 OCR 解析队列。
- **Service (`standard_service.go`)**:
  - 重构 `UploadFile`，引入任务表支持异步进度管理；`ProcessFile` 要支持持久化写入任务状态（pending, processing, completed, failed）。
**前端修改方案**：
- **API (`api/upload.ts`)**: 
  - 废弃与 `document.ts` 冲突的配置。统一用 `api/document.ts` 主持上传。
  - 将 OCR 轮询 API 调整为 `/tasks/{task_id}/status`。

### 3. 数据大盘模块 (Dashboard)
**目标**：统一仪表盘的数据源路径。
**后端与前端均只做路径调整**：
- 后端 `router.go`：将 `v1.GET("/stats/dashboard", ...)` 变更为 `/admin/dashboard`。
- 前端 `api/stats.ts`：请求地址同步改为 `/admin/dashboard`。

### 4. 系统设置与日志模块 (Settings & Audit)
**目标**：从零补全缺失的管理端基础接口。
**后端实现方案**：
- **数据表与 Model**: 增加 `SystemSetting` 和 `AuditLog` GORM 模型。
- **Repository, Service & Handler**: 
  - Setting 模块：支持对系统的基础设置和 OCR Key 等进行增删改查；提供测试百度 OCR Key 有效性的验证逻辑 (`/settings/ocr-test`)。
  - Audit 模块：利用中间件捕获对关键 API (POST, PUT, DELETE) 的操作日志保存到了 `audit_logs` 表里；提供分页查询逻辑。
- **Router**: 注册 `/settings`, `/settings/ocr-test`, `/audit-logs`。
**前端修改方案**：
- **API**: 将 `api/settings.ts` 中 `POST /settings/test-ocr` 修正为 `/settings/ocr-test`，其他调用保持不变（当前已定义就绪）。

### 5. 用户管理与偏好 (User & Auth)
**目标**：补齐管理端的用户 CRUD 接口，并提供用户偏好的修改支持。
**后端实现方案**：
- **Handler & Service**: 
  - 增加用户管理 CRUD 功能 (针对前端发送至 `/admin/users` 的请求拉取分页列表和启停用/删数据)。
  - 在 `auth.go` 中增加 `UpdateTheme` 接口响应 `PUT /users/me/theme` 并写入 user 记录的 Theme 字段。
- **Router**: 在后台模块下注册 `/admin/users` 相关的对应路由。

---

## 下一步行动指南
由于修改范围较为庞大，为确保安全和避免冲突，建议按以下 **三个阶段 (Phase)** 依序开展：
1. [x] **第一阶段 (Phase 1)**：重构基础骨架与文件/回收站模块（已完成）。
2. [x] **第二阶段 (Phase 2)**：重构上传逻辑及接入真实的异步 OCR 轮询支持（已完成）。
3. [x] **第三阶段 (Phase 3)**：补齐系统支撑类功能：日志、设置、大盘、用户管理（已完成）。
4. [x] **补丁阶段 (Bug Fixes)**：修复 Nginx 代理斜杠、补全 OCR 任务列表、实现真实 AuthMe 接口（已完成）。

---

## 交付确认
- [x] 前后端接口路径全量对齐（`/categories`, `/documents`, `/admin` 等）。
- [x] 修复 Nginx 404 代理问题（移除 `/api/v1/` 末尾斜杠）。
- [x] 实现 AuthMe 角色校验逻辑，解决前端非法重定向。
- [x] OCR 任务异步化全流程打通，支持列表查询。
- [x] Docker 镜像构建优化（Node 20 + ESM 兼容）。

*(备注：本计划涉及的所有模块整改及线上 Bug 修复均已验证通过。)*
