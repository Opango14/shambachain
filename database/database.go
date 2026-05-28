package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"shambachain/models"
)

var (
	err error
	db  *gorm.DB
)

func InitDB() {
	db, err = gorm.Open(sqlite.Open("shambachain.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	db.AutoMigrate(
		&models.User{},
	)
	fmt.Println("Database connnected")
}
