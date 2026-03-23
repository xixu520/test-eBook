# PLAN-06: 核心业务逻辑修复与安全加固

## 目标
解决在前后端架构审查中发现的 4 个核心业务逻辑与权限安全防范问题。

## 状态
- [x] 修复 1.1：(分类删除) 拦截后端 400 错误码并给予分类删除的针对性提示（在实现 CategoryPage 时应用与 Mock）。
- [x] 修复 1.2：(回收站) 严格限制 `empty_all: true` 的触发场景，避免批量删除误触。
- [x] 修复 1.3：(权限冗余) 剥离 `HomePage.vue` 中的高危管理权限（删除、编辑、重新 OCR），交由后台统一闭环。
- [x] 修复 2.1：(OCR 状态隔离) 文件若处于 `pending/processing` 阶段，前端需强禁用“预览”及“历史版本”功能，防止查空或崩溃。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/home/HomePage.vue` | 移除管理按钮，增加 OCR 状态前端校验机制 |
| `src/views/admin/document/AdminDocumentsPage.vue` | 增加 OCR 状态关联拦截 |
| `src/views/admin/recycle/RecycleBinPage.vue` | 审查清空与批量删除的传参逻辑 |
| `mock/document.ts` | 增强对删除清空的 Mock 拦截验证 |
| `mock/category.ts` | 增加带有子分类或文档关联时的模拟拦截 |
