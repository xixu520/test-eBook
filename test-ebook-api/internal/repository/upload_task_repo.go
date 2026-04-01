package repository

import (
	"test-ebook-api/internal/model"
	"time"

	"gorm.io/gorm"
)

type UploadTaskRepository struct {
	db *gorm.DB
}

func NewUploadTaskRepository(db *gorm.DB) *UploadTaskRepository {
	return &UploadTaskRepository{db: db}
}

func (r *UploadTaskRepository) Create(task *model.UploadTask) error {
	return r.db.Create(task).Error
}

func (r *UploadTaskRepository) GetByID(id uint) (*model.UploadTask, error) {
	var task model.UploadTask
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *UploadTaskRepository) Update(task *model.UploadTask) error {
	return r.db.Save(task).Error
}

// GetPendingTasks 获取所有待处理的任务（用于进程启动时恢复队列）
func (r *UploadTaskRepository) GetPendingTasks() ([]model.UploadTask, error) {
	var tasks []model.UploadTask
	err := r.db.Where("status IN ?", []string{"pending", "uploading"}).Find(&tasks).Error
	return tasks, err
}

// GetRetryableTasks 获取可以重试的失败任务（未超最大重试次数且到达重试时间）
func (r *UploadTaskRepository) GetRetryableTasks() ([]model.UploadTask, error) {
	var tasks []model.UploadTask
	err := r.db.Where("status = ? AND retry_count < max_retry AND (next_retry_at IS NULL OR next_retry_at <= ?)",
		"failed", time.Now()).Find(&tasks).Error
	return tasks, err
}

// GetFailedTasks 获取所有最终失败的任务（已达最大重试次数）
func (r *UploadTaskRepository) GetFailedTasks() ([]model.UploadTask, error) {
	var tasks []model.UploadTask
	err := r.db.Where("status = ? AND retry_count >= max_retry", "failed").Find(&tasks).Error
	return tasks, err
}

func (r *UploadTaskRepository) GetByDocumentID(docID uint) (*model.UploadTask, error) {
	var task model.UploadTask
	err := r.db.Where("document_id = ?", docID).Order("created_at DESC").First(&task).Error
	return &task, err
}

// GetAllTasks 获取所有上传任务（管理员查看，带分页）
func (r *UploadTaskRepository) GetAllTasks(limit int) ([]model.UploadTask, error) {
	var tasks []model.UploadTask
	err := r.db.Order("created_at DESC").Limit(limit).Find(&tasks).Error
	return tasks, err
}

// GetSyncedTasksWithLocalFiles 获取已同步但暂存文件仍然存在的任务（用于清理）
func (r *UploadTaskRepository) GetSyncedTasksWithLocalFiles() ([]model.UploadTask, error) {
	var tasks []model.UploadTask
	err := r.db.Where("status = ? AND local_path != ''", "synced").Find(&tasks).Error
	return tasks, err
}

// GetTaskByLocalPath 根据本地路径查找任务
func (r *UploadTaskRepository) GetTaskByLocalPath(localPath string) (*model.UploadTask, error) {
	var task model.UploadTask
	err := r.db.Where("local_path = ?", localPath).First(&task).Error
	return &task, err
}
