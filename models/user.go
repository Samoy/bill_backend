package models

type User struct {
	BaseModel
	Username  string `json:"username" gorm:"not null;unique_index"`
	Password  string `json:"-" gorm:"not null"`
	Telephone string `json:"telephone" gorm:"not null;unique_index"`
	Nickname  string `json:"nickname"`
}
