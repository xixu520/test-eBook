package repository

import (
	"test-ebook-api/internal/model"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) GetLogs(page, pageSize int, action string) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	
	db := r.db.Model(&model.AuditLog{})
	if action != "" {
		db = db.Where("action = ?", action)
	}
	db.Count(&total)
	
	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

func (r *AuditRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}
