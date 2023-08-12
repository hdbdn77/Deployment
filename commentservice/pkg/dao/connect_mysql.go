package dao

import (
	"fmt"
	"github.com/micro/simplifiedTikTok/commentservice/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlConfig = configs.MySQL
	db          *gorm.DB
)

func init() {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DBName)
	con, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("无法连接数据库")
	}

	sqlDB, _ := con.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)

	db = con
}

func GetDB() *gorm.DB {
	return db
}
