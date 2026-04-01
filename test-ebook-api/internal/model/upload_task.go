package model

import (
	"time"

	"gorm.io/gorm"
)

// UploadTask 文件上传同步任务（staging → 最终存储）
type UploadTask struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	DocumentID  uint       `json:"document_id" gorm:"index"`
	LocalPath   string     `json:"local_path" gorm:"type:varchar(500)"`               // staging 暂存路径
	RemotePath  string     `json:"remote_path" gorm:"type:varchar(500)"`              // 最终存储路径
	Status      string     `json:"status" gorm:"type:varchar(20);default:'pending'"`  // pending, uploading, synced, failed
	RetryCount  int        `json:"retry_count" gorm:"default:0"`
	MaxRetry    int        `json:"max_retry" gorm:"default:5"`
	Error       string     `json:"error" gorm:"type:text"`
	NextRetryAt *time.Time `json:"next_retry_at"` // 下次重试时间（指数退避）
}
