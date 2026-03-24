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
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// 放回请求体供后续 Handler 使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		// 只有成功的操作才记录（可选，也可以记录失败的）
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			//TODO: 从上下文获取真实 UserID (目前简单处理)
			username := "admin" // 占位符，待 AuthMiddleware 完善
			
			audit := &model.AuditLog{
				Username: username,
				Action:   method + " " + c.Request.URL.Path,
				Details:  string(bodyBytes),
				IP:       c.ClientIP(),
			}
			
			if err := db.Create(audit).Error; err != nil {
				log.Printf("[Audit] Failed to save audit log: %v", err)
			}
		}
	}
}
