package handler

import (
	"test-ebook-api/internal/pkg"

	"github.com/gin-gonic/gin"
)

type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) GetDashboardStats(c *gin.Context) {
	pkg.Success(c, gin.H{
		"total_files":    0,
		"pending_ocr":    0,
		"categories":     0,
		"recent_updates": []interface{}{},
	})
}

func (h *MockHandler) GetActiveAnnouncements(c *gin.Context) {
	pkg.Success(c, []interface{}{})
}
