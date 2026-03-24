package service

import (
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/repository"
)

type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) GetLogs(page, pageSize int) ([]model.AuditLog, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 20 }
	return s.repo.GetLogs(page, pageSize)
}
