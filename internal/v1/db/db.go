package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(dsn string, models ...interface{}) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if len(models) > 0 {
		err = db.AutoMigrate(models...)
		if err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}
	return db
}
