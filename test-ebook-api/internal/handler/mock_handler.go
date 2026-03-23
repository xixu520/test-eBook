package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) GetDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_files":    0,
			"pending_ocr":    0,
			"categories":     0,
			"recent_updates": []interface{}{},
		},
		"message": "success (mock)",
	})
}

func (h *MockHandler) GetActiveAnnouncements(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    []interface{}{},
		"message": "success (mock)",
	})
}
