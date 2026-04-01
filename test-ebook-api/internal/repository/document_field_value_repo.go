package repository

import (
	"test-ebook-api/internal/model"

	"gorm.io/gorm"
)

type DocumentFieldValueRepository struct {
	db *gorm.DB
}

func NewDocumentFieldValueRepository(db *gorm.DB) *DocumentFieldValueRepository {
	return &DocumentFieldValueRepository{db: db}
}

func (r *DocumentFieldValueRepository) GetByDocumentID(docID uint) ([]model.DocumentFieldValue, error) {
	var values []model.DocumentFieldValue
	err := r.db.Preload("Field").Where("document_id = ?", docID).Find(&values).Error
	return values, err
}

func (r *DocumentFieldValueRepository) BatchSave(docID uint, values []model.DocumentFieldValue) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除该文档的所有旧动态属性值
		if err := tx.Where("document_id = ?", docID).Delete(&model.DocumentFieldValue{}).Error; err != nil {
			return err
		}
		// 插入新值
		if len(values) > 0 {
			for i := range values {
				values[i].DocumentID = docID
				values[i].ID = 0 // 确保创建新记录
			}
			if err := tx.Create(&values).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *DocumentFieldValueRepository) DeleteByDocumentID(docID uint) error {
	return r.db.Where("document_id = ?", docID).Delete(&model.DocumentFieldValue{}).Error
}

// GetDocIDsByFieldValue 根据字段 ID 和过滤值查找匹配的文档 ID 列表
func (r *DocumentFieldValueRepository) GetDocIDsByFieldValue(fieldIDs []uint, filterValue string) ([]uint, error) {
	var docIDs []uint
	query := r.db.Model(&model.DocumentFieldValue{}).
		Where("value LIKE ?", "%"+filterValue+"%")
	
	if len(fieldIDs) > 0 {
		query = query.Where("field_id IN ?", fieldIDs)
	}
	
	err := query.Pluck("document_id", &docIDs).Error
	return docIDs, err
}
