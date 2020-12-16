package billservice

import (
	"fmt"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
	"github.com/Samoy/bill_backend/utils"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"time"
)

var timeFormat = "2006-01-02 15:04:05"

func AddBill(bill *models.Bill) error {
	return dao.DB.Create(&bill).Error
}

func GetBill(billID uint, userID uint) (models.Bill, error) {
	var bill models.Bill
	err := dao.DB.Preload("BillType").First(&bill, "id = ? and user_id = ?", billID, userID).Error
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
		db = db.Where("date <= ? and date >= ?", et, st)
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

func GetRecentBillList() ([]models.Bill, error) {
	recentSt, recentEt := utils.GetRecentRange()
	var billList []models.Bill
	err := dao.DB.Preload("BillType").Where("date >= ? and date <= ?", recentSt, recentEt).Order("date desc").Find(&billList).Error
	return billList, err
}

func GetBillOverview() (map[string]decimal.Decimal, error) {
	todaySt, todayEt := utils.GetToday()
	todayAmount, err := getBillAmount(todaySt, todayEt)
	if err != nil {
		return nil, err
	}
	weekSt, weekEt := utils.GetWeekRang()
	weekAmount, err := getBillAmount(weekSt, weekEt)
	if err != nil {
		return nil, err
	}
	monthSt, monthEt := utils.GetMonthRange()
	monthAmount, err := getBillAmount(monthSt, monthEt)
	if err != nil {
		return nil, err
	}
	annualSt, annualEt := utils.GetAnnualRange()
	annualAmount, err := getBillAmount(annualSt, annualEt)
	if err != nil {
		return nil, err
	}
	return map[string]decimal.Decimal{
		"today_amount":  todayAmount,
		"week_amount":   weekAmount,
		"month_amount":  monthAmount,
		"annual_amount": annualAmount,
	}, nil
}

func getBillAmount(st, et time.Time) (decimal.Decimal, error) {
	amount := 0.00
	count := 0
	rows, err := dao.DB.Model(&models.Bill{}).Select("sum(amount)").Where("date >= ? and date <= ?", st, et).Count(&count).Rows()
	if count <= 0 {
		return decimal.NewFromFloat(0), nil
	}

	if err != nil {
		logrus.Error(err)
		return decimal.NewFromFloat(0), err
	}
	if rows.Next() {
		err := rows.Scan(&amount)
		if err != nil {
			logrus.Error(err)
			return decimal.NewFromFloat(0), err
		}
	}
	return decimal.NewFromFloat(amount), nil
}

func DeleteBill(billID, userID uint) (err error) {
	var bill models.Bill
	err = dao.DB.Where("id = ? and user_id = ?", billID, userID).First(&bill).Error
	if err != nil {
		return
	}
	return dao.DB.Delete(&bill, billID).Error
}
