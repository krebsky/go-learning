package database

import (
	"blog/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	dsn := "root:Zhaoyang@100297@tcp(localhost:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	log.Println("Database connected successfully")

	// 删除旧表（如果存在）以避免外键约束冲突
	// 注意：这会删除所有数据，仅用于开发环境
	DB.Migrator().DropTable(&models.Comment{}, &models.Post{}, &models.User{})

	// 自动迁移模型
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	log.Println("Database migration completed")
}
