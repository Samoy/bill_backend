package dao

import (
	"fmt"
	"github.com/Samoy/bill_backend/models"
	"github.com/jinzhu/gorm"
	"log"

	// MySQL连接器
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

// Setup 初始化数据库连接
func Setup(dbType, dbUser, dbPassword, dbHost, dbName, dbTablePrefix string) {
	var err error
	DB, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser,
		dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatalf("Failed to connect database:%v\n", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return dbTablePrefix + defaultTableName
	}
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.AutoMigrate(&models.BillType{}, &models.Bill{}, &models.User{})
}

func CloseDB() error {
	return DB.Close()
}
