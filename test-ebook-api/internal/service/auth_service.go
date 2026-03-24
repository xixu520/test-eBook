package service

import (
	"errors"
	"test-ebook-api/internal/config"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	if !user.IsActive {
		return "", nil, errors.New("账号已禁用，请联系管理员")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	cfg := config.GlobalConfig.JWT
	token, err := pkg.GenerateToken(user.ID, user.Username, user.Role, cfg.Secret, cfg.ExpireHours)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) GetUserInfo(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) SeedAdmin() error {
	_, err := s.userRepo.FindByUsername("admin")
	if err == nil {
		return nil
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	user := &model.User{
		Username:     "admin",
		PasswordHash: string(hash),
		Role:         "admin",
		Permissions:  `["upload", "download", "delete", "verify", "manage_category"]`,
	}
	return s.userRepo.Create(user)
}
