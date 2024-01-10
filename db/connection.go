package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PRIVATE
var db *gorm.DB

func DataBase() *gorm.DB {
	return db
}

func Connection() {
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("DNS")), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error conecting to database %v", err)
	}
}
