# PLAN-09: 完善 PaddleOCR 相关功能设置

## 目标
补全针对 PaddleOCR 引擎的专属后台配置项，完善其联调与配置交互闭环。

## 状态
- [x] 1. 在 `SettingsPage.vue` 添加针对 PaddleOCR 引擎特有的官方 API 配置（API Key 和 Secret Key）。
- [x] 2. 在原有测试连接框架中，打通 PaddleOCR 官方接口配置的存活性联调测试按钮。
- [x] 3. 补齐 `mock/settings.ts` 数据结构，使得前端可以正常回显并校验模拟的 PaddleOCR 官方调用入参。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/settings/SettingsPage.vue` | 增加 PaddleOCR 专属配置输入项与验证 |
| `mock/settings.ts` | 添加 PaddleOCR `server_url` 等相关数据 Mock 支持 |
