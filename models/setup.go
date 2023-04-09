package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost password=postres dbname=postgres port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to the database")
	}

	err = database.AutoMigrate(&User{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&Group{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&UnsentMessage{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&Membership{})
	if err != nil {
		return
	}

	DB = database
}
