# PLAN-25: 重建分类管理模块

你好 xixu520。本计划旨在根据后端 API 实际接口结构，重建前端分类管理页面。

## 第一轮：初次重建（已完成）
- [x] `src/api/category.ts` - 增加 TypeScript 接口定义
- [x] `src/views/admin/category/CategoryPage.vue` - 初次重写

## 第二轮：交互修复（已完成）
- [x] 改用 `el-tree` 树状展示，删除行内"添加子类"按钮

## 第三轮：核心 BUG 修复 + 视觉强化（已完成）
- [x] `request.ts` 响应拦截器 `res.message` → `res.msg` 对齐后端字段
- [x] `category.ts` 接口 `id` → `ID` 对齐 GORM `gorm.Model` 大写字段
- [x] `CategoryPage.vue` 6 处 `row.id` → `row.ID` 修正（node-key/props/filter/edit/delete）
- [x] 树结构连接线样式（竖线 + 横线 + 渐变 + 悬停高亮）
- [x] 层级字体分级（L1粗/L2中/L3正常）、层级标签（L2/L3）

## 完成情况
- [x] 删除空分类 → 成功
- [x] 删除有文档分类 → 正确显示"该分类下有关联文件，无法删除"
- [x] 添加子分类 → 父级选择器正常
- [x] 树连接线 → 有子节点时显示
