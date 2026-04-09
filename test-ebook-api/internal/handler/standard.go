package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"test-ebook-api/internal/model"
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

// --- Recycle Bin ---

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
	publisher := c.PostForm("publisher")
	implDate := c.PostForm("implementation_date")
	implStatus := c.PostForm("implementation_status")
	dynamicFieldsStr := c.PostForm("dynamic_fields") // 获取动态字段 JSON

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

	// 解析动态字段
	dynamicFields := make(map[uint]string)
	if dynamicFieldsStr != "" {
		_ = json.Unmarshal([]byte(dynamicFieldsStr), &dynamicFields)
	}

	fileModel, taskID, err := h.svc.UploadFile(title, number, year, version, publisher, implDate, implStatus, uint(catID), dynamicFields, f, file.Filename, file.Size)

	pkg.Success(c, gin.H{
		"document": fileModel,
		"task_id":  taskID,
	})
}

func (h *StandardHandler) ListFiles(c *gin.Context) {
	catID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	year := c.Query("year")
	keyword := c.Query("keyword")
	publisher := c.Query("publisher")
	implStatus := c.Query("implementation_status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 解析动态属性过滤: filter[1]=value
	dynamicFilters := make(map[uint]string)
	queries := c.Request.URL.Query()
	for k, v := range queries {
		if len(v) > 0 && v[0] != "" {
			if len(k) > 7 && k[:7] == "filter[" && k[len(k)-1] == ']' {
				fieldID, _ := strconv.Atoi(k[7 : len(k)-1])
				if fieldID > 0 {
					dynamicFilters[uint(fieldID)] = v[0]
				}
			}
		}
	}

	files, total, err := h.svc.SearchFiles(uint(catID), year, keyword, publisher, implStatus, dynamicFilters, page, pageSize)
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
	file, fields, err := h.svc.GetFileDetailWithFields(uint(id))
	if err != nil {
		pkg.Error(c, http.StatusNotFound, 404, "文件不存在")
		return
	}
	pkg.Success(c, gin.H{
		"document":     file,
		"field_values": fields,
	})
}

func (h *StandardHandler) GetDocumentFields(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的文档ID")
		return
	}
	_, fv, err := h.svc.GetFileDetailWithFields(uint(id))
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, fv)
}

func (h *StandardHandler) SaveDocumentFields(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的文档ID")
		return
	}
	var input []struct {
		FieldID uint   `json:"field_id"`
		Value   string `json:"value"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	var dynFields []model.DocumentFieldValue
	for _, df := range input {
		dynFields = append(dynFields, model.DocumentFieldValue{
			DocumentID: uint(id),
			FieldID:    df.FieldID,
			Value:      df.Value,
		})
	}
	if err := h.svc.SaveDocumentFields(uint(id), dynFields); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *StandardHandler) UpdateFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的文档ID")
		return
	}

	var input struct {
		Title                string `json:"title" binding:"required"`
		Number               string `json:"number"`
		Version              string `json:"version"`
		Publisher            string `json:"publisher"`
		ImplementationDate   string `json:"implementation_date"`
		ImplementationStatus string `json:"implementation_status"`
		CategoryID           uint   `json:"category_id"`
		Status               int    `json:"status"`
		VerifyStatus         string `json:"verify_status"`
		DynamicFields        []struct {
			FieldID uint   `json:"field_id"`
			Value   string `json:"value"`
		} `json:"dynamic_fields"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	var dynFields []model.DocumentFieldValue
	for _, df := range input.DynamicFields {
		dynFields = append(dynFields, model.DocumentFieldValue{
			DocumentID: uint(id),
			FieldID:    df.FieldID,
			Value:      df.Value,
		})
	}

	if err := h.svc.UpdateFile(uint(id), input.Title, input.CategoryID, input.Status, input.VerifyStatus, dynFields); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
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

func (h *StandardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, stats)
}

func (h *StandardHandler) DownloadFile(c *gin.Context) {
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
	// Trigger file download
	stream, err := h.svc.GetFileStream(file.FilePath, file.SyncStatus)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "文件打开失败: "+err.Error())
		return
	}
	defer stream.Close()

	c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(file.Title))
	c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, stream)
}

func (h *StandardHandler) PreviewFile(c *gin.Context) {
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
	// Return file inline for preview (e.g. PDF)
	stream, err := h.svc.GetFileStream(file.FilePath, file.SyncStatus)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "文件打开失败: "+err.Error())
		return
	}
	defer stream.Close()

	c.Header("Content-Disposition", "inline; filename*=utf-8''"+url.PathEscape(file.Title))
	c.Header("Content-Type", "application/pdf")
	io.Copy(c.Writer, stream)
}

// --- Sync Task Handlers ---

func (h *StandardHandler) RetrySync(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, 400, "无效的文档ID")
		return
	}
	if err := h.svc.RetrySync(uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *StandardHandler) GetUploadTasks(c *gin.Context) {
	tasks, err := h.svc.GetUploadTasks(100)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	pkg.Success(c, tasks)
}
