package model

import (
	"time"

	"gorm.io/gorm"
)

// OCRTask 异步 OCR 任务模型
type OCRTask struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	
	TaskID     string         `json:"task_id" gorm:"type:varchar(100);index;unique"`
	DocumentID uint           `json:"document_id" gorm:"index"`
	Status     string         `json:"status" gorm:"type:varchar(20);default:'pending'"` // pending, processing, completed, failed
	Progress   int            `json:"progress" gorm:"default:0"`
	Result     string         `json:"result" gorm:"type:text"` // JSON 字符串存储识别结果
	Error      string         `json:"error" gorm:"type:text"`
}
