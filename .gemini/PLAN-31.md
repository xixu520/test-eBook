# PLAN-31: 后台文档管理列表增加属性修改功能

## 目标
检查后台管理/文档管理界面（`AdminDocumentsPage.vue`），在操作栏为每条文档记录增加"编辑"功能，支持动态修改文档关联的核心属性。**仅创建计划，不直接修改代码。**

## 涉及页面及组件
- **前端页面**：`src/views/admin/document/AdminDocumentsPage.vue`
- **前端 API**：`src/api/document.ts`
- **组件库**：Element Plus

## 计划步骤

### 1. 后端接口开发 (test-ebook-api)
- [x] **路由层**：在 `internal/router/router.go` 的 `documents` 路由组中，新增 `PUT /:id` 路由配置。
- [x] **Handler 层**：在 `internal/handler/standard.go` 中，新增 `UpdateFile` 方法，接收修改参数并做校验。
- [x] **Service 层**：在 `internal/service/standard_service.go` 中，新增 `UpdateFile(id, payload)` 逻辑。
- [x] **Repo 层**：在 `internal/repository/standard_repo.go` 中，实现 `Updates` 方法，更新指定的 `StandardFile` 记录字段。

### 2. 前端 API 对接 (test-ebook-web)
- [x] **API 封装**：在 `src/api/document.ts` 中新增 `updateDocument(id: number, data: any)` 请求方法，对应 `PUT /api/v1/documents/${id}`。

### 3. 前端 UI 改造 (test-ebook-web)
- [x] **操作按钮**：在 `AdminDocumentsPage.vue` 的 `<el-table-column label="操作">` 内，添加"编辑"的 `<el-button>`。
- [x] **编辑弹窗组件**：实现一个包含 `<el-form>` 的 Dialog 弹窗。
    - 表单字段包括：名称 (Title)、标准号 (Number)、版本 (Version)、发布机构 (Publisher)、发布日期 (IssueDate)、实施状态 (ImplementationStatus)、所属分类 (CategoryID)。
    - 根据提取到的 Element Plus 设计风格，使用漂亮的交互表单（带有验证规则）。
- [x] **功能联调**：点击编辑时完成数据回显；点击保存时请求 `updateDocument` 并重新加载 `loadData()`。

## 备注
- 遵循统一的 Element Plus 样式和现有的接口规范。
- 确保有适当的鉴权 (Token 需要包含在请求头，由于已有拦截器通常不用特殊处理)。
- 暂未开始代码修改，等待指令进行开发。
