package handler

import (
	"net/http"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/pkg/worker"
	"test-ebook-api/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	svc           *service.SettingService
	orphanCleaner *worker.OrphanCleaner
}

func NewSettingHandler(svc *service.SettingService, orphanCleaner *worker.OrphanCleaner) *SettingHandler {
	return &SettingHandler{svc: svc, orphanCleaner: orphanCleaner}
}

func (h *SettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.svc.GetSettings()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, settings)
}

func (h *SettingHandler) SaveSettings(c *gin.Context) {
	var settings []model.SystemSetting
	if err := c.ShouldBindJSON(&settings); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.SaveSettings(settings); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *SettingHandler) TestOCR(c *gin.Context) {
	var req struct {
		Engine string                 `json:"engine"`
		Config map[string]interface{} `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.TestOCRConnection(req.Engine, req.Config); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "连接失败: "+err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *SettingHandler) TestStorage(c *gin.Context) {
	var req struct {
		Type   string                 `json:"type"`
		Config map[string]interface{} `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.TestStorageConnection(req.Type, req.Config); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "连接失败: "+err.Error())
		return
	}
	pkg.Success(c, nil)
}

// OrphanScan 手动触发孤儿文件扫描（管理员功能）
func (h *SettingHandler) OrphanScan(c *gin.Context) {
	if h.orphanCleaner == nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "孤儿清理器未初始化")
		return
	}
	result := h.orphanCleaner.RunManualScan()
	pkg.Success(c, result)
}
