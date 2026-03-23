# PLAN-08: 完善后台 OCR 引擎设置

## 目标
检查并完善后台管理系统中的“系统设置”页面，特别是围绕“OCR引擎”与“百度云 OCR”的配置项，保证设置界面的全功能覆盖与 Mock 数据支撑。

## 状态
- [x] 1. 审查 `SettingsPage.vue`，评估现有 OCR 表单是否覆盖百度云的必要配置（如 API Key, Secret Key 等）。
- [x] 2. 完善 `SettingsPage.vue` 的界面与交互（密码遮盖、连接测试等）。
- [x] 3. 补充或完善 `src/api/settings.ts` 以及 `mock/settings.ts` 的接口，打通前端状态流转。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/settings/SettingsPage.vue` | 设置界面核心代码 |
| `src/api/settings.ts` | 接口调用封装 |
| `mock/settings.ts` | 模拟后端接口返回 |
