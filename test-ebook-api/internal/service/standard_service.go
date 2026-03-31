package service

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"test-ebook-api/internal/config"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/pkg/ocr"
	"test-ebook-api/internal/pkg/storage"
	"test-ebook-api/internal/repository"
	"time"

	"github.com/google/uuid"
)

type StandardService struct {
	repo      *repository.StandardRepository
	ocrClient ocr.Client
	storage   storage.Storage
}

func NewStandardService(repo *repository.StandardRepository, ocrClient ocr.Client, storage storage.Storage) *StandardService {
	return &StandardService{
		repo:      repo,
		ocrClient: ocrClient,
		storage:   storage,
	}
}

// --- Category Logic ---

func (s *StandardService) GetCategoryTree() ([]model.Category, error) {
	return s.repo.GetCategoryTree()
}

func (s *StandardService) AddCategory(name string, parentID uint, order int) error {
	if parentID > 0 {
		_, err := s.repo.FindCategoryByID(parentID)
		if err != nil {
			return errors.New("父分类不存在")
		}
	}
	cat := &model.Category{
		Name:     name,
		ParentID: parentID,
		Order:    order,
	}
	return s.repo.CreateCategory(cat)
}

func (s *StandardService) UpdateCategory(id uint, name string, parentID uint, order int) error {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	if parentID > 0 {
		if parentID == id {
			return errors.New("父分类不能是自己")
		}
		currentParentID := parentID
		for currentParentID > 0 {
			if currentParentID == id {
				return errors.New("父分类不能是自己的子孙分类")
			}
			p, err := s.repo.FindCategoryByID(currentParentID)
			if err != nil {
				return errors.New("父分类不存在")
			}
			currentParentID = p.ParentID
		}
	}

	cat.Name = name
	cat.ParentID = parentID
	cat.Order = order

	return s.repo.UpdateCategory(cat)
}

func (s *StandardService) DeleteCategory(id uint) error {
	// 1. Check sub-categories
	subCount, err := s.repo.CountSubCategories(id)
	if err != nil {
		return err
	}
	if subCount > 0 {
		return errors.New("存在子分类，无法删除")
	}

	// 2. Check associated files
	fileCount, err := s.repo.CountFilesByCategory(id)
	if err != nil {
		return err
	}
	if fileCount > 0 {
		return errors.New("该分类下有关联文件，无法删除")
	}

	return s.repo.DeleteCategory(id)
}

// --- File Logic ---

func (s *StandardService) UploadFile(title, number, year, version, publisher, issueDate, implStatus string, categoryID uint, fileReader io.Reader, fileName string, fileSize int64) (*model.StandardFile, string, error) {
	// 1. Verify Category
	_, err := s.repo.FindCategoryByID(categoryID)
	if err != nil {
		return nil, "", errors.New("分类不存在")
	}

	ext := filepath.Ext(fileName)
	newFileName := uuid.New().String() + ext
	
	savePath, err := s.storage.Save(newFileName, fileReader)
	if err != nil {
		return nil, "", err
	}

	// 3. Save to DB
	standardFile := &model.StandardFile{
		Title:                title,
		Number:               number,
		Year:                 year,
		Version:              version,
		Publisher:            publisher,
		IssueDate:            issueDate,
		ImplementationStatus: implStatus,
		CategoryID:           categoryID,
		FilePath:             savePath,
		FileSize:             fileSize,
		Status:               0,
	}

	if err := s.repo.CreateFile(standardFile); err != nil {
		return nil, "", err
	}

	// 4. Create OCR Task
	taskID := uuid.New().String()
	task := &model.OCRTask{
		TaskID:     taskID,
		DocumentID: standardFile.ID,
		Status:     "pending",
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, "", err
	}

	// 5. Async processing (OCR Real)
	go s.ProcessFile(standardFile.ID, taskID)

	return standardFile, taskID, nil
}

func (s *StandardService) ProcessFile(fileID uint, taskID string) {
	file, err := s.repo.FindFileByID(fileID)
	if err != nil {
		return
	}

	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return
	}

	task.Status = "processing"
	task.Progress = 10
	s.repo.UpdateTask(task)

	log.Printf("[OCR] Starting process for file %d, task %s: %s", fileID, taskID, file.FilePath)

	localFilePath := file.FilePath

	// 如果当前存储非本地，需要先下载到临时文件
	if config.GlobalConfig.Storage.Type != "local" {
		reader, err := s.storage.Get(file.FilePath)
		if err != nil {
			log.Printf("[OCR] 无法从云存储获取文件 %s: %v", file.FilePath, err)
			file.Status = 2
			s.repo.UpdateFile(file)
			task.Status = "failed"
			task.Error = "无法从云存储获取文件: " + err.Error()
			s.repo.UpdateTask(task)
			return
		}

		tmpFile, err := os.CreateTemp("", "ocr-*"+filepath.Ext(file.FilePath))
		if err != nil {
			reader.Close()
			log.Printf("[OCR] 创建临时文件失败: %v", err)
			file.Status = 2
			s.repo.UpdateFile(file)
			task.Status = "failed"
			task.Error = "创建临时文件失败: " + err.Error()
			s.repo.UpdateTask(task)
			return
		}

		if _, err := io.Copy(tmpFile, reader); err != nil {
			reader.Close()
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			log.Printf("[OCR] 下载文件到临时目录失败: %v", err)
			file.Status = 2
			s.repo.UpdateFile(file)
			task.Status = "failed"
			task.Error = "下载文件到临时目录失败: " + err.Error()
			s.repo.UpdateTask(task)
			return
		}
		reader.Close()
		tmpFile.Close()
		localFilePath = tmpFile.Name()
		defer os.Remove(localFilePath) // 函数结束后清理临时文件
	}

	// 1. Submit to OCR
	jobID, err := s.ocrClient.SubmitTask(localFilePath)
	if err != nil {
		log.Printf("[OCR] Submission failed for file %d, task %s: %v", fileID, taskID, err)
		file.Status = 2 // Failed
		s.repo.UpdateFile(file)

		task.Status = "failed"
		task.Error = "OCR提交失败: " + err.Error()
		s.repo.UpdateTask(task)
		return
	}
	log.Printf("[OCR] Task submitted to provider, JobID: %s", jobID)
	task.Progress = 30
	s.repo.UpdateTask(task)

	// 2. Poll for result
	maxRetries := config.GlobalConfig.OCR.TaskTimeoutMinutes * 6 // 10s interval
	for i := 0; i < maxRetries; i++ {
		time.Sleep(10 * time.Second)

		content, status, err := s.ocrClient.GetResult(jobID)
		log.Printf("[OCR] Polling JobID %s for task %s, attempt %d, status: %s", jobID, taskID, i+1, status)
		
		task.Progress = 30 + (i+1)*2
		if task.Progress > 95 { task.Progress = 95 }
		s.repo.UpdateTask(task)

		if err != nil {
			file.Status = 2
			s.repo.UpdateFile(file)

			task.Status = "failed"
			task.Error = "OCR轮询失败: " + err.Error()
			s.repo.UpdateTask(task)
			return
		}

		if status == "success" {
			log.Printf("[OCR] Successfully processed JobID %s for task %s", jobID, taskID)
			file.Status = 1
			file.OCRContent = content
			s.repo.UpdateFile(file)

			task.Status = "completed"
			task.Progress = 100
			task.Result = content // In real case, we might parse specific fields
			s.repo.UpdateTask(task)
			return
		}

		if status == "failed" {
			log.Printf("[OCR] Processing failed for JobID %s for task %s", jobID, taskID)
			file.Status = 2
			s.repo.UpdateFile(file)

			task.Status = "failed"
			task.Error = "OCR处理失败"
			s.repo.UpdateTask(task)
			return
		}
	}

	log.Printf("[OCR] Timeout reached for JobID %s for task %s", jobID, taskID)
	file.Status = 2
	s.repo.UpdateFile(file)

	task.Status = "failed"
	task.Error = "OCR处理超时"
	s.repo.UpdateTask(task)
}

func (s *StandardService) SearchFiles(categoryID uint, year, keyword, publisher, implStatus string, page, pageSize int) ([]model.StandardFile, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }
	return s.repo.ListFiles(categoryID, year, keyword, publisher, implStatus, page, pageSize)
}

func (s *StandardService) UpdateFile(id uint, title, number, version, publisher, issueDate, implStatus string, categoryID uint) error {
	file, err := s.repo.FindFileByID(id)
	if err != nil {
		return errors.New("文件不存在")
	}

	if categoryID > 0 {
		_, err := s.repo.FindCategoryByID(categoryID)
		if err != nil {
			return errors.New("分类不存在")
		}
		file.CategoryID = categoryID
	} else {
		file.CategoryID = 0 // Allow unclassified
	}

	file.Title = title
	file.Number = number
	file.Version = version
	file.Publisher = publisher
	file.IssueDate = issueDate
	file.ImplementationStatus = implStatus
	file.Category = model.Category{} // Clear association so Gorm uses the updated CategoryID

	return s.repo.UpdateFile(file)
}

func (s *StandardService) GetFileDetail(id uint) (*model.StandardFile, error) {
	return s.repo.FindFileByID(id)
}

func (s *StandardService) DeleteFile(id uint) error {
	return s.repo.DeleteFile(id)
}

func (s *StandardService) GetTaskStatus(taskID string) (*model.OCRTask, error) {
	return s.repo.GetTaskByID(taskID)
}

func (s *StandardService) GetOCRTasks() ([]model.OCRTask, error) {
	return s.repo.GetTasks(100)
}

func (s *StandardService) RetryOCR(docID uint) (string, error) {
	file, err := s.repo.FindFileByID(docID)
	if err != nil {
		return "", errors.New("文件不存在")
	}

	// Create a new task
	taskID := uuid.New().String()
	task := &model.OCRTask{
		TaskID:     taskID,
		DocumentID: docID,
		Status:     "pending",
	}
	if err := s.repo.CreateTask(task); err != nil {
		return "", err
	}

	// Update file status to pending
	file.Status = 0
	s.repo.UpdateFile(file)

	// Restart process
	go s.ProcessFile(docID, taskID)

	return taskID, nil
}

func (s *StandardService) GetFileHistory(number string) ([]model.StandardFile, error) {
	if number == "" {
		return nil, errors.New("标准编号不能为空")
	}
	return s.repo.GetFileHistory(number)
}

func (s *StandardService) GetFileStream(path string) (io.ReadCloser, error) {
	return s.storage.Get(path)
}

func (s *StandardService) GetRecycleBin() ([]model.StandardFile, error) {
	return s.repo.GetRecycleBinFiles()
}

func (s *StandardService) RestoreDocuments(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.RestoreFiles(ids)
}

func (s *StandardService) HardDeleteDocuments(ids []uint, emptyAll bool) error {
	var toDeleteIDs []uint
	if emptyAll {
		files, err := s.repo.GetRecycleBinFiles()
		if err != nil {
			return err
		}
		for _, f := range files {
			toDeleteIDs = append(toDeleteIDs, f.ID)
		}
	} else {
		toDeleteIDs = ids
	}

	if len(toDeleteIDs) == 0 {
		return nil
	}

	// 1. Delete physical files
	// Fetch files including deleted ones to get their paths
	var files []model.StandardFile
	if err := s.repo.UnscopedFindFiles(toDeleteIDs, &files); err == nil {
		for _, f := range files {
			if f.FilePath != "" {
				s.storage.Delete(f.FilePath)
			}
		}
	}


	return s.repo.HardDeleteFiles(toDeleteIDs)
}

func (s *StandardService) GetDashboardStats() (map[string]interface{}, error) {
	totalFiles, _ := s.repo.TotalFilesCount()
	totalCats, _ := s.repo.TotalCategoriesCount()
	todayFiles, _ := s.repo.TodayUploadedCount()
	pendingOCR, _ := s.repo.PendingOCRCount()
	storageUsed, _ := s.repo.TotalStorageUsed()
	recentFiles, _ := s.repo.GetRecentFiles(5)

	recentActivities := make([]map[string]interface{}, 0)
	for _, f := range recentFiles {
		recentActivities = append(recentActivities, map[string]interface{}{
			"id":      f.ID,
			"content": "更新了文档: " + f.Title,
			"time":    f.UpdatedAt.Format("2006-01-02 15:04:05"),
			"type":    "upload",
		})
	}

	return map[string]interface{}{
		"total_documents":   totalFiles,
		"total_categories":  totalCats,
		"today_uploaded":    todayFiles,
		"pending_ocr":       pendingOCR,
		"storage_used":      storageUsed,
		"recent_activities": recentActivities,
	}, nil
}
