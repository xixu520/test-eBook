package service

import (
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/repository"
)

type DocumentFieldService struct {
	repo         *repository.DocumentFieldValueRepository
	formRepo     *repository.FormRepository
}

func NewDocumentFieldService(repo *repository.DocumentFieldValueRepository, formRepo *repository.FormRepository) *DocumentFieldService {
	return &DocumentFieldService{
		repo:     repo,
		formRepo: formRepo,
	}
}

func (s *DocumentFieldService) GetFieldValues(docID uint) ([]model.DocumentFieldValue, error) {
	return s.repo.GetByDocumentID(docID)
}

func (s *DocumentFieldService) SaveFieldValues(docID uint, values []model.DocumentFieldValue) error {
	return s.repo.BatchSave(docID, values)
}

func (s *DocumentFieldService) GetListColumns(formID uint) ([]model.FormField, error) {
	form, err := s.formRepo.FindFormByID(formID)
	if err != nil {
		return nil, err
	}
	var cols []model.FormField
	for _, f := range form.Fields {
		if f.ShowInList {
			cols = append(cols, f)
		}
	}
	return cols, nil
}

func (s *DocumentFieldService) GetFilterFields(formID uint) ([]model.FormField, error) {
	form, err := s.formRepo.FindFormByID(formID)
	if err != nil {
		return nil, err
	}
	var fields []model.FormField
	for _, f := range form.Fields {
		if f.ShowInFilter {
			fields = append(fields, f)
		}
	}
	return fields, nil
}
