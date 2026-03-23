package pkg

import (
	"test-ebook-api/internal/model"

	"github.com/gin-gonic/gin"
)

func Result(c *gin.Context, httpCode int, code int, msg string, data interface{}) {
	c.JSON(httpCode, model.Response{
		Code:    code,
		Msg:     msg,
		Data:    data,
		Success: code == 200,
	})
}

func Success(c *gin.Context, data interface{}) {
	Result(c, 200, 200, "success", data)
}

func Error(c *gin.Context, httpCode int, code int, msg string) {
	Result(c, httpCode, code, msg, nil)
}
