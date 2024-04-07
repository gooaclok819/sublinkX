package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitSqlite() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("./db/test.db"), &gorm.Config{})
	if err != nil {
		log.Println("连接数据库失败")
	}
	DB = db
	err = db.AutoMigrate(&User{}, &Subcription{}, &Node{}, &IPLogs{})
	if err != nil {
		log.Println("数据表迁移失败")
	}
	// 初始化用户数据
	err = db.First(&User{}).Error
	if err == gorm.ErrRecordNotFound {
		admin := &User{
			Username: "admin",
			Password: "123456",
			Role:     "admin",
			Nickname: "管理员",
		}
		err = admin.Create()
		if err != nil {
			log.Println("初始化添加用户数据失败")
		}
	}
	log.Println("数据库初始化成功")
}
