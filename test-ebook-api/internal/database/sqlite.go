package database

import (
	"log"
	"os"
	"path/filepath"
	"test-ebook-api/internal/config"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	WriteDB *gorm.DB
	ReadDB  *gorm.DB
)

func InitDB() error {
	cfg := config.GlobalConfig.Database
	dbDir := filepath.Dir(cfg.Path)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		_ = os.MkdirAll(dbDir, 0755)
	}

	// Write DB - Single connection for safety
	dsn := cfg.Path + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL"
	var err error
	WriteDB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, _ := WriteDB.DB()
	sqlDB.SetMaxOpenConns(1)

	// Explicitly set WAL mode and other pragmas for safety and performance
	sqlDB.Exec("PRAGMA journal_mode=WAL;")
	sqlDB.Exec("PRAGMA synchronous=NORMAL;")
	sqlDB.Exec("PRAGMA busy_timeout=5000;")

	// Read DB - Connection pool
	readDSN := "file:" + cfg.Path + "?mode=ro&_journal_mode=WAL"
	ReadDB, err = gorm.Open(sqlite.Open(readDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Warning: Failed to open read-only connection, fallback to WriteDB: %v", err)
		ReadDB = WriteDB
	} else {
		sqlDB2, _ := ReadDB.DB()
		sqlDB2.SetMaxOpenConns(cfg.MaxReadConns)
		sqlDB2.Exec("PRAGMA journal_mode=WAL;")
		sqlDB2.Exec("PRAGMA synchronous=NORMAL;")
		sqlDB2.Exec("PRAGMA busy_timeout=5000;")
	}

	return nil
}
