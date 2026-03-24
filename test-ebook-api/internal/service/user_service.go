package service

import (
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.GetAllUsers(page, pageSize)
}

func (s *UserService) UpdateUserStatus(id uint, isActive bool) error {
	return s.userRepo.UpdateUserStatus(id, isActive)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.DeleteUser(id)
}

func (s *UserService) UpdateTheme(id uint, theme string) error {
	return s.userRepo.UpdateUserTheme(id, theme)
}
