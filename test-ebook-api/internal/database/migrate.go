package database

import (
	"test-ebook-api/internal/model"
)

func AutoMigrate() error {
	return WriteDB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.StandardFile{},
		&model.OCRTask{},
		&model.SystemSetting{},
		&model.AuditLog{},
		&model.UploadTask{},
		&model.Form{},
		&model.FormField{},
	)
}
