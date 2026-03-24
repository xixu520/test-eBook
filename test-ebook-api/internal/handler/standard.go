package handler

import (
	"net/http"
	"strconv"
	"test-ebook-api/internal/pkg"
	"test-ebook-api/internal/service"

	"github.com/gin-gonic/gin"
)

type StandardHandler struct {
	svc *service.StandardService
}

func NewStandardHandler(svc *service.StandardService) *StandardHandler {
	return &StandardHandler{svc: svc}
}

// --- Category Handlers ---

func (h *StandardHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.svc.GetCategoryTree()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, tree)
}

func (h *StandardHandler) AddCategory(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		ParentID uint   `json:"parent_id"`
		Order    int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}

	if err := h.svc.AddCategory(input.Name, input.ParentID, input.Order); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *StandardHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的分类ID")
		return
	}
	var input struct {
		Name     string `json:"name" binding:"required"`
		ParentID uint   `json:"parent_id"`
		Order    int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}

	if err := h.svc.UpdateCategory(uint(id), input.Name, input.ParentID, input.Order); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *StandardHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的分类ID")
		return
	}
	if err := h.svc.DeleteCategory(uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

// --- History & Recycle Bin ---

func (h *StandardHandler) GetFileHistory(c *gin.Context) {
	number := c.Query("standard_no")
	history, err := h.svc.GetFileHistory(number)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, history)
}

func (h *StandardHandler) GetRecycleBin(c *gin.Context) {
	files, err := h.svc.GetRecycleBin()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, gin.H{
		"list": files,
	})
}

func (h *StandardHandler) RestoreDocuments(c *gin.Context) {
	var input struct {
		IDs []uint `json:"document_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.RestoreDocuments(input.IDs); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *StandardHandler) BatchDeleteDocuments(c *gin.Context) {
	var input struct {
		IDs      []uint `json:"document_ids"`
		EmptyAll bool   `json:"empty_all"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.HardDeleteDocuments(input.IDs, input.EmptyAll); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *StandardHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	task, err := h.svc.GetTaskStatus(taskID)
	if err != nil {
		pkg.Error(c, http.StatusNotFound, 404, "任务不存在")
		return
	}
	pkg.Success(c, task)
}

func (h *StandardHandler) GetOCRTasks(c *gin.Context) {
	tasks, err := h.svc.GetOCRTasks()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, tasks)
}

func (h *StandardHandler) RetryOCR(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的文档ID")
		return
	}
	taskID, err := h.svc.RetryOCR(uint(id))
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, gin.H{
		"task_id": taskID,
	})
}

// --- File Handlers ---

func (h *StandardHandler) UploadFile(c *gin.Context) {
	title := c.PostForm("title")
	number := c.PostForm("number")
	year := c.PostForm("year")
	version := c.PostForm("version")
	catIDStr := c.PostForm("category_id")

	if title == "" || catIDStr == "" {
		pkg.Error(c, http.StatusBadRequest, 400, "标题和分类不能为空")
		return
	}

	catID, _ := strconv.Atoi(catIDStr)
	file, err := c.FormFile("file")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "未上传文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "文件打开失败")
		return
	}
	defer f.Close()

	fileModel, taskID, err := h.svc.UploadFile(title, number, year, version, uint(catID), f, file.Filename, file.Size)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"document": fileModel,
		"task_id":  taskID,
	})
}

func (h *StandardHandler) ListFiles(c *gin.Context) {
	catID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	year := c.Query("year")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	files, total, err := h.svc.SearchFiles(uint(catID), year, page, pageSize)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"list":  files,
		"total": total,
		"page":  page,
	})
}

func (h *StandardHandler) GetFileDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的文档ID")
		return
	}
	file, err := h.svc.GetFileDetail(uint(id))
	if err != nil {
		pkg.Error(c, http.StatusNotFound, 404, "文件不存在")
		return
	}
	pkg.Success(c, file)
}

func (h *StandardHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的文档ID")
		return
	}
	if err := h.svc.DeleteFile(uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}
