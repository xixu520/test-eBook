# PLAN-32: 审计日志系统重构计划书

## 背景
现有审计日志系统存在严重的前后端逻辑断层与潜在的性能隐患（例如全局读取文件上传等流数据的 Body 导致内存溢出）。目前列表查询条件未接通，记录的操作类型（Action）与前端期待的简写标签脱节，且缺少关键的“登录(LOGIN)”日志记录。
此计划书已完整执行。系统重构已完成。**（已于 2026-03-30 执行完毕）**

## 缺陷与现状分析

### 1. 前端期待 vs 后端实际记录断层
- **UI 期待**：`AuditLogPage.vue` 的 `filters.action` 期待的是 `LOGIN`, `UPLOAD`, `DELETE`, `EDIT`, `VERIFY`。
- **后端现状**：`middleware/audit.go` 中的 `AuditMiddleware` 暴力将 `Action` 赋值为 `c.Request.Method + " " + c.Request.URL.Path` (例如 `POST /api/v1/documents/upload`)。
- **结果**：前端下拉筛选失效，标签变色匹配失败。

### 2. 查询入参被抛弃
- **UI 逻辑**：前端在请求 `getAuditLogs` 时传入了 `action` 查询参数。
- **后端现状**：`internal/handler/audit.go` 下的 `GetAuditLogs` **根本没有接收 `action` 参数**，Service 和 Repo 同样没有开放此条件过滤。

### 3. 未屏蔽 Multipart / 流媒体带来的内存泄漏
- **后端现状**：所有 POST/PUT 请求均会通过 `io.ReadAll(c.Request.Body)` 缓冲 Body。一旦遇到大文件上传，会导致内存占满（OOM）引发程序崩溃。

### 4. 无法抓取 Login 登录日志
- **后端现状**：登录路由 `POST /api/v1/auth/login` 属于公开路由，位于 `protected.Use(middleware.AuditMiddleware)` 保护组之外，所以没有任何登录日志。

---

## 详细功能级重构指南 (精确到函数)

### 阶段一：修复安全隐患与日志生成定义

#### [1] `internal/middleware/audit.go` -> `AuditMiddleware(db *gorm.DB)` 函数
- **修改内容**：
  1. 引入 Content-Type 检查：`contentType := c.Request.Header.Get("Content-Type")`，遇到 `multipart/form-data` 时，**跳过** `io.ReadAll(c.Request.Body)` 的读取，强行将 `Details` 设置为 `"文件上传请求"`。
  2. 新编写一个独立函数 `mapActionName(method, path string) string`，使用正则或前缀匹配：
     - 如果是 `/auth/login` ->返回 `LOGIN`
     - 如果包含 `/documents/upload` ->返回 `UPLOAD`
     - 如果碰到 `DELETE` ->返回 `DELETE`
     - 如果碰到 `PUT` ->返回 `EDIT`
     - 其余保持原样或标为 `OTHER`。将该函数的返回值赋给 `AuditLog.Action`。

#### [2] `internal/handler/auth.go` -> `Login(c *gin.Context)` 函数
- **修改内容**：
  - 登录接口是公开的，无法走全局 Audit 中间件。在返回 Token `pkg.Success(c, ...)` 的同时，**手动向 `model.AuditLog` 表内 Insert 一条记录**。指定 Action 为 `"LOGIN"`，Username 为请求用户的名称。

### 阶段二：打通全链路过滤查询机制

#### [1] `internal/handler/audit.go` -> `GetAuditLogs(c *gin.Context)` 函数
- **修改内容**：
  - 增加：`action := c.Query("action")`。
  - 更改参数传递：调用 `h.svc.GetLogs(page, pageSize, action)`。

#### [2] `internal/service/audit_service.go` -> `GetLogs(page, pageSize int, action string)` 函数
- **修改内容**：
  - 更新函数签名，接收 `action`。
  - 传递给 `s.repo.GetLogs(page, pageSize, action)`。

#### [3] `internal/repository/audit_repo.go` -> `GetLogs(page, pageSize int, action string)` 函数
- **修改内容**：
  - 函数签名增加 `action string`。
  - `db := r.db.Model(&model.AuditLog{})` 之后，补充 IF 判断：
    ```go
    if action != "" {
       db = db.Where("action = ?", action)
    }
    ```

### 阶段三：前端调优校对

#### [1] `src/api/audit.ts` -> `getAuditLogs`
- **检查内容**：确认是否支持传入 `action` 作为 params 供后端使用（目前已支持，但应严谨核对类型）。

#### [2] `src/views/admin/audit/AuditLogPage.vue`
- **检查与修改内容**：
  - 原逻辑 `getActionType` 拥有 `UPLOAD/DELETE/VERIFY/EDIT/LOGIN` 的染色能力。确保后端传输的正是这些魔法字符串，无需修改样式逻辑。
  - 需要在 `el-table-column prop="details"` 这个列中增加对过长 JSON 字符串的处理方案（如 `<el-tooltip>` 或截断），防止因上传长参数引起布局撑破。

## 结论
按此计划书操作，项目将完美接通审计日志前后端查询功能，并消除原有的文件 Body 读取崩溃隐患。已将此方案交接给 Gemini Flash 进行落地修改。
