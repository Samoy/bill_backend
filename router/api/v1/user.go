package v1

import (
	"github.com/Samoy/bill_backend/models"
	"github.com/Samoy/bill_backend/router/api"
	"github.com/Samoy/bill_backend/service/userservice"
	"github.com/Samoy/bill_backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// LoginBody 登录实体
type LoginBody struct {
	Username string `json:"username" binding:"required,min=6,max=16"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

//Login 用户登录
func Login(c *gin.Context) {
	l := &LoginBody{}
	if err := c.ShouldBindJSON(l); err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := userservice.Login(l.Username, l.Password)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	token, err := utils.GenerateToken(l.Username, l.Password)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "token生成失败")
		return
	}
	api.Success(c, "登录成功", map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// RegisterBody 注册实体
type RegisterBody struct {
	Username  string `json:"username" binding:"required,min=6,max=16"`
	Password  string `json:"password" binding:"required,min=6,max=16"`
	Telephone string `json:"telephone" binding:"required,telephone"`
	Nickname  string `json:"nickname" binding:"omitempty,min=2,max=10"`
}

// Register 注册
func Register(c *gin.Context) {
	r := &RegisterBody{}
	if err := c.ShouldBindJSON(r); err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	u := models.User{
		Username:  r.Username,
		Password:  r.Password,
		Telephone: r.Telephone,
		Nickname:  r.Nickname,
	}
	if err := userservice.Register(&u); err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "注册失败")
		return
	}
	api.Success(c, "注册成功", u)
}
