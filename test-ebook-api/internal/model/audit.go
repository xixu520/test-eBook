package model

import (
	"time"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	
	UserID    uint      `json:"user_id" gorm:"index"`
	Username  string    `json:"username" gorm:"type:varchar(100)"`
	Action    string    `json:"action" gorm:"type:varchar(255)"` // METHOD + PATH
	Details   string    `json:"details" gorm:"type:text"`       // 请求参数或变更详情
	IP        string    `json:"ip" gorm:"type:varchar(50)"`
}
