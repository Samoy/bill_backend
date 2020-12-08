package billservice

import (
	"fmt"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
	"time"
)

var timeFormat = "2006-01-02 15:04:05"

func AddBill(bill *models.Bill) error {
	return dao.DB.Preload("BillType").Create(&bill).Error
}

func GetBill(billID uint, userID uint) (models.Bill, error) {
	var bill models.Bill
	err := dao.DB.First(&bill, "id = ? and user_id = ?", billID, userID).Error
	return bill, err
}

func UpdateBill(billID, userID uint, data map[string]interface{}) (models.Bill, error) {
	var bill models.Bill
	err := dao.DB.Model(&bill).Where("id = ? and user_id = ?", billID, userID).Updates(data).First(&bill).Error
	return bill, err
}

func GetBillList(
	userID uint,
	startTime string,
	endTime string,
	page int,
	pageSize int,
	category uint,
	income string,
	sortKey string,
	asc string,
) ([]models.Bill, error) {
	var (
		st, et      time.Time
		billList    []models.Bill
		incomeValue bool
		ascValue    string
	)
	db := dao.DB.Where("user_id = ?", userID)
	if income != "" {
		if income == "0" {
			incomeValue = false
		} else {
			incomeValue = true
		}
		db = db.Where("income = ?", incomeValue)
	}
	if startTime != "" && endTime != "" {
		st, _ = time.Parse(timeFormat, startTime)
		et, _ = time.Parse(timeFormat, endTime)
		db = db.Where("date < ? and date > ?", et, st)
	}
	if category != 0 {
		db = db.Where("bill_type_id = ?", category)
	}
	if sortKey != "" {
		if asc == "1" {
			ascValue = "asc"
		} else {
			ascValue = "desc"
		}
		db = db.Order(fmt.Sprintf("%s %s", sortKey, ascValue))
	}
	if page > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	err := db.Preload("BillType").Find(&billList).Error
	return billList, err
}

func DeleteBill(billID, userID uint) (err error) {
	var bill models.Bill
	err = dao.DB.Where("id = ? and user_id = ?", billID, userID).First(&bill).Error
	if err != nil {
		return
	}
	return dao.DB.Delete(&bill, billID).Error
}
