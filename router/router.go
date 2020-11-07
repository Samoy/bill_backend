package router

import (
	v1 "github.com/Samoy/bill_backend/router/api/v1"
	"github.com/Samoy/bill_backend/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// InitRouter 初始化路由 mode:构建模式,debug or release
func InitRouter(mode string) *gin.Engine {
	binding.Validator = new(validator.DefaultValidator)
	gin.SetMode(mode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	apiV1 := r.Group("/api/v1")
	{
		//注册
		apiV1.POST("/user/register", v1.Register)
		//登录
		apiV1.POST("/user/login", v1.Login)
	}

	return r
}
