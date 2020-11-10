package userservice

import (
	"errors"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
)

// Register 用户注册
// username 用户名
// password 密码
// telephone 手机号
// nickname 昵称
func Register(u *models.User) (err error) {
	err = existUser(u.Username, u.Password)
	if err == nil {
		err = dao.DB.Create(&u).Error
		if err != nil {
			err = errors.New("用户注册失败")
		}
	}
	return
}

// Login 用户登录
// username  用户名
// password 密码
func Login(username, password string) (user models.User, err error) {
	dao.DB.Where("username = ?", username).First(&user)
	if user.ID < 0 {
		err = errors.New("该用户不存在")
	}
	dao.DB.Where("username = ? and password = ?", username, password).First(&user)
	if user.ID < 0 {
		err = errors.New("用户名和密码不匹配")
	}
	return
}

func GetUser(username string) (models.User, error) {
	var user models.User
	err := dao.DB.First(&user, "username = ?", username).Error
	return user, err
}

func existUser(username, telephone string) (err error) {
	var user models.User
	dao.DB.Where("username = ?", username).First(&user)
	if user.ID > 0 {
		err = errors.New("该用户名已被注册")
	}
	dao.DB.Where("telephone = ?", telephone).First(&user)
	if user.ID > 0 {
		err = errors.New("该手机号已被注册")
	}
	return
}
