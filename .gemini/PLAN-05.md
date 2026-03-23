# PLAN-05: 后台管理增强（文档库与回收站）

## 目标
完善系统后台管理功能，提供专门的文档库管理页面、回收站功能以及审计日志，实现更专业的内容管理能力。

## 状态
- [x] 注册 `/admin/documents`, `/admin/recycle`, `/admin/audit-logs` 路由
- [x] 实现 `AdminDocumentsPage.vue`（含上传者/下载统计等）
- [x] 实现 `RecycleBinPage.vue`（支持还原与彻底删除）
- [x] 实现 `AuditLogPage.vue`（操作日志展示）
- [x] 优化 `AdminLayout` 侧边栏菜单跳转
- [x] 完善 Mock 接口数据支持

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/document/AdminDocumentsPage.vue` | 后台文档管理主页 |
| `src/views/admin/recycle/RecycleBinPage.vue` | 回收站页面 |
| `src/views/admin/audit/AuditLogPage.vue` | 审计日志页面 |
| `.gemini/PLAN-05.md` | 本计划文件 |

## 设计依据
- `test eBook - API.md` — 模块 6 (回收站) 与 模块 7 (审计日志)
- `前端架构书.md` — 5.6, 5.8 章节建议
