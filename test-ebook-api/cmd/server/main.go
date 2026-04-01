package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-ebook-api/internal/config"
	"test-ebook-api/internal/database"
	"test-ebook-api/internal/handler"
	"test-ebook-api/internal/pkg/ocr"
	"test-ebook-api/internal/pkg/queue"
	"test-ebook-api/internal/pkg/storage"
	"test-ebook-api/internal/pkg/worker"
	"test-ebook-api/internal/repository"
	"test-ebook-api/internal/router"
	"test-ebook-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 2. Initialize Database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. Auto Migrate
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// 4. Setup Storage Infrastructure
	storageEngine, err := storage.NewStorage(config.GlobalConfig.Storage)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	stagingPath := config.GlobalConfig.Storage.StagingPath
	if stagingPath == "" {
		stagingPath = "uploads/staging"
	}
	stagingStorage := storage.NewStagingStorage(stagingPath)

	// 5. Setup Queue & Repositories
	uploadQueue := queue.NewQueue("upload_sync", 100)
	ocrQueue := queue.NewQueue("ocr", 100)

	standardRepo := repository.NewStandardRepository(database.WriteDB)
	uploadTaskRepo := repository.NewUploadTaskRepository(database.WriteDB)

	paddleOCR := ocr.NewPaddleClient()

	// 6. Setup Services & Handlers
	standardService := service.NewStandardService(
		standardRepo, paddleOCR, storageEngine,
		stagingStorage, uploadTaskRepo, uploadQueue,
	)
	standardHandler := handler.NewStandardHandler(standardService)

	userRepo := repository.NewUserRepository(database.WriteDB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, userService, database.WriteDB)

	settingRepo := repository.NewSettingRepository(database.WriteDB)
	settingService := service.NewSettingService(settingRepo)

	auditRepo := repository.NewAuditRepository(database.WriteDB)
	auditService := service.NewAuditService(auditRepo)
	auditHandler := handler.NewAuditHandler(auditService)

	mockHandler := handler.NewMockHandler()
	systemHandler := handler.NewSystemHandler()

	formRepo := repository.NewFormRepository(database.WriteDB)
	formService := service.NewFormService(formRepo, standardRepo)
	formHandler := handler.NewFormHandler(formService)

	// 7. Setup Workers
	retryPolicy := worker.NewRetryPolicy(
		config.GlobalConfig.Storage.RetryMax,
		config.GlobalConfig.Storage.RetryInitialDelaySec,
	)

	syncWorker := worker.NewSyncWorker(worker.SyncWorkerConfig{
		UploadQueue: uploadQueue,
		OCRQueue:    ocrQueue,
		Staging:     stagingStorage,
		Remote:      storageEngine,
		Repo:        standardRepo,
		UploadRepo:  uploadTaskRepo,
		Concurrency: config.GlobalConfig.Storage.SyncConcurrency,
		RetryPolicy: retryPolicy,
		StorageType: config.GlobalConfig.Storage.Type,
	})

	orphanCleaner := worker.NewOrphanCleaner(worker.OrphanCleanerConfig{
		Staging:     stagingStorage,
		Remote:      storageEngine,
		Repo:        standardRepo,
		UploadRepo:  uploadTaskRepo,
		StorageType: config.GlobalConfig.Storage.Type,
		IntervalH:   config.GlobalConfig.Storage.OrphanCleanIntervalH,
		StaleHours:  config.GlobalConfig.Storage.OrphanStaleHours,
	})

	// Setting handler needs orphanCleaner reference
	settingHandler := handler.NewSettingHandler(settingService, orphanCleaner)

	// 8. Seed Initial Data
	if err := authService.SeedAdmin(); err != nil {
		log.Printf("Warning: Failed to seed admin user: %v", err)
	}

	// 9. Start Workers
	workerCtx, workerCancel := context.WithCancel(context.Background())
	syncWorker.Start(workerCtx)
	orphanCleaner.Start(workerCtx)

	// 10. Setup Router
	gin.SetMode(config.GlobalConfig.Server.Mode)
	r := router.InitRouter(
		authHandler,
		standardHandler,
		mockHandler,
		userHandler,
		settingHandler,
		auditHandler,
		systemHandler,
		formHandler,
		database.WriteDB,
	)

	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Stop accepting new tasks
	uploadQueue.Close()
	ocrQueue.Close()

	// Stop workers (waits for current tasks to complete)
	workerCancel()
	syncWorker.Stop()
	orphanCleaner.Stop()
	log.Println("Workers stopped")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close DB connections safely
	if sqlDB, err := database.WriteDB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing DB connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}

	log.Println("Server exiting")
}

