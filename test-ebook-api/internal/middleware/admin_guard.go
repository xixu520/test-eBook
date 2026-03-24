package middleware

import (
	"net/http"
	"test-ebook-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

func AdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			pkg.Error(c, http.StatusForbidden, 403, "无权限操作")
			c.Abort()
			return
		}
		c.Next()
	}
}
