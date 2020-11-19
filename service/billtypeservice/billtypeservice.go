package billtypeservice

import (
	"errors"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
)

func AddBillType(billType *models.BillType) error {
	return dao.DB.Create(&billType).Error
}

func UpdateBillType(billTypeID, owner uint, data map[string]interface{}) (models.BillType, error) {
	var billType models.BillType
	err := dao.DB.Model(&billType).Where("id = ? and owner = ?", billTypeID, owner).Updates(data).First(&billType).Error
	return billType, err
}

func GetBillType(billTypeID uint) (models.BillType, error) {
	var billType models.BillType
	err := dao.DB.First(&billType, billTypeID).Error
	return billType, err
}

func DeleteBillType(billTypeID uint, ownerID uint) error {
	var billType models.BillType
	dao.DB.Where("id = ? and owner = ?", billTypeID, ownerID).First(&billType)
	if billType.ID < 0 {
		return errors.New("未找到该账单类型")
	}
	if dao.DB.Delete(&billType, billTypeID).Error != nil {
		return errors.New("账单类型删除失败")
	}
	return nil
}

func GetBillTypeList(ownerID uint) ([]models.BillType, error) {
	var billTypeList []models.BillType
	err := dao.DB.Where("owner = ? or owner is null", ownerID).Find(&billTypeList).Error
	return billTypeList, err
}
