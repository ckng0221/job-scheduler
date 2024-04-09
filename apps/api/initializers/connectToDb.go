package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func ConnectToDb() {
	env := os.Getenv("ENV")
	logLevel := logger.Silent
	if env == "development" {
		// logLevel = logger.Info
	}

	var err error
	dsn := os.Getenv("DB_URL")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		panic("Failed to connnect to DB")
	}
	fmt.Println("Connected to DB")
}
