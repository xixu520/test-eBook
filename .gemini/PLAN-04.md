# PLAN-04: 前端多端响应式优化

## 目标
实现“建筑标准文件管理系统”前端页面的多端（手机、平板、桌面）适配，确保在不同屏幕尺寸下均有出色的用户体验。

## 状态
- [x] 定义全局响应式断点与 Mixins
- [x] 优化 `TopNavBar`（移动端隐藏标题、搜索框适配）
- [x] 优化 `SideMenu`（移动端改为 Drawer 模式，平板端自动折叠）
- [x] 优化 `MainLayout` / `AdminLayout` 布局逻辑
- [x] 优化 `HomePage`（表格滚动、筛选区折叠）
- [x] 优化 `LoginPage` 卡片适配
- [x] 全设备模拟验证

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `test-ebook-web/src/styles/variables.scss` | 全局样式变量与 Mixins |
| `test-ebook-web/src/components/layout/*` | 布局相关组件 |
| `test-ebook-web/src/views/*` | 主要页面组件 |
| `.gemini/PLAN-04.md` | 本计划文件 |
