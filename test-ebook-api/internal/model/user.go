package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"not null;default:'user'" json:"role"` // admin / user
	Theme        string         `gorm:"default:'light'" json:"theme"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	Permissions  string         `json:"permissions"` // JSON string ["upload", "download"]
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
