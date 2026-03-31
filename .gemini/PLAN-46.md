# PLAN-46: OCR重构与PP-OCRv5模型切换

## 目标
根据 `/.gemini/OCR/OCR.MD` 技术文档，重构并实践（验证）OCR功能，并将模型名称正式修改为 `PP-OCRv5`。

## 当前状态
- [x] 后端 `paddle_client.go` 中针对错误状态码的捕获重构已在之前的会话中完成。
- [x] 切换 OCR 配置文件中配置模型为 `PP-OCRv5`
- [x] 测试后端的 OCR 连接验证接口
- [x] (可选) 实际发起一次上传或者 OCR 解析测试验证其正常可用。

## 具体变更记录
1. `test-ebook-api/config.yaml` 
   - 将 `paddle_model` 的值由 `"PaddleOCR-VL"` 修改为 `"PP-OCRv5"`。
2. 验证并实践
   - 测试调用测试连接的接口或尝试执行一次标准的文档上传 OCR 解析，验证配置正确。
