package middleware

import (
	"bytes"
	"io"
	"log"
	"test-ebook-api/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditMiddleware 记录变更操作的审计日志
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 只记录变更操作
		if method != "POST" && method != "PUT" && method != "DELETE" {
			c.Next()
			return
		}

		// 读取请求体（为了后续记录详情）
		var bodyBytes []byte
		contentType := c.Request.Header.Get("Content-Type")
		
		// 遇到大文件上传或 Multipart，不要读取请求体以防止 OOM 风险
		if contentType != "" && bytes.HasPrefix([]byte(contentType), []byte("multipart/form-data")) {
			bodyBytes = []byte("文件上传请求或 Multipart 负载")
		} else if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// 放回请求体供后续 Handler 使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		// 只有成功的操作才记录
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			username, _ := c.Get("username")
			usernameStr, ok := username.(string)
			if !ok {
				usernameStr = "anonymous"
			}

			userIDVal, _ := c.Get("userID")
			userID, _ := userIDVal.(uint)
			
			action := mapActionName(method, c.Request.URL.Path)
			
			audit := &model.AuditLog{
				UserID:   userID,
				Username: usernameStr,
				Action:   action,
				Details:  string(bodyBytes),
				IP:       c.ClientIP(),
			}
			
			if err := db.Create(audit).Error; err != nil {
				log.Printf("[Audit] Failed to save audit log: %v", err)
			}
		}
	}
}

// mapActionName 将请求方法和路径映射为前端期待的简写标签 (LOGIN/UPLOAD/DELETE/EDIT/VERIFY)
func mapActionName(method, path string) string {
	pathBytes := []byte(path)
	
	if bytes.Contains(pathBytes, []byte("/auth/login")) {
		return "LOGIN"
	}
	if bytes.Contains(pathBytes, []byte("/documents/upload")) {
		return "UPLOAD"
	}
	if bytes.Contains(pathBytes, []byte("/documents") ) && method == "DELETE" {
		return "DELETE"
	}
	if (bytes.Contains(pathBytes, []byte("/documents")) || bytes.Contains(pathBytes, []byte("/categories"))) && method == "PUT" {
		return "EDIT"
	}
	if bytes.Contains(pathBytes, []byte("/verify")) {
		return "VERIFY"
	}
	
	return method + " " + path
}
