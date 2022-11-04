package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}
}

func Cleanup() {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get *sql.DB: %w", err))
	}

	if err = sqlDB.Close(); err != nil {
		panic(fmt.Errorf("failed to close database: %w", err))
	}
}
