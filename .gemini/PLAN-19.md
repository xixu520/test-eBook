# PLAN-19: 前后端功能与数据一致性检查

## 目标
检查前后端各项功能，发现并记录互相矛盾的地方（如 API 路径不一致，请求参数、响应数据格式不一致，或是功能逻辑冲突）。

## 发现的问题清单
通过对比发现，`test eBook - API.md` 文档、前端请求定义 (`test-ebook-web/src/api/*.ts`) 以及后端路由实现 (`test-ebook-api/internal/router/router.go`) 之间存在大量不一致或遗漏的地方：

### 1. 分类管理模块 (Category)
- **前端定义**: 
  - `GET /standards/categories`
  - `POST /standards/categories`
  - `DELETE /standards/categories/:id`
- **后端实现**:
  - `GET /standards/categories` (不符合API文档的 `/categories`)
  - `POST /standards/categories`
  - 缺失 `DELETE` 路由，但在文件模块里错误地写成了 `standards.DELETE("/files/:id", standardHandler.DeleteFile)`
- **API文档**: 规定为 `/categories`，并有 `PUT` 接口但在前后端均未实现完全。

### 2. 标准文件模块 (Documents/Files)
- **前端定义**: `GET /standards/files`, `POST /standards/files`, `GET /standards/files/:id`, `DELETE /standards/files/:id`
- **后端实现**: 与前端调用一致，均为 `/standards/files` 相关。
- **API文档**: 规定路径为 `/documents`，上传为 `/documents/upload`。
- **缺失项**: 前端调用了 `GET /documents/history`, `GET /recycle-bin/documents`, `PUT /recycle-bin/documents/restore`, `POST /recycle-bin/documents/batch-delete`，但后端路由 **完全未实现** 历史记录和回收站模块。

### 3. 上传与 OCR 模块
- **前端定义**: `api/upload.ts` 定义了 `POST /upload`, `GET /ocr/tasks`，但 `document.ts` 中也定义了 `POST /standards/files` 作为上传。
- **后端实现**: 仅实现了 `POST /standards/files`。
- **API文档**: 规定上传为 `/documents/upload`，获取状态为 `/tasks/{task_id}/status`，重试为 `/documents/{id}/ocr/retry`。前后端与文档均不一致，且后端未实现 OCR 任务状态查询。

### 4. 日志、设置、用户管理 (后台功能)
- **前端定义**: 
  - `GET /admin/users` 及相关增删改
  - `GET /audit-logs`
  - `GET /settings`, `POST /settings`, `GET /system/status`, `POST /settings/test-ocr`
- **后端实现**: **完全没有** 实现上述任何路由。
- **API文档**: 部分存在（如设置和日志），但用户管理 `/admin/users` 在文档中也未提及。

### 5. 看板数据 (Dashboard)
- **前端与后端一致**: 使用 `/stats/dashboard`。
- **API文档**: 规定为 `/admin/dashboard`。

### 6. 用户认证与个人信息 (Auth)
- **前端调用**: `PUT /users/me/theme`
- **后端实现**: 缺失该更新主题的路由。

## 总结与建议
当前项目的前后端 API 完全脱节。
1. **重新对齐 API 文档**：建议先对 `test eBook - API.md` 进行修订和确认，确定最终的接口规范。
2. **重构后端路由**：按照规范重写 `test-ebook-api/internal/router/router.go` 及相关 Handler，补齐缺失的路由（特别是回收站、OCR查询、系统设置、审查日志、用户管理等）。
3. **修复前端请求**：更新 `src/api` 下的文件，让 URL 路径与文档保持绝对一致，消除重复或矛盾的函数。
