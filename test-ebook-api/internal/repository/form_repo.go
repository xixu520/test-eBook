package repository

import (
	"test-ebook-api/internal/model"

	"gorm.io/gorm"
)

type FormRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) *FormRepository {
	return &FormRepository{db: db}
}

func (r *FormRepository) CreateForm(form *model.Form) error {
	return r.db.Create(form).Error
}

func (r *FormRepository) GetFormsWithFields() ([]model.Form, error) {
	var forms []model.Form
	err := r.db.Preload("Fields", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).Find(&forms).Error
	return forms, err
}

func (r *FormRepository) FindFormByID(id uint) (*model.Form, error) {
	var form model.Form
	err := r.db.Preload("Fields", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).First(&form, id).Error
	return &form, err
}

func (r *FormRepository) UpdateForm(form *model.Form) error {
	return r.db.Save(form).Error
}

func (r *FormRepository) DeleteForm(id uint) error {
	return r.db.Delete(&model.Form{}, id).Error
}

// Field Operations

func (r *FormRepository) CreateFormField(field *model.FormField) error {
	return r.db.Create(field).Error
}

func (r *FormRepository) UpdateFormField(field *model.FormField) error {
	return r.db.Save(field).Error
}

func (r *FormRepository) DeleteFormField(id uint) error {
	return r.db.Delete(&model.FormField{}, id).Error
}

func (r *FormRepository) UpdateFormFields(formID uint, fields []model.FormField) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 简单的做法：先删除所有旧的，再插入新的（或者根据ID更新）
		// 考虑到可能涉及ID变化，这里采用删除+重新创建或者显式更新
		if err := tx.Where("form_id = ?", formID).Delete(&model.FormField{}).Error; err != nil {
			return err
		}
		for i := range fields {
			fields[i].FormID = formID
			fields[i].ID = 0 // 重置 ID 确保创建
		}
		if len(fields) > 0 {
			if err := tx.Create(&fields).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
