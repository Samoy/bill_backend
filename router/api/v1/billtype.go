package v1

import (
	"fmt"
	"github.com/Samoy/bill_backend/config"
	"github.com/Samoy/bill_backend/middleware/jwt"
	"github.com/Samoy/bill_backend/models"
	"github.com/Samoy/bill_backend/router/api"
	"github.com/Samoy/bill_backend/service/billtypeservice"
	"github.com/Samoy/bill_backend/service/userservice"
	"github.com/Samoy/bill_backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/unknwon/com"
	_ "golang.org/x/image/bmp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"
)

func AddBillType(c *gin.Context) {
	user, err := userservice.GetUser(jwt.Username)
	billTypeName := c.Request.FormValue("name")
	if len(billTypeName) == 0 {
		api.Fail(c, http.StatusBadRequest, "账单类型名称不能为空")
		return
	}
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusUnauthorized, "未找到该用户")
		return
	}
	file, img, err := c.Request.FormFile("image")
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "无效的image参数")
		return
	}
	handleImageForm(c, file, img, func(imagePath string) {
		billType := models.BillType{
			Name:  billTypeName,
			Owner: user.ID,
			Image: imagePath,
		}
		err := billtypeservice.AddBillType(&billType)
		if err != nil {
			logrus.Error(err)
			api.Fail(c, http.StatusInternalServerError, "添加账单类型失败")
		} else {
			api.Success(c, "账单类型添加成功", billType)
		}
	})
}

func GetBillType(c *gin.Context) {
	BillTypeID := c.Query("bill_type_id")
	if len(BillTypeID) == 0 {
		api.Fail(c, http.StatusBadRequest, "bill_type_id不能为空")
		return
	}
	billType, err := billtypeservice.GetBillType(uint(com.StrTo(BillTypeID).MustUint8()))
	if err != nil {
		logrus.Error(err)
		api.Fail(c, http.StatusInternalServerError, "未找到该账单类型")
	} else {
		api.Success(c, "获取账单类型成功", billType)
	}
}

func UpdateBillType(c *gin.Context) {
	user, err := userservice.GetUser(jwt.Username)
	if err != nil {
		api.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	BillTypeID := c.Request.FormValue("bill_type_id")
	if len(BillTypeID) == 0 {
		api.Fail(c, http.StatusBadRequest, "bill_type_id不能为空")
		return
	}
	_, err = billtypeservice.GetBillType(uint(com.StrTo(BillTypeID).MustUint8()))
	if err != nil {
		logrus.Error(err)
		api.Fail(c, http.StatusInternalServerError, "未获取到该账单类型")
		return
	}
	billTypeName := c.Request.FormValue("name")
	file, img, err := c.Request.FormFile("image")
	if (file != nil && img != nil) || len(billTypeName) != 0 {
		handleImageForm(c, file, img, func(imagePath string) {
			billTypeData := make(map[string]interface{})
			if billTypeName != "" {
				billTypeData["name"] = billTypeName
			}
			if imagePath != "" {
				billTypeData["image"] = imagePath
			}
			billType, err := billtypeservice.UpdateBillType(uint(com.StrTo(BillTypeID).MustUint8()), user.ID, billTypeData)
			if err != nil {
				logrus.Error(err)
				api.Fail(c, http.StatusInternalServerError, "账单类型修改失败")
			} else {
				api.Success(c, "账单类型更新成功", billType)
			}
		})
	} else {
		api.Fail(c, http.StatusBadRequest, "缺少参数")
	}
}

func GetBillTypeList(c *gin.Context) {
	user, err := userservice.GetUser(jwt.Username)
	if err != nil {
		api.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	billTypeList, err := billtypeservice.GetBillTypeList(user.ID)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, "获取账单类型列表失败")
		return
	} else {
		api.Success(c, "获取账单类型列表成功", billTypeList)
	}
}

type DeleteBillTypeForm struct {
	BillTypeID uint `json:"bill_type_id" binding:"required"`
}

func DeleteBillType(c *gin.Context) {
	deleteBillTypeForm := &DeleteBillTypeForm{}
	if err := c.ShouldBindJSON(deleteBillTypeForm); err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := userservice.GetUser(jwt.Username)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusUnauthorized, "未找到该用户")
		return
	}
	err = billtypeservice.DeleteBillType(deleteBillTypeForm.BillTypeID, user.ID)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "账单类型删除失败")
		return
	}
	api.Success(c, "账单类型删除成功", nil)
}

func handleImageForm(c *gin.Context, file multipart.File, img *multipart.FileHeader, resultHandle func(imagePath string)) {
	if file == nil || img == nil {
		resultHandle("")
		return
	}
	if !utils.CheckImageSize(file, 100*1024) {
		api.Fail(c, http.StatusBadRequest, "图片大小不能超过100k")
		return
	}
	if !utils.CheckImageExt(img.Filename, ".png", ".jpg", ".bmp", ".jpeg") {
		api.Fail(c, http.StatusBadRequest, "图片格式只能为png,bmp,jpg和jpeg")
		return
	}
	_, err := file.Seek(0, 0)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "图片上传失败")
		return
	}
	im, _, err := image.DecodeConfig(file)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "图片解析失败")
		return
	}
	if im.Width > 128 || im.Height > 128 {
		api.Fail(c, http.StatusBadRequest, "图片尺寸应该为128*128以下")
		return
	}
	if err := utils.IsNotExistMkDir(config.AppConf.ImageSavePath); err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
		return
	}
	tempParentPath := fmt.Sprintf("%stemp/", config.AppConf.ImageSavePath)
	if err := utils.IsNotExistMkDir(tempParentPath); err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
		return
	}
	timestamp := time.Now().Unix()
	ext := path.Ext(img.Filename)
	tempPath := fmt.Sprintf("%s%d%s", tempParentPath, timestamp, ext)
	if err := c.SaveUploadedFile(img, tempPath); err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
		return
	}
	filePath := fmt.Sprintf("%s%d%s", config.AppConf.ImageSavePath, timestamp, ext)
	if utils.CheckExist(tempPath) {
		if err := c.SaveUploadedFile(img, filePath); err != nil {
			logrus.Error(err.Error())
			api.Fail(c, http.StatusBadRequest, "图片上传失败")
			return
		}
		resultHandle(fmt.Sprintf("/images/%d%s", timestamp, ext))
		_ = os.Remove(tempPath)
	} else {
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
	}
}
