package router

import (
	"github.com/Samoy/bill_backend/config"
	"github.com/Samoy/bill_backend/middleware/jwt"
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
	r := gin.New()
	r.Use(gin.Recovery())
	if mode == "debug" {
		r.Use(gin.Logger())
	}
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	print(config.AppConf.ImageSavePath)
	r.StaticFS("/images/", http.Dir(config.AppConf.ImageSavePath))

	apiV1 := r.Group("/api/v1")
	{
		//注册
		apiV1.POST("/user/register", v1.Register)
		//登录
		apiV1.POST("/user/login", v1.Login)

		//以下路由需要token验证
		apiV1.Use(jwt.Jwt())
		{
			//获取用户信息
			apiV1.GET("/user/profile", v1.GetProfile)

			//账单类型
			apiV1.POST("/bill_type", v1.AddBillType)
			apiV1.GET("/bill_type", v1.GetBillType)
			apiV1.PUT("/bill_type", v1.UpdateBillType)
			apiV1.GET("/bill_type_list", v1.GetBillTypeList)
			apiV1.DELETE("/bill_type", v1.DeleteBillType)

			//账单
			apiV1.POST("/bill", v1.AddBill)
			apiV1.GET("/bill", v1.GetBill)
			apiV1.PUT("/bill", v1.UpdateBill)
			apiV1.GET("/bill_list", v1.GetBillList)
			apiV1.DELETE("/bill", v1.DeleteBill)
		}
	}

	return r
}
