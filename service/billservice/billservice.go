package billservice

import (
	"fmt"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
	"time"
)

var timeFormat = "2006-01-02 15:04:05"

func AddBill(bill *models.Bill) error {
	return dao.DB.Create(&bill).Error
}

func GetBill(billID uint, userID uint) (models.Bill, error) {
	var bill models.Bill
	err := dao.DB.First(&bill, "id = ? and user_id = ?", billID, userID).Error
	return bill, err
}

func UpdateBill(billID, userID uint, bill *models.Bill) error {
	return dao.DB.Model(&models.Bill{}).Where("id = ? and user_id = ?", billID, userID).Updates(&bill).Error
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
	if income == "1" {
		incomeValue = true
	} else {
		incomeValue = false
	}
	db := dao.DB.Where("income = ? and user_id = ?", incomeValue, userID)
	if startTime != "" && endTime != "" {
		st, _ = time.Parse(timeFormat, startTime)
		et, _ = time.Parse(timeFormat, endTime)
		db.Where("updated_at < ? and updated_at > ?", et, st)
	}
	if category != 0 {
		db.Where("bill_type_id = ?", category)
	}
	if sortKey != "" {
		if asc == "1" {
			ascValue = "asc"
		} else {
			ascValue = "desc"
		}
		db.Order(fmt.Sprintf("%s %s", sortKey, ascValue))
	}
	if page > 0 && pageSize > 0 {
		db.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	err := db.Find(&billList).Error
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
