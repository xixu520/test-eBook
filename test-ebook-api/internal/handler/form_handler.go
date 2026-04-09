package handler

import (
	"net/http"
	"strconv"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/service"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	svc *service.FormService
}

func NewFormHandler(svc *service.FormService) *FormHandler {
	return &FormHandler{svc: svc}
}

func (h *FormHandler) GetForms(c *gin.Context) {
	forms, err := h.svc.GetForms()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "加载表单失败")
		return
	}
	pkg.Success(c, forms)
}

func (h *FormHandler) GetGlobalForm(c *gin.Context) {
	form, err := h.svc.GetOrCreateGlobalForm()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "加载全局属性失败")
		return
	}
	pkg.Success(c, form)
}

func (h *FormHandler) CreateForm(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	form, err := h.svc.CreateForm(input.Name, input.Description)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, form)
}

func (h *FormHandler) UpdateForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的ID")
		return
	}
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.UpdateForm(uint(id), input.Name, input.Description); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *FormHandler) DeleteForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的ID")
		return
	}
	if err := h.svc.DeleteForm(uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *FormHandler) SaveFormFields(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Fields []model.FormField `json:"fields"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.SaveFormFields(uint(id), input.Fields); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *FormHandler) BindCategoriesToForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的表单ID")
		return
	}
	var input struct {
		CategoryIDs []uint `json:"category_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.BindCategoriesToForm(uint(id), input.CategoryIDs); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}
