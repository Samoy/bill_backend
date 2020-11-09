package billtypeservice

import (
	"errors"
	"github.com/Samoy/bill_backend/dao"
	"github.com/Samoy/bill_backend/models"
)

func AddBillType(billType *models.BillType) error {
	return dao.DB.Create(&billType).Error
}

func UpdateBillType(billID uint, billType *models.BillType) error {
	return dao.DB.Model(&billType).Where("id = ?", billID).Updates(models.BillType{
		Name:  billType.Name,
		Image: billType.Image,
	}).Error
}

func GetBillType(billID uint) (models.BillType, error) {
	var billType models.BillType
	err := dao.DB.First(&billType, billID).Error
	return billType, err
}

func DeleteBillType(billID uint, ownerID uint) error {
	var billType models.BillType
	dao.DB.Where("billID = ? and owner = ?", billID, ownerID).First(&billType)
	if billType.ID < 0 {
		return errors.New("未找到该账单类型")
	}
	if dao.DB.Delete(&billType, billID).Error != nil {
		return errors.New("账单类型删除失败")
	}
	return nil
}

func GetBillTypeList(ownerID uint) ([]models.BillType, error) {
	var billTypeList []models.BillType
	err := dao.DB.Where("owner = ? or owner = ''", ownerID).Find(&billTypeList).Error
	return billTypeList, err
}
