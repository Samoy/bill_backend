package v1

import (
	"github.com/Samoy/bill_backend/models"
	"github.com/Samoy/bill_backend/router/api"
	"github.com/Samoy/bill_backend/service/userservice"
	"github.com/gin-gonic/gin"
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
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := userservice.Login(l.Username, l.Password); err != nil {
		api.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	api.Success(c, "登录成功", nil)
}

// RegisterBody 注册实体
type RegisterBody struct {
	Username  string `json:"username" binding:"required,min=6,max=16"`
	Password  string `json:"password" binding:"required,min=6,max=16"`
	Telephone string `json:"telephone" binding:"required,telephone"`
	Nickname  string `json:"nickname" binding:"omitempty,min=6,max=10"`
}

// Register 注册
func Register(c *gin.Context) {
	r := &RegisterBody{}
	if err := c.ShouldBindJSON(r); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	u := models.User{
		Username:  r.Username,
		Password:  r.Password,
		Telephone: r.Telephone,
		Nickname:  r.Nickname,
	}
	if err := userservice.Register(&u); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	api.Success(c, "注册成功", u)
}
