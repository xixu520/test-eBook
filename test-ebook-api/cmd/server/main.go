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

	// 4. Setup Dependencies
	userRepo := repository.NewUserRepository(database.WriteDB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	standardRepo := repository.NewStandardRepository(database.WriteDB)
	paddleOCR := ocr.NewPaddleClient()
	standardService := service.NewStandardService(standardRepo, paddleOCR)
	standardHandler := handler.NewStandardHandler(standardService)
	mockHandler := handler.NewMockHandler()

	// 5. Seed Initial Data
	if err := authService.SeedAdmin(); err != nil {
		log.Printf("Warning: Failed to seed admin user: %v", err)
	}

	// 6. Setup Router
	gin.SetMode(config.GlobalConfig.Server.Mode)
	r := router.InitRouter(authHandler, standardHandler, mockHandler)

	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
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
