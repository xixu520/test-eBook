package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 文件分类
type Category struct {
	gorm.Model
	Name     string     `json:"name" gorm:"type:varchar(100);not null"`
	ParentID uint       `json:"parent_id" gorm:"index"`
	Order    int        `json:"order" gorm:"default:0"`
	Children []Category `json:"children" gorm:"foreignKey:ParentID"`
}

// StandardFile 建筑标准文件
type StandardFile struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string   `json:"title" gorm:"type:varchar(255);not null"`
	Number      string   `json:"number" gorm:"type:varchar(100);index"`   // 标准号，如 GB/T 50311
	Year        string   `json:"year" gorm:"type:varchar(10);index"`     // 年份
	Version     string   `json:"version" gorm:"type:varchar(50)"`        // 版本/修订号
	CategoryID  uint     `json:"category_id" gorm:"index"`
	Category    Category `json:"category" gorm:"references:ID"`
	FilePath    string   `json:"file_path" gorm:"type:varchar(500)"`
	FileSize    int64    `json:"file_size"`
	Status      int      `json:"status" gorm:"default:0"`                // 0: 未处理, 1: 已处理/OCR完毕
	OCRContent  string   `json:"ocr_content" gorm:"type:text"`           // 全文索引预留
	Tags        string   `json:"tags" gorm:"type:varchar(500)"`          // 逗号分隔的标签
}
