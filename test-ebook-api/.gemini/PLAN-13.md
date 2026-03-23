# PLAN-13: 真实 OCR 接口集成开发

## 目标
将现有的 Mock 异步处理逻辑替换为真实的 OCR 接口调用（优先集成 PaddleOCR官方 API），实现建筑标准文件的文字内容自动提取。

## 状态
- [x] 1. 完善 `config.yaml` 与 `config.go`，增加 PaddleOCR (AI Studio) 的鉴权配置。
- [x] 2. 在 `internal/pkg/ocr` 下实现 `PaddleClient`，支持任务提交与结果轮询。
- [x] 3. 修改 `StandardService.ProcessFile`，接入真实的 OCR 调用链。
- [x] 4. 实现错误重试机制与状态更新（如：解析中、解析成功、解析失败）。
- [ ] 5. 验证真实 PDF/图片文件的文字提取效果及存库准确性。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `config.yaml` | 存放 OCR Token 等敏感信息 |
| `internal/config/config.go` | 映射 OCR 配置结构 |
| `internal/pkg/ocr/paddle.go` | PaddleOCR API 客户端实现 |
| `internal/service/standard_service.go` | 更新异步处理逻辑 |
