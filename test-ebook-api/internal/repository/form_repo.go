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
		// 获取当前已有字段
		var existing []model.FormField
		if err := tx.Where("form_id = ?", formID).Find(&existing).Error; err != nil {
			return err
		}

		existingMap := make(map[uint]bool)
		for _, e := range existing {
			existingMap[e.ID] = true
		}

		// 收集本次提交中保留的 ID
		keepIDs := make(map[uint]bool)
		for i := range fields {
			fields[i].FormID = formID
			if fields[i].ID > 0 {
				keepIDs[fields[i].ID] = true
				// 更新已有字段
				if err := tx.Save(&fields[i]).Error; err != nil {
					return err
				}
			} else {
				// 新增字段
				if err := tx.Create(&fields[i]).Error; err != nil {
					return err
				}
			}
		}

		// 删除本次未提交的旧字段
		for id := range existingMap {
			if !keepIDs[id] {
				if err := tx.Delete(&model.FormField{}, id).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
