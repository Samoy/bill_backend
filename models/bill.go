package models

import "github.com/shopspring/decimal"

// Bill 账单Model
type Bill struct {
	BaseModel
	Name       string          `json:"name" gorm:"not null"`               //账单名称
	Amount     decimal.Decimal `json:"amount" gorm:"notnull;type:decimal"` //账单金额
	BillTypeID uint            `json:"bill_type_id" gorm:"not null"`       //账单类型ID，外键
	Type       BillType        `json:"-" gorm:"foreignKey:BillTypeID"`     //账单类型
	Remark     string          `json:"remark"`                             //账单备注
	UserID     string          `json:"user_id" gorm:"not null"`            //用户ID
	Owner      User            `json:"-" gorm:"foreignKey:UserID"`         //用户
}
