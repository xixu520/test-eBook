package queue

import (
	"errors"
	"sync"
)

// TaskType 任务类型
type TaskType string

const (
	TaskUploadSync TaskType = "upload_sync" // 文件从 staging 同步到最终存储
	TaskOCR        TaskType = "ocr"         // OCR 识别任务
)

// Task 队列中的任务单元
type Task struct {
	Type       TaskType
	DocumentID uint
	UploadID   uint   // UploadTask.ID（仅 upload_sync 类型使用）
	TaskID     string // OCRTask.TaskID（仅 ocr 类型使用）
	RetryCount int
}

// Queue 基于 Go channel 的内存任务队列
type Queue struct {
	ch     chan Task
	closed bool
	mu     sync.Mutex
	name   string
}

// NewQueue 创建指定缓冲大小的队列
func NewQueue(name string, bufferSize int) *Queue {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &Queue{
		ch:   make(chan Task, bufferSize),
		name: name,
	}
}

// Push 推入任务到队列，队列已关闭时返回错误
func (q *Queue) Push(task Task) error {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return errors.New("队列已关闭，无法推入任务")
	}
	q.mu.Unlock()

	select {
	case q.ch <- task:
		return nil
	default:
		return errors.New("队列已满，任务被拒绝")
	}
}

// Pop 返回任务消费通道
func (q *Queue) Pop() <-chan Task {
	return q.ch
}

// Close 关闭队列，已推入的任务仍可被消费完毕
func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if !q.closed {
		q.closed = true
		close(q.ch)
	}
}

// Len 返回当前队列中的任务数量
func (q *Queue) Len() int {
	return len(q.ch)
}

// Name 返回队列名称
func (q *Queue) Name() string {
	return q.name
}

// IsClosed 返回队列是否已关闭
func (q *Queue) IsClosed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.closed
}
