package middleware

import (
	"net/http"
	"strings"
	"test-ebook-api/internal/config"
	"test-ebook-api/internal/pkg"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.Error(c, http.StatusUnauthorized, 401, "未登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			pkg.Error(c, http.StatusUnauthorized, 401, "Token 格式错误")
			c.Abort()
			return
		}

		claims, err := pkg.ParseToken(parts[1], config.GlobalConfig.JWT.Secret)
		if err != nil {
			pkg.Error(c, http.StatusUnauthorized, 401, "Token 无效或已过期")
			c.Abort()
			return
		}

		// Set user info to context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
