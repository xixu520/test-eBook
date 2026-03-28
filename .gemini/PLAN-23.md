# PLAN-23: 登录页面快捷键功能优化

你好 xixu520。本计划旨在为登录页面添加键盘快捷键支持，提升用户操作体验。

## 目标
- [x] 在登录页面的用户名和密码输入框添加回车键（Enter）监听，触发登录逻辑。

## 变更记录
- [x] 修改 `src/views/login/LoginPage.vue`，为 `el-input` 添加 `@keyup.enter` 事件。

## 完成情况
- [x] 登录页面 Enter 键快捷登录功能已实现

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `test-ebook-web/src/views/login/LoginPage.vue` | 登录页面组件 |
