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
		&models.Profile{},
		&models.Batch{},
		&models.Block{},
	)
	fmt.Println("Database connnected")
}

// GetDB returns the database connection instance
func GetDB() *gorm.DB {
	return db
}
