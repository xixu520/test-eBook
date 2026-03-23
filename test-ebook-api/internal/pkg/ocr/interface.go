package ocr

// Client OCR 客户端通用接口
type Client interface {
	// SubmitTask 提交 OCR 任务，返回任务 ID
	SubmitTask(filePath string) (string, error)
	
	// GetResult 获取 OCR 结果，返回内容、状态（success/processing/failed）和错误
	GetResult(taskID string) (string, string, error)
}
