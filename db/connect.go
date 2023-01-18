package db

import (
	"fmt"
	"fun_server/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	url := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, url, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

}
