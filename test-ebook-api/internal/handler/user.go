package handler

import (
	"net/http"
	"strconv"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	users, total, err := h.svc.GetAllUsers(page, pageSize)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	
	pkg.Success(c, gin.H{
		"list":  users,
		"total": total,
		"page":  page,
	})
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.UpdateUserStatus(uint(id), req.IsActive); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteUser(uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *UserHandler) UpdateTheme(c *gin.Context) {
	// TODO: Get real user ID from context
	userID := uint(1) // Placeholder
	var req struct {
		Theme string `json:"theme"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.UpdateTheme(userID, req.Theme); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}
