# PLAN-29: 前端首页全面修复

## 目标
修复前端首页的全部严重问题，包括表格字段不匹配、搜索断裂、上传表单错误、侧边栏崩溃等。

## 计划步骤

### P0 — 严重问题
- [x] #1 表格字段名映射修复 (HomePage.vue prop 对齐后端 StandardFile 模型)
- [x] #4 侧边栏 Category ID 大小写 + 树构建修复
- [x] #3 上传表单字段对齐后端
- [x] #2 搜索查询参数修复 (前后端对接)
- [x] #10 分页参数名修复

### P1 — 中等问题
- [x] #5/#6 下载/预览功能对接 (在后端立即实现 /documents/:id/download 和 /documents/:id/preview 端点)
- [x] #7 分类数据去重 + 统一管理
- [x] #9 移除"导出 Excel"和"列配置"按钮
- [x] #12/#13 清理脚手架残留 CSS 和 HelloWorld.vue

### P2 — 低级别
- [x] #15 适配发布机构字段 (配合后端新加字段)
- [x] #17 公告栏 API 类型修复

### 新增后端字段任务
- [x] 后端 StandardFile 模型增加字段：publisher(发布机构), issue_date(发布日期), implementation_status(实施状态), verify_status(核验状态)
- [x] 后端 ListFiles 支持 keyword, publisher, implementation_status过滤
