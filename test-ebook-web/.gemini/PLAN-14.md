# PLAN-14: 前端功能对接开发 (Frontend Integration)

## 目标
将前端的“文档管理”与“分类管理”模块正式接入后端的真实 API（PLAN-12/13 成果），实现标准文件的真实上传、层级展示及 OCR 解析状态同步。

## 状态
- [ ] 1. 更新 `src/utils/request.ts`，确保 `baseURL` 正确指向 `http://localhost:8181/api/v1`。
- [ ] 2. 重构 `src/api/document.ts` 与 `src/api/category.ts`，适配 `/standards` 路径及新的参数结构。
- [ ] 3. 适配 `AdminDocumentsPage` (通常位于 `src/views/admin/document/index.vue`)：
    - [ ] 接入真实的分页查询。
    - [ ] 字段映射补全：`standard_no` -> `number`, `name` -> `title`。
    - [ ] 状态逻辑映射：UI `ocr_status === 'completed'` -> Backend `Status === 1`。
    - [ ] 实现真实的文件上传（使用 `FormData` 提交至 `/standards/files`）。
    - [ ] 动态显示 OCR 解析进度 (Status: 0, 1, 2)。
勾选
- [ ] 3. 适配 `AdminDocumentsPage` (通常位于 `src/views/admin/document/index.vue`)：
    - [ ] 接入真实的分页查询。
    - [ ] 字段映射补全：`standard_no` -> `number`, `name` -> `title`。
    - [ ] 状态逻辑映射：UI `ocr_status === 'completed'` -> Backend `Status === 1`。
    - [ ] 实现真实的文件上传（使用 `FormData` 提交至 `/standards/files`）。
    - [ ] 动态显示 OCR 解析进度 (Status: 0, 1, 2)。
- [ ] 4. 适配 `CategoryManagement`：
    - [ ] 接入递归分类树展示。
    - [ ] 实现新增分类功能。
- [ ] 5. 浏览器验证：完成上传 -> 轮询 -> 结果显示的完整链路。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/api/document.ts` | 后端接口适配 |
| `src/api/category.ts` | 分类接口适配 |
| `src/views/admin/document/index.vue` | 文档管理页面数据驱动更新 |
| `src/views/admin/category/index.vue` | 分类管理页面数据驱动更新 |
