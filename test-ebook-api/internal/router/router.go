package router

import (
	"test-ebook-api/internal/handler"
	"test-ebook-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(authHandler *handler.AuthHandler, standardHandler *handler.StandardHandler, mockHandler *handler.MockHandler) *gin.Engine {
	r := gin.Default()

	// Global Middlewares
	// CORS config... (existing or handled by middleware)
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authHandler.Me)
		}

		// standards components...
		standards := v1.Group("/standards")
		{
			standards.GET("/categories", standardHandler.GetCategoryTree)
			standards.POST("/categories", standardHandler.AddCategory)
			standards.POST("/files", standardHandler.UploadFile)
			standards.GET("/files", standardHandler.ListFiles)
			standards.GET("/files/:id", standardHandler.GetFileDetail)
			standards.DELETE("/files/:id", standardHandler.DeleteFile)
		}

		// Mock routes for dashboard and announcements
		v1.GET("/stats/dashboard", mockHandler.GetDashboardStats)
		v1.GET("/announcements/active", mockHandler.GetActiveAnnouncements)
	}

	return r
}
