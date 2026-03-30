package service

import (
	"errors"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
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

// Register 用户自助注册
func (s *UserService) Register(username, password string) error {
	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         "user",
		IsActive:     true,
	}

	return s.userRepo.Create(user)
}

// CreateUser 管理员手动创建用户
func (s *UserService) CreateUser(username, password, role string) error {
	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
	}

	return s.userRepo.Create(user)
}

// UpdateUser 管理员修改用户信息 (如角色)
func (s *UserService) UpdateUser(id uint, role string) error {
	return s.userRepo.UpdateUserRole(id, role)
}

// ResetPassword 管理员强制重置用户密码
func (s *UserService) ResetPassword(id uint, newPassword string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	return s.userRepo.UpdateUserPassword(id, string(hash))
}

