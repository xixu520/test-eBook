package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemSetting 系统配置模型
type SystemSetting struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	Key         string         `json:"key" gorm:"type:varchar(100);uniqueIndex;not null"`
	Value       string         `json:"value" gorm:"type:text"`
	Description string         `json:"description" gorm:"type:varchar(255)"`
}
