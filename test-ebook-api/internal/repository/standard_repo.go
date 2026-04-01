package repository

import (
	"test-ebook-api/internal/model"
	"time"

	"gorm.io/gorm"
)

type StandardRepository struct {
	db *gorm.DB
}

func NewStandardRepository(db *gorm.DB) *StandardRepository {
	return &StandardRepository{db: db}
}

func (r *StandardRepository) GetDB() *gorm.DB {
	return r.db
}

// --- Category Operations ---

func (r *StandardRepository) CreateCategory(cat *model.Category) error {
	return r.db.Create(cat).Error
}

func (r *StandardRepository) UpdateCategory(cat *model.Category) error {
	return r.db.Save(cat).Error
}

func (r *StandardRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *StandardRepository) CountSubCategories(parentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

func (r *StandardRepository) CountFilesByCategory(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.StandardFile{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *StandardRepository) CountCategoriesByFormID(formID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("form_id = ?", formID).Count(&count).Error
	return count, err
}

func (r *StandardRepository) GetCategoryTree() ([]model.Category, error) {
	var results []model.Category
	err := r.db.Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).Where("parent_id = 0").Order("\"order\" ASC").Find(&results).Error

	if err != nil {
		return nil, err
	}

	var counts []struct {
		CategoryID uint
		Count      int
	}
	r.db.Model(&model.StandardFile{}).Select("category_id, count(id) as count").Group("category_id").Find(&counts)

	countMap := make(map[uint]int)
	for _, c := range counts {
		countMap[c.CategoryID] = c.Count
	}

	var assignCounts func(cats []model.Category)
	assignCounts = func(cats []model.Category) {
		for i := range cats {
			cats[i].DocCount = countMap[cats[i].ID]
			if len(cats[i].Children) > 0 {
				assignCounts(cats[i].Children)
			}
		}
	}
	assignCounts(results)

	return results, nil
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

func (r *StandardRepository) ListFiles(categoryID uint, year, keyword, publisher, implStatus string, dynamicFilters map[uint]string, page, pageSize int) ([]model.StandardFile, int64, error) {
	var files []model.StandardFile
	var total int64

	db := r.db.Model(&model.StandardFile{})

	if categoryID > 0 {
		db = db.Where("category_id = ?", categoryID)
	}
	if year != "" {
		db = db.Where("year = ?", year)
	}
	if keyword != "" {
		db = db.Where("title LIKE ? OR number LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if publisher != "" {
		db = db.Where("publisher = ?", publisher)
	}
	if implStatus != "" {
		db = db.Where("implementation_status = ?", implStatus)
	}

	// 动态属性过滤逻辑
	if len(dynamicFilters) > 0 {
		for fieldID, value := range dynamicFilters {
			if value == "" {
				continue
			}
			// 子查询寻找匹配该属性值的文档 ID
			subQuery := r.db.Model(&model.DocumentFieldValue{}).
				Select("document_id").
				Where("field_id = ?", fieldID).
				Where("value LIKE ?", "%"+value+"%")
			db = db.Where("id IN (?)", subQuery)
		}
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

func (r *StandardRepository) UpdateFileFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.StandardFile{}).Where("id = ?", id).Updates(fields).Error
}

func (r *StandardRepository) DeleteFile(id uint) error {
	return r.db.Delete(&model.StandardFile{}, id).Error
}

func (r *StandardRepository) GetFileHistory(number string) ([]model.StandardFile, error) {
	var files []model.StandardFile
	err := r.db.Where("number = ?", number).Order("created_at DESC").Find(&files).Error
	return files, err
}

func (r *StandardRepository) GetRecycleBinFiles() ([]model.StandardFile, error) {
	var files []model.StandardFile
	err := r.db.Unscoped().Where("deleted_at IS NOT NULL").Order("deleted_at DESC").Find(&files).Error
	return files, err
}

func (r *StandardRepository) RestoreFiles(ids []uint) error {
	return r.db.Unscoped().Model(&model.StandardFile{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}

func (r *StandardRepository) HardDeleteFiles(ids []uint) error {
	return r.db.Unscoped().Delete(&model.StandardFile{}, ids).Error
}

func (r *StandardRepository) UnscopedFindFiles(ids []uint, files *[]model.StandardFile) error {
	return r.db.Unscoped().Where("id IN ?", ids).Find(files).Error
}

// --- Dashboard Statistics ---

func (r *StandardRepository) TotalFilesCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.StandardFile{}).Count(&count).Error
	return count, err
}

func (r *StandardRepository) TotalCategoriesCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Count(&count).Error
	return count, err
}

func (r *StandardRepository) TodayUploadedCount() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.StandardFile{}).Where("created_at >= ?", today).Count(&count).Error
	return count, err
}

func (r *StandardRepository) PendingOCRCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.StandardFile{}).Where("status = ?", 0).Count(&count).Error
	return count, err
}

func (r *StandardRepository) TotalStorageUsed() (int64, error) {
	var total int64
	err := r.db.Model(&model.StandardFile{}).Select("COALESCE(SUM(file_size), 0)").Scan(&total).Error
	return total, err
}

func (r *StandardRepository) GetRecentFiles(limit int) ([]model.StandardFile, error) {
	var files []model.StandardFile
	err := r.db.Order("updated_at DESC").Limit(limit).Find(&files).Error
	return files, err
}

// --- OCR Task Operations ---

func (r *StandardRepository) CreateTask(task *model.OCRTask) error {
	return r.db.Create(task).Error
}

func (r *StandardRepository) GetTaskByID(taskID string) (*model.OCRTask, error) {
	var task model.OCRTask
	err := r.db.Where("task_id = ?", taskID).First(&task).Error
	return &task, err
}

func (r *StandardRepository) UpdateTask(task *model.OCRTask) error {
	return r.db.Save(task).Error
}

func (r *StandardRepository) GetTasks(limit int) ([]model.OCRTask, error) {
	var tasks []model.OCRTask
	err := r.db.Order("created_at DESC").Limit(limit).Find(&tasks).Error
	return tasks, err
}

// --- Orphan Cleaner Operations ---

// GetExpiredSoftDeletedFiles 获取超过指定时间的软删除文件
func (r *StandardRepository) GetExpiredSoftDeletedFiles(cutoff time.Time) ([]model.StandardFile, error) {
	var files []model.StandardFile
	err := r.db.Unscoped().Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoff).Find(&files).Error
	return files, err
}

// HardDeleteFileByID 按 ID 硬删除单个文件记录
func (r *StandardRepository) HardDeleteFileByID(id uint) error {
	return r.db.Unscoped().Delete(&model.StandardFile{}, id).Error
}

// GetAllFilePaths 获取所有活跃文档的文件路径（用于云端孤儿比对）
func (r *StandardRepository) GetAllFilePaths() ([]string, error) {
	var paths []string
	err := r.db.Model(&model.StandardFile{}).Pluck("file_path", &paths).Error
	return paths, err
}

