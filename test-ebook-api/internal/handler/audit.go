package handler

import (
	"net/http"
	"strconv"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	svc *service.AuditService
}

func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	logs, total, err := h.svc.GetLogs(page, pageSize)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	
	pkg.Success(c, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	})
}
