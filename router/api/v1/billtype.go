package v1

import (
	"fmt"
	"github.com/Samoy/bill_backend/middleware/jwt"
	"github.com/Samoy/bill_backend/models"
	"github.com/Samoy/bill_backend/router/api"
	"github.com/Samoy/bill_backend/service/billtypeservice"
	"github.com/Samoy/bill_backend/service/userservice"
	"github.com/Samoy/bill_backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"image"
	_ "image/jpeg"
	"net/http"
	"os"
	"path"
	"time"
)

var imageSavePath = "upload/images/"

func AddBillType(c *gin.Context) {
	user, err := userservice.GetUser(jwt.Username)
	billTypeName := c.Request.FormValue("name")
	if len(billTypeName) == 0 {
		api.Fail(c, http.StatusBadRequest, "账单类型名称不能为空")
		return
	}
	if err != nil {
		api.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	handleImageForm(c, func(imagePath string) {
		billType := models.BillType{
			Name:  billTypeName,
			Owner: user.ID,
			Image: imageSavePath,
		}
		err := billtypeservice.AddBillType(&billType)
		if err != nil {
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
	}
	_, err = billtypeservice.GetBillType(uint(com.StrTo(BillTypeID).MustUint8()))
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, "未获取到该账单类型")
	}
	billTypeName := c.Request.FormValue("name")
	if len(billTypeName) == 0 {
		api.Fail(c, http.StatusBadRequest, "账单类型名称不能为空")
		return
	}
	handleImageForm(c, func(imagePath string) {
		billType := models.BillType{
			Name:  billTypeName,
			Owner: user.ID,
			Image: imageSavePath,
		}
		err := billtypeservice.UpdateBillType(uint(com.StrTo(BillTypeID).MustUint8()), user.ID, &billType)
		if err != nil {
			api.Fail(c, http.StatusInternalServerError, "账单类型修改失败")
		} else {
			api.Success(c, "账单类型更新成功", billType)
		}
	})
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
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := userservice.GetUser(jwt.Username)
	if err != nil {
		api.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = billtypeservice.DeleteBillType(deleteBillTypeForm.BillTypeID, user.ID)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	api.Success(c, "账单类型删除成功", nil)
}

func handleImageForm(c *gin.Context, successHandle func(imagePath string)) {
	file, img, err := c.Request.FormFile("image")
	if err != nil {
		api.Fail(c, http.StatusBadRequest, "无效的image参数")
		return
	}
	if !utils.CheckImageSize(file, 100*1024) {
		api.Fail(c, http.StatusBadRequest, "图片大小不能超过100k")
		return
	}
	if !utils.CheckImageExt(img.Filename, ".png,.bmp,.jpg,.jpeg") {
		api.Fail(c, http.StatusBadRequest, "图片格式只能为png,bmp,jpg和jpeg")
		return
	}
	im, _, err := image.DecodeConfig(file)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, "图片解析失败")
		return
	}
	if im.Width > 128 || im.Height > 128 {
		api.Fail(c, http.StatusBadRequest, "图片尺寸应该为128*128以下")
		return
	}
	if err := utils.IsNotExistMkDir(imageSavePath); err != nil {
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
		return
	}
	tempParentPath := fmt.Sprintf("%stemp/", imageSavePath)
	if err := utils.IsNotExistMkDir(tempParentPath); err != nil {
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
		return
	}
	timestamp := time.Now().Unix()
	ext := path.Ext(img.Filename)
	tempPath := fmt.Sprintf("%s%d%s", tempParentPath, timestamp, ext)
	if err := c.SaveUploadedFile(img, tempPath); err != nil {
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
		return
	}
	filePath := fmt.Sprintf("%s%d%s", imageSavePath, timestamp, ext)
	if utils.CheckExist(tempPath) {
		if err := c.SaveUploadedFile(img, filePath); err != nil {
			api.Fail(c, http.StatusBadRequest, "图片上传失败")
			return
		} else {
			//处理上传成功
			successHandle(filePath)
		}
		_ = os.Remove(tempPath)
	} else {
		api.Fail(c, http.StatusBadRequest, "图片上传失败")
	}
}
