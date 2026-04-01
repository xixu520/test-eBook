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
	err := r.db.Create(log).Error
	if err != nil {
		return err
	}

	// 限制最多保留 1000 条审计记录，异步执行删除以防阻塞
	go func() {
		// 删除 id <= (SELECT id FROM audit_logs ORDER BY id DESC LIMIT 1 OFFSET 1000) 的记录
		// 这样确保永远只有最新的 1000 条保留
		subQuery := r.db.Model(&model.AuditLog{}).Select("id").Order("id DESC").Offset(1000).Limit(1)
		r.db.Where("id <= (?)", subQuery).Delete(&model.AuditLog{})
	}()

	return nil
}
