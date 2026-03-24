package handler

import (
	"net/http"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.Error(c, http.StatusUnauthorized, 401, "未授权")
		return
	}

	user, err := h.authService.GetUserInfo(userID.(uint))
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	pkg.Success(c, user)
}
