package models

import (
	"github.com/shopspring/decimal"
	"time"
)

// Bill 账单Model
type Bill struct {
	BaseModel
	Name       string          `json:"name" gorm:"not null"`                     //账单名称
	Amount     decimal.Decimal `json:"amount" gorm:"notnull;type:decimal(10,2)"` //账单金额
	BillTypeID uint            `json:"-" gorm:"not null"`                        //账单类型ID，外键
	BillType   *BillType       `json:"billType" gorm:"foreignKey:BillTypeID"`    //账单类型
	Date       time.Time       `json:"date" gorm:"not null"`                     //账单日期
	Remark     string          `json:"remark"`                                   //账单备注
	UserID     uint            `json:"user_id" gorm:"not null"`                  //用户ID
	Owner      User            `json:"-" gorm:"foreignKey:UserID"`               //用户
	Income     bool            `json:"income" gorm:"not null"`                   //是否是收入项目
}
