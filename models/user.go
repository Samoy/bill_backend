package models

type User struct {
	BaseModel
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"not null"`
	Telephone string `json:"telephone" gorm:"unique;not null"`
	Nickname  string `json:"nickname"`
}
