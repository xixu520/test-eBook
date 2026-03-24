package repository

import (
	"test-ebook-api/internal/model"

	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) GetByKey(key string) (*model.SystemSetting, error) {
	var setting model.SystemSetting
	err := r.db.Where("key = ?", key).First(&setting).Error
	return &setting, err
}

func (r *SettingRepository) GetAll() ([]model.SystemSetting, error) {
	var settings []model.SystemSetting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *SettingRepository) Save(setting *model.SystemSetting) error {
	return r.db.Save(setting).Error
}

func (r *SettingRepository) BatchSave(settings []model.SystemSetting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, s := range settings {
			var existing model.SystemSetting
			if err := tx.Where("key = ?", s.Key).First(&existing).Error; err == nil {
				existing.Value = s.Value
				existing.Description = s.Description
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
