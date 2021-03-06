package models

// BillType 账单类型Model
type BillType struct {
	BaseModel
	Name  string `json:"name" gorm:"not_null"`  //账单类型名称
	Image string `json:"image"`                 //账单类型图片
	Owner uint   `json:"owner" gorm:"not_null"` //账单类型拥有者，可为空，表示共有，此类账单类型不可删除
}
