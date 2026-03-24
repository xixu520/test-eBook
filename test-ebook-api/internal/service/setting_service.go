package service

import (
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/repository"
)

type SettingService struct {
	repo *repository.SettingRepository
}

func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetSettings() ([]model.SystemSetting, error) {
	return s.repo.GetAll()
}

func (s *SettingService) SaveSettings(settings []model.SystemSetting) error {
	return s.repo.BatchSave(settings)
}

func (s *SettingService) TestOCRConnection(apiKey, secretKey string) error {
	// 实际调用 OCR Client 的测试方法
	// 目前简单模拟
	return nil
}
