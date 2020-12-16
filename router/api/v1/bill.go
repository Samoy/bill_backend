package v1

import (
	"github.com/Samoy/bill_backend/middleware/jwt"
	"github.com/Samoy/bill_backend/models"
	"github.com/Samoy/bill_backend/router/api"
	"github.com/Samoy/bill_backend/service/billservice"
	"github.com/Samoy/bill_backend/service/userservice"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/unknwon/com"
	"net/http"
	"strconv"
	"time"
)

type BillBody struct {
	Name   string          `json:"name" binding:"required"`
	Amount decimal.Decimal `json:"amount" binding:"required,gt>0"`
	Remark string          `json:"remark" binding:"omitempty,max=100"`
	Date   string          `json:"date" binding:"required,datetime"`
	TypeID uint            `json:"type_id" binding:"required"`
	Income bool            `json:"income" binding:"required"`
}

var timeFormat = "2006-01-02 15:04:05"

func AddBill(c *gin.Context) {
	b := &BillBody{}
	if err := c.ShouldBindJSON(&b); err != nil {
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
	date, _ := time.Parse(timeFormat, b.Date)
	//FIXME:账单类型没有更新，一般用不到，但是确实存在问题。
	l := &models.Bill{
		Name:       b.Name,
		Amount:     b.Amount,
		Remark:     b.Remark,
		Date:       date,
		UserID:     user.ID,
		BillTypeID: b.TypeID,
		Income:     b.Income,
	}
	err = billservice.AddBill(l)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusBadRequest, "添加账单失败")
		return
	} else {
		api.Success(c, "添加账单成功", l)
	}
}

func GetBill(c *gin.Context) {
	billID := c.Query("bill_id")
	if len(billID) == 0 {
		api.Fail(c, http.StatusBadRequest, "bill_id不能为空")
		return
	}
	user, err := userservice.GetUser(jwt.Username)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusUnauthorized, "未找到该用户")
		return
	}
	bill, err := billservice.GetBill(uint(com.StrTo(billID).MustUint8()), user.ID)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "获取账单失败")
		return
	}
	api.Success(c, "获取账单成功", bill)
}

type UpdateBillBody struct {
	BillID uint            `json:"bill_id" binding:"required"`
	Name   string          `json:"name"`
	Amount decimal.Decimal `json:"amount" binding:"omitempty,gt>0"`
	Remark string          `json:"remark" binding:"omitempty,max=100"`
	Date   string          `json:"date" binding:"omitempty,datetime"`
	TypeID uint            `json:"type_id"`
	Income string          `json:"income"`
}

func UpdateBill(c *gin.Context) {
	ubb := &UpdateBillBody{}
	err := c.ShouldBindJSON(&ubb)
	if err != nil {
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
	billData := make(map[string]interface{})
	if ubb.Name != "" {
		billData["name"] = ubb.Name
	}
	if !ubb.Amount.IsZero() {
		billData["amount"] = ubb.Amount
	}
	if ubb.Remark != "" {
		billData["remark"] = ubb.Remark
	}
	if ubb.TypeID != 0 {
		billData["bill_type_id"] = ubb.TypeID
	}
	if ubb.Income != "" {
		income, _ := strconv.ParseBool(ubb.Income)
		billData["income"] = income
	}
	if ubb.Date != "" {
		date, _ := time.Parse(timeFormat, ubb.Date)
		billData["date"] = date
	}

	bill, err := billservice.UpdateBill(ubb.BillID, user.ID, billData)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "更新账单失败")
	} else {
		api.Success(c, "账单更新成功", bill)
	}
}

func GetBillList(c *gin.Context) {
	// 日期范围
	starTime := c.Query("start_time")
	endTime := c.Query("end_time")
	// 分页
	page := c.DefaultQuery("page", "0")
	pageSize := c.DefaultQuery("page_size", "10")
	// 账单类型
	category := c.Query("type")
	// 收入还是支出
	income := c.DefaultQuery("income", "0")
	// 排序(金额，时间远近)
	sortKey := c.DefaultQuery("sort_key", "date")
	asc := c.DefaultQuery("asc", "0")
	user, err := userservice.GetUser(jwt.Username)
	if err != nil {
		api.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	billList, err := billservice.GetBillList(user.ID, starTime, endTime,
		com.StrTo(page).MustInt(), com.StrTo(pageSize).MustInt(),
		uint(com.StrTo(category).MustUint8()), income, sortKey, asc,
	)
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, "查询账单失败")
	} else {
		api.Success(c, "查询账单列表成功", billList)
	}
}

func GetBillOverview(c *gin.Context) {
	billOverview, err := billservice.GetBillOverview()
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, "查询账单概览失败")
	} else {
		api.Success(c, "查询账单概览成功", billOverview)
	}
}

func GetRecentBill(c *gin.Context) {
	billList, err := billservice.GetRecentBillList()
	if err != nil {
		api.Fail(c, http.StatusInternalServerError, "查询最近7天账单失败")
	} else {
		api.Success(c, "查询最近7天账单成功", billList)
	}
}

type DeleteBillBody struct {
	BillID uint `json:"bill_id" binding:"required"`
}

func DeleteBill(c *gin.Context) {
	dbb := &DeleteBillBody{}
	err := c.ShouldBindJSON(dbb)
	if err != nil {
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
	err = billservice.DeleteBill(dbb.BillID, user.ID)
	if err != nil {
		logrus.Error(err.Error())
		api.Fail(c, http.StatusInternalServerError, "删除账单失败")
		return
	}
	api.Success(c, "删除账单成功", nil)
}
