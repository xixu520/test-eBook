package router

import (
	"test-ebook-api/internal/handler"
	"test-ebook-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(
	authHandler *handler.AuthHandler,
	standardHandler *handler.StandardHandler,
	mockHandler *handler.MockHandler,
	userHandler *handler.UserHandler,
	settingHandler *handler.SettingHandler,
	auditHandler *handler.AuditHandler,
	systemHandler *handler.SystemHandler,
	formHandler *handler.FormHandler,
	db *gorm.DB,
) *gin.Engine {
	r := gin.Default()

	// Global Middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/register", authHandler.Register)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		protected.Use(middleware.AuditMiddleware(db)) // Record all state-changing ops
		{
			// Auth
			protected.GET("/auth/me", authHandler.Me)

			// User preferences
			protected.PUT("/users/me/theme", userHandler.UpdateTheme)

			// Category components
			categories := protected.Group("/categories")
			{
				categories.GET("", standardHandler.GetCategoryTree)
				categories.POST("", standardHandler.AddCategory)
				categories.PUT("/:id", standardHandler.UpdateCategory)
				categories.DELETE("/:id", standardHandler.DeleteCategory)
			}

			// Document components
			documents := protected.Group("/documents")
			{
				documents.GET("", standardHandler.ListFiles)
				documents.POST("/upload", standardHandler.UploadFile)
				documents.GET("/history", standardHandler.GetFileHistory)
				documents.GET("/:id", standardHandler.GetFileDetail)
				documents.PUT("/:id", standardHandler.UpdateFile)
				documents.GET("/:id/download", standardHandler.DownloadFile)
				documents.GET("/:id/preview", standardHandler.PreviewFile)
				documents.GET("/:id/fields", standardHandler.GetDocumentFields)
				documents.PUT("/:id/fields", standardHandler.SaveDocumentFields)
				documents.DELETE("/:id", standardHandler.DeleteFile)
				documents.POST("/:id/ocr/retry", standardHandler.RetryOCR)
				documents.POST("/:id/retry-sync", standardHandler.RetrySync)
			}

			// Tasks
			tasks := protected.Group("/tasks")
			{
				tasks.GET("/:task_id/status", standardHandler.GetTaskStatus)
			}
			protected.GET("/ocr/tasks", standardHandler.GetOCRTasks)

			// Recycle Bin
			recycle := protected.Group("/recycle-bin/documents")
			{
				recycle.GET("", standardHandler.GetRecycleBin)
				recycle.PUT("/restore", standardHandler.RestoreDocuments)
				recycle.POST("/batch-delete", standardHandler.BatchDeleteDocuments)
			}

			// Admin Management
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminGuard())
			{
				admin.GET("/dashboard", standardHandler.GetDashboardStats)
				admin.GET("/system/status", systemHandler.GetSystemStatus)
				admin.GET("/users", userHandler.GetUsers)
				admin.POST("/users", userHandler.CreateUser)
				admin.PUT("/users/:id/status", userHandler.UpdateStatus)
				admin.PUT("/users/:id", userHandler.UpdateUser)
				admin.PUT("/users/:id/password", userHandler.ResetPassword)
				admin.DELETE("/users/:id", userHandler.DeleteUser)
				admin.GET("/upload-tasks", standardHandler.GetUploadTasks)

				// Forms Management
				forms := admin.Group("/forms")
				{
					forms.GET("", formHandler.GetForms)
					forms.POST("", formHandler.CreateForm)
					forms.PUT("/:id", formHandler.UpdateForm)
					forms.DELETE("/:id", formHandler.DeleteForm)
					forms.POST("/:id/fields", formHandler.SaveFormFields)
				}
			}

			// Settings & Audit
			protected.GET("/settings", settingHandler.GetSettings)
			protected.PUT("/settings", settingHandler.SaveSettings)
			protected.POST("/settings/ocr-test", settingHandler.TestOCR)
			protected.POST("/settings/storage-test", settingHandler.TestStorage)
			protected.POST("/settings/orphan-scan", settingHandler.OrphanScan)
			protected.GET("/audit-logs", auditHandler.GetAuditLogs)

			// Announcements
			protected.GET("/announcements/active", mockHandler.GetActiveAnnouncements)
		}
	}

	return r
}
