package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 给出响应，可自定义code和message
func Response(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, gin.H{
		"message": message,
		"data":    data,
	})
}

// Success 给出成功响应
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    data,
	})
}

// Fail 给出失败响应
func Fail(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, gin.H{
		"message": message,
		"data":    nil,
	})
}
