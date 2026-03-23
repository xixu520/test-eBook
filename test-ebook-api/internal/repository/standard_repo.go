package repository

import (
	"test-ebook-api/internal/model"

	"gorm.io/gorm"
)

type StandardRepository struct {
	db *gorm.DB
}

func NewStandardRepository(db *gorm.DB) *StandardRepository {
	return &StandardRepository{db: db}
}

// --- Category Operations ---

func (r *StandardRepository) CreateCategory(cat *model.Category) error {
	return r.db.Create(cat).Error
}

func (r *StandardRepository) GetCategoryTree() ([]model.Category, error) {
	var results []model.Category
	err := r.db.Preload("Children").Where("parent_id = 0").Order("\"order\" ASC").Find(&results).Error
	return results, err
}

func (r *StandardRepository) FindCategoryByID(id uint) (*model.Category, error) {
	var cat model.Category
	err := r.db.First(&cat, id).Error
	return &cat, err
}

// --- File Operations ---

func (r *StandardRepository) CreateFile(file *model.StandardFile) error {
	return r.db.Create(file).Error
}

func (r *StandardRepository) ListFiles(categoryID uint, year string, page, pageSize int) ([]model.StandardFile, int64, error) {
	var files []model.StandardFile
	var total int64

	db := r.db.Model(&model.StandardFile{})

	if categoryID > 0 {
		db = db.Where("category_id = ?", categoryID)
	}
	if year != "" {
		db = db.Where("year = ?", year)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = db.Preload("Category").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&files).Error
	return files, total, err
}

func (r *StandardRepository) FindFileByID(id uint) (*model.StandardFile, error) {
	var file model.StandardFile
	err := r.db.Preload("Category").First(&file, id).Error
	return &file, err
}

func (r *StandardRepository) UpdateFile(file *model.StandardFile) error {
	return r.db.Save(file).Error
}

func (r *StandardRepository) DeleteFile(id uint) error {
	return r.db.Delete(&model.StandardFile{}, id).Error
}
