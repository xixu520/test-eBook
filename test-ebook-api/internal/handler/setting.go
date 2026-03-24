package handler

import (
	"net/http"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	svc *service.SettingService
}

func NewSettingHandler(svc *service.SettingService) *SettingHandler {
	return &SettingHandler{svc: svc}
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
		APIKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.TestOCRConnection(req.APIKey, req.SecretKey); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "连接失败: "+err.Error())
		return
	}
	pkg.Success(c, nil)
}
