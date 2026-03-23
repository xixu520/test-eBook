# PLAN-10: 适配 PaddleOCR 官方异步 API 参数

## 目标
根据官方 Python 异步解析脚本分析，重构目前预留的 PaddleOCR 官方接口设置参数形式。从原本假定的 API/Secret Key 形式转换成为原生符合要求的单 `Token` 形式，并暴露官方特有的配置项模型设置与图像优化设置。

## 状态
- [x] 1. 将鉴权参数由 `API Key` + `Secret Key` 简化修改为 `Token` 认证。
- [x] 2. 增加 `Model` 字段配置项（默认：`PaddleOCR-VL-1.5`）。
- [x] 3. 增加布尔值高级可选参数（文档方向分类、文档去畸变、图表识别）作为可选开关选项。
- [x] 4. 同步修正 `mock/settings.ts` 数据结构和测试鉴权规则（已增加对特定合法 Token 的等值校验）。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/settings/SettingsPage.vue` | 移除旧字段，引入 Token 及 OCR 选型表单 |
| `mock/settings.ts` | 更新保存和测试接口 Mock |
