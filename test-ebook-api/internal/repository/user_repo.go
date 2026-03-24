package repository

import (
	"test-ebook-api/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetAllUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	offset := (page - 1) * pageSize
	_ = r.db.Model(&model.User{}).Count(&total)
	err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) UpdateUserStatus(id uint, isActive bool) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *UserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) UpdateUserTheme(id uint, theme string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("theme", theme).Error
}
