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
	"test-ebook-api/internal/repository"
	"time"

	"github.com/google/uuid"
)

type StandardService struct {
	repo      *repository.StandardRepository
	ocrClient ocr.Client
}

func NewStandardService(repo *repository.StandardRepository, ocrClient ocr.Client) *StandardService {
	return &StandardService{
		repo:      repo,
		ocrClient: ocrClient,
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

// --- File Logic ---

func (s *StandardService) UploadFile(title, number, year, version string, categoryID uint, fileReader io.Reader, fileName string, fileSize int64) (*model.StandardFile, error) {
	// 1. Verify Category
	_, err := s.repo.FindCategoryByID(categoryID)
	if err != nil {
		return nil, errors.New("分类不存在")
	}

	// 2. Prepare Storage
	uploadDir := config.GlobalConfig.Upload.Dir
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	ext := filepath.Ext(fileName)
	newFileName := uuid.New().String() + ext
	savePath := filepath.Join(uploadDir, newFileName)

	out, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, fileReader)
	if err != nil {
		return nil, err
	}

	// 3. Save to DB
	standardFile := &model.StandardFile{
		Title:      title,
		Number:     number,
		Year:       year,
		Version:    version,
		CategoryID: categoryID,
		FilePath:   savePath,
		FileSize:   fileSize,
		Status:     0,
	}

	if err := s.repo.CreateFile(standardFile); err != nil {
		return nil, err
	}

	// 4. Async processing (OCR Real)
	go s.ProcessFile(standardFile.ID)

	return standardFile, nil
}

func (s *StandardService) ProcessFile(fileID uint) {
	file, err := s.repo.FindFileByID(fileID)
	if err != nil {
		return
	}

	log.Printf("[OCR] Starting process for file %d: %s", fileID, file.FilePath)

	// 1. Submit to OCR
	jobID, err := s.ocrClient.SubmitTask(file.FilePath)
	if err != nil {
		log.Printf("[OCR] Submission failed for file %d: %v", fileID, err)
		file.Status = 2 // Failed
		file.OCRContent = "OCR提交失败: " + err.Error()
		s.repo.UpdateFile(file)
		return
	}
	log.Printf("[OCR] Task submitted, JobID: %s", jobID)

	// 2. Poll for result
	maxRetries := config.GlobalConfig.OCR.TaskTimeoutMinutes * 6 // 10s interval
	for i := 0; i < maxRetries; i++ {
		time.Sleep(10 * time.Second)

		content, status, err := s.ocrClient.GetResult(jobID)
		log.Printf("[OCR] Polling JobID %s, attempt %d, status: %s, err: %v", jobID, i+1, status, err)
		if err != nil {
			file.Status = 2
			file.OCRContent = "OCR轮询失败: " + err.Error()
			s.repo.UpdateFile(file)
			return
		}

		if status == "success" {
			log.Printf("[OCR] Successfully processed JobID %s", jobID)
			file.Status = 1
			file.OCRContent = content
			s.repo.UpdateFile(file)
			return
		}

		if status == "failed" {
			log.Printf("[OCR] Processing failed for JobID %s", jobID)
			file.Status = 2
			file.OCRContent = "OCR处理失败"
			s.repo.UpdateFile(file)
			return
		}
		// status == "processing", continue
	}

	log.Printf("[OCR] Timeout reached for JobID %s", jobID)
	// Timeout
	file.Status = 2
	file.OCRContent = "OCR处理超时"
	s.repo.UpdateFile(file)
}

func (s *StandardService) SearchFiles(categoryID uint, year string, page, pageSize int) ([]model.StandardFile, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }
	return s.repo.ListFiles(categoryID, year, page, pageSize)
}

func (s *StandardService) GetFileDetail(id uint) (*model.StandardFile, error) {
	return s.repo.FindFileByID(id)
}

func (s *StandardService) DeleteFile(id uint) error {
	file, err := s.repo.FindFileByID(id)
	if err != nil {
		return err
	}
	// Delete physical file
	os.Remove(file.FilePath)
	return s.repo.DeleteFile(id)
}
